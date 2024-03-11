package database

import (
	"database/sql"
)

func (db *appdbimpl) AddPhoto(p Photo) error {
    var lastPhotoID int

    // Query the last inserted photo ID
    err := db.c.QueryRow("SELECT MAX(photoid) FROM photos").Scan(&lastPhotoID)
    if err != nil && err != sql.ErrNoRows {
        // Error occurred while querying
        return err
    }

    // Increment the last photo ID to get the new photo ID
    newPhotoID := lastPhotoID + 1

    // Utilize a SQL INSERT query to insert the photo into the database
    _, err = db.c.Exec("INSERT INTO photos (photoid, username, date) VALUES (?, ?, ?)",
    	newPhotoID, p.Username, p.Date)

    if err != nil {
        // Error executing query
        return err
    }

    return nil
}

func (db *appdbimpl) AddLike(photoId string, likerUsername string) error {
    _, err := db.c.Exec("INSERT INTO photos (username, photoid) VALUES (?, ?, ?)",
    	photoId, likerUsername)

    if err != nil {
        // Error executing query
        return err
    }

    return nil
}

/*
func (db *appdbimpl) GetUsernameFromPhotoId(photoid string) (string, error) {

	var username string

	// Utilizza una query SQL SELECT per cercare il nickname dell'utente nella tabella users utilizzando l'identificativo dell'utente (id_user).
	err := db.c.QueryRow(`SELECT username FROM photos WHERE photoid = ?`, photoid).Scan(&username)
	if err != nil {
		// Error during the execution of the query
		return "", err
	}
	return username, err
}
*/


func (db *appdbimpl) DeletePhoto(photoId string) (error) {
	// Utilizza una query SQL INSERT per inserire la foto nel database.
	_, err := db.c.Exec("DELETE FROM photos WHERE photoid = ?",
		photoId)

	if err != nil {
		// Error executing query
		return err
	}

	return err
}

