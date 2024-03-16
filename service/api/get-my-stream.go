package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/zoematr/WASAPhoto/service/api/reqcontext"
	"net/http"
)

func (rt *_router) GetStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// verify user identity
	w.Header().Set("Content-Type", "application/json")
	authToken := r.Header.Get("Authorization")
	allowedUsername := ps.ByName("username")
	allowedToken, err := rt.db.GetTokenFromUsername(allowedUsername)
	if err != nil {
		ctx.Logger.WithError(err).Error("error retrieving username from token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	allowed := validateRequestingUser(allowedToken, authToken)
	if allowed != 0 {
		w.WriteHeader(allowed)
		return
	}

	photos, err := rt.db.GetStream(allowedUsername)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// respond 200 and give back the list of photos
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(photos)
}
