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
#include <fstream>
#include <algorithm>

#include "lte/gateway/c/core/oai/test/spgw_task/mock_spgw_op.h"
#include "lte/gateway/c/core/oai/include/spgw_state.h"
#include <google/protobuf/text_format.h>

namespace magma {
namespace lte {

std::vector<std::string> load_file_into_vector_of_line_content(
    const std::string& data_folder_path, const std::string& file_name) {
  std::fstream file_content(file_name.c_str(), std::ios_base::in);
  std::string line;
  std::vector<std::string> vector_of_lines;
  if (file_content) {
    while (std::getline(file_content, line)) {
      line.erase(std::remove(line.begin(), line.end(), '\n'), line.end());
      vector_of_lines.push_back(data_folder_path + "/" + line);
    }
  } else {
    std::cerr << "couldn't open file: " << file_name << std::endl;
  }

  file_content.close();
  return vector_of_lines;
}

status_code_e mock_read_spgw_ue_state_db(
    const std::vector<std::string>& ue_samples) {
  for (const auto& name_of_sample : ue_samples) {
    oai::SpgwUeContext ue_proto = oai::SpgwUeContext();
    std::fstream input(name_of_sample.c_str(), std::ios::in | std::ios::binary);
    if (!ue_proto.ParseFromIstream(&input)) {
      std::cerr << "Failed to parse the sample: " << name_of_sample
                << std::endl;
      return RETURNerror;
    }

    spgw_ue_context_t* ue_context_p =
        (spgw_ue_context_t*) calloc(1, sizeof(spgw_ue_context_t));
    SpgwStateConverter::proto_to_ue(ue_proto, ue_context_p);
  }
  return RETURNok;
}

}  // namespace lte
}  // namespace magma
