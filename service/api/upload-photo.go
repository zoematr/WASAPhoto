package api

import (
	"encoding/base64"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/zoematr/WASAPhoto/service/api/reqcontext"
	"net/http"
	"time"
)

// Funzione che gestisce l'upload di una foto
func (rt *_router) uplaodPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")
	authToken := r.Header.Get("Authorization")
	pathUsername := ps.ByName("username")
	// get the username from path and then get the token from the db because i did not manage to do it inside of validaterequestingUser
	tokenDbPath, err := rt.db.GetTokenFromUsername(pathUsername)
	// Verifica l'identità dell'utente che effettua la richiesta
	valid := validateRequestingUser(tokenDbPath, authToken)
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}
	var photoInput PhotoInput
	err = json.NewDecoder(r.Body).Decode(&photoInput)
	if err != nil {
		ctx.Logger.WithError(err).Error("photo-upload: error decoding json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Decode the base64-encoded photo data into binary
	photoFile, err := base64.StdEncoding.DecodeString(photoInput.PhotoFile)
	if err != nil {
		ctx.Logger.WithError(err).Error("photo-upload: error decoding json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	photo := Photo{
		Username:  pathUsername,
		Date:      time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		PhotoFile: photoFile,
	}

	err = rt.db.AddPhoto(photo.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("photo-upload: error executing db function call")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Invia una risposta con stato "Created" e un oggetto JSON che rappresenta la foto appena caricata.
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(photo)

}
