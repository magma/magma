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

#include "lte/gateway/c/core/oai/test/ngap/util_ngap_pkt.hpp"
extern "C" {
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/tasks/ngap/ngap_amf_handlers.h"
#include "lte/gateway/c/core/oai/include/amf_config.hpp"
}
#include "lte/gateway/c/core/oai/tasks/ngap/ngap_state_manager.hpp"

using ::testing::Test;

namespace magma5g {
TEST(test_ngap_handle_new_association, empty_initial_state) {
  ngap_state_t* state = create_ngap_state(2, 2);

  bstring ran_cp_ipaddr = bfromcstr("\xc0\xa8\x3c\x8d");
  sctp_new_peer_t p = {
      .instreams = 1,
      .outstreams = 2,
      .assoc_id = 3,
      .ran_cp_ipaddr = ran_cp_ipaddr,
  };

  EXPECT_EQ(ngap_handle_new_association(state, &p), RETURNok);
  EXPECT_EQ(state->gnbs.num_elements, 1);

  gnb_description_t* gnbd = nullptr;
  EXPECT_EQ(hashtable_ts_get(&state->gnbs, (const hash_key_t)p.assoc_id,
                             reinterpret_cast<void**>(&gnbd)),
            HASH_TABLE_OK);
  EXPECT_EQ(gnbd->sctp_assoc_id, 3);
  EXPECT_EQ(gnbd->instreams, 1);
  EXPECT_EQ(gnbd->outstreams, 2);
  EXPECT_EQ(gnbd->gnb_id, 0);
  EXPECT_EQ(gnbd->ng_state, NGAP_INIT);
  EXPECT_EQ(gnbd->next_sctp_stream, 1);

  EXPECT_EQ(state->num_gnbs, 1);

  bdestroy(ran_cp_ipaddr);

  free_ngap_state(state);
}

}  // namespace magma5g
