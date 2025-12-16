package main

import (
	"encoding/json"
	"log"
	"server/common/config"
	"server/common/dto"
	"server/common/services"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	hub            *Hub
	conn           *websocket.Conn
	send           chan []byte
	userID         uint
	teamIDs        map[uint]bool
	blockedUserIDs map[uint]bool
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)

	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Printf("Error setting initial read deadline: %v", err)
		return
	}

	c.conn.SetPongHandler(func(string) error {
		if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			log.Printf("Error setting read deadline in PongHandler: %v", err)
			return err
		}
		return nil
	})

	for {
		_, messageBytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var req dto.IncomingMessage
		if err := json.Unmarshal(messageBytes, &req); err != nil {
			log.Printf("Invalid JSON: %v", err)
			continue
		}

		if req.Content == "" {
			continue
		}

		// New conversation-based messaging
		if req.ConversationID != nil {
			// Use conversation service
			msg, err := services.SendMessage(*req.ConversationID, c.userID, req.Content)
			if err != nil {
				log.Printf("Error sending message via conversation: %v", err)
				continue
			}
			c.hub.broadcast <- dto.ToMessageResponseDto(*msg)
		} else if req.TeamID != nil || req.RecipientID != nil {
			// Legacy messaging (backward compatibility)
			// --- Blocking Check (Incoming DM) ---
			// Prevent user from sending a DM to someone who has blocked them
			if req.RecipientID != nil {
				if services.IsBlocked(*req.RecipientID, c.userID) {
					// User is blocked by recipient. Ignore message.
					continue
				}
			}

			// Direct DB write for legacy messages
			if err := config.DB.Exec(
				"INSERT INTO messages (sender_id, team_id, recipient_id, content, created_at) VALUES (?, ?, ?, ?, ?)",
				c.userID, req.TeamID, req.RecipientID, req.Content, time.Now(),
			).Error; err != nil {
				log.Printf("Error saving message: %v", err)
				continue
			}

			// Fetch the message back with sender info
			var msg dto.MessageResponseDto
			config.DB.Raw(`
				SELECT m.id, m.sender_id, m.team_id, m.recipient_id, m.content, m.created_at,
					   u.id as "sender.id", u.email as "sender.email", u.first_name as "sender.first_name",
					   u.last_name as "sender.last_name", u.profile_picture as "sender.profile_picture"
				FROM messages m
				JOIN users u ON m.sender_id = u.id
				WHERE m.sender_id = ? AND m.created_at >= ?
				ORDER BY m.created_at DESC
				LIMIT 1
			`, c.userID, time.Now().Add(-1*time.Second)).Scan(&msg)

			c.hub.broadcast <- msg
		} else {
			// No valid target specified
			continue
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Printf("Error setting write deadline: %v", err)
				return
			}

			if !ok {
				if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					log.Printf("Error sending close message: %v", err)
				}
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, err = w.Write(message)
			if err != nil {
				log.Printf("Error writing message: %v", err)
				return
			}

			n := len(c.send)
			for range n {
				_, err := w.Write(<-c.send)
				if err != nil {
					log.Printf("Error writing message: %v", err)
					return
				}
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				log.Printf("Error setting write deadline: %v", err)
				return
			}

			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
