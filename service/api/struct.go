package api

import (
	"time"

	"github.com/zoematr/serivce/database"
)

type User struct{
	Username   string
	UserId     string
	Followers  string
	Following  string
	Banned     string
}

type Photo struct{
	PhotoId    string
	UserId     string
	PhotoFile  string
	Date       time.Time
}

type Like struct {
	PhotoId    string
	UserId     string
	LikeId     string
}

type Comment struct {
	PhotoId    string
	UserId     string
	CommentId  string
	Date       time.Time
	CommentContent
}