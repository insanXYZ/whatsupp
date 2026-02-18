// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"whatsupp-backend/dto"
	"whatsupp-backend/dto/converter"
	"whatsupp-backend/entity"
	"whatsupp-backend/util"

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

type HandlerIncomingMessage func(ctx context.Context, msg *dto.BroadcastMessageWs, hub *Hub) error

// Client is a middleman between the websocket connection and the Hub.
type Client struct {
	User *entity.User

	Hub *Hub

	// The websocket connection.
	Conn *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan *dto.BroadcastMessageWs

	// Handler for save message to db
	HandlerIncomingMessage HandlerIncomingMessage
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
			ctx := context.Background()

			event := new(dto.EventMessageWs)
			err = json.Unmarshal(message, event)
			if err != nil {
				fmt.Println("error unmarshaling message:", err.Error())
				continue
			}

			if event.Event != string(dto.EVENT_SEND_MESSAGE) {
				//handle invalid event
				fmt.Println("not send message event:", event.Event)
				continue
			}

			dataByte, err := json.Marshal(event.Data)
			if err != nil {
				// handle

				fmt.Println("error marshal:", err.Error())
				continue
			}

			var req dto.SendMessageRequestWs
			err = json.Unmarshal(dataByte, &req)
			if err != nil {
				fmt.Println("error umarshal:", err.Error())
				continue
			}

			broadcast := &dto.BroadcastMessageWs{
				Request: &req,
				Sender:  c.User,
			}

			if bc, err := util.MarshalIndent(broadcast); err == nil {
				fmt.Println(bc)
			}

			err = c.HandlerIncomingMessage(ctx, broadcast, c.Hub)
			if err != nil {
				fmt.Println("error handling incoming message:", err.Error())
				continue
			}

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
		case broadcast, ok := <-c.Send:
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

			var tmpConversationId *string

			if broadcast.Request.TmpConversationID != nil && c.User.ID == broadcast.Sender.ID {
				tmpConversationId = broadcast.Request.TmpConversationID
			}

			response := &dto.EventMessageWs{
				Event: string(dto.EVENT_NEW_MESSAGE),
				Data: &dto.NewMessageResponse{
					IsMe:              broadcast.Message.UserID == c.User.ID,
					ConversationID:    *broadcast.Request.ConversationID,
					TmpConversationID: tmpConversationId,
					Message:           converter.MessageEntityToDto(broadcast.Message),
				},
			}

			msgByte, err := json.Marshal(response)
			if err != nil {
				log.Println("error marshal message on writepump, ", err.Error())
				return
			}

			w.Write(msgByte)

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for range n {
				broadcast, ok := <-c.Send
				if !ok {
					continue
				}

				response := &dto.EventMessageWs{
					Event: string(dto.EVENT_NEW_MESSAGE),
					Data: &dto.NewMessageResponse{
						IsMe:              broadcast.Message.UserID == c.User.ID,
						ConversationID:    *broadcast.Request.ConversationID,
						TmpConversationID: broadcast.Request.TmpConversationID,
						Message:           converter.MessageEntityToDto(broadcast.Message),
					},
				}

				msgByte, err := json.Marshal(response)
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
