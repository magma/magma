/**
 * Copyright 2022 The Magma Authors.
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
#include "lte/gateway/c/core/oai/include/s1ap_state.h"
}

#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_state_manager.h"
#include "lte/gateway/c/core/oai/test/mock_tasks/mock_tasks.h"

namespace magma {
namespace lte {

class S1apMmeHandlersTest : public ::testing::Test {
  virtual void SetUp();
  virtual void TearDown();

 protected:
  std::shared_ptr<MockMmeAppHandler> mme_app_handler;
  std::shared_ptr<MockSctpHandler> sctp_handler;
  s1ap_state_t* state;
  sctp_assoc_id_t assoc_id;
  sctp_stream_id_t stream_id;
};

}  // namespace lte
}  // namespace magma