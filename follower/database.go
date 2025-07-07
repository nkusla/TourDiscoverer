package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Database struct {
	Driver neo4j.DriverWithContext
}

func (db *Database) InitDatabase() error {
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

	db.Driver = driver

	if err := db.createConstraints(); err != nil {
		return errors.New("failed to create constraints: " + err.Error())
	}

	log.Println("Successfully connected to Neo4j")
	return nil
}

func (db *Database) createConstraints() error {
	ctx := context.Background()
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	// Create unique constraint for User username
	_, err := session.Run(ctx, `
        CREATE CONSTRAINT user_username_unique IF NOT EXISTS
        FOR (u:User) REQUIRE u.username IS UNIQUE
    `, nil)
	if err != nil {
		return errors.New("failed to create username constraint: " + err.Error())
	}

	log.Println("Database constraints created successfully")
	return nil
}

func (db *Database) CloseDatabase() error {
	if db.Driver != nil {
		return db.Driver.Close(context.Background())
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
