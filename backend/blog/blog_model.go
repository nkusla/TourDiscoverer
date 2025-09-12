package main

type Blog struct {
	ID          string   `json:"id" bson:"_id,omitempty"`
	Title       string   `json:"title" bson:"title"`
	Description string   `json:"description" bson:"description"` // markdown podrška
	CreatedAt   int64    `json:"created_at" bson:"created_at"`
	Images      []string `json:"images,omitempty" bson:"images,omitempty"`
	Likes       []string `json:"likes,omitempty" bson:"likes,omitempty"` // lista username-ova koji su lajkovali
	LikeCount   int      `json:"like_count" bson:"like_count"`           // broj lajkova
}

type LikeRequest struct {
	BlogID string `json:"blog_id"`
}
