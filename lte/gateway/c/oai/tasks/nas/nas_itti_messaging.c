/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the Apache License, Version 2.0  (the "License"); you may not use this file
 * except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
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

/*! \file nas_itti_messaging.c
   \brief
   \author  Sebastien ROUX, Lionel GAUTHIER
   \date
   \email: lionel.gauthier@eurecom.fr
*/

#include <ctype.h>
#include <stdio.h>
#include <string.h>
#include <stdbool.h>
#include <stdint.h>
#include <time.h>

#include "bstrlib.h"
#include "log.h"
#include "assertions.h"
#include "conversions.h"
#include "intertask_interface.h"
#include "common_defs.h"
#include "secu_defs.h"
#include "mme_app_ue_context.h"
#include "esm_proc.h"
#include "nas_itti_messaging.h"
#include "nas_proc.h"
#include "emm_proc.h"
#include "3gpp_24.008.h"
#include "3gpp_24.301.h"
#include "3gpp_29.274.h"
#include "3gpp_33.401.h"
#include "EpsAttachType.h"
#include "common_ies.h"
#include "emm_data.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
#include "mme_app_state.h"
#include "mme_app_messages_types.h"
#include "nas_messages_types.h"
#include "nas_procedures.h"
#include "nas_timer.h"
#include "s6a_messages_types.h"
#include "nas/securityDef.h"
#include "sgs_messages_types.h"

#define TASK_ORIGIN TASK_NAS_MME

//***************************************************************************
void s6a_auth_info_rsp_timer_expiry_handler(void *args)
{
  OAILOG_FUNC_IN(LOG_NAS_EMM);

  emm_context_t *emm_ctx = (emm_context_t *) (args);

  if (emm_ctx) {
    nas_auth_info_proc_t *auth_info_proc =
      get_nas_cn_procedure_auth_info(emm_ctx);
    if (!auth_info_proc) {
      OAILOG_FUNC_OUT(LOG_NAS_EMM);
    }

    void *timer_callback_args = NULL;
    nas_stop_Ts6a_auth_info(
      auth_info_proc->ue_id, &auth_info_proc->timer_s6a, timer_callback_args);

    auth_info_proc->timer_s6a.id = NAS_TIMER_INACTIVE_ID;
    if (auth_info_proc->resync) {
      OAILOG_ERROR(
        LOG_NAS_EMM,
        "EMM-PROC  - Timer timer_s6_auth_info_rsp expired. Resync auth "
        "procedure was in progress. Aborting attach procedure. UE "
        "id " MME_UE_S1AP_ID_FMT "\n",
        auth_info_proc->ue_id);
    } else {
      OAILOG_ERROR(
        LOG_NAS_EMM,
        "EMM-PROC  - Timer timer_s6_auth_info_rsp expired. Initial auth "
        "procedure was in progress. Aborting attach procedure. UE "
        "id " MME_UE_S1AP_ID_FMT "\n",
        auth_info_proc->ue_id);
    }

    // Send Attach Reject with cause NETWORK FAILURE and delete UE context
    nas_proc_auth_param_fail(auth_info_proc->ue_id, NAS_CAUSE_NETWORK_FAILURE);
  } else {
    OAILOG_ERROR(
      LOG_NAS_EMM,
      "EMM-PROC  - Timer timer_s6_auth_info_rsp expired. Null EMM Context for "
      "UE \n");
  }

  OAILOG_FUNC_OUT(LOG_NAS_EMM);
}
