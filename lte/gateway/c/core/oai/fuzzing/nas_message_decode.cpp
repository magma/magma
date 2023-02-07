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
#include "lte/gateway/c/core/oai/tasks/nas/api/network/nas_message.hpp"


#define kMinInputLength 10
#define kMaxInputLength 5120

extern "C" int LLVMFuzzerTestOneInput(const uint8_t *Data, size_t Size) {
/*
* mme_app_task/mme_procedure_test_fixture.cpp
* mme_app_task/test_mme_app_esm_encode_decode.cpp
*/

    if (Size < kMinInputLength || Size > kMaxInputLength) {
        return 0;
    }

    nas_message_t nas_msg_decoded = {};
    emm_security_context_t emm_security_context;
    nas_message_decode_status_t decode_status;
    int decoder_rc = 0;

    decoder_rc = nas_message_decode(Data, &nas_msg_decoded, Size,
                                            reinterpret_cast<void*>(&emm_security_context), &decode_status);

    return decoder_rc;
}
