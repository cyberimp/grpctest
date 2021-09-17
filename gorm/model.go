package gorm

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title string
	Text  string
}

type Conn struct {
	db *gorm.DB
}

func (s *Conn) Create(newPost Post) (error, uint) {
	post := newPost
	result := s.db.Create(&post)
	return result.Error, post.ID
}

func (s *Conn) Read(id uint) (error, Post) {
	post := Post{}
	result := s.db.First(&post, id)
	return result.Error, post
}

func (s *Conn) Update(newPost Post) (error, uint) {
	post := Post{}
	find := s.db.First(&post, newPost.ID)
	if find.Error != nil {
		return find.Error, 0
	}

	post.Text = newPost.Text
	post.Title = newPost.Title
	result := s.db.Save(&post)

	return result.Error, post.ID
}

func (s *Conn) Delete(id uint) error {
	result := s.db.Delete(&Post{}, id)
	return result.Error
}

func (s *Conn) List() (error, []Post) {
	posts := make([]Post, 0)
	result := s.db.Find(&posts)
	return result.Error, posts
}
