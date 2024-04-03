package api

import (
	"bytes"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/zoematr/WASAPhoto/service/api/reqcontext"
	"io"
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

	// Legge il body della richiesta e verifica se ci sono errori durante la lettura.
	photoFile, err := io.ReadAll(r.Body)
	if err != nil {
		ctx.Logger.WithError(err).Error("photo-upload: error reading image data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	r.Body = io.NopCloser(bytes.NewBuffer(photoFile))

	photo := Photo{
		Username:  pathUsername,
		Date:      time.Now().UTC(),
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

