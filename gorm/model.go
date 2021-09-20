package gorm

import (
	"gorm.io/gorm"
	pb "grpctest/grpc"
)

type Post struct {
	gorm.Model
	Title string
	Text  string
}

type Conn struct {
	db *gorm.DB
}

func (p *Post) ParseRPCPost(post *pb.Post) {
	p.Title = post.Title
	p.Text = post.Text
}

func (p *Post) ParseFullRPCPost(post *pb.FullPostInfo) {
	p.Title = post.NewPost.Title
	p.Text = post.NewPost.Text
	p.ID = uint(post.Id.Id)
}

func (p *Post) MakeRPCPost() *pb.Post {
	return &pb.Post{Title: p.Title, Text: p.Text}
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
