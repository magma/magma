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

/*! \file s11_mme_session_manager.hpp
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#pragma once

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/lib/gtpv2-c/gtpv2c_ie_formatter/shared/gtpv2c_ie_formatter.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#ifdef __cplusplus
}
#endif

/* @brief Create a new Create Session Request and send it to provided S-GW. */
status_code_e s11_mme_create_session_request(
    nw_gtpv2c_stack_handle_t* stack_p,
    itti_s11_create_session_request_t* create_session_p);

/* @brief Handle a Create Session Response received from S-GW. */
status_code_e s11_mme_handle_create_session_response(
    nw_gtpv2c_stack_handle_t* stack_p, nw_gtpv2c_ulp_api_t* pUlpApi);

/* @brief Create a new Delete Session Request and send it to provided S-GW. */
status_code_e s11_mme_delete_session_request(
    nw_gtpv2c_stack_handle_t* stack_p,
    itti_s11_delete_session_request_t* delete_session_p);

status_code_e s11_mme_handle_delete_session_response(
    nw_gtpv2c_stack_handle_t* stack_p, nw_gtpv2c_ulp_api_t* pUlpApi);

/* @brief Create a new Modify Bearer Request and send it to provided S-GW. */
status_code_e s11_mme_modify_bearer_request(
    nw_gtpv2c_stack_handle_t* stack_p,
    itti_s11_modify_bearer_request_t* modify_bearer_p);

status_code_e s11_mme_handle_ulp_error_indicatior(
    nw_gtpv2c_stack_handle_t* stack_p, nw_gtpv2c_ulp_api_t* pUlpApi);
