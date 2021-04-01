package main

import (
	"google.golang.org/grpc"
	"io"
	pb "khalid.foundation/heartbeat/proto"
	"log"
	"net"
	"strconv"
	"time"
)


type server struct {
	pb.UnimplementedStreamServer
	set map[string]struct{}
}

func newServer() *server {
	s := &server{
		set: make(map[string]struct{}, 0),
	}
	return s
}

func (s *server) HeartbeatStream(stream pb.Stream_HeartbeatStreamServer) error {
	//接收消息
	errCh := make(chan error, 2)
	go func() {
		for {
			in, err := stream.Recv()
			// 结束
			if err == io.EOF {
				errCh <- err
				log.Println("finish recv message")
				return
			}
			if err != nil {
				errCh <- err
				return
			}
			s.set[in.Node] = struct{}{}
			log.Println(s.set)
		}
	}()
	//定时发送消息
	go func() {
		for {
			select {
			case <-time.After(2*time.Second):
				msg := pb.StreamMsg{
					Node:   "master",
					Type:   "info",
					Number:  strconv.Itoa(len(s.set)),
				}
				if err := stream.Send(&msg); err != nil {
					errCh <-err
					return
				}
			}
		}
	}()
	err := <-errCh
	log.Print("err: ", err)
	return err
}


func main() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("fail to listen %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterStreamServer(grpcServer, newServer())
	log.Println("grpc server listen in localhost:50051")
	grpcServer.Serve(lis)
}
