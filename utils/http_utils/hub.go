package http_utils

// TODO: remove topic from hub if no clients are connected for that topic

var Hub_global = NewHub()

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool
	// Inbound messages from the clients.
	broadcast chan *Broadcast
	// Register requests from the clients.
	register chan *Client
	// Unregister requests from clients.
	unregister chan *Client
	// topics
	topics []string
}

//
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan *Broadcast),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

//
func (h *Hub) Run() {
	for {
		select {
		// register
		case client := <-h.register:
			h.clients[client] = true
		// unregister
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		// broadcast
		case broadcast_data := <-h.broadcast:
			for client := range h.clients {
				if client.topic == broadcast_data.topic {
					if client.conn.UnderlyingConn() != broadcast_data.client {
						select {
						case client.send <- broadcast_data.message:
						default:
							close(client.send)
							delete(h.clients, client)
						}
					}
				}
			}
		}
	}
}
