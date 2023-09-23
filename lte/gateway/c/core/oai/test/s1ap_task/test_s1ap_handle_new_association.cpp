/**
 * Copyright 2021 The Magma Authors.
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
#include <gtest/gtest.h>

extern "C" {
#include "lte/gateway/c/core/oai/common/log.h"
#include "S1ap_S1AP-PDU.h"
}

#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_mme_handlers.hpp"
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_state_manager.hpp"

using ::testing::Test;

namespace magma {
namespace lte {

using oai::EnbDescription;
TEST(test_s1ap_handle_new_association, empty_initial_state) {
  oai::S1apState* s = create_s1ap_state();
  // 192.168.60.141 as network bytes
  bstring ran_cp_ipaddr = bfromcstr("\xc0\xa8\x3c\x8d");
  sctp_new_peer_t p = {
      .instreams = 1,
      .outstreams = 2,
      .assoc_id = 3,
      .ran_cp_ipaddr = ran_cp_ipaddr,
  };
  EXPECT_EQ(s1ap_handle_new_association(s, &p), RETURNok);

  EXPECT_EQ(s->enbs_size(), 1);

  EnbDescription enbd;
  proto_map_uint32_enb_description_t enb_map;
  enb_map.map = s->mutable_enbs();
  EXPECT_EQ(enb_map.get(p.assoc_id, &enbd), magma::PROTO_MAP_OK);
  EXPECT_EQ(enbd.sctp_assoc_id(), 3);
  EXPECT_EQ(enbd.instreams(), 1);
  EXPECT_EQ(enbd.outstreams(), 2);
  EXPECT_EQ(enbd.enb_id(), 0xFFFFFFFF);
  EXPECT_EQ(enbd.s1_enb_state(), oai::S1AP_INIT);
  EXPECT_EQ(enbd.next_sctp_stream(), 1);
  EXPECT_STREQ(enbd.ran_cp_ipaddr().c_str(),
               "\300\250<\215\0\0\0\0\0\0\0\0\0\0\0\0");
  EXPECT_EQ(enbd.ran_cp_ipaddr_sz(), 4);

  // association is created, but S1Setup has not yet occurred
  EXPECT_EQ(s->num_enbs(), 0);

  bdestroy(ran_cp_ipaddr);
  free_s1ap_state(s);
}

TEST(test_s1ap_handle_new_association, shutdown) {
  oai::S1apState* s = create_s1ap_state();
  sctp_new_peer_t p = {.assoc_id = 1};
  EXPECT_EQ(s1ap_handle_new_association(s, &p), RETURNok);

  // set enb to shutdown state
  EnbDescription enbd;
  proto_map_uint32_enb_description_t enb_map;
  enb_map.map = s->mutable_enbs();
  EXPECT_EQ(enb_map.get(p.assoc_id, &enbd), magma::PROTO_MAP_OK);
  enbd.set_s1_enb_state(oai::S1AP_SHUTDOWN);
  s1ap_state_update_enb_map(s, p.assoc_id, &enbd);

  // expect error
  EXPECT_EQ(s1ap_handle_new_association(s, &p), RETURNerror);

  free_s1ap_state(s);
}

TEST(test_s1ap_handle_new_association, resetting) {
  oai::S1apState* s = create_s1ap_state();
  sctp_new_peer_t p = {.assoc_id = 1};
  EXPECT_EQ(s1ap_handle_new_association(s, &p), RETURNok);

  // set enb to shutdown state
  EnbDescription enbd;
  proto_map_uint32_enb_description_t enb_map;
  enb_map.map = s->mutable_enbs();
  EXPECT_EQ(enb_map.get(p.assoc_id, &enbd), magma::PROTO_MAP_OK);
  enbd.set_s1_enb_state(oai::S1AP_RESETING);
  s1ap_state_update_enb_map(s, p.assoc_id, &enbd);

  // expect error
  EXPECT_EQ(s1ap_handle_new_association(s, &p), RETURNerror);

  free_s1ap_state(s);
}

TEST(test_s1ap_handle_new_association, reassociate) {
  oai::S1apState* s = create_s1ap_state();
  sctp_new_peer_t p = {.assoc_id = 1};
  EXPECT_EQ(s1ap_handle_new_association(s, &p), RETURNok);

  // make sure first association worked
  EnbDescription enbd;
  proto_map_uint32_enb_description_t enb_map;
  enb_map.map = s->mutable_enbs();
  EXPECT_EQ(enb_map.get(p.assoc_id, &enbd), magma::PROTO_MAP_OK);
  EXPECT_EQ(enbd.sctp_assoc_id(), 1);
  EXPECT_EQ(enbd.instreams(), 0);
  EXPECT_EQ(enbd.outstreams(), 0);
  EXPECT_STREQ(enbd.ran_cp_ipaddr().c_str(), "");
  EXPECT_EQ(enbd.ran_cp_ipaddr_sz(), 0);
  // should be OK if enb status is READY
  enbd.set_s1_enb_state(oai::S1AP_READY);
  s1ap_state_update_enb_map(s, p.assoc_id, &enbd);

  // new assoc with same id should overwrite
  bstring ran_cp_ipaddr = bfromcstr("\xc0\xa8\x3c\x8d");
  sctp_new_peer_t p2 = {
      .instreams = 10,
      .outstreams = 20,
      .assoc_id = 1,
      .ran_cp_ipaddr = ran_cp_ipaddr,
  };
  EXPECT_EQ(s1ap_handle_new_association(s, &p2), RETURNok);
  EXPECT_EQ(enb_map.get(p2.assoc_id, &enbd), magma::PROTO_MAP_OK);

  EXPECT_EQ(enbd.sctp_assoc_id(), 1);
  EXPECT_EQ(enbd.instreams(), 10);
  EXPECT_EQ(enbd.outstreams(), 20);
  EXPECT_STREQ(enbd.ran_cp_ipaddr().c_str(), "\300\250<\215");
  EXPECT_EQ(enbd.ran_cp_ipaddr_sz(), 4);
  EXPECT_EQ(enbd.s1_enb_state(), oai::S1AP_INIT);

  bdestroy(ran_cp_ipaddr);
  free_s1ap_state(s);
}

TEST(test_s1ap_handle_new_association, clean_stale_association) {
  oai::S1apState* s = create_s1ap_state();
  // 192.168.60.141 as network bytes
  bstring ran_cp_ipaddr = bfromcstr("\xc0\xa8\x3c\x8d");
  sctp_new_peer_t p = {
      .instreams = 1,
      .outstreams = 2,
      .assoc_id = 3,
      .ran_cp_ipaddr = ran_cp_ipaddr,
  };
  EXPECT_EQ(s1ap_handle_new_association(s, &p), RETURNok);

  EXPECT_EQ(s->enbs_size(), 1);

  EnbDescription* enb_ref = new EnbDescription();

  EnbDescription enb_associated;
  proto_map_uint32_enb_description_t enb_map;
  enb_map.map = s->mutable_enbs();
  enb_map.get(p.assoc_id, &enb_associated);

  enb_ref->set_enb_id(enb_associated.enb_id());
  clean_stale_enb_state(s, enb_ref);
  EXPECT_EQ(s->enbs_size(), 0);

  bdestroy(ran_cp_ipaddr);
  free_enb_description(reinterpret_cast<void**>(&enb_ref));
  free_s1ap_state(s);
}

}  // namespace lte
}  // namespace magma
