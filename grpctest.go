package main

import (
	"github.com/cyberimp/grpctest/gorm"
	"github.com/cyberimp/grpctest/grpcserver"
	"github.com/cyberimp/grpctest/rest"
	"log"
	"time"
)

func main() {
	logChannel := make(chan string)
	r := &rest.LogServer{}
	g := &grpcserver.GRPCServer{}
	provider := &gorm.Conn{}
	err := provider.ConnectDB()
	if err != nil {
		log.Fatalf("could not connect to db, error:%q", err)
	}
	r.Serve(logChannel)
	g.Serve(provider, logChannel)

	time.Sleep(time.Second * 1000)
}
