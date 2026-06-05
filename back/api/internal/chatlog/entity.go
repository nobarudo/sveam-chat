package chatlog

import "gorm.io/gorm"

type Chatlog struct {
	gorm.Model
	userID  string
	message string
}
