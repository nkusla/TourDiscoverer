package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

// RPC Request/Response strukture
type CheckoutRPCRequest struct {
	Username string `json:"username"`
}

type CheckoutRPCResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Tokens  []string `json:"tokens"`
}

type GetPurchasedToursRPCRequest struct {
	Username string `json:"username"`
}

type TokenRPCMessage struct {
	Token       string `json:"token"`
	TourID      uint   `json:"tour_id"`
	TourName    string `json:"tour_name"`
	PurchasedAt string `json:"purchased_at"`
}

type GetPurchasedToursRPCResponse struct {
	Tokens []TokenRPCMessage `json:"tokens"`
}

// RPC Server
type PurchaseRPCServer struct {
	service *PurchaseService
}

func NewPurchaseRPCServer(service *PurchaseService) *PurchaseRPCServer {
	return &PurchaseRPCServer{service: service}
}

func (s *PurchaseRPCServer) Checkout(ctx context.Context, req *CheckoutRPCRequest) (*CheckoutRPCResponse, error) {
	tokens, err := s.service.Checkout(req.Username)
	if err != nil {
		return &CheckoutRPCResponse{
			Success: false,
			Message: "Failed to checkout",
			Tokens:  []string{},
		}, err
	}

	var tokenStrings []string
	for _, token := range tokens {
		tokenStrings = append(tokenStrings, token.Token)
	}

	return &CheckoutRPCResponse{
		Success: true,
		Message: "Checkout successful",
		Tokens:  tokenStrings,
	}, nil
}

func (s *PurchaseRPCServer) GetPurchasedTours(ctx context.Context, req *GetPurchasedToursRPCRequest) (*GetPurchasedToursRPCResponse, error) {
	tokens, err := s.service.GetUserTokens(req.Username)
	if err != nil {
		return nil, err
	}

	var rpcTokens []TokenRPCMessage
	for _, token := range tokens {
		rpcTokens = append(rpcTokens, TokenRPCMessage{
			Token:       token.Token,
			TourID:      token.TourID,
			TourName:    token.TourName,
			PurchasedAt: token.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &GetPurchasedToursRPCResponse{
		Tokens: rpcTokens,
	}, nil
}

// TCP RPC Server implementacija
func (s *PurchaseRPCServer) handleConnection(conn net.Conn) {
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
		case "Checkout":
			var req CheckoutRPCRequest
			paramsBytes, _ := json.Marshal(rpcCall["params"])
			json.Unmarshal(paramsBytes, &req)

			resp, err := s.Checkout(context.Background(), &req)
			if err != nil {
				encoder.Encode(map[string]interface{}{
					"error": err.Error(),
				})
			} else {
				encoder.Encode(map[string]interface{}{
					"result": resp,
				})
			}

		case "GetPurchasedTours":
			var req GetPurchasedToursRPCRequest
			paramsBytes, _ := json.Marshal(rpcCall["params"])
			json.Unmarshal(paramsBytes, &req)

			resp, err := s.GetPurchasedTours(context.Background(), &req)
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

func (s *PurchaseRPCServer) StartRPCServer(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to start RPC server: %v", err)
	}
	defer listener.Close()

	log.Printf("Purchase RPC server listening on port %s", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go s.handleConnection(conn)
	}
}