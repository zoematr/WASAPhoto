package api

import (
	"github.com/zoematr/WASAPhoto/service/api/reqcontext"
	"encoding/json"
	"net/http"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) GetStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// imposto il tipo di contenuto della risposta http in json
	// estraggo l'id del'utente dal  token Bearer nell'header di autorizzazione della richiesta HTTP.
	w.Header().Set("Content-Type", "application/json")
	token := extractToken(r.Header.Get("Authorization"))

	requestingUsername, err := rt.db.GetUsernameFromToken(token)
	if err != nil {
		ctx.Logger.WithError(err).Error("error retrieving username from token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Verifica se è l'utente stesso a vedere la propria home
	allowed := validateRequestingUser(token, ps.ByName("username"))
	if allowed != 0 {
		w.WriteHeader(allowed)
		return
	}

	photos, err := rt.db.GetStream(requestingUsername)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Imposta lo stato della risposta HTTP come 200 OK. Codifica l'elenco di foto in formato JSON e lo invia come corpo della risposta HTTP.
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(photos)
}