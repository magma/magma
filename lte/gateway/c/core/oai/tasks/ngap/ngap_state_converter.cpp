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
/****************************************************************************
  Source      ngap_state_converter.cpp
  Date        2020/07/28
  Author      Ashish Prajapati
  Subsystem   Access and Mobility Management Function
  Description Defines NG Application Protocol Messages

*****************************************************************************/

#include "ngap_state_converter.h"
using namespace std;
using namespace magma::lte;
using namespace magma::lte::oai;

using magma::lte::oai::Ngap_UeDescription;
using magma::lte::oai::NgapState;

namespace magma5g {

NgapStateConverter::~NgapStateConverter() = default;
NgapStateConverter::NgapStateConverter()  = default;

void NgapStateConverter::state_to_proto(ngap_state_t* state, NgapState* proto) {
  proto->Clear();

  // copy over gnbs
  hashtable_ts_to_proto<gnb_description_t, GnbDescription>(
      &state->gnbs, proto->mutable_gnbs(), gnb_to_proto, LOG_NGAP);

  // copy over amfid2associd
  hashtable_rc_t ht_rc;
  amf_ue_ngap_id_t amfid;
  // Helper ptr so sctp_assoc_id can be casted from double ptr on
  // hashtable_ts_get
  void* sctp_id_ptr  = nullptr;
  auto amfid2associd = proto->mutable_amfid2associd();

  hashtable_key_array_t* keys = hashtable_ts_get_keys(&state->amfid2associd);
  if (!keys) {
    OAILOG_DEBUG(LOG_NGAP, "No keys in amfid2associd hashtable");
  } else {
    for (int i = 0; i < keys->num_keys; i++) {
      amfid = (amf_ue_ngap_id_t) keys->keys[i];
      ht_rc = hashtable_ts_get(
          &state->amfid2associd, (hash_key_t) amfid, (void**) &sctp_id_ptr);
      AssertFatal(ht_rc == HASH_TABLE_OK, "amfid not in amfid2associd");
      if (sctp_id_ptr) {
        sctp_assoc_id_t sctp_assoc_id =
            (sctp_assoc_id_t)(uintptr_t) sctp_id_ptr;
        (*amfid2associd)[amfid] = sctp_assoc_id;
      }
    }
    FREE_HASHTABLE_KEY_ARRAY(keys);
  }

  proto->set_num_gnbs(state->num_gnbs);
}

void NgapStateConverter::proto_to_state(
    const NgapState& proto, ngap_state_t* state) {
  proto_to_hashtable_ts<GnbDescription, gnb_description_t>(
      proto.gnbs(), &state->gnbs, proto_to_gnb, LOG_NGAP);

  hashtable_rc_t ht_rc;
  auto amfid2associd = proto.amfid2associd();
  for (auto const& kv : amfid2associd) {
    amf_ue_ngap_id_t amfid  = (amf_ue_ngap_id_t) kv.first;
    sctp_assoc_id_t associd = (sctp_assoc_id_t) kv.second;

    ht_rc = hashtable_ts_insert(
        &state->amfid2associd, (hash_key_t) amfid, (void*) (uintptr_t) associd);
    AssertFatal(ht_rc == HASH_TABLE_OK, "failed to insert associd");
  }

  state->num_gnbs = proto.num_gnbs();
}

void NgapStateConverter::gnb_to_proto(
    gnb_description_t* gnb, oai::GnbDescription* proto) {
  proto->Clear();

  proto->set_gnb_id(gnb->gnb_id);
  proto->set_ng_state(gnb->ng_state);
  proto->set_gnb_name(gnb->gnb_name);
  proto->set_default_paging_drx(gnb->default_paging_drx);
  proto->set_nb_ue_associated(gnb->nb_ue_associated);
  proto->set_sctp_assoc_id(gnb->sctp_assoc_id);
  proto->set_next_sctp_stream(gnb->next_sctp_stream);
  proto->set_instreams(gnb->instreams);
  proto->set_outstreams(gnb->outstreams);

  // store ue_ids
  hashtable_uint64_ts_to_proto(&gnb->ue_id_coll, proto->mutable_ue_ids());
}

void NgapStateConverter::proto_to_gnb(
    const oai::GnbDescription& proto, gnb_description_t* gnb) {
  memset(gnb, 0, sizeof(*gnb));

  gnb->gnb_id   = proto.gnb_id();
  gnb->ng_state = (amf_ng_gnb_state_s) proto.ng_state();
  strncpy(gnb->gnb_name, proto.gnb_name().c_str(), sizeof(gnb->gnb_name));
  gnb->default_paging_drx = proto.default_paging_drx();
  gnb->nb_ue_associated   = proto.nb_ue_associated();
  gnb->sctp_assoc_id      = proto.sctp_assoc_id();
  gnb->next_sctp_stream   = proto.next_sctp_stream();
  gnb->instreams          = proto.instreams();
  gnb->outstreams         = proto.outstreams();

  // load ues
  hashtable_rc_t ht_rc;
  auto ht_name = bfromcstr("ngap_ue_coll");

  hashtable_uint64_ts_init(
      &gnb->ue_id_coll, amf_config.max_ues, nullptr, ht_name);
  bdestroy(ht_name);

  auto ue_ids = proto.ue_ids();
  for (auto const& kv : ue_ids) {
    amf_ue_ngap_id_t amf_ue_ngap_id = kv.first;
    uint64_t comp_ngap_id           = kv.second;

    ht_rc = hashtable_uint64_ts_insert(
        &gnb->ue_id_coll, (hash_key_t) amf_ue_ngap_id, comp_ngap_id);
    if (ht_rc != HASH_TABLE_OK) {
      OAILOG_DEBUG(
          LOG_NGAP, "Failed to insert amf_ue_ngap_id in ue_coll_id hashtable");
    }
  }
}
void NgapStateConverter::ue_to_proto(
    const m5g_ue_description_t* ue, oai::Ngap_UeDescription* proto) {
  proto->Clear();

  proto->set_ng_ue_state(ue->ng_ue_state);
  proto->set_gnb_ue_ngap_id(ue->gnb_ue_ngap_id);
  proto->set_amf_ue_ngap_id(ue->amf_ue_ngap_id);
  proto->set_sctp_assoc_id(ue->sctp_assoc_id);
  proto->set_sctp_stream_recv(ue->sctp_stream_recv);
  proto->set_sctp_stream_send(ue->sctp_stream_send);
  proto->mutable_ngap_ue_context_rel_timer()->set_id(
      ue->ngap_ue_context_rel_timer.id);
  proto->mutable_ngap_ue_context_rel_timer()->set_sec(
      ue->ngap_ue_context_rel_timer.sec);
}
void NgapStateConverter::proto_to_ue(
    const oai::Ngap_UeDescription& proto, m5g_ue_description_t* ue) {
  memset(ue, 0, sizeof(*ue));

  ue->ng_ue_state                   = (ng_ue_state_s) proto.ng_ue_state();
  ue->gnb_ue_ngap_id                = proto.gnb_ue_ngap_id();
  ue->amf_ue_ngap_id                = proto.amf_ue_ngap_id();
  ue->sctp_assoc_id                 = proto.sctp_assoc_id();
  ue->sctp_stream_recv              = proto.sctp_stream_recv();
  ue->sctp_stream_send              = proto.sctp_stream_send();
  ue->ngap_ue_context_rel_timer.id  = proto.ngap_ue_context_rel_timer().id();
  ue->ngap_ue_context_rel_timer.sec = proto.ngap_ue_context_rel_timer().sec();
}

void NgapStateConverter::ngap_imsi_map_to_proto(
    const ngap_imsi_map_t* ngap_imsi_map, oai::NgapImsiMap* ngap_imsi_proto) {
  hashtable_uint64_ts_to_proto(
      ngap_imsi_map->amf_ue_id_imsi_htbl,
      ngap_imsi_proto->mutable_amf_ue_id_imsi_map());
}
void NgapStateConverter::proto_to_ngap_imsi_map(
    const oai::NgapImsiMap& ngap_imsi_proto, ngap_imsi_map_t* ngap_imsi_map) {
  proto_to_hashtable_uint64_ts(
      ngap_imsi_proto.amf_ue_id_imsi_map(), ngap_imsi_map->amf_ue_id_imsi_htbl);
}

}  // namespace magma5g
