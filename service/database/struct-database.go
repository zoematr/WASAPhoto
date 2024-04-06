package database


type User struct {
	Username  string   `json: username`
	Followers []string `json: following`
	Following []string `json: userid`
	Banned    []string `json: userid`
	Token     string   `json: token`
}

// user + photos
type UserProfile struct {
	Username        string          `json: username`
	Followers       []string        `json: followers`
	Following       []string        `json: followingb`
	Photos          []CompletePhoto `json: photos`
	AlreadyFollowed bool            `json: alreadyfollowed`
	AlreadyBanned   bool            `json: alreadybanned`
	OwnProfile      bool            `json: ownprofile`
}

// type Username struct{
// 	Username string `json: username`
// }

type Photo struct {
	PhotoId   string    `json: photoid`
	Username  string    `json: username`
	PhotoFile []byte    `json: photofile`
	Date      string    `json: datetime`
}

type CompletePhoto struct {
	PhotoId      string    `json: photoid`
	Username     string    `json: username`
	PhotoFile    []byte    `json: photofile`
	Date         string    `json: datetime`
	AlreadyLiked bool      `json: alreadyliked`
	Likes        []Like    `json: likes`
	Comments     []Comment `json: comments`
}

type Like struct {
	PhotoId  string `json: photoid`
	Username string `json: username`
	LikeId   string `json: likeid`
}

type Comment struct {
	PhotoId        string `json: photoid`
	Username       string `json: username`
	CommentId      string `json: commentid`
	CommentContent string `json: commentcontent`
}

/*
type LikeId struct {
	LikeId string `json: likeid`
}

type CommentId struct {
	CommentId string `json: commentid`
}

type CommentContent struct {
	CommentContent string `json: commentcontent`
}

type PhotoId struct {
	PhotoId string `json: photoid`
}
*/
