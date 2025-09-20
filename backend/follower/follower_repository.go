package main

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type FollowerRepository struct {
	db *Database
}

func (r *FollowerRepository) CreateUser(username string, role string) error {
	ctx := context.Background()
	session := r.db.Driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})

	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := "MERGE (u:User {username: $username}) SET u.role = $role"
		params := map[string]any{"username": username, "role": role}
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

func (r *FollowerRepository) GetFollowers(username string) ([]User, error) {
	ctx := context.Background()
	session := r.db.Driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})

	defer session.Close(ctx)

	followers := []User{}

	_, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (u:User {username: $username})<-[:FOLLOWS]-(f:User)
			RETURN f.username AS follower, f.role AS role
    `
		params := map[string]any{"username": username}
		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		for result.Next(ctx) {
			record := result.Record()
			followerUsername, _ := record.Get("follower")
			followerRole, _ := record.Get("role")

			role := ""
			if followerRole != nil {
				role = followerRole.(string)
			}

			followers = append(followers, User{
				Username: followerUsername.(string),
				Role:     role,
			})
		}

		return nil, result.Err()
	})

	return followers, err
}

func (r *FollowerRepository) GetFollowing(username string) ([]User, error) {
	ctx := context.Background()
	session := r.db.Driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})

	defer session.Close(ctx)

	var following []User

	_, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (u:User {username: $username})-[:FOLLOWS]->(f:User)
			RETURN f.username AS followee, f.role AS role
		`
		params := map[string]any{"username": username}
		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		for result.Next(ctx) {
			record := result.Record()
			followeeUsername, _ := record.Get("followee")
			followeeRole, _ := record.Get("role")

			role := ""
			if followeeRole != nil {
				role = followeeRole.(string)
			}

			following = append(following, User{
				Username: followeeUsername.(string),
				Role:     role,
			})
		}

		return nil, result.Err()
	})

	return following, err
}

func (r *FollowerRepository) IsFollowing(follower string, followee string) (bool, error) {
	ctx := context.Background()
	session := r.db.Driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})

	defer session.Close(ctx)

	var exists bool

	_, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
            MATCH (f:User {username: $follower})-[:FOLLOWS]->(t:User {username: $followee})
            RETURN COUNT(*) > 0 AS exists
        `
		params := map[string]any{"follower": follower, "followee": followee}
		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			record := result.Record()
			existsValue, _ := record.Get("exists")
			exists = existsValue.(bool)
		}

		return nil, result.Err()
	})

	return exists, err
}

func (r *FollowerRepository) GetRecommendations(username string) ([]User, error) {
	ctx := context.Background()
	session := r.db.Driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})

	defer session.Close(ctx)

	var recommendations []User

	_, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (me:User {username: $username})-[:FOLLOWS]->(friend:User)-[:FOLLOWS]->(recommendation:User)
			WHERE NOT (me)-[:FOLLOWS]->(recommendation) AND recommendation.username <> $username
			RETURN DISTINCT recommendation.username AS username, recommendation.role AS role, COUNT(*) AS mutualConnections
			ORDER BY mutualConnections DESC
			LIMIT 10
		`
		params := map[string]any{"username": username}
		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		for result.Next(ctx) {
			record := result.Record()
			recommendationUsername, _ := record.Get("username")
			recommendationRole, _ := record.Get("role")

			role := ""
			if recommendationRole != nil {
				role = recommendationRole.(string)
			}

			recommendations = append(recommendations, User{
				Username: recommendationUsername.(string),
				Role:     role,
			})
		}

		return nil, result.Err()
	})

	return recommendations, err
}
