package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/andreasgarvik/grpc_blog/blogpb"
	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to server: %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	blog := &blogpb.Blog{
		AuthorId: "Andreas",
		Title:    "My first blog",
		Content:  "My content",
	}
	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{
		Blog: blog,
	})
	if err != nil {
		log.Fatalf("Could not create blog: %v", err)
	}
	fmt.Printf("Blog created: %v", res)

	_, err = c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: res.GetBlog().GetId(),
	})
	if err != nil {
		log.Fatalf("Error while reading: %v", err)
	}

	uptBlog := &blogpb.Blog{
		Id:       res.GetBlog().GetId(),
		AuthorId: "Changed Auther",
		Title:    "My first blog (edited)",
		Content:  "My content",
	}

	_, err = c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{
		Blog: uptBlog,
	})
	if err != nil {
		log.Fatalf("Error while updating: %v", err)
	}

	_, err = c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{
		BlogId: uptBlog.GetId(),
	})
	if err != nil {
		log.Fatalf("Error while deleting: %v", err)
	}
	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("Error streaming: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while getting blog: %v", err)
		}
		fmt.Println(res.GetBlog())
	}

}
