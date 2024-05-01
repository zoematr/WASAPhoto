package api

import (
	"github.com/zoematr/WASAPhoto/service/database"
)

type UsernameUpdate struct {
	NewUsername string `json:"newusername"`
}

type PhotoInput struct {
	PhotoFile string `json:"photofile"`
}

type UsernameReqBody struct {
	Username string `json:"username"`
}

/*
type Content struct {
	CommentContent string `json:"content"`
}
*/

type User struct {
	Username  string   `json: username`
	Followers []string `json: following`
	Following []string `json: username`
	Banned    []string `json: username`
	Token     string   `json: token`
}

// user + photos + if already following or already banned
type UserProfile struct {
	Username        string                   `json: username`
	Followers       []string                 `json: followers`
	Following       []string                 `json: followingb`
	Photos          []database.CompletePhoto `json: photos`
	AlreadyFollowed bool                     `json: alreadyfollowed`
	AlreadyBanned   bool                     `json: alreadybanned`
	OwnProfile      bool                     `json: ownprofile`
}

type Photo struct {
	PhotoId   string `json: photoid`
	Username  string `json: username`
	PhotoFile []byte `json: photofile`
	Date      string `json: datetime`
}

type CompletePhoto struct {
	PhotoId      string             `json: photoid`
	Username     string             `json: username`
	PhotoFile    []byte             `json: photofile`
	Date         string             `json: datetime`
	AlreadyLiked bool               `json: alreadyliked`
	Likes        []database.Like    `json: likes`
	Comments     []database.Comment `json: comments`
}

type Like struct {
	PhotoId  string `json: photoid`
	Username string `json: username`
	LikeId   string `json: likeid`
}

type Comment struct {
	PhotoId        string `json: photoid`
	Username       string `json: username`
	CommentId      int64  `json: commentid`
	CommentContent string `json: commentcontent`
}

/*
type CommentId struct {
	CommentId string `json: commentid`
}

type LikeId struct {
	LikeId string `json: likeid`
}

type PhotoId struct {
	PhotoId string `json: photoid`
}

type CommentContent struct {
	CommentContent string `json: commentcontent`
}
*/

// now functions to convert the types defined before into a type of the database package

func (u User) ToDatabase() database.User {
	return database.User{
		Username:  u.Username,
		Followers: u.Followers,
		Following: u.Following,
		Banned:    u.Banned,
		Token:     u.Token,
	}
}

func (up UserProfile) ToDatabase() database.UserProfile {
	return database.UserProfile{
		Username:        up.Username,
		Followers:       up.Followers,
		Following:       up.Following,
		Photos:          up.Photos,
		AlreadyFollowed: up.AlreadyFollowed,
		AlreadyBanned:   up.AlreadyBanned,
		OwnProfile:      up.OwnProfile,
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

func (cph CompletePhoto) ToDatabase() database.CompletePhoto {
	return database.CompletePhoto{
		PhotoId:      cph.PhotoId,
		Date:         cph.Date,
		Username:     cph.Username,
		PhotoFile:    cph.PhotoFile,
		AlreadyLiked: cph.AlreadyLiked,
		Likes:        cph.Likes,
		Comments:     cph.Comments,
	}
}

func (ph Like) ToDatabase() database.Like {
	return database.Like{
		PhotoId:  ph.PhotoId,
		Username: ph.Username,
		LikeId:   ph.LikeId,
	}
}

// func (uid Username) ToDatabase() database.Username {
// 	return database.Username{
//		Username: uid.Username,
//	 }
// }

/*
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

func (phid PhotoId) ToDatabase() database.PhotoId {
	return database.PhotoId{
		PhotoId: phid.PhotoId,
	}
}
*/

func (c Comment) ToDatabase() database.Comment {
	return database.Comment{
		CommentId:      c.CommentId,
		PhotoId:        c.PhotoId,
		Username:       c.Username,
		CommentContent: c.CommentContent,
	}
}
