package api

import (
	"github.com/zoematr/WASAPhoto/service/api/reqcontext"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"encoding/json"
)


func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	usernameRequestUser := ps.ByName("username")
	tokenRequestUser, err := rt.db.GetTokenFromUsername(usernameRequestUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bearerToken := r.Header.Get("Authorization")
	w.Header().Set("Content-Type", "application/json")
	var usernameTargetUser string
	err := json.NewDecoder(r.Body).Decode(&usernameTargetUser)
	_, err = rt.db.ExistsUser(usernameTargetUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Controlla se l'utente che effettua la richiesta ha il permesso di bannare l'utente specificato. 
	// (solo l'owner dell'account puo aggiugnere un banned user nel suo account list)
	allow := validateRequestingUser(tokenRequestUser, bearerToken)
	if valid != 0 {
		w.WriteHeader(allow)
		return
	}

	// Controlla se l'utente sta cercando di bannare se stesso.se si(400:"Bad Request")
	if usernameRequestUser == usernameTargetUser {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Chiama una funzione del database per aggiungere l'utente specificato alla lista degli utenti bannati.
	err := rt.db.BanUser(usernameRequestUser, usernameTargetUser)
	if err != nil {
		ctx.Logger.WithError(err).Error("ban user/db.BanUser: error executing insert query")

		// C'è stato un errore interno,restituisco(error:500)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Respond with 204 http status
	w.WriteHeader(http.StatusNoContent)
}