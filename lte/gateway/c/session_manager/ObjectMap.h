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

#include <string>
#include <vector>

namespace magma {

enum ObjectMapResult {
  SUCCESS = 0,
  CLIENT_ERROR = 1,
  KEY_NOT_FOUND = 2,
  INCORRECT_VALUE_TYPE = 3,
  DESERIALIZE_FAIL = 4,
  SERIALIZE_FAIL = 5,
};

/**
 * ObjectMap is an abstract class to represent any class that can store key,
 * object pairs
 */
template <typename ObjectType>
class ObjectMap {
  virtual ObjectMapResult set(const std::string& key,
                              const ObjectType& object) = 0;

  virtual ObjectMapResult get(const std::string& key,
                              ObjectType& object_out) = 0;

  virtual ObjectMapResult getall(std::vector<ObjectType>& values_out) = 0;
};

}  // namespace magma
