package oauth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/matac42/etChat/database"
)

// CredentialInfo implements a oauth2 access token etc...
type CredentialInfo struct {
	ID          int
	Name        string `json:"login"`
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

// CreateCredentialInfo create an instance of CredentialInfo.
func CreateCredentialInfo() *CredentialInfo {
	cre := CredentialInfo{}
	return &cre
}

// RedirectAuthenticateClient fires fn when a /oauth connection.
func RedirectAuthenticateClient(c *gin.Context) {
	authURL := "https://github.com/login/oauth/authorize?client_id=" + os.Getenv("GithubClientID")
	c.Redirect(http.StatusMovedPermanently, authURL)
}

// LogInClient fires fn when a /login connection.
func LogInClient(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, "html/login.html")
}

// AccessTokenNotFound checks if an access token exists in the db.
func NameNotFound(t string) bool {
	//SQLConnectはこの関数外でやって受け取る形が良い
	db, err := database.SQLConnect()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	creEX := CredentialInfo{}
	find := db.First(&creEX, "name=?", t).RecordNotFound()

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
	values.Add("client_id", os.Getenv("GithubClientID"))
	values.Add("client_secret", os.Getenv("GithubClientSecret"))
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

	GetGithubUserData(cre)

	db, err := database.SQLConnect()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	creEX := CredentialInfo{}
	find := db.First(&creEX, "name=?", cre.Name)

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

func GetGithubUserData(c *CredentialInfo) {

	values := url.Values{}
	req, err := http.NewRequest(
		"POST",
		"https://api.github.com/user",
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "token "+c.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(byteArray, &c)
}
