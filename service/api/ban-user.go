package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/zoematr/WASAPhoto/service/api/reqcontext"
	"io/ioutil"
	"net/http"
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
	var alreadybanned bool
	var alreadyfollowing bool
	var owner bool
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ctx.Logger.WithError(err).Error("ban user: error reading request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx.Logger.Infof("this is the request body: %s", string(body))

	var target UsernameReqBody
	err = json.Unmarshal(body, &target)
	if err != nil {
		ctx.Logger.WithError(err).Error("ban user: error decoding json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	usernameTargetUser = target.Username
	exists, err := rt.db.ExistsUser(usernameTargetUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !exists {
		ctx.Logger.WithError(err).Error("ban user: the photo does not exist")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	allow := validateRequestingUser(tokenRequestUser, bearerToken)
	if allow != 0 {
		w.WriteHeader(allow)
		return
	}
	if usernameRequestUser == usernameTargetUser {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	ctx.Logger.Infof("checks done, user not banned yet")
	err = rt.db.BanUser(usernameRequestUser, usernameTargetUser)
	if err != nil {
		ctx.Logger.WithError(err).Error("ban user, db.BanUser: error executing insert query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	alreadybanned = true

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

	wasFollowed, err := rt.db.WasTargetFollowed(usernameRequestUser, usernameTargetUser)
	if err != nil {
		ctx.Logger.WithError(err).Error("error checking if the requesting user follows the target")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if wasFollowed {
		alreadyfollowing = true
	}

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
