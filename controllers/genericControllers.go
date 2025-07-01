package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/conzorkingkong/conazon-users-and-auth/types"
)

func Root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(types.Response{Status: http.StatusNotFound, Message: "invalid path " + r.URL.RequestURI(), Data: ""})
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}
