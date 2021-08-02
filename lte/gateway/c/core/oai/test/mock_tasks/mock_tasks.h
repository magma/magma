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

extern "C" {
#include "dynamic_memory_check.h"
#define CHECK_PROTOTYPE_ONLY
#include "intertask_interface_init.h"
#undef CHECK_PROTOTYPE_ONLY
#include "intertask_interface.h"
#include "intertask_interface_types.h"
#include "itti_free_defined_msg.h"
}

const task_info_t tasks_info[] = {
    {THREAD_NULL, "TASK_UNKNOWN", "ipc://IPC_TASK_UNKNOWN"},
#define TASK_DEF(tHREADiD)                                                     \
  {THREAD_##tHREADiD, #tHREADiD, "ipc://IPC_" #tHREADiD},
#include <tasks_def.h>
#undef TASK_DEF
};

/* Map message id to message information */
const message_info_t messages_info[] = {
#define MESSAGE_DEF(iD, sTRUCT, fIELDnAME) {iD, sizeof(sTRUCT), #iD},
#include <messages_def.h>
#undef MESSAGE_DEF
};

void start_mock_ha_task();
void start_mock_s1ap_task();
void start_mock_s6a_task();
void start_mock_s11_task();
void start_mock_service303_task();
void start_mock_sgs_task();
void start_mock_sgw_s8_task();
void start_mock_sms_orc8r_task();
void start_mock_spgw_task();