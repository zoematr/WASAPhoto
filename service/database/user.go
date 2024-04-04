// all the user DB functions

package database

// function that gets token linked to a username from the DB
func (db *appdbimpl) GetTokenFromUsername(username string) (int, error) {

	row := db.c.QueryRow("SELECT token FROM users WHERE username = ?", username)
	var token int
	err := row.Scan(&token)
	if err != nil {
		return 0, err
	}
	// Return the token
	return token, err
}

// creates a new user and puts it in the database
func (db *appdbimpl) CreateUser(username string) (int, error) {
	// Insert the user into the database
	_, err := db.c.Exec("INSERT INTO users (username) VALUES (?)", username)
	if err != nil {
		return 0, err // can't be 0, sqlite autoincrement starts from 1
	}
	// Retrieve the token for the inserted user
	token, err := db.GetTokenFromUsername(username)
	// fmt.Println(err)
	if err != nil {
		return 0, err
	}
	// Return the retrieved token
	return token, nil
}

// checks if user exists if someone looks.
func (db *appdbimpl) ExistsUser(searcheduser string) (bool, error) {
	var cnt int
	err := db.c.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?",
		searcheduser).Scan(&cnt)

	if err != nil {
		// Count always returns a row thanks to COUNT(*), so this situation should not happen
		return false, err
	}

	// If counter 1 then the user exists
	if cnt == 1 {
		return true, nil
	}
	return false, nil
}

// from the token, gets username in DB
func (db *appdbimpl) GetUsernameFromToken(token int) (string, error) {
	var username string
	err := db.c.QueryRow(`SELECT username FROM users WHERE token = ?`, token).Scan(&username)
	if err != nil {
		// Error during the execution of the query
		return "", err
	}
	return username, nil
}

// function to change username.
func (db *appdbimpl) ChangeUsername(token int, newusername string) error {
	// query update using the token
	// doesn't have to be changed in every table because it's a foreign key
	_, err := db.c.Exec(`UPDATE users SET username = ? WHERE token = ?`, newusername, token)
	if err != nil {
		// Error during the execution of the query
		return err
	}
	return nil
}

// checks if a banner banned banned
func (db *appdbimpl) CheckBanned(banner string, banned string) (bool, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM banned WHERE username = ? AND bannedusername = ?", banner, banned).Scan(&count)
	if err != nil {
		return false, err
	}
	// If count > 0, it means "username" has banned "bannedUsername"
	if count > 0 {
		return true, nil
	}

	return false, nil
}

// function that returns all the followers of a follower
func (db *appdbimpl) GetFollowers(followed string) ([]string, error) {
	var followers []string
	rows, err := db.c.Query(`SELECT followerusername FROM followers WHERE username = ?`, followed)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var follower string
		if err := rows.Scan(&follower); err != nil {
			return nil, err
		}
		followers = append(followers, follower)
	}
	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return followers, nil
}

// function that returns all the people followed by a user
func (db *appdbimpl) GetFollowing(follower string) ([]string, error) {
	var following []string
	rows, err := db.c.Query(`SELECT username FROM followers WHERE followerusername = ?`, follower)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var follow string
		if err := rows.Scan(&follow); err != nil {
			return nil, err
		}
		following = append(following, follow)
	}
	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return following, nil
}

// function to start follow a user
func (db *appdbimpl) FollowUser(requesting string, target string) error {

	_, err := db.c.Exec("INSERT INTO followers (username, followerusername) VALUES (?,?)", target, requesting)
	if err != nil {
		return err
	}
	// Return the username of the user followed
	return nil
}

// function to check if a user was already followed
func (db *appdbimpl) WasTargetFollowed(requesting string, target string) (bool, error) {
	var cnt int
	err := db.c.QueryRow("SELECT COUNT(*) FROM followers WHERE username = ? AND followerusername = ?", target, requesting).Scan(&cnt)

	if err != nil {
		// Count always returns a row thanks to COUNT(*), so this situation should not happen
		return true, err
	}
	// If counter 1 then the target was followed
	if cnt == 1 {
		return true, nil
	}
	return false, nil
}

// check if a user is banned by another user
func (db *appdbimpl) WasTargetBanned(requesting string, target string) (bool, error) {
	var cnt int
	err := db.c.QueryRow("SELECT COUNT(*) FROM banned WHERE username = ? AND bannedusername = ?", requesting, target).Scan(&cnt)

	if err != nil {
		// Count always returns a row thanks to COUNT(*), so this situation should not happen
		return true, err
	}
	// If counter 1 then the target was followed
	if cnt == 1 {
		return true, nil
	}
	return false, nil
}

// function to unfollow a user
func (db *appdbimpl) UnfollowUser(requesting string, target string) error {
	_, err := db.c.Exec("DELETE FROM followers WHERE username = ? AND followerusername = ?", target, requesting)
	if err != nil {
		return err
	}
	return nil
}

// function to unban a user
func (db *appdbimpl) UnbanUser(requesting string, target string) error {
	_, err := db.c.Exec("DELETE FROM banned WHERE username = ? AND bannedusername = ?", requesting, target)
	if err != nil {
		return err
	}
	return nil
}

// function to ban a user
func (db *appdbimpl) BanUser(requesting string, target string) error {
	wasbanned, err := db.WasTargetBanned(requesting, target)
	if err != nil {
		return err
	}
	if wasbanned == true {
		return err
	}
	wasfollowed, err := db.WasTargetFollowed(requesting, target)
	if err != nil {
		return err
	}
	if wasfollowed == true {
		err = db.UnfollowUser(requesting, target)
		if err != nil {
			return err
		}
	}
	wasfollowed, err = db.WasTargetFollowed(target, requesting)
	if err != nil {
		return err
	}
	if wasfollowed == true {
		err = db.UnfollowUser(target, requesting)
		if err != nil {
			return err
		}
	}
	_, err = db.c.Exec("INSERT INTO banned (username, bannedusername) VALUES (?,?)", requesting, target)
	if err != nil {
		return err
	}
	return nil
}
