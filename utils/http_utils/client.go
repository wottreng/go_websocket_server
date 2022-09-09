package http_utils

import (
	"bytes"
	"file_utils"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Broadcast struct
type Broadcast struct {
	topic   string
	message []byte
	client  net.Conn
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	// hub the client is connected to.
	hub *Hub
	// topic subscribed
	topic string
	// The websocket connection.
	conn *websocket.Conn
	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		err := c.conn.Close()
		if err != nil {
			return
		}
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	var broadcast_data Broadcast
	for {
		//fmt.Println(c.hub.clients)
		//fmt.Println(c.conn.UnderlyingConn())
		broadcast_data.client = c.conn.UnderlyingConn()
		broadcast_data.topic = c.topic
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				break
			}
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				file_utils.Log_error_to_file(err, "[client.go][readPump]")
			}
			break
		}
		broadcast_data.message = bytes.TrimSpace(bytes.Replace(message, []byte{'\n'}, []byte{' '}, -1))

		c.hub.broadcast <- &broadcast_data
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			//fmt.Println("client topic: " + c.topic)
			//if c.topic == broadcast_topic {
			//	if c.conn.UnderlyingConn() != broadcast_client {
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}
			if err := w.Close(); err != nil {
				return
			}
			//	}
			//}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serve_ws_request(hub *Hub, w http.ResponseWriter, r *http.Request, topic string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		file_utils.Log_error_to_file(err, "[client.go][serveWs]")
		return
	}
	client := &Client{hub: hub, topic: topic, conn: conn, send: make(chan []byte, 256)}
	//fmt.Println(client.conn)
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in new goroutines
	go client.writePump()
	go client.readPump()
}
