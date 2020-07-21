package main

import (
	"context"
	"grpc_go_course/calculator/calcpb"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect %v", err)
	}

	defer cc.Close()

	c := calcpb.NewCalculatorServiceClient(cc)
	sum(c, 3, 10)
	pnd(c, 120)
}

func sum(c calcpb.CalculatorServiceClient, num1, num2 int64) {
	req := &calcpb.SumRequest{
		Num_1: num1,
		Num_2: num2,
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Sum RPC: %v\n", err)
	}
	log.Printf("Response from Sum: %v", res.Result)
}

func pnd(c calcpb.CalculatorServiceClient, num int32) {
	req := &calcpb.PrimeNumberDecompositionRequest{
		Num: num,
	}
	resStream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling PND: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			//we've reached the end of the stream (the stream was closed)
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}
		log.Printf("Response from PND: %v", msg.GetResult())
	}
}
