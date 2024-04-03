package api

import (
	"net/http"
	"strconv"
	"strings"
)

// check correct length of a username
func validUsername(username string) bool {
	// remove white spaces if not empty and does not contain "?" o "_".
	var trimmedusername = strings.TrimSpace(username)
	return len(trimmedusername) >= 3 && len(trimmedusername) <= 31 && trimmedusername != "" && !strings.ContainsAny(trimmedusername, "?_")
}

// checks if user requesting is the one logged in
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

// checks if logged in, gives true if authorization is "", i.e. user not logged in
func isNotLogged(auth string) bool {

	return auth == ""
}

// function that extracts the token (int) from authorization header
func extractToken(authorization string) int {
	// Divide the authorization header in token utilizing space as divider
	var tokens = strings.Split(authorization, " ")
	// if there is exactly 1 token, return it as the bearer token
	if len(tokens) == 1 {
		token, err := strconv.Atoi(tokens[0])
		if err != nil {
			return 0
		}
		return token
	}
	// if there are exactly 2 tokens, give the second token (bearer token) after removing spaces if there are
	if len(tokens) == 2 {
		tokenstr := strings.TrimSpace(tokens[1])
		// Convert the token string to an integer
		token, err := strconv.Atoi(tokenstr)
		if err != nil {
			return 0
		}
		return token
	}
	// if there are not 1 or 2 tokens, return 0
	return 0
}
