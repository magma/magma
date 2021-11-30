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

#include <glog/logging.h>
#include <gtest/gtest.h>

#include "lte/gateway/c/sctpd/src/sctp_assoc.h"
#include "lte/gateway/c/sctpd/src/sctp_desc.h"

using ::testing::Test;

namespace magma {
namespace sctpd {

const int DESC_SD = 3;

const int ASSOC_1_ASSOC_ID = 1;
const int ASSOC_1_SD       = 4;
const int ASSOC_2_ASSOC_ID = 2;
const int ASSOC_2_SD       = 5;

class SctpdDescTest : public ::testing::Test {
 protected:
  virtual void SetUp() {
    assoc_1.assoc_id = ASSOC_1_ASSOC_ID;
    assoc_1.sd       = ASSOC_1_SD;
    assoc_2.assoc_id = ASSOC_2_ASSOC_ID;
    assoc_2.sd       = ASSOC_2_SD;
  }

  void check_assoc(int assoc_id, int sd, const SctpAssoc& assoc) {
    EXPECT_EQ(assoc_id, assoc.assoc_id);
    EXPECT_EQ(sd, assoc.sd);
  }

  SctpAssoc assoc_1;
  SctpAssoc assoc_2;
};

TEST_F(SctpdDescTest, test_sctpd_desc) {
  SctpDesc desc(DESC_SD);

  EXPECT_EQ(DESC_SD, desc.sd());

  // check addition an retreival of associations
  desc.addAssoc(assoc_1);
  EXPECT_NO_THROW(desc.getAssoc(ASSOC_1_ASSOC_ID));
  check_assoc(ASSOC_1_ASSOC_ID, ASSOC_1_SD, assoc_1);

  EXPECT_THROW(desc.getAssoc(ASSOC_2_ASSOC_ID), std::out_of_range);

  desc.addAssoc(assoc_2);
  EXPECT_NO_THROW(desc.getAssoc(ASSOC_1_ASSOC_ID));
  check_assoc(ASSOC_1_ASSOC_ID, ASSOC_1_SD, assoc_1);

  EXPECT_NO_THROW(desc.getAssoc(ASSOC_2_ASSOC_ID));
  check_assoc(ASSOC_2_ASSOC_ID, ASSOC_2_SD, assoc_2);

  // check iteration
  bool found_1 = false;
  bool found_2 = false;
  for (auto kv : desc) {
    auto assoc_id = kv.first;
    auto assoc    = kv.second;

    EXPECT_EQ(assoc_id, assoc.assoc_id);
    if (assoc_id == ASSOC_1_ASSOC_ID) {
      found_1 = true;
      check_assoc(ASSOC_1_ASSOC_ID, ASSOC_1_SD, assoc);
    } else if (assoc_id == ASSOC_2_ASSOC_ID) {
      found_2 = true;
      check_assoc(ASSOC_2_ASSOC_ID, ASSOC_2_SD, assoc);
    } else {
      FAIL();
    }
  }

  EXPECT_TRUE(found_1 && found_2);

  // check deletion of associations
  desc.delAssoc(assoc_1.assoc_id);
  EXPECT_THROW(desc.getAssoc(ASSOC_1_ASSOC_ID), std::out_of_range);

  EXPECT_NO_THROW(desc.getAssoc(ASSOC_2_ASSOC_ID));
  check_assoc(ASSOC_2_ASSOC_ID, ASSOC_2_SD, assoc_2);

  desc.delAssoc(assoc_2.assoc_id);
  EXPECT_THROW(desc.getAssoc(ASSOC_1_ASSOC_ID), std::out_of_range);
  EXPECT_THROW(desc.getAssoc(ASSOC_2_ASSOC_ID), std::out_of_range);
}

}  // namespace sctpd
}  // namespace magma

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  FLAGS_logtostderr = 1;
  FLAGS_v           = 10;
  return RUN_ALL_TESTS();
}
