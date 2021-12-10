/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#pragma once
#include <pthread.h>

#include <iostream>
#include <unordered_map>
#include <string>

#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_23.003.h"

namespace std {
typedef size_t hash_size_t;
/*Default hash function*/
static hash_size_t def_hashfunc(const void* const keyP, const int key_sizeP) {
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

/*specialise std::hash function for type guti_m5_t
  HashFunction for Guti as key and unit64 data
  Note: This HashFunc in turn calls the def_hashfunc which supports hashing
  only for unit64*/
template<>
struct hash<guti_m5_t> {
  size_t operator()(const guti_m5_t& k) const {
    return def_hashfunc(&k, sizeof(k));
  }
};
}  // namespace std
namespace magma {

/*Enum for Map Return code
  Note: If new named constant is added to the enumeration list, add the new case
  in map_rc_code2string()*/
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
** Name:    map_rc_code2string()                                          **
**                                                                        **
** Description: This converts the map_rc_t, return code to string         **
**                                                                        **
***************************************************************************/

static std::string map_rc_code2string(map_rc_t rc) {
  switch (rc) {
    case MAP_OK:
      return "MAP_OK";
      break;

    case MAP_KEY_NOT_EXISTS:
      return "MAP_KEY_NOT_EXISTS";
      break;

    case MAP_SEARCH_NO_RESULT:
      return "MAP_SEARCH_NO_RESULT";
      break;

    case MAP_KEY_ALREADY_EXISTS:
      return "MAP_KEY_ALREADY_EXISTS";
      break;

    case MAP_BAD_PARAMETER_KEY:
      return "MAP_BAD_PARAMETER_KEY";
      break;

    case MAP_BAD_PARAMETER_VALUE:
      return "MAP_BAD_PARAMETER_VALUE";
      break;

    case MAP_EMPTY:
      return "MAP_EMPTY";
      break;

    case MAP_DUMP_FAIL:
      return "MAP_DUMP_FAIL";
      break;

    default:
      return "UNKNOWN map_rc_t";
  }
}

/***************************************************************************
**                                                                        **
** Name:    map_s                                                         **
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
struct map_s {
  std::unordered_map<keyT, valueT, Hash, KeyEqual> umap;
  std::string name;
  bool log_enabled = false;

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

    if (umap.erase(key)) {
      return MAP_OK;
    } else {
      return MAP_KEY_NOT_EXISTS;
    }
  }

  /***************************************************************************
  **                                                                        **
  ** Name:    isEmpty()                                                     **
  **                                                                        **
  ** Description: Returns true if map is empty, else returns false           **
  **                                                                        **
  ***************************************************************************/
  bool isEmpty() {
    if (umap.empty()) {
      return true;
    }
    return false;
  }

  /***************************************************************************
  **                                                                        **
  ** Name:    clear()                                                       **
  **                                                                        **
  ** Description: Clears the contents of the map                            **
  **                                                                        **
  ***************************************************************************/
  void clear() {
    umap.clear();
    name.clear();
  }
};

// Amf-Map Declarations:
// Map- Key: uint64_t , Data: uint64_t
typedef magma::map_s<uint64_t, uint64_t> map_uint64_uint64_t;

}  // namespace magma
