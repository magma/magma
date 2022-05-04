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
#include <iostream>
#include <fstream>
#include <string>
#include <vector>

#include "lte/gateway/c/core/oai/tasks/sgw/spgw_state_converter.hpp"

namespace magma {
namespace lte {

#define DEFAULT_SPGW_CONTEXT_DATA_PATH \
  "lte/gateway/c/core/oai/test/spgw_task/data/"

std::vector<std::string> load_file_into_vector_of_line_content(
    const std::string& data_folder_path, const std::string& file_name);
status_code_e mock_read_spgw_ue_state_db(
    const std::vector<std::string>& ue_samples);
}  // namespace lte
}  // namespace magma
