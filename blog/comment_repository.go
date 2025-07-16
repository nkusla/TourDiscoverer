package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type CommentRepository struct {
	collection *mongo.Collection
}

func (r *CommentRepository) Create(comment *Comment) error {
	comment.CreatedAt = time.Now().Unix()
	comment.UpdatedAt = comment.CreatedAt
	_, err := r.collection.InsertOne(context.Background(), comment)
	return err
}

func (r *CommentRepository) GetByBlogID(blogID string) ([]Comment, error) {
	cursor, err := r.collection.Find(context.Background(), map[string]interface{}{"blog_id": blogID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var comments []Comment
	for cursor.Next(context.Background()) {
		var comment Comment
		if err := cursor.Decode(&comment); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
