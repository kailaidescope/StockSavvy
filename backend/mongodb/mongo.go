package mongodb

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// GetMongoDBInstance connects to a local MongoDB instance on the given port using the provided
// username and password. It uses authSource=admin and SCRAM-SHA-256 by default and returns a
// connected *mongo.Client (or an error).
func GetMongoDBInstance(username, password string, port int) (*mongo.Client, error) {
	// connection timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// escape username/password
	u := url.QueryEscape(username)
	p := url.QueryEscape(password)

	// build URI (authSource=admin by default)
	uri := fmt.Sprintf("mongodb://%s:%s@localhost:%d/?authSource=admin&authMechanism=SCRAM-SHA-256", u, p, port)

	clientOpts := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	// verify connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		// best-effort disconnect on failure
		_ = client.Disconnect(ctx)
		return nil, err
	}

	return client, nil
}
