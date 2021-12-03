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

#include <gtest/gtest.h>
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_defs.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_state_converter.h"
#include "lte/gateway/c/core/oai/include/map.h"
using ::testing::Test;

namespace magma5g {
// Note this UT Has minimal checks and uses hard coded values
// TODO: Minimize hardcoding in UT
TEST(test_guti_to_string, test_guti_to_string) {
  guti_m5_t guti1, guti2;
  guti1.guamfi.plmn.mcc_digit1 = 2;
  guti1.guamfi.plmn.mcc_digit2 = 2;
  guti1.guamfi.plmn.mcc_digit3 = 2;
  guti1.guamfi.plmn.mnc_digit1 = 4;
  guti1.guamfi.plmn.mnc_digit2 = 5;
  guti1.guamfi.plmn.mnc_digit3 = 6;
  guti1.guamfi.amf_regionid    = 1;
  guti1.guamfi.amf_set_id      = 1;
  guti1.guamfi.amf_pointer     = 0;
  guti1.m_tmsi                 = 0X212e5025;

  std::string guti1_str =
      AmfNasStateConverter::amf_app_convert_guti_m5_to_string(guti1);

  AmfNasStateConverter::amf_app_convert_string_to_guti_m5(&guti2, guti1_str);

  EXPECT_EQ(guti1.guamfi.plmn.mcc_digit1, guti2.guamfi.plmn.mcc_digit1);
  EXPECT_EQ(guti1.guamfi.plmn.mcc_digit2, guti2.guamfi.plmn.mcc_digit2);
  EXPECT_EQ(guti1.guamfi.plmn.mcc_digit3, guti2.guamfi.plmn.mcc_digit3);
  EXPECT_EQ(guti1.guamfi.plmn.mnc_digit1, guti2.guamfi.plmn.mnc_digit1);
  EXPECT_EQ(guti1.guamfi.plmn.mnc_digit2, guti2.guamfi.plmn.mnc_digit2);
  EXPECT_EQ(guti1.guamfi.plmn.mnc_digit3, guti2.guamfi.plmn.mnc_digit3);
  EXPECT_EQ(guti1.guamfi.amf_regionid, guti2.guamfi.amf_regionid);
  EXPECT_EQ(guti1.guamfi.amf_set_id, guti2.guamfi.amf_set_id);
  EXPECT_EQ(guti1.guamfi.amf_pointer, guti2.guamfi.amf_pointer);
  EXPECT_EQ(guti1.m_tmsi, guti2.m_tmsi);
}

// Note this UT Has minimal checks and uses hard coded values
TEST(test_state_to_proto, test_state_to_proto) {
  // Guti setup
  guti_m5_t guti1;
  memset(&guti1, 0, sizeof(guti1));

  guti1.guamfi.plmn.mcc_digit1 = 2;
  guti1.guamfi.plmn.mcc_digit2 = 2;
  guti1.guamfi.plmn.mcc_digit3 = 2;
  guti1.guamfi.plmn.mnc_digit1 = 4;
  guti1.guamfi.plmn.mnc_digit2 = 5;
  guti1.guamfi.plmn.mnc_digit3 = 6;
  guti1.guamfi.amf_regionid    = 1;
  guti1.guamfi.amf_set_id      = 1;
  guti1.guamfi.amf_pointer     = 0;
  guti1.m_tmsi                 = 556683301;

  amf_app_desc_t amf_app_desc1 = {}, amf_app_desc2 = {};
  magma::lte::oai::MmeNasState state_proto = magma::lte::oai::MmeNasState();
  uint64_t data                            = 0;

  amf_app_desc1.amf_app_ue_ngap_id_generator = 0x05;
  amf_app_desc1.amf_ue_contexts.imsi_amf_ue_id_htbl.insert(1, 10);
  amf_app_desc1.amf_ue_contexts.tun11_ue_context_htbl.insert(2, 20);
  amf_app_desc1.amf_ue_contexts.gnb_ue_ngap_id_ue_context_htbl.insert(3, 30);
  amf_app_desc1.amf_ue_contexts.guti_ue_context_htbl.insert(guti1, 40);

  AmfNasStateConverter::state_to_proto(&amf_app_desc1, &state_proto);

  AmfNasStateConverter::proto_to_state(state_proto, &amf_app_desc2);

  EXPECT_EQ(
      amf_app_desc1.amf_app_ue_ngap_id_generator,
      amf_app_desc2.amf_app_ue_ngap_id_generator);

  EXPECT_EQ(
      amf_app_desc2.amf_ue_contexts.imsi_amf_ue_id_htbl.get(1, &data),
      magma::MAP_OK);
  EXPECT_EQ(data, 10);
  data = 0;

  EXPECT_EQ(
      amf_app_desc2.amf_ue_contexts.tun11_ue_context_htbl.get(2, &data),
      magma::MAP_OK);
  EXPECT_EQ(data, 20);
  data = 0;

  EXPECT_EQ(
      amf_app_desc2.amf_ue_contexts.gnb_ue_ngap_id_ue_context_htbl.get(
          3, &data),
      magma::MAP_OK);
  EXPECT_EQ(data, 30);
  data = 0;

  EXPECT_EQ(
      amf_app_desc2.amf_ue_contexts.guti_ue_context_htbl.get(guti1, &data),
      magma::MAP_OK);
  EXPECT_EQ(data, 40);
}
}  // namespace magma5g
