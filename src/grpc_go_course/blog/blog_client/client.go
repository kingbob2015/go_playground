package main

import (
	"context"
	"fmt"
	"grpc_go_course/blog/blogpb"
	"io"
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
	blogID, err := createPost(c)
	if err != nil {
		fmt.Println("Create post failed, exiting...")
		return
	}
	readPost(c, blogID)
	updateBlog(c, blogID)
	deleteBlog(c, blogID)
	listBlog(c)
}

func createPost(c blogpb.BlogServiceClient) (string, error) {
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
		return "", err
	}
	fmt.Printf("Blog has been created: %v\n", res)
	return res.GetBlog().GetId(), nil
}

func readPost(c blogpb.BlogServiceClient, blogID string) {
	//One that will error, random blog id
	_, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: "5f1ca5bc99c3066789143072",
	})
	if err != nil {
		fmt.Printf("Error happened while reading: %v\n", err)
	}

	//Pulled ID from Mongo, know it exists
	res, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: blogID,
	})
	if err != nil {
		fmt.Printf("Error happened while reading: %v\n", err)
		return
	}

	fmt.Printf("We have read the blog: %v\n", res)
}

func updateBlog(c blogpb.BlogServiceClient, blogID string) {
	blog := &blogpb.Blog{
		Id:       blogID,
		AuthorId: "Bob King",
		Title:    "My Second Blog",
		Content:  "Content of the second blog",
	}
	res, err := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{
		Blog: blog,
	})
	if err != nil {
		log.Fatalf("Error happened while updating: %v", err)
		return
	}
	fmt.Printf("Blog update was successful: %v\n", res)
}

func deleteBlog(c blogpb.BlogServiceClient, blogID string) {
	res, err := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{
		BlogId: blogID,
	})
	if err != nil {
		log.Fatalf("Error happened while deleting: %v", err)
		return
	}
	fmt.Printf("%v was deleted\n", res)
}

func listBlog(c blogpb.BlogServiceClient) {
	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("Error while calling list blog: %v", err)
		return
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happened receiving stream: %v", err)
		}
		fmt.Printf("Blog Entry: %v\n", res)
	}
}
