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
#ifndef FILE_MME_UE_CONTEXT_H_SEEN
#define FILE_MME_UE_CONTEXT_H_SEEN
#include "lte/gateway/c/core/oai/common/common_utility_funs.hpp"

namespace magma {
namespace lte {

struct MmeUeContext : public magma::utils::AppTimerUeContext<timer_arg_t> {
  static MmeUeContext& Instance();
  explicit MmeUeContext(task_zmq_ctx_s* zctx)
      : magma::utils::AppTimerUeContext<timer_arg_t>{zctx} {}
};

}  // namespace lte
}  // namespace magma
#endif /* FILE_MME_UE_CONTEXT_H_SEEN */
