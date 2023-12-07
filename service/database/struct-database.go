package database

import "time"

type Photo struct {
	PhotoId    string
	UserId     string
	PhotoFile  string
	ListLikes  []UserId
	Date       time.Time
}


type UserId struct {
	UserId     
}
