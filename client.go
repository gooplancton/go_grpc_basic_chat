package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"playground/gen/chatrpc"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func mainClient(username string) {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := chatrpc.NewChatServiceClient(conn)
	stream, err := client.GetMessages(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer stream.CloseSend()

	go func() {
		for {
			in, _ := stream.Recv()
			fmt.Println("\033[34m", in.Username, ":", "\033[0m", in.Msg)
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		msg, _ := reader.ReadString('\n')
		client.SendMessage(context.Background(), &chatrpc.ChatMessage{Msg: msg, Username: username})
	}

}
