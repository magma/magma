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
#include "lte/gateway/c/core/oai/tasks/amf/include/amf_smf_session_context.h"

namespace magma5g {
status_code_e amf_smf_context_ue_aggregate_max_bit_rate_set(
    amf_context_s* amf_ctxt_p, ambr_t subscribed_ue_ambr) {
  memcpy(&amf_ctxt_p->subscribed_ue_ambr, &subscribed_ue_ambr, sizeof(ambr_t));

  return RETURNok;
}

status_code_e amf_smf_context_ue_aggregate_max_bit_rate_get(
    const amf_context_s* amf_ctxt_p, bit_rate_t* subscriber_ambr_dl,
    bit_rate_t* subscriber_ambr_ul) {
  if (amf_ctxt_p == NULL) {
    return RETURNerror;
  }

  *subscriber_ambr_dl = amf_ctxt_p->subscribed_ue_ambr.br_dl;
  *subscriber_ambr_ul = amf_ctxt_p->subscribed_ue_ambr.br_ul;

  return RETURNok;
}

}  // namespace magma5g
