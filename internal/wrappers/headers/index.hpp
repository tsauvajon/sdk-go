#ifndef _KUZZLE_INDEX_HPP
#define _KUZZLE_INDEX_HPP

#include "exceptions.hpp"
#include "core.hpp"

namespace kuzzleio {
  class Index {
    kuzzle_index *_index;
    Index();

    public:
      Index(Kuzzle* kuzzle);
      Index(Kuzzle* kuzzle, kuzzle_index *index);
      virtual ~Index();
      void create(const std::string& index, query_options *options=NULL) Kuz_Throw_KuzzleException;
      void delete_(const std::string& index, query_options *options=NULL) Kuz_Throw_KuzzleException;
      std::vector<std::string> mDelete(const std::vector<std::string>& indexes, query_options *options=NULL) Kuz_Throw_KuzzleException;
      bool exists(const std::string& index, query_options *options=NULL) Kuz_Throw_KuzzleException;
      void refresh(const std::string& index, query_options *options=NULL) Kuz_Throw_KuzzleException;
      void refreshInternal(query_options *options=NULL) Kuz_Throw_KuzzleException;
      void setAutoRefresh(const std::string& index, bool autoRefresh, query_options *options=NULL) Kuz_Throw_KuzzleException;
      bool getAutoRefresh(const std::string& index, query_options *options=NULL) Kuz_Throw_KuzzleException;
      std::vector<std::string> list(query_options *options=NULL) Kuz_Throw_KuzzleException;
  };
}

#endif