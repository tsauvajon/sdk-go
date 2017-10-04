package ms

import (
	"encoding/json"
	"errors"
	"github.com/kuzzleio/sdk-go/types"
)

// Lrange returns the list elements between the start and stop positions (inclusive).
func (ms Ms) Lrange(key string, start int, stop int, options types.QueryOptions) ([]string, error) {
	if key == "" {
		return nil, errors.New("Ms.Lrange: key required")
	}

	result := make(chan *types.KuzzleResponse)

	query := &types.KuzzleRequest{
		Controller: "ms",
		Action:     "lrange",
		Id:         key,
		Start:      start,
		Stop:       stop,
	}

	go ms.Kuzzle.Query(query, options, result)

	res := <-result

	if res.Error != nil {
		return nil, errors.New(res.Error.Message)
	}
	var returnedResult []string
	json.Unmarshal(res.Result, &returnedResult)

	return returnedResult, nil
}
