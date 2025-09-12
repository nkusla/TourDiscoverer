package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BlogRepository struct {
	collection *mongo.Collection
}

func (r *BlogRepository) Create(blog *Blog) error {
	fmt.Println("Creating blog:", blog)
	blog.CreatedAt = time.Now().Unix()
	_, err := r.collection.InsertOne(context.Background(), blog)
	return err
}
func (r *BlogRepository) GetAll() ([]Blog, error) {
	cursor, err := r.collection.Find(context.Background(), map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var blogs []Blog
	for cursor.Next(context.Background()) {
		var blog Blog
		if err := cursor.Decode(&blog); err != nil {
			return nil, err
		}
		blogs = append(blogs, blog)
	}
	return blogs, nil
}

func (r *BlogRepository) GetByID(id string) (*Blog, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var blog Blog
	err = r.collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&blog)
	if err != nil {
		return nil, err
	}
	return &blog, nil
}

func (r *BlogRepository) AddLike(blogID, username string) error {
	objectID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return err
	}

	_, err = r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": objectID},
		bson.M{
			"$addToSet": bson.M{"likes": username},
			"$inc":      bson.M{"like_count": 1},
		},
	)
	return err
}

func (r *BlogRepository) RemoveLike(blogID, username string) error {
	objectID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return err
	}

	_, err = r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": objectID},
		bson.M{
			"$pull": bson.M{"likes": username},
			"$inc":  bson.M{"like_count": -1},
		},
	)
	return err
}

func (r *BlogRepository) IsLikedByUser(blogID, username string) (bool, error) {
	objectID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return false, err
	}

	count, err := r.collection.CountDocuments(
		context.Background(),
		bson.M{
			"_id":   objectID,
			"likes": username,
		},
	)
	return count > 0, err
}
