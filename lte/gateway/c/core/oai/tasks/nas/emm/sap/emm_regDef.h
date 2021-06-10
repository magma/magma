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

Source      emm_regDef.h

Version     0.1

Date        2012/10/16

Product     NAS stack

Subsystem   EPS Mobility Management

Author      Frederic Maurel

Description Defines the EMMREG Service Access Point that provides
        registration services for location updating and attach/detach
        procedures.

*****************************************************************************/
#ifndef FILE_EMM_REGDEF_SEEN
#define FILE_EMM_REGDEF_SEEN

#include "common_types.h"
#include "nas/commonDef.h"
#include "3gpp_36.401.h"
#include "emm_fsm.h"
#include <stdbool.h>

/****************************************************************************/
/*********************  G L O B A L    C O N S T A N T S  *******************/
/****************************************************************************/

/*
 * EMMREG-SAP primitives
 */
typedef enum {
  _EMMREG_START = 0,
  _EMMREG_COMMON_PROC_REQ,   /* EMM common procedure requested   */
  _EMMREG_COMMON_PROC_CNF,   /* EMM common procedure successful  */
  _EMMREG_COMMON_PROC_REJ,   /* EMM common procedure failed, CN send REJECT */
  _EMMREG_COMMON_PROC_ABORT, /* EMM common procedure aborted     */
  _EMMREG_ATTACH_CNF,        /* EPS network attach accepted      */
  _EMMREG_ATTACH_REJ,        /* EPS network attach rejected      */
  _EMMREG_ATTACH_ABORT,      /* EPS network attach aborted      */
  _EMMREG_DETACH_INIT,       /* Network detach initiated         */
  _EMMREG_DETACH_REQ,        /* Network detach requested         */
  _EMMREG_DETACH_FAILED,     /* Network detach attempt failed    */
  _EMMREG_DETACH_CNF,        /* Network detach accepted          */
  _EMMREG_TAU_REQ,
  _EMMREG_TAU_CNF,
  _EMMREG_TAU_REJ,
  _EMMREG_SERVICE_REQ,
  _EMMREG_SERVICE_CNF,
  _EMMREG_SERVICE_REJ,
  _EMMREG_LOWERLAYER_SUCCESS,      /* Data successfully delivered      */
  _EMMREG_LOWERLAYER_FAILURE,      /* Lower layer failure indication   */
  _EMMREG_LOWERLAYER_RELEASE,      /* NAS signalling connection released   */
  _EMMREG_LOWERLAYER_NON_DELIVERY, /*  remote Lower layer failure indication  */
  _EMMREG_END
} emm_reg_primitive_t;

/****************************************************************************/
/************************  G L O B A L    T Y P E S  ************************/
/****************************************************************************/

/*
 * EMMREG primitive for attach procedure
 * -------------------------------------
 */
typedef struct emm_reg_attach_s {
  bool is_emergency; /* true if the UE was attempting to register to
                      * the network for emergency services only  */
  struct nas_emm_attach_proc_s* proc;
} emm_reg_attach_t;

/*
 * EMMREG primitive for attach procedure
 * -------------------------------------
 */
typedef struct emm_reg_tau_s {
  struct nas_emm_tau_proc_s* proc;
} emm_reg_tau_t;

/*
 * EMMREG primitive for detach procedure
 * -------------------------------------
 */
typedef struct emm_reg_detach_s {
  bool switch_off; /* true if the UE is switched off       */
  int type;        /* Network detach type              */
} emm_reg_detach_t;
/*
 * EMMREG primitive for service request procedure
 * -------------------------------------
 */
typedef struct emm_reg_sr_s {
  bool is_emergency; /* true if the UE was attempting to register to
                      * the network for emergency services only  */
  struct nas_sr_proc_s* proc;
} emm_reg_sr_t;

/*
 * EMMREG primitive for EMM common procedures
 * ------------------------------------------
 */
struct nas_emm_common_proc_s;
typedef struct emm_reg_common_s {
  emm_fsm_state_t previous_emm_fsm_state;
  struct nas_emm_common_proc_s* common_proc;
} emm_reg_common_t;
/*
 * EMMREG primitive for Lower Layer success
 * ------------------------------------------
 */
typedef struct emm_reg_ll_success_s {
  uint64_t puid;
  uint16_t msg_len;
#define EMM_REG_MSG_DIGEST_SIZE 16
  size_t digest_len;
  uint8_t msg_digest[EMM_REG_MSG_DIGEST_SIZE];
} emm_reg_ll_sucess_t;

/*
 * EMMREG primitive for Lower Layer failure
 * ------------------------------------------
 */
typedef struct emm_reg_ll_failure_s {
  emm_fsm_state_t previous_emm_fsm_state;
  size_t msg_len;
  size_t digest_len;
  uint8_t msg_digest[EMM_REG_MSG_DIGEST_SIZE];
} emm_reg_ll_failure_t;

/*
 * EMMREG primitive for SDU non delivery due to HO
 * ------------------------------------------
 */
typedef struct emm_reg_sdu_non_delivery_ho_s {
  emm_fsm_state_t previous_emm_fsm_state;
  size_t msg_len;
  size_t digest_len;
  uint8_t msg_digest[EMM_REG_MSG_DIGEST_SIZE];
} emm_reg_sdu_non_delivery_ho_t;

/*
 * Structure of EMMREG-SAP primitive
 */
typedef struct emm_reg_s {
  emm_reg_primitive_t primitive;
  mme_ue_s1ap_id_t ue_id;
  struct emm_context_s* ctx;
  bool notify;  // notify through call-backs
  bool free_proc;

  union {
    emm_reg_attach_t attach;
    emm_reg_detach_t detach;
    emm_reg_tau_t tau;
    emm_reg_sr_t sr;
    emm_reg_common_t common;
    emm_reg_ll_failure_t ll_failure;
    emm_reg_ll_sucess_t ll_success;
    emm_reg_sdu_non_delivery_ho_t non_delivery_ho;
  } u;
} emm_reg_t;

/****************************************************************************/
/********************  G L O B A L    V A R I A B L E S  ********************/
/****************************************************************************/

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

#endif /* FILE_EMM_REGDEF_SEEN*/
