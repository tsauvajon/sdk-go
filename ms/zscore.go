package ms

import (
	"encoding/json"
	"errors"
	"github.com/kuzzleio/sdk-go/types"
	"strconv"
)

/*
  Returns the score of a member in a sorted set.
*/
func (ms Ms) Zscore(key string, member string, options types.QueryOptions) (float64, error) {
	if key == "" {
		return 0, errors.New("Ms.Zscore: key required")
	}
	if member == "" {
		return 0, errors.New("Ms.Zscore: member required")
	}

	result := make(chan types.KuzzleResponse)

	query := types.KuzzleRequest{
		Controller: "ms",
		Action:     "zscore",
		Id:         key,
		Member:     member,
	}

	go ms.Kuzzle.Query(query, options, result)

	res := <-result

	if res.Error.Message != "" {
		return 0, errors.New(res.Error.Message)
	}

	var scanResponse string
	json.Unmarshal(res.Result, &scanResponse)

	floatResult, _ := strconv.ParseFloat(scanResponse, 64)

	return floatResult, nil
}