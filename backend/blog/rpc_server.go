package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

// RPC Request/Response strukture
type CreateBlogRPCRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
	Author      string   `json:"author"`
}

type CreateBlogRPCResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	BlogID  string `json:"blog_id"`
}

type GetPersonalizedBlogsRPCRequest struct {
	Username string `json:"username"`
}

type BlogRPCMessage struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Images        []string `json:"images"`
	Author        string   `json:"author"`
	CreatedAt     int64    `json:"created_at"`
	LikeCount     int      `json:"like_count"`
	IsLikedByUser bool     `json:"is_liked_by_user"`
}

type GetPersonalizedBlogsRPCResponse struct {
	Blogs []BlogRPCMessage `json:"blogs"`
}

// RPC Server
type BlogRPCServer struct {
	service *BlogService
}

func NewBlogRPCServer(service *BlogService) *BlogRPCServer {
	return &BlogRPCServer{service: service}
}

func (s *BlogRPCServer) CreateBlog(ctx context.Context, req *CreateBlogRPCRequest) (*CreateBlogRPCResponse, error) {
	blog := &Blog{
		Title:       req.Title,
		Description: req.Description,
		Images:      req.Images,
		Author:      req.Author,
		CreatedAt:   time.Now().Unix(),
		LikeCount:   0,
		Likes:       []string{},
	}

	err := s.service.CreateBlog(blog)
	if err != nil {
		return &CreateBlogRPCResponse{
			Success: false,
			Message: "Failed to create blog",
			BlogID:  "",
		}, err
	}

	return &CreateBlogRPCResponse{
		Success: true,
		Message: "Blog created successfully",
		BlogID:  blog.ID,
	}, nil
}

func (s *BlogRPCServer) GetPersonalizedBlogs(ctx context.Context, req *GetPersonalizedBlogsRPCRequest) (*GetPersonalizedBlogsRPCResponse, error) {
	blogs, err := s.service.GetAllBlogs(req.Username)
	if err != nil {
		return nil, err
	}

	var rpcBlogs []BlogRPCMessage
	for _, blog := range blogs {
		rpcBlogs = append(rpcBlogs, BlogRPCMessage{
			ID:            blog.ID,
			Title:         blog.Title,
			Description:   blog.Description,
			Images:        blog.Images,
			Author:        blog.Author,
			CreatedAt:     blog.CreatedAt,
			LikeCount:     blog.LikeCount,
			IsLikedByUser: blog.IsLikedByUser,
		})
	}

	return &GetPersonalizedBlogsRPCResponse{
		Blogs: rpcBlogs,
	}, nil
}

// TCP RPC Server implementacija
func (s *BlogRPCServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	for {
		var rpcCall map[string]interface{}
		if err := decoder.Decode(&rpcCall); err != nil {
			log.Printf("Error decoding RPC call: %v", err)
			return
		}

		method, ok := rpcCall["method"].(string)
		if !ok {
			log.Printf("Missing or invalid method in RPC call")
			continue
		}

		switch method {
		case "CreateBlog":
			var req CreateBlogRPCRequest
			paramsBytes, _ := json.Marshal(rpcCall["params"])
			json.Unmarshal(paramsBytes, &req)

			resp, err := s.CreateBlog(context.Background(), &req)
			if err != nil {
				encoder.Encode(map[string]interface{}{
					"error": err.Error(),
				})
			} else {
				encoder.Encode(map[string]interface{}{
					"result": resp,
				})
			}

		case "GetPersonalizedBlogs":
			var req GetPersonalizedBlogsRPCRequest
			paramsBytes, _ := json.Marshal(rpcCall["params"])
			json.Unmarshal(paramsBytes, &req)

			resp, err := s.GetPersonalizedBlogs(context.Background(), &req)
			if err != nil {
				encoder.Encode(map[string]interface{}{
					"error": err.Error(),
				})
			} else {
				encoder.Encode(map[string]interface{}{
					"result": resp,
				})
			}

		default:
			encoder.Encode(map[string]interface{}{
				"error": fmt.Sprintf("Unknown method: %s", method),
			})
		}
	}
}

func (s *BlogRPCServer) StartRPCServer(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to start RPC server: %v", err)
	}
	defer listener.Close()

	log.Printf("Blog RPC server listening on port %s", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go s.handleConnection(conn)
	}
}
