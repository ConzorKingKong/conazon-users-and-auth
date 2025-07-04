package controllers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/conzorkingkong/conazon-users-and-auth/config"
	"github.com/conzorkingkong/conazon-users-and-auth/token"
	"github.com/conzorkingkong/conazon-users-and-auth/types"
	"github.com/jackc/pgx/v5"
)

func Verify(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "GET" {
		log.Println("Method Not Allowed")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(types.Response{Status: http.StatusMethodNotAllowed, Message: "Method Not Allowed", Data: ""})
		return
	}

	jwt, err := token.ValidateAndReturnSession(w, r)

	if err != nil {
		return
	}

	json.NewEncoder(w).Encode(types.TokenResponse{Status: http.StatusOK, Message: "Success", Data: jwt})
}

func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		log.Println("Method Not Allowed")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(types.Response{Status: http.StatusMethodNotAllowed, Message: "Method Not Allowed", Data: ""})
		return
	}

	// generate random state for CSRF protection in google oauth
	b := make([]byte, 128)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	state := base64.URLEncoding.EncodeToString(b)

	cookie := http.Cookie{
		Name:     "state",
		Value:    state,
		HttpOnly: true,
		Secure:   config.SecureCookie,
		Path:     "/auth/google",
	}

	http.SetCookie(w, &cookie)

	http.Redirect(w, r, config.GoogleOauthConfig.AuthCodeURL(state), http.StatusTemporaryRedirect)
}

// This function handles creating a user if they dont exist
// This needs to be changed to posting to the /users endpoint
// To handle that later
func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	state := r.URL.Query().Get("state")

	// if state from /auth/google/login doesnt match reject request
	cookie, err := r.Cookie("state")
	if err != nil {
		log.Println("Error getting cookie")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(types.Response{Status: http.StatusForbidden, Message: "State Mismatch", Data: ""})
		return
	}

	if state != cookie.Value {
		log.Println("Request came in with wrong state")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(types.Response{Status: http.StatusForbidden, Message: "State Mismatch", Data: ""})
		return
	}

	code := r.URL.Query().Get("code")

	oauthToken, err := config.GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Error exchanging token: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.Response{Status: http.StatusInternalServerError, Message: "Internal Service Error", Data: ""})
		return
	}

	idTokenPayload := strings.Split(oauthToken.Extra("id_token").(string), ".")[1]

	value, err := base64.RawStdEncoding.DecodeString(idTokenPayload)
	if err != nil {
		log.Printf("Error decoding token: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.Response{Status: http.StatusInternalServerError, Message: "Internal Service Error", Data: ""})
		return
	}

	TokenData := types.IdTokenPayload{}
	json.Unmarshal(value, &TokenData)

	conn, err := pgx.Connect(context.Background(), config.DatabaseURLEnv)
	if err != nil {
		log.Printf("Error connecting to database: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.Response{Status: http.StatusInternalServerError, Message: "Internal Service Error", Data: ""})
		return
	}

	defer conn.Close(context.Background())

	var id int

	err = conn.QueryRow(context.Background(), "select id from users.users where google_id=$1", TokenData.Sub).Scan(&id)
	if err != nil {
		// user does not exist, save them
		_, err2 := conn.Exec(context.Background(), "insert into users.users (name, email, picture, google_id) values ($1, $2, $3, $4)", TokenData.Name, TokenData.Email, TokenData.Picture, TokenData.Sub)
		if err2 != nil {
			log.Printf("Error saving user: %s", err2)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(types.Response{Status: http.StatusInternalServerError, Message: "Internal Service Error", Data: ""})
			return
		}
		err3 := conn.QueryRow(context.Background(), "select id from users.users where google_id=$1", TokenData.Sub).Scan(&id)
		if err3 != nil {
			log.Printf("Error getting user id: %s", err3)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(types.Response{Status: http.StatusInternalServerError, Message: "Internal Service Error", Data: ""})
			return
		}
		// create session with custom session jwt
		jwt, err := token.CreateToken(id)
		if err != nil {
			log.Printf("Error creating token: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(types.Response{Status: http.StatusInternalServerError, Message: "Internal Service Error", Data: ""})
			return
		}

		config.SetCookieSession(w, jwt)

		http.Redirect(w, r, config.ProtocolEnv+"://"+config.HostnameEnv, http.StatusTemporaryRedirect)
	} else {
		// user exists, return custom session jwt
		jwt, err := token.CreateToken(id)
		if err != nil {
			log.Printf("Error creating token: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(types.Response{Status: http.StatusInternalServerError, Message: "Internal Service Error", Data: ""})
			return
		}

		config.SetCookieSession(w, jwt)

		http.Redirect(w, r, "http://localhost", http.StatusTemporaryRedirect)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "DELETE" {
		log.Println("Method Not Allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(types.Response{Status: http.StatusMethodNotAllowed, Message: "Method Not Allowed", Data: ""})
		return
	}

	_, err := r.Cookie("JWTToken")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			log.Printf("cookie not found")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(types.Response{Status: http.StatusBadRequest, Message: "cookie not found", Data: ""})
		default:
			log.Printf("Cookie err: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(types.Response{Status: http.StatusInternalServerError, Message: "server error", Data: ""})
		}
		return
	}

	config.DeleteCookieSession(w)

	json.NewEncoder(w).Encode(types.Response{Status: http.StatusOK, Message: "Logged out", Data: ""})
}
