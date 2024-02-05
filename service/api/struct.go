package api

import (
	"time"
	"github.com/zoematr/WASAPhoto/service/database"
)

type User struct{ 
	Username   string   `json: username`
	Followers  []string `json: following`
	Following  []string `json: username`
	Banned     []string `json: username`
	Token      string   `json: token`
}

// type Username struct{
// 	Username string `json: username`
// }

type Photo struct{
	PhotoId    string    `json: photoid`
	Username   string    `json: username`       
	PhotoFile  string    `json: photofile`
	Date       time.Time `json: datetime`
}

type PhotoId struct{
	PhotoId string `json: photoid`
}

type Like struct {
	PhotoId    string `json: photoid`
	Username   string `json: username`
	LikeId     string `json: likeid`
}

type LikeId struct{
	LikeId string `json: likeid`
}

type Comment struct {
	PhotoId        string    `json: photoid`
	Username       string    `json: username`
	CommentId      string    `json: commentid`
	Date           time.Time `json: datetime`
	CommentContent string    `json: commentcontent`
}

type CommentId struct{
	CommentId string `json: commentid`
}

type CommentContent struct{
	CommentContent string `json: commentcontent`
}

// now functions to convert the types defined before into a type of the database package

func (u User) ToDatabase() database.User {
	return database.User{
		Username:  u.Username,
		Followers: u.Followers,
		Following: u.Following,
		Banned:    u.Banned,
	}
}

func (ph Photo) ToDatabase() database.Photo {
	return database.Photo{
		PhotoId:   ph.PhotoId,
		Date:      ph.Date,
		Username:  ph.Username,
		PhotoFile: ph.PhotoFile,
	}
}

func (ph Like) ToDatabase() database.Like {
	return database.Like{
		PhotoId:  ph.PhotoId,
		Username: ph.Username,
		LikeId:   ph.LikeId,
	}
}

func (phid PhotoId) ToDatabase() database.PhotoId {
	return database.PhotoId{
		PhotoId: phid.PhotoId,
	}
}

// func (uid Username) ToDatabase() database.Username {
// 	return database.Username{
//		Username: uid.Username,
//	 }
// }

func (cc CommentContent) ToDatabase() database.CommentContent {
	return database.CommentContent{
		CommentContent: cc.CommentContent,
	}
}

func (cid CommentId) ToDatabase() database.CommentId {
	return database.CommentId{
		CommentId: cid.CommentId,
	}
}

func (cid LikeId) ToDatabase() database.LikeId {
	return database.LikeId{
		LikeId: cid.LikeId,
	}
}

func (c Comment) ToDatabase() database.Comment {
	return database.Comment{
		CommentId:      c.CommentId,
		PhotoId:        c.PhotoId,
		Username:       c.Username,
		CommentContent: c.CommentContent,
	}
}