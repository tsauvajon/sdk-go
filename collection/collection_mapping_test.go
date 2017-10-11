package collection_test

import (
	"encoding/json"
	"fmt"
	"github.com/kuzzleio/sdk-go/collection"
	"github.com/kuzzleio/sdk-go/internal"
	"github.com/kuzzleio/sdk-go/kuzzle"
	"github.com/kuzzleio/sdk-go/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCollectionMappingApplyError(t *testing.T) {
	c := &internal.MockedConnection{
		MockSend: func(query []byte, options types.QueryOptions) *types.KuzzleResponse {
			return &types.KuzzleResponse{Error: &types.MessageError{Message: "Unit test error"}}
		},
	}
	k, _ := kuzzle.NewKuzzle(c, nil)
	cl := collection.NewCollection(k, "collection", "index")
	fields := make(map[string]interface{})
	fields["type"] = interface{}("keyword")

	cm := collection.CollectionMapping{
		Mapping: &types.KuzzleFieldsMapping{
			"foo": {
				Type:       "text",
				Properties: fields,
			},
		},
		Collection: cl,
	}

	_, err := cm.Apply(nil)
	assert.NotNil(t, err)
}

func TestCollectionMappingApply(t *testing.T) {
	callCount := 0

	c := &internal.MockedConnection{
		MockSend: func(query []byte, options types.QueryOptions) *types.KuzzleResponse {
			parsedQuery := &types.KuzzleRequest{}
			json.Unmarshal(query, parsedQuery)
			callCount++

			if callCount == 1 {
				assert.Equal(t, "collection", parsedQuery.Controller)
				assert.Equal(t, "getMapping", parsedQuery.Action)
				assert.Equal(t, "index", parsedQuery.Index)
				assert.Equal(t, "collection", parsedQuery.Collection)

				res := types.KuzzleResponse{Result: []byte(`{"index":{"mappings":{"collection":{"properties":{"foo":{"type":"text","properties":{"type":"keyword","ignore_above":256}}}}}}}`)}
				r, _ := json.Marshal(res.Result)
				return &types.KuzzleResponse{Result: r}
			}
			if callCount == 2 {
				assert.Equal(t, "collection", parsedQuery.Controller)
				assert.Equal(t, "updateMapping", parsedQuery.Action)
				assert.Equal(t, "index", parsedQuery.Index)
				assert.Equal(t, "collection", parsedQuery.Collection)

				res := types.KuzzleResponse{Result: []byte(`{"index":{"mappings":{"collection":{"properties":{"foo":{"type":"text","properties":{"type":"keyword","ignore_above":100}}}}}}}`)}
				r, _ := json.Marshal(res.Result)
				return &types.KuzzleResponse{Result: r}
			}
			if callCount == 3 {
				assert.Equal(t, "collection", parsedQuery.Controller)
				assert.Equal(t, "getMapping", parsedQuery.Action)
				assert.Equal(t, "index", parsedQuery.Index)
				assert.Equal(t, "collection", parsedQuery.Collection)

				res := types.KuzzleResponse{Result: []byte(`{"index":{"mappings":{"collection":{"properties":{"foo":{"type":"text","properties":{"type":"keyword","ignore_above":100}}}}}}}`)}
				r, _ := json.Marshal(res.Result)
				return &types.KuzzleResponse{Result: r}
			}

			return &types.KuzzleResponse{Result: nil}
		},
	}
	k, _ := kuzzle.NewKuzzle(c, nil)
	cl := collection.NewCollection(k, "collection", "index")
	cm, _ := cl.GetMapping(nil)
	fields := make(map[string]interface{})
	fields["type"] = interface{}("keyword")
	fields["ignore_above"] = interface{}(100.0)

	fieldMapping := &types.KuzzleFieldsMapping{
		"foo": {
			Type:       "text",
			Properties: fields,
		},
	}

	res, _ := cm.Set(fieldMapping).Apply(nil)

	assert.Equal(t, cm, res)
}

func ExampleCollectionMapping_Apply() {
	c := &internal.MockedConnection{}
	k, _ := kuzzle.NewKuzzle(c, nil)
	cl := collection.NewCollection(k, "collection", "index")
	cm, _ := cl.GetMapping(nil)
	qo := types.NewQueryOptions()
	fields := make(map[string]interface{})
	fields["type"] = interface{}("keyword")
	fields["ignore_above"] = interface{}(100.0)

	var fieldMapping = &types.KuzzleFieldsMapping{
		"foo": {
			Type:       "text",
			Properties: fields,
		},
	}

	res, err := cm.Set(fieldMapping).Apply(qo)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(res.Collection, res.Mapping)
}

func TestCollectionMappingRefreshError(t *testing.T) {
	c := &internal.MockedConnection{
		MockSend: func(query []byte, options types.QueryOptions) *types.KuzzleResponse {
			return &types.KuzzleResponse{Error: &types.MessageError{Message: "Unit test error"}}
		},
	}
	k, _ := kuzzle.NewKuzzle(c, nil)
	cl := collection.NewCollection(k, "collection", "index")
	fields := make(map[string]interface{})
	fields["type"] = interface{}("keyword")
	fields["ignore_above"] = interface{}(100.0)

	cm := collection.CollectionMapping{
		Mapping: &types.KuzzleFieldsMapping{
			"foo": {
				Type:       "text",
				Properties: fields,
			},
		},
		Collection: cl,
	}

	_, err := cm.Refresh(nil)
	assert.NotNil(t, err)
}

func TestCollectionMappingRefreshUnknownIndex(t *testing.T) {
	c := &internal.MockedConnection{
		MockSend: func(query []byte, options types.QueryOptions) *types.KuzzleResponse {
			parsedQuery := &types.KuzzleRequest{}
			json.Unmarshal(query, parsedQuery)

			assert.Equal(t, "collection", parsedQuery.Controller)
			assert.Equal(t, "getMapping", parsedQuery.Action)
			assert.Equal(t, "wrong-index", parsedQuery.Index)
			assert.Equal(t, "collection", parsedQuery.Collection)

			res := types.KuzzleResponse{Result: []byte(`{"index":{"mappings":{"collection":{"properties":{"foo":{"type":"text","properties":{"type":"keyword","ignore_above":256}}}}}}}`)}
			r, _ := json.Marshal(res.Result)
			return &types.KuzzleResponse{Result: r}
		},
	}
	k, _ := kuzzle.NewKuzzle(c, nil)
	cl := collection.NewCollection(k, "collection", "wrong-index")
	fields := make(map[string]interface{})
	fields["type"] = interface{}("keyword")
	fields["ignore_above"] = interface{}(100.0)

	cm := collection.CollectionMapping{
		Mapping: &types.KuzzleFieldsMapping{
			"foo": {
				Type:       "text",
				Properties: fields,
			},
		},
		Collection: cl,
	}

	_, err := cm.Refresh(nil)

	assert.Equal(t, "No mapping found for index wrong-index", fmt.Sprint(err))
}

func TestCollectionMappingRefreshUnknownCollection(t *testing.T) {
	c := &internal.MockedConnection{
		MockSend: func(query []byte, options types.QueryOptions) *types.KuzzleResponse {
			parsedQuery := &types.KuzzleRequest{}
			json.Unmarshal(query, parsedQuery)

			assert.Equal(t, "collection", parsedQuery.Controller)
			assert.Equal(t, "getMapping", parsedQuery.Action)
			assert.Equal(t, "index", parsedQuery.Index)
			assert.Equal(t, "wrong-collection", parsedQuery.Collection)

			res := types.KuzzleResponse{Result: []byte(`{"index":{"mappings":{"collection":{"properties":{"foo":{"type":"text","properties":{"type":"keyword","ignore_above":256}}}}}}}`)}
			r, _ := json.Marshal(res.Result)
			return &types.KuzzleResponse{Result: r}
		},
	}
	k, _ := kuzzle.NewKuzzle(c, nil)
	cl := collection.NewCollection(k, "wrong-collection", "index")
	fields := make(map[string]interface{})
	fields["type"] = interface{}("keyword")
	fields["ignore_above"] = interface{}(100.0)

	cm := collection.CollectionMapping{
		Mapping: &types.KuzzleFieldsMapping{
			"foo": {
				Type:       "text",
				Properties: fields,
			},
		},
		Collection: cl,
	}

	_, err := cm.Refresh(nil)

	assert.Equal(t, "No mapping found for collection wrong-collection", fmt.Sprint(err))
}

func TestCollectionMappingRefresh(t *testing.T) {
	c := &internal.MockedConnection{
		MockSend: func(query []byte, options types.QueryOptions) *types.KuzzleResponse {
			parsedQuery := &types.KuzzleRequest{}
			json.Unmarshal(query, parsedQuery)

			assert.Equal(t, "collection", parsedQuery.Controller)
			assert.Equal(t, "getMapping", parsedQuery.Action)
			assert.Equal(t, "index", parsedQuery.Index)
			assert.Equal(t, "collection", parsedQuery.Collection)

			res := types.KuzzleResponse{Result: []byte(`{"index":{"mappings":{"collection":{"properties":{"foo":{"type":"text","properties":{"type":"keyword","ignore_above":255}}}}}}}`)}
			r, _ := json.Marshal(res.Result)
			return &types.KuzzleResponse{Result: r}
		},
	}
	k, _ := kuzzle.NewKuzzle(c, nil)
	cl := collection.NewCollection(k, "collection", "index")
	fields := make(map[string]interface{})
	fields["type"] = interface{}("keyword")
	fields["ignore_above"] = interface{}(100)
	fields2 := make(map[string]interface{})
	fields2["type"] = interface{}("keyword")
	fields2["ignore_above"] = interface{}(255.0)

	cm := collection.CollectionMapping{
		Mapping: &types.KuzzleFieldsMapping{
			"foo": {
				Type:       "text",
				Properties: fields,
			},
		},
		Collection: cl,
	}
	updatedCm := collection.CollectionMapping{
		Mapping: &types.KuzzleFieldsMapping{
			"foo": {
				Type:       "text",
				Properties: fields2,
			},
		},
		Collection: cl,
	}

	cm.Refresh(nil)
	assert.Equal(t, *updatedCm.Mapping, *cm.Mapping)
}

func ExampleCollectionMapping_Refresh() {
	c := &internal.MockedConnection{}
	k, _ := kuzzle.NewKuzzle(c, nil)
	cl := collection.NewCollection(k, "collection", "index")
	qo := types.NewQueryOptions()
	fields := make(map[string]interface{})
	fields["type"] = interface{}("keyword")
	fields["ignore_above"] = interface{}(100)

	cm := collection.CollectionMapping{
		Mapping: &types.KuzzleFieldsMapping{
			"foo": {
				Type:       "text",
				Properties: fields,
			},
		},
		Collection: cl,
	}

	res, err := cm.Refresh(qo)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(res.Collection, res.Mapping)
}

func TestCollectionMappingSet(t *testing.T) {
	c := &internal.MockedConnection{
		MockSend: func(query []byte, options types.QueryOptions) *types.KuzzleResponse {
			parsedQuery := &types.KuzzleRequest{}
			json.Unmarshal(query, parsedQuery)

			assert.Equal(t, "collection", parsedQuery.Controller)
			assert.Equal(t, "getMapping", parsedQuery.Action)
			assert.Equal(t, "index", parsedQuery.Index)
			assert.Equal(t, "collection", parsedQuery.Collection)

			res := types.KuzzleResponse{Result: []byte(`{"index":{"mappings":{"collection":{"properties":{"foo":{"type":"text","properties":{"type":"keyword","ignore_above":256}}}}}}}`)}
			r, _ := json.Marshal(res.Result)
			return &types.KuzzleResponse{Result: r}
		},
	}
	k, _ := kuzzle.NewKuzzle(c, nil)
	cl := collection.NewCollection(k, "collection", "index")
	cm, _ := cl.GetMapping(nil)
	fields := make(map[string]interface{})
	fields["type"] = interface{}("keywords")
	fields["ignore_above"] = interface{}(100.0)

	fieldMapping := &types.KuzzleFieldsMapping{
		"foo": {
			Type:       "text",
			Properties: fields,
		},
	}

	cm.Set(fieldMapping)

	assert.Equal(t, "keywords", (*cm.Mapping)["foo"].Properties["type"].(string))
}

func ExampleCollectionMapping_Set() {
	c := &internal.MockedConnection{}
	k, _ := kuzzle.NewKuzzle(c, nil)
	cl := collection.NewCollection(k, "collection", "index")
	cm, _ := cl.GetMapping(nil)
	fields := make(map[string]interface{})
	fields["type"] = interface{}("keyword")
	fields["ignore_above"] = interface{}(100)

	fieldMapping := &types.KuzzleFieldsMapping{
		"foo": {
			Type:       "text",
			Properties: fields,
		},
	}

	res := cm.Set(fieldMapping)

	fmt.Println(res.Collection, res.Mapping)
}

func TestCollectionMappingSetHeaders(t *testing.T) {
	k, _ := kuzzle.NewKuzzle(&internal.MockedConnection{}, nil)
	cl := collection.NewCollection(k, "collection", "index")
	fields := make(map[string]interface{})
	fields["type"] = interface{}("keyword")
	fields["ignore_above"] = interface{}(100.0)

	cm := collection.CollectionMapping{
		Mapping: &types.KuzzleFieldsMapping{
			"foo": {
				Type:       "text",
				Properties: fields,
			},
		},
		Collection: cl,
	}

	headers := make(map[string]interface{}, 0)

	assert.Equal(t, headers, k.GetHeaders())

	headers["foo"] = "bar"
	headers["bar"] = "foo"

	cm.SetHeaders(headers, false)

	newHeaders := make(map[string]interface{}, 0)
	newHeaders["foo"] = "rab"

	cm.SetHeaders(newHeaders, false)

	headers["foo"] = "rab"

	assert.Equal(t, headers, k.GetHeaders())
	assert.NotEqual(t, newHeaders, k.GetHeaders())
}

func TestCollectionMappingSetHeadersReplace(t *testing.T) {
	k, _ := kuzzle.NewKuzzle(&internal.MockedConnection{}, nil)
	cl := collection.NewCollection(k, "collection", "index")
	fields := make(map[string]interface{})
	fields["type"] = interface{}("keyword")
	fields["ignore_above"] = interface{}(100.0)
	fields["ignore_above"] = interface{}(100.0)

	cm := collection.CollectionMapping{
		Mapping: &types.KuzzleFieldsMapping{
			"foo": {
				Type:       "text",
				Properties: fields,
			},
		},
		Collection: cl,
	}

	headers := make(map[string]interface{}, 0)

	assert.Equal(t, headers, k.GetHeaders())

	headers["foo"] = "bar"
	headers["bar"] = "foo"

	cm.SetHeaders(headers, false)

	newHeaders := make(map[string]interface{}, 0)
	newHeaders["foo"] = "rab"

	cm.SetHeaders(newHeaders, true)

	headers["foo"] = "rab"

	assert.Equal(t, newHeaders, k.GetHeaders())
	assert.NotEqual(t, headers, k.GetHeaders())
}

func ExampleCollectionMapping_SetHeaders() {
	c := &internal.MockedConnection{}
	k, _ := kuzzle.NewKuzzle(c, nil)
	cl := collection.NewCollection(k, "collection", "index")
	fields := make(map[string]interface{})
	fields["type"] = interface{}("keyword")
	fields["ignore_above"] = interface{}(100)

	cm := collection.CollectionMapping{
		Mapping: &types.KuzzleFieldsMapping{
			"foo": {
				Type:       "text",
				Properties: fields,
			},
		},
		Collection: cl,
	}

	headers := make(map[string]interface{}, 0)

	headers["foo"] = "bar"
	headers["bar"] = "foo"

	cm.SetHeaders(headers, false)

	fmt.Println(k.GetHeaders())
}
