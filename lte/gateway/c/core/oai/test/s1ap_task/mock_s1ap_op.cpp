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

#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_state_converter.hpp"

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
  if (!state_ue_ht) {
    std::cerr << "Cannot get S1AP UE State" << std::endl;
    return RETURNerror;
  }

  for (const auto& name_of_sample_file : ue_samples) {
    UeDescription ue_proto = UeDescription();
    std::fstream input(name_of_sample_file.c_str(),
                       std::ios::in | std::ios::binary);
    if (!ue_proto.ParseFromIstream(&input)) {
      std::cerr << "Failed to parse the sample: " << name_of_sample_file
                << std::endl;
      return RETURNerror;
    }

    ue_description_t* ue_context_p = reinterpret_cast<ue_description_t*>(
        calloc(1, sizeof(ue_description_t)));
    S1apStateConverter::proto_to_ue(ue_proto, ue_context_p);

    hashtable_rc_t h_rc =
        hashtable_ts_insert(state_ue_ht, ue_context_p->comp_s1ap_id,
                            reinterpret_cast<void*>(ue_context_p));

    if (HASH_TABLE_OK != h_rc) {
      std::cerr << "Failed to insert UE state :" << name_of_sample_file
                << std::endl;
      return RETURNerror;
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

  std::ifstream input(file_name_state_sample.c_str(),
                      std::ios::in | std::ios::binary);
  std::vector<char> data = std::vector<char>(
      std::istreambuf_iterator<char>(input), std::istreambuf_iterator<char>());
  std::string decoded_msg = decode_msg(data);
  std::istringstream input_stream(decoded_msg);
  if (!state_proto.ParseFromIstream(&input_stream)) {
    std::cerr << "Failed to decode and parse the sample: "
              << file_name_state_sample << std::endl;
    return RETURNerror;
  }
  S1apStateConverter::proto_to_state(state_proto, state_cache_p);
  return RETURNok;
}

// putting data in the buffer to decode and update the size of the buffer
void add_data_to_buffer(unsigned int& buf, int& buf_size, int data,
                        const std::vector<int>& decode_table) {
  buf = buf << ENCODER_BLOCK_SIZE;
  buf += data;
  buf_size += ENCODER_BLOCK_SIZE;
}

// getting the last byte from the buffer
int get_last_decoder_block_size_bit(int num) {
  return num & ((1 << DECODER_BLOCK_SIZE) - 1);
}

// decoding each byte and putting it into the buffer of the decoded message
void add_decoded_data_to_output(unsigned int& buf, int& buf_size,
                                std::string& buf_decoded_msg) {
  if (buf_size >= DECODER_BLOCK_SIZE) {
    buf_size -= DECODER_BLOCK_SIZE;
    char decoded_char =
        static_cast<char>(get_last_decoder_block_size_bit(buf >> buf_size));
    buf_decoded_msg.push_back(decoded_char);
  }
}

// decoding an encoded message
std::string decode_msg(const std::vector<char>& encoded_msg) {
  if (encoded_msg.size() % 4 != 0) {
    std::cerr << " This data is wrongly encoded." << std::endl;
    return "";
  }
  std::vector<int> decode_table(MAX_ASCII_VAL_OF_DECODE_TABLE, NO_CODE_VAL);
  for (int i = 0; i < LENGTH_OF_STR_DECODE_TABLE; i++)
    decode_table[STR_DECODE_TABLE[i]] = i;

  int pad_size = 0;
  int cur_len_buffer = 0;
  unsigned int buffer = 0;
  std::string decoded_msg;
  for (int i = 0; i < encoded_msg.size(); ++i) {
    char ch = encoded_msg[i];
    ch -= SHIFT;
    int data = decode_table[ch];
    if (data == NO_CODE_VAL) {
      if (ch != PAD_SYMBOL) {
        std::cerr << " Invalid char in the encoded message." << std::endl;
        return "";
      }
      break;
    }
    add_data_to_buffer(buffer, cur_len_buffer, data, decode_table);
    add_decoded_data_to_output(buffer, cur_len_buffer, decoded_msg);
  }
  return decoded_msg;
}

}  // namespace lte
}  // namespace magma