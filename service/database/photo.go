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

func (db *appdbimpl) DeleteLike(photoId string, likerUsername string) error {
    _, err := db.c.Exec("DELETE FROM photos WHERE photoid = ? AND username = ?",
    	photoId, likerUsername)

    if err != nil {
        // Error executing query
        return err
    }

    return nil
}


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

func (db *appdbimpl) GetUsernameFromCommentId(commentid string) (string, error) {

	var username string

	// Utilizza una query SQL SELECT per cercare il nickname dell'utente nella tabella users utilizzando l'identificativo dell'utente (id_user).
	err := db.c.QueryRow(`SELECT username FROM comments WHERE commentid = ?`, commentid).Scan(&username)
	if err != nil {
		// Error during the execution of the query
		return "", err
	}
	return username, err
}



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

func (db *appdbimpl) AddComment(c Comment) error {
    var lastCommentID int

    // Query the last inserted photo ID
    err := db.c.QueryRow("SELECT MAX(commentid) FROM comments").Scan(&lastCommentID)
    if err != nil && err != sql.ErrNoRows {
        // Error occurred while querying
        return err
    }

    // Increment the last photo ID to get the new photo ID
    newCommentID := lastCommentID + 1

    // Utilize a SQL INSERT query to insert the photo into the database
    _, err = db.c.Exec("INSERT INTO photos (photoid, username, date, content, commentid) VALUES (?, ?, ?)",
    	newCommentID, c.Username, c.Date, c.PhotoId, c.CommentContent)

    if err != nil {
        // Error executing query
        return err
    }

    return nil
}

func (db *appdbimpl) DeleteComment(commentid string) (error) {
	// Utilizza una query SQL INSERT per inserire la foto nel database.
	_, err := db.c.Exec("DELETE FROM comments WHERE commentid = ?",
		commentid)

	if err != nil {
		// Error executing query
		return err
	}

	return err
}

func (db *appdbimpl) PhotoExists(searchedphotoid string) (bool, error) {
	// checks if a photo exists
	var cnt int
	err := db.c.QueryRow("SELECT COUNT(*) FROM photos WHERE photoid = ?",
		searchedphotoid).Scan(&cnt)

	if err != nil {
		// Count always returns a row thanks to COUNT(*), so this situation should not happen
		return false, err
	}

	// If counter 1 then the photo exists
	if cnt == 1 {
		return true, nil
	}
	return false, nil
}

func (db *appdbimpl) CommentExists(commentid string) (bool, error) {
	// checks if a photo exists
	var cnt int
	err := db.c.QueryRow("SELECT COUNT(*) FROM comments WHERE commentid = ?",
		commentid).Scan(&cnt)

	if err != nil {
		// Count always returns a row thanks to COUNT(*), so this situation should not happen
		return false, err
	}

	// If counter 1 then the photo exists
	if cnt == 1 {
		return true, nil
	}
	return false, nil
}

func (db *appdbimpl) DoesUserLikePhoto(photoid string, likerusername string) (bool, error) {
	// checks if a photo exists
	var cnt int
	err := db.c.QueryRow("SELECT COUNT(*) FROM likes WHERE photoid = ? AND username = ?",
		photoid, likerusername).Scan(&cnt)

	if err != nil {
		// Count always returns a row thanks to COUNT(*), so this situation should not happen
		return false, err
	}

	// If counter 1 then the photo exists
	if cnt == 1 {
		return true, nil
	}
	return false, nil
}


