package api

import (
	"github.com/zoematr/WASAPhoto/service/api/reqcontext"
	"encoding/json"
	"net/http"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setMyUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	pathUsername := ps.ByName("username")

	// get the username from path and then get the token from the db because i did not manage to do it inside of validaterequestingUser
	tokenDbPath, err := rt.db.GetTokenFromUsername(pathUsername)
	if err != nil {
		ctx.Logger.WithError(err).Error("set my username: error retrieving token")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	authToken := r.Header.Get("Authorization")
	allowed := validateRequestingUser(tokenDbPath, authToken)
	if allowed != 0 {
		return
	}

	// Estraggo il nuovo nickname dal corpo della richiesta e decodifica del JSON
	var newusername string
	err = json.NewDecoder(r.Body).Decode(&newusername)
	// Se c'è un errore nella decodifica del JSON, si risponde con un codice di stato HTTP 400 (Bad Request).
	if err != nil {
		ctx.Logger.WithError(err).Error("set my username: error decoding json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token := extractToken(r.Header.Get("Authorization"))
	err = rt.db.ChangeUsername(token, newusername)
	if err != nil {
		ctx.Logger.WithError(err).Error("set my username: error executing update query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{"username": newusername})
	if err != nil {
		ctx.Logger.WithError(err).Error("set my username: error returning the new username")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}