// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"encoding/json"
	"whatsupp-backend/dto"
	"whatsupp-backend/entity"

	"github.com/gorilla/websocket"
)

type SyncConversationFunc func(conversationId int) error

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	conversations map[int]map[int]bool

	// Registered clients.
	clients map[int]*Client

	// Inbound messages from the clients.
	broadcast chan *dto.BroadcastMessageWs

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	SyncConversation SyncConversationFunc
}

func NewHub() *Hub {
	return &Hub{
		conversations: make(map[int]map[int]bool),
		clients:       make(map[int]*Client),
		broadcast:     make(chan *dto.BroadcastMessageWs),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
	}
}

func (h *Hub) GetClient(clientId int) *entity.User {

	client, ok := h.clients[clientId]
	if !ok {
		return nil
	}

	return client.User
}

func (h *Hub) DeleteClientConversation(conversationId, clientId int) {
	clients, ok := h.conversations[conversationId]
	if ok {
		delete(clients, clientId)
	}
}

func (h *Hub) UpdateClient(id int, user *entity.User) {
	clients, ok := h.clients[id]
	if ok {
		clients.User = user
	}
}

func (h *Hub) sendMessageClient(clientId int, eventName string, data any) error {

	client, ok := h.clients[clientId]
	if !ok {
		return nil
	}

	event := &dto.EventMessageWs{
		Event: eventName,
		Data:  data,
	}

	dataByte, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return client.Conn.WriteMessage(websocket.TextMessage, dataByte)
}

func (h *Hub) sendMessageConversations(conversationId int, eventName string, data any) error {

	clients, ok := h.conversations[conversationId]
	if !ok {
		return nil
	}

	event := &dto.EventMessageWs{
		Event: eventName,
		Data:  data,
	}

	dataByte, err := json.Marshal(event)
	if err != nil {
		return err
	}

	for clientId, ok := range clients {
		if ok {
			client, ok := h.clients[clientId]
			if ok {
				client.Conn.WriteMessage(websocket.TextMessage, dataByte)
			}
		}
	}

	return nil
}

func (h *Hub) SendLeaveConversation(clientId int, data *dto.LeaveConversationResponse) error {
	return h.sendMessageClient(clientId, string(dto.EVENT_LEAVE_CONVERSATION), data)
}

func (h *Hub) SendMemberLeaveConversation(data *dto.MemberLeaveConversationResponse) error {
	return h.sendMessageConversations(data.ConversationId, string(dto.EVENT_MEMBER_LEAVE_CONVERSATION), data)
}

func (h *Hub) SendMemberJoinConversation(conversationId int, data *dto.MemberJoinConversationResponse) error {
	return h.sendMessageConversations(conversationId, string(dto.EVENT_MEMBER_JOIN_CONVERSATION), data)
}

func (h *Hub) SendNewConversation(clientId int, data *dto.NewConversationResponse) error {
	return h.sendMessageClient(clientId, string(dto.EVENT_NEW_CONVERSATION), data)
}

func (h *Hub) CreateConversation(conversationID int, members []int) {
	_, exist := h.conversations[conversationID]

	if exist {
		return
	}

	membersMap := make(map[int]bool)
	for _, v := range members {
		membersMap[v] = true
	}

	h.conversations[conversationID] = membersMap
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) IsExistConversation(conversationID int) bool {
	_, exist := h.conversations[conversationID]
	return exist
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.User.ID] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client.User.ID]; ok {
				delete(h.clients, client.User.ID)
				close(client.Send)
			}
		case message := <-h.broadcast:

			clients, ok := h.conversations[*message.Request.ConversationID]
			if !ok {
				continue
			}

			isMember := clients[message.Sender.ID]

			if !isMember {
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
					close(client.Send)
					delete(h.clients, clientID)
				}
			}

		}
	}
}
