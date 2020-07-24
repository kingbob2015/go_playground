package main

import (
	"context"
	"fmt"
	"grpc_go_course/calculator/calcpb"
	"io"
	"log"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect %v", err)
	}

	defer cc.Close()

	c := calcpb.NewCalculatorServiceClient(cc)
	// sum(c, 3, 10)
	// pnd(c, 120)
	// avg(c, []int32{1, 2, 3, 4})
	// maxStream(c, []int64{1, 5, 3, 6, 2, 20})

	//Valid no errors
	squareRoot(c, 10)
	//Errors
	squareRoot(c, -5)
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

func avg(c calcpb.CalculatorServiceClient, nums []int32) {
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while calling average: %v", err)
	}
	for _, num := range nums {
		stream.Send(&calcpb.ComputeAverageRequest{
			Num: num,
		})
	}
	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error while receiving response from computer average: %v", err)
	}
	fmt.Printf("Computer Average Response: %v\n", res)
}

func maxStream(c calcpb.CalculatorServiceClient, nums []int64) {
	stream, err := c.FindMaxmimum(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	waitc := make(chan struct{})

	go func() {
		for _, num := range nums {
			stream.Send(&calcpb.FindMaxmimumRequest{
				Num: num,
			})
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				//we've reached the end of the stream (the stream was closed)
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
				break
			}
			fmt.Printf("Received: %v\n", res.GetResult())
		}
		close(waitc)
	}()
	<-waitc
}

func squareRoot(c calcpb.CalculatorServiceClient, num int32) {
	res, err := c.SquareRoot(context.Background(), &calcpb.SquareRootRequest{
		Num: num,
	})
	if err != nil {
		respError, ok := status.FromError(err)
		if ok {
			//actual error from gRPC (user error from server grpc)
			fmt.Println("Error Message From Server: ", respError.Message())
			fmt.Println(respError.Code())
			if respError.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative number")
			}
		} else {
			log.Fatalf("Big system error calling squareroot: %v", err)
		}
		return
	}
	fmt.Printf("Square root of %v is %v\n", num, res.GetResult())
}
