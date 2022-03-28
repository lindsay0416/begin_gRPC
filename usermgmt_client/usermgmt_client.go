package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/lindsay0416/begin_grpc/usermgmt"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Connection Failed: %v", err)
	}
	defer conn.Close()

	//Create New User
	c := pb.NewUserManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var new_Users = make(map[string]int32)
	new_Users["Alice"] = 43
	new_Users["Bob"] = 30

	for name, age := range new_Users {
		r, err := c.CreateNewUser(ctx, &pb.NewUser{Name: name, Age: int32(age)})
		if err != nil {
			log.Fatalf("Fail  create user: %v", err)
		}
		log.Printf(`User Details:
		NAME:%s
		AGE:%d
		ID: %d`, r.GetName(), r.GetAge(), r.GetId())
	}
	params := &pb.GetUsersParams{}
	r, err := c.GetUsers(ctx, params)
	if err != nil {
		log.Fatalf("Could not retrieve users: %v", err)
	}
	// print the user list we received from the server.
	log.Print("\nUser List:\n")
	fmt.Printf("r.GetUsers(): %v\n", r.GetUsers())
}
