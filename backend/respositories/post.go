package respositories

import (
	"time"

	"gallery/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

//CRUD logic (your Get, Update, Delete functions)

func GetAllPosts(db *gorm.DB) ([]models.Post, error) {
	var posts []models.Post
	if err := db.Order("date_created desc").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func CreatePost(db *gorm.DB, post *models.Post) error {
	post.ID = uuid.New()
	post.DateCreated = time.Now()

	if err := db.Create(post).Error; err != nil {
		return err
	}
	return nil
}
