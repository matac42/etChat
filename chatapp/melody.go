package chatapp

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matac42/etChat/oauth"
	"gopkg.in/olahol/melody.v1"
)

// MelodyHandler implements a websocket manager.
type MelodyHandler struct {
	melo *melody.Melody
}

// ChatClient fires fn when a /chat connection.
func ChatClient(c *gin.Context) {
	cre := oauth.GetCredentialInfo(c)

	// セッションの存在で判断するように書き直す．
	if oauth.NameNotFound(cre.Name) {
		c.Redirect(http.StatusMovedPermanently, "/login")
	} else {
		http.ServeFile(c.Writer, c.Request, "html/chat.html")
	}
}

// MelodyClient fires fn when a /ws connects.
func (e *MelodyHandler) MelodyClient(c *gin.Context) {
	e.melo.HandleRequest(c.Writer, c.Request)
}

// CreateMelodyHandler establishes a websocket connection.
func CreateMelodyHandler() MelodyHandler {
	mel := MelodyHandler{}
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
