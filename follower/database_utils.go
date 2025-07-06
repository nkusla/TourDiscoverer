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

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
