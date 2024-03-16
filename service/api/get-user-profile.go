package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/zoematr/WASAPhoto/service/api/reqcontext"
	"github.com/zoematr/WASAPhoto/service/database"
	"net/http"
)

// Funzione che ritrova tutte le info necessarie del profilo
func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Estraggo l'id dell'utente che fa la richiesta e l'id dell'utente richiesto
	tokenUserRequesting := extractToken(r.Header.Get("Authorization"))
	userRequesting, err := rt.db.GetUsernameFromToken(tokenUserRequesting)
	targetUser := ps.ByName("username")

	var followers []string
	var following []string
	var photos []database.Photo

	// check if targetuser exists
	exists, err := rt.db.ExistsUser(targetUser)
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserProfile : error executing query")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if !exists {
		ctx.Logger.WithError(err).Error("getUserProfile : the user does not exist")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// check if the 2 users banned eachother
	userBanned, err := rt.db.CheckBanned(targetUser, userRequesting)
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserProfile - userBanned: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if userBanned {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	userBanned, err = rt.db.CheckBanned(userRequesting, targetUser)
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserProfile - userBanned: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if userBanned {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// get followers and following of the requested user
	followers, err = rt.db.GetFollowers(targetUser)
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserProfile - GetFollowers: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	following, err = rt.db.GetFollowing(targetUser)
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserProfilen - GetFollowing: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// get also list of photos of the user from db
	photos, err = rt.db.GetPhotos(targetUser)
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserProfile - GetPhotosList: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// send code 200 and returm the user profile
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(UserProfile{
		Username:  targetUser,
		Followers: followers,
		Following: following,
		Photos:    photos,
	})

}
