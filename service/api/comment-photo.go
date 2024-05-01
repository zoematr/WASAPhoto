package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/zoematr/WASAPhoto/service/api/reqcontext"
	"net/http"
)

func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	ctx.Logger.Infof("Comment-photo is being called in the backend")
	authToken := r.Header.Get("Authorization")
	targetUsername := ps.ByName("username")
	targetPhotoId := ps.ByName("photoid")
	requestUsername, err := rt.db.GetUsernameFromToken(extractToken(authToken))
	// Verifica l'identità dell'utente che effettua la richiesta
	ctx.Logger.Infof("this is authToken", authToken)
	ctx.Logger.Infof("this is the comment requestingUser %s", requestUsername)

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
		ctx.Logger.WithError(err).Error("photo-like error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if exists != true {
		ctx.Logger.WithError(err).Error("delete-photo: the photo does not exist")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var commentContent string
	err = json.NewDecoder(r.Body).Decode(&commentContent)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("comment photo: failed to decode request body json")
		return
	}
	// Controllo la lunghezza del comment(<=400)
	if len(commentContent) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("comment is too short")
		return
	}
	if len(commentContent) > 400 {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("comment is too long")
		return
	}

	comment := Comment{
		Username:       requestUsername,
		PhotoId:        targetPhotoId,
		CommentContent: commentContent,
	}
	ctx.Logger.Infof("this is the content of the comment", commentContent)
	commentid, err := rt.db.AddComment(comment.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("comment photo: error adding comment to db")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	comment.CommentId = commentid
	ctx.Logger.Infof("this is the commentid", comment.CommentId)
	// return 201 and comment
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(comment)

}
