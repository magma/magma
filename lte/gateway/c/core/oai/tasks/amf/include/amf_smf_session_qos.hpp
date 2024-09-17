/**
 * Copyright 2022 The Magma Authors.
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
#include <memory>

#ifdef __cplusplus
extern "C" {
#endif

#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/common/dynamic_memory_check.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GQOSRules.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_ue_context_and_proc.hpp"

namespace magma5g {

// Fill the qos rule for the filter to be deleted
void amf_smf_session_api_fill_delete_packet_filter(
    delete_packet_filter_t* delete_pkt_filter, QOSRule* qos_rule);

// Fill the complete ie buffer
int amf_smf_session_api_fill_qos_ie_info(std::shared_ptr<smf_context_t> smf_ctx,
                                         bstring* authorized_qosrules,
                                         bstring* qos_flow_descriptors);

// Default qos to be filled in PDU session creation
void amf_smf_session_set_default_qos_info(
    std::shared_ptr<smf_context_t> smf_ctx);

// Get the qos flow map entry by flow identifier
std::shared_ptr<qos_flow_setup_request_item>
amf_smf_get_qos_flow_map_entry_by_flow_identifier(
    std::shared_ptr<smf_context_t> smf_context, uint32_t qos_flow_identifier);

// Insert flow map entry into the smf_context map of qos
std::shared_ptr<qos_flow_setup_request_item>
amf_smf_get_qos_flow_map_entry_insert(
    std::shared_ptr<smf_context_t> smf_context, uint32_t qos_flow_identifier);

// Remove the flow map entry
void amf_smf_qos_flow_map_entry_remove(
    std::shared_ptr<smf_context_t> smf_context, uint32_t qos_flow_identifier);
}  // namespace magma5g
