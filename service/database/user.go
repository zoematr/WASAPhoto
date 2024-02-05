// all the user DB functions

package database

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
// 
// 	return nil
// }

//insert user in DB
func (db *appdbimpl) CreateUser(username string) (int, error) {
    // Insert the user into the database
    _, err := db.c.Exec("INSERT INTO users (username) VALUES (?)", username)
    if err != nil {
        return 0, err // can be 0, sqlite autoincrement starts from 1
    }

    // Retrieve the token for the inserted user
    row := db.c.QueryRow("SELECT token FROM users WHERE username = ?", username)
    var token int
    err = row.Scan(&token)
    if err != nil {
        return 0, err
    }

    // Return the retrieved token
    return token, nil
}

func (db *appdbimpl) GetToken(username string) (int, error) {
    row := db.c.QueryRow("SELECT token FROM users WHERE username = ?", username)
    var token int
    err := row.Scan(&token)
    if err != nil {
        return 0, err
    }

    // Return the retrieved token
    return token, nil
}


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