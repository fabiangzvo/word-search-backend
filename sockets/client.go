package sockets

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"word-search/pkg/logger"
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
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//Client struct for websocket connection and message sending
type Client struct {
	ID   string
	Conn *websocket.Conn
	send chan Message
	hub  *Hub
}

//NewClient creates a new client
func NewClient(id string, conn *websocket.Conn, hub *Hub) *Client {
	return &Client{ID: id, Conn: conn, send: make(chan Message, 256), hub: hub}
}

//Client goroutine to read messages from client
func (c *Client) Read() {
	const section = "client.Read"
	logger.Log.Infoln(section, "starting")

	defer func() {
		c.hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var msg Message
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("Error: ", err)
			break
		}
		c.hub.broadcast <- msg
	}

	logger.Log.Infoln(section, "finished")
}

//Client goroutine to write messages to client
func (c *Client) Write() {
	const section = "client.Write"
	logger.Log.Infoln(section, "starting")

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()

		logger.Log.Infoln(section, "finished")
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			} else {
				err := c.Conn.WriteJSON(message)
				if err != nil {
					fmt.Println("Error: ", err)
					break
				}
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

//Close Client closing channel to unregister client
func (c *Client) Close() {
	close(c.send)
}

//ServeWS Function to handle websocket connection and register client to hub and start goroutines
func ServeWS(ctx *gin.Context, gameID string, hub *Hub) {
	const section = "client.ServeWS"
	logger.Log.Infoln(section, gameID,"starting")

	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := NewClient(gameID, ws, hub)
	hub.register <- client

	go client.Write()
	go client.Read()

	logger.Log.Infoln(section, gameID,"finished")
}