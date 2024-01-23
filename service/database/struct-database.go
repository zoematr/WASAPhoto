package database

import (
	"time"
	"WASAPhoto/serivce/database"
)

type User struct{ 
	Username   string   `json: username`
	UserId     string   `json: userid`
	Followers  []UserId `json: following`
	Following  []UserId `json: userid`
	Banned     []UserId `json: userid`
}

type UserId struct{
	UserId string `json: userid`
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
	CommentId string `json: commentcontent`
}