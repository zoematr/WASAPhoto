package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/zoematr/WASAPhoto/service/api/reqcontext"
	"io/ioutil"
	"net/http"
)

func (rt *_router) setMyUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	pathUsername := ps.ByName("username")
	ctx.Logger.Infof("the backend function is being called for %s", pathUsername)
	ctx.Logger.Infof("this is the username from the path: %s", pathUsername)
	// get the username from path and then get the token from the db because i did not manage to do it inside of validaterequestingUser
	tokenDbPath, err := rt.db.GetTokenFromUsername(pathUsername)
	ctx.Logger.Infof("this is the token from the path: %v", tokenDbPath)
	authToken := r.Header.Get("Authorization")
	if err != nil {
		ctx.Logger.WithError(err).Error("set my username: error retrieving token")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx.Logger.Infof("this is the dbtoken %d", tokenDbPath)
	ctx.Logger.Infof("this is the auth string %s", authToken)
	ctx.Logger.Infof("this is the extraced token %d", extractToken(authToken))
	allowed := validateRequestingUser(tokenDbPath, authToken)
	if allowed != 0 {
		ctx.Logger.Infof("the user is not allowed to change the username")
		return
	}

	// get new username from request body
	ctx.Logger.Infof("the request body is about to be read")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ctx.Logger.WithError(err).Error("set my username: error reading request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx.Logger.Infof("this is the request body: %s", string(body))

	var newusername string
	err = json.Unmarshal(body, &newusername)
	if err != nil {
		ctx.Logger.WithError(err).Error("set my username: error decoding json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// change the username in the DB
	token := extractToken(r.Header.Get("Authorization"))
	err = rt.db.ChangeUsername(token, newusername)
	if err != nil {
		ctx.Logger.WithError(err).Error("set my username: error executing update query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ctx.Logger.Infof("this is the new username %s", newusername)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{"username": newusername})
	if err != nil {
		ctx.Logger.WithError(err).Error("set my username: error returning the new username")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
