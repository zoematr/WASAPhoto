package database

import (
	"time"
)

type User struct {
	Username  string   `json: username`
	Followers []string `json: following`
	Following []string `json: userid`
	Banned    []string `json: userid`
	Token     string   `json: token`
}

// user + photos
type UserProfile struct {
	Username  string   `json: username`
	Followers []string `json: following`
	Following []string `json: username`
	Photos    []Photo  `json: photos`
}

// type Username struct{
// 	Username string `json: username`
// }

type Photo struct {
	PhotoId   string    `json: photoid`
	Username  string    `json: username`
	PhotoFile []byte    `json: photofile`
	Date      time.Time `json: datetime`
}

type PhotoId struct {
	PhotoId string `json: photoid`
}

type Like struct {
	PhotoId  string `json: photoid`
	Username string `json: username`
	LikeId   string `json: likeid`
}

type LikeId struct {
	LikeId string `json: likeid`
}

type Comment struct {
	PhotoId        string    `json: photoid`
	Username       string    `json: username`
	CommentId      string    `json: commentid`
	Date           time.Time `json: datetime`
	CommentContent string    `json: commentcontent`
}

type CommentId struct {
	CommentId string `json: commentid`
}

type CommentContent struct {
	CommentContent string `json: commentcontent`
}
