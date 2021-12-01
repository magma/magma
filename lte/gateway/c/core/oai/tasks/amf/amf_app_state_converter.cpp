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

namespace magma5g {

AmfNasStateConverter::AmfNasStateConverter()  = default;
AmfNasStateConverter::~AmfNasStateConverter() = default;

// Converts Map<uint64_t,uint64_t> to Proto (replaces hastable to proto)
void AmfNasStateConverter::map_uint64_uint64_to_proto(
    map_uint64_uint64_t map,
    google::protobuf::Map<uint64_t, uint64_t>* proto_map) {
  for (auto& elm : map.umap) {
    (*proto_map)[elm.first] = elm.second;
  }
}

// HelperFunction: Converts guti_m5_t to std::string
// TODO: Implement with C++ equivalent of snprintf, get rid of calloc
std::string AmfNasStateConverter::amf_app_convert_guti_m5_to_string(
    guti_m5_t guti) {
#define GUTI_STRING_LEN 25
  char* str = reinterpret_cast<char*>(calloc(1, sizeof(char) * GUTI_STRING_LEN));
  snprintf(
      str, GUTI_STRING_LEN, "%x%x%x%x%x%x%02x%04x%04x%08x",
      guti.guamfi.plmn.mcc_digit1, guti.guamfi.plmn.mcc_digit2,
      guti.guamfi.plmn.mcc_digit3, guti.guamfi.plmn.mnc_digit1,
      guti.guamfi.plmn.mnc_digit2, guti.guamfi.plmn.mnc_digit3,
      guti.guamfi.amf_regionid, guti.guamfi.amf_set_id, guti.guamfi.amf_pointer,
      guti.m_tmsi);
  std::string guti_str(str);
  free(str);
  return (guti_str);
}

// HelperFunction: Converts std:: string back to guti_m5_t
void AmfNasStateConverter::amf_app_convert_string_to_guti_m5(
    guti_m5_t* guti_m5_p, const std::string& guti_str) {
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
  // TODO: Logs
  // OAILOG_DEBUG_GUTI(guti_m5_p);
}

// Converts Map<guti_m5_t,uint64_t> to proto
void AmfNasStateConverter::map_guti_uint64_to_proto(
    map_guti_m5_uint64_t guti_map,
    google::protobuf::Map<std::string, uint64_t>* proto_map) {
  std::string guti_str;
  for (auto& elm : guti_map.umap) {
    guti_str               = amf_app_convert_guti_m5_to_string(elm.first);
    (*proto_map)[guti_str] = elm.second;
  }
}

// Converts Proto to Map<uint64_t,uint64_t>
void AmfNasStateConverter::proto_to_map_uint64_uint64(
    const google::protobuf::Map<uint64_t, uint64_t>& proto_map,
    map_uint64_uint64_t* map) {
  for (auto const& kv : proto_map) {
    uint64_t id          = kv.first;
    uint64_t val         = kv.second;
    magma::map_rc_t m_rc = map->insert(kv.first, kv.second);

    if (m_rc != magma::MAP_OK) {
      OAILOG_ERROR(
          LOG_UTIL, "Failed to insert value %lu in table %s: error: %s\n", val,
          map->name.c_str(), map_rc_code2string(m_rc).c_str());
    }
  }
}

// Converts Proto to Map<guti_m5_t,uint64_t> [Needs recheck]
void AmfNasStateConverter::proto_to_guti_map(
    const google::protobuf::Map<std::string, uint64_t>& proto_map,
    map_guti_m5_uint64_t* guti_map) {
  for (auto const& kv : proto_map) {
    amf_ue_ngap_id_t amf_ue_ngap_id = kv.second;
    // May remove unique pointer
    std::unique_ptr<guti_m5_t> guti = std::make_unique<guti_m5_t>();
    memset(guti.get(), 0, sizeof(guti_m5_t));
    // Converts guti to string.
    amf_app_convert_string_to_guti_m5(guti.get(), kv.first);

    guti_m5_t guti_received = *guti.get();
    magma::map_rc_t m_rc    = guti_map->insert(guti_received, amf_ue_ngap_id);
    if (m_rc != magma::MAP_OK) {
      OAILOG_ERROR(
          LOG_AMF_APP,
          "Failed to insert amf_ue_ngap_id %lu in GUTI table, error: %s\n",
          amf_ue_ngap_id, map_rc_code2string(m_rc).c_str());
    }
  }
}

// /*********************************************************
//  *                AMF app state<-> Proto                  *
//  * Functions to serialize/desearialize AMF app state      *
//  * The caller is responsible for all memory management    *
//  **********************************************************/

void AmfNasStateConverter::state_to_proto(
    const amf_app_desc_t* amf_nas_state_p,
    magma::lte::oai::MmeNasState* state_proto) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  state_proto->set_mme_app_ue_s1ap_id_generator(
      amf_nas_state_p->amf_app_ue_ngap_id_generator);
  // COMMENTED NOW:
  // state_proto->set_statistic_timer_id(amf_nas_state_p->m5_statistic_timer_id);

  // maps to proto
  auto amf_ue_ctxts_proto = state_proto->mutable_mme_ue_contexts();
  OAILOG_DEBUG(LOG_AMF_APP, "IMSI table to proto");
  map_uint64_uint64_to_proto(
      amf_nas_state_p->amf_ue_contexts.imsi_amf_ue_id_htbl,
      amf_ue_ctxts_proto->mutable_imsi_ue_id_htbl());
  map_uint64_uint64_to_proto(
      amf_nas_state_p->amf_ue_contexts.tun11_ue_context_htbl,
      amf_ue_ctxts_proto->mutable_tun11_ue_id_htbl());
  map_uint64_uint64_to_proto(
      amf_nas_state_p->amf_ue_contexts.gnb_ue_ngap_id_ue_context_htbl,
      amf_ue_ctxts_proto->mutable_enb_ue_id_ue_id_htbl());
  map_guti_uint64_to_proto(
      amf_nas_state_p->amf_ue_contexts.guti_ue_context_htbl,
      amf_ue_ctxts_proto->mutable_guti_ue_id_htbl());
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

void AmfNasStateConverter::proto_to_state(
    const magma::lte::oai::MmeNasState& state_proto,
    amf_app_desc_t* amf_nas_state_p) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  amf_nas_state_p->amf_app_ue_ngap_id_generator =
      state_proto.mme_app_ue_s1ap_id_generator();
  // COMMENTED NOW: amf_nas_state_p->m5_statistic_timer_id =
  // state_proto.statistic_timer_id();

  if (amf_nas_state_p->amf_app_ue_ngap_id_generator == 0) {  // uninitialized
    amf_nas_state_p->amf_app_ue_ngap_id_generator = 1;
  }
  OAILOG_INFO(LOG_AMF_APP, "Done reading AMF statistics from data store");

  // copy mme_ue_contexts
  magma::lte::oai::MmeUeContext amf_ue_ctxts_proto =
      state_proto.mme_ue_contexts();

  amf_ue_context_t* amf_ue_ctxt_state = &amf_nas_state_p->amf_ue_contexts;

  // proto to maps
  OAILOG_INFO(LOG_AMF_APP, "Hashtable AMF UE ID => IMSI");
  proto_to_map_uint64_uint64(
      amf_ue_ctxts_proto.imsi_ue_id_htbl(),
      &amf_ue_ctxt_state->imsi_amf_ue_id_htbl);
  proto_to_map_uint64_uint64(
      amf_ue_ctxts_proto.tun11_ue_id_htbl(),
      &amf_ue_ctxt_state->tun11_ue_context_htbl);
  proto_to_map_uint64_uint64(
      amf_ue_ctxts_proto.enb_ue_id_ue_id_htbl(),
      &amf_ue_ctxt_state->gnb_ue_ngap_id_ue_context_htbl);

  proto_to_guti_map(
      amf_ue_ctxts_proto.guti_ue_id_htbl(),
      &amf_ue_ctxt_state->guti_ue_context_htbl);
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

void AmfNasStateConverter::ue_to_proto(
    const ue_m5gmm_context_t* ue_ctxt,
    magma::lte::oai::UeContext* ue_ctxt_proto) {
  // ue_context_to_proto(ue_ctxt, ue_ctxt_proto);
}

void AmfNasStateConverter::proto_to_ue(
    const magma::lte::oai::UeContext& ue_ctxt_proto,
    ue_m5gmm_context_t* ue_ctxt) {
  // proto_to_ue_mm_context(ue_ctxt_proto, ue_ctxt);
}
}  // namespace magma5g
