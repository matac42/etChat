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
	"github.com/matac42/etChat/database"
)

// RedirectAuthenticateClient fires fn when a /oauth connection.
func RedirectAuthenticateClient(c *gin.Context) {
	authURL := "https://github.com/login/oauth/authorize?client_id=" + os.Getenv("GithubClientID")
	c.Redirect(http.StatusMovedPermanently, authURL)
}

// LogInClient fires fn when a /login connection.
func LogInClient(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, "html/login.html")
}

// NameNotFound checks if an name exists in the db.
func NameNotFound(t string) bool {
	//SQLConnectはこの関数外でやって受け取る形が良い
	db, err := database.SQLConnect()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	creEX := database.CredentialInfo{}
	find := db.First(&creEX, "login=?", t).RecordNotFound()

	return find
}

// GetCredentialInfo gets a CredentialInfo from token end point.
func GetCredentialInfo(c *gin.Context) *database.CredentialInfo {
	// この関数でかい.

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

	cre := database.CredentialInfo{}
	json.Unmarshal(byteArray, &cre)

	return &cre
}

// GetAccessTokenClient deal with callback.
func GetAccessTokenClient(c *gin.Context) {
	cre := GetCredentialInfo(c)

	cre.GetGithubUserData()

	db, err := database.SQLConnect()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	creEX := database.CredentialInfo{}

	db.AutoMigrate(creEX)

	find := db.First(&creEX, "login=?", cre.Login)

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
