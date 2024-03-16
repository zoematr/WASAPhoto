package api

import (
	"github.com/zoematr/WASAPhoto/service/api/reqcontext"
	"net/http"
	"github.com/julienschmidt/httprouter"
)

// UNLIKE PHOTO
func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// retrieve data from header and path -> "/users/:username/photos/:photoid/likes/:likingusername"
	w.Header().Set("Content-Type", "application/json")
	authToken := r.Header.Get("Authorization")
	pathUsername := ps.ByName("username")
	targetPhotoId := ps.ByName("photoid")
	targetLikingUsername := ps.ByName("likingusername")

	// verify identity of the user, aka if the user is logged in
	tokenDbPath, err := rt.db.GetTokenFromUsername(pathUsername)
	valid := validateRequestingUser(tokenDbPath, authToken)
	if valid != 0 {
		ctx.Logger.WithError(err).Error("unlike-photo: user is not allowed to perform actions if not logged as the author of the action")
		w.WriteHeader(valid)
		return
	}

	// check if the user is trying to remove someome else's like
	if targetLikingUsername != pathUsername {
		ctx.Logger.WithError(err).Error("unlike-photo: you can't remove someone else's like")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// check if the photo exists
	exists, err := rt.db.PhotoExists(targetPhotoId)
	if exists != true {
		ctx.Logger.WithError(err).Error("delete-photo: the photo does not exist")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	liked, err := rt.db.DoesUserLikePhoto(targetPhotoId, pathUsername)
	if err != nil {
		ctx.Logger.WithError(err).Error("unlike-photo: error executing db function call")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if liked != true {
		ctx.Logger.WithError(err).Error("delete-photo: you cannot unlike an unliked photo")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// delete the like
	err = rt.db.DeleteLike(targetPhotoId, targetLikingUsername)
	if err != nil {
		ctx.Logger.WithError(err).Error("unlike-photo: error executing db function call")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// return status no content 204
	w.WriteHeader(http.StatusNoContent)

}
