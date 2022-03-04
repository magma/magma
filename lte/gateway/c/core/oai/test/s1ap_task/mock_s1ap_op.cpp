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

#include <fstream>
#include <algorithm>

#include <google/protobuf/text_format.h>
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_state_converter.h"

using magma::lte::oai::S1apState;
using magma::lte::oai::UeDescription;

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
  hash_table_ts_t* state_ue_ht = get_s1ap_ue_state();

  for (const auto& name_of_sample_file : ue_samples) {
    UeDescription ue_proto = UeDescription();
    std::fstream input(
        name_of_sample_file.c_str(), std::ios::in | std::ios::binary);
    if (!ue_proto.ParseFromIstream(&input)) {
      std::cerr << "Failed to parse the sample: " << name_of_sample_file
                << std::endl;
      return RETURNerror;
    }

    ue_description_t* ue_context_p = reinterpret_cast<ue_description_t*>(
        calloc(1, sizeof(ue_description_t)));
    S1apStateConverter::proto_to_ue(ue_proto, ue_context_p);

    hashtable_rc_t h_rc = hashtable_ts_insert(
        state_ue_ht, ue_context_p->comp_s1ap_id, (void*) ue_context_p);

    if (HASH_TABLE_OK != h_rc) {
      std::cerr << "Failed to insert UE state :" << name_of_sample_file
                << std::endl;
    }
  }
  return RETURNok;
}

// mocking the reading s1ap ue states from redis database by injecting local
// samples
status_code_e mock_read_s1ap_state_db(
    const std::string& file_name_state_sample) {
  s1ap_state_t* state_cache_p = get_s1ap_state(false);

  S1apState state_proto = S1apState();

  std::fstream input(
      file_name_state_sample.c_str(), std::ios::in | std::ios::binary);
  if (!state_proto.ParseFromIstream(&input)) {
    std::cerr << "Failed to parse the sample: " << file_name_state_sample
              << std::endl;
    return RETURNerror;
  }
  S1apStateConverter::proto_to_state(state_proto, state_cache_p);
  return RETURNok;
}

}  // namespace lte
}  // namespace magma