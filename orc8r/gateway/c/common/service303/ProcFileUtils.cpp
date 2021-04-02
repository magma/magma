
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

#include "ProcFileUtils.h"
#include <fstream>  // for basic_istream, ifstream
#include <string>   // for string, operator>>, stod

namespace magma {
namespace service303 {

const std::string ProcFileUtils::STATUS_FILE         = "/proc/self/status";
const std::string ProcFileUtils::VIRTUAL_MEM_PREFIX  = "VmSize:";
const std::string ProcFileUtils::PHYSICAL_MEM_PREFIX = "VmRSS:";

double ProcFileUtils::parseForPrefix(
    std::ifstream& infile, const std::string& to_compare,
    const std::string& prefix_name) {
  if (to_compare.compare(prefix_name) == 0) {
    std::string value_string;
    infile >> value_string;
    // KiB -> bytes
    return std::stod(value_string) * 1024;
  }
  return -1;
}

const ProcFileUtils::memory_info_t ProcFileUtils::getMemoryInfo() {
  std::ifstream infile(ProcFileUtils::STATUS_FILE);
  ProcFileUtils::memory_info_t info;
  std::string content;
  // Parse file token by token until prefixes are found
  while (infile >> content) {
    double value;
    // look for and set virtual_mem
    value = ProcFileUtils::parseForPrefix(
        infile, content, ProcFileUtils::VIRTUAL_MEM_PREFIX);
    if (value >= 0) {
      info.virtual_mem = value;
    }
    // look for and set physical_mem
    value = ProcFileUtils::parseForPrefix(
        infile, content, ProcFileUtils::PHYSICAL_MEM_PREFIX);
    if (value >= 0) {
      info.physical_mem = value;
    }

    if (info.virtual_mem >= 0 && info.physical_mem >= 0) {
      break;
    }
  }
  return info;
}

}  // namespace service303
}  // namespace magma
