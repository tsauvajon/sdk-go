package ms

import (
	"encoding/json"
	"errors"
	"github.com/kuzzleio/sdk-go/types"
)

/*
  Sets a field and its value in a hash, only if the field does not already exist.
*/
func (ms Ms) Hsetnx(key string, field string, value string, options types.QueryOptions) (int, error) {
	if key == "" {
		return 0, errors.New("Ms.Hsetnx: key required")
	}
	if field == "" {
		return 0, errors.New("Ms.Hsetnx: field required")
	}

	result := make(chan types.KuzzleResponse)

	type body struct {
		Field string `json:"field"`
		Value string `json:"value"`
	}

	query := types.KuzzleRequest{
		Controller: "ms",
		Action:     "hsetnx",
		Id:         key,
		Body: 			&body{Field: field, Value: value},
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
