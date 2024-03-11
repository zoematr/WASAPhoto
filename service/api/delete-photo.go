package api

import (
	"github.com/zoematr/WASAPhoto/service/api/reqcontext"
	"encoding/json"
	"net/http"
	"github.com/julienschmidt/httprouter"
)

// Funzione che gestisce l'upload di una foto
func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")
	authToken := r.Header.Get("Authorization")
	pathUsername := ps.ByName("username")
	targetPhotoId := ps.ByName("photoid")
	// get the username from path and then get the token from the db because i did not manage to do it inside of validaterequestingUser
	tokenDbPath, err := rt.db.GetTokenFromUsername(pathUsername)
	// Verifica l'identità dell'utente che effettua la richiesta
	valid := validateRequestingUser(tokenDbPath, authToken)
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}

	err = rt.db.DeletePhoto(targetPhotoId)
	if err != nil {
		ctx.Logger.WithError(err).Error("photo-delete: error executing db function call")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Invia una risposta con stato "Created" e un oggetto JSON che rappresenta la foto appena caricata.
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(pathUsername)

}
