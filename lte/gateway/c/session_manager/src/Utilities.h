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

#include <google/protobuf/timestamp.pb.h>
#include <google/protobuf/util/time_util.h>

#include <string>

namespace magma {
std::string bytes_to_hex(const std::string& s);
uint64_t get_time_in_sec_since_epoch();
std::chrono::milliseconds time_difference_from_now(
    const google::protobuf::Timestamp& timestamp);
std::chrono::milliseconds time_difference_from_now(const std::time_t timestamp);
}  // namespace magma
