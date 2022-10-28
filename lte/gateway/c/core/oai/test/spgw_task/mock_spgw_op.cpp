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

#include "lte/gateway/c/core/oai/test/spgw_task/mock_spgw_op.hpp"
#include "lte/gateway/c/core/oai/include/spgw_state.hpp"
#include <google/protobuf/text_format.h>

namespace magma {
namespace lte {

// loads paths of data samples of spgw states (storing in each line) from a file
std::vector<std::string> load_file_into_vector_of_line_content(
    const std::string& data_folder_path, const std::string& file_name) {
  std::fstream file_content(file_name.c_str(), std::ios_base::in);
  std::string data_file_name;
  std::vector<std::string> vector_of_file_names;
  if (file_content) {
    while (std::getline(file_content, data_file_name)) {
      data_file_name.erase(
          std::remove(data_file_name.begin(), data_file_name.end(), '\n'),
          data_file_name.end());
      vector_of_file_names.push_back(data_folder_path + "/" + data_file_name);
    }
  } else {
    std::cerr << "couldn't open file: " << file_name << std::endl;
  }

  file_content.close();
  return vector_of_file_names;
}

// mocking the reading spgw ue states from redis database by injecting local
// samples
status_code_e mock_read_spgw_ue_state_db(
    const std::vector<std::string>& ue_samples) {
  for (const auto& name_of_sample_file : ue_samples) {
    oai::SpgwUeContext ue_proto = oai::SpgwUeContext();
    std::fstream input(name_of_sample_file.c_str(),
                       std::ios::in | std::ios::binary);
    if (!ue_proto.ParseFromIstream(&input)) {
      std::cerr << "Failed to parse the sample: " << name_of_sample_file
                << std::endl;
      return RETURNerror;
    }

    spgw_ue_context_t* ue_context_p = new spgw_ue_context_t();
    SpgwStateConverter::proto_to_ue(ue_proto, ue_context_p);
  }
  return RETURNok;
}

}  // namespace lte
}  // namespace magma
