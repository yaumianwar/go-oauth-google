package handler

import (
	"crypto/rand"
	"golang.org/x/oauth2"
	"encoding/base64"
	//"github.com/spf13/viper"
	"github.com/gin-gonic/gin"
	"net/http"
	"io/ioutil"
	"github.com/yaumianwar/go-oauth-google/model"
	"encoding/json"
	googlepkg "golang.org/x/oauth2/google"
	"github.com/yaumianwar/go-oauth-google/google"
	"github.com/gin-gonic/contrib/sessions"
	"log"
	"os"
)

var cred Credentials
var conf *oauth2.Config

func RandToken(l int) string {
	b := make([]byte, l)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func getLoginURL(state string) string {
	return conf.AuthCodeURL(state)
}

type Credentials struct {
	Cid     string `json:"cid"`
	Csecret string `json:"csecret"`
}


func init() {
	file, err := ioutil.ReadFile("./creds.json")
	if err != nil {
		log.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	json.Unmarshal(file, &cred)

	conf = &oauth2.Config{
		ClientID:     cred.Cid,
		ClientSecret: cred.Csecret,
		RedirectURL:  "http://127.0.0.1:3000/GoogleCallback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
		},
		Endpoint: googlepkg.Endpoint,
	}
}

// IndexHandler handels /.
func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{})
}

func AuthHandler(svc *google.Service) func(*gin.Context) {
	return func(c *gin.Context) {
		// Handle the exchange code to initiate a transport.
		session := sessions.Default(c)
		retrievedState := session.Get("state")
		queryState := c.Request.URL.Query().Get("state")
		if retrievedState != queryState {
			log.Printf("Invalid session state: retrieved: %s; Param: %s", retrievedState, queryState)
			c.HTML(http.StatusUnauthorized, "error.tmpl", gin.H{"message": "Invalid session state."})
			return
		}
		code := c.Request.URL.Query().Get("code")
		tok, err := conf.Exchange(oauth2.NoContext, code)
		if err != nil {
			log.Println(err)
			c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Login failed. Please try again."})
			return
		}

		client := conf.Client(oauth2.NoContext, tok)
		userinfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		defer userinfo.Body.Close()
		data, _ := ioutil.ReadAll(userinfo.Body)
		u := model.User{}
		if err = json.Unmarshal(data, &u); err != nil {
			log.Println(err)
			c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Error marshalling response. Please try agian."})
			return
		}
		session.Set("user-id", u.Email)
		err = session.Save()
		if err != nil {
			log.Println(err)
			c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Error while saving session. Please try again."})
			return
		}
		seen := false
		if _, err := svc.LoadUser(u.Email); err == nil {
			seen = true
		} else {
			err = svc.SaveUser(&u)
			if err != nil {
				log.Println(err)
				c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Error while saving user. Please try again."})
				return
			}
		}
		c.HTML(http.StatusOK, "battle.tmpl", gin.H{"email": u.Email, "seen": seen})
	}
}


// LoginHandler handles the login procedure.
func LoginHandler(c *gin.Context) {
	state := RandToken(32)
	session := sessions.Default(c)
	session.Set("state", state)
	session.Save()
	link := getLoginURL(state)
	c.HTML(http.StatusOK, "auth.tmpl", gin.H{"link": link})
}

// FieldHandler is a rudementary handler for logged in users.
func FieldHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user-id")
	c.HTML(http.StatusOK, "field.tmpl", gin.H{"user": userID})
}

