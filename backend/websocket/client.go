// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"whatsupp-backend/dto"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type HandlerSendMessage func(msg *dto.BroadcastMessageWS) (int, error)

// Client is a middleman between the websocket connection and the Hub.
type Client struct {
	Id int

	Hub *Hub

	// The websocket connection.
	Conn *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan *dto.BroadcastMessageWS

	HandlerSendMessage HandlerSendMessage
}

// readPump pumps messages from the websocket connection to the Hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		mt, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		switch mt {
		case websocket.TextMessage:
			req := new(dto.MessageWS)
			err = json.Unmarshal(message, req)
			if err != nil {
				log.Printf("error unmarshaling message: %v", err)
				break
			}

			broadcast := &dto.BroadcastMessageWS{
				MessageWS: req,
				ClientID:  c.Id,
			}

			messageId, err := c.HandlerSendMessage(broadcast)
			if err != nil {
				log.Printf("error send message message: %v", err)
				break
			}

			broadcast.MessageID = messageId

			c.Hub.broadcast <- broadcast
		default:
			log.Println("unsupported message type:", mt)
			continue
		}

	}
}

// writePump pumps messages from the Hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The Hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			msgByte, err := json.Marshal(message)
			if err != nil {
				log.Println("error marshal message on writepump, ", err.Error())
				return
			}

			w.Write(msgByte)

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for range n {
				msgByte, err := json.Marshal(<-c.Send)
				if err != nil {
					continue
				}

				w.Write(newline)
				w.Write(msgByte)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
