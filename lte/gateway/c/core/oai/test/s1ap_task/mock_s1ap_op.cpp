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

#include "lte/gateway/c/core/oai/test/s1ap_task/mock_s1ap_op.h"

#include <google/protobuf/text_format.h>

#include <fstream>
#include <algorithm>
#include <sstream>

#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_mme.hpp"

namespace magma {
namespace lte {

// loads paths of data samples of s1ap states (storing in each line) from a file
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

// mocking the reading s1ap ue states from redis database by injecting local
// samples
status_code_e mock_read_s1ap_ue_state_db(
    const std::vector<std::string>& ue_samples) {
  map_uint64_ue_description_t* state_ue_map = get_s1ap_ue_state();
  if (!state_ue_map) {
    std::cerr << "Cannot get S1AP UE State" << std::endl;
    return RETURNerror;
  }

  for (const auto& name_of_sample_file : ue_samples) {
    oai::UeDescription ue_proto = oai::UeDescription();
    std::fstream input(name_of_sample_file.c_str(),
                       std::ios::in | std::ios::binary);
    if (!ue_proto.ParseFromIstream(&input)) {
      std::cerr << "Failed to parse the sample: " << name_of_sample_file
                << std::endl;
      return RETURNerror;
    }

    oai::UeDescription* ue_context_p = new oai::UeDescription();
    if (!ue_context_p) {
      std::cerr << "Failed to allocate memory for ue_context_p" << std::endl;
      return RETURNerror;
    }
    ue_context_p->MergeFrom(ue_proto);

    proto_map_rc_t rc = state_ue_map->insert(
        ue_context_p->comp_s1ap_id(),
        reinterpret_cast<oai::UeDescription*>(ue_context_p));

    if (rc != magma::PROTO_MAP_OK) {
      std::cerr << "Failed to insert UE state :" << name_of_sample_file
                << std::endl;
      free_ue_description(reinterpret_cast<void**>(&ue_context_p));
      return RETURNerror;
    }
  }
  return RETURNok;
}

// mocking the reading s1ap ue states from redis database by injecting local
// samples
status_code_e mock_read_s1ap_state_db(
    const std::string& file_name_state_sample) {
  oai::S1apState* state_cache_p = get_s1ap_state(false);

  oai::S1apState state_proto = oai::S1apState();

  std::ifstream input(file_name_state_sample.c_str(),
                      std::ios::in | std::ios::binary);
  if (!state_proto.ParseFromIstream(&input)) {
    std::cerr << "Failed to decode and parse the sample: "
              << file_name_state_sample << std::endl;
    return RETURNerror;
  }
  state_cache_p->Clear();
  state_cache_p->MergeFrom(state_proto);
  return RETURNok;
}

}  // namespace lte
}  // namespace magma
