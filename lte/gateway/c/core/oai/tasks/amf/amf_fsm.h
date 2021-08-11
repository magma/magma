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

#pragma once

#include <sstream>
#include <thread>
#ifdef __cplusplus
extern "C" {
#endif
#include "bstrlib.h"

#ifdef __cplusplus
};
#endif
namespace magma5g {
/*
 * State Events to Trigger UE states and PDU Session States
 */
typedef enum {
  STATE_EVENT_REG_REQUEST,
  STATE_EVENT_SEC_MODE_COMPLETE,
  STATE_EVENT_REG_COMPLETE,
  STATE_EVENT_DEREGISTER,  // Handling Deregister and Deregister init
  STATE_PDU_SESSION_ESTABLISHMENT_REQUEST,
  STATE_PDU_SESSION_ESTABLISHMENT_ACCEPT,
  STATE_PDU_SESSION_RELEASE_COMPLETE,
  STATE_EVENT_CONTEXT_RELEASE,
  STATE_EVENT_MAX,
} state_events;

/*PDU session states*/
typedef enum {
  SESSION_NULL,
  CREATING,
  CREATE,
  ACTIVE,
  INACTIVE,
  PENDING_RELEASE,
  RELEASED,
  SESSION_MAX
} SMSessionFSMState;

/* UE states */
enum m5gmm_state_t {
  DEREGISTERED = 0,
  REGISTERED_IDLE,
  REGISTERED_CONNECTED,
  DEREGISTERED_INITIATED,
  PENDING_RELEASE_RESPONSE,
  COMMON_PROCEDURE_INITIATED1,
  COMMON_PROCEDURE_INITIATED2,
  UE_STATE_MAX
};

typedef struct UE_Handlers_s {
  const char* name;
  void (*func)(void);
} UE_Handlers_t;

typedef struct ue_state_transition_s {
  UE_Handlers_t handler;
  m5gmm_state_t next_state;
  SMSessionFSMState next_sess_state;
} ue_state_transition_t;

void create_state_matrix();
// TODO in upcoming PR of FSM this enum and got changed
enum amf_fsm_state_t {
  AMF_STATE_MIN = 0,
  AMF_DEREGISTERED,
  AMF_REGISTERED,
  AMF_DEREGISTERED_INITIATED,
  AMF_COMMON_PROCEDURE_INITIATED,
  AMF_STATE_MAX
};
}  // namespace magma5g
