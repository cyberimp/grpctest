package gorm

import (
	pb "github.com/cyberimp/grpctest/grpc"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title string
	Text  string
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
