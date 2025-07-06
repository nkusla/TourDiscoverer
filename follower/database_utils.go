package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var Driver neo4j.DriverWithContext

func InitDatabase() error {
	url := GetEnv("FOLLOWER_DB_URL", "bolt://localhost:7687")
	username := GetEnv("FOLLOWER_DB_USER", "neo4j")
	password := GetEnv("FOLLOWER_DB_PASSWORD", "password")

	driver, err := neo4j.NewDriverWithContext(url, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return errors.New("failed to create Neo4j driver: " + err.Error())
	}

	// Verify connectivity
	ctx := context.Background()
	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		return errors.New("failed to connect to Neo4j: " + err.Error())
	}

	Driver = driver
	log.Println("Successfully connected to Neo4j")
	return nil
}

func CloseDatabase() error {
	if Driver != nil {
		return Driver.Close(context.Background())
	}
	return nil
}

func SeedDatabase() error {
	ctx := context.Background()
	session := Driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})

	defer session.Close(ctx)

	// Example seed data
	executeSeedTransaction(ctx, session, "tourist123")
	executeSeedTransaction(ctx, session, "guide123")

	log.Println("Database seeded successfully")
	return nil
}

func executeSeedTransaction(ctx context.Context, session neo4j.SessionWithContext, username string) error {
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := "CREATE (u:User {username: $username})"
		params := map[string]any{"username": username}
		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}
		return result.Consume(ctx)
	})
	return err
}

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
