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

#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_ue_context_and_proc.hpp"

#define PDU_SESSION_DEFAULT_AMBR 1

namespace magma5g {
// Set the ue aggregate max bit rate
status_code_e amf_smf_context_ue_aggregate_max_bit_rate_set(
    amf_context_s* amf_ctxt_p, ambr_t subscribed_ue_ambr);

// Get the ue aggregate max bit rate
status_code_e amf_smf_context_ue_aggregate_max_bit_rate_get(
    const amf_context_s* amf_ctxt_p, bit_rate_t* subscriber_ambr_dl,
    bit_rate_t* subscriber_ambr_ul);

void amf_smf_get_slice_configuration(std::shared_ptr<smf_context_t> smf_ctx,
                                     s_nssai_t* slice_config);
}  // namespace magma5g
