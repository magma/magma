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
#ifndef HA_DEFS_H_
#define HA_DEFS_H_

#include "intertask_interface.h"
#include "mme_config.h"

extern task_zmq_ctx_t ha_task_zmq_ctx;

int ha_init(const mme_config_t* mme_config);

/*
 * Syncs up with orc8r and fetches eNB connection state
 * to its primary AGW. Inside the function:
 * state-0: No action taken for such eNBs
 * state-1: eNB has established S1 connection with the primary;
 *          AGW will try to offload one of the UEs to test waters.
 * state-2: eNB has established S1 connection with the primary AGW and
 *          at least one UE served by eNB is camped on the primary.
 *          Offload remaining UEs from this gateway.
 */
bool sync_up_with_orc8r(void);

/*
 * Sends a S1AP_UE_CONTEXT_RELEASE_REQ message to MME.
 */
void handle_agw_offload_req(ha_agw_offload_req_t* offload_req);

#endif /* HA_DEFS_H_ */
