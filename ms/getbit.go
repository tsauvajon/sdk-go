package ms

import (
	"encoding/json"
	"errors"
	"github.com/kuzzleio/sdk-go/types"
)

// Getbit returns the bit value at offset, in the string value stored in a key.
func (ms Ms) Getbit(key string, offset int, options types.QueryOptions) (int, error) {
	if key == "" {
		return 0, errors.New("Ms.Getbit: key required")
	}

	result := make(chan *types.KuzzleResponse)

	query := &types.KuzzleRequest{
		Controller: "ms",
		Action:     "getbit",
		Id:         key,
		Offset:     offset,
	}

	go ms.Kuzzle.Query(query, options, result)

	res := <-result

	if res.Error != nil {
		return 0, errors.New(res.Error.Message)
	}

	var returnedResult int
	json.Unmarshal(res.Result, &returnedResult)

	return returnedResult, nil
}
