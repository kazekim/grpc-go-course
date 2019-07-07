package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc-go-course/greet/greetpb"
	"io"
	"log"
	"time"
)

func main() {

	fmt.Println("Hello, I'm a client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Couldn't connect: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	//fmt.Printf("Created client: %f", c)

	doUnary(c)
	doServerStream(c)
	doClientStream(c)
	doBidirectionalStream(c)
	doUnaryWithDeadline(c, 5 * time.Second) // should complete
	doUnaryWithDeadline(c, 1 * time.Second) // should timeout
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Start to do a Unary RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Jirawat",
			LastName: "Harnsiriwatanakit",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	log.Println("Response from Greet: ", res.Result)
}

func doServerStream(c greetpb.GreetServiceClient) {
	fmt.Println("Start to do a Server Streaming RPC...")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Jirawat",
			LastName: "Harnsiriwatanakit",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we've reached the end of stream
			break
		}

		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Println("Response from GreetManyTimes: ", msg.GetResult())
	}
}

func doClientStream(c greetpb.GreetServiceClient) {
	fmt.Println("starting to do a Client Streaming RPC...")

	requests := []*greetpb.LongGreetRequest {
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Kim",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Papure",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Oat",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Amy",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Max",
			},
		},
	}

	resStream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while calling LongGreet : %v", err)
	}

	//we iterate over our slice and send each message individually
	for _, req := range requests {
		fmt.Printf("Sending request : %v\n", req)
		resStream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	res, err := resStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from LongGreet; %v", err)
	}

	fmt.Printf("LongGreet Response: %v\n", res)
}

func doBidirectionalStream(c greetpb.GreetServiceClient) {
	fmt.Println("starting to do a Bidirectional Streaming RPC...")

	requests := []*greetpb.GreetEveryoneRequest {
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Kim",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Papure",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Oat",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Amy",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Max",
			},
		},
	}

	// We create a stream by invoking the client

	stream, err := c.GreetEveryone(context.Background())

	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	waitC := make(chan struct{})

	// we send a bunch of messages to the client (go routine)

	go func() {

		// function to send a bunch of messages

		for _, req := range requests {
			fmt.Println("Send message : ", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}

		stream.CloseSend()
	}()

	// we receive a bunch of messages from the client (go routine)

	go func() {

		// function to receive a bunch of messages

		for {
			res, err := stream.Recv()

			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
				close(waitC)
			}

			fmt.Println("Received: ", res.GetResult())
		}


		close(waitC)
	}()
	// block until everything is done

	<-waitC
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	fmt.Println("Start to do a Unary With Deadline RPC...")
	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Jirawat",
			LastName: "Harnsiriwatanakit",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {

		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit! Deadline was exceeded")
			}else{
				fmt.Printf("Unexpected error: %v\n", statusErr)
			}
		}else{
			log.Fatalf("error while calling Greet With Deadline RPC: %v", err)
		}

		return
	}
	log.Println("Response from Greet With Deadline: ", res.Result)
}