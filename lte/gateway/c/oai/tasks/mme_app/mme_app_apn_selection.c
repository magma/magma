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

/*! \file mme_app_apn_selection.c
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#include <stdio.h>
#include <string.h>

#include "bstrlib.h"
#include "log.h"
#include "common_types.h"
#include "mme_app_ue_context.h"
#include "mme_app_apn_selection.h"
#include "emm_data.h"
#include "conversions.h"

//------------------------------------------------------------------------------
int select_pdn_type(
    struct apn_configuration_s* apn_config,
    esm_proc_pdn_type_t ue_selected_pdn_type, esm_cause_t* esm_cause) {
  /* Overwrite apn_config->pdn_type based on the PDN type sent by UE and the PDN
   * Type received in subscription profile
   */
  switch (ue_selected_pdn_type) {
    case ESM_PDN_TYPE_IPV4:
      if ((apn_config->pdn_type == IPv4_AND_v6) ||
          (apn_config->pdn_type == IPv4_OR_v6)) {
        apn_config->pdn_type = IPv4;
      } else if (apn_config->pdn_type == IPv6) {
        /* As per 3gpp 23.401-cb0 sec 5.3.1, If the requested PDN type is IPv4
         * or IPv6, and either the requested PDN type or PDN type IPv4v6 are
         * subscribed, the MME sets the PDN type as requested. Otherwise the PDN
         * connection request is rejected
         */
        OAILOG_ERROR(
            LOG_MME_APP,
            " Sending PDN Connectivity Reject with cause "
            "ESM_CAUSE_UNKNOWN_PDN_TYPE,"
            " UE requested PDN Type %d, subscribed PDN Type %d \n",
            ue_selected_pdn_type, apn_config->pdn_type);

        *esm_cause = ESM_CAUSE_UNKNOWN_PDN_TYPE;
        OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
      }
      break;

    case ESM_PDN_TYPE_IPV6:
      if ((apn_config->pdn_type == IPv4_AND_v6) ||
          (apn_config->pdn_type == IPv4_OR_v6)) {
        apn_config->pdn_type = IPv6;
      } else if (apn_config->pdn_type == IPv4) {
        /* As per 3gpp 23.401-cb0 sec 5.3.1, If the requested PDN type is IPv4
         * or IPv6, and either the requested PDN type or PDN type IPv4v6 are
         * subscribed, the MME sets the PDN type as requested. Otherwise the PDN
         * connection request is rejected
         */
        OAILOG_ERROR(
            LOG_MME_APP,
            " Sending PDN Connectivity Reject with cause "
            "ESM_CAUSE_UNKNOWN_PDN_TYPE,"
            " UE requested PDN Type %d, subscribed PDN Type %d \n",
            ue_selected_pdn_type, apn_config->pdn_type);

        *esm_cause = ESM_CAUSE_UNKNOWN_PDN_TYPE;
        OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
      }
      break;

    case ESM_PDN_TYPE_IPV4V6:
      if (apn_config->pdn_type == IPv4_OR_v6) {
        apn_config->pdn_type = IPv4;
      } else if (apn_config->pdn_type == IPv4) {
        *esm_cause = ESM_CAUSE_PDN_TYPE_IPV4_ONLY_ALLOWED;
      } else if (apn_config->pdn_type == IPv6) {
        *esm_cause = ESM_CAUSE_PDN_TYPE_IPV6_ONLY_ALLOWED;
      }
      break;

    default:
      OAILOG_ERROR(
          LOG_MME_APP,
          " Sending PDN Connectivity Reject with cause "
          "ESM_CAUSE_UNKNOWN_PDN_TYPE,"
          " UE requested PDN Type %d, subscribed PDN Type %d \n",
          ue_selected_pdn_type, apn_config->pdn_type);

      OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
      break;
  }
  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
}

//------------------------------------------------------------------------------
struct apn_configuration_s* mme_app_select_apn(
    ue_mm_context_t* const ue_context, int* esm_cause) {
  context_identifier_t default_context_identifier =
      ue_context->apn_config_profile.context_identifier;
  int index;
  int rc = RETURNok;

  const_bstring const ue_selected_apn =
      ue_context->emm_context.esm_ctx.esm_proc_data->apn;
  esm_proc_pdn_type_t ue_selected_pdn_type =
      ue_context->emm_context.esm_ctx.esm_proc_data->pdn_type;

  for (index = 0; index < ue_context->apn_config_profile.nb_apns; index++) {
    if (!ue_selected_apn) {
      /*
       * OK we got our default APN
       */
      if (ue_context->apn_config_profile.apn_configuration[index]
              .context_identifier == default_context_identifier) {
        // Select PDN Type
        rc = select_pdn_type(
            &ue_context->apn_config_profile.apn_configuration[index],
            ue_selected_pdn_type, esm_cause);
        if (*esm_cause == ESM_CAUSE_UNKNOWN_PDN_TYPE || rc == RETURNerror) {
          return NULL;
        }
        OAILOG_INFO(
            LOG_MME_APP,
            "Selected APN <%s>, PDN Type <%d> for UE " IMSI_64_FMT "\n",
            ue_context->apn_config_profile.apn_configuration[index]
                .service_selection,
            ue_context->apn_config_profile.apn_configuration[index].pdn_type,
            ue_context->emm_context._imsi64);
        return &ue_context->apn_config_profile.apn_configuration[index];
      }
    } else {
      /*
       * OK we got the UE selected APN
       */
      if (biseqcaselessblk(
              ue_selected_apn,
              ue_context->apn_config_profile.apn_configuration[index]
                  .service_selection,
              strlen(ue_context->apn_config_profile.apn_configuration[index]
                         .service_selection)) == 1) {
        // Select PDN Type
        rc = select_pdn_type(
            &ue_context->apn_config_profile.apn_configuration[index],
            ue_selected_pdn_type, esm_cause);
        if (*esm_cause == ESM_CAUSE_UNKNOWN_PDN_TYPE || rc == RETURNerror) {
          return NULL;
        }
        OAILOG_INFO(
            LOG_MME_APP,
            "Selected APN <%s>, PDN Type <%d> for UE " IMSI_64_FMT "\n",
            ue_context->apn_config_profile.apn_configuration[index]
                .service_selection,
            ue_context->apn_config_profile.apn_configuration[index].pdn_type,
            ue_context->emm_context._imsi64);

        return &ue_context->apn_config_profile.apn_configuration[index];
      }
    }
  }
  *esm_cause = ESM_CAUSE_UNKNOWN_ACCESS_POINT_NAME;
  return NULL;
}

//------------------------------------------------------------------------------
struct apn_configuration_s* mme_app_get_apn_config(
    ue_mm_context_t* const ue_context,
    const context_identifier_t context_identifier) {
  int index;

  for (index = 0; index < ue_context->apn_config_profile.nb_apns; index++) {
    if (ue_context->apn_config_profile.apn_configuration[index]
            .context_identifier == context_identifier) {
      return &ue_context->apn_config_profile.apn_configuration[index];
    }
  }
  return NULL;
}

bstring mme_app_process_apn_correction(imsi_t* imsi, bstring accesspointname) {
  int i;
  char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
  apn_map_config_t config = mme_config.nas_config.apn_map_config;

  IMSI_TO_STRING(imsi, imsi_str, IMSI_BCD_DIGITS_MAX + 1);
  for (i = 0; i < config.nb; i++) {
    const char* imsi_prefix = bdata(config.apn_map[i].imsi_prefix);
    int imsi_prefix_len     = strlen(imsi_prefix);
    if ((imsi_prefix_len <= IMSI_BCD_DIGITS_MAX) &&
        !strncmp(imsi_prefix, imsi_str, imsi_prefix_len)) {
      return config.apn_map[i].apn_override;
    }
  }
  return accesspointname;
}
