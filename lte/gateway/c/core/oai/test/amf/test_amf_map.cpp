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
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_ue_context_and_proc.h"
#include "lte/gateway/c/core/oai/tasks/amf/include/amf_ue_context_storage.h"

using ::testing::Test;

namespace magma5g {
// auto& context_storage = AmfUeContextStorage::getUeContextStorage();
extern AmfUeContextStorage& context_storage;

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

  delete state_cache_p;
}

TEST(test_map, test_amf_state_ue_ht) {
  // Initializations for Map: Key-uint64_t Data-void*
  amf_ue_ngap_id_t ue_id = 1;
  auto ue_context        = context_storage.amf_create_new_ue_context();
}
}  // namespace magma5g
