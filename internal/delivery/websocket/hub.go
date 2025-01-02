package websocket

import (
	"chat-app/internal/domain/entities"
	domain "chat-app/internal/domain/interfaces"
	"chat-app/internal/dto"
	"encoding/json"
	"log"
)

type Hub struct {
	clients        map[string]*Client
	groups         map[string]map[string]*Client
	broadcast      chan BroadcastMessage
	register       chan *Client
	unregister     chan *Client
	messageUseCase domain.MessageUseCase
}

func NewHub(messageUseCase domain.MessageUseCase) *Hub {
	return &Hub{
		clients:        make(map[string]*Client),
		groups:         make(map[string]map[string]*Client),
		broadcast:      make(chan BroadcastMessage),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		messageUseCase: messageUseCase,
	}
}

func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.UserID] = client

			messages, err := h.messageUseCase.HandleNewConnection(client.UserID)
			if err != nil {
				log.Printf("Error handling new connection: %v", err)
				continue
			}

			if messages != nil {
				for _, msg := range messages {
					msgBytes, err := json.Marshal(msg)
					if err != nil {
						continue
					}
					client.Send <- msgBytes
				}
			}

		case client := <-h.unregister:
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)
			}

		case bm := <-h.broadcast:
			var wsMsg dto.WSMessage

			if err := json.Unmarshal(bm.Message, &wsMsg); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}

			// Handle read receipt melalui usecase
			if wsMsg.Action == "read" {
				err := h.messageUseCase.HandleReadReceipt(wsMsg.MessageId, wsMsg.RecipientId)
				if err != nil {
					log.Printf("Error handling read receipt: %v", err)
				}
				continue
			}

			// Handle new message melalui usecase
			msg, err := h.messageUseCase.HandleNewMessage(bm.UserId, wsMsg)
			if err != nil {
				log.Printf("Error handling new message: %v", err)
				continue
			}

			// Handle pengiriman pesan berdasarkan tipe
			switch msg.Type {
			case entities.TypePrivate:
				if recipient, ok := h.clients[msg.RecipientId]; ok {
					err := h.messageUseCase.HandleMessageDelivery(msg, msg.RecipientId)
					if err != nil {
						log.Printf("Error handling message delivery: %v", err)
					}

					msgBytes, err := json.Marshal(msg)
					if err != nil {
						log.Printf("Error marshaling message: %v", err)
						continue
					}

					select {
					case recipient.Send <- msgBytes:
					default:
						close(recipient.Send)
						delete(h.clients, recipient.UserID)
					}
				}

			case entities.TypeGroup:
				group, err := h.messageUseCase.GetGroup(msg.RecipientId)
				if err != nil {
					log.Printf("Error getting group: %v", err)
					continue
				}

				msgBytes, err := json.Marshal(msg)
				if err != nil {
					log.Printf("Error marshaling message: %v", err)
					continue
				}

				for _, member := range group.Members {
					if recipient, ok := h.clients[member.UserId]; ok {
						if member.UserId == msg.SenderId {
							continue
						}

						err := h.messageUseCase.HandleMessageDelivery(msg, member.UserId)
						if err != nil {
							log.Printf("Failed to deliver message to user %s: %v", member.UserId, err)
							continue
						}

						select {
						case recipient.Send <- msgBytes:
							log.Printf("Message sent to user %s", member.UserId)
						default:
							close(recipient.Send)
							delete(h.clients, member.UserId)
						}
					}
				}
			}
		}
	}
}
