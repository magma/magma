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
#include <sstream>
#include <string>

#include "SessionID.h"

SessionIDGenerator::SessionIDGenerator() {
  std::random_device rseed;
  rgen_ = std::mt19937(rseed());  // mersenne_twister
  idist_ = std::uniform_int_distribution<int>(0, 999999);
}

std::string SessionIDGenerator::gen_session_id(const std::string& imsi) {
  // imsi- + random 6 digit number
  return imsi + "-" + std::to_string(idist_(rgen_));
}

bool SessionIDGenerator::get_imsi_from_session_id(const std::string& session_id,
                                                  std::string& imsi_out) {
  std::istringstream ss(session_id);
  if (std::getline(ss, imsi_out, '-')) {
    return true;
  }
  return false;
}
