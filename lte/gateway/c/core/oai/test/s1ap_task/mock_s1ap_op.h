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

#include "lte/gateway/c/core/oai/common/common_defs.h"

namespace magma {
namespace lte {

#define DEFAULT_S1AP_CONTEXT_DATA_PATH                                         \
  "lte/gateway/c/core/oai/test/s1ap_task/data/"

#define DEFAULT_S1AP_STATE_DATA_PATH                                         \
  "lte/gateway/c/core/oai/test/s1ap_task/data/s1ap_state_ATTACHED"

std::vector<std::string> load_file_into_vector_of_line_content(
    const std::string& data_folder_path, const std::string& file_name);
status_code_e mock_read_s1ap_ue_state_db(
    const std::vector<std::string>& ue_samples);
status_code_e mock_read_s1ap_state_db(
  const std::string& file_name_state_sample);

}  // namespace lte
}  // namespace magma
