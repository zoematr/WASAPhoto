package api

import (
	"github.com/zoematr/WASAPhoto/service/api/reqcontext"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"encoding/json"
)


func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	ctx.Logger.Infof("start")
	usernameRequestUser := ps.ByName("username")
	tokenRequestUser, err := rt.db.GetTokenFromUsername(usernameRequestUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx.Logger.Infof("found user from the path")
	bearerToken := r.Header.Get("Authorization")
	w.Header().Set("Content-Type", "application/json")
	var usernameTargetUser string
	err = json.NewDecoder(r.Body).Decode(&usernameTargetUser)
	_, err = rt.db.ExistsUser(usernameTargetUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	allow := validateRequestingUser(tokenRequestUser, bearerToken)
	if allow != 0 {
		w.WriteHeader(allow)
		return
	}

	if usernameRequestUser == usernameTargetUser {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	ctx.Logger.Infof("checks done, user not banned yet")
	err = rt.db.BanUser(usernameRequestUser, usernameTargetUser)
	if err != nil {
		ctx.Logger.WithError(err).Error("ban user/db.BanUser: error executing insert query")

		// C'è stato un errore interno,restituisco(error:500)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Respond with 204 http status
	w.WriteHeader(http.StatusNoContent)
}