/*
 * Copyright (c) 2015, EURECOM (www.eurecom.fr)
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice,
 * this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 *
 * The views and conclusions contained in the software and documentation are
 * those of the authors and should not be interpreted as representing official
 * policies, either expressed or implied, of the FreeBSD Project.
 */

/*! \file 3gpp_29.274.h
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#ifndef FILE_3GPP_29_274_SEEN
#define FILE_3GPP_29_274_SEEN

#include <arpa/inet.h>
#include <stdint.h>

#include "3gpp_24.008.h"
#include "common_types.h"

//-------------------------------------
// 8.4 Cause

typedef enum gtpv2c_cause_value_e {
  /* Request / Initial message */
  LOCAL_DETACH                = 2,
  COMPLETE_DETACH             = 3,
  RAT_CHANGE_3GPP_TO_NON_3GPP = 4,  ///< RAT changed from 3GPP to Non-3GPP
  ISR_DEACTIVATION            = 5,
  ERROR_IND_FROM_RNC_ENB_SGSN = 6,
  IMSI_DETACH_ONLY            = 7,
  REACTIVATION_REQUESTED      = 8,
  PDN_RECONNECTION_TO_THIS_APN_DISALLOWED = 9,
  ACCESS_CHANGED_FROM_NON_3GPP_TO_3GPP    = 10,
  PDN_CONNECTION_INACTIVITY_TIMER_EXPIRES = 11,

  /* Acceptance in a Response/Triggered message */
  REQUEST_ACCEPTED           = 16,
  REQUEST_ACCEPTED_PARTIALLY = 17,
  NEW_PDN_TYPE_NW_PREF       = 18,  ///< New PDN type due to network preference
  NEW_PDN_TYPE_SAB_ONLY =
      19,  ///< New PDN type due to single address bearer only
  /* Rejection in a Response triggered message. */
  CONTEXT_NOT_FOUND                                                      = 64,
  INVALID_MESSAGE_FORMAT                                                 = 65,
  VERSION_NOT_SUPPORTED_BY_NEXT_PEER                                     = 66,
  INVALID_LENGTH                                                         = 67,
  SERVICE_NOT_SUPPORTED                                                  = 68,
  MANDATORY_IE_INCORRECT                                                 = 69,
  MANDATORY_IE_MISSING                                                   = 70,
  SYSTEM_FAILURE                                                         = 72,
  NO_RESOURCES_AVAILABLE                                                 = 73,
  SEMANTIC_ERROR_IN_TFT                                                  = 74,
  SYNTACTIC_ERROR_IN_TFT                                                 = 75,
  SEMANTIC_ERRORS_IN_PF                                                  = 76,
  SYNTACTIC_ERRORS_IN_PF                                                 = 77,
  MISSING_OR_UNKNOWN_APN                                                 = 78,
  GRE_KEY_NOT_FOUND                                                      = 80,
  RELOCATION_FAILURE                                                     = 81,
  DENIED_IN_RAT                                                          = 82,
  PREFERRED_PDN_TYPE_NOT_SUPPORTED                                       = 83,
  ALL_DYNAMIC_ADDRESSES_ARE_OCCUPIED                                     = 84,
  UE_CONTEXT_WITHOUT_TFT_ALREADY_ACTIVATED                               = 85,
  PROTOCOL_TYPE_NOT_SUPPORTED                                            = 86,
  UE_NOT_RESPONDING                                                      = 87,
  UE_REFUSES                                                             = 88,
  SERVICE_DENIED                                                         = 89,
  UNABLE_TO_PAGE_UE                                                      = 90,
  NO_MEMORY_AVAILABLE                                                    = 91,
  USER_AUTHENTICATION_FAILED                                             = 92,
  APN_ACCESS_DENIED_NO_SUBSCRIPTION                                      = 93,
  REQUEST_REJECTED                                                       = 94,
  P_TMSI_SIGNATURE_MISMATCH                                              = 95,
  IMSI_IMEI_NOT_KNOWN                                                    = 96,
  SEMANTIC_ERROR_IN_THE_TAD_OPERATION                                    = 97,
  SYNTACTIC_ERROR_IN_THE_TAD_OPERATION                                   = 98,
  REMOTE_PEER_NOT_RESPONDING                                             = 100,
  COLLISION_WITH_NETWORK_INITIATED_REQUEST                               = 101,
  UNABLE_TO_PAGE_UE_DUE_TO_SUSPENSION                                    = 102,
  CONDITIONAL_IE_MISSING                                                 = 103,
  APN_RESTRICTION_TYPE_INCOMPATIBLE_WITH_CURRENTLY_ACTIVE_PDN_CONNECTION = 104,
  INVALID_OVERALL_LENGTH_OF_THE_TRIGGERED_RESPONSE_MESSAGE_AND_A_PIGGYBACKED_INITIAL_MESSAGE =
      105,
  DATA_FORWARDING_NOT_SUPPORTED  = 106,
  INVALID_REPLY_FROM_REMOTE_PEER = 107,
  FALLBACK_TO_GTPV1              = 108,
  INVALID_PEER                   = 109,
  TEMP_REJECT_HO_IN_PROGRESS =
      110,  ///< Temporarily rejected due to handover procedure in progress
  MODIFICATIONS_NOT_LIMITED_TO_S1_U_BEARERS = 111,
  REJECTED_FOR_PMIPv6_REASON =
      112,  ///< Request rejected for a PMIPv6 reason (see 3GPP TS 29.275 [26]).
  APN_CONGESTION                = 113,
  BEARER_HANDLING_NOT_SUPPORTED = 114,
  UE_ALREADY_RE_ATTACHED        = 115,
  M_PDN_APN_NOT_ALLOWED =
      116,  ///< Multiple PDN connections for a given APN not allowed.
  LATE_OVERLAPPING_REQUEST =
      121,  ///< If the response message has not been received yet..
  SGW_CAUSE_MAX
} gtpv2c_cause_value_t;

typedef struct {
  gtpv2c_cause_value_t cause_value;
  uint8_t pce : 1;
  uint8_t bce : 1;
  uint8_t cs : 1;

  uint8_t offending_ie_type;
  uint16_t offending_ie_length;
  uint8_t offending_ie_instance;
} gtpv2c_cause_t;

//-------------------------------------
// 8.15 Bearer Quality of Service (Bearer QoS)
#define PRE_EMPTION_CAPABILITY_ENABLED (0x0)
#define PRE_EMPTION_CAPABILITY_DISABLED (0x1)
#define PRE_EMPTION_VULNERABILITY_ENABLED (0x0)
#define PRE_EMPTION_VULNERABILITY_DISABLED (0x1)

typedef struct bearer_qos_s {
  /* PCI (Pre-emption Capability)
   * The following values are defined:
   * - PRE-EMPTION_CAPABILITY_ENABLED (0)
   *    This value indicates that the service data flow or bearer is allowed
   *    to get resources that were already assigned to another service data
   *    flow or bearer with a lower priority level.
   * - PRE-EMPTION_CAPABILITY_DISABLED (1)
   *    This value indicates that the service data flow or bearer is not
   *    allowed to get resources that were already assigned to another service
   *    data flow or bearer with a lower priority level.
   * Default value: PRE-EMPTION_CAPABILITY_DISABLED
   */
  unsigned pci : 1;
  /* PL (Priority Level): defined in 3GPP TS.29.212 #5.3.45
   * Values 1 to 15 are defined, with value 1 as the highest level of priority.
   * Values 1 to 8 should only be assigned for services that are authorized to
   * receive prioritized treatment within an operator domain. Values 9 to 15
   * may be assigned to resources that are authorized by the home network and
   * thus applicable when a UE is roaming.
   */
  unsigned pl : 4;
  /* PVI (Pre-emption Vulnerability): defined in 3GPP TS.29.212 #5.3.47
   * Defines whether a service data flow can lose the resources assigned to it
   * in order to admit a service data flow with higher priority level.
   * The following values are defined:
   * - PRE-EMPTION_VULNERABILITY_ENABLED (0)
   *   This value indicates that the resources assigned to the service data
   *   flow or bearer can be pre-empted and allocated to a service data flow
   *   or bearer with a higher priority level.
   * - PRE-EMPTION_VULNERABILITY_DISABLED (1)
   *   This value indicates that the resources assigned to the service data
   *   flow or bearer shall not be pre-empted and allocated to a service data
   *   flow or bearer with a higher priority level.
   * Default value: EMPTION_VULNERABILITY_ENABLED
   */
  unsigned pvi : 1;
  uint8_t qci;
  ambr_t gbr;  ///< Guaranteed bit rate
  ambr_t mbr;  ///< Maximum bit rate
} bearer_qos_t;

//-------------------------------------
// 8.22 Fully Qualified TEID (F-TEID)

/* WARNING: not complete... */
typedef enum interface_type_e {
  INTERFACE_TYPE_MIN = 0,
  S1_U_ENODEB_GTP_U  = INTERFACE_TYPE_MIN,
  S1_U_SGW_GTP_U,
  S12_RNC_GTP_U,
  S12_SGW_GTP_U,
  S5_S8_SGW_GTP_U,
  S5_S8_PGW_GTP_U,
  S5_S8_SGW_GTP_C,
  S5_S8_PGW_GTP_C,
  S5_S8_SGW_PMIPv6,
  S5_S8_PGW_PMIPv6,
  S11_MME_GTP_C,
  S11_SGW_GTP_C,
  S10_MME_GTP_C,
  S3_MME_GTP_C,
  S3_SGSN_GTP_C,
  S4_SGSN_GTP_U,
  S4_SGW_GTP_U,
  S4_SGSN_GTP_C,
  S16_SGSN_GTP_C,
  ENODEB_GTP_U_DL_DATA_FORWARDING,
  ENODEB_GTP_U_UL_DATA_FORWARDING,
  RNC_GTP_U_DATA_FORWARDING,
  SGSN_GTP_U_DATA_FORWARDING,
  SGW_GTP_U_DL_DATA_FORWARDING,
  SM_MBMS_GW_GTP_C,
  SN_MBMS_GW_GTP_C,
  SM_MME_GTP_C,
  SN_SGSN_GTP_C,
  SGW_GTP_U_UL_DATA_FORWARDING,
  SN_SGSN_GTP_U,
  S2B_EPDG_GTP_C,
  INTERFACE_TYPE_MAX = S2B_EPDG_GTP_C
} interface_type_t;

typedef struct fteid_s {
  unsigned ipv4 : 1;
  unsigned ipv6 : 1;
  interface_type_t interface_type;
  teid_t teid;  ///< TEID or GRE Key
  struct in_addr ipv4_address;
  struct in6_addr ipv6_address;
} fteid_t;

//-------------------------------------
// 7.2.3-2: Bearer Context within Create Bearer Request

typedef struct bearer_context_within_create_bearer_request_s {
  uint8_t eps_bearer_id;  ///< EBI,  Mandatory CSR
  traffic_flow_template_t
      tft;  ///< Bearer TFT, Optional CSR, This IE may be included on the S4/S11
            ///< and S5/S8 interfaces.
  fteid_t s1u_sgw_fteid;  ///< This IE shall be sent on the S11 interface if the
                          ///< S1-U interface is used.
  ///< If SGW supports both IPv4 and IPv6, it shall send both an
  ///< IPv4 address and an IPv6 address within the S1-U SGW F-TEID IE.
  fteid_t s5_s8_u_pgw_fteid;  ///< This IE shall be sent on the S4, S5/S8 and
                              ///< S11 interfaces for GTP-based S5/S8 interface.
                              ///< The MME/SGSN shall
  ///< ignore the IE on S11/S4 for PMIP-based S5/S8 interface.
  fteid_t s12_sgw_fteid;   ///< This IE shall be sent on the S4 interface if the
                           ///< S12 interface is used. See NOTE 1.
  fteid_t s4_u_sgw_fteid;  ///< This IE shall be sent on the S4 interface if the
                           ///< S4-U interface is used. See NOTE 1.
  fteid_t s2b_u_pgw_fteid;  ///< This IE (for user plane) shall be sent on the
                            ///< S2b interface.
  bearer_qos_t bearer_level_qos;  ///<
  protocol_configuration_options_t
      pco;  ///< This IE may be sent on the S5/S8 and S4/S11 interfaces
  ///< if ePCO is not supported by the UE or the network. This bearer level IE
  ///< takes precedence over the PCO IE in the message body if they both exist.
} bearer_context_within_create_bearer_request_t;

//-------------------------------------
// 7.2.4-2: Bearer Context within Create Bearer Response

typedef struct bearer_context_within_create_bearer_response_s {
  uint8_t eps_bearer_id;  ///< EBI
  gtpv2c_cause_t
      cause;  ///< This IE shall indicate if the bearer handling was successful,
              ///< and if not, it gives information on the reason.
  fteid_t s1u_enb_fteid;  ///< This IE shall be sent on the S11 interface if the
                          ///< S1-U interface is used.
  fteid_t s1u_sgw_fteid;  ///< This IE shall be sent on the S11 interface. It
                          ///< shall be used to correlate the bearers with those
                          ///< in the Create Bearer Request.
  fteid_t
      s5_s8_u_sgw_fteid;  ///< This IE shall be sent on the S5/S8 interfaces.
  fteid_t s5_s8_u_pgw_fteid;  ///< This IE shall be sent on the S5/S8
                              ///< interfaces. It shall be
  ///< used to correlate the bearers with those in the Create
  ///< Bearer Request.
  fteid_t s12_rnc_fteid;  ///< C This IE shall be sent on the S4 interface if
                          ///< the S12 interface is used. See NOTE 1.
  fteid_t s12_sgw_fteid;  ///< C This IE shall be sent on the S4 interface. It
                          ///< shall be used to correlate the bearers with those
                          ///< in the Create Bearer Request. See NOTE1.
  fteid_t s4_u_sgsn_fteid;  ///< C This IE shall be sent on the S4 interface if
                            ///< the S4-U interface is used. See NOTE1.
  fteid_t s4_u_sgw_fteid;   ///< C This IE shall be sent on the S4 interface. It
                            ///< shall be used to correlate the bearers with
                            ///< those in the Create Bearer Request. See NOTE1.
  fteid_t s2b_u_epdg_fteid;  ///<  C This IE shall be sent on the S2b interface.
  fteid_t s2b_u_pgw_fteid;   ///<  C This IE shall be sent on the S2b interface.
                             ///<  It shall be used
  ///< to correlate the bearers with those in the Create Bearer
  ///<   Request.
  protocol_configuration_options_t
      pco;  ///< If the UE includes the PCO IE in the corresponding
            ///< message, then the MME/SGSN shall copy the content of
            ///< this IE transparently from the PCO IE included by the UE.
            ///< If the SGW receives PCO from MME/SGSN, SGW shall
            ///< forward it to the PGW. This bearer level IE takes
            ///< precedence over the PCO IE in the message body if they
            ///< both exist.
} bearer_context_within_create_bearer_response_t;

//-------------------------------------
// 8.16 Flow Quality of Service (Flow QoS)

typedef struct flow_qos_s {
  uint8_t qci;
  ambr_t gbr;  ///< Guaranteed bit rate
  ambr_t mbr;  ///< Maximum bit rate
} flow_qos_t;

#endif /* FILE_3GPP_29_274_SEEN */
