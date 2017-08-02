package ms_test

import (
	"encoding/json"
	"github.com/kuzzleio/sdk-go/internal"
	"github.com/kuzzleio/sdk-go/kuzzle"
	"github.com/kuzzleio/sdk-go/types"
	"github.com/stretchr/testify/assert"
	MemoryStorage "github.com/kuzzleio/sdk-go/ms"
	"testing"
	"fmt"
)

func TestHsetnxEmptyKey(t *testing.T) {
	k, _ := kuzzle.NewKuzzle(&internal.MockedConnection{}, nil)
	memoryStorage := MemoryStorage.NewMs(k)
	qo := types.NewQueryOptions()

	_, err := memoryStorage.Hsetnx("", "bar", "barbar", qo)

	assert.NotNil(t, err)
	assert.Equal(t, "Ms.Hsetnx: key required", fmt.Sprint(err))
}

func TestHsetnxEmptyField(t *testing.T) {
	k, _ := kuzzle.NewKuzzle(&internal.MockedConnection{}, nil)
	memoryStorage := MemoryStorage.NewMs(k)
	qo := types.NewQueryOptions()

	_, err := memoryStorage.Hsetnx("foo", "", "barbar", qo)

	assert.NotNil(t, err)
	assert.Equal(t, "Ms.Hsetnx: field required", fmt.Sprint(err))
}

func TestHsetnxError(t *testing.T) {
	c := &internal.MockedConnection{
		MockSend: func(query []byte, options types.QueryOptions) types.KuzzleResponse {
			return types.KuzzleResponse{Error: types.MessageError{Message: "Unit test error"}}
		},
	}
	k, _ := kuzzle.NewKuzzle(c, nil)
	memoryStorage := MemoryStorage.NewMs(k)
	qo := types.NewQueryOptions()

	_, err := memoryStorage.Hsetnx("foo", "bar", "barbar", qo)

	assert.NotNil(t, err)
}

func TestHsetnx(t *testing.T) {
	c := &internal.MockedConnection{
		MockSend: func(query []byte, options types.QueryOptions) types.KuzzleResponse {
			parsedQuery := &types.KuzzleRequest{}
			json.Unmarshal(query, parsedQuery)

			assert.Equal(t, "ms", parsedQuery.Controller)
			assert.Equal(t, "hsetnx", parsedQuery.Action)
			assert.Equal(t, "foo", parsedQuery.Id)
			assert.Equal(t, "bar", parsedQuery.Body.(map[string]interface{})["field"].(string))
			assert.Equal(t, "barbar", parsedQuery.Body.(map[string]interface{})["value"].(string))

			r, _ := json.Marshal(1)
			return types.KuzzleResponse{Result: r}
		},
	}
	k, _ := kuzzle.NewKuzzle(c, nil)
	memoryStorage := MemoryStorage.NewMs(k)
	qo := types.NewQueryOptions()

	res, _ := memoryStorage.Hsetnx("foo", "bar", "barbar", qo)

	assert.Equal(t, 1, res)
}
