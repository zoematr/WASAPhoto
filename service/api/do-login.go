package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/zoematr/WASAPhoto/service/api/reqcontext"
	"net/http"
	//	"fmt"
	//	"log"
	//	"strconv"
)

// func that handles user login
func (rt *_router) handleLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// type of content will be json
	ctx.Logger.Infof("handleLogin was called in the backend")

	// init var User and decode body of request
	var username string
	err := json.NewDecoder(r.Body).Decode(&username)

	// if error during decoding, like not parseable JSON or invalid username, respond with 400 bad request
	if err != nil {
		ctx.Logger.WithError(err).Error("login: error decoding body")
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if !validUsername(username) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userexists, err := rt.db.ExistsUser(username)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Create usert bc it does not exist
	if !userexists {
		token, err := rt.db.CreateUser(username) // TODO TOKEN IS REWRITTEN BEFORE CHECKING
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "plain/text")
		err = json.NewEncoder(w).Encode(token)
		if err != nil {
			ctx.Logger.WithError(err).WithField("username", username).Error("Can't login user")
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	ctx.Logger.Infof("The user exists and it's being logged in ")

	// if the user exists, returns token
	// user exists, token returned
	token, err := rt.db.GetTokenFromUsername(username)
	ctx.Logger.Infof("The value of the login token is: %d", token)
	if err != nil {
		ctx.Logger.WithError(err).WithField("username", username).Error("Can't login user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	/*
		// w.Header().Set("Content-Type", "plain/text")
		// err = json.NewEncoder(w).Encode(token)
		// consider if there is an error, like the user can't be logged in
		ctx.Logger.Infof("This is the token %d", token)
		w.Header().Set("Authorization", "Bearer "+strconv.Itoa(token))
		ctx.Logger.Infof("Authorization header set to: Bearer %d", token)
		ctx.Logger.Infof("Authorization header set to: Bearer %s", r.Header.Get("Authorization"))
		if r.Header.Get("Authorization") != ""{
			ctx.Logger.Infof("The auth header exists")
		} else {
			ctx.Logger.Infof("The auth header doesnt exists")
		}*/
	w.WriteHeader(http.StatusOK)
	ctx.Logger.Infof("The token is being encoded in the response body")
	_ = json.NewEncoder(w).Encode(token)
	return

}
