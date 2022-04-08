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
#include <google/protobuf/map.h>
#include "lte/protos/oai/mme_nas_state.pb.h"
#include "lte/gateway/c/core/oai/include/proto_map.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_23.003.h"

using ::testing::Test;
TEST(test_map, test_map) {
  magma::lte::oai::MmeUeContext mme_ue_context_proto =
      magma::lte::oai::MmeUeContext::default_instance();
  magma::map_uint64_uint64_t imsi_map;
  imsi_map.map = mme_ue_context_proto.mutable_imsi_ue_id_htbl();
  imsi_map.set_name("IMSI HASHTABLE");

  EXPECT_EQ(imsi_map.get_name(), "IMSI HASHTABLE");
  // Trying to get from an empty map
  uint64_t data;
  EXPECT_EQ(imsi_map.get(1, &data), magma::MAP_EMPTY);

  // Inserting new <key,value> pair
  EXPECT_EQ(imsi_map.insert(1, 10), magma::MAP_OK);
  EXPECT_EQ(imsi_map.insert(2, 20), magma::MAP_OK);

  // Inserting already existing key.Expected: failure.
  EXPECT_EQ(imsi_map.insert(1, 10), magma::MAP_KEY_ALREADY_EXISTS);

  // Getting data from map
  EXPECT_EQ(imsi_map.get(1, &data), magma::MAP_OK);
  EXPECT_EQ(data, 10);

  // Getting with invalid key
  EXPECT_EQ(imsi_map.get(7, &data), magma::MAP_KEY_NOT_EXISTS);

  // Removing entry from table
  EXPECT_EQ(imsi_map.remove(1), magma::MAP_OK);

  // Trying to remove from invalid key
  EXPECT_EQ(imsi_map.remove(9), magma::MAP_KEY_NOT_EXISTS);

  // Test for clear() and isEmpty()
  imsi_map.clear();
  EXPECT_TRUE(imsi_map.isEmpty());

  // Object table
  std::string guti_1, guti_2, guti_3;
  // GUTI: mme-group-id = 0x0001, mme-code = 01 and mtmsi = remaining string
  guti_1 = "0x0001011a82a179";
  guti_2 = "0x0001011a82a190";
  guti_3 = "0x0001011a82a222";
  uint64_t gutiData;

  magma::map_string_uint64_t guti_map;
  guti_map.map = mme_ue_context_proto.mutable_guti_ue_id_htbl();
  guti_map.set_name("GUTI HASHTABLE");

  EXPECT_EQ(guti_map.get_name(), "GUTI HASHTABLE");
  // Trying to get from an empty map
  EXPECT_EQ(guti_map.get(guti_1, &gutiData), magma::MAP_EMPTY);

  // Inserting new <key,value> pair
  EXPECT_EQ(guti_map.insert(guti_1, 10), magma::MAP_OK);
  EXPECT_EQ(guti_map.insert(guti_2, 20), magma::MAP_OK);

  // Inserting already existing key.Expected: failure.
  EXPECT_EQ(guti_map.insert(guti_1, 10), magma::MAP_KEY_ALREADY_EXISTS);

  // Getting data from map
  EXPECT_EQ(guti_map.get(guti_1, &gutiData), magma::MAP_OK);
  EXPECT_EQ(gutiData, 10);
  // Getting with invalid key
  EXPECT_EQ(guti_map.get(guti_3, &gutiData), magma::MAP_KEY_NOT_EXISTS);

  // Removing entry from table
  EXPECT_EQ(guti_map.remove(guti_1), magma::MAP_OK);

  // Trying to remove from invalid key
  EXPECT_EQ(guti_map.remove(guti_3), magma::MAP_KEY_NOT_EXISTS);
  // Test for clear() and isEmpty()
  guti_map.clear();
  EXPECT_TRUE(guti_map.isEmpty());
}
