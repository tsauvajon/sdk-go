package document

import (
	"encoding/json"

	"github.com/kuzzleio/sdk-go/types"
)

func (d *Document) Validate(index string, collection string, body string) (bool, error) {
	if index == "" {
		return false, types.NewError("Document.Validate: index required", 400)
	}

	if collection == "" {
		return false, types.NewError("Document.Validate: collection required", 400)
	}

	if body == "" {
		return false, types.NewError("Document.Validate: body required", 400)
	}

	ch := make(chan *types.KuzzleResponse)

	query := &types.KuzzleRequest{
		Collection: collection,
		Index:      index,
		Controller: "document",
		Action:     "validate",
		Body:       body,
	}

	go d.Kuzzle.Query(query, nil, ch)

	res := <-ch

	if res.Error != nil {
		return false, res.Error
	}

	var valid bool
	json.Unmarshal(res.Result.valid, &valid)

	return valid, nil
}
