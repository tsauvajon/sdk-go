package collection_test

import (
	"encoding/json"
	"github.com/kuzzleio/sdk-go/collection"
	"github.com/kuzzleio/sdk-go/internal"
	"github.com/kuzzleio/sdk-go/kuzzle"
	"github.com/kuzzleio/sdk-go/types"
	"github.com/stretchr/testify/assert"
	"testing"
	"fmt"
)

type ExistsFilter struct {
	Field string `json:"field"`
}
type QueryFilters struct {
	Exists ExistsFilter `json:"exists"`
}

func TestFetchNextError(t *testing.T) {
	c := &internal.MockedConnection{
		MockSend: func(query []byte, options types.QueryOptions) types.KuzzleResponse {
			return types.KuzzleResponse{Error: types.MessageError{Message: "Unit test error"}}
		},
	}
	k, _ := kuzzle.NewKuzzle(c, nil)
	cl := collection.NewCollection(k, "collection", "index")
	ksr := collection.KuzzleSearchResult{Collection: *cl}

	_, err := ksr.FetchNext()

	assert.NotNil(t, err)
}

func TestFetchNextNotPossible(t *testing.T) {
	k, _ := kuzzle.NewKuzzle(&internal.MockedConnection{}, nil)
	cl := collection.NewCollection(k, "collection", "index")
	ksr := collection.KuzzleSearchResult{Collection: *cl}

	_, err := ksr.FetchNext()

	assert.NotNil(t, err)
	assert.Equal(t, "KuzzleSearchResult.FetchNext: Unable to retrieve next results from search: missing scrollId or from/size parameters", fmt.Sprint(err))
}

func TestFetchNextWithScroll(t *testing.T) {
	requestCount := 0

	var sort = []map[string]string{}
	sort = append(sort, map[string]string{"price": "asc"})

	var filters = types.SearchFilters{
		Query: QueryFilters{Exists: ExistsFilter{Field: "price"}},
		Sort: sort,
	}

	var options = types.NewQueryOptions()
	options.SetSize(2)
	options.SetScroll("1m")

	c := &internal.MockedConnection{
		MockSend: func(query []byte, options types.QueryOptions) types.KuzzleResponse {
			parsedQuery := &types.KuzzleRequest{}
			json.Unmarshal(query, parsedQuery)

			if requestCount == 0 {
				requestCount++
				assert.Equal(t, "document", parsedQuery.Controller)
				assert.Equal(t, "search", parsedQuery.Action)
				assert.Equal(t, "index", parsedQuery.Index)
				assert.Equal(t, "collection", parsedQuery.Collection)

				results := []collection.Document{
					{Id: "product1", Content: []byte(`{"label":"Foo1","price":1200}`)},
					{Id: "product2", Content: []byte(`{"label":"Foo2","price":800}`)},
				}

				k, _ := kuzzle.NewKuzzle(&internal.MockedConnection{}, nil)
				cl := collection.NewCollection(k, "collection", "index")

				res := collection.KuzzleSearchResult{
					Total:      4,
					Hits:       results,
					ScrollId:   "f00b4r",
					Filters:    filters,
					Options:    options,
					Collection: *cl,
				}
				r, _ := json.Marshal(res)
				return types.KuzzleResponse{Result: r}
			}
			if requestCount == 1 {
				requestCount++
				assert.Equal(t, "document", parsedQuery.Controller)
				assert.Equal(t, "scroll", parsedQuery.Action)
				assert.Equal(t, "1m", parsedQuery.Scroll)
				assert.Equal(t, "f00b4r", parsedQuery.ScrollId)

				results := []collection.Document{
					{Id: "product3", Content: []byte(`{"label":"Foo3","price":400}`)},
					{Id: "product4", Content: []byte(`{"label":"Foo4","price":200}`)},
				}

				res := collection.KuzzleSearchResult{
					Total: 4,
					Hits:  results,
				}
				r, _ := json.Marshal(res)
				return types.KuzzleResponse{Result: r}
			}

			return types.KuzzleResponse{}
		},
	}
	k, _ := kuzzle.NewKuzzle(c, nil)
	cl := collection.NewCollection(k, "collection", "index")

	ksr, _ := cl.Search(filters, options)
	fetchNextRes, _ := ksr.FetchNext()

	assert.Equal(t, "f00b4r", ksr.ScrollId)
	assert.Equal(t, 4, fetchNextRes.Total)
	assert.Equal(t, 2, len(fetchNextRes.Hits))
	assert.Equal(t, "Foo4", fetchNextRes.Hits[1].SourceToMap()["label"])
}

func TestFetchNextWithSearchAfter(t *testing.T) {
	requestCount := 0

	var sort = []map[string]string{}
	sort = append(sort, map[string]string{"price": "desc"})
	sort = append(sort, map[string]string{"label": "asc"})

	var filters = types.SearchFilters{
		Query: QueryFilters{Exists: ExistsFilter{Field: "price"}},
		Sort: sort,
	}

	var options = types.NewQueryOptions()
	options.SetSize(2)

	c := &internal.MockedConnection{
		MockSend: func(query []byte, options types.QueryOptions) types.KuzzleResponse {
			parsedQuery := &types.KuzzleRequest{}
			json.Unmarshal(query, parsedQuery)

			if requestCount == 0 {
				requestCount++
				assert.Equal(t, "document", parsedQuery.Controller)
				assert.Equal(t, "search", parsedQuery.Action)
				assert.Equal(t, "index", parsedQuery.Index)
				assert.Equal(t, "collection", parsedQuery.Collection)

				results := []collection.Document{
					{Id: "product1", Content: []byte(`{"label":"Foo1","price":"1200"}`)},
					{Id: "product2", Content: []byte(`{"label":"Foo2","price":"800"}`)},
				}

				k, _ := kuzzle.NewKuzzle(&internal.MockedConnection{}, nil)
				cl := collection.NewCollection(k, "collection", "index")

				res := collection.KuzzleSearchResult{
					Total:      4,
					Hits:       results,
					Filters:    filters,
					Options:    options,
					Collection: *cl,
				}
				r, _ := json.Marshal(res)
				return types.KuzzleResponse{Result: r}
			}
			if requestCount == 1 {
				requestCount++
				assert.Equal(t, "document", parsedQuery.Controller)
				assert.Equal(t, "search", parsedQuery.Action)
				assert.Equal(t, "index", parsedQuery.Index)
				assert.Equal(t, "collection", parsedQuery.Collection)
				assert.Equal(t, []interface {}([]interface {}{"800", "Foo2"}), parsedQuery.Body.(map[string]interface{})["search_after"])

				results := []collection.Document{
					{Id: "product3", Content: []byte(`{"label":"Foo3","price":"400"}`)},
					{Id: "product4", Content: []byte(`{"label":"Foo4","price":"200"}`)},
				}

				res := collection.KuzzleSearchResult{
					Total: 4,
					Hits:  results,
				}
				r, _ := json.Marshal(res)
				return types.KuzzleResponse{Result: r}
			}

			return types.KuzzleResponse{}
		},
	}
	k, _ := kuzzle.NewKuzzle(c, nil)
	cl := collection.NewCollection(k, "collection", "index")

	ksr, _ := cl.Search(filters, options)
	fetchNextRes, _ := ksr.FetchNext()

	assert.Equal(t, 4, fetchNextRes.Total)
	assert.Equal(t, 2, len(fetchNextRes.Hits))
	assert.Equal(t, "Foo4", fetchNextRes.Hits[1].SourceToMap()["label"])
}

func TestFetchNextWithSizeFrom(t *testing.T) {
	requestCount := 0

	var filters = types.SearchFilters{
		Query: QueryFilters{Exists: ExistsFilter{Field: "price"}},
	}

	var options = types.NewQueryOptions()
	options.SetSize(2)

	c := &internal.MockedConnection{
		MockSend: func(query []byte, options types.QueryOptions) types.KuzzleResponse {
			parsedQuery := &types.KuzzleRequest{}
			json.Unmarshal(query, parsedQuery)

			if requestCount == 0 {
				requestCount++
				assert.Equal(t, "document", parsedQuery.Controller)
				assert.Equal(t, "search", parsedQuery.Action)
				assert.Equal(t, "index", parsedQuery.Index)
				assert.Equal(t, "collection", parsedQuery.Collection)

				results := []collection.Document{
					{Id: "product1", Content: []byte(`{"label":"Foo1","price":1200}`)},
					{Id: "product2", Content: []byte(`{"label":"Foo2","price":800}`)},
				}

				k, _ := kuzzle.NewKuzzle(&internal.MockedConnection{}, nil)
				cl := collection.NewCollection(k, "collection", "index")

				res := collection.KuzzleSearchResult{
					Total:      4,
					Hits:       results,
					Filters:    filters,
					Options:    options,
					Collection: *cl,
				}
				r, _ := json.Marshal(res)
				return types.KuzzleResponse{Result: r}
			}
			if requestCount == 1 {
				requestCount++
				assert.Equal(t, "document", parsedQuery.Controller)
				assert.Equal(t, "search", parsedQuery.Action)
				assert.Equal(t, "index", parsedQuery.Index)
				assert.Equal(t, "collection", parsedQuery.Collection)

				results := []collection.Document{
					{Id: "product3", Content: []byte(`{"label":"Foo3","price":400}`)},
					{Id: "product4", Content: []byte(`{"label":"Foo4","price":200}`)},
				}

				res := collection.KuzzleSearchResult{
					Total:   4,
					Hits:    results,
					Options: options,
				}
				r, _ := json.Marshal(res)
				return types.KuzzleResponse{Result: r}
			}

			return types.KuzzleResponse{}
		},
	}
	k, _ := kuzzle.NewKuzzle(c, nil)
	cl := collection.NewCollection(k, "collection", "index")

	ksr, _ := cl.Search(filters, options)
	fetchNextRes, _ := ksr.FetchNext()

	assert.Equal(t, 4, fetchNextRes.Total)
	assert.Equal(t, 2, len(fetchNextRes.Hits))
	assert.Equal(t, "Foo4", fetchNextRes.Hits[1].SourceToMap()["label"])

	tooFarRes, _ := fetchNextRes.FetchNext()

	assert.Equal(t, collection.KuzzleSearchResult{}, tooFarRes)
}