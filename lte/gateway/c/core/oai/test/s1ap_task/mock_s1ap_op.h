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

#define STR_DECODE_TABLE                                                       \
  "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
#define MAX_ASCII_VAL_OF_DECODE_TABLE 123
#define NO_CODE_VAL 128
#define SHIFT 128
#define LENGTH_OF_STR_DECODE_TABLE 64
#define ENCODER_BLOCK_SIZE 6
#define DECODER_BLOCK_SIZE 8
#define A_UPPER_ASCII 65
#define Z_UPPER_ASCII 90
#define A_LOWER_ASCII 97
#define Z_LOWER_ASCII 122
#define SLASH_ASCII 47
#define EQUAL_ASCII 61
#define PLUS_ASCII 43
#define ZERO_ASCII 48
#define NINE_ASCII 57
#define IGNORE 128
#define MAX_PAD_SIZE 2
#define PAD_SYMBOL '='

namespace magma {
namespace lte {

#define DEFAULT_S1AP_CONTEXT_DATA_PATH                                         \
  "lte/gateway/c/core/oai/test/s1ap_task/data/"

#define DEFAULT_S1AP_STATE_DATA_PATH                                           \
  "lte/gateway/c/core/oai/test/s1ap_task/data/s1ap_state_ATTACHED_encoded"

std::vector<std::string> load_file_into_vector_of_line_content(
    const std::string& data_folder_path, const std::string& file_name);
status_code_e mock_read_s1ap_ue_state_db(
    const std::vector<std::string>& ue_samples);
status_code_e mock_read_s1ap_state_db(
    const std::string& file_name_state_sample);
void add_data_to_buffer(
    unsigned int& buf, int& buf_size, int data,
    const std::vector<int>& decode_table);
int get_last_decoder_block_size_bit(int num);
void add_decoded_data_to_output(
    unsigned int& buf, int& buf_size, std::string& buf_decoded_msg);
std::string decode_msg(const std::vector<char>& encoded_msg);

}  // namespace lte
}  // namespace magma
