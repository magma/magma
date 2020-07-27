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

#include <chrono>
#include <ctime>
#include <iomanip>
#include <ostream>

namespace devmand {
namespace test {
namespace utils {
namespace json {

using namespace std;

extern string sortJson(const string& json);

} // namespace json
} // namespace utils
} // namespace test
} // namespace devmand
