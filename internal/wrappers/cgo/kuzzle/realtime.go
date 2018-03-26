package main

/*
	#cgo CFLAGS: -I../../headers
	#include "kuzzlesdk.h"
    #include "sdk_wrappers_internal.h"
*/
import "C"

import (
	"encoding/json"
	"unsafe"

	"github.com/kuzzleio/sdk-go/kuzzle"
	"github.com/kuzzleio/sdk-go/realtime"
)

// map which stores instances to keep references in case the gc passes
var realtimeInstances map[interface{}]bool

// register new instance to the instances map
func registerRealtime(instance interface{}) {
	realtimeInstances[instance] = true
}

// unregister an instance from the instances map
//export unregisterRealtime
func unregisterRealtime(rt *C.realtime) {
	delete(realtimeInstances, (*realtime.Realtime)(rt.instance))
}

// Allocates memory
//export kuzzle_new_realtime
func kuzzle_new_realtime(rt *C.realtime, k *C.kuzzle) {
	kuz := (*kuzzle.Kuzzle)(k.instance)
	gort := realtime.NewRealtime(kuz)

	if realtimeInstances == nil {
		realtimeInstances = make(map[interface{}]bool)
	}

	rt.instance = unsafe.Pointer(gort)
	rt.kuzzle = k

	registerRealtime(rt)
}

//export kuzzle_realtime_count
func kuzzle_realtime_count(rt *C.realtime, index, collection, roomId *C.char) *C.int_result {
	res, err := (*realtime.Realtime)(rt.instance).Count(C.GoString(index), C.GoString(collection), C.GoString(roomId))
	return goToCIntResult(res, err)
}

//export kuzzle_realtime_list
func kuzzle_realtime_list(rt *C.realtime, index, collection *C.char) *C.string_result {
	res, err := (*realtime.Realtime)(rt.instance).List(C.GoString(index), C.GoString(collection))
	var stringResult string
	json.Unmarshal(res, &stringResult)
	return goToCStringResult(&stringResult, err)
}

//export kuzzle_realtime_publish
func kuzzle_realtime_publish(rt *C.realtime, index, collection, body *C.char) *C.void_result {
	err := (*realtime.Realtime)(rt.instance).Publish(C.GoString(index), C.GoString(collection), C.GoString(body))
	return goToCVoidResult(err)
}

//export kuzzle_realtime_unsubscribe
func kuzzle_realtime_unsubscribe(rt *C.realtime, roomId *C.char) *C.void_result {
	err := (*realtime.Realtime)(rt.instance).Unsubscribe(C.GoString(roomId))
	return goToCVoidResult(err)
}

//export kuzzle_realtime_subscribe
func kuzzle_realtime_subscribe(rt *C.realtime, index, collection, body *C.char, callback *C.callback, options *C.room_options) *C.string_result {
	res, err := (*realtime.Realtime)(rt.instance).Subscribe(C.GoString(index), C.GoString(collection), json.RawMessage(C.GoString(body)), cToGoKuzzleNotificationChannel(callback), cToGoRoomOptions(options))
	return goToCStringResult(&res, err)
}

//export kuzzle_realtime_join
func kuzzle_realtime_join(rt *C.realtime, index, collection, roomId *C.char, callback *C.callback) *C.void_result {
	err := (*realtime.Realtime)(rt.instance).Join(C.GoString(index), C.GoString(collection), C.GoString(roomId), cToGoKuzzleNotificationChannel(callback))
	return goToCVoidResult(err)
}

//export kuzzle_realtime_validate
func kuzzle_realtime_validate(rt *C.realtime, index, collection, body *C.char) *C.bool_result {
	res, err := (*realtime.Realtime)(rt.instance).Validate(C.GoString(index), C.GoString(collection), C.GoString(body))
	return goToCBoolResult(res, err)
}
