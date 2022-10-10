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

/*! \file s6a_subscription_data.cpp
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#include <stdint.h>
#include <netinet/in.h>
#include <stdio.h>
#include <string.h>

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/common/assertions.h"
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.008.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/oai/tasks/s6a/s6a_defs.hpp"

struct avp;

static inline int s6a_parse_subscriber_status(struct avp_hdr* hdr_sub_status,
                                              subscriber_status_t* sub_status) {
  DevCheck(hdr_sub_status->avp_value->u32 < SS_MAX,
           hdr_sub_status->avp_value->u32, SS_MAX, 0);
  *sub_status = (subscriber_status_t)hdr_sub_status->avp_value->u32;
  return RETURNok;
}

static inline int s6a_parse_msisdn(struct avp_hdr* hdr_msisdn, char* msisdn,
                                   uint8_t* length) {
  int i;

  DevCheck(hdr_msisdn->avp_value->os.len <= MSISDN_LENGTH,
           hdr_msisdn->avp_value->os.len, MSISDN_LENGTH, 0);

  if (hdr_msisdn->avp_value->os.len == 0) return RETURNok;

  *length = (int)hdr_msisdn->avp_value->os.len;

  for (i = 0; i < (*length); i++) {
    msisdn[i] = hdr_msisdn->avp_value->os.data[i];
  }
  return RETURNok;
}

static inline int s6a_parse_network_access_mode(
    struct avp_hdr* hdr_network_am, network_access_mode_t* access_mode) {
  DevCheck(hdr_network_am->avp_value->u32 < NAM_MAX &&
               hdr_network_am->avp_value->u32 != NAM_RESERVED,
           hdr_network_am->avp_value->u32, NAM_MAX, NAM_RESERVED);
  *access_mode = (network_access_mode_t)hdr_network_am->avp_value->u32;
  return RETURNok;
}

static inline int s6a_parse_access_restriction_data(
    struct avp_hdr* hdr_access_restriction,
    access_restriction_t* access_restriction) {
  DevCheck(hdr_access_restriction->avp_value->u32 < ARD_MAX,
           hdr_access_restriction->avp_value->u32, ARD_MAX, 0);
  *access_restriction = hdr_access_restriction->avp_value->u32;
  return RETURNok;
}

static inline int s6a_parse_bitrate(struct avp_hdr* hdr_bitrate,
                                    bitrate_t* bitrate) {
  *bitrate = hdr_bitrate->avp_value->u32;
  return RETURNok;
}

static inline int s6a_parse_ambr(struct avp* avp_ambr, ambr_t* ambr) {
  struct avp* avp = NULL;
  struct avp_hdr* hdr;
  bitrate_t ext_max_req_bw_ul = ULONG_MAX;
  bitrate_t ext_max_req_bw_dl = ULONG_MAX;

  CHECK_FCT(fd_msg_browse_internal(avp_ambr, MSG_BRW_FIRST_CHILD,
                                   reinterpret_cast<msg_or_avp**>(&avp), NULL));

  if (!avp) {
    /*
     * First Child avps for ambr are mandatory
     */
    return RETURNerror;
  }

  while (avp) {
    CHECK_FCT(fd_msg_avp_hdr(avp, &hdr));
    if (hdr) {
      switch (hdr->avp_code) {
        case AVP_CODE_MAX_REQUESTED_BANDWIDTH_UL:
          CHECK_FCT(s6a_parse_bitrate(hdr, &ambr->br_ul));
          break;

        case AVP_CODE_MAX_REQUESTED_BANDWIDTH_DL:
          CHECK_FCT(s6a_parse_bitrate(hdr, &ambr->br_dl));
          break;

        case AVP_CODE_EXTENDED_MAX_REQUESTED_BW_UL:
          CHECK_FCT(s6a_parse_bitrate(hdr, &ext_max_req_bw_ul));
          break;

        case AVP_CODE_EXTENDED_MAX_REQUESTED_BW_DL:
          CHECK_FCT(s6a_parse_bitrate(hdr, &ext_max_req_bw_dl));
          break;

        default:
          OAILOG_DEBUG(LOG_S6A, "AMBR child AVP %u is silently discarded\n",
                       hdr->avp_code);
          return RETURNerror;
      }
    } else {
      OAILOG_DEBUG(LOG_S6A, "AMBR child AVP header error\n");
    }
    /*
     * Go to next AVP in the grouped AVP
     */
    CHECK_FCT(fd_msg_browse_internal(
        avp, MSG_BRW_NEXT, reinterpret_cast<msg_or_avp**>(&avp), NULL));
  }
  if ((ambr->br_ul != 4294967295) && (ext_max_req_bw_ul == ULONG_MAX)) {
    ambr->br_unit = BPS;
  } else if ((ambr->br_ul == 4294967295)) {
    ambr->br_unit = KBPS;
    ambr->br_ul = ext_max_req_bw_ul;
  } else {
    OAILOG_DEBUG(LOG_S6A, "AMBR UL parsing error\n");
    return RETURNerror;
  }
  if ((ambr->br_dl != 4294967295) && (ext_max_req_bw_dl == ULONG_MAX)) {
    // harmonize
    if (ambr->br_unit == KBPS) {
      ambr->br_dl = ambr->br_dl / 1000;
    }
  } else if ((ambr->br_dl == 4294967295)) {
    if (ambr->br_unit == BPS) {
      ambr->br_unit = KBPS;
      ambr->br_ul = ambr->br_ul / 1000;
    }
    ambr->br_dl = ext_max_req_bw_dl;
  } else {
    OAILOG_DEBUG(LOG_S6A, "AMBR DL parsing error\n");
    return RETURNerror;
  }
  return RETURNok;
}

static inline int s6a_parse_all_apn_conf_inc_ind(struct avp_hdr* hdr,
                                                 all_apn_conf_ind_t* ptr) {
  DevCheck(hdr->avp_value->u32 < ALL_APN_MAX, hdr->avp_value->u32, ALL_APN_MAX,
           0);
  *ptr = (all_apn_conf_ind_t)hdr->avp_value->u32;
  return RETURNok;
}

static inline int s6a_parse_pdn_type(struct avp_hdr* hdr,
                                     pdn_type_t* pdn_type) {
  DevCheck(hdr->avp_value->u32 < IP_MAX, hdr->avp_value->u32, IP_MAX, 0);
  *pdn_type = hdr->avp_value->u32;
  return RETURNok;
}

static inline int s6a_parse_service_selection(
    struct avp_hdr* hdr_service_selection, char* service_selection,
    int* length) {
  DevCheck(
      hdr_service_selection->avp_value->os.len <= SERVICE_SELECTION_MAX_LENGTH,
      hdr_service_selection->avp_value->os.len, ACCESS_POINT_NAME_MAX_LENGTH,
      0);
  *length = snprintf(service_selection, SERVICE_SELECTION_MAX_LENGTH, "%*s",
                     (int)hdr_service_selection->avp_value->os.len,
                     hdr_service_selection->avp_value->os.data);

  return RETURNok;
}

static inline int s6a_parse_qci(struct avp_hdr* hdr, qci_t* qci) {
  DevCheck(hdr->avp_value->u32 < QCI_MAX, hdr->avp_value->u32, QCI_MAX, 0);
  *qci = hdr->avp_value->u32;
  return RETURNok;
}

static inline int s6a_parse_priority_level(struct avp_hdr* hdr,
                                           priority_level_t* priority_level) {
  DevCheck(hdr->avp_value->u32 <= PRIORITY_LEVEL_MAX &&
               hdr->avp_value->u32 >= PRIORITY_LEVEL_MIN,
           hdr->avp_value->u32, PRIORITY_LEVEL_MAX, PRIORITY_LEVEL_MIN);
  *priority_level = (priority_level_t)hdr->avp_value->u32;
  return RETURNok;
}

static inline int s6a_parse_pre_emp_capability(
    struct avp_hdr* hdr, pre_emption_capability_t* pre_emp_capability) {
  DevCheck(hdr->avp_value->u32 < PRE_EMPTION_CAPABILITY_MAX,
           hdr->avp_value->u32, PRE_EMPTION_CAPABILITY_MAX, 0);
  *pre_emp_capability = (pre_emption_capability_t)hdr->avp_value->u32;
  return RETURNok;
}

static inline int s6a_parse_pre_emp_vulnerability(
    struct avp_hdr* hdr, pre_emption_vulnerability_t* pre_emp_vulnerability) {
  DevCheck(hdr->avp_value->u32 < PRE_EMPTION_VULNERABILITY_MAX,
           hdr->avp_value->u32, PRE_EMPTION_VULNERABILITY_MAX, 0);
  *pre_emp_vulnerability = (pre_emption_vulnerability_t)hdr->avp_value->u32;
  return RETURNok;
}

static inline int s6a_parse_allocation_retention_priority(
    struct avp* avp_arp, allocation_retention_priority_t* ptr) {
  struct avp* avp = NULL;
  struct avp_hdr* hdr;

  /*
   * If the Pre-emption-Capability AVP is not present in the
   * * * * Allocation-Retention-Priority AVP, the default value shall be
   * * * * PRE-EMPTION_CAPABILITY_DISABLED (1).
   */
  ptr->pre_emp_capability =
      (pre_emption_capability_t)PRE_EMPTION_CAPABILITY_DISABLED;
  /*
   * If the Pre-emption-Vulnerability AVP is not present in the
   * * * * Allocation-Retention-Priority AVP, the default value shall be
   * * * * PRE-EMPTION_VULNERABILITY_ENABLED (0).
   */
  ptr->pre_emp_vulnerability =
      (pre_emption_vulnerability_t)PRE_EMPTION_VULNERABILITY_ENABLED;
  CHECK_FCT(fd_msg_browse_internal(avp_arp, MSG_BRW_FIRST_CHILD,
                                   reinterpret_cast<msg_or_avp**>(&avp), NULL));

  while (avp) {
    CHECK_FCT(fd_msg_avp_hdr(avp, &hdr));

    switch (hdr->avp_code) {
      case AVP_CODE_PRIORITY_LEVEL:
        CHECK_FCT(s6a_parse_priority_level(hdr, &ptr->priority_level));
        break;

      case AVP_CODE_PRE_EMPTION_CAPABILITY:
        CHECK_FCT(s6a_parse_pre_emp_capability(hdr, &ptr->pre_emp_capability));
        break;

      case AVP_CODE_PRE_EMPTION_VULNERABILITY:
        CHECK_FCT(
            s6a_parse_pre_emp_vulnerability(hdr, &ptr->pre_emp_vulnerability));
        break;

      default:
        return RETURNerror;
    }

    /*
     * Go to next AVP in the grouped AVP
     */
    CHECK_FCT(fd_msg_browse_internal(
        avp, MSG_BRW_NEXT, reinterpret_cast<msg_or_avp**>(&avp), NULL));
  }

  return RETURNok;
}

static inline int s6a_parse_eps_subscribed_qos_profile(
    struct avp* avp_qos, eps_subscribed_qos_profile_t* ptr) {
  struct avp* avp = NULL;
  struct avp_hdr* hdr;

  CHECK_FCT(fd_msg_browse_internal(avp_qos, MSG_BRW_FIRST_CHILD,
                                   reinterpret_cast<msg_or_avp**>(&avp), NULL));

  while (avp) {
    CHECK_FCT(fd_msg_avp_hdr(avp, &hdr));

    switch (hdr->avp_code) {
      case AVP_CODE_QCI:
        CHECK_FCT(s6a_parse_qci(hdr, &ptr->qci));
        break;

      case AVP_CODE_ALLOCATION_RETENTION_PRIORITY:
        CHECK_FCT(s6a_parse_allocation_retention_priority(
            avp, &ptr->allocation_retention_priority));
        break;

      default:
        return RETURNerror;
    }

    /*
     * Go to next AVP in the grouped AVP
     */
    CHECK_FCT(fd_msg_browse_internal(
        avp, MSG_BRW_NEXT, reinterpret_cast<msg_or_avp**>(&avp), NULL));
  }

  return RETURNok;
}

static inline int s6a_parse_ip_address(struct avp_hdr* hdr,
                                       ip_address_t* ip_address) {
  uint16_t ip_type;

  DevCheck(hdr->avp_value->os.len >= 2, hdr->avp_value->os.len, 2, 0);
  ip_type = (hdr->avp_value->os.data[0] << 8) | (hdr->avp_value->os.data[1]);

  if (ip_type == IANA_IPV4) {
    /*
     * This is an IPv4 address
     */
    ip_address->pdn_type = IPv4;
    DevCheck(hdr->avp_value->os.len == 6, hdr->avp_value->os.len, 6, ip_type);
    uint32_t ip = (((uint32_t)hdr->avp_value->os.data[2]) << 24) |
                  (((uint32_t)hdr->avp_value->os.data[3]) << 16) |
                  (((uint32_t)hdr->avp_value->os.data[4]) << 8) |
                  ((uint32_t)hdr->avp_value->os.data[5]);

    ip_address->address.ipv4_address.s_addr = htonl(ip);
  } else if (ip_type == IANA_IPV6) {
    /*
     * This is an IPv6 address
     */
    ip_address->pdn_type = IPv6;
    DevCheck(hdr->avp_value->os.len == 18, hdr->avp_value->os.len, 18, ip_type);
    memcpy(ip_address->address.ipv6_address.__in6_u.__u6_addr8,
           &hdr->avp_value->os.data[2], 16);
  } else {
    /*
     * unhandled case...
     */
    return RETURNerror;
  }

  return RETURNok;
}

static inline int s6a_parse_apn_configuration(struct avp* avp_apn_conf_prof,
                                              apn_configuration_t* apn_config) {
  struct avp* avp = NULL;
  struct avp_hdr* hdr;

  CHECK_FCT(fd_msg_browse_internal(avp_apn_conf_prof, MSG_BRW_FIRST_CHILD,
                                   reinterpret_cast<msg_or_avp**>(&avp), NULL));
  memset(apn_config, 0, sizeof *apn_config);
  DevAssert(apn_config->nb_ip_address == 0);

  while (avp) {
    CHECK_FCT(fd_msg_avp_hdr(avp, &hdr));

    switch (hdr->avp_code) {
      case AVP_CODE_CONTEXT_IDENTIFIER:
        apn_config->context_identifier = hdr->avp_value->u32;
        break;

      case AVP_CODE_SERVED_PARTY_IP_ADDRESS:
        if (apn_config->nb_ip_address == 2) {
          DevMessage("Only two IP addresses can be provided");
        }

        CHECK_FCT(s6a_parse_ip_address(
            hdr, &apn_config->ip_address[apn_config->nb_ip_address]));
        apn_config->nb_ip_address++;
        break;

      case AVP_CODE_PDN_TYPE:
        CHECK_FCT(s6a_parse_pdn_type(hdr, &apn_config->pdn_type));
        break;

      case AVP_CODE_SERVICE_SELECTION:
        CHECK_FCT(
            s6a_parse_service_selection(hdr, apn_config->service_selection,
                                        &apn_config->service_selection_length));
        break;

      case AVP_CODE_EPS_SUBSCRIBED_QOS_PROFILE:
        CHECK_FCT(s6a_parse_eps_subscribed_qos_profile(
            avp, &apn_config->subscribed_qos));
        break;

      case AVP_CODE_AMBR:
        CHECK_FCT(s6a_parse_ambr(avp, &apn_config->ambr));
        break;
      case AVP_CODE_PDN_GW_ALLOCATION_TYPE:
        // Assuming PDN GW ALLOCATION TYPE is actually static
        break;

      case AVP_CODE_MIP6_AGENT_INFO:
        // TODO with AVP_CODE_PDN_GW_ALLOCATION_TYPE when splitting S and P-GW
        break;

      case AVP_CODE_3GPP_CHARGING_CHARACTERISTICS:
        OAILOG_INFO(LOG_S6A,
                    "AVP_CODE_3GPP_CHARGING_CHARACTERISTICS %d not processed\n",
                    hdr->avp_code);
        break;

      case AVP_CODE_VPLMN_DYNAMIC_ADDRESS_ALLOWED:
        OAILOG_INFO(LOG_S6A,
                    "AVP_CODE_VPLMN_DYNAMIC_ADDRESS_ALLOWED %d not processed\n",
                    hdr->avp_code);
        break;

      default:
        OAILOG_ERROR(LOG_S6A,
                     "Unknownn AVP code %d while parsing APN configuration\n",
                     hdr->avp_code);
    }

    /*
     * Go to next AVP in the grouped AVP
     */
    CHECK_FCT(fd_msg_browse_internal(
        avp, MSG_BRW_NEXT, reinterpret_cast<msg_or_avp**>(&avp), NULL));
  }

  return RETURNok;
}

static inline int s6a_parse_apn_configuration_profile(
    struct avp* avp_apn_conf_prof, apn_config_profile_t* apn_config_profile) {
  struct avp* avp = NULL;
  struct avp_hdr* hdr;

  CHECK_FCT(fd_msg_browse_internal(avp_apn_conf_prof, MSG_BRW_FIRST_CHILD,
                                   reinterpret_cast<msg_or_avp**>(&avp), NULL));

  apn_config_profile->nb_apns = 0;

  while (avp) {
    CHECK_FCT(fd_msg_avp_hdr(avp, &hdr));

    switch (hdr->avp_code) {
      case AVP_CODE_CONTEXT_IDENTIFIER:
        apn_config_profile->context_identifier = hdr->avp_value->u32;
        break;

      case AVP_CODE_ALL_APN_CONFIG_INC_IND:
        CHECK_FCT(s6a_parse_all_apn_conf_inc_ind(
            hdr, &apn_config_profile->all_apn_conf_ind));
        break;

      case AVP_CODE_APN_CONFIGURATION: {
        DevCheck(apn_config_profile->nb_apns < MAX_APN_PER_UE,
                 apn_config_profile->nb_apns, MAX_APN_PER_UE, 0);
        CHECK_FCT(s6a_parse_apn_configuration(
            avp, &apn_config_profile
                      ->apn_configuration[apn_config_profile->nb_apns]));
        apn_config_profile->nb_apns++;
      } break;
    }

    /*
     * Go to next AVP in the grouped AVP
     */
    CHECK_FCT(fd_msg_browse_internal(
        avp, MSG_BRW_NEXT, reinterpret_cast<msg_or_avp**>(&avp), NULL));
  }

  return RETURNok;
}

int s6a_parse_subscription_data(struct avp* avp_subscription_data,
                                subscription_data_t* subscription_data) {
  struct avp* avp = NULL;
  struct avp_hdr* hdr;

  CHECK_FCT(fd_msg_browse_internal(avp_subscription_data, MSG_BRW_FIRST_CHILD,
                                   reinterpret_cast<msg_or_avp**>(&avp), NULL));

  while (avp) {
    hdr = NULL;
    CHECK_FCT(fd_msg_avp_hdr(avp, &hdr));

    if (hdr) {
      switch (hdr->avp_code) {
        case AVP_CODE_SUBSCRIBER_STATUS:
          CHECK_FCT(s6a_parse_subscriber_status(
              hdr, &subscription_data->subscriber_status));
          break;

        case AVP_CODE_MSISDN:
          CHECK_FCT(s6a_parse_msisdn(hdr, subscription_data->msisdn,
                                     &subscription_data->msisdn_length));
          break;

        case AVP_CODE_NETWORK_ACCESS_MODE:
          CHECK_FCT(s6a_parse_network_access_mode(
              hdr, &subscription_data->access_mode));
          break;

        case AVP_CODE_ACCESS_RESTRICTION_DATA:
          CHECK_FCT(s6a_parse_access_restriction_data(
              hdr, &subscription_data->access_restriction));
          break;

        case AVP_CODE_AMBR:
          CHECK_FCT(s6a_parse_ambr(avp, &subscription_data->subscribed_ambr));
          break;

        case AVP_CODE_APN_CONFIGURATION_PROFILE:
          CHECK_FCT(s6a_parse_apn_configuration_profile(
              avp, &subscription_data->apn_config_profile));
          break;

        case AVP_CODE_SUBSCRIBED_PERIODIC_RAU_TAU_TIMER:
          subscription_data->rau_tau_timer = hdr->avp_value->u32;
          break;

        case AVP_CODE_APN_OI_REPLACEMENT:
          OAILOG_DEBUG(LOG_S6A,
                       "AVP code %d APN-OI-Replacement not processed\n",
                       hdr->avp_code);
          break;

        case AVP_CODE_3GPP_CHARGING_CHARACTERISTICS:
          OAILOG_DEBUG(
              LOG_S6A,
              "AVP code %d 3GPP-Charging Characteristics not processed\n",
              hdr->avp_code);
          break;

        case AVP_CODE_REGIONAL_SUBSCRIPTION_ZONE_CODE:
          OAILOG_DEBUG(
              LOG_S6A,
              "AVP code %d Regional-Subscription-Zone=Code not processed\n",
              hdr->avp_code);
          break;

        default:
          OAILOG_DEBUG(LOG_S6A, "Unknown AVP code %d not processed\n",
                       hdr->avp_code);
          return RETURNerror;
      }
    } else {
      OAILOG_DEBUG(LOG_S6A, "Subscription Data parsing Error\n");
      return RETURNerror;
    }

    /*
     * Go to next AVP in the grouped AVP
     */
    CHECK_FCT(fd_msg_browse_internal(
        avp, MSG_BRW_NEXT, reinterpret_cast<msg_or_avp**>(&avp), NULL));
  }
  return RETURNok;
}
