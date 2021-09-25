package gorm

import (
	"gorm.io/gorm"
)

type Conn struct {
	db *gorm.DB
}

func (s *Conn) Create(newPost *Post) (uint, error) {
	post := newPost
	result := s.db.Create(post)
	return post.ID, result.Error
}

func (s *Conn) Read(id uint) (*Post, error) {
	post := new(Post)
	result := s.db.First(post, id)
	return post, result.Error
}

func (s *Conn) Update(newPost *Post) (uint, error) {
	post := Post{}
	find := s.db.First(&post, newPost.ID)
	if find.Error != nil {
		return 0, find.Error
	}

	post.Text = newPost.Text
	post.Title = newPost.Title
	result := s.db.Save(&post)

	return post.ID, result.Error
}

func (s *Conn) Delete(id uint) error {
	result := s.db.Delete(&Post{}, id)
	return result.Error
}

func (s *Conn) List(page int, size int) ([]Post, error) {
	posts := make([]Post, 0)
	result := s.db.Limit(size).Offset(page * size).Find(&posts)
	return posts, result.Error
}

func (s *Conn) migrate() error {
	return s.db.AutoMigrate(&Post{})
}
