package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	
	"time"
	_"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

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

	// db接続
    db, err := sqlConnect()
    if err != nil {
        panic(err.Error())
    }
    defer db.Close()

    error := db.Create(&Users{
        Name:     "testtarou",
        Age:      18,
        Address:  "tokyo",
        UpdateAt: getDate(),
    }).Error
    if error != nil {
        fmt.Println(error)
    } else {
        fmt.Println("データ追加成功")
    }	
		
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




func getDate() string {
    const layout = "2006-01-02 15:04:05"
    now := time.Now()
    return now.Format(layout)
}

// SQLConnect DB接続
func sqlConnect() (database *gorm.DB, err error) {
    DBMS := "mysql"
    USER := "jb5"
    PASS := "h19life"
    PROTOCOL := "tcp(localhost:3306)"
    DBNAME := "et"

    CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?charset=utf8&parseTime=true&loc=Asia%2FTokyo"
    return gorm.Open(DBMS, CONNECT)
}
// Users ユーザー情報のテーブル情報
type Users struct {
    ID       int
    Name     string `json:"name"`
    Age      int    `json:"age"`
    Address  string `json:"address"`
    UpdateAt string `json:"updateAt" sql:"not null;type:date"`
}
