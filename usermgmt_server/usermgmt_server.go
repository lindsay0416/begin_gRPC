package main

import (
	"context"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"

	pb "github.com/lindsay0416/begin_grpc/usermgmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson" // Store users into a json file
)

const (
	port = ":50051"
)

func NewUserManagementServer() *UserManagementServer {
	return &UserManagementServer{
		// user_list: &pb.UserList{},
	}
}

type UserManagementServer struct {
	//This UserManagementServer is the implementations of our gRPC service
	pb.UnimplementedUserManagementServer
	// user_list *pb.UserList
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
	readBytes, err := ioutil.ReadFile("user.json")
	var users_list *pb.UserList = &pb.UserList{}
	var user_id int32 = int32(rand.Intn(100))
	create_user := &pb.User{Name: in.GetName(), Age: in.GetAge(), Id: user_id}

	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("File is not found. Creating a new file.")
			users_list.Users = append(users_list.Users, create_user)
			jsonBytes, err := protojson.Marshal(users_list)
			if err != nil {
				log.Fatalf("JSON Marshal failed")
			}
			if err = ioutil.WriteFile("user.json", jsonBytes, 0664); err != nil {
				log.Fatalf("Failed write to thw file: %v", err)
			}
			return create_user, nil
		} else {
			log.Fatalln("Error reading file: ", err)
		}
	}

	if err = protojson.Unmarshal(readBytes, users_list); err != nil {
		log.Fatalf("Failed to parse user list: %v", err)
	}
	users_list.Users = append(users_list.Users, create_user)
	jsonBytes, err := protojson.Marshal(users_list)
	if err != nil {
		log.Fatalf("JSON Marshal failed")
	}
	if err = ioutil.WriteFile("user.json", jsonBytes, 0664); err != nil {
		log.Fatalf("Failed write to thw file: %v", err)
	}
	// s.user_list.Users = append(s.user_list.Users, create_user)
	return create_user, nil
}

//Implememt the getUser function
// define a new Receiver function
func (s *UserManagementServer) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UserList, error) {
	jsonBytes, err := ioutil.ReadFile("user.json")
	if err != nil {
		log.Fatalf("Failed read from the file: %v", err)
	}
	var users_list *pb.UserList = &pb.UserList{}
	if err := protojson.Unmarshal(jsonBytes, users_list); err != nil {
		log.Fatalf("Unmarshal failed: %v", err)
	}
	return users_list, nil
}

func main() {
	// instantiate a new user managemnet server
	var user_mgmt_server *UserManagementServer = NewUserManagementServer()
	// Call run function and check the error
	if err := user_mgmt_server.Run(); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
