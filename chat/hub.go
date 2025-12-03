package main

import (
	"encoding/json"
	"log"
	"server/common/config"
	"server/common/dto"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan dto.MessageResponseDto
	register   chan *Client
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan dto.MessageResponseDto),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("User %d connected", client.userID)

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}

		case msgDto := <-h.broadcast:
			payload, err := json.Marshal(msgDto)
			if err != nil {
				log.Println("Error marshaling message:", err)
				continue
			}

			// Optimization: Pre-fetch chat members if it's a chat message
			// This ensures we don't rely on stale client.chatIDs maps
			allowedUserIDs := make(map[uint]bool)

			if msgDto.ChatID != nil {
				var userIDs []uint
				// Query the join table directly to get current members
				if err := config.DB.Table("user_chats").
					Where("chat_id = ?", *msgDto.ChatID).
					Pluck("user_id", &userIDs).Error; err == nil {
					for _, uid := range userIDs {
						allowedUserIDs[uid] = true
					}
				}
			}

			for client := range h.clients {
				shouldSend := false

				// Check blocking
				if client.blockedUserIDs[msgDto.SenderID] {
					continue
				}

				// 1. Team Chat Broadcast
				if msgDto.TeamID != nil {
					if client.teamIDs[*msgDto.TeamID] {
						shouldSend = true
					}
				}

				// 2. Group/DM Chat Broadcast
				if msgDto.ChatID != nil {
					// Check local cache OR strict DB check results
					if client.chatIDs[*msgDto.ChatID] {
						shouldSend = true
					} else if allowedUserIDs[client.userID] {
						// Self-healing: User is in DB but not in cache. Update cache and send.
						client.chatIDs[*msgDto.ChatID] = true
						shouldSend = true
					}
				}

				if shouldSend {
					select {
					case client.send <- payload:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
		}
	}
}
