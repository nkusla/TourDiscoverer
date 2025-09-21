package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/rs/cors"
)

func main() {
    // Initialize database
    db, err := InitDatabase()
    if err != nil {
        log.Fatal("Failed to initialize database:", err)
    }
    
    // Initialize repositories
    cartRepo := NewCartRepository(db)
    purchaseRepo := NewPurchaseRepository(db)
    
    // Initialize services
    cartService := NewCartService(cartRepo)
    purchaseService := NewPurchaseService(purchaseRepo, cartRepo)
    
    // Initialize handlers
    cartHandler := NewCartHandler(cartService)
    purchaseHandler := NewPurchaseHandler(purchaseService)
    
    // Setup router
    router := mux.NewRouter()
    
    // ========== CART ROUTES ==========
    // Gateway: /api/purchases/cart/* -> strips /api/purchases -> /cart/*
    router.HandleFunc("/cart", cartHandler.GetCart).Methods("GET")                    // /api/purchases/cart
    router.HandleFunc("/cart/items", cartHandler.AddToCart).Methods("POST")           // /api/purchases/cart/items
    router.HandleFunc("/cart/items/{tourId}", cartHandler.RemoveFromCart).Methods("DELETE") // /api/purchases/cart/items/{tourId}
    router.HandleFunc("/cart", cartHandler.ClearCart).Methods("DELETE")              // /api/purchases/cart
    router.HandleFunc("/cart/checkout", purchaseHandler.Checkout).Methods("POST")     // /api/purchases/cart/checkout
    
    // ========== PURCHASE ROUTES ==========  
    // Gateway: /api/purchases/* -> strips /api/purchases -> /*
    router.HandleFunc("/tokens", purchaseHandler.GetUserTokens).Methods("GET")         // /api/purchases/tokens
    router.HandleFunc("/tokens/{token}", purchaseHandler.GetTokenDetails).Methods("GET") // /api/purchases/tokens/{token}
    router.HandleFunc("/validate/{tourId}", purchaseHandler.ValidateAccess).Methods("GET") // /api/purchases/validate/{tourId}
    
    // Health check
    router.HandleFunc("/ping", purchaseHandler.Ping).Methods("GET")
    
    // Setup CORS
    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"*"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"*"},
        AllowCredentials: true,
    })
    
    // Apply CORS middleware
    handler := c.Handler(router)
    
    // Get port from environment or default
    port := GetEnv("PORT", "8084")
    log.Printf("🚀 Purchase service starting on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, handler))
}