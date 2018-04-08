package realtime_test

import (
	"encoding/json"
	"testing"

	"github.com/kuzzleio/sdk-go/realtime"

	"github.com/kuzzleio/sdk-go/internal"
	"github.com/kuzzleio/sdk-go/kuzzle"
	"github.com/kuzzleio/sdk-go/types"
	"github.com/stretchr/testify/assert"
)

func TestCountIndexNull(t *testing.T) {
	k, _ := kuzzle.NewKuzzle(&internal.MockedConnection{}, nil)
	nr := realtime.NewRealtime(k)

	_, err := nr.Count("", "collection", "roomID")

	assert.NotNil(t, err)
}

func TestCountCollectionNull(t *testing.T) {
	k, _ := kuzzle.NewKuzzle(&internal.MockedConnection{}, nil)
	nr := realtime.NewRealtime(k)

	_, err := nr.Count("index", "", "roomID")

	assert.NotNil(t, err)
}

func TestCountRoomIDNull(t *testing.T) {
	k, _ := kuzzle.NewKuzzle(&internal.MockedConnection{}, nil)
	nr := realtime.NewRealtime(k)

	_, err := nr.Count("index", "collection", "")

	assert.NotNil(t, err)
}

func TestCountError(t *testing.T) {
	c := &internal.MockedConnection{
		MockSend: func(query []byte, options types.QueryOptions) *types.KuzzleResponse {
			return &types.KuzzleResponse{Error: &types.KuzzleError{Message: "Unit test error"}}
		},
	}
	k, _ := kuzzle.NewKuzzle(c, nil)
	realTime := realtime.NewRealtime(k)

	_, err := realTime.Count("index", "collection", "42")
	assert.NotNil(t, err)
}

func TestCount(t *testing.T) {
	type result struct {
		Count int `json:"count"`
	}

	c := &internal.MockedConnection{
		MockSend: func(query []byte, options types.QueryOptions) *types.KuzzleResponse {
			res := result{Count: 10}
			r, _ := json.Marshal(res)
			return &types.KuzzleResponse{Result: r}
		},
	}
	k, _ := kuzzle.NewKuzzle(c, nil)
	realTime := realtime.NewRealtime(k)

	res, _ := realTime.Count("index", "collection", "42")
	assert.Equal(t, 10, res)
}