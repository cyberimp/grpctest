package gorm

import (
	pb "github.com/cyberimp/grpctest/grpc"
	"gorm.io/gorm"
	"reflect"
	"testing"
	"time"
)

func TestPost_MakeRPCPost(t *testing.T) {
	type fields struct {
		Model gorm.Model
		Title string
		Text  string
	}
	tests := []struct {
		name   string
		fields fields
		want   *pb.Post
	}{
		{"should return RPC Post",
			fields{
				gorm.Model{
					ID:        0,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				"name",
				"lorem ipsum"},
			&pb.Post{Title: "name", Text: "lorem ipsum"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Post{
				Model: tt.fields.Model,
				Title: tt.fields.Title,
				Text:  tt.fields.Text,
			}
			if got := p.MakeRPCPost(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MakeRPCPost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPost_ParseFullRPCPost(t *testing.T) {
	tests := []struct {
		name   string
		fields *Post
		args   *pb.FullPostInfo
	}{
		{
			"Should parse post",
			&Post{Model: gorm.Model{ID: 1}, Title: "name", Text: "lorem ipsum"},
			&pb.FullPostInfo{
				NewPost: &pb.Post{Title: "name", Text: "lorem ipsum"},
				Id:      &pb.Id{Id: 1}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Post{}
			p.ParseFullRPCPost(tt.args)
			if !reflect.DeepEqual(p, tt.fields) {
				t.Errorf("MakeRPCPost() = %v, want %v", p, tt.fields)
			}
		})
	}
}

func TestPost_ParseRPCPost(t *testing.T) {
	tests := []struct {
		name   string
		fields *Post
		args   *pb.Post
	}{
		{
			"should parse simple post",
			&Post{Title: "name", Text: "lorem ipsum"},
			&pb.Post{Title: "name", Text: "lorem ipsum"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Post{}
			p.ParseRPCPost(tt.args)
			if !reflect.DeepEqual(p, tt.fields) {
				t.Errorf("MakeRPCPost() = %v, want %v", p, tt.fields)
			}
		})
	}
}
