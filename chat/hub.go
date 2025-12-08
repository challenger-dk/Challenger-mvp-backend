package main

import (
	"encoding/json"
	"log"
	"server/common/dto"
)

type Hub struct {
	clients map[*Client]bool

	// Changed to handle DTOs
	broadcast chan dto.MessageResponseDto

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
				// If the recipient (client) has blocked the sender, they should NOT receive the message
				if client.blockedUserIDs[msgDto.SenderID] {
					continue
				}

				if msgDto.TeamID != nil {
					if _, isMember := client.teamIDs[*msgDto.TeamID]; isMember {
						shouldSend = true
					}
				}

				if msgDto.RecipientID != nil {
					if client.userID == *msgDto.RecipientID || client.userID == msgDto.SenderID {
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
