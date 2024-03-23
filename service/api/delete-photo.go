package api

import (
	"github.com/julienschmidt/httprouter"
	"github.com/zoematr/WASAPhoto/service/api/reqcontext"
	"net/http"
)

// DELETE PHOTO
func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// get data from header and path
	// get the username from path and then get the token from the db because i did not manage to do it inside of validaterequestingUser
	w.Header().Set("Content-Type", "application/json")
	authToken := r.Header.Get("Authorization")
	pathUsername := ps.ByName("username")
	targetPhotoId := ps.ByName("photoid")
	tokenDbPath, err := rt.db.GetTokenFromUsername(pathUsername)
	// verify identity of the user
	valid := validateRequestingUser(tokenDbPath, authToken)
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}
	// check if the photo exists
	exists, err := rt.db.PhotoExists(targetPhotoId)
	if exists != true {
		ctx.Logger.WithError(err).Error("delete-photo: the photo does not exist")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// check if the photo belongs to the user requesting the action
	usernameTarget, err := rt.db.GetUsernameFromPhotoId(targetPhotoId)
	if err != nil {
		ctx.Logger.WithError(err).Error("delete-photo: error checking the author of the photo")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if usernameTarget != pathUsername {
		ctx.Logger.WithError(err).Error("delete-photo: you cannot remove someone else's photo")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// delete the photo from the db
	err = rt.db.DeletePhoto(targetPhotoId)
	if err != nil {
		ctx.Logger.WithError(err).Error("delete-photo: error removing the photo from the DB")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// send code 204
	w.WriteHeader(http.StatusNoContent)

}
