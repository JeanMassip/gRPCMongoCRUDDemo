package main

import (
	"context"
	"fmt"
	"jeanmassip/gRPCMongoCRUDDemo/database"
	"jeanmassip/gRPCMongoCRUDDemo/database/documents"
	"jeanmassip/gRPCMongoCRUDDemo/post/postpb"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	db *database.Mongo
	postpb.UnimplementedPostsServer
}

func (server *server) AddPost(ctx context.Context, req *postpb.Post) (*postpb.Post, error) {
	mongoPost := documents.Post{
		AuthorID: req.AuthorID,
		Title:    req.Title,
		Content:  req.Content,
	}

	oid, err := mongoPost.InsertOne(*server.db.Database)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Unable to process request : %v", err))
	}
	mongoPost.ID = oid

	return mongoPost.ToPostPB(), nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Error while creating listener : %v", err)
	}

	postServer := &server{
		db: database.NewMongoConnection(),
	}

	err = postServer.db.ConnectToDB("mongodb+srv://<user>:<passwd>@<server>", "posts")
	if err != nil {
		log.Fatalf("Unable to connect to db : %v", err)
	}
	gRPCServer := grpc.NewServer()
	postpb.RegisterPostsServer(gRPCServer, postServer)

	if err := gRPCServer.Serve(lis); err != nil {
		log.Fatalf("Error while serving : %v", err)
	}
}
