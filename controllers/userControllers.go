package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/conzorkingkong/conazon-users-and-auth/config"
	"github.com/conzorkingkong/conazon-users-and-auth/helpers"
	"github.com/conzorkingkong/conazon-users-and-auth/token"
	"github.com/conzorkingkong/conazon-users-and-auth/types"
	"github.com/jackc/pgx/v5"
)

func Users(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// GET ALL USERS
	if r.Method == "GET" {
		page := 1
		limit := 10

		// Get page from query params
		if pageParam := r.URL.Query().Get("page"); pageParam != "" {
			if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
				page = p
			}
		}

		// Get limit from query params (max 100)
		if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
			if l, err := strconv.Atoi(limitParam); err == nil && l > 0 && l <= 100 {
				limit = l
			}
		}

		// Calculate offset
		offset := (page - 1) * limit

		conn, err := pgx.Connect(context.Background(), config.DatabaseURLEnv)
		if err != nil {
			log.Printf("Error connecting to database: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(types.Response{Status: http.StatusInternalServerError, Message: "Internal Service Error", Data: ""})
			return
		}

		defer conn.Close(context.Background())

		// Get total count for pagination metadata
		var totalCount int
		err = conn.QueryRow(context.Background(), "SELECT COUNT(*) FROM users.users").Scan(&totalCount)
		if err != nil {
			log.Printf("Error counting users: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(types.Response{Status: http.StatusInternalServerError, Message: "Internal Service Error", Data: ""})
			return
		}

		// Get paginated users
		rows, err := conn.Query(context.Background(),
			"SELECT id, name, picture FROM users.users ORDER BY id LIMIT $1 OFFSET $2",
			limit, offset)
		if err != nil {
			log.Printf("Error getting users: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(types.Response{Status: http.StatusInternalServerError, Message: "Internal Service Error", Data: ""})
			return
		}
		defer rows.Close()

		var users []types.User
		for rows.Next() {
			var user types.User
			err := rows.Scan(&user.ID, &user.Name, &user.Picture)
			if err != nil {
				log.Printf("Error scanning user: %s", err)
				continue
			}
			users = append(users, user)
		}

		// Create pagination metadata
		totalPages := (totalCount + limit - 1) / limit

		response := types.PaginatedUsersResponse{
			Status:  http.StatusOK,
			Message: "Success",
		}
		response.Data.Users = users
		response.Data.Pagination.Page = page
		response.Data.Pagination.Limit = limit
		response.Data.Pagination.TotalItems = totalCount
		response.Data.Pagination.TotalPages = totalPages
		response.Data.Pagination.HasNext = page < totalPages
		response.Data.Pagination.HasPrev = page > 1

		json.NewEncoder(w).Encode(response)

		// CREATE A NEW USER
		// TO BE IMPLEMENTED
	} else if r.Method == "POST" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.Response{Status: http.StatusInternalServerError, Message: "Coming Soon", Data: ""})
		return
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(types.Response{Status: http.StatusMethodNotAllowed, Message: "Method Not Allowed", Data: ""})
		return
	}

}

func UserId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	routeId, parsedRouteId, err := helpers.RouteIdHelper(w, r)
	if err != nil {
		return
	}

	// GETS ONLY PUBLIC INFORMATION
	// USE ME ENDPOINT TO GET ALL USERS INFO
	// WILL BE UPDATED LATER
	if r.Method == "GET" {
		conn, err := pgx.Connect(context.Background(), config.DatabaseURLEnv)
		if err != nil {
			log.Printf("Error connecting to database: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(types.Response{Status: http.StatusInternalServerError, Message: "Internal Service Error", Data: ""})
			return
		}

		defer conn.Close(context.Background())

		user := types.User{}

		err = conn.QueryRow(context.Background(), "select name, picture from users.users where id=$1", routeId).Scan(&user.Name, &user.Picture)
		if err != nil {
			log.Printf("Error getting user with id %s - %s", routeId, err)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(types.Response{Status: http.StatusNotFound, Message: "User not found", Data: ""})
			return
		}

		json.NewEncoder(w).Encode(types.UserResponse{Status: http.StatusOK, Message: "Success", Data: user})

		// UPDATE WHOLE USER
		// MUST POST NAME, EMAIL, AND PICTURE URL
	} else if r.Method == "POST" {

		TokenData, err := token.ValidateAndReturnSession(w, r)
		if err != nil {
			return
		}

		// CHECK IF URL AND TOKEN MATCH
		if TokenData.Id != parsedRouteId {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(types.Response{Status: http.StatusForbidden, Message: "You are not authorized to complete this action", Data: ""})
			return
		}

		// Pull data from request
		updateUser := types.User{}

		err = json.NewDecoder(r.Body).Decode(&updateUser)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(types.Response{Status: http.StatusBadRequest, Message: "Invalid JSON", Data: ""})
			return
		}

		conn, err := pgx.Connect(context.Background(), config.DatabaseURLEnv)
		if err != nil {
			log.Printf("Error connecting to database: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(types.Response{Status: http.StatusInternalServerError, Message: "internal service error", Data: ""})
			return
		}

		defer conn.Close(context.Background())

		// PULL CURRENT USER DATA TO FILL IN NON CHANGED DATA
		var currentUser types.User

		err = conn.QueryRow(context.Background(), "SELECT name, email, picture FROM users.users WHERE id=$1", TokenData.Id).Scan(&currentUser.Name, &currentUser.Email, &currentUser.Picture)
		if err != nil {
			log.Printf("Error reading user with id %d - %s", TokenData.Id, err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(types.Response{Status: http.StatusInternalServerError, Message: "internal service error", Data: ""})
			return
		}

		// Merge updates (only update non-empty fields)
		if updateUser.Name != "" {
			currentUser.Name = updateUser.Name
		}
		if updateUser.Email != "" {
			currentUser.Email = updateUser.Email
		}
		if updateUser.Picture != "" {
			currentUser.Picture = updateUser.Picture
		}

		_, err = conn.Exec(context.Background(), "update users.users set name=$1, email=$2, picture=$3 where id=$4", currentUser.Name, currentUser.Email, currentUser.Picture, TokenData.Id)
		if err != nil {
			log.Printf("Error updating user with id %d - %s", TokenData.Id, err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(types.Response{Status: http.StatusInternalServerError, Message: "internal service error", Data: ""})
			return
		}

		json.NewEncoder(w).Encode(types.UserResponse{Status: http.StatusOK, Message: fmt.Sprintf("user %d updated", TokenData.Id), Data: updateUser})

		// DELETE USER IF ROUTE AND TOKEN ID MATCH
	} else if r.Method == "DELETE" {

		TokenData, err := token.ValidateAndReturnSession(w, r)
		if err != nil {
			return
		}

		// CHECK IF URL AND TOKEN MATCH
		if TokenData.Id != parsedRouteId {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(types.Response{Status: http.StatusForbidden, Message: "internal service error", Data: ""})
			return
		}

		conn, err := pgx.Connect(context.Background(), config.DatabaseURLEnv)
		if err != nil {
			log.Printf("Error connecting to database: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(types.Response{Status: http.StatusInternalServerError, Message: "internal service error", Data: ""})
			return
		}

		defer conn.Close(context.Background())

		_, err = conn.Exec(context.Background(), "delete from users.users where id=$1", TokenData.Id)
		if err != nil {
			log.Printf("Error deleting user with id %d - %s", TokenData.Id, err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(types.Response{Status: http.StatusInternalServerError, Message: "internal service error", Data: ""})
			return
		}

		config.DeleteCookieSession(w)

		json.NewEncoder(w).Encode(types.Response{Status: http.StatusOK, Message: fmt.Sprintf("user %d deleted", TokenData.Id), Data: ""})
	} else {
		log.Println("Method Not Allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(types.Response{Status: http.StatusMethodNotAllowed, Message: "method not allowed", Data: ""})
		return
	}
}

func Me(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	TokenData, err := token.ValidateAndReturnSession(w, r)
	if err != nil {
		return
	}

	if r.Method == "GET" {
		conn, err := pgx.Connect(context.Background(), config.DatabaseURLEnv)
		if err != nil {
			log.Printf("Error connecting to database: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(types.Response{Status: http.StatusInternalServerError, Message: "Internal Service Error", Data: ""})
			return
		}

		defer conn.Close(context.Background())

		user := types.User{}

		err = conn.QueryRow(context.Background(), "select name, email, picture from users.users where id=$1", TokenData.Id).Scan(&user.Name, &user.Email, &user.Picture)
		if err != nil {
			log.Printf("Error getting user with id %d - %s", TokenData.Id, err)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(types.Response{Status: http.StatusNotFound, Message: "User not found", Data: ""})
			return
		}

		fmt.Printf("USER INFO %+v\n", user)

		json.NewEncoder(w).Encode(types.UserResponse{Status: http.StatusOK, Message: "Success", Data: user})
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(types.Response{Status: http.StatusMethodNotAllowed, Message: "Method not allowed", Data: ""})
		return
	}
}
