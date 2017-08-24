package collection_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/kuzzleio/sdk-go/kuzzle"
	"github.com/kuzzleio/sdk-go/internal"
	"github.com/kuzzleio/sdk-go/collection"
	"github.com/kuzzleio/sdk-go/types"
	"encoding/json"
	"fmt"
)

func TestCollectionMappingApplyError(t *testing.T) {
	c := &internal.MockedConnection{
		MockSend: func(query []byte, options types.QueryOptions) types.KuzzleResponse {
			return types.KuzzleResponse{Error: types.MessageError{Message: "Unit test error"}}
		},
	}
	k, _ := kuzzle.NewKuzzle(c, nil)
	cl := collection.NewCollection(k, "collection", "index")

	cm := collection.CollectionMapping{
		Mapping: types.KuzzleFieldMapping{
			"foo": {
				Type: "text",
				Fields: []byte(`{"type":"keyword","ignore_above":256}`),
			},
		},
		Collection: *cl,
	}

	_, err := cm.Apply(nil)
	assert.NotNil(t, err)
}

func TestCollectionMappingApply(t *testing.T) {
	callCount := 0

	c := &internal.MockedConnection{
		MockSend: func(query []byte, options types.QueryOptions) types.KuzzleResponse {
			parsedQuery := &types.KuzzleRequest{}
			json.Unmarshal(query, parsedQuery)

			if callCount == 0 {
				callCount++
				assert.Equal(t, "collection", parsedQuery.Controller)
				assert.Equal(t, "getMapping", parsedQuery.Action)
				assert.Equal(t, "index", parsedQuery.Index)
				assert.Equal(t, "collection", parsedQuery.Collection)

				res := types.KuzzleResponse{Result: []byte(`{"index":{"mappings":{"collection":{"properties":{"foo":{"type":"text","fields":{"type":"keyword","ignore_above":256}}}}}}}`)}
				r, _ := json.Marshal(res.Result)
				return types.KuzzleResponse{Result: r}
			}
			if callCount == 1 {
				callCount++
				assert.Equal(t, "collection", parsedQuery.Controller)
				assert.Equal(t, "updateMapping", parsedQuery.Action)
				assert.Equal(t, "index", parsedQuery.Index)
				assert.Equal(t, "collection", parsedQuery.Collection)

				res := types.KuzzleResponse{Result: []byte(`{"index":{"mappings":{"collection":{"properties":{"foo":{"type":"text","fields":{"type":"keyword","ignore_above":100}}}}}}}`)}
				r, _ := json.Marshal(res.Result)
				return types.KuzzleResponse{Result: r}
			}
			if callCount == 2 {
				callCount++
				assert.Equal(t, "collection", parsedQuery.Controller)
				assert.Equal(t, "getMapping", parsedQuery.Action)
				assert.Equal(t, "index", parsedQuery.Index)
				assert.Equal(t, "collection", parsedQuery.Collection)

				res := types.KuzzleResponse{Result: []byte(`{"index":{"mappings":{"collection":{"properties":{"foo":{"type":"text","fields":{"type":"keyword","ignore_above":100}}}}}}}`)}
				r, _ := json.Marshal(res.Result)
				return types.KuzzleResponse{Result: r}
			}

			return types.KuzzleResponse{Result: nil}
		},
	}
	k, _ := kuzzle.NewKuzzle(c, nil)
	cl := collection.NewCollection(k, "collection", "index")
	cm, _ := cl.GetMapping(nil)

	var fieldMapping = types.KuzzleFieldMapping{
		"foo": {
			Type: "text",
			Fields: []byte(`{"type":"keyword","ignore_above":100}`),
		},
	}

	res, _ := cm.Set(fieldMapping).Apply(nil)

	assert.Equal(t, cm, res)
}

func TestCollectionMappingRefreshError(t *testing.T) {
	c := &internal.MockedConnection{
		MockSend: func(query []byte, options types.QueryOptions) types.KuzzleResponse {
			return types.KuzzleResponse{Error: types.MessageError{Message: "Unit test error"}}
		},
	}
	k, _ := kuzzle.NewKuzzle(c, nil)
	cl := collection.NewCollection(k, "collection", "index")

	cm := collection.CollectionMapping{
		Mapping: types.KuzzleFieldMapping{
			"foo": {
				Type: "text",
				Fields: []byte(`{"type":"keyword","ignore_above":256}`),
			},
		},
		Collection: *cl,
	}

	_, err := cm.Refresh(nil)
	assert.NotNil(t, err)
}

func TestCollectionMappingRefreshUnknownIndex(t *testing.T) {
	c := &internal.MockedConnection{
		MockSend: func(query []byte, options types.QueryOptions) types.KuzzleResponse {
			parsedQuery := &types.KuzzleRequest{}
			json.Unmarshal(query, parsedQuery)

			assert.Equal(t, "collection", parsedQuery.Controller)
			assert.Equal(t, "getMapping", parsedQuery.Action)
			assert.Equal(t, "wrong-index", parsedQuery.Index)
			assert.Equal(t, "collection", parsedQuery.Collection)

			res := types.KuzzleResponse{Result: []byte(`{"index":{"mappings":{"collection":{"properties":{"foo":{"type":"text","fields":{"type":"keyword","ignore_above":256}}}}}}}`)}
			r, _ := json.Marshal(res.Result)
			return types.KuzzleResponse{Result: r}
		},
	}
	k, _ := kuzzle.NewKuzzle(c, nil)
	cl := collection.NewCollection(k, "collection", "wrong-index")

	cm := collection.CollectionMapping{
		Mapping: types.KuzzleFieldMapping{
			"foo": {
				Type: "text",
				Fields: []byte(`{"type":"keyword","ignore_above":256}`),
			},
		},
		Collection: *cl,
	}

	_, err := cm.Refresh(nil)

	assert.Equal(t, "No mapping found for index wrong-index", fmt.Sprint(err))
}

func TestCollectionMappingRefreshUnknownCollection(t *testing.T) {
	c := &internal.MockedConnection{
		MockSend: func(query []byte, options types.QueryOptions) types.KuzzleResponse {
			parsedQuery := &types.KuzzleRequest{}
			json.Unmarshal(query, parsedQuery)

			assert.Equal(t, "collection", parsedQuery.Controller)
			assert.Equal(t, "getMapping", parsedQuery.Action)
			assert.Equal(t, "index", parsedQuery.Index)
			assert.Equal(t, "wrong-collection", parsedQuery.Collection)

			res := types.KuzzleResponse{Result: []byte(`{"index":{"mappings":{"collection":{"properties":{"foo":{"type":"text","fields":{"type":"keyword","ignore_above":256}}}}}}}`)}
			r, _ := json.Marshal(res.Result)
			return types.KuzzleResponse{Result: r}
		},
	}
	k, _ := kuzzle.NewKuzzle(c, nil)
	cl := collection.NewCollection(k, "wrong-collection", "index")

	cm := collection.CollectionMapping{
		Mapping: types.KuzzleFieldMapping{
			"foo": {
				Type: "text",
				Fields: []byte(`{"type":"keyword","ignore_above":256}`),
			},
		},
		Collection: *cl,
	}

	_, err := cm.Refresh(nil)

	assert.Equal(t, "No mapping found for collection wrong-collection", fmt.Sprint(err))
}

func TestCollectionMappingRefresh(t *testing.T) {
	c := &internal.MockedConnection{
		MockSend: func(query []byte, options types.QueryOptions) types.KuzzleResponse {
			parsedQuery := &types.KuzzleRequest{}
			json.Unmarshal(query, parsedQuery)

			assert.Equal(t, "collection", parsedQuery.Controller)
			assert.Equal(t, "getMapping", parsedQuery.Action)
			assert.Equal(t, "index", parsedQuery.Index)
			assert.Equal(t, "collection", parsedQuery.Collection)

			res := types.KuzzleResponse{Result: []byte(`{"index":{"mappings":{"collection":{"properties":{"foo":{"type":"text","fields":{"type":"keyword","ignore_above":256}}}}}}}`)}
			r, _ := json.Marshal(res.Result)
			return types.KuzzleResponse{Result: r}
		},
	}
	k, _ := kuzzle.NewKuzzle(c, nil)
	cl := collection.NewCollection(k, "collection", "index")

	cm := collection.CollectionMapping{
		Mapping: types.KuzzleFieldMapping{
			"foo": {
				Type: "text",
				Fields: []byte(`{"type":"keyword","ignore_above":100}`),
			},
		},
		Collection: *cl,
	}
	updatedCm := collection.CollectionMapping{
		Mapping: types.KuzzleFieldMapping{
			"foo": {
				Type: "text",
				Fields: []byte(`{"type":"keyword","ignore_above":256}`),
			},
		},
		Collection: *cl,
	}

	res, _ := cm.Refresh(nil)

	assert.NotEqual(t, cm, res)
	assert.Equal(t, updatedCm, res)
}

func TestCollectionMappingSet(t *testing.T) {
	c := &internal.MockedConnection{
		MockSend: func(query []byte, options types.QueryOptions) types.KuzzleResponse {
			parsedQuery := &types.KuzzleRequest{}
			json.Unmarshal(query, parsedQuery)

			assert.Equal(t, "collection", parsedQuery.Controller)
			assert.Equal(t, "getMapping", parsedQuery.Action)
			assert.Equal(t, "index", parsedQuery.Index)
			assert.Equal(t, "collection", parsedQuery.Collection)

			res := types.KuzzleResponse{Result: []byte(`{"index":{"mappings":{"collection":{"properties":{"foo":{"type":"text","fields":{"type":"keyword","ignore_above":256}}}}}}}`)}
			r, _ := json.Marshal(res.Result)
			return types.KuzzleResponse{Result: r}
		},
	}
	k, _ := kuzzle.NewKuzzle(c, nil)
	cl := collection.NewCollection(k, "collection", "index")
	cm, _ := cl.GetMapping(nil)

	var fieldMapping = types.KuzzleFieldMapping{
		"foo": {
			Type: "text",
			Fields: []byte(`{"type":"keyword","ignore_above":100}`),
		},
	}

	cm.Set(fieldMapping)

	type FieldsStruct struct {
		Type        string `json:"type"`
		IgnoreAbove int    `json:"ignore_above"`
	}
	fieldStruct := FieldsStruct{}
	json.Unmarshal(cm.Mapping["foo"].Fields, &fieldStruct)

	assert.Equal(t, 100, fieldStruct.IgnoreAbove)
}

func TestCollectionMappingSetHeaders(t *testing.T) {
	k, _ := kuzzle.NewKuzzle(&internal.MockedConnection{}, nil)
	cl := collection.NewCollection(k, "collection", "index")

	cm := collection.CollectionMapping{
		Mapping: types.KuzzleFieldMapping{
			"foo": {
				Type: "text",
				Fields: []byte(`{"type":"keyword","ignore_above":100}`),
			},
		},
		Collection: *cl,
	}

	var headers = make(map[string]interface{}, 0)

	assert.Equal(t, headers, k.GetHeaders())

	headers["foo"] = "bar"
	headers["bar"] = "foo"

	cm.SetHeaders(headers, false)

	var newHeaders = make(map[string]interface{}, 0)
	newHeaders["foo"] = "rab"

	cm.SetHeaders(newHeaders, false)

	headers["foo"] = "rab"

	assert.Equal(t, headers, k.GetHeaders())
	assert.NotEqual(t, newHeaders, k.GetHeaders())
}

func TestCollectionMappingSetHeadersReplace(t *testing.T) {
	k, _ := kuzzle.NewKuzzle(&internal.MockedConnection{}, nil)
	cl := collection.NewCollection(k, "collection", "index")

	cm := collection.CollectionMapping{
		Mapping: types.KuzzleFieldMapping{
			"foo": {
				Type: "text",
				Fields: []byte(`{"type":"keyword","ignore_above":100}`),
			},
		},
		Collection: *cl,
	}

	var headers = make(map[string]interface{}, 0)

	assert.Equal(t, headers, k.GetHeaders())

	headers["foo"] = "bar"
	headers["bar"] = "foo"

	cm.SetHeaders(headers, false)

	var newHeaders = make(map[string]interface{}, 0)
	newHeaders["foo"] = "rab"

	cm.SetHeaders(newHeaders, true)

	headers["foo"] = "rab"

	assert.Equal(t, newHeaders, k.GetHeaders())
	assert.NotEqual(t, headers, k.GetHeaders())
}