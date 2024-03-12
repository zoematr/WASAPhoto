package api

import (
	"github.com/zoematr/WASAPhoto/service/api/reqcontext"
	"net/http"
	"time"
	"github.com/julienschmidt/httprouter"
	"encoding/json"
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
		Username: pathRequestUsername,
		Date:  time.Now().UTC(),
		PhotoId: targetPhotoId,
		CommentContent: commentContent,
	}

	err = rt.db.AddComment(comment.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("comment photo: error addind comment to db")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

