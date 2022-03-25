package main

import (
	"context"
	"log"
	"math/rand"
	"net"

	pb "github.com/lindsay0416/begin_grpc/usermgmt"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type UserManagementServer struct {
	//This UserManagementServer is the implementations of our gRPC service
	pb.UnimplementedUserManagementServer
}

func (s *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	log.Printf("Received: %v", in.GetName())
	var user_id int32 = int32(rand.Intn(100))
	return &pb.User{Name: in.GetName(), Age: in.GetAge(), Id: user_id}, nil
}

func main() {
	// Listen the port
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Fail to listen the port: %v", err)
	}

	// init a new server
	s := grpc.NewServer()
	// Register the new Server
	pb.RegisterUserManagementServer(s, &UserManagementServer{})
	log.Printf("Server listening at: %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
