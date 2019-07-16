package main

import (
	"fmt"
	"google.golang.org/grpc"
	"grpc-go-course/blog/blogpb"
	"log"
	"net"
	"os"
	"os/signal"
)

type server struct {

}

func main() {
	//if we crash the go code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Blog Service Started")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	opts := []grpc.ServerOption{}

	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, &server{})

	go func() {

		fmt.Println("Starting server....")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	//Wait for CTRL+c to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	//Block until a signal is received
	<-ch
	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Closing the listener")
	lis.Close()
	fmt.Println("End of program")
}

