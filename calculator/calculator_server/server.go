package main

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"grpc-go-course/calculator/calculatorpb"
	"log"
	"net"
)

type server struct{}

func (*server) Sum(c context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error){

	fmt.Println("Received Sum RPC %v", req)
	firstNumber := req.GetFirstNumber()
	secondNumber := req.GetSecondNumber()
	sum := firstNumber + secondNumber
	res := &calculatorpb.SumResponse{
		SumResult: sum,
	}
	return res, nil
}

func main() {
	fmt.Println("Calculator Server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
