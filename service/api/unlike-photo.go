package api

import (
	"github.com/julienschmidt/httprouter"
	"github.com/zoematr/WASAPhoto/service/api/reqcontext"
	"net/http"
)

// UNLIKE PHOTO
func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// retrieve data from body and path -> "/users/:username/photos/:photoid/likes/:likingusername"
	w.Header().Set("Content-Type", "application/json")
	authToken := r.Header.Get("Authorization")
	targetPhotoId := ps.ByName("photoid")
	likingUsername := ps.ByName("likingusername")

	// verify identity of the user, aka if the user is logged in
	likingToken, err := rt.db.GetTokenFromUsername(likingUsername)
	valid := validateRequestingUser(likingToken, authToken)
	if valid != 0 {
		ctx.Logger.WithError(err).Error("unlike-photo: user is not allowed to perform actions if not logged as the author of the action")
		w.WriteHeader(valid)
		return
	}

	// check if the photo exists
	exists, err := rt.db.PhotoExists(targetPhotoId)
	if !exists {
		ctx.Logger.WithError(err).Error("delete-photo: the photo does not exist")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	liked, err := rt.db.DoesUserLikePhoto(targetPhotoId, likingUsername)
	if err != nil {
		ctx.Logger.WithError(err).Error("unlike-photo: error executing db function call")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !liked {
		ctx.Logger.WithError(err).Error("delete-photo: you cannot unlike an unliked photo")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// delete the like
	err = rt.db.DeleteLike(targetPhotoId, likingUsername)
	if err != nil {
		ctx.Logger.WithError(err).Error("unlike-photo: error executing db function call")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// return status no content 204
	w.WriteHeader(http.StatusNoContent)

}
