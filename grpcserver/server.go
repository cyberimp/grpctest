package grpcserver

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	db "grpctest/gorm"
	pb "grpctest/grpc"
	"log"
	"net"
)

type GRPCServer struct {
	pb.UnimplementedPostsServer
	db      *db.Conn
	logChan chan string
}

func (s GRPCServer) Serve(provider *db.Conn, logProvider chan string) {
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

func getIP(ctx context.Context) string {
	p, _ := peer.FromContext(ctx)
	return p.Addr.String()
}

func (s GRPCServer) log(command string, ctx context.Context) {
	user := getIP(ctx)
	s.logChan <- fmt.Sprintf("command %q from %s", command, user)
}

func (s *GRPCServer) CreatePost(ctx context.Context, post *pb.Post) (*pb.Id, error) {
	s.log("CREATE", ctx)
	newPost := new(db.Post)
	newPost.ParseRPCPost(post)
	result, err := s.db.Create(newPost)
	return &pb.Id{Id: uint32(result)}, err
}

func (s *GRPCServer) ReadPost(ctx context.Context, id *pb.Id) (*pb.Post, error) {
	s.log("READ", ctx)
	numId := id.Id
	result, err := s.db.Read(uint(numId))
	return result.MakeRPCPost(), err
}

func (s *GRPCServer) UpdatePost(ctx context.Context, fpi *pb.FullPostInfo) (*pb.Ok, error) {
	s.log("UPDATE", ctx)
	newPost := new(db.Post)
	newPost.ParseFullRPCPost(fpi)
	_, err := s.db.Update(newPost)
	ok := &pb.Ok{Ok: err == nil, Message: err.Error()}
	return ok, err
}

func (s *GRPCServer) DeletePost(ctx context.Context, id *pb.Id) (*pb.Ok, error) {
	s.log("DELETE", ctx)
	err := s.db.Delete(uint(id.Id))
	ok := &pb.Ok{Ok: err == nil, Message: err.Error()}
	return ok, err
}

func (s *GRPCServer) ListPosts(page *pb.Pagination, srv pb.Posts_ListPostsServer) error {
	s.log("LIST", srv.Context())
	posts, err := s.db.List(int(page.Page), int(page.Size))
	if err != nil {
		return err
	}
	for _, post := range posts {
		err := srv.Send(post.MakeRPCPost())
		if err != nil {
			return err
		}
	}
	return nil
}
