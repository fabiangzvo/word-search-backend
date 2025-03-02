package sockets

import (
	"github.com/gin-gonic/gin"
	socketIo "github.com/googollee/go-socket.io"

	"word-search/pkg/logger"
)

// ServeWS Function to handle websocket connection and register client to hub and start goroutines
func ServeWS(router *gin.Engine) *socketIo.Server {
	const section = "client.ServeWS"
	logger.Log.Infoln(section, "starting")

	server := socketIo.NewServer(nil)

	server.OnConnect("/", func(s socketIo.Conn) error {
		s.Emit("welcome", "welcome to server Socket.IO")
		return nil
	})

	// handle custom events
	server.OnEvent("/", "sendMessage", func(s socketIo.Conn, msg string) {
		s.Emit("receiveMessage", "server say: "+msg)
	})

	server.OnDisconnect("/", func(s socketIo.Conn, reason string) {
		logger.Log.Println("Cliente desconectado:", reason)
	})

	server.OnError("/", func(s socketIo.Conn, e error) {
		logger.Log.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketIo.Conn, reason string) {
		logger.Log.Println("closed", reason)
	})

	go func() {
		if err := server.Serve(); err != nil {
			logger.Log.Fatalf("socketio listen error: %s\n", err)
		}
	}()

	router.GET("/socket.io/*any", gin.WrapH(server))
	router.POST("/socket.io/*any", gin.WrapH(server))

	logger.Log.Infoln(section, "finished")

	return server
}
