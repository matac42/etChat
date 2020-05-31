package database

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// CredentialInfo implements a oauth2 access token etc...
type CredentialInfo struct {
	gorm.Model
	Login       string `json:"login"`
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

// CreateCredentialInfo create an instance of CredentialInfo.
func CreateCredentialInfo() *CredentialInfo {
	cre := CredentialInfo{}
	return &cre
}

// SQLConnect establishes a mysql connection.
func SQLConnect() (database *gorm.DB, err error) {
	DBMS := os.Getenv("DBMS")
	USER := os.Getenv("USER")
	PASS := os.Getenv("PASS")
	PROTOCOL := os.Getenv("PROTOCOL")
	DBNAME := os.Getenv("DBNAME")

	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?charset=utf8&parseTime=true&loc=Asia%2FTokyo"
	return gorm.Open(DBMS, CONNECT)
}

// GetGithubUserData get a user data via github api.
func (c *CredentialInfo) GetGithubUserData() {

	req, err := http.NewRequest(
		"POST",
		"https://api.github.com/user",
		nil,
	)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", "token "+c.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(byteArray, &c)
}
