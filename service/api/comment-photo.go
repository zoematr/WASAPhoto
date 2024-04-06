package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/zoematr/WASAPhoto/service/api/reqcontext"
	"net/http"
)

func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")
	authToken := r.Header.Get("Authorization")
	pathRequestUsername := ps.ByName("username")
	targetPhotoId := ps.ByName("photoid")
	tokenDbPath, err := rt.db.GetTokenFromUsername(pathRequestUsername)
	// Verifica l'identità dell'utente che effettua la richiesta
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
	if len(commentContent) > 400 {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("comment is too long")
		return
	}

	comment := Comment{
		Username:       pathRequestUsername,
		PhotoId:        targetPhotoId,
		CommentContent: commentContent,
	}

	err = rt.db.AddComment(comment.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("comment photo: error adding comment to db")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// return 201 and comment
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(comment)

}
