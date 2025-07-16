package main

type Blog struct {
    ID          string   `json:"id" bson:"_id,omitempty"`
    Title       string   `json:"title" bson:"title"`
    Description string   `json:"description" bson:"description"` // markdown podrška
    CreatedAt   int64    `json:"created_at" bson:"created_at"`
    Images      []string `json:"images,omitempty" bson:"images,omitempty"`
}