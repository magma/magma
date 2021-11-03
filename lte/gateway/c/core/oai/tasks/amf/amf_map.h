
#pragma once
#include <pthread.h>

#include <iostream>
#include <unordered_map>
#include <mutex>
#include <string>

#include "3gpp_23.003.h"

namespace std {
typedef size_t hash_size_t;
/*Default hash function*/
static hash_size_t amf_def_hashfunc(
    const void* const keyP, const int key_sizeP) {
  hash_size_t hash = 0;
  int key_size     = key_sizeP;

  // may use MD4 ?
  while (key_size > 0) {
    uint32_t val = 0;
    int size     = sizeof(val);
    while ((size > 0) && (key_size > 0)) {
      val = val << 8;
      val |= (reinterpret_cast<const uint8_t*>(keyP))[key_size - 1];
      size--;
      key_size--;
    }
    hash ^= val;
  }

  return hash;
}
// specialise std::equal_to function for type guti_m5_t
template<>
struct equal_to<guti_m5_t> {
  bool operator()(const guti_m5_t& guti1, const guti_m5_t& guti2) const {
    return (guti1.m_tmsi == guti2.m_tmsi);
  }
};

// specialise std::hash function for type guti_m5_t
// HashFunction for Guti as key and unit64 data
// Note: This HashFunc in turn calls the def_hashfunc which supports hashing
// only for unit64
template<>
struct hash<guti_m5_t> {
  size_t operator()(const guti_m5_t& k) const {
    return amf_def_hashfunc(&k, sizeof(k));
  }
};
}  // namespace std
namespace magma5g {

/*Enum for Map Return code
  Note: If new named constant is added to the enumeration list, add the new case
  in map_rc_code2string() located in amf_map.cpp*/
typedef enum map_return_code_e {
  MAP_OK = 0,
  MAP_KEY_NOT_EXISTS,
  MAP_SEARCH_NO_RESULT,
  MAP_KEY_ALREADY_EXISTS,
  MAP_BAD_PARAMETER_KEY,
  MAP_BAD_PARAMETER_VALUE,
  MAP_EMPTY,
  MAP_DUMP_FAIL
} map_rc_t;

/***************************************************************************
**                                                                        **
** Name:    map_ts_s                                                      **
**                                                                        **
** Description: This is a generic structure for maps.It is implemented    **
**              using template struct definitions.                        **
**                                                                        **
** Parameters:  keyT     - data type of key                               **
**              valueT   - data type of value                             **
**              Hash     - used to pass a custom hash function.           **
**                         Default: std::hash<keyT>                       **
**              KeyEqual - it is used to overload the EqualTo operator.   **
**                         When using a userdefined struct as key,        **
**                         pass a suitable EqualTo method.                **
**                            Default: std::equal_to<keyT>                **
**                                                                        **
** APIs:       set_name() , get_name(), get(), insert(), delete()        **
**                                                                        **
***************************************************************************/
template<
    typename keyT, typename valueT, class Hash = std::hash<keyT>,
    class KeyEqual = std::equal_to<keyT>>
struct map_ts_s {
  std::mutex mutex_obj;
  std::unordered_map<keyT, valueT, Hash, KeyEqual> umap;
  std::string name;
  bool log_enabled;

  void set_name(std::string umap_name) { name = umap_name; }
  std::string get_name() { return name; }

  /***************************************************************************
  **                                                                        **
  ** Name:    get()                                                         **
  **                                                                        **
  ** Description: Takes key and valueP as parameters.If the key exists,     **
  **              corresponding value is returned through the valueP,       **
  **              else returns error.                                       **
  **                                                                        **
  ***************************************************************************/
  map_rc_t get(const keyT key, valueT* valueP) {
    if (umap.empty()) {
      return MAP_EMPTY;
    }
    if (valueP == nullptr) {
      return MAP_BAD_PARAMETER_VALUE;
    }
    std::lock_guard<std::mutex> lock(mutex_obj);
    auto search_result = umap.find(key);
    if (search_result != umap.end()) {
      *valueP = search_result->second;
      return MAP_OK;
    } else {
      return MAP_KEY_NOT_EXISTS;
    }
  }

  /***************************************************************************
  **                                                                        **
  ** Name:    insert()                                                      **
  **                                                                        **
  ** Description: Takes key and value as parameters.Inserts the <key,value> **
  **              pair into the map. If identical key already exists then   **
  **              returns error.                                            **
  **                                                                        **
  ***************************************************************************/
  map_rc_t insert(const keyT key, const valueT value) {
    typedef typename std::unordered_map<keyT, valueT>::iterator itr;
    std::lock_guard<std::mutex> lock(mutex_obj);
    std::pair<itr, bool> insert_response =
        umap.insert(std::make_pair(key, value));
    if (insert_response.second) {
      return MAP_OK;
    } else {
      return MAP_KEY_ALREADY_EXISTS;
    }
  }

  /***************************************************************************
  **                                                                        **
  ** Name:    remove()                                                      **
  **                                                                        **
  ** Description: Takes key parameter.Removes the corresponding entry from  **
  **              the map. If key does not exists returns error             **
  **                                                                        **
  ***************************************************************************/
  map_rc_t remove(const keyT key) {
    if (umap.empty()) {
      return MAP_EMPTY;
    }

    std::lock_guard<std::mutex> lock(mutex_obj);
    if (umap.erase(key)) {
      return MAP_OK;
    } else {
      return MAP_KEY_NOT_EXISTS;
    }
  }
};

// Helper Function Declarations:
// Function to convert the return codes into string before logging them.
std::string map_rc_code2string(map_rc_t rc);

// Map Declarations:
// Map- Key: uint64_t , Data: uint64_t
typedef map_ts_s<uint64_t, uint64_t> map_uint64_ts_t;

// Map Key: guti_m5_t Data: uint64_t;
typedef map_ts_s<guti_m5_t, uint64_t> obj_map_uint64_ts_t;

}  // namespace magma5g
