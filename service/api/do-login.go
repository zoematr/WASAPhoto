package api

import (
	"github.com/zoematr/WASAPhoto/service/api/reqcontext"
	"encoding/json"
	"net/http"
	"github.com/julienschmidt/httprouter"
)

// func that handles user login
func (rt *_router) handleLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// type of content will be json
	w.Header().Set("Content-Type", "application/json")

	// init var User and decode body of request
	var username string
	err := json.NewDecoder(r.Body).Decode(&username)

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
		resp := map[string]int{"token": token}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(resp)
		//consider if there is an error, like the user can't be logged in
		if err != nil {
			ctx.Logger.WithError(err).WithField("username", username).Error("Can't login user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
//	w.WriteHeader(http.StatusOK)
	return
}

