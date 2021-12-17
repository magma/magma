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

#include "lte/gateway/c/core/oai/tasks/amf/amf_app_state_converter.h"
#include <vector>
#include <memory>
extern "C" {
#include "lte/gateway/c/core/oai/lib/message_utils/bytes_to_ie.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/lib/message_utils/ie_to_bytes.h"
#include "lte/gateway/c/core/oai/common/log.h"
}

using magma::lte::oai::MmeNasState;
namespace magma5g {

AmfNasStateConverter::AmfNasStateConverter()  = default;
AmfNasStateConverter::~AmfNasStateConverter() = default;

// HelperFunction: Converts guti_m5_t to std::string
std::string AmfNasStateConverter::amf_app_convert_guti_m5_to_string(
    const guti_m5_t& guti) {
#define GUTI_STRING_LEN 25
  char* temp_str =
      reinterpret_cast<char*>(calloc(1, sizeof(char) * GUTI_STRING_LEN));
  snprintf(
      temp_str, GUTI_STRING_LEN, "%x%x%x%x%x%x%02x%04x%04x%08x",
      guti.guamfi.plmn.mcc_digit1, guti.guamfi.plmn.mcc_digit2,
      guti.guamfi.plmn.mcc_digit3, guti.guamfi.plmn.mnc_digit1,
      guti.guamfi.plmn.mnc_digit2, guti.guamfi.plmn.mnc_digit3,
      guti.guamfi.amf_regionid, guti.guamfi.amf_set_id, guti.guamfi.amf_pointer,
      guti.m_tmsi);
  std::string guti_str(temp_str);
  free(temp_str);
  return guti_str;
}

// HelperFunction: Converts std:: string back to guti_m5_t
void AmfNasStateConverter::amf_app_convert_string_to_guti_m5(
    const std::string& guti_str, guti_m5_t* guti_m5_p) {
  int idx                   = 0;
  std::size_t chars_to_read = 1;
#define HEX_BASE_VAL 16
  guti_m5_p->guamfi.plmn.mcc_digit1 = std::stoul(
      guti_str.substr(idx++, chars_to_read), &chars_to_read, HEX_BASE_VAL);
  guti_m5_p->guamfi.plmn.mcc_digit2 = std::stoul(
      guti_str.substr(idx++, chars_to_read), &chars_to_read, HEX_BASE_VAL);
  guti_m5_p->guamfi.plmn.mcc_digit3 = std::stoul(
      guti_str.substr(idx++, chars_to_read), &chars_to_read, HEX_BASE_VAL);
  guti_m5_p->guamfi.plmn.mnc_digit1 = std::stoul(
      guti_str.substr(idx++, chars_to_read), &chars_to_read, HEX_BASE_VAL);
  guti_m5_p->guamfi.plmn.mnc_digit2 = std::stoul(
      guti_str.substr(idx++, chars_to_read), &chars_to_read, HEX_BASE_VAL);
  guti_m5_p->guamfi.plmn.mnc_digit3 = std::stoul(
      guti_str.substr(idx++, chars_to_read), &chars_to_read, HEX_BASE_VAL);
  chars_to_read                  = 2;
  guti_m5_p->guamfi.amf_regionid = std::stoul(
      guti_str.substr(idx, chars_to_read), &chars_to_read, HEX_BASE_VAL);
  idx += chars_to_read;
  chars_to_read                = 4;
  guti_m5_p->guamfi.amf_set_id = std::stoul(
      guti_str.substr(idx, chars_to_read), &chars_to_read, HEX_BASE_VAL);
  idx += chars_to_read;
  chars_to_read                 = 4;
  guti_m5_p->guamfi.amf_pointer = std::stoul(
      guti_str.substr(idx, chars_to_read), &chars_to_read, HEX_BASE_VAL);
  idx += chars_to_read;
  chars_to_read     = 8;
  guti_m5_p->m_tmsi = std::stoul(
      guti_str.substr(idx, chars_to_read), &chars_to_read, HEX_BASE_VAL);
}

/*********************************************************
 *                AMF app state<-> Proto                  *
 * Functions to serialize/desearialize AMF app state      *
 * The caller is responsible for all memory management    *
 **********************************************************/

void AmfNasStateConverter::state_to_proto(
    const amf_app_desc_t* amf_nas_state_p, MmeNasState* state_proto) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  state_proto->set_mme_app_ue_s1ap_id_generator(
      amf_nas_state_p->amf_app_ue_ngap_id_generator);

  // These Functions are to be removed as part of the stateless enhancement
  // maps to proto
  auto amf_ue_ctxts_proto = state_proto->mutable_mme_ue_contexts();
  OAILOG_DEBUG(LOG_AMF_APP, "IMSI table to proto");
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

void AmfNasStateConverter::proto_to_state(
    const MmeNasState& state_proto, amf_app_desc_t* amf_nas_state_p) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  amf_nas_state_p->amf_app_ue_ngap_id_generator =
      state_proto.mme_app_ue_s1ap_id_generator();

  if (amf_nas_state_p->amf_app_ue_ngap_id_generator == 0) {  // uninitialized
    amf_nas_state_p->amf_app_ue_ngap_id_generator = 1;
  }
  OAILOG_INFO(LOG_AMF_APP, "Done reading AMF statistics from data store");

  magma::lte::oai::MmeUeContext amf_ue_ctxts_proto =
      state_proto.mme_ue_contexts();

  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

void AmfNasStateConverter::ue_to_proto(
    const ue_m5gmm_context_t* ue_ctxt,
    magma::lte::oai::UeContext* ue_ctxt_proto) {
  ue_m5gmm_context_to_proto(ue_ctxt, ue_ctxt_proto);
}

void AmfNasStateConverter::proto_to_ue(
    const magma::lte::oai::UeContext& ue_ctxt_proto,
    ue_m5gmm_context_t* ue_ctxt) {
  proto_to_ue_m5gmm_context(ue_ctxt_proto, ue_ctxt);
}

/*********************************************************
 *                UE Context <-> Proto                    *
 * Functions to serialize/desearialize UE context         *
 * The caller needs to acquire a lock on UE context       *
 **********************************************************/

void AmfNasStateConverter::ue_m5gmm_context_to_proto(
    const ue_m5gmm_context_t* state_ue_m5gmm_context,
    magma::lte::oai::UeContext* ue_context_proto) {
  // Actual implementation logic will be added as part of upcoming pr
}

void AmfNasStateConverter::proto_to_ue_m5gmm_context(
    const magma::lte::oai::UeContext& ue_context_proto,
    ue_m5gmm_context_t* state_ue_m5gmm_context) {
  // Actual implementation logic will be added as part of upcoming pr
}

}  // namespace magma5g
