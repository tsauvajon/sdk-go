#include "realtime.hpp"

namespace kuzzleio {
  Realtime::Realtime(Kuzzle *kuzzle) {
    _realtime = new realtime();
    kuzzle_new_realtime(_realtime, kuzzle->_kuzzle);
  }

  Realtime::Realtime(Kuzzle *kuzzle, realtime *realtime) {
    _realtime = realtime;
    kuzzle_new_realtime(_realtime, kuzzle->_kuzzle);
  }

  Realtime::~Realtime() {
    unregisterRealtime(_realtime);
    delete(_realtime);
  }

  int Realtime::count(const std::string& index, const std::string collection, const std::string roomId) Kuz_Throw_KuzzleException {
    int_result *r = kuzzle_realtime_count(_realtime, const_cast<char*>(index.c_str()), const_cast<char*>(collection.c_str()), const_cast<char*>(roomId.c_str()));
    if (r->error != NULL)
        throwExceptionFromStatus(r);
    int ret = r->result;
    kuzzle_free_int_result(r);
    return ret;
  }

  void Realtime::join(const std::string& index, const std::string collection, const std::string roomId, callback cb) Kuz_Throw_KuzzleException {
    void_result *r = kuzzle_realtime_join(_realtime, const_cast<char*>(index.c_str()), const_cast<char*>(collection.c_str()), const_cast<char*>(roomId.c_str()), cb);
    if (r->error != NULL)
        throwExceptionFromStatus(r);
    kuzzle_free_void_result(r);
  }

  std::string Realtime::list(const std::string& index, const std::string collection) Kuz_Throw_KuzzleException {
    string_result *r = kuzzle_realtime_list(_realtime, const_cast<char*>(index.c_str()), const_cast<char*>(collection.c_str()));
    if (r->error != NULL)
        throwExceptionFromStatus(r);
    std::string ret = r->result;
    kuzzle_free_string_result(r);
    return ret;
  }

  void Realtime::publish(const std::string& index, const std::string collection, const std::string body) Kuz_Throw_KuzzleException {
    void_result *r = kuzzle_realtime_publish(_realtime, const_cast<char*>(index.c_str()), const_cast<char*>(collection.c_str()), const_cast<char*>(body.c_str()));
    if (r->error != NULL)
        throwExceptionFromStatus(r);
    kuzzle_free_void_result(r);
  }

  std::string Realtime::subscribe(const std::string& index, const std::string collection, const std::string body, callback cb, room_options* options) Kuz_Throw_KuzzleException {
    string_result *r = kuzzle_realtime_subscribe(_realtime, const_cast<char*>(index.c_str()), const_cast<char*>(collection.c_str()),  const_cast<char*>(body.c_str()), cb, options);
    if (r->error != NULL)
        throwExceptionFromStatus(r);
    std::string ret = r->result;
    kuzzle_free_string_result(r);
    return ret;
  }

  void Realtime::unsubscribe(const std::string& roomId) Kuz_Throw_KuzzleException {
    void_result *r = kuzzle_realtime_unsubscribe(_realtime, const_cast<char*>(roomId.c_str()));
    if (r->error != NULL)
        throwExceptionFromStatus(r);
    kuzzle_free_void_result(r);
  }

  bool Realtime::validate(const std::string& index, const std::string collection, const std::string body) Kuz_Throw_KuzzleException {
    bool_result *r = kuzzle_realtime_validate(_realtime, const_cast<char*>(index.c_str()), const_cast<char*>(collection.c_str()),  const_cast<char*>(body.c_str()));
    if (r->error != NULL)
        throwExceptionFromStatus(r);
    bool ret = r->result;
    kuzzle_free_bool_result(r);
    return ret;
}
}
