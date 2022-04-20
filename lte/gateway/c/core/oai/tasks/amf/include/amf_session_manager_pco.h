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
#pragma once

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.008.h"
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
#ifdef __cplusplus
};
#endif

namespace magma5g {

#define SM_PCO_IPCP_HDR_LENGTH 7

void sm_clear_protocol_configuration_options(
    protocol_configuration_options_t* const pco);

void sm_free_protocol_configuration_options(
    protocol_configuration_options_t** const protocol_configuration_options);

void sm_copy_protocol_configuration_options(
    protocol_configuration_options_t* const pco_dst,
    const protocol_configuration_options_t* const pco_src);

uint16_t sm_process_pco_request(protocol_configuration_options_t* pco_req,
                                protocol_configuration_options_t* pco_resp);
}  // namespace magma5g
