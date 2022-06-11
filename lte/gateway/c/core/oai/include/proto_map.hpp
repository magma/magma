/**
 * Copyright 2022 The Magma Authors.
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
#include <string>
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_23.003.h"

#include <iostream>
#include <google/protobuf/map.h>

namespace magma {

/*Enum for Map Return code
  Note: If new named constant is added to the enumeration list, add the new case
  in map_rc_code2string()*/
typedef enum proto_map_return_code_e {
  PROTO_MAP_OK = 0,
  PROTO_MAP_KEY_NOT_EXISTS,
  PROTO_MAP_SEARCH_NO_RESULT,
  PROTO_MAP_KEY_ALREADY_EXISTS,
  PROTO_MAP_BAD_PARAMETER_KEY,
  PROTO_MAP_BAD_PARAMETER_VALUE,
  PROTO_MAP_EMPTY,
  PROTO_MAP_DUMP_FAIL
} proto_map_rc_t;

/***************************************************************************
**                                                                        **
** Name:    map_rc_code2string()                                          **
**                                                                        **
** Description: This converts the proto_map_rc_t, return code to string   **
**                                                                        **
***************************************************************************/

static std::string map_rc_code2string(proto_map_rc_t rc) {
  switch (rc) {
    case PROTO_MAP_OK:
      return "MAP_OK";
      break;

    case PROTO_MAP_KEY_NOT_EXISTS:
      return "MAP_KEY_NOT_EXISTS";
      break;

    case PROTO_MAP_SEARCH_NO_RESULT:
      return "MAP_SEARCH_NO_RESULT";
      break;

    case PROTO_MAP_KEY_ALREADY_EXISTS:
      return "MAP_KEY_ALREADY_EXISTS";
      break;

    case PROTO_MAP_BAD_PARAMETER_KEY:
      return "MAP_BAD_PARAMETER_KEY";
      break;

    case PROTO_MAP_BAD_PARAMETER_VALUE:
      return "MAP_BAD_PARAMETER_VALUE";
      break;

    case PROTO_MAP_EMPTY:
      return "MAP_EMPTY";
      break;

    case PROTO_MAP_DUMP_FAIL:
      return "MAP_DUMP_FAIL";
      break;

    default:
      return "UNKNOWN proto_map_rc_t";
  }
}

/***************************************************************************
**                                                                        **
** Name:    proto_map_s                                                   **
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
struct proto_map_s {
  google::protobuf::Map<keyT, valueT>* map;
  /* TODO (rsarwad): on final conversion to cpp,
   replace char array with std::string */
  char name[1024];

  void set_name(char* umap_name) {
    strncpy(name, umap_name, strlen(umap_name));
  }
  char* get_name() { return name; }
  /***************************************************************************
  **                                                                        **
  ** Name:    get()                                                         **
  **                                                                        **
  ** Description: Takes key and valueP as parameters.If the key exists,     **
  **              corresponding value is returned through the valueP,       **
  **              else returns error.                                       **
  **                                                                        **
  ***************************************************************************/

  proto_map_rc_t get(const keyT key, valueT* valueP) {
    if (map->empty()) {
      return PROTO_MAP_EMPTY;
    }
    if (valueP == nullptr) {
      return PROTO_MAP_BAD_PARAMETER_VALUE;
    }
    auto search_result = map->find(key);
    if (search_result != map->end()) {
      *valueP = search_result->second;
      return PROTO_MAP_OK;
    } else {
      return PROTO_MAP_KEY_NOT_EXISTS;
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
  proto_map_rc_t insert(const keyT key, const valueT value) {
    typedef typename google::protobuf::Map<keyT, valueT>::iterator itr;
    std::pair<itr, bool> insert_response =
        map->insert(google::protobuf::MapPair<keyT, valueT>(key, value));
    if (insert_response.second) {
      return PROTO_MAP_OK;
    } else {
      return PROTO_MAP_KEY_ALREADY_EXISTS;
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
  proto_map_rc_t remove(const keyT key) {
    if (map->empty()) {
      return PROTO_MAP_EMPTY;
    }

    if (map->erase(key)) {
      return PROTO_MAP_OK;
    } else {
      return PROTO_MAP_KEY_NOT_EXISTS;
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
    if (!(map->empty())) {
      map->clear();
    }
  }
  /***************************************************************************
  **                                                                        **
  ** Name:    size()                                                        **
  **                                                                        **
  ** Description: size the contents of the map                              **
  **                                                                        **
  ***************************************************************************/
  size_t size() { return map->size(); }

  /***************************************************************************
  **                                                                        **
  ** Name:    destroy_map()                                                 **
  **                                                                        **
  ** Description: Clears the contents of the map and also delete map        **
  **                                                                        **
  ***************************************************************************/
  void destroy_map() {
    map->clear();
    delete map;
  }

  /***************************************************************************
  **                                                                        **
  ** Name:    map_apply_callback_on_all_elements()                          **
  **                                                                        **
  ** Description: Traverses through map and call callback function to be    **
  **              executed on each node                                     **
  **                                                                        **
  ***************************************************************************/
  proto_map_rc_t map_apply_callback_on_all_elements(
      bool funct_cb(const keyT key, const valueT value, void* parameterP,
                    void** resultP),
      void* parameterP, void** resultP) {
    if (map->empty()) {
      return PROTO_MAP_EMPTY;
    }
    for (auto itr = map->begin(); itr != map->end(); itr++) {
      if (funct_cb(itr->first, itr->second, parameterP, resultP)) {
        return PROTO_MAP_OK;
      }
    }
    return PROTO_MAP_DUMP_FAIL;
  }
};

// Map- Key: uint64_t, Data: uint64_t
typedef magma::proto_map_s<uint64_t, uint64_t> proto_map_uint64_uint64_t;
// Map- Key: string, Data: uint64_t
typedef magma::proto_map_s<std::string, uint64_t> proto_map_string_uint64_t;
// Map- Key: uint32_t, Data: uint64_t
typedef magma::proto_map_s<uint32_t, uint64_t> proto_map_uint32_uint64_t;

}  // namespace magma
