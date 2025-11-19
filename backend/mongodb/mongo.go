package mongodb

import (
	"context"
	"errors"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// GetMongoDBInstance connects to a local MongoDB instance on the given port using the provided
// username and password. It uses authSource=admin and SCRAM-SHA-256 by default and returns a
// connected *mongo.Client (or an error).
func GetMongoDBInstance(username, password, host string, port int) (*mongo.Client, error) {
	// connection timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// build URI (authSource=admin by default)
	uri := "mongodb://" + host + ":" + strconv.Itoa(port)

	clientOpts := options.Client().ApplyURI(uri)
	if username != "" && password != "" {
		cred := options.Credential{
			Username: username,
			Password: password,
		}
		clientOpts = clientOpts.SetAuth(cred)
	} else {
		return nil, errors.New("username and password not found when creating MongoDB client")
	}

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
