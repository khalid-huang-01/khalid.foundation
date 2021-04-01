package main

import (
	"context"
	"google.golang.org/grpc"
	"io"
	pb "khalid.foundation/heartbeat/proto"
	"log"
	"math/rand"
	"strconv"
	"time"
)

const (
	serverAddr = "localhost:50051"
)
var nodeName int

func runHeartbeat(client pb.StreamClient) {
	// 这里要改成长连接
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	stream, err := client.HeartbeatStream(ctx)
	if err != nil {
		log.Fatalf("%v Heartbeat err %v", nodeName, err)
	}
	waitCh := make(chan struct{})
	//定时发送心跳
	msg := pb.StreamMsg{
		Node:   strconv.Itoa(nodeName),
		Type:   "heartbeat",
	}
	go func() {
		for {
			select {
			case <-time.After(5*time.Second):
				if err := stream.Send(&msg); err != nil {
					log.Println("the stream disconnect err = %v", err)
					close(waitCh)
					return
				}

			}
		}
	}()
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitCh)
				return
			}
			if err != nil {
				log.Println("receive error: ", err)
			}
			log.Println("get message from server, node nums: ", in.Number)
		}
	}()
	<-waitCh
}

func main() {
	rand.Seed(time.Now().UnixNano())
	nodeName = rand.Int()

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatal("fail to dial", err)
	}
	defer conn.Close()
	client := pb.NewStreamClient(conn)
	log.Println("nodeName: ", nodeName)
	runHeartbeat(client)
}
