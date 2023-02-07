/*
 * Copyright 2023 The Magma Authors.
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_ue_context_and_proc.hpp"

#define kMinInputLength 10
#define kMaxInputLength 5120

extern "C" int LLVMFuzzerTestOneInput(const uint8_t *Data, size_t Size) {
/*
* amf/test_amf_encode_decode.cpp
*/

    if (Size < kMinInputLength || Size > kMaxInputLength) {
        return 0;
    }

    magma5g::amf_nas_message_t decode_msg = {};
    magma5g::amf_nas_message_decode_status_t decode_status = {};
    int status = RETURNerror;

    status = magma5g::nas5g_message_decode(Data, &decode_msg, Size, nullptr, &decode_status);

    return status;
}
