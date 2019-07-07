package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc-go-course/calculator/calculatorpb"
	"io"
	"log"
	"time"
)

func main() {

	fmt.Println("Calculator client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Couldn't connect: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)
	//fmt.Printf("Created client: %f", c)

	doUnary(c)

	doServerStreaming(c)

	doClientStreaming(c)

	doBiDirectionalStreaming(c)

	doErrorUnary(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Start to do a Sum Unary RPC...")
	req := &calculatorpb.SumRequest{
		FirstNumber: 5,
		SecondNumber: 40,
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	log.Println("Response from Greet: ", res.SumResult)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Start to do a Prime Decomposition Server Streaming RPC...")
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 2100,
	}
	resStream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Prime Decomposition RPC: %v", err)
	}

	for {
		res, err := resStream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("something happened: %v", err)
		}
		fmt.Println(res.GetPrimeFactor())
	}
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Start to do a Compute Average Client Streaming RPC...")

	resStream, err := c.ComputeAverage(context.Background())

	if err != nil {
		log.Fatalf("Error while opening stream: %v", err)
	}

	numbers := []int32{ 3, 5, 9, 27, 35 }

	for _, number := range numbers {
		fmt.Println("Sending number: ", number)
		resStream.Send(&calculatorpb.ComputeAverageRequest{
			Number: number,
		})
	}

	res, err := resStream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error while receiving response: %v", err)
	}

	fmt.Printf("The average is: %v\n", res.GetAverage())
}

func doBiDirectionalStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Start to do a Find Maximum BiDirectional Streaming RPC...")

	stream, err := c.FindMaximum(context.Background())

	if err != nil {
		log.Fatalf("Error while opening stream and calling FindMaximum: %v", err)
	}

	waitC := make(chan struct{})

	// send go routine

	go func() {

		numbers := []int32{4, 7, 11, 9, 2, 6, 33}

		for _, number := range numbers {

			fmt.Println("Sending number: ", number)
			stream.Send(&calculatorpb.FindMaximumRequest{
				Number: number,
			})

			time.Sleep(1000 * time.Millisecond)
		}

		stream.CloseSend()
	}()
	// receive go routine

	go func() {

		for {

			res, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Problem while reading server stream: %v", err)
				break
			}

			maximum := res.GetMaximumNumber()
			fmt.Println("Received a new maximum number of ...: ", maximum)
		}

		close(waitC)
	}()

	<-waitC
}

func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Start to do SquareRoot Unary RPC...")


	//correct call
	doErrorCall(c, 10)

	//error call
	doErrorCall(c, -2)

}

func doErrorCall(c calculatorpb.CalculatorServiceClient, n int32) {
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: n})

	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			//actual error from gRPC (user error)
			fmt.Printf("Error message from server : %v\n", respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We Probably sent a negative number!")
			}
		}else{
			log.Fatalf("Big error calling SquareRoot: %v", err)
		}
		return
	}
	fmt.Printf("Result of square root of %v: %v\n", n, res.GetNumberRoot())
}