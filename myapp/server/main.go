package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"myapp/config"
	"myapp/handler"
	"myapp/inits"
	__ "myapp/proto"
	"net"
)

func main() {
	config.Inits()
	inits.InitMysql()
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 7777))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	__.RegisterMyappServer(s, &handler.Server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
