package api

import (
	"net/http"
	"strings"
)

// check correct length
func validUsername(username string) bool {
	// remove white spaces if not empty and does not contain "?" o "_".
	var trimmedusername = strings.TrimSpace(username)
	return len(trimmedusername) >= 3 && len(trimmedusername) <= 31 && trimmedusername != "" && !strings.ContainsAny(trimmedusername, "?_")
}


// Funzione che verifica se l'utente che effettua la richiesta ha un token valido per l'endpoint specificato.Restituisce 0 se è valido,o errore
func validateRequestingUser(identifier string, bearerToken string) int {

	// Se l'utente che effettua la richiesta ha un token non valido, restituisci un codice di stato HTTP 403
	if isNotLogged(bearerToken) {
		return http.StatusForbidden
	}

	// Se l'ID dell'utente che effettua la richiesta è diverso da quello nel percorso, restituisci un codice di stato HTTP 401
	if identifier != bearerToken {
		return http.StatusUnauthorized
	}
	return 0
}

// funzione che verifica se un utente è loggato.
// Restituisci true se la stringa di autenticazione è vuota (cioè l'utente non è loggato), altrimenti restituisci false
func isNotLogged(auth string) bool {

	return auth == ""
}

func extractBearer(authorization string) string {
	// Divido l'intestazione di autorizzazione in token utilizzando lo spazio come delimitatore.
	var tokens = strings.Split(authorization, " ")
	// Se ci sono esattamente due token, restituisci il secondo token (il token Bearer) dopo aver rimosso eventuali spazi.
	if len(tokens) == 2 {
		return strings.Trim(tokens[1], " ")
	}
	// Se non ci sono 2 token restituisco stringa vuota
	return ""
}