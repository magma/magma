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
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_ue_context_and_proc.hpp"

extern "C" {
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/include/amf_config.hpp"
}

using ::testing::Test;

namespace magma5g {

/* Test for algorithm selection */
TEST(test_algorithm_selection, ue_security_capabilities) {
  // Priority order for integrity algorithm : { IA2 IA1 IA0 }
  amf_config.nas_config.preferred_integrity_algorithm[0] = 2;
  amf_config.nas_config.preferred_integrity_algorithm[1] = 1;
  amf_config.nas_config.preferred_integrity_algorithm[2] = 0;
  // Priority order for ciphering algorithm : { EA0 EA1 EA2}
  amf_config.nas_config.preferred_ciphering_algorithm[0] = 0;
  amf_config.nas_config.preferred_ciphering_algorithm[1] = 1;
  amf_config.nas_config.preferred_ciphering_algorithm[2] = 2;

  int in_IA, in_EA;
  int out_IA = 0, out_EA = 0;

  // If UE supports all the above mentioned algortithms, then AMF should select
  // IA2 and EA0.
  in_IA = 0xE0;  // 1110 0000
  in_EA = 0xE0;  // 1110 0000
  EXPECT_EQ(m5g_security_select_algorithms(in_IA, in_EA, &out_IA, &out_EA),
            RETURNok);
  EXPECT_EQ(out_IA, 2);
  EXPECT_EQ(out_EA, 0);

  // If UE supports IA0, IA1 and EA0, EA1 , then AMF should select IA1 and EA0.
  in_IA = 0xC0;  // 1100 0000
  in_EA = 0xC0;  // 1100 0000
  EXPECT_EQ(m5g_security_select_algorithms(in_IA, in_EA, &out_IA, &out_EA),
            RETURNok);
  EXPECT_EQ(out_IA, 1);
  EXPECT_EQ(out_EA, 0);

  // If UE supports IA0, IA1, IA2 and EA1 EA2, then AMF should select IA2 and
  // EA1.
  in_IA = 0x60;  // 0110 0000
  in_EA = 0x60;  // 0110 0000
  EXPECT_EQ(m5g_security_select_algorithms(in_IA, in_EA, &out_IA, &out_EA),
            RETURNok);
  EXPECT_EQ(out_IA, 2);
  EXPECT_EQ(out_EA, 1);
}

}  // namespace magma5g
