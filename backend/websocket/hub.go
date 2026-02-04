// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"whatsupp-backend/dto"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	groups map[int]map[int]bool

	// Registered clients.
	clients map[int]*Client

	// Inbound messages from the clients.
	broadcast chan *dto.BroadcastMessageWS

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		groups:     make(map[int]map[int]bool),
		clients:    make(map[int]*Client),
		broadcast:  make(chan *dto.BroadcastMessageWS),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.Id] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client.Id]; ok {
				delete(h.clients, client.Id)
				close(client.Send)
			}
		case message := <-h.broadcast:

			if message.Event == dto.SEND_MESSAGE {

				data, ok := message.Data.(dto.SendMessageWS)
				if !ok {
					// handle unavailable schema
				}

				clients, ok := h.groups[data.GroupID]
				if !ok {
					// handle missing group id
				}

				isMember := clients[message.ClientID]

				if !isMember {
					// handle is not member group
				}

				for clientID, _ := range clients {
					client, ok := h.clients[clientID]
					if ok {
						client.Send <- 
					}
				}

			}


			//
			// for client := range h.clients {
			// 	select {
			// 	case client.Send <- message:
			// 	default:
			// 		close(client.Send)
			// 		delete(h.clients, client)
			// 	}
			// }
		}
	}
}
