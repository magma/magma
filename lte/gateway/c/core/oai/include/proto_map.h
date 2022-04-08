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
#include <google/protobuf/map.h>
#include <string>

#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_23.003.h"

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
** APIs:       set_name() , get_name(), get(), insert(), delete()         **
**                                                                        **
***************************************************************************/

template <typename keyT, typename valueT>
struct map_s {
  google::protobuf::Map<keyT, valueT>* map;
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
    if (map->empty()) {
      return MAP_EMPTY;
    }
    if (valueP == nullptr) {
      return MAP_BAD_PARAMETER_VALUE;
    }
    auto search_result = map->find(key);
    if (search_result != map->end()) {
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
    typedef typename google::protobuf::Map<keyT, valueT>::iterator itr;
    std::pair<itr, bool> insert_response =
        map->insert(google::protobuf::MapPair<keyT, valueT>(key, value));
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
    if (map->empty()) {
      return MAP_EMPTY;
    }

    if (map->erase(key)) {
      return MAP_OK;
    } else {
      return MAP_KEY_NOT_EXISTS;
    }
  }

  /***************************************************************************
  **                                                                        **
  ** Name:    isEmpty()                                                     **
  **                                                                        **
  ** Description: Returns true if map is empty, else returns false          **
  **                                                                        **
  ***************************************************************************/
  bool isEmpty() { return map->empty(); }

  /***************************************************************************
  **                                                                        **
  ** Name:    clear()                                                       **
  **                                                                        **
  ** Description: Clears the contents of the map                            **
  **                                                                        **
  ***************************************************************************/
  void clear() {
    map->clear();
    name.clear();
  }
  /***************************************************************************
  **                                                                        **
  ** Name:    size()                                                        **
  **                                                                        **
  ** Description: size the contents of the map                              **
  **                                                                        **
  ***************************************************************************/
  size_t size() { return map->size(); }
};

// Map- Key: uint64_t, Data: uint64_t
typedef magma::map_s<uint64_t, uint64_t> map_uint64_uint64_t;
// Map- Key: string, Data: uint64_t
typedef magma::map_s<std::string, uint64_t> map_string_uint64_t;

}  // namespace magma
