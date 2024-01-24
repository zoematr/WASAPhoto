package api

import (
	"encoding/json"
	"net/http"

	"github.com/zoematr/WASAPhoto/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

// Login handles the user login or registration.
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Parse the request body to get the username.
	var usernameReq Username
	if err := json.NewDecoder(r.Body).Decode(&usernameReq); err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	// Check if the user with the given username already exists.
	user, err := rt.db.GetUserByUsername(usernameReq.Username)
	if err != nil {
		// Handle the error (e.g., database error).
		http.Error(w, "Error checking user existence", http.StatusInternalServerError)
		return
	}

	if user == nil {
		// If the user does not exist, register a new user.
		newUser := User{
			Username: usernameReq.Username,
			// Generate a unique user ID (you might have your own logic for this).
			UserId: generateUser(),
		}

		// Save the new user to the database.
		if err := rt.db.RegisterUser(newUser); err != nil {
			// Handle the error (e.g., database error).
			http.Error(w, "Error registering new user", http.StatusInternalServerError)
			return
		}

		user = &newUser
	}

	// Respond with the user identifier.
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(UserId{UserId: user.UserId})
}

// You need a function to generate a unique user ID.
// You might use a library like UUID or any other method based on your requirements.
func generateUniqueUserID() string {
	// Implement your logic to generate a unique user ID.
	// Example: return uuid.New().String()
	return "example_user_id"
}
