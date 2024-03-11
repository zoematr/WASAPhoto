package api

import (
	"github.com/zoematr/WASAPhoto/service/api/reqcontext"
	"bytes"
	"encoding/json"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"time"
	"github.com/julienschmidt/httprouter"
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

	// Reimposta il body della richiesta in modo da poterlo leggere di nuovo in seguito
	// Dopo aver letto il body bisogna riassegnare un io.ReadCloser per poterlo rileggere
	r.Body = io.NopCloser(bytes.NewBuffer(photoFile))

	// verifico se il contenuto del body è una immagine png o jpeg (in caso di errore:400 badrequest)
	err = checkFormatPhoto(r.Body, io.NopCloser(bytes.NewBuffer(photoFile)), ctx)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("photo-upload: body contains file that is neither jpg or png")
		// controllaerrore
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	// Reimposta nuovamente il corpo della richiesta per poterlo leggere di nuovo.
	r.Body = io.NopCloser(bytes.NewBuffer(photoFile))

	photo := Photo{
		Username: pathUsername,
		Date:  time.Now().UTC(),
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
	_ = json.NewEncoder(w).Encode(pathUsername)

}

// Funzione per controllare se il formato della foto è png o jpeg.Ritorno l'estenzione del formato e un errore 
func checkFormatPhoto(body io.ReadCloser, newReader io.ReadCloser, ctx reqcontext.RequestContext) error {

	_, errJpg := jpeg.Decode(body)
	if errJpg != nil {

		body = newReader
		_, errPng := png.Decode(body)
		if errPng != nil {
			return errPng
		}
		return nil
	}
	return nil
}
