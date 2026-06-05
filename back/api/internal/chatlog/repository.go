package chatlog

import "gorm.io/gorm"

type ChatlogRepository struct {
	db *gorm.DB
}

func NewChatlogRepository(db *gorm.DB) *ChatlogRepository {
	return &ChatlogRepository{db: db}
}

func (r *ChatlogRepository) createChatlog(chatlog ChatLog) error {
	err := r.db.Create(&chatlog).Error
	return err
}
