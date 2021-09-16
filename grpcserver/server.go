package grpcserver

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"gorm.io/gorm"
	pb "grpctest/grpc"
	"log"
	"net"
)

type GRPCServer struct {
	pb.UnimplementedPostsServer
	db *gorm.DB
	logChan chan string
}

func (s GRPCServer) Serve(provider *gorm.DB, logProvider chan string)  {
	s.db = provider
	s.logChan = logProvider

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 3009))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPostsServer(grpcServer, &s)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *GRPCServer) CreatePost(ctx context.Context, post *pb.Post) (*pb.Id, error) {
	p, _ := peer.FromContext(ctx)
	remoteUser := p.Addr.String()
	remoteCommand := "CREATE"
	s.logChan <- fmt.Sprintf("command %q from %s", remoteCommand, remoteUser)
	return nil, nil
}

func (s *GRPCServer) ReadPost(context.Context, *pb.Id) (*pb.Post, error){
	return nil, nil

}
func (s *GRPCServer) UpdatePost(context.Context, *pb.FullPostInfo) (*pb.Ok, error){
	return nil, nil

}
func (s *GRPCServer) DeletePost(context.Context, *pb.Id) (*pb.Ok, error){
	return nil, nil

}
func (s *GRPCServer) ListPosts(*pb.Pagination, pb.Posts_ListPostsServer) error{
	return nil
}