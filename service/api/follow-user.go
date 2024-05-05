package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/zoematr/WASAPhoto/service/api/reqcontext"
	"io/ioutil"
	"net/http"
)

// Funzione per mettere nella lista dei follow di un utente il follow di un'altro utente
func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	ctx.Logger.Infof("followUser is being called")
	w.Header().Set("Content-Type", "application/json")
	var usernameTargetUser string
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ctx.Logger.WithError(err).Error("set my username: error reading request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx.Logger.Infof("this is the request body: %s", string(body))

	var target UsernameReqBody
	err = json.Unmarshal(body, &target)
	if err != nil {
		ctx.Logger.WithError(err).Error("set my username: error decoding json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	usernameTargetUser = target.Username

	ctx.Logger.Infof("this is the usernameTargetUser of followUser %s", usernameTargetUser)
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
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// add follower in DB
	err = rt.db.FollowUser(usernameRequestUser, usernameTargetUser)
	if err != nil {
		ctx.Logger.WithError(err).Error("follow-user error in query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// get followers and following of the requested user
	followers, err := rt.db.GetFollowers(usernameTargetUser)
	if err != nil {
		ctx.Logger.WithError(err).Error("follow-user - GetFollowers: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	following, err := rt.db.GetFollowing(usernameTargetUser)
	if err != nil {
		ctx.Logger.WithError(err).Error("follow-user - GetFollowing: error executing query")
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
