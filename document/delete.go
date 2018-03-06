package document

import (
	"encoding/json"

	"github.com/kuzzleio/sdk-go/types"
)

func (d *Document) Delete(index string, collection string, _id string, options *DocumentOptions) (string, error) {
	if index == "" {
		return "", types.NewError("Document.Delete: index required", 400)
	}

	if collection == "" {
		return "", types.NewError("Document.Delete: collection required", 400)
	}

	if _id == "" {
		return "", types.NewError("Document.Delete: id required", 400)
	}

	ch := make(chan *types.KuzzleResponse)

	query := &types.KuzzleRequest{
		Index:      index,
		Collection: collection,
		Controller: "document",
		Action:     "delete",
		Id:         _id,
	}

	queryOpts := types.NewQueryOptions()

	if options != nil {
		queryOpts.SetVolatile(options.Volatile)
		queryOpts.SetRefresh(options.WaitFor)
	}

	go d.Kuzzle.Query(query, queryOpts, ch)

	res := <-ch

	if res.Error != nil {
		return "", res.Error
	}

	var deleted string
	json.Unmarshal(res.Result, &deleted)

	return deleted, nil
}
