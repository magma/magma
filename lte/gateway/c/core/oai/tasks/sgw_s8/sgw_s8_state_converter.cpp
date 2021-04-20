/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

extern "C" {
#include "dynamic_memory_check.h"
#include "sgw_context_manager.h"
}

#include "sgw_s8_state_converter.h"

using magma::lte::oai::SgwState;
using magma::lte::oai::SgwUeContext;

namespace magma {
namespace lte {

SgwStateConverter::SgwStateConverter()  = default;
SgwStateConverter::~SgwStateConverter() = default;

void SgwStateConverter::state_to_proto(
    const sgw_state_t* sgw_state, SgwState* proto) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  proto->Clear();
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

void SgwStateConverter::proto_to_state(
    const SgwState& proto, sgw_state_t* sgw_state) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

void SgwStateConverter::ue_to_proto(
    const spgw_ue_context_t* ue_state, oai::SgwUeContext* ue_proto) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

void SgwStateConverter::proto_to_ue(
    const oai::SgwUeContext& ue_proto, spgw_ue_context_t* ue_context_p) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

}  // namespace lte
}  // namespace magma
