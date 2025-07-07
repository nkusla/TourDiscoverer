package main

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type FollowerRepository struct {
	db *Database
}

func (r *FollowerRepository) CreateUser(username string) error {
	ctx := context.Background()
	session := r.db.Driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})

	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := "MERGE (u:User {username: $username})"
		params := map[string]any{"username": username}
		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}
		return result.Consume(ctx)
	})

	return err
}

func (r *FollowerRepository) CreateFollowRelationship(follower string, followee string) error {
	ctx := context.Background()
	session := r.db.Driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})

	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
            MERGE (f:User {username: $follower})
            MERGE (t:User {username: $followee})
            MERGE (f)-[:FOLLOWS]->(t)
        `
		params := map[string]any{"follower": follower, "followee": followee}
		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}
		return result.Consume(ctx)
	})

	return err
}

func (r *FollowerRepository) DeleteFollowRelationship(follower string, followee string) error {
	ctx := context.Background()
	session := r.db.Driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})

	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
						MATCH (f:User {username: $follower})-[r:FOLLOWS]->(t:User {username: $followee})
						DELETE r
				`
		params := map[string]any{"follower": follower, "followee": followee}
		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}
		return result.Consume(ctx)
	})

	return err
}
