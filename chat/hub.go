package main

import (
	"encoding/json"
	"log"
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
					if client.chatIDs[*msgDto.ChatID] {
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
