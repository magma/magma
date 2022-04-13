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

#include <iosfwd>  // for ifstream
#include <string>  // for stringnamespace magma {

namespace magma {
namespace service303 {

/**
 * ProcFileUtils is a helper class to parse proc files for process information
 */
class ProcFileUtils final {
 public:
  /*
   * memory_info_t wraps the needed information from the /status file
   */
  typedef struct memory_info_s {
    double physical_mem = -1;
    double virtual_mem = -1;
  } memory_info_t;

 public:
  /*
   * Parses the /proc/self/status file for information on memory usage
   *
   * @return memory_info_t containing virtual and physical memory usage
   */
  static const memory_info_t getMemoryInfo();

 private:
  /*
   * Helper function to read from the proc file stream and output the value if
   * the prefix is found
   *
   * @return -1 if the string isn't the prefix we're looking for, otherwise
   *    return the actual value
   */
  static double parseForPrefix(std::ifstream& infile,
                               const std::string& to_compare,
                               const std::string& prefix_name);

 private:
  static const std::string STATUS_FILE;
  // status file labels
  static const std::string VIRTUAL_MEM_PREFIX;
  static const std::string PHYSICAL_MEM_PREFIX;
};

}  // namespace service303
}  // namespace magma
