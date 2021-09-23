package main

import (
	"github.com/cyberimp/grpctest/grpcserver"
	"github.com/cyberimp/grpctest/rest"
	"time"
)

func main() {
	logChannel := make(chan string)
	r := rest.LogServer{}
	g := grpcserver.GRPCServer{}
	r.Serve(logChannel)
	g.Serve(nil, logChannel)

	time.Sleep(time.Second * 1000)
}
