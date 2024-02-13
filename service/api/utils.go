package api

import (
	"net/http"
	"strconv"
	"strings"
)

// check correct length
func validUsername(username string) bool {
	// remove white spaces if not empty and does not contain "?" o "_".
	var trimmedusername = strings.TrimSpace(username)
	return len(trimmedusername) >= 3 && len(trimmedusername) <= 31 && trimmedusername != "" && !strings.ContainsAny(trimmedusername, "?_")
}

// Funzione che verifica se l'utente che effettua la richiesta ha un token valido per l'endpoint specificato.Restituisce 0 se è valido,o errore
func validateRequestingUser(dbToken int, auth string) int {

	// if the user is not logged he is not allowed to perform operation
	if isNotLogged(auth) {
		return http.StatusForbidden
	}
	bearerToken := extractToken(auth)
	// if token in the header is different from the token linked to the username in the path, the user requesting is not authorized
	if dbToken != bearerToken {
		return http.StatusUnauthorized
	}
	return 0
}

// gives true if authorization is "", i.e. user not logged in
func isNotLogged(auth string) bool {

	return auth == ""
}

func extractToken(authorization string) int {
	// Divide the authorization header in token utilizing space as divider
	var tokens = strings.Split(authorization, " ")
	// if there are exactly 2 token, give the second token (bearer token) after removing spaces if there are
	if len(tokens) == 2 {
		tokenstr := strings.TrimSpace(tokens[1])
		// Convert the token string to an integer
		token, err := strconv.Atoi(tokenstr)
		if err != nil {
			return 0
		}
		return token
	}
	// if there are not 2 spaces i give back empty strings
	return 0
}
