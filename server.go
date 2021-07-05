package main

import (
	"context"
	"log"
	"net"
	"playground/gen/chatrpc"

	"github.com/guiguan/caster"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	chatrpc.UnimplementedChatServiceServer
}

type chatMsg struct {
	content  string
	username string
}

var msgCaster = caster.New(context.Background())

func mainServer() {

	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	chatrpc.RegisterChatServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) GetMessages(_ *emptypb.Empty, stream chatrpc.ChatService_GetMessagesServer) error {

	allMsgs, _ := msgCaster.Sub(context.Background(), 1)
	for in := range allMsgs {
		msg := in.(chatMsg)
		stream.Send(&chatrpc.ChatMessage{Msg: msg.content, Username: msg.username})
	}
	return nil
}

func (s *server) SendMessage(ctx context.Context, msg *chatrpc.ChatMessage) (*emptypb.Empty, error) {
	msgContent := msg.GetMsg()
	msgUsername := msg.GetUsername()
	msgCaster.Pub(chatMsg{content: msgContent, username: msgUsername})
	return &emptypb.Empty{}, nil
}
