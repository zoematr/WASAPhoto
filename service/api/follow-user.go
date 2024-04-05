package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/zoematr/WASAPhoto/service/api/reqcontext"
	"net/http"
)

// Funzione per mettere nella lista dei follow di un utente il follow di un'altro utente
func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")
	var usernameTargetUser string
	err := json.NewDecoder(r.Body).Decode(&usernameTargetUser)
	exists, err := rt.db.ExistsUser(usernameTargetUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if exists != true {
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

	if usernameRequestUser == usernameTargetUser { // user can't follow himself
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	alreadybanned := false // used to give back to the frontend
	owner := false         // if the user that is searched is also the owner of the profile, then they can't follow or ban themselves

	// check if banned
	banned, err := rt.db.CheckBanned(usernameTargetUser, usernameRequestUser)
	if err != nil {
		ctx.Logger.WithError(err).Error("post-comment/rt.db.BannedUserCheck: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if banned {
		// User was banned, can't perform the follow action
		w.WriteHeader(http.StatusForbidden)
		return
	}

	banned, err = rt.db.CheckBanned(usernameRequestUser, usernameTargetUser)
	if err != nil {
		ctx.Logger.WithError(err).Error("follow-user/rt.db.BannedUserCheck: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if banned {
		alreadybanned = true
	}
	var alreadyfollowing bool
	wasFollowed, err := rt.db.WasTargetFollowed(usernameRequestUser, usernameTargetUser)
	if err != nil {
		ctx.Logger.WithError(err).Error("error getting token requesting user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if wasFollowed {
		ctx.Logger.Infof("follow-user: you cannot follow a user that you already follow!")
		return
	}

	// add follower in DB
	err = rt.db.FollowUser(usernameRequestUser, usernameTargetUser)
	if err != nil {
		ctx.Logger.WithError(err).Error("follow user error in query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// get followers and following of the requested user
	followers, err := rt.db.GetFollowers(usernameTargetUser)
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserProfile - GetFollowers: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	following, err := rt.db.GetFollowing(usernameTargetUser)
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserProfilen - GetFollowing: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// get also list of photos of the user from db
	photos, err := rt.db.GetPhotos(usernameTargetUser, usernameRequestUser)
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserProfile - GetPhotosList: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	alreadyfollowing = true

	// Respond with 201 http status
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(UserProfile{
		Username:        usernameTargetUser,
		Followers:       followers,
		Following:       following,
		Photos:          photos,
		AlreadyFollowed: alreadyfollowing,
		AlreadyBanned:   alreadybanned,
		OwnProfile:      owner,
	})

}
