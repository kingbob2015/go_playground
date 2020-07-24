package main

import (
	"context"
	"fmt"
	"grpc_go_course/greet/greetpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I am a client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect %v", err)
	}

	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	//doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	doBiDiStreaming(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Bob",
			LastName:  "King",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC: %v\n", err)
	}
	log.Printf("Response from Greet: %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Bob",
			LastName:  "King",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling server stream greet many times RPC: %v", err)
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
		log.Printf("Response from Greet Many Times: %v", msg.GetResult())
	}

}

func doClientStreaming(c greetpb.GreetServiceClient) {
	//Since its a stream this time we do not have to pass in a request because we will push on it
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Err while calling LongGreet: %v", err)
	}
	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Bob",
				LastName:  "King",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Riley",
				LastName:  "Hammer",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Tim",
				LastName:  "King",
			},
		},
	}

	//Iterate over slice and send each message individually
	for _, req := range requests {
		stream.Send(req)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error while receiving response from Long Greet: %v", err)
	}
	fmt.Printf("Long Greet Response: %v\n", res)
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	// we create a stream by invoking client
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Bob",
				LastName:  "King",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Riley",
				LastName:  "Hammer",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Tim",
				LastName:  "King",
			},
		},
	}

	waitc := make(chan struct{})
	// we send a bunch of messages to the server
	go func() {
		//function to send a bunch of messages
		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// we receive a bunch of messages
	go func() {
		for {
			//function to receive a bunch of messages
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

	//block until everything is done
	<-waitc
}
