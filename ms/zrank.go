package ms

import (
	"encoding/json"
	"errors"
	"github.com/kuzzleio/sdk-go/types"
)

/*
  Returns the position of an element in a sorted set, with scores in ascending order. The index returned is 0-based (the lowest score member has an index of 0).
*/
func (ms Ms) Zrank(key string, member string, options types.QueryOptions) (int, error) {
	if key == "" {
		return 0, errors.New("Ms.Zrank: key required")
	}
	if member == "" {
		return 0, errors.New("Ms.Zrank: member required")
	}

	result := make(chan types.KuzzleResponse)

	query := types.KuzzleRequest{
		Controller: "ms",
		Action:     "zrank",
		Id:         key,
		Member:     member,
	}

	go ms.Kuzzle.Query(query, options, result)

	res := <-result

	if res.Error.Message != "" {
		return 0, errors.New(res.Error.Message)
	}

	var returnedResult int
	json.Unmarshal(res.Result, &returnedResult)

	return returnedResult, nil
}