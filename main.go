package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var PORT, PORTExists = "", false
var JwtSecret, jwtSecretExists = "", false
var ClientIDEnv, ClientIDExists = "", false
var ClientSecretEnv, ClientSecretExists = "", false
var RedirectURLEnv, RedirectURLExists = "", false
var SecureCookieEnv, secureCookieExists = "", false
var DatabaseURLEnv, DatabaseURLExists = "", false

var GoogleOauthConfig = oauth2.Config{
	Scopes:   []string{"email", "profile", "openid"},
	Endpoint: google.Endpoint,
}

var SecureCookie = false

var SECRETKEY []byte

func main() {

	godotenv.Load()

	PORT, PORTExists = os.LookupEnv("PORT")
	JwtSecret, jwtSecretExists = os.LookupEnv("JWTSECRET")
	ClientIDEnv, ClientIDExists = os.LookupEnv("CLIENTID")
	ClientSecretEnv, ClientSecretExists = os.LookupEnv("CLIENTSECRET")
	RedirectURLEnv, RedirectURLExists = os.LookupEnv("REDIRECTURL")
	SecureCookieEnv, secureCookieExists = os.LookupEnv("SECURECOOKIE")
	DatabaseURLEnv, DatabaseURLExists = os.LookupEnv("DATABASEURL")

	SECRETKEY = []byte(JwtSecret)

	if !jwtSecretExists || !ClientIDExists || !ClientSecretExists || !RedirectURLExists || !DatabaseURLExists {
		log.Fatal("Required environment variable not set")
	}

	GoogleOauthConfig.ClientID = ClientIDEnv
	GoogleOauthConfig.ClientSecret = ClientSecretEnv
	GoogleOauthConfig.RedirectURL = RedirectURLEnv

	if secureCookieExists {
		SecureCookie, _ = strconv.ParseBool(SecureCookieEnv)
	}

	if !PORTExists {
		PORT = "8080"
	}

	http.HandleFunc("/", Root)

	http.HandleFunc("/auth/google/login", googleLogin)
	http.HandleFunc("/auth/google/callback", googleCallback)
	http.HandleFunc("/logout", logout)

	http.HandleFunc("/verify", verify)
	http.HandleFunc("/me", me)

	http.HandleFunc("/users", users)
	http.HandleFunc("/users/{id}", userId)

	fmt.Println("server starting on port", PORT)
	http.ListenAndServe(":"+PORT, nil)
}
