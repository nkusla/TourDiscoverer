package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

// RPC Request/Response strukture
type GetProfileRPCRequest struct {
	Username string `json:"username"`
}

type GetProfileRPCResponse struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Email    string `json:"email"`
	UserRole string `json:"user_role"`
}

type GetRecommendationsRPCRequest struct {
	Username string `json:"username"`
}

type RecommendationRPCMessage struct {
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	EntityID    uint   `json:"entity_id"`
}

type GetRecommendationsRPCResponse struct {
	Recommendations []RecommendationRPCMessage `json:"recommendations"`
}

// RPC Server
type StakeholderRPCServer struct {
	service *StakeholderService
}

func NewStakeholderRPCServer(service *StakeholderService) *StakeholderRPCServer {
	return &StakeholderRPCServer{service: service}
}

func (s *StakeholderRPCServer) GetProfile(ctx context.Context, req *GetProfileRPCRequest) (*GetProfileRPCResponse, error) {
	stakeholder, err := s.service.GetStakeholderProfile(req.Username)
	if err != nil {
		return nil, err
	}

	return &GetProfileRPCResponse{
		Username: stakeholder.Username,
		Name:     stakeholder.FirstName,
		Surname:  stakeholder.LastName,
		Email:    "",
		UserRole: "",
	}, nil
}

func (s *StakeholderRPCServer) GetRecommendations(ctx context.Context, req *GetRecommendationsRPCRequest) (*GetRecommendationsRPCResponse, error) {
	// Mock implementacija preporuka
	recommendations := []RecommendationRPCMessage{
		{
			Type:        "tour",
			Title:       "Recommended Tour",
			Description: "Based on your preferences",
			EntityID:    1,
		},
		{
			Type:        "blog",
			Title:       "Trending Blog",
			Description: "Popular in your area",
			EntityID:    2,
		},
	}

	return &GetRecommendationsRPCResponse{
		Recommendations: recommendations,
	}, nil
}

// TCP RPC Server implementacija
func (s *StakeholderRPCServer) handleConnection(conn net.Conn) {
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
		case "GetProfile":
			var req GetProfileRPCRequest
			paramsBytes, _ := json.Marshal(rpcCall["params"])
			json.Unmarshal(paramsBytes, &req)

			resp, err := s.GetProfile(context.Background(), &req)
			if err != nil {
				encoder.Encode(map[string]interface{}{
					"error": err.Error(),
				})
			} else {
				encoder.Encode(map[string]interface{}{
					"result": resp,
				})
			}

		case "GetRecommendations":
			var req GetRecommendationsRPCRequest
			paramsBytes, _ := json.Marshal(rpcCall["params"])
			json.Unmarshal(paramsBytes, &req)

			resp, err := s.GetRecommendations(context.Background(), &req)
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

func (s *StakeholderRPCServer) StartRPCServer(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to start RPC server: %v", err)
	}
	defer listener.Close()

	log.Printf("Stakeholder RPC server listening on port %s", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go s.handleConnection(conn)
	}
}