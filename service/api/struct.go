package api

import (
	"time"
	"github.com/zoematr/WASAPhoto/service/database"
)

type UserId struct{
	UserId string `json: userid`
}

type User struct{ 
	Username   string   `json: username`
	UserId     string   `json: userid`
	Followers  []UserId `json: following`
	Following  []UserId `json: userid`
	Banned     []UserId `json: userid`
}

type Username struct{
	Username string `json: username`
}

type Photo struct{
	PhotoId    string    `json: photoid`
	UserId     string    `json: userid`       
	PhotoFile  string    `json: photofile`
	Date       time.Time `json: datetime`
}

type PhotoId struct{
	PhotoId string `json: photoid`
}

type Like struct {
	PhotoId    string `json: photoid`
	UserId     string `json: userid`
	LikeId     string `json: likeid`
}

type LikeId struct{
	LikeId string `json: likeid`
}

type Comment struct {
	PhotoId        string    `json: photoid`
	UserId         string    `json: userid`
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
	return database.UserId{
		Username:  u.Username,
		UserId:    u.UserId,
		Followers: u.Followers,
		Following: u.Following,
		Banned:    u.Banned,
	}
}

func (ph Photo) ToDatabase() database.Photo {
	return database.Photo{
		PhotoId:   ph.PhotoId,
		Date:      ph.Date,
		UserId:    ph.UserId,
		PhotoFile: ph.PhotoFile,
	}
}

func (ph Like) ToDatabase() database.Like {
	return database.Like{
		PhotoId: ph.PhotoId,
		UserId:  ph.UserId,
		LikeId:  h.LikeId,
	}
}

func (phid PhotoId) ToDatabase() database.PhotoId {
	return database.PhotoId{
		PhotoId: phid.PhotoId,
	}
}

func (uid UserId) ToDatabase() database.UserId {
	return database.UserId{
		UserId: uid.UserId,
	}
}

func (uid Username) ToDatabase() database.Username {
	return database.Username{
		Username: uid.Username,
	}
}

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
	return database.CompleteComment{
		CommentId:      c.CommentId,
		PhotoId:        c.PhotoId,
		UserId:         c.UserId,
		CommentContent: c.CommentContent,
	}
}