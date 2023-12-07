package database

import "time"

type Photo struct {
	PhotoId  string      `json:"photoid"` 
	UserId   string     `json:"owner"`
	Likes    []UserId   `json:"likes"`
	Date     time.Time  `json:"date"`     // Date in which the photo was uploaded
}