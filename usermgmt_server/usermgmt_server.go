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

func NewUserManagementServer() *UserManagementServer {
	return &UserManagementServer{
		user_list: &pb.UserList{},
	}
}

type UserManagementServer struct {
	//This UserManagementServer is the implementations of our gRPC service
	pb.UnimplementedUserManagementServer
	user_list *pb.UserList
}

func (server *UserManagementServer) Run() error {
	// Listen the port
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Fail to listen the port: %v", err)
	}

	// init a new server
	s := grpc.NewServer()
	// Register the new Server
	pb.RegisterUserManagementServer(s, server)
	log.Printf("Server listening at: %v", lis.Addr())
	return s.Serve(lis)

}

func (s *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	log.Printf("Received: %v", in.GetName())
	var user_id int32 = int32(rand.Intn(100))
	create_user := &pb.User{Name: in.GetName(), Age: in.GetAge(), Id: user_id}
	s.user_list.Users = append(s.user_list.Users, create_user)
	return create_user, nil
}

//Implememt the getUser function
// define a new Receiver function
func (s *UserManagementServer) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UserList, error) {
	return s.user_list, nil
}

func main() {
	// instantiate a new user managemnet server
	var user_mgmt_server *UserManagementServer = NewUserManagementServer()
	// Call run function and check the error
	if err := user_mgmt_server.Run(); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
