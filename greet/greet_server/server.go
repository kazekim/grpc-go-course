package main

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"grpc-go-course/greet/greetpbreetpb"
	"log"
	"net"
)

type server struct{}

func (*server) Greet(c context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error){

	fmt.Println("Greet function was invoked with %v", req)
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()
	result := "Hello, " + firstName + " " + lastName
	res := &greetpb.GreetResponse{
		Result:result,
	}
	return res, nil
}

func main() {
	fmt.Println("Hello world")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
