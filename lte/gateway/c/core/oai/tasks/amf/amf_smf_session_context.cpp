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

#include <sstream>
#include "lte/gateway/c/core/oai/tasks/amf/include/amf_smf_session_context.hpp"

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/include/amf_config.hpp"

#ifdef __cplusplus
}
#endif

namespace magma5g {
status_code_e amf_smf_context_ue_aggregate_max_bit_rate_set(
    amf_context_s* amf_ctxt_p, ambr_t subscribed_ue_ambr) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  memcpy(&amf_ctxt_p->subscribed_ue_ambr, &subscribed_ue_ambr, sizeof(ambr_t));

  OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNok);
}

status_code_e amf_smf_context_ue_aggregate_max_bit_rate_get(
    const amf_context_s* amf_ctxt_p, bit_rate_t* subscriber_ambr_dl,
    bit_rate_t* subscriber_ambr_ul) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  if (amf_ctxt_p == NULL) {
    OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
  }

  *subscriber_ambr_dl = amf_ctxt_p->subscribed_ue_ambr.br_dl;
  *subscriber_ambr_ul = amf_ctxt_p->subscribed_ue_ambr.br_ul;

  OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNok);
}

/***************************************************************************
**                                                                        **
** Name:    amf_config_get_default_sst_config()                           **
**                                                                        **
** Description: Get default sst value from amf config                     **
**                                                                        **
**                                                                        **
***************************************************************************/
void amf_config_get_default_slice_config(uint8_t* slice_type,
                                         uint8_t* slice_differentiator) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  amf_config_read_lock(&amf_config);

  /* Validate the input parameter */
  if ((slice_type == NULL) || (slice_differentiator == NULL)) {
    OAILOG_FUNC_OUT(LOG_AMF_APP);
  }

  /* Get the default ST value */
  *slice_type = amf_config.plmn_support_list.plmn_support[0].s_nssai.sst;

  /* Get the default SD value */
  if (amf_config.plmn_support_list.plmn_support[0].s_nssai.sd.v !=
      AMF_S_NSSAI_SD_INVALID_VALUE) {
    INT24_TO_BUFFER(amf_config.plmn_support_list.plmn_support[0].s_nssai.sd.v,
                    slice_differentiator);
  }

  amf_config_unlock(&amf_config);
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

// Function to fill the slice information in pdu session establishment
// accept message
void amf_smf_get_slice_configuration(std::shared_ptr<smf_context_t> smf_ctx,
                                     s_nssai_t* slice_config) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  if (slice_config == NULL) {
    OAILOG_FUNC_OUT(LOG_AMF_APP);
  }

  // Check if there is an requested slice value
  if (smf_ctx->requested_nssai.sst) {
    slice_config->sst = smf_ctx->requested_nssai.sst;

    // Check if slice descriptor is found
    if (smf_ctx->requested_nssai.sd[0]) {
      memcpy(slice_config->sd, smf_ctx->requested_nssai.sd, SD_LENGTH);
    }
  } else {
    // Fill in the default if requested slice information is not found
    amf_config_get_default_slice_config(&(slice_config->sst), slice_config->sd);
  }
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

}  // namespace magma5g
