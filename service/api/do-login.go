package api

import (
	"github.com/zoematr/WASAPhoto/service/api/reqcontext"
	"encoding/json"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"fmt"
)

// func that handles user login
func (rt *_router) handleLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// type of content will be json
	w.Header().Set("Content-Type", "application/json")

	// init var User and decode body of request
	var username string
	err := json.NewDecoder(r.Body).Decode(&username)
	fmt.Println(err)

	// if error during decoding, like not parseable JSON or invalid username, respond with 400 bad request
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if !validUsername(username) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Create user, insert in DB
	token, err := rt.db.CreateUser(username)
	// if there is error, like the user exists, returns token
	if err != nil {
		// user exists, token returned
		token, err = rt.db.GetToken(username)
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(map[string]int{"token": token})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("session: can't create response json")
		}
		return
	}
}
