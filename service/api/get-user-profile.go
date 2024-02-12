package api

import (
	"github.com/zoematr/WASAPhoto/service/api/reqcontext"
	"github.com/zoematr/WASAPhoto/service/database"
	"encoding/json"
	"net/http"
	"github.com/julienschmidt/httprouter"
)

// Funzione che ritrova tutte le info necessarie del profilo
func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Estraggo l'id dell'utente che fa la richiesta e l'id dell'utente richiesto
	tokenUserRequesting := extractToken(r.Header.Get("Authorization"))
	userRequesting, err := rt.db.GetUsernameFromToken(tokenUserRequesting)
	targetUser := ps.ByName("username")

	var followers []string
	var following []string
	var photos []database.Photo

	// check if targetuser exists
	exists, err := rt.db.ExistsUser(targetUser)
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserProfile : error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !exists {
		ctx.Logger.WithError(err).Error("getUserProfile : the user does not exist")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// check if the 2 users banned eachother
	userBanned, err := rt.db.CheckBanned(targetUser, userRequesting)
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserProfile - userBanned: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if userBanned {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	userBanned, err = rt.db.CheckBanned(userRequesting, targetUser)
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserProfile - userBanned: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if userBanned {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Recupero la lista dei followers dell'utente richiesto
	followers, err = rt.db.GetFollowers(targetUser)
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserProfile - GetFollowers: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Recupero la lista dei utenti seguiti dall'utente richiesto
	following, err = rt.db.GetFollowing(targetUser)
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserProfilen - GetFollowing: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Recupera la lista delle foto dell'utente richiesto dal database.
	photos, err = rt.db.GetPhotos(targetUser)
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserProfile - GetPhotosList: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Imposta il codice di stato della risposta HTTP come 200 (OK) e invia un oggetto
	// JSON che rappresenta il profilo completo dell'utente richiesto come corpo della risposta.
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(UserProfile{
		Username:  targetUser,
		Followers: followers,
		Following: following,
		Photos:    photos,
	})

}