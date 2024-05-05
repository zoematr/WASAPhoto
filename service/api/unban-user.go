package api

import (
	"github.com/julienschmidt/httprouter"
	"github.com/zoematr/WASAPhoto/service/api/reqcontext"
	"net/http"
)

// Funzione per mettere nella lista dei follow di un utente il follow di un'altro utente
func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	usernameTargetUser := ps.ByName("bannedusername")
	exists, err := rt.db.ExistsUser(usernameTargetUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !exists {
		ctx.Logger.WithError(err).Error("delete-photo: the photo does not exist")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	usernameRequestUser := ps.ByName("username")
	tokenRequestUser, err := rt.db.GetTokenFromUsername(usernameRequestUser)
	if err != nil {
		ctx.Logger.WithError(err).Error("error getting token requesting user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	bearerToken := extractToken(r.Header.Get("Authorization"))
	if tokenRequestUser != bearerToken { // not person logged in
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// CHECK IF THE TARGET WAS BANNED ALREADY BY THE REQUESTING USER
	wasBanned, err := rt.db.WasTargetBanned(usernameRequestUser, usernameTargetUser)
	if err != nil {
		ctx.Logger.WithError(err).Error("error getting token requesting user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !wasBanned {
		w.WriteHeader(http.StatusForbidden)
	}
	// UNBAN
	_ = rt.db.UnbanUser(usernameRequestUser, usernameTargetUser)
	// Respond with 204 http status
	w.WriteHeader(http.StatusNoContent)
}
