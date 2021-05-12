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

#include <google/protobuf/util/message_differencer.h>
#include <gtest/gtest.h>
#include <mutator.h>

#include "bstrlib.h"
#include "lte/protos/oai/mme_nas_state.pb.h"
#include "mme_app_state_converter.h"

constexpr unsigned int kSeed           = 17;
constexpr unsigned int kMaxUeHtblLists = 1000;
constexpr size_t kMaxSizeHint          = 1500;

void populate_state_struct(mme_app_desc_t* mme_nas_state) {
  mme_nas_state->mme_app_ue_s1ap_id_generator = 1;
  bstring b = bfromcstr("IMSI_UE_ID_TABLE_NAME");
  mme_nas_state->mme_ue_contexts.imsi_mme_ue_id_htbl =
      hashtable_uint64_ts_create(kMaxUeHtblLists, nullptr, b);
  btrunc(b, 0);
  bassigncstr(b, "TUN_UE_ID_TABLE_NAME");
  mme_nas_state->mme_ue_contexts.tun11_ue_context_htbl =
      hashtable_uint64_ts_create(kMaxUeHtblLists, nullptr, b);
  btrunc(b, 0);
  bassigncstr(b, "ENB_UE_ID_MME_UE_ID_TABLE_NAME");
  mme_nas_state->mme_ue_contexts.enb_ue_s1ap_id_ue_context_htbl =
      hashtable_uint64_ts_create(kMaxUeHtblLists, nullptr, b);
  btrunc(b, 0);
  bassigncstr(b, "GUTI_UE_ID_TABLE_NAME");
  mme_nas_state->mme_ue_contexts.guti_ue_context_htbl =
      obj_hashtable_uint64_ts_create(kMaxUeHtblLists, nullptr, nullptr, b);
  bdestroy(b);
}

void destroy_state_struct(mme_app_desc_t* mme_nas_state) {
  hashtable_uint64_ts_destroy(
      mme_nas_state->mme_ue_contexts.imsi_mme_ue_id_htbl);
  hashtable_uint64_ts_destroy(
      mme_nas_state->mme_ue_contexts.tun11_ue_context_htbl);
  hashtable_uint64_ts_destroy(
      mme_nas_state->mme_ue_contexts.enb_ue_s1ap_id_ue_context_htbl);
  obj_hashtable_uint64_ts_destroy(
      mme_nas_state->mme_ue_contexts.guti_ue_context_htbl);
  free(mme_nas_state);
}

TEST(MmeAppStateConverter, DoubleConversionPropertyTest) {
  magma::lte::oai::MmeNasState nas_state_proto;

  protobuf_mutator::Mutator mutator;
  mutator.Seed(kSeed);

  std::stringstream result_text;
  for (int i = 0; i < 2000; ++i) {
    std::cout << "Beginning Mutation #" << i << std::endl;
    mutator.Mutate(&nas_state_proto, kMaxSizeHint);

    // TODO: This default-to-one behavior is only valid if this is an
    // uninitialized state, we should fix the system behavior.
    nas_state_proto.set_mme_app_ue_s1ap_id_generator(1);

    result_text << "IN:" << nas_state_proto.DebugString();

    // Create a state struct
    mme_app_desc_t* mme_nas_state =
        (mme_app_desc_t*) calloc(1, sizeof(mme_app_desc_t));
    populate_state_struct(mme_nas_state);

    magma::lte::MmeNasStateConverter converter;

    try {
      converter.proto_to_state(nas_state_proto, mme_nas_state);
    } catch (std::invalid_argument const& err) {
      EXPECT_EQ(err.what(), std::string("stoul"));
      // If exception was thrown, do not compare outputs
      destroy_state_struct(mme_nas_state);
      continue;
    } catch (...) {
      FAIL() << "Expected std::invalid_argument exception if any";
    }

    magma::lte::oai::MmeNasState result_proto;
    converter.state_to_proto(mme_nas_state, &result_proto);

    google::protobuf::util::MessageDifferencer diff;
    result_text << "OUT:" << result_proto.DebugString();
    EXPECT_TRUE(diff.Equivalent(nas_state_proto, result_proto))
        << result_text.str();

    // Cleanup the state struct
    destroy_state_struct(mme_nas_state);
  }
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}