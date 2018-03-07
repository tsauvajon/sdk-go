package main

/*
	#cgo CFLAGS: -I../../headers
	#include "kuzzlesdk.h"
    #include "sdk_wrappers_internal.h"
*/
import "C"
import (
	"unsafe"

	"github.com/kuzzleio/sdk-go/document"
	"github.com/kuzzleio/sdk-go/kuzzle"
)

// map which stores instances to keep references in case the gc passes
var documentInstances map[interface{}]bool

// register new instance to the instances map
func registerDocument(instance interface{}) {
	documentInstances[instance] = true
}

// unregister an instance from the instances map
//export unregisterDocument
func unregisterDocument(d *C.document) {
	delete(documentInstances, (*document.Document)(d.instance))
}

//export kuzzle_new_document
func kuzzle_new_document(d *C.document, k *C.kuzzle) {
	kuz := (*kuzzle.Kuzzle)(k.instance)
	doc := document.NewDocument(kuz)

	if documentInstances == nil {
		documentInstances = make(map[interface{}]bool)
	}

	d.instance = unsafe.Pointer(doc)
	d.kuzzle = k
	registerDocument(doc)
}

//export kuzzle_document_count
func kuzzle_document_count(d *C.document, index *C.char, collection *C.char, body *C.char) *C.int_result {
	res, err := (*document.Document)(d.instance).Count(C.GoString(index), C.GoString(collection), C.GoString(body))
	return goToCIntResult(res, err)
}

//export kuzzle_document_exists
func kuzzle_document_exists(d *C.document, index *C.char, collection *C.char, id *C.char) *C.bool_result {
	res, err := (*document.Document)(d.instance).Exists(C.GoString(index), C.GoString(collection), C.GoString(id))
	return goToCBoolResult(res, err)
}

//export kuzzle_document_create
func kuzzle_document_create(d *C.document, index *C.char, collection *C.char, id *C.char, body *C.char, options *C.document_options) *C.string_result {
	res, err := (*document.Document)(d.instance).Create(C.GoString(index), C.GoString(collection), C.GoString(id), C.GoString(body), cToGoDocumentOptions(options))
	return goToCStringResult(&res, err)
}

//export kuzzle_document_create_or_replace
func kuzzle_document_create_or_replace(d *C.document, index *C.char, collection *C.char, id *C.char, body *C.char, options *C.document_options) *C.string_result {
	res, err := (*document.Document)(d.instance).CreateOrReplace(C.GoString(index), C.GoString(collection), C.GoString(id), C.GoString(body), cToGoDocumentOptions(options))
	return goToCStringResult(&res, err)
}

//export kuzzle_document_delete
func kuzzle_document_delete(d *C.document, index *C.char, collection *C.char, id *C.char, options *C.document_options) *C.string_result {
	res, err := (*document.Document)(d.instance).Delete(C.GoString(index), C.GoString(collection), C.GoString(id), cToGoDocumentOptions(options))
	return goToCStringResult(&res, err)
}

//export kuzzle_document_delete_by_query
func kuzzle_document_delete_by_query(d *C.document, index *C.char, collection *C.char, body *C.char, options *C.document_options) *C.string_array_result {
	res, err := (*document.Document)(d.instance).DeleteByQuery(C.GoString(index), C.GoString(collection), C.GoString(body), cToGoDocumentOptions(options))
	return goToCStringArrayResult(res, err)
}

//export kuzzle_document_get
func kuzzle_document_get(d *C.document, index *C.char, collection *C.char, id *C.char) *C.string_result {
	res, err := (*document.Document)(d.instance).Get(C.GoString(index), C.GoString(collection), C.GoString(id))
	return goToCStringResult(&res, err)
}

//export kuzzle_document_replace
func kuzzle_document_replace(d *C.document, index *C.char, collection *C.char, id *C.char, body *C.char, options *C.document_options) *C.string_result {
	res, err := (*document.Document)(d.instance).Replace(C.GoString(index), C.GoString(collection), C.GoString(id), C.GoString(body), cToGoDocumentOptions(options))
	return goToCStringResult(&res, err)
}

//export kuzzle_document_update
func kuzzle_document_update(d *C.document, index *C.char, collection *C.char, id *C.char, body *C.char, options *C.document_options) *C.string_result {
	res, err := (*document.Document)(d.instance).Update(C.GoString(index), C.GoString(collection), C.GoString(id), C.GoString(body), cToGoDocumentOptions(options))
	return goToCStringResult(&res, err)
}

//export kuzzle_document_validate
func kuzzle_document_validate(d *C.document, index *C.char, collection *C.char, body *C.char) *C.bool_result {
	res, err := (*document.Document)(d.instance).Validate(C.GoString(index), C.GoString(collection), C.GoString(body))
	return goToCBoolResult(res, err)
}

//export kuzzle_document_search
func kuzzle_document_search(d *C.document, index *C.char, collection *C.char, body *C.char, options *C.search_options) *C.search_result {
	res, err := (*document.Document)(d.instance).Search(C.GoString(index), C.GoString(collection), C.GoString(body), cToGoSearchOptions(options))
	return goToCSearchResult(res, err)
}

//export kuzzle_document_mcreate
func kuzzle_document_mcreate(d *C.document, index *C.char, collection *C.char, body *C.char, options *C.document_options) *C.string_result {
	res, err := (*document.Document)(d.instance).MCreate(C.GoString(index), C.GoString(collection), C.GoString(body), cToGoDocumentOptions(options))
	return goToCStringResult(&res, err)
}

//export kuzzle_document_mcreate_or_replace
func kuzzle_document_mcreate_or_replace(d *C.document, index *C.char, collection *C.char, body *C.char, options *C.document_options) *C.string_result {
	res, err := (*document.Document)(d.instance).MCreateOrReplace(C.GoString(index), C.GoString(collection), C.GoString(body), cToGoDocumentOptions(options))
	return goToCStringResult(&res, err)
}

//export kuzzle_document_mdelete
func kuzzle_document_mdelete(d *C.document, index *C.char, collection *C.char, ids **C.char, l C.size_t, options *C.document_options) *C.string_array_result {
	res, err := (*document.Document)(d.instance).MDelete(C.GoString(index), C.GoString(collection), cToGoStrings(ids, l), cToGoDocumentOptions(options))
	return goToCStringArrayResult(res, err)
}

//export kuzzle_document_mget
func kuzzle_document_mget(d *C.document, index *C.char, collection *C.char, ids **C.char, l C.size_t, includeTrash C.bool) *C.string_result {
	res, err := (*document.Document)(d.instance).MGet(C.GoString(index), C.GoString(collection), cToGoStrings(ids, l), bool(includeTrash))
	return goToCStringResult(&res, err)
}

//export kuzzle_document_mreplace
func kuzzle_document_mreplace(d *C.document, index *C.char, collection *C.char, body *C.char, options *C.document_options) *C.string_result {
	res, err := (*document.Document)(d.instance).MReplace(C.GoString(index), C.GoString(collection), C.GoString(body), cToGoDocumentOptions(options))
	return goToCStringResult(&res, err)
}

//export kuzzle_document_mupdate
func kuzzle_document_mupdate(d *C.document, index *C.char, collection *C.char, body *C.char, options *C.document_options) *C.string_result {
	res, err := (*document.Document)(d.instance).MUpdate(C.GoString(index), C.GoString(collection), C.GoString(body), cToGoDocumentOptions(options))
	return goToCStringResult(&res, err)
}
