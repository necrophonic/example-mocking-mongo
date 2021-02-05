package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/necrophonic/mocking-mongo/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	database   = "my-database"
	collection = "my-collection"
)

// DatabaseDoer is the interface our database client needs to satisfy
type DatabaseDoer interface {
	Insert(context.Context, string, string, interface{}) (interface{}, error)
	Fetch(context.Context, string, string, interface{}, interface{}) error
}

// API ...
type API struct {
	DBClient DatabaseDoer
}

// ConnectDB creates a new database client and adds it to the api
func (a *API) ConnectDB(ctx context.Context, uri string) (err error) {
	cctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	a.DBClient, err = db.New(cctx, uri)
	return
}

// ThingHandler ...
func (a *API) ThingHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// For arbitrary reasons, this handler just sets
	// and then retrieves a record.
	document := bson.D{
		{"_id", "some-id"},
		{"name", "some-name"},
	}

	insertedID, err := a.DBClient.Insert(ctx, database, collection, document)
	if err != nil {
		http.Error(w, "failed to insert", http.StatusInternalServerError)
		return
	}

	var fetched map[string]interface{}
	if err := a.DBClient.Fetch(ctx, database, collection, bson.D{{"_id", insertedID}}, &fetched); err != nil {
		http.Error(w, "failed to fetch", http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(fetched)
	if err != nil {
		http.Error(w, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
