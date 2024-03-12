package api

import (
	"github.com/zoematr/WASAPhoto/service/api/reqcontext"
	"encoding/json"
	"net/http"
	"github.com/julienschmidt/httprouter"
)

// Funzione che gestisce l'upload di una foto
func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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
	if exists != true {
		ctx.Logger.WithError(err).Error("delete-photo: the photo does not exist")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	liked, err := rt.db.DoesUserLikePhoto(targetPhotoId, pathRequestUsername)
	if liked != false {
		ctx.Logger.WithError(err).Error("delete-photo: you already liked this photo")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}


	err = rt.db.AddLike(targetPhotoId, pathRequestUsername)
	if err != nil {
		ctx.Logger.WithError(err).Error("photo-like error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Invia una risposta con stato "Created" e un oggetto JSON che rappresenta la foto appena caricata.
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(pathRequestUsername)

}
