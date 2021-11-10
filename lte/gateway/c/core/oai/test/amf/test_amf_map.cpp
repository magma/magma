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
#include "lte/gateway/c/core/oai/include/map.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_defs.h"

using ::testing::Test;

namespace magma5g {
TEST(test_map, test_map) {
  // Initializations for Map: Key-uint64_t Data-uint64_t
  uint64_t data;
  amf_app_desc_t* state_cache_p = new (amf_app_desc_t);

  // Initializations for Map: Key-guti_m5_t Data-uint64_t
  guti_m5_t guti_1;
  guti_1.m_tmsi = 0x2bfb815f;
  guti_m5_t guti_2;
  guti_2.m_tmsi = 0x1afc831f;
  uint64_t gutiData;

  state_cache_p->amf_ue_contexts.imsi_amf_ue_id_htbl.set_name("IMSI HASHTABLE");
  EXPECT_EQ(
      state_cache_p->amf_ue_contexts.imsi_amf_ue_id_htbl.get_name(),
      "IMSI HASHTABLE");

  // Trying to get from an empty map
  EXPECT_EQ(
      state_cache_p->amf_ue_contexts.imsi_amf_ue_id_htbl.get(2, &data),
      MAP_EMPTY);

  // Inserting new <key,value> pair
  EXPECT_EQ(
      state_cache_p->amf_ue_contexts.imsi_amf_ue_id_htbl.insert(1, 10), MAP_OK);
  EXPECT_EQ(
      state_cache_p->amf_ue_contexts.imsi_amf_ue_id_htbl.insert(2, 20), MAP_OK);

  // Inserting already existing key.Expected: failure.
  EXPECT_EQ(
      state_cache_p->amf_ue_contexts.imsi_amf_ue_id_htbl.insert(1, 20),
      MAP_KEY_ALREADY_EXISTS);

  // //Getting data from map
  EXPECT_EQ(
      state_cache_p->amf_ue_contexts.imsi_amf_ue_id_htbl.get(1, &data), MAP_OK);
  EXPECT_EQ(data, 10);

  // Getting with invalid key
  EXPECT_EQ(
      state_cache_p->amf_ue_contexts.imsi_amf_ue_id_htbl.get(5, &data),
      MAP_KEY_NOT_EXISTS);

  // Removing entry from table
  EXPECT_EQ(
      state_cache_p->amf_ue_contexts.imsi_amf_ue_id_htbl.remove(1), MAP_OK);

  // Trying to remove from invalid key
  EXPECT_EQ(
      state_cache_p->amf_ue_contexts.imsi_amf_ue_id_htbl.remove(5),
      MAP_KEY_NOT_EXISTS);

  // Object table
  EXPECT_EQ(
      state_cache_p->amf_ue_contexts.guti_ue_context_htbl.insert(guti_1, 100),
      MAP_OK);
  EXPECT_EQ(
      state_cache_p->amf_ue_contexts.guti_ue_context_htbl.insert(guti_1, 400),
      MAP_KEY_ALREADY_EXISTS);

  EXPECT_EQ(
      state_cache_p->amf_ue_contexts.guti_ue_context_htbl.insert(guti_2, 200),
      MAP_OK);

  EXPECT_EQ(
      state_cache_p->amf_ue_contexts.guti_ue_context_htbl.get(
          guti_1, &gutiData),
      MAP_OK);
  EXPECT_EQ(gutiData, 100);

  EXPECT_EQ(
      state_cache_p->amf_ue_contexts.guti_ue_context_htbl.remove(guti_1),
      MAP_OK);
  EXPECT_EQ(
      state_cache_p->amf_ue_contexts.guti_ue_context_htbl.get(
          guti_1, &gutiData),
      MAP_KEY_NOT_EXISTS);

  delete state_cache_p;
}
}  // namespace magma5g
