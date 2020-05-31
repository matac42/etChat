package oauth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/matac42/etChat/database"
)

// CredentialInfo implements a oauth2 access token etc...
type CredentialInfo struct {
	ID          int
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

// あとで治そう
const (
	GithubClientID     = "A"
	GithubClientSecret = "B"
)

// CreateCredentialInfo create an instance of CredentialInfo.
func CreateCredentialInfo() *CredentialInfo {
	cre := CredentialInfo{}
	return &cre
}

// RedirectAuthenticateClient fires fn when a /oauth connection.
func RedirectAuthenticateClient(c *gin.Context) {
	authURL := "https://github.com/login/oauth/authorize?client_id=" + GithubClientID
	c.Redirect(http.StatusMovedPermanently, authURL)
}

// LogInClient fires fn when a /login connection.
func LogInClient(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, "html/login.html")
}

// AccessTokenNotFound checks if an access token exists in the db.
func AccessTokenNotFound(t string) bool {
	db, err := database.SQLConnect()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	creEX := CredentialInfo{}
	find := db.First(&creEX, "access_token=?", t).RecordNotFound()

	return find
}

// GetCredentialInfo gets a CredentialInfo from token end point.
func GetCredentialInfo(c *gin.Context) *CredentialInfo {
	// first, get the authentication code.
	code := c.Request.URL.Query().Get("code")
	state := c.Request.URL.Query().Get("state")
	if state == "" {
		fmt.Println("state is empty")
	}

	// second, get the access token using authentication code.
	values := url.Values{}
	values.Add("code", code)
	values.Add("client_id", GithubClientID)
	values.Add("client_secret", GithubClientSecret)
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

	cre := CredentialInfo{}
	json.Unmarshal(byteArray, &cre)

	return &cre
}

// GetAccessTokenClient deal with callback.
func GetAccessTokenClient(c *gin.Context) {
	cre := GetCredentialInfo(c)

	db, err := database.SQLConnect()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	db.AutoMigrate(&CredentialInfo{})

	creEX := CredentialInfo{}
	find := db.First(&creEX, "access_token=?", cre.AccessToken)

	if find.RecordNotFound() {
		error := db.Create(&cre).Error
		if error != nil {
			fmt.Println(error)
		} else {
			fmt.Println("success addition access token to db!!!")
		}
	}

	c.Redirect(http.StatusMovedPermanently, "/chat")
}