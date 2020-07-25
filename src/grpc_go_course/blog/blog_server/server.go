package main

import (
	"context"
	"fmt"
	"grpc_go_course/blog/blogpb"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"google.golang.org/grpc/codes"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var collection *mongo.Collection

type blogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

type server struct{}

func (*server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	blog := req.GetBlog()

	data := blogItem{
		AuthorID: blog.GetAuthorId(),
		Content:  blog.GetContent(),
		Title:    blog.GetTitle(),
	}

	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		log.Fatalf("Error inserting blog record: %v", err)
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unable to insert data in MongoDB: %v", err),
		)
	}
	oid, ok := res.InsertedID.(primitive.ObjectID) //Cast insert one result to object id (checks the underlying type!)
	if !ok {
		log.Fatalf("Cannot convert to oid: %v", err)
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot convert to oid: %v", err),
		)
	}

	blog.Id = oid.Hex()
	return &blogpb.CreateBlogResponse{
		Blog: blog,
	}, nil
}

func main() {
	// if we crash the go code, we get the file name and line number in error message when we use log
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//Connect to MongoDB
	fmt.Println("Connecting to MongoDB")
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Unable to create new mongo client: %v", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Unable to connect to new mongo client: %v", err)
	}

	collection = client.Database("mydb").Collection("blog")

	fmt.Println("Blog Service Started")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, &server{})
	reflection.Register(s)

	go func() {
		fmt.Println("Starting Server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve %v", err)
		}
	}()

	//Make a channel that takes in when ctrl+c is hit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	//Block until a signal is received
	<-ch
	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Closing the listener")
	lis.Close()
	fmt.Println("Closing MongoDB Connection")
	client.Disconnect(context.Background())
	fmt.Println("End of Program")
}
