// Package db provides datastore access functions
package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client is a resuable client for performing database operations
type Client struct {
	c *mongo.Client
}

// New creates and connects a new database client
func New(ctx context.Context, uri string) (client *Client, err error) {
	client = &Client{}

	client.c, err = mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err = client.c.Connect(ctx); err != nil {
		return nil, err
	}

	return client, nil
}

// Insert will insert a single document into the given database and collection
func (db *Client) Insert(ctx context.Context, database, collection string, document interface{}) (interface{}, error) {
	res, err := db.c.Database(database).Collection(collection).InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}
	return res.InsertedID, nil
}

// Fetch attempts to fetch a single record from the database. The result (if any)
// will be unmarshaled into the given v.
func (db *Client) Fetch(ctx context.Context, database, collection string, filter, v interface{}) error {
	res := db.c.Database(database).Collection(collection).FindOne(ctx, filter)
	if res == nil {
		// Might want to do something if nothing is returned,
		// but here we just leave the v as nil and return.
		return nil
	}

	if err := res.Decode(&v); err != nil {
		return err
	}
	return nil
}

// Other functions for setting / retrieving etc as required
// ...
// ...
