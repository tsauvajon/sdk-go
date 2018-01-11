#ifndef _DOCUMENT_HPP_
#define _DOCUMENT_HPP_

#include "exceptions.hpp"
#include "collection.hpp"
#include "core.hpp"

#include <string>
#include <iostream>


namespace kuzzleio {
    class Document {
        Document(){};
        document *_document;
        Collection *_collection;

        public:
            Document(Collection *collection, const std::string& id=NULL, json_object* content=NULL) Kuz_Throw_KuzzleException;
            virtual ~Document();
            std::string delete_(query_options* options=NULL) Kuz_Throw_KuzzleException;
            bool exists(query_options* options=NULL) Kuz_Throw_KuzzleException;
            bool publish(query_options* options=NULL) Kuz_Throw_KuzzleException;
            Document* refresh(query_options* options=NULL) Kuz_Throw_KuzzleException;
            Document* save(query_options* options=NULL) Kuz_Throw_KuzzleException;
            Document* setContent(json_object* content, bool replace=false);
    };
}

#endif