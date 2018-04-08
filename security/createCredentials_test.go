package security_test

import (
	"encoding/json"
	"testing"

	"github.com/kuzzleio/sdk-go/internal"
	"github.com/kuzzleio/sdk-go/kuzzle"
	"github.com/kuzzleio/sdk-go/types"
	"github.com/stretchr/testify/assert"
)

func TestCreateCredentialsStrategyNull(t *testing.T) {
	k, _ := kuzzle.NewKuzzle(&internal.MockedConnection{}, nil)
	_, err := k.Security.CreateCredentials("", "userid", []byte(`{"body": "test"}`), nil)
	assert.Error(t, err)
}

func TestCreateCredentialsIDNull(t *testing.T) {
	k, _ := kuzzle.NewKuzzle(&internal.MockedConnection{}, nil)
	_, err := k.Security.CreateCredentials("strategy", "", []byte(`{"body": "test"}`), nil)
	assert.Error(t, err)
}

func TestCreateCredentialsBodyNull(t *testing.T) {
	k, _ := kuzzle.NewKuzzle(&internal.MockedConnection{}, nil)
	_, err := k.Security.CreateCredentials("strategy", "userId", nil, nil)
	assert.Error(t, err)
}

func TestCreateCredentialsError(t *testing.T) {
	c := &internal.MockedConnection{
		MockSend: func(query []byte, options types.QueryOptions) *types.KuzzleResponse {
			return &types.KuzzleResponse{Error: &types.KuzzleError{Message: "Unit test error"}}
		},
	}
	k, _ := kuzzle.NewKuzzle(c, nil)

	_, err := k.Security.CreateCredentials("strategy", "userId", []byte(`{"body": "test"}`), nil)
	assert.NotNil(t, err)
}

func TestCreateCredentials(t *testing.T) {
	c := &internal.MockedConnection{
		MockSend: func(query []byte, options types.QueryOptions) *types.KuzzleResponse {
			parsedQuery := &types.KuzzleRequest{}
			json.Unmarshal(query, parsedQuery)

			assert.Equal(t, "security", parsedQuery.Controller)
			assert.Equal(t, "createCredentials", parsedQuery.Action)
			assert.Equal(t, "strategy", parsedQuery.Strategy)
			assert.Equal(t, "userid", parsedQuery.Id)

			return &types.KuzzleResponse{Result: []byte{}}
		},
	}

	k, _ := kuzzle.NewKuzzle(c, nil)
	_, err := k.Security.CreateCredentials("strategy", "userid", []byte(`{"body": "test"}`), nil)

	assert.Nil(t, err)
}