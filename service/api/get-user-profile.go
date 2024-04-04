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
	authString := (r.Header.Get("Authorization"))
	if isNotLogged(authString) {
		ctx.Logger.Infof("get-user-profile: log in to see the profiles of other users!")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	tokenUserRequesting := extractToken(authString)
	userRequesting, err := rt.db.GetUsernameFromToken(tokenUserRequesting)

	targetUser := ps.ByName("username")

	var followers []string
	var following []string
	var photos []database.Photo
	alreadyfollowing := false // used later to check if the requesting user already follows the target, then he can't be able to refollow him but just to unfollow
	alreadybanned := false    // same as above but for banned
	owner := false            // if the user that is searched is also the owner of the profile, then they can't follow or ban themselves

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

	// check if the requesting is banned, then he can't see the profile users banned eachother
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
		alreadybanned = true
	}

	// get followers and following of the requested user
	followers, err = rt.db.GetFollowers(targetUser)
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserProfile - GetFollowers: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, follower := range followers {
		if follower == userRequesting {
			alreadyfollowing = true
			break
		}
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

	if userRequesting == targetUser {
		owner = true
	}

	// send code 200 and returm the user profile
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(UserProfile{
		Username:        targetUser,
		Followers:       followers,
		Following:       following,
		Photos:          photos,
		AlreadyBanned:   alreadybanned,
		AlreadyFollowed: alreadyfollowing,
		OwnProfile:      owner,
	})

}
