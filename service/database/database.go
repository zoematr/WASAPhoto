/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {

	// the structure seems to be
	// nameFunction(input) (output)
	GetName() (string, error)
	SetName(name string) error
	CreateUser(string) (int, error)
	GetStream(string) ([]CompletePhoto, error)
	ExistsUser(string) (bool, error)
	GetTokenFromUsername(string) (int, error)
	ChangeUsername(int, string) error
	GetUsernameFromToken(int) (string, error)
	GetFollowers(string) ([]string, error)
	GetFollowing(string) ([]string, error)
	GetPhotos(string, string) ([]CompletePhoto, error)
	CheckBanned(string, string) (bool, error)
	FollowUser(string, string) error
	UnfollowUser(string, string) error
	WasTargetFollowed(string, string) (bool, error)
	WasTargetBanned(string, string) (bool, error)
	BanUser(string, string) error
	UnbanUser(string, string) error
	AddPhoto(Photo) error
	DeletePhoto(string) error
	AddLike(string, string) error
	GetUsernameFromPhotoId(string) (string, error)
	AddComment(Comment) error
	DeleteLike(string, string) error
	PhotoExists(string) (bool, error)
	DoesUserLikePhoto(string, string) (bool, error)
	GetUsernameFromCommentId(string) (string, error)
	DeleteComment(string) error
	CommentExists(string) (bool, error)
	GetPhotoFromPhotoId(string) (CompletePhoto, error)
	GetLikes(CompletePhoto) error
	GetComments(CompletePhoto) error

	// Ping checks availability of the database, if not it returns an error.
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}
	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&tableName)
	// fmt.Println("this is the db error")
	// fmt.Println(err)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `

		CREATE TABLE IF NOT EXISTS example_table (
			id INTEGER NOT NULL, 
			name TEXT,
			UNIQUE(id)
			);
		
		CREATE TABLE IF NOT EXISTS users (
			token INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			username TEXT UNIQUE NOT NULL
			);

		CREATE TABLE IF NOT EXISTS photos (
			photoid TEXT NOT NULL PRIMARY KEY,
			username TEXT NOT NULL,
			photofile VARBINARY(100000) NOT NULL,
			datetime TEXT NOT NULL DEFAULT '0000-01-01T00:00:00Z',
			FOREIGN KEY(username) REFERENCES users(username) ON DELETE CASCADE
			);
		
		CREATE TABLE IF NOT EXISTS likes (
			photoid TEXT NOT NULL, 
			username TEXT NOT NULL,
			FOREIGN KEY(username) REFERENCES users(username) ON DELETE CASCADE,
			FOREIGN KEY(photoid) REFERENCES photos(photoid) ON DELETE CASCADE,
			PRIMARY KEY (photoid, username)
			);
		
		CREATE TABLE IF NOT EXISTS comments (
			commentid TEXT NOT NULL PRIMARY KEY,
			photoid TEXT NOT NULL, 
			username TEXT NOT NULL,
			content TEXT NOT NULL,
			FOREIGN KEY(photoid) REFERENCES photos(photoid) ON DELETE CASCADE,
			FOREIGN KEY(username) REFERENCES users(username) ON DELETE CASCADE
			);
		
		CREATE TABLE IF NOT EXISTS banned (
			username TEXT NOT NULL, 
			bannedusername TEXT NOT NULL,
			FOREIGN KEY(username) REFERENCES users(username) ON DELETE CASCADE,
			FOREIGN KEY(bannedusername) REFERENCES users(username) ON DELETE CASCADE,
			PRIMARY KEY (username, bannedusername)
			);

		CREATE TABLE IF NOT EXISTS followers (
			username TEXT NOT NULL, 
			followerusername TEXT NOT NULL,
			FOREIGN KEY(username) REFERENCES users(username) ON DELETE CASCADE,
			FOREIGN KEY(followerusername) REFERENCES users(username) ON DELETE CASCADE,
			PRIMARY KEY (username, followerusername)
			);
		
		`

		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
