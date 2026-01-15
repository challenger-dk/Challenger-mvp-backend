package main

import (
	"encoding/json"
	"log"
	"server/common/dto"
	"server/common/services"
)

type Hub struct {
	clients map[*Client]bool

	// ✅ Now broadcast realtime events (message + typing)
	broadcast chan dto.RealtimeEventDto

	register   chan *Client
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan dto.RealtimeEventDto),
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

		case evt := <-h.broadcast:
			payload, err := json.Marshal(evt)
			if err != nil {
				log.Println("Error marshaling realtime event:", err)
				continue
			}

			// ✅ Pre-fetch participant IDs for conversation routing (once per event)
			var participantIDs map[uint]bool
			if evt.ConversationID != nil {
				ids, err := services.GetConversationParticipantIDs(*evt.ConversationID)
				if err != nil {
					log.Println("Error fetching conversation participants:", err)
					continue
				}
				participantIDs = make(map[uint]bool, len(ids))
				for _, id := range ids {
					participantIDs[id] = true
				}
			}

			for client := range h.clients {
				shouldSend := false

				// If the receiving client has blocked the triggering user, skip
				if client.blockedUserIDs[evt.UserID] {
					continue
				}

				// Conversation routing
				if evt.ConversationID != nil {
					if participantIDs != nil && participantIDs[client.userID] {
						shouldSend = true
					}
				}

				// Legacy team routing
				if evt.TeamID != nil {
					if _, isMember := client.teamIDs[*evt.TeamID]; isMember {
						shouldSend = true
					}
				}

				// Legacy DM routing
				if evt.RecipientID != nil {
					if client.userID == *evt.RecipientID || client.userID == evt.UserID {
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
