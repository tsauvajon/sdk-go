// Copyright 2015-2018 Kuzzle
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ms_test

import (
	"encoding/json"
	"fmt"
	"github.com/kuzzleio/sdk-go/internal"
	"github.com/kuzzleio/sdk-go/kuzzle"
	"github.com/kuzzleio/sdk-go/protocol/websocket"
	"github.com/kuzzleio/sdk-go/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestZrevrangebylexEmptyMin(t *testing.T) {
	k, _ := kuzzle.NewKuzzle(&internal.MockedConnection{}, nil)

	_, err := k.MemoryStorage.Zrevrangebylex("foo", "", "(g", nil)

	assert.NotNil(t, err)
	assert.Equal(t, "[400] Ms.Zrevrangebylex: an empty string is not a valid string range item", fmt.Sprint(err))
}

func TestZrevrangebylexEmptyMax(t *testing.T) {
	k, _ := kuzzle.NewKuzzle(&internal.MockedConnection{}, nil)

	_, err := k.MemoryStorage.Zrevrangebylex("foo", "-", "", nil)

	assert.NotNil(t, err)
	assert.Equal(t, "[400] Ms.Zrevrangebylex: an empty string is not a valid string range item", fmt.Sprint(err))
}

func TestZrevrangebylexError(t *testing.T) {
	c := &internal.MockedConnection{
		MockSend: func(query []byte, options types.QueryOptions) *types.KuzzleResponse {
			return &types.KuzzleResponse{Error: types.KuzzleError{Message: "Unit test error"}}
		},
	}
	k, _ := kuzzle.NewKuzzle(c, nil)

	_, err := k.MemoryStorage.Zrevrangebylex("foo", "-", "(g", nil)

	assert.NotNil(t, err)
}

func TestZrevrangebylex(t *testing.T) {
	c := &internal.MockedConnection{
		MockSend: func(query []byte, options types.QueryOptions) *types.KuzzleResponse {
			parsedQuery := &types.KuzzleRequest{}
			json.Unmarshal(query, parsedQuery)

			assert.Equal(t, "ms", parsedQuery.Controller)
			assert.Equal(t, "zrevrangebylex", parsedQuery.Action)
			assert.Equal(t, []interface{}([]interface{}{"withscores"}), parsedQuery.Options)

			r, _ := json.Marshal([]string{"bar", "rab"})
			return &types.KuzzleResponse{Result: r}
		},
	}
	k, _ := kuzzle.NewKuzzle(c, nil)

	res, _ := k.MemoryStorage.Zrevrangebylex("foo", "-", "(g", nil)

	assert.Equal(t, []string{"bar", "rab"}, res)
}

func TestZrevrangebylexWithLimits(t *testing.T) {
	c := &internal.MockedConnection{
		MockSend: func(query []byte, options types.QueryOptions) *types.KuzzleResponse {
			parsedQuery := &types.KuzzleRequest{}
			json.Unmarshal(query, parsedQuery)

			assert.Equal(t, "ms", parsedQuery.Controller)
			assert.Equal(t, "zrevrangebylex", parsedQuery.Action)
			assert.Equal(t, []int{0, 1}, options.Limit())
			assert.Equal(t, []interface{}([]interface{}{"withscores"}), parsedQuery.Options)

			r, _ := json.Marshal([]string{"bar", "rab"})
			return &types.KuzzleResponse{Result: r}
		},
	}
	k, _ := kuzzle.NewKuzzle(c, nil)
	qo := types.NewQueryOptions()

	qo.SetLimit([]int{0, 1})
	res, _ := k.MemoryStorage.Zrevrangebylex("foo", "-", "(g", qo)

	assert.Equal(t, []string{"bar", "rab"}, res)
}

func ExampleMs_Zrevrangebylex() {
	c := websocket.NewWebSocket("localhost:7512", nil)
	k, _ := kuzzle.NewKuzzle(c, nil)
	qo := types.NewQueryOptions()

	qo.SetLimit([]int{0, 1})
	res, err := k.MemoryStorage.Zrevrangebylex("foo", "-", "(g", qo)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(res)
}
