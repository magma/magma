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

#include <random>
#include <string>

class SessionIDGenerator {
 public:
  SessionIDGenerator();

  /**
   * Generates a random session id from the IMSI in the form
   * "<IMSI>-<RANDOM NUM>"
   */
  std::string gen_session_id(const std::string& imsi);

  /**
   * Parses an IMSI value from a session_id
   */
  bool get_imsi_from_session_id(const std::string& session_id,
                                std::string& imsi_out);

 private:
  std::mt19937 rgen_;
  std::uniform_int_distribution<int> idist_;
};
