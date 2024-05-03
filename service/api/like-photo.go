package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/zoematr/WASAPhoto/service/api/reqcontext"
	"net/http"
)

// LIKE PHOTO
func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	ctx.Logger.Infof("like-photo is called")
	// PROBLEM SO FAR, THE AUTHSTRING IS EMPTY
	w.Header().Set("Content-Type", "application/json")
	var authToken string
	err := json.NewDecoder(r.Body).Decode(&authToken)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("like photo: failed to decode authorization string")
		return
	}
	ctx.Logger.Infof("this is authtoken %s", authToken)
	targetUsername := ps.ByName("username")
	targetPhotoId := ps.ByName("photoid")
	requestUsername, err := rt.db.GetUsernameFromToken(extractToken(authToken))
	if err != nil {
		ctx.Logger.WithError(err).Error("like photo: error authenticating user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ctx.Logger.Infof("this is requestUsername", requestUsername)
	banned, err := rt.db.CheckBanned(targetUsername, requestUsername)
	if err != nil {
		ctx.Logger.WithError(err).Error("like photo: error executing db function call")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if banned != false {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// check if the photo exists
	exists, err := rt.db.PhotoExists(targetPhotoId)
	if err != nil {
		ctx.Logger.WithError(err).Error("photo-like error: PhotoExists")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if exists != true {
		ctx.Logger.WithError(err).Error("like-photo: the photo does not exist")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	liked, err := rt.db.DoesUserLikePhoto(targetPhotoId, requestUsername)
	if err != nil {
		ctx.Logger.WithError(err).Error("photo-like error: DoesUserLikePhoto")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if liked != false {
		ctx.Logger.WithError(err).Error("like-photo: you already liked this photo")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err = rt.db.AddLike(targetPhotoId, requestUsername)
	if err != nil {
		ctx.Logger.WithError(err).Error("photo-like error: AddLike")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// send code 201 (like created) and send back the liked photo
	photo, err := rt.db.GetPhotoFromPhotoId(targetPhotoId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error retriving the photo from the DB")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	photo.AlreadyLiked = true
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(photo)
}
