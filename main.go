package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

func chatFunc(c *gin.Context) {
	//chat processing
	log.Println("Websocket App start.")

	router := gin.Default()
	m := melody.New()

	rg := router.Group("/chat")
	rg.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "html/chat.html")
	})

	rg.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		m.Broadcast(msg)
	})

	m.HandleConnect(func(s *melody.Session) {
		log.Printf("websocket connection open. [session: %#v]\n", s)
	})

	m.HandleDisconnect(func(s *melody.Session) {
		log.Printf("websocket connection close. [session: %#v]\n", s)
	})

	// Listen and server on 0.0.0.0:8080
	router.Run(":8080")

	fmt.Println("Websocket App End.")
}

func main() {
	r := gin.Default()
	v1 := r.Group("/")
	{
		v1.GET("/chat", chatFunc)
	}
	r.Run(":8080")
}
