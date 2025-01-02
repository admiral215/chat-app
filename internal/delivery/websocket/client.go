package websocket

import (
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type Client struct {
	Hub    *Hub
	Conn   *websocket.Conn
	Send   chan []byte
	UserID string
	Groups map[string]struct{}
}

const (
	// writeWait adalah waktu yang diizinkan untuk menulis pesan ke peer
	writeWait = 10 * time.Second

	// pongWait adalah waktu yang diizinkan untuk membaca pesan pong berikutnya dari peer
	pongWait = 60 * time.Second

	// pingPeriod adalah periode pengiriman ping ke peer
	// Harus kurang dari pongWait
	pingPeriod = (pongWait * 9) / 10

	// maxMessageSize adalah ukuran maksimum pesan yang diizinkan dari peer
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		err := c.Conn.Close()
		if err != nil {
			return
		}
	}()

	c.Conn.SetReadLimit(maxMessageSize)

	err := c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		return
	}

	c.Conn.SetPongHandler(func(string) error {
		err := c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			return err
		}
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		bm := BroadcastMessage{
			UserId:  c.UserID,
			Message: message,
		}

		c.Hub.broadcast <- bm
	}
}

// Implementasi WritePump
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()

		err := c.Conn.Close()
		if err != nil {
			return
		}
	}()

	for {
		select {
		case message, ok := <-c.Send:
			log.Printf("Received message, channel ok: %v, message length: %d", ok, len(message))
			err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				return
			}
			if !ok {
				err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					return
				}

				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			_, err = w.Write(message)
			if err != nil {
				return
			}

			n := len(c.Send)
			for i := 0; i < n; i++ {
				select {
				case additionalMsg := <-c.Send:
					_, err = w.Write(newline)
					if err != nil {
						return
					}
					_, err = w.Write(additionalMsg)
					if err != nil {
						return
					}
				default:
					break
				}
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				return
			}
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
