package respositories

import (
	"log"
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

func GetPostById(db *gorm.DB, uuid uuid.UUID) (*models.Post, error) {
	var post *models.Post
	if err := db.First(&post, "id = ?", uuid).Error; err != nil {
		log.Println("Post not found:", err)
		return nil, err
	}

	return post, nil
}

func CreatePost(db *gorm.DB, post *models.Post) error {
	post.ID = uuid.New()
	post.DateCreated = time.Now()

	if err := db.Create(post).Error; err != nil {
		return err
	}
	return nil
}

func UpdatePost(db *gorm.DB, post *models.Post) error {
	if err := db.Save(&post).Error; err != nil {
		log.Println("Failed to update post:", err)
		return err
	}
	return nil
}

func DeletePost(db *gorm.DB, uuid uuid.UUID) error {
	if err := db.Delete(&models.Post{}, uuid).Error; err != nil {
		log.Println("Failed to delete post from database:", err)
		return err
	}

	return nil
}
