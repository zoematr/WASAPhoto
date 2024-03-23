package api

import (
	"github.com/julienschmidt/httprouter"
	"github.com/zoematr/WASAPhoto/service/api/reqcontext"
	"net/http"
)

// DELETE PHOTO
func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// get data from header and path
	w.Header().Set("Content-Type", "application/json")
	authToken := r.Header.Get("Authorization")
	pathUsername := ps.ByName("username")
	targetPhotoId := ps.ByName("photoid")
	targetCommentId := ps.ByName("commenid")
	tokenDbPath, err := rt.db.GetTokenFromUsername(pathUsername)
	// verify identity of the user
	valid := validateRequestingUser(tokenDbPath, authToken)
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}
	// check if the photo exists
	exists, err := rt.db.PhotoExists(targetPhotoId)
	if err != nil {
		ctx.Logger.WithError(err).Error("uncomment-photo: error checking author of the comment")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if exists != true {
		ctx.Logger.WithError(err).Error("uncomment-photo: the photo does not exist")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// check if the comment exists
	exists, err = rt.db.CommentExists(targetCommentId)
	if err != nil {
		ctx.Logger.WithError(err).Error("uncomment-photo: error checking author of the comment")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if exists != true {
		ctx.Logger.WithError(err).Error("uncomment-photo: the comment does not exist")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// check if the comment belongs to the user requesting the action
	usernameTarget, err := rt.db.GetUsernameFromCommentId(targetCommentId)
	if err != nil {
		ctx.Logger.WithError(err).Error("uncomment-photo: error checking author of the comment")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if usernameTarget != pathUsername {
		ctx.Logger.WithError(err).Error("uncomment-photo: you cannot remove someone else's comment")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// delete the comment from the db
	err = rt.db.DeleteComment(targetCommentId)
	if err != nil {
		ctx.Logger.WithError(err).Error("uncomment-photo: error removing the comment from the DB")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// send code 204
	w.WriteHeader(http.StatusNoContent)

}
