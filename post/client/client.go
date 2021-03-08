package main

import (
	"context"
	"fmt"
	"jeanmassip/gRPCMongoCRUDDemo/post/postpb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to connect to server : %v", err)
	}

	postClient := postpb.NewPostsClient(cc)
	res, err := postClient.AddPost(context.Background(), &postpb.Post{
		ID:       "",
		AuthorID: "Jean",
		Title:    "Mon premier post",
		Content:  "Youpi !",
	})
	if err != nil {
		log.Fatalf("Unable to create Post : %v", err)
	}

	fmt.Printf("Created Post : %v\n", res)
}
