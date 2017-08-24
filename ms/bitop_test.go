package ms_test

import (
	"encoding/json"
	"fmt"
	"github.com/kuzzleio/sdk-go/internal"
	"github.com/kuzzleio/sdk-go/kuzzle"
	MemoryStorage "github.com/kuzzleio/sdk-go/ms"
	"github.com/kuzzleio/sdk-go/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBitopEmptyKey(t *testing.T) {
	k, _ := kuzzle.NewKuzzle(&internal.MockedConnection{}, nil)
	memoryStorage := MemoryStorage.NewMs(k)
	qo := types.NewQueryOptions()

	_, err := memoryStorage.Bitop("", "", []string{}, qo)

	assert.NotNil(t, err)
	assert.Equal(t, "Ms.Bitop: key required", fmt.Sprint(err))
}

func TestBitopError(t *testing.T) {
	c := &internal.MockedConnection{
		MockSend: func(query []byte, options types.QueryOptions) types.KuzzleResponse {
			return types.KuzzleResponse{Error: types.MessageError{Message: "Unit test error"}}
		},
	}
	k, _ := kuzzle.NewKuzzle(c, nil)
	memoryStorage := MemoryStorage.NewMs(k)
	qo := types.NewQueryOptions()

	_, err := memoryStorage.Bitop("foo", "", []string{}, qo)

	assert.NotNil(t, err)
}

func TestBitop(t *testing.T) {
	c := &internal.MockedConnection{
		MockSend: func(query []byte, options types.QueryOptions) types.KuzzleResponse {
			parsedQuery := &types.KuzzleRequest{}
			json.Unmarshal(query, parsedQuery)

			assert.Equal(t, "ms", parsedQuery.Controller)
			assert.Equal(t, "bitop", parsedQuery.Action)

			assert.Equal(t, "operation", parsedQuery.Body.(map[string]interface{})["operation"].(string))
			assert.Equal(t, "some", parsedQuery.Body.(map[string]interface{})["keys"].([]interface{})[0].(string))
			assert.Equal(t, "keys", parsedQuery.Body.(map[string]interface{})["keys"].([]interface{})[1].(string))

			r, _ := json.Marshal(1)
			return types.KuzzleResponse{Result: r}
		},
	}
	k, _ := kuzzle.NewKuzzle(c, nil)
	memoryStorage := MemoryStorage.NewMs(k)
	qo := types.NewQueryOptions()

	res, _ := memoryStorage.Bitop("foo", "operation", []string{"some", "keys"}, qo)

	assert.Equal(t, 1, res)
}