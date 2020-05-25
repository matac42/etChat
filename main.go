package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

type melodyHandler struct {
	melo *melody.Melody
}

type credentialInfo struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
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

func redirectAuthrizeClient(c *gin.Context) {
	authURL := "https://github.com/login/oauth/authorize?client_id=" + githubClientID
	c.Redirect(http.StatusMovedPermanently, authURL)
}

func getAccessTokenClient(c *gin.Context) {
	code := c.Request.URL.Query().Get("code")
	state := c.Request.URL.Query().Get("state")
	if state == "" {
		fmt.Println("state is empty")
	}
	values := url.Values{}
	values.Add("code", code)
	values.Add("client_id", githubClientID)
	values.Add("client_secret", githubClientSecret)
	req, err := http.NewRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)

	var cre *credentialInfo
	json.Unmarshal(byteArray, &cre)

	c.Redirect(http.StatusMovedPermanently, "/chat")

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
		v1.GET("oauth", redirectAuthrizeClient)
		v1.GET("callback", getAccessTokenClient)
	}
	r.Run(":8080")
}
