// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"fmt"
	"log"
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

func (h *Hub) CreateGroup(groupId int, members []int) {
	membersMap := make(map[int]bool)
	for _, v := range members {
		membersMap[v] = true
	}

	h.groups[groupId] = membersMap
}

func (h *Hub) Register(client *Client) {
	h.register <- client
	fmt.Println("success register client with id:", client.Id)
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

			clients, ok := h.groups[*message.GroupID]
			if !ok {
				log.Println("error missing group id")
				continue
			}

			isMember := clients[message.ClientID]

			if !isMember {
				log.Println("error is not member")
				continue
			}

			for clientID := range clients {
				client, ok := h.clients[clientID]
				if !ok {
					continue
				}

				select {
				case client.Send <- message:
				default:
					log.Println("dropping message, client too slow:", clientID)

					close(client.Send)
					delete(h.clients, clientID)
				}
			}

		}
	}
}
