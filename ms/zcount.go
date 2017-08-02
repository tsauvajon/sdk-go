package ms

import (
	"encoding/json"
	"errors"
	"github.com/kuzzleio/sdk-go/types"
	"strconv"
)

/*
  Returns the number of elements held by a sorted set with a score between the provided min and max values.

  By default, the provided min and max values are inclusive. This behavior can be changed using the syntax described in the Redis ZRANGEBYSCORE documentation.
*/
func (ms Ms) Zcount(key string, min int, max int, options types.QueryOptions) (int, error) {
	if key == "" {
		return 0, errors.New("Ms.Zcount: key required")
	}

	result := make(chan types.KuzzleResponse)

	query := types.KuzzleRequest{
		Controller: "ms",
		Action:     "zcount",
		Id:         key,
		Min:        strconv.Itoa(min),
		Max:        strconv.Itoa(max),
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
