/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

#pragma once

#include <string>

namespace devmand {
namespace devices {
namespace cli {

using namespace std;

class Model {
 public:
  Model() = delete;
  ~Model() = default;

 protected:
  explicit Model(string _dir) : dir(_dir) {}

 private:
  const string dir;

 public:
  const string& getDir() const {
    return dir;
  }
  bool operator<(const Model& x) const {
    return dir < x.dir;
  }

  static const Model OPENCONFIG_2_4_3;
};

} // namespace cli
} // namespace devices
} // namespace devmand
