package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/zoematr/WASAPhoto/service/api/reqcontext"
	"net/http"
)

// Funzione per mettere nella lista dei follow di un utente il follow di un'altro utente
func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	usernameTargetUser := ps.ByName("followingusername")
	_, err = rt.db.ExistsUser(usernameTargetUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
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

	// CHECK IF THE TARGET WAS FOLLOWED BY THE REQUESTING USER
	wasFollowed, err := rt.db.WasTargetFollowed(usernameRequestUser, usernameTargetUser)
	if err != nil {
		ctx.Logger.WithError(err).Error("error getting token requesting user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !wasFollowed {
		w.WriteHeader(http.StatusForbidden)
	}
	// UNFOLLOW
	_ = rt.db.UnfollowUser(usernameRequestUser, usernameTargetUser)
	// Respond with 204 http status
	w.WriteHeader(http.StatusNoContent)
}
