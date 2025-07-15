package main

import (
	"context"
	"fmt"
	"time"

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
