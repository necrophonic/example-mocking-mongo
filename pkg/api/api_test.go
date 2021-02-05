package api_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/necrophonic/mocking-mongo/pkg/api"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

type MockDB struct {
	lastInsertedDoc interface{}
}

type insertReturn struct {
	InsertedID interface{}
}

func (m *MockDB) Insert(ctx context.Context, dbn, cn string, doc interface{}) (interface{}, error) {
	// Store the doc we inserted so we can test if it happened
	m.lastInsertedDoc = doc
	return insertReturn{
		InsertedID: "some-id", // TODO need to look at a more "real" implementation
	}, nil
}

func (m *MockDB) Fetch(ctx context.Context, dbn, cn string, filter, v interface{}) error {
	// Marshal some fake data into the provided interface
	// ...
	return nil
}

func TestThingHandler(t *testing.T) {

	// Instantiate an API struct
	a := &api.API{}

	// Now instead of calling Connect() we set the
	// DBClient to be our mocked client directly.
	dbc := &MockDB{}
	a.DBClient = dbc

	// Now we can do tests on the handler as normal
	req, _ := http.NewRequest("GET", "https://localhost", nil)
	rr := httptest.NewRecorder()
	a.ThingHandler(rr, req)

	expectedInserted := bson.D{
		{Key: "_id", Value: "some-id"},
		{Key: "name", Value: "some-name"},
	}

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, expectedInserted, dbc.lastInsertedDoc)
}
