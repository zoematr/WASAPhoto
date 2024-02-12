// all the user DB functions

package database

// import (
// 	"fmt"
// )

func (db *appdbimpl) GetTokenFromUsername(username string) (int, error) {
    // fmt.Println("executing get token")
	row := db.c.QueryRow("SELECT token FROM users WHERE username = ?", username)
	// fmt.Println("this is the row")
	// fmt.Println(row)
	var token int
    err := row.Scan(&token)
	// fmt.Println("this is the db error")
	// fmt.Println(err)
	// fmt.Println("this is the token:")
	// fmt.Println(token)
    if err != nil {
        return 0, err
    }
    // Return the retrieved token
    return token, nil
}

func (db *appdbimpl) CreateUser(username string) (int, error) {
    // Insert the user into the database
	// fmt.Println("executing create user")
    _, err := db.c.Exec("INSERT INTO users (username) VALUES (?)", username)
	// fmt.Println("this is the db error")
	// fmt.Println(err)
    if err != nil {
        return 0, err // can be 0, sqlite autoincrement starts from 1
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

// get user stream
func (db *appdbimpl) GetStream(user User) ([]Photo, error) {
	rows, err := db.c.Query(`SELECT * FROM photos WHERE username IN (SELECT username FROM followers WHERE followerusername = ?) ORDER BY datetime DESC`,
		user.Username)
	if err != nil {
		return nil, err
	}
	// Wait for the func to finish before closing rows
	defer func() { _ = rows.Close() }()

	// Read all the users in the resulset
	var res []Photo
	for rows.Next() {
		var photo Photo
		err = rows.Scan(&photo.PhotoId, &photo.Username, &photo.Date) //  &photo.Comments, &photo.Likes,
		if err != nil {
			return nil, err
		}
		res = append(res, photo)
	}

	if rows.Err() != nil {
		return nil, err
	}
	// gives back slice of Photo which is the stream.
	return res, nil
}


// func (db *appdbimpl) CreateUser(username string) error {
// 	_, err := db.c.Exec("INSERT INTO users (username) VALUES (?)",
// 		username)
// 
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }





// checks if user exists if someone looks.
func (db *appdbimpl) ExistsUser(searcheduser string) (bool, error) {
	//  Esegue una query SQL per contare il numero di righe nella tabella degli utenti che corrispondono all'ID utente specificato.
	var cnt int
	err := db.c.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?",
		searcheduser).Scan(&cnt)

	if err != nil {
		// Count always returns a row thanks to COUNT(*), so this situation should not happen
		return true, err
	}

	// If counter 1 then the user exists
	if cnt == 1 {
		return true, nil
	}
	return false, nil
}

// gives nickname of a user given the token
func (db *appdbimpl) GetUsernameFromToken(token int) (string, error) {

	var username string

	// Utilizza una query SQL SELECT per cercare il nickname dell'utente nella tabella users utilizzando l'identificativo dell'utente (id_user).
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
	_, err := db.c.Exec(`UPDATE users SET username = ? WHERE token = ?`, newusername, token)
	if err != nil {
		// Error during the execution of the query
		return err
	}
	return nil
}

func (db *appdbimpl) CheckBanned(banner string, banned string) (bool, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM banned WHERE username = ? AND bannedusername = ?", banner, banned).Scan(&count)
    if err != nil {
        return false, err
    }
    // If count > 0, it means "username" has banned "bannedUsername"
	if count>0 {
		return true, nil
	}

    return false, nil
}

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

func (db *appdbimpl) GetPhotos(username string) ([]Photo, error) {
    var photos []Photo
    rows, err := db.c.Query(`SELECT * FROM photos WHERE username = ?`, username)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var photo Photo
        if err := rows.Scan(&photo.PhotoId, &photo.Username, &photo.PhotoFile, &photo.Date); err != nil {
            return nil, err
        }
        photos = append(photos, photo)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return photos, nil
}


