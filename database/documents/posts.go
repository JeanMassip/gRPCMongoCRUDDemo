package documents

import (
	"context"

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
func (blog *Post) InsertOne(db mongo.Database) (primitive.ObjectID, error) {
	collection := db.Collection(postCollection)
	result, err := collection.InsertOne(context.Background(), blog)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

//FindOne returns the post with the specified ID from the database
func (blog *Post) FindOne(db mongo.Database) error {
	collection := db.Collection(postCollection)
	filter := bson.M{"_id": blog.ID}

	err := collection.FindOne(context.Background(), filter).Decode(blog)
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
func (blog *Post) Update(db mongo.Database) error {
	collection := db.Collection(postCollection)
	update := bson.M{
		"$set": bson.M{
			"author_id": blog.AuthorID,
			"title":     blog.Title,
			"content":   blog.Content,
		},
	}

	_, err := collection.UpdateOne(context.Background(), bson.M{"_id": blog.ID}, update)
	if err != nil {
		return err
	}

	return nil
}
