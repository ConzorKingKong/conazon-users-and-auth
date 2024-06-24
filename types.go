package main

type User struct {
	ID       int    `json:"-"`
	Name     string `json:"name"`
	Email    string `json:"email,omitempty"`
	Picture  string `json:"picture"`
	GoogleID int    `json:"-"`
}

type IdTokenPayload struct {
	Iss            string `json:"iss"`
	Azp            string `json:"azp"`
	Aud            string `json:"aud"`
	Sub            string `json:"sub"`
	Hd             string `json:"hd"`
	Email          string `json:"email"`
	Email_verified bool   `json:"email_verified"`
	At_hash        string `json:"at_hash"`
	Name           string `json:"name"`
	Picture        string `json:"picture"`
	Given_name     string `json:"given_name"`
	Family_name    string `json:"family_name"`
	Iat            int    `json:"iat"`
	Exp            int    `json:"exp"`
}

type MyJWT struct {
	Id int `json:"id"`
}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

type TokenResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    MyJWT  `json:"data"`
}

type UserResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    User   `json:"data"`
}
