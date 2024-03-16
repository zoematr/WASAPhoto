package api

import (
	"github.com/zoematr/WASAPhoto/service/api/reqcontext"
	"encoding/json"
	"net/http"
	"github.com/julienschmidt/httprouter"
)

// LIKE PHOTO
func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")
	authToken := r.Header.Get("Authorization")
	pathRequestUsername := ps.ByName("username")
	targetPhotoId := ps.ByName("photoid")
	tokenDbPath, err := rt.db.GetTokenFromUsername(pathRequestUsername)
	valid := validateRequestingUser(tokenDbPath, authToken)
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}
	targetUsername, err := rt.db.GetUsernameFromPhotoId(targetPhotoId)
	if err != nil {
		ctx.Logger.WithError(err).Error("like photo: error executing db function call")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	banned, err := rt.db.CheckBanned(targetUsername, pathRequestUsername)
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
	if exists != true {
		ctx.Logger.WithError(err).Error("delete-photo: the photo does not exist")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	liked, err := rt.db.DoesUserLikePhoto(targetPhotoId, pathRequestUsername)
	if liked != false {
		ctx.Logger.WithError(err).Error("delete-photo: you already liked this photo")
		w.WriteHeader(http.StatusForbidden)
		return
	}


	err = rt.db.AddLike(targetPhotoId, pathRequestUsername)
	if err != nil {
		ctx.Logger.WithError(err).Error("photo-like error")
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
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(photo)
}
