package sockets

import (
	"word-search/pkg/logger"
)

//Hub is a struct that holds all the clients and the messages that are sent to them
type Hub struct {
	// Registered clients.
	clients map[string]map[*Client]bool
	//Unregistered clients.
	unregister chan *Client
	// Register requests from the clients.
	register chan *Client
	// Inbound messages from the clients.
	broadcast chan Message
}

//Message struct to hold message data
type Message struct {
	Type      string  `json:"type"`
	Sender    string  `json:"sender"`
	Recipient string  `json:"recipient"`
	Content   string  `json:"content"`
	ID        string  `json:"id"`
}

//NewHub function to create a new hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		unregister: make(chan *Client),
		register:   make(chan *Client),
		broadcast:  make(chan Message),
	}
}

//Run function to run the hub
func (h *Hub) Run() {
	for {
		select {
		// Register a client.
		case client := <-h.register:
			h.RegisterNewClient(client)
			// Unregister a client.
		case client := <-h.unregister:
			h.RemoveClient(client)
			// Broadcast a message to all clients.
		case message := <-h.broadcast:

			//Check if the message is a type of "message"
			h.HandleMessage(message)

		}
	}
}

//RegisterNewClient check if room exists and if not create it and add client to it
func (h *Hub) RegisterNewClient(client *Client) {
	const section = "hub.RegisterNewClient"
	logger.Log.Infoln(section, "starting")

	connections := h.clients[client.ID]
	if connections == nil {
		connections = make(map[*Client]bool)
		h.clients[client.ID] = connections
	}
	h.clients[client.ID][client] = true

	logger.Log.Infoln("Size of clients: ", len(h.clients[client.ID]))
	logger.Log.Infoln(section, "finished")
}

//RemoveClient function to remvoe client from room
func (h *Hub) RemoveClient(client *Client) {
	const section = "hub.RemoveClient"
	logger.Log.Infoln(section, "starting")
	if _, ok := h.clients[client.ID]; ok {
		delete(h.clients[client.ID], client)
		close(client.send)

		logger.Log.Infoln(section, "removed client", client.ID)
	}

	logger.Log.Infoln(section, "finished")
}

//HandleMessage function to handle message based on type of message
func (h *Hub) HandleMessage(message Message) {
const section = "hub.HandleMessage"
	logger.Log.Infoln(section, "starting")
	logger.Log.Infof("message: %+v", message)
	//Check if the message is a type of "message"
	if message.Type == "message" {
		clients := h.clients[message.ID]
		logger.Log.Infof("clients: %+v", clients)
		for client := range clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(h.clients[message.ID], client)
			}
		}
	}

	//Check if the message is a type of "notification"
	if message.Type == "notification" {
		logger.Log.Infoln(section, "Notification: ", message.Content)
		clients := h.clients[message.Recipient]
		for client := range clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(h.clients[message.Recipient], client)
			}
		}
	}
logger.Log.Infoln(section, "finished")
}