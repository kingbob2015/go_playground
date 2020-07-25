package main

import (
	"context"
	"fmt"
	"grpc_go_course/blog/blogpb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Blog client")

	opts := grpc.WithInsecure()
	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect %v", err)
	}

	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)
	createPost(c)
}

func createPost(c blogpb.BlogServiceClient) {
	blog := &blogpb.Blog{
		AuthorId: "Bob King",
		Title:    "My First Blog",
		Content:  "Content of the first blog",
	}
	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{
		Blog: blog,
	})
	if err != nil {
		log.Fatalf("Error creating blog: %v", err)
	}
	fmt.Printf("Blog has been created: %v\n", res)
}
