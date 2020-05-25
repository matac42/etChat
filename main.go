package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

type melodyHandler struct {
	melo *melody.Melody
}

func createMelodyHandler() melodyHandler {
	mel := melodyHandler{}
	m := melody.New()

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		m.Broadcast(msg)
	})

	m.HandleConnect(func(s *melody.Session) {
		log.Printf("websocket connection open. [session: %#v]\n", s)
	})

	m.HandleDisconnect(func(s *melody.Session) {
		log.Printf("websocket connection close. [session: %#v]\n", s)
	})

	mel.melo = m
	return mel
}

func chatFunc(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, "html/chat.html")
}

func (e *melodyHandler) wsHandler(c *gin.Context) {
	e.melo.HandleRequest(c.Writer, c.Request)
}

func logInHandler(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, "html/login.html")
}

func main() {
	r := gin.Default() //ginは基本的にgin.Default()の返す構造体のメソッド経由でないと操作できない．
	r.LoadHTMLGlob("html/*.html")

	cmelody := createMelodyHandler()

	v1 := r.Group("/")
	{
		v1.GET("chat", chatFunc)
		v1.GET("ws", cmelody.wsHandler)
		v1.GET("login", logInHandler)
	}
	r.Run(":8080")
}
