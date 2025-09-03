package main

// Comment represents a comment on a blog post
type Comment struct {
	ID        string `json:"id" bson:"_id,omitempty"`
	BlogID    string `json:"blog_id" bson:"blog_id"`
	Author    string `json:"author" bson:"author"`
	CreatedAt int64  `json:"created_at" bson:"created_at"`
	Text      string `json:"text" bson:"text"`
	UpdatedAt int64  `json:"updated_at" bson:"updated_at"`
}
