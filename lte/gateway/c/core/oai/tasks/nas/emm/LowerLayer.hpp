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
Source      lowerlayer.hpp

Version     0.1

Date        2013/06/19

Product     NAS stack

Subsystem   EPS Mobility Management

Author      Frederic Maurel, Lionel GAUTHIER

Description Defines EMM procedures executed by the Non-Access Stratum
        upon receiving notifications from lower layers so that data
        transfer succeed or failed, or NAS signalling connection is
        released, or ESM unit data has been received from under layer,
        and to request ESM unit data transfer to under layer.

*****************************************************************************/
#pragma once

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.007.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_36.401.h"

/****************************************************************************/
/*********************  G L O B A L    C O N S T A N T S  *******************/
/****************************************************************************/

/****************************************************************************/
/************************  G L O B A L    T Y P E S  ************************/
/****************************************************************************/

/****************************************************************************/
/********************  G L O B A L    V A R I A B L E S  ********************/
/****************************************************************************/

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

status_code_e lowerlayer_success(mme_ue_s1ap_id_t ue_id, bstring* nas_msg);
status_code_e lowerlayer_failure(mme_ue_s1ap_id_t ueid,
                                 STOLEN_REF bstring* nas_msg);
status_code_e lowerlayer_non_delivery_indication(mme_ue_s1ap_id_t ue_id,
                                                 STOLEN_REF bstring* nas_msg);
status_code_e lowerlayer_establish(void);
status_code_e lowerlayer_release(mme_ue_s1ap_id_t ue_id, int cause);

status_code_e lowerlayer_data_ind(mme_ue_s1ap_id_t ueid, const_bstring data);
status_code_e lowerlayer_data_req(mme_ue_s1ap_id_t ueid, bstring data);
status_code_e lowerlayer_activate_bearer_req(
    const mme_ue_s1ap_id_t ue_id, const ebi_t ebi, const bitrate_t mbr_dl,
    const bitrate_t mbr_ul, const bitrate_t gbr_dl, const bitrate_t gbr_ul,
    bstring data);
status_code_e lowerlayer_deactivate_bearer_req(const mme_ue_s1ap_id_t ue_id,
                                               const ebi_t ebi, bstring data);
