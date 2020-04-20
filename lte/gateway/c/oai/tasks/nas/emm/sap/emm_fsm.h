/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

/*****************************************************************************

Source      emm_fsm.h

Version     0.1

Date        2012/10/03

Product     NAS stack

Subsystem   EPS Mobility Management

Author      Frederic Maurel

Description Defines the EPS Mobility Management procedures executed at
        the EMMREG Service Access Point.

*****************************************************************************/
#ifndef FILE_EMM_FSM_SEEN
#define FILE_EMM_FSM_SEEN

/****************************************************************************/
/*********************  G L O B A L    C O N S T A N T S  *******************/
/****************************************************************************/

/****************************************************************************/
/************************  G L O B A L    T Y P E S  ************************/
/****************************************************************************/

/*
 * States of the EPS Mobility Management sublayer
 * ----------------------------------------------
 * The EMM protocol of the UE and the network is described by means of two
 * different state machines.
 */
typedef enum {
  EMM_STATE_MIN = 0,
  EMM_INVALID   = EMM_STATE_MIN,
  EMM_DEREGISTERED,
  EMM_REGISTERED,
  EMM_DEREGISTERED_INITIATED,
  EMM_COMMON_PROCEDURE_INITIATED,
  EMM_STATE_MAX
} emm_fsm_state_t;

#include "emm_regDef.h"
#include "3gpp_36.401.h"

struct emm_context_s;
struct emm_reg_s;
/****************************************************************************/
/********************  G L O B A L    V A R I A B L E S  ********************/
/****************************************************************************/

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

void emm_fsm_initialize(void);

struct emm_context_s;
struct emm_reg_s;

int emm_fsm_set_state(
    const mme_ue_s1ap_id_t ueid, struct emm_context_s* const emm_context,
    const emm_fsm_state_t status);
emm_fsm_state_t emm_fsm_get_state(
    const struct emm_context_s* const emm_context);
const char* emm_fsm_get_state_str(
    const struct emm_context_s* const emm_context);
int emm_fsm_process(struct emm_reg_s* const evt);

#endif /* FILE_EMM_FSM_SEEN*/
