package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/conzorkingkong/conazon-users-and-auth/config"
	"github.com/conzorkingkong/conazon-users-and-auth/controllers"
	"github.com/joho/godotenv"
)

var PORT, PORTExists = "", false
var JwtSecret, jwtSecretExists = "", false
var ClientIDEnv, ClientIDExists = "", false
var ClientSecretEnv, ClientSecretExists = "", false
var RedirectURLEnv, RedirectURLExists = "", false
var SecureCookieEnv, secureCookieExists = "", false
var DatabaseURLEnv, DatabaseURLExists = "", false


func main() {

	godotenv.Load()

	PORT, PORTExists = os.LookupEnv("PORT")
	JwtSecret, jwtSecretExists = os.LookupEnv("JWTSECRET")
	ClientIDEnv, ClientIDExists = os.LookupEnv("CLIENTID")
	ClientSecretEnv, ClientSecretExists = os.LookupEnv("CLIENTSECRET")
	RedirectURLEnv, RedirectURLExists = os.LookupEnv("REDIRECTURL")
	SecureCookieEnv, secureCookieExists = os.LookupEnv("SECURECOOKIE")
	DatabaseURLEnv, DatabaseURLExists = os.LookupEnv("DATABASEURL")

	config.SECRETKEY = []byte(JwtSecret)

	if !jwtSecretExists || !ClientIDExists || !ClientSecretExists || !RedirectURLExists || !DatabaseURLExists {
		log.Fatal("Required environment variable not set")
	}

	config.GoogleOauthConfig.ClientID = ClientIDEnv
	config.GoogleOauthConfig.ClientSecret = ClientSecretEnv
	config.GoogleOauthConfig.RedirectURL = RedirectURLEnv
	config.DatabaseURLEnv = DatabaseURLEnv

	if secureCookieExists {
		config.SecureCookie, _ = strconv.ParseBool(SecureCookieEnv)
	}

	if !PORTExists {
		PORT = "8080"
	}

	http.HandleFunc("/", controllers.Root)

	http.HandleFunc("/auth/google/login", controllers.GoogleLogin)
	http.HandleFunc("/auth/google/callback", controllers.GoogleCallback)
	http.HandleFunc("/logout", controllers.Logout)

	http.HandleFunc("/verify", controllers.Verify)
	http.HandleFunc("/me", controllers.Me)

	http.HandleFunc("/users", controllers.Users)
	http.HandleFunc("/users/{id}", controllers.UserId)

	fmt.Println("server starting on port", PORT)
	http.ListenAndServe(":"+PORT, nil)
}
