package database

import (
	"database/sql"
)

func (db *appdbimpl) AddPhoto(p Photo) error {
	// function to add the photo with the correct photo id
	var lastPhotoID int

	// Query the last inserted photo ID
	err := db.c.QueryRow("SELECT COALESCE(MAX(photoid), 0) FROM photos").Scan(&lastPhotoID)
	if err != nil && err != sql.ErrNoRows {
		// Error occurred while querying
		return err
	}

	// Increment the last photo ID to get the new photo ID
	newPhotoID := lastPhotoID + 1

	// add db
	_, err = db.c.Exec("INSERT INTO photos (photoid, username, datetime, photofile) VALUES (?, ?, ?, ?)",
		newPhotoID, p.Username, p.Date, p.PhotoFile)

	if err != nil {
		// Error executing query
		return err
	}

	return nil
}

// TODO CHANGE TO COMPLETE PHOTO
func (db *appdbimpl) GetPhotoFromPhotoId(photoid string) (CompletePhoto, error) {
	// function to get username-> author of a picture from the photo id
	var photo CompletePhoto

	// look for username where id of the photo is the input
	err := db.c.QueryRow(`SELECT username, photofile, datetime FROM photos WHERE photoid = ?`, photoid).Scan(&photo.Username, &photo.PhotoFile, &photo.Date)
	if err != nil {
		// Error during the execution of the query
		return photo, err
	}
	err = db.GetLikes(photo)
	if err != nil {
		// Error during the execution of the query
		return photo, err
	}
	err = db.GetComments(photo)
	if err != nil {
		// Error during the execution of the query
		return photo, err
	}
	return photo, nil
}

func (db *appdbimpl) AddLike(photoId string, likerUsername string) error {
	// function to add like
	_, err := db.c.Exec("INSERT INTO photos (username, photoid) VALUES (?, ?, ?)",
		photoId, likerUsername)

	if err != nil {
		// Error executing query
		return err
	}

	return nil
}

func (db *appdbimpl) DeleteLike(photoId string, likerUsername string) error {
	// function to unlike
	_, err := db.c.Exec("DELETE FROM photos WHERE photoid = ? AND username = ?",
		photoId, likerUsername)

	if err != nil {
		// Error executing query
		return err
	}

	return nil
}

func (db *appdbimpl) GetUsernameFromPhotoId(photoid string) (string, error) {
	// function to get username-> author of a picture from the photo id
	var username string

	// look for username where id of the photo is the input
	err := db.c.QueryRow(`SELECT username FROM photos WHERE photoid = ?`, photoid).Scan(&username)
	if err != nil {
		// Error during the execution of the query
		return "", err
	}
	return username, err
}

func (db *appdbimpl) GetUsernameFromCommentId(commentid string) (string, error) {
	// function to get username-> author of a comment from the comment id
	var username string

	err := db.c.QueryRow(`SELECT username FROM comments WHERE commentid = ?`, commentid).Scan(&username)
	if err != nil {
		// Error during the execution of the query
		return "", err
	}
	return username, err
}

func (db *appdbimpl) DeletePhoto(photoId string) error {
	// function to delete photo from db
	_, err := db.c.Exec("DELETE FROM photos WHERE photoid = ?",
		photoId)

	if err != nil {
		// Error executing query
		return err
	}

	return err
}

func (db *appdbimpl) AddComment(c Comment) error {
	// function to comment a photo
	// data is passed in the struct from the backend
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
		newCommentID, c.Username, c.PhotoId, c.CommentContent)

	if err != nil {
		// Error executing query
		return err
	}

	return nil
}

func (db *appdbimpl) DeleteComment(commentid string) error {
	// function to delete a comment
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
	// checks if a comment exists
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
	// checks if a user has liked a photo
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

// get user stream
func (db *appdbimpl) GetStream(username string) ([]CompletePhoto, error) {
	rows, err := db.c.Query(`SELECT * FROM photos WHERE username IN (SELECT username FROM followers WHERE followerusername = ?) ORDER BY datetime DESC`,
		username)
	if err != nil {
		return nil, err
	}
	// Wait for the func to finish before closing rows
	defer func() { _ = rows.Close() }()

	// Read all the users in the resulset
	var res []CompletePhoto
	for rows.Next() {
		var photo CompletePhoto
		photo.AlreadyLiked = false
		err = rows.Scan(&photo.PhotoId, &photo.Username, &photo.PhotoFile, &photo.Date)
		if err != nil {
			return nil, err
		}
		err = db.GetLikes(photo)
		if err != nil {
			return nil, err
		}
		err = db.GetComments(photo)
		if err != nil {
			return nil, err
		}
		isliked, err := db.DoesUserLikePhoto(photo.PhotoId, username)
		if err != nil {
			return nil, err
		}
		if isliked {
			photo.AlreadyLiked = true
		}
		res = append(res, photo)
	}

	if rows.Err() != nil {
		return nil, err
	}
	// gives back slice of Photo which is the stream.
	return res, nil
}

func (db *appdbimpl) GetLikes(photo CompletePhoto) error {
	rows, err := db.c.Query(`SELECT * FROM likes WHERE photoid = ?`,
		photo.PhotoId)
	if err != nil {
		return err
	}
	// Wait for the func to finish before closing rows
	defer func() { _ = rows.Close() }()

	// Read all the likes in the resulset
	var likes []Like
	for rows.Next() {
		var l Like
		err = rows.Scan(&l.PhotoId, &l.Username)
		if err != nil {
			return err
		}
		likes = append(likes, l)
	}

	if rows.Err() != nil {
		return err
	}
	photo.Likes = likes
	return nil
}

func (db *appdbimpl) GetComments(photo CompletePhoto) error {
	rows, err := db.c.Query(`SELECT * FROM comments WHERE photoid = ?`,
		photo.PhotoId)
	if err != nil {
		return err
	}
	// Wait for the func to finish before closing rows
	defer func() { _ = rows.Close() }()

	// Read all the likes in the resulset
	var comments []Comment
	for rows.Next() {
		var c Comment
		err = rows.Scan(&c.CommentId, &c.PhotoId, &c.Username, &c.CommentContent)
		if err != nil {
			return err
		}
		comments = append(comments, c)
	}

	if rows.Err() != nil {
		return err
	}
	photo.Comments = comments
	return nil
}

// function that gets all the photos of a user
func (db *appdbimpl) GetPhotos(username string, requesting string) ([]CompletePhoto, error) { // TO DO check if the photos are liked by a user

	var photos []CompletePhoto
	rows, err := db.c.Query(`SELECT * FROM photos WHERE username = ?`, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var photo CompletePhoto
		if err := rows.Scan(&photo.PhotoId, &photo.Username, &photo.PhotoFile, &photo.Date); err != nil {
			return nil, err
		}
		err = db.GetLikes(photo)
		if err != nil {
			return nil, err
		}
		err = db.GetComments(photo)
		if err != nil {
			return nil, err
		}
		isliked, err := db.DoesUserLikePhoto(photo.PhotoId, requesting)
		if err != nil {
			return nil, err
		}
		if isliked {
			photo.AlreadyLiked = true
		}
		photos = append(photos, photo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return photos, nil
}
