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

Source      mme_app_sgs_fsm.h

Version

Date

Product

Subsystem

Author

Description Defines the SGS State Machine handling

*****************************************************************************/
#ifndef FILE_SGS_FSM_SEEN
#define FILE_SGS_FSM_SEEN

#include "common_types.h"
#include "sgs_messages_types.h"
#include "3gpp_36.401.h"
/****************************************************************************/
/*********************  G L O B A L    C O N S T A N T S  *******************/
/****************************************************************************/

/****************************************************************************/
/************************  G L O B A L    T Y P E S  ************************/
/****************************************************************************/

/****************************************************************************/
/********************  G L O B A L    V A R I A B L E S  ********************/
/****************************************************************************/

typedef enum {
  _SGS_LOCATION_UPDATE_ACCEPT,
  _SGS_LOCATION_UPDATE_REJECT,
  _SGS_PAGING_REQUEST,
  _SGS_SERVICE_ABORT_REQUEST,
  _SGS_EPS_DETACH_IND,
  _SGS_IMSI_DETACH_IND,
  _SGS_RESET_INDICATION,
} sgs_primitive_t;

typedef enum sgs_fsm_state_e {
  SGS_STATE_MIN = 0,
  SGS_INVALID   = SGS_STATE_MIN,
  SGS_NULL,
  SGS_LA_UPDATE_REQUESTED,
  SGS_ASSOCIATED,
  SGS_STATE_MAX
} sgs_fsm_state_t;

/*
 * Structure of SGS-AP primitive
 */
typedef struct {
  sgs_primitive_t primitive; /* Primitive to identify SGSAP messages */
  mme_ue_s1ap_id_t ue_id;    /* mme_ue_s1ap_id to uniquely identify the UE */
  void*
      ctx; /* void pointer to point context and holds sgsap received message */
} sgs_fsm_t;

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

int sgs_fsm_process(const sgs_fsm_t* sgs_evt);
int sgs_fsm_set_status(
    mme_ue_s1ap_id_t ue_id, void* ctx, sgs_fsm_state_t state);
void sgs_fsm_initialize(void);
sgs_fsm_state_t sgs_fsm_get_status(mme_ue_s1ap_id_t ueid, void* ctx);
int sgs_fsm_process(const sgs_fsm_t* sgs_evt);

int sgs_fsm_null_loc_updt_acc(const sgs_fsm_t* fsm_evt);
int sgs_fsm_null_loc_updt_rej(const sgs_fsm_t* fsm_evt);
int sgs_fsm_la_updt_req_loc_updt_acc(const sgs_fsm_t* fsm_evt);
int sgs_fsm_la_updt_req_loc_updt_rej(const sgs_fsm_t* fsm_evt);
int sgs_fsm_associated_loc_updt_acc(const sgs_fsm_t* fsm_evt);
int sgs_fsm_associated_loc_updt_rej(const sgs_fsm_t* fsm_evt);
int sgs_handle_associated_paging_request(const sgs_fsm_t* sgs_evt);
int sgs_handle_null_paging_request(const sgs_fsm_t* sgs_evt);
int sgs_fsm_associated_service_abort_request(const sgs_fsm_t* fsm_evt);
#endif /* FILE_SGS_FSM_SEEN*/
