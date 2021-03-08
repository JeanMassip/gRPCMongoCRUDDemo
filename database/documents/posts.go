package documents

import (
	"context"
	"jeanmassip/gRPCMongoCRUDDemo/post/postpb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//postCollection is the name of the collection storing our blog documents within the mongo database
const postCollection = "posts"

//Post defines the structure defining the blog document within our mongo database
type Post struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

//InsertOne inserts one post in the database
func (post *Post) InsertOne(db mongo.Database) (primitive.ObjectID, error) {
	collection := db.Collection(postCollection)
	result, err := collection.InsertOne(context.Background(), post)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

//FindOne returns the post with the specified ID from the database
func (post *Post) FindOne(db mongo.Database) error {
	collection := db.Collection(postCollection)
	filter := bson.M{"_id": post.ID}

	err := collection.FindOne(context.Background(), filter).Decode(post)
	if err != nil {
		return err
	}

	return nil
}

//Find returns a cursor pointin to all the posts in the db
func Find(db mongo.Database) (*mongo.Cursor, error) {
	collection := db.Collection(postCollection)
	return collection.Find(context.Background(), bson.D{{}})
}

//Update updates the specified post within the database
func (post *Post) Update(db mongo.Database) error {
	collection := db.Collection(postCollection)
	update := bson.M{
		"$set": bson.M{
			"author_id": post.AuthorID,
			"title":     post.Title,
			"content":   post.Content,
		},
	}

	_, err := collection.UpdateOne(context.Background(), bson.M{"_id": post.ID}, update)
	if err != nil {
		return err
	}

	return nil
}

//FromPostPB parses a post defined by the protobuff into a mongo post document
func FromPostPB(postProto *postpb.Post) (*Post, error) {
	oid, err := primitive.ObjectIDFromHex(postProto.ID)
	if err != nil {
		return nil, err
	}

	return &Post{
		ID:       oid,
		AuthorID: postProto.AuthorID,
		Title:    postProto.Title,
		Content:  postProto.Content,
	}, nil
}

//ToPostPB parses a mongo post document into a post defined by the protobuff
func (post *Post) ToPostPB() *postpb.Post {
	return &postpb.Post{
		ID:       post.ID.Hex(),
		AuthorID: post.AuthorID,
		Title:    post.Title,
		Content:  post.Content,
	}
}
