package api

import (
	"github.com/zoematr/WASAPhoto/service/api/reqcontext"
	"encoding/json"
	"net/http"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) SetMyUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	pathUsername := ps.ByName("username")

	// get the username from path
	// Verifica l'identità dell'utente per l'operazione, confrontando l'ID dell'utente con l'ID dell'utente nel token Bearer.
	valid := validateRequestingUser(pathUsername, extractBearer(r.Header.Get("Authorization")))
	if valid != 0 {
		return
	}

	// Estraggo il nuovo nickname dal corpo della richiesta e decodifica del JSON
	var newusername string
	err := json.NewDecoder(r.Body).Decode(&newusername)
	// Se c'è un errore nella decodifica del JSON, si risponde con un codice di stato HTTP 400 (Bad Request).
	if err != nil {
		ctx.Logger.WithError(err).Error("changeusername: error decoding json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token := extractBearer(r.Header.Get("Authorization"))
	err = rt.db.ChangeUsername(token, newusername)
	if err != nil {
		ctx.Logger.WithError(err).Error("update-username: error executing update query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Respond with 204 http status
	w.WriteHeader(http.StatusNoContent)
}