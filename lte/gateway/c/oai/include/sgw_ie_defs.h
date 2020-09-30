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

/*! \file sgw_ie_defs.h
 * \brief
 * \author Lionel Gauthier
 * \company Eurecom
 * \email: lionel.gauthier@eurecom.fr
 */

#ifndef FILE_SGW_IE_DEFS_SEEN
#define FILE_SGW_IE_DEFS_SEEN
#include "common_types.h"
#include "3gpp_23.003.h"
#include "3gpp_24.007.h"
#include "3gpp_24.008.h"
#include "3gpp_29.274.h"
#include "TrackingAreaIdentity.h"

typedef uint8_t DelayValue_t;
typedef uint32_t SequenceNumber_t;

/* Only one type of address can be present at the same time
 * This type is applicable to IP address Information Element defined
 * in 3GPP TS 29.274 #8.9
 */
typedef struct {
#define GTP_IP_ADDR_v4 0x0
#define GTP_IP_ADDR_v6 0x1
  unsigned present : 1;
  union {
    uint8_t v4[4];
    uint8_t v6[16];
  } address;
} gtp_ip_address_t;

/* 3GPP TS 29.274 Figure 8.12 */

typedef struct indication_flags_s {
  uint8_t daf : 1;
  uint8_t dtf : 1;
  uint8_t hi : 1;
  uint8_t dfi : 1;
  uint8_t oi : 1;
  uint8_t isrsi : 1;
  uint8_t israi : 1;
  uint8_t sgwci : 1;

  uint8_t sqci : 1;
  uint8_t uimsi : 1;
  uint8_t cfsi : 1;
  uint8_t crsi : 1;
  uint8_t p : 1;
  uint8_t pt : 1;
  uint8_t si : 1;
  uint8_t msv : 1;

  uint8_t spare1 : 1;
  uint8_t spare2 : 1;
  uint8_t spare3 : 1;
  uint8_t s6af : 1;
  uint8_t s4af : 1;
  uint8_t mbmdt : 1;
  uint8_t israu : 1;
  uint8_t ccrsi : 1;
} indication_flags_t;

/* Bit mask for octet 7 in indication IE */
// UPDATE RELEASE 10
#define S6AF_FLAG_BIT_POS 4
#define S4AF_FLAG_BIT_POS 3
#define MBMDT_FLAG_BIT_POS 2
#define ISRAU_FLAG_BIT_POS 1
#define CCRSI_FLAG_BIT_POS 0

/* Bit mask for octet 6 in indication IE */
#define SQSI_FLAG_BIT_POS 7
#define UIMSI_FLAG_BIT_POS 6
#define CFSI_FLAG_BIT_POS 5
#define CRSI_FLAG_BIT_POS 4
#define P_FLAG_BIT_POS 3
#define PT_FLAG_BIT_POS 2
#define SI_FLAG_BIT_POS 1
#define MSV_FLAG_BIT_POS 0

/* Bit mask for octet 5 in indication IE */
#define DAF_FLAG_BIT_POS 7
#define DTF_FLAG_BIT_POS 6
#define HI_FLAG_BIT_POS 5
#define DFI_FLAG_BIT_POS 4
#define OI_FLAG_BIT_POS 3
#define ISRSI_FLAG_BIT_POS 2
#define ISRAI_FLAG_BIT_POS 1
#define SGWCI_FLAG_BIT_POS 0

typedef struct {
  uint8_t digit[MSISDN_LENGTH];
  uint8_t length;
} Msisdn_t;

#define MEI_IMEI 0x0
#define MEI_IMEISV 0x1

typedef struct {
  uint8_t present;
  union {
    imei_t imei;
    imeisv_t imeisv;
  } choice;
} Mei_t;

typedef struct {
  uint8_t mcc[3];
  uint8_t mnc[3];
  uint16_t lac;
  uint16_t ci;
} Cgi_t;

typedef struct {
  uint8_t mcc[3];
  uint8_t mnc[3];
  uint16_t lac;
  uint16_t sac;
} Sai_t;

typedef struct {
  uint8_t mcc[3];
  uint8_t mnc[3];
  uint16_t lac;
  uint16_t rac;
} Rai_t;

typedef struct {
  uint8_t mcc[3];
  uint8_t mnc[3];
  uint16_t lac;
} Lai_t;

#define ULI_CGI 0x01
#define ULI_SAI 0x02
#define ULI_RAI 0x04
#define ULI_TAI 0x08
#define ULI_ECGI 0x10
#define ULI_LAI 0x20

typedef struct {
  uint8_t present;
  struct {
    Cgi_t cgi;
    Sai_t sai;
    Rai_t rai;
    tai_t tai;
    ecgi_t ecgi;
    Lai_t lai;
  } s;
} Uli_t;

typedef struct {
  uint8_t mcc[3];
  uint8_t mnc[3];
} ServingNetwork_t;

#define FTEID_T_2_IP_ADDRESS_T(fte_p, ip_p)                                    \
  do {                                                                         \
    if ((fte_p)->ipv4) {                                                       \
      (ip_p)->pdn_type                    = IPv4;                              \
      (ip_p)->address.ipv4_address.s_addr = (fte_p)->ipv4_address.s_addr;      \
    }                                                                          \
    if ((fte_p)->ipv6) {                                                       \
      if ((fte_p)->ipv4) {                                                     \
        (ip_p)->pdn_type = IPv4_AND_v6;                                        \
      } else {                                                                 \
        (ip_p)->pdn_type = IPv6;                                               \
      }                                                                        \
      memcpy(                                                                  \
          &(ip_p)->address.ipv6_address, &(fte_p)->ipv6_address,               \
          sizeof((fte_p)->ipv6_address));                                      \
    }                                                                          \
  } while (0)

typedef enum {
  TARGET_ID_RNC_ID       = 0,
  TARGET_ID_MACRO_ENB_ID = 1,
  TARGET_ID_CELL_ID      = 2,
  TARGET_ID_HOME_ENB_ID  = 3
  /* Other values are spare */
} target_type_t;

typedef struct {
  uint16_t lac;
  uint8_t rac;
  uint16_t id;
  uint16_t xid;
  /* Length of RNC Id can be 2 bytes if length of element is 8
   * or 4 bytes long if length is 10.
   */
  uint32_t rnc_id;
} rnc_id_t;

typedef struct {
  unsigned enb_id : 20;
  uint16_t tac;
} macro_enb_id_t;

typedef struct {
  unsigned enb_id : 28;
  uint16_t tac;
} home_enb_id_t;

typedef struct {
  /* Common part */
  uint8_t target_type;

  uint8_t mcc[3];
  uint8_t mnc[3];
  union {
    rnc_id_t rnc_id;
    macro_enb_id_t macro_enb_id;
    home_enb_id_t home_enb_id;
  } target_id;
} target_identification_t;

typedef enum SelectionMode_e {
  MS_O_N_P_APN_S_V = 0,  ///< MS or network provided APN, subscribed verified
  MS_P_APN_S_N_V   = 1,  ///< MS provided APN, subscription not verified
  N_P_APN_S_N_V    = 2,  ///< Network provided APN, subscription not verified
} SelectionMode_t;

typedef struct {
  uint32_t uplink_ambr;
  uint32_t downlink_ambr;
} AMBR_t;

typedef enum node_id_type_e {
  GLOBAL_UNICAST_IPv4 = 0,
  GLOBAL_UNICAST_IPv6 = 1,
  TYPE_EXOTIC = 2,  ///< (MCC * 1000 + MNC) << 12 + Integer value assigned to
                    ///< MME by operator
} node_id_type_t;

typedef struct {
  node_id_type_t node_id_type;
  uint16_t csid;  ///< Connection Set Identifier
  union {
    struct in_addr unicast_ipv4;
    struct in6_addr unicast_ipv6;
    struct {
      uint16_t mcc;
      uint16_t mnc;
      uint16_t operator_specific_id;
    } exotic;
  } node_id;
} FQ_CSID_t;

typedef struct {
  uint8_t time_zone;
  unsigned daylight_saving_time : 2;
} UETimeZone_t;

typedef enum AccessMode_e {
  CLOSED_MODE = 0,
  HYBRID_MODE = 1,
} AccessMode_t;

typedef struct {
  uint8_t mcc[3];
  uint8_t mnc[3];
  uint32_t csg_id;
  AccessMode_t access_mode;
  unsigned lcsg : 1;
  unsigned cmi : 1;
} UCI_t;

typedef struct {
  /* PPC (Prohibit Payload Compression):
   * This flag is used to determine whether an SGSN should attempt to
   * compress the payload of user data when the users asks for it
   * to be compressed (PPC = 0), or not (PPC = 1).
   */
  unsigned ppc : 1;

  /* VB (Voice Bearer):
   * This flag is used to indicate a voice bearer when doing PS-to-CS
   * SRVCC handover.
   */
  unsigned vb : 1;
} bearer_flags_t;

typedef enum node_type_e { NODE_TYPE_MME = 0, NODE_TYPE_SGSN = 1 } node_type_t;

typedef struct {
  uint8_t eps_bearer_id;  ///< EBI,  Mandatory CSR
  bearer_qos_t bearer_level_qos;
  traffic_flow_template_t
      tft;  ///< Bearer TFT, Optional CSR, This IE may be included on the S4/S11
            ///< and S5/S8 interfaces.
} bearer_to_create_t;

//-----------------
typedef struct bearer_context_to_be_created_s {
  uint8_t eps_bearer_id;  ///< EBI,  Mandatory CSR
  traffic_flow_template_t
      tft;  ///< Bearer TFT, Optional CSR, This IE may be included on the S4/S11
            ///< and S5/S8 interfaces.
  fteid_t s1u_enb_fteid;  ///< S1-U eNodeB F-TEID, Conditional CSR, This IE
                          ///< shall be included on the S11 interface for
                          ///< X2-based handover with SGW relocation.
  fteid_t s1u_sgw_fteid;  ///< S1-U SGW F-TEID, Conditional CSR, This IE shall
                          ///< be included on the S11 interface for X2-based
                          ///< handover with SGW relocation.fteid_t

  fteid_t s4u_sgsn_fteid;  ///< S4-U SGSN F-TEID, Conditional CSR, This IE shall
                           ///< be included on the S4 interface if the S4-U
                           ///< interface is used.
  fteid_t s5_s8_u_sgw_fteid;  ///< S5/S8-U SGW F-TEID, Conditional CSR, This IE
                              ///< shall be included on the S5/S8 interface for
                              ///< an "eUTRAN Initial Attach",
  ///  a "PDP Context Activation" or a "UE Requested PDN Connectivity".
  fteid_t s5_s8_u_pgw_fteid;  ///< S5/S8-U PGW F-TEID, Conditional CSR, This IE
                              ///< shall be included on the S4 and S11
                              ///< interfaces for the TAU/RAU/Handover
                              /// cases when the GTP-based S5/S8 is used.
  fteid_t s12_rnc_fteid;  ///< S12 RNC F-TEID, Conditional Optional CSR, This IE
                          ///< shall be included on the S4 interface if the S12
  /// interface is used in the Enhanced serving RNS relocation with SGW
  /// relocation procedure.
  fteid_t
      s2b_u_epdg_fteid;  ///< S2b-U ePDG F-TEID, Conditional CSR, This IE shall
                         ///< be included on the S2b interface for an Attach
  /// with GTP on S2b, a UE initiated Connectivity to Additional PDN with GTP on
  /// S2b and a Handover to Untrusted Non- 3GPP IP Access with GTP on S2b.
  /* This parameter is received only if the QoS parameters have been modified */
  bearer_qos_t bearer_level_qos;  ///< Bearer QoS, Mandatory CSR
  protocol_configuration_options_t
      pco;  ///< This IE may be sent on the S5/S8 and S4/S11 interfaces
            ///< if ePCO is not supported by the UE or the network. This bearer
            ///< level IE takes precedence over the PCO IE in the message body
            ///< if they both exist.
  gtpv2c_cause_t cause;
} bearer_context_to_be_created_t;

typedef struct bearer_contexts_to_be_created_s {
#define MSG_CREATE_SESSION_REQUEST_MAX_BEARER_CONTEXTS 11
  uint8_t num_bearer_context;
  bearer_context_to_be_created_t bearer_contexts
      [MSG_CREATE_SESSION_REQUEST_MAX_BEARER_CONTEXTS];  ///< Bearer Contexts to
                                                         ///< be created
  ///< Several IEs with the same type and instance value shall be
  ///< included on the S4/S11 and S5/S8 interfaces as necessary
  ///< to represent a list of Bearers. One single IE shall be
  ///< included on the S2b interface.
  ///< One bearer shall be included for an E-UTRAN Initial
  ///< Attach, a PDP Context Activation, a UE requested PDN
  ///< Connectivity, an Attach with GTP on S2b, a UE initiated
  ///< Connectivity to Additional PDN with GTP on S2b and a
  ///< Handover to Untrusted Non-3GPP IP Access with GTP on
  ///< S2b.
  ///< One or more bearers shall be included for a
  ///< Handover/TAU/RAU with an SGW change.
} bearer_contexts_to_be_created_t;

//-----------------
typedef struct bearer_context_created_s {
  uint8_t eps_bearer_id;  ///< EPS Bearer ID
  gtpv2c_cause_t cause;

  /* This parameter is used on S11 interface only */
  fteid_t s1u_sgw_fteid;  ///< S1-U SGW F-TEID

  /* This parameter is used on S4 interface only */
  fteid_t s4u_sgw_fteid;  ///< S4-U SGW F-TEID

  /* This parameter is used on S11 and S5/S8 interface only for a
   * GTP-based S5/S8 interface and during:
   * - E-UTRAN Inintial attch
   * - PDP Context Activation
   * - UE requested PDN connectivity
   */
  fteid_t s5_s8_u_pgw_fteid;  ///< S4-U SGW F-TEID

  /* This parameter is used on S4 interface only and when S12 interface is used
   */
  fteid_t s12_sgw_fteid;  ///< S12 SGW F-TEID

  /* This parameter is received only if the QoS parameters have been modified */
  bearer_qos_t* bearer_level_qos;

  traffic_flow_template_t tft;  ///< Bearer TFT
} bearer_context_created_t;

typedef struct bearer_contexts_created_s {
  uint8_t num_bearer_context;
  bearer_context_created_t
      bearer_contexts[MSG_CREATE_SESSION_REQUEST_MAX_BEARER_CONTEXTS];
} bearer_contexts_created_t;

//-----------------
typedef struct bearer_context_to_be_updated_s {
  uint8_t eps_bearer_id;  ///< EBI,  Mandatory CSR
  traffic_flow_template_t*
      tft;  ///< Bearer TFT, Optional CSR, This IE may be included on the S4/S11
            ///< and S5/S8 interfaces.
  /* This parameter is received only if the QoS parameters have been modified */
  bearer_qos_t* bearer_level_qos;  ///< Bearer QoS, Mandatory CSR
  protocol_configuration_options_t
      pco;  ///< This IE may be sent on the S5/S8 and S4/S11 interfaces
            ///< if ePCO is not supported by the UE or the network. This bearer
            ///< level IE takes precedence over the PCO IE in the message body
            ///< if they both exist.
  gtpv2c_cause_t cause;
} bearer_context_to_be_updated_t;

typedef struct bearer_contexts_to_be_updated_s {
#define MSG_UPDATE_BEARER_REQUEST_MAX_BEARER_CONTEXTS 11
  uint8_t num_bearer_context;
  bearer_context_to_be_updated_t bearer_context
      [MSG_UPDATE_BEARER_REQUEST_MAX_BEARER_CONTEXTS];  ///< Bearer Contexts to
                                                        ///< be created
  ///< Several IEs with the same type and instance value shall be
  ///< included on the S4/S11 and S5/S8 interfaces as necessary
  ///< to represent a list of Bearers. One single IE shall be
  ///< included on the S2b interface.
  ///< One bearer shall be included for an E-UTRAN Initial
  ///< Attach, a PDP Context Activation, a UE requested PDN
  ///< Connectivity, an Attach with GTP on S2b, a UE initiated
  ///< Connectivity to Additional PDN with GTP on S2b and a
  ///< Handover to Untrusted Non-3GPP IP Access with GTP on
  ///< S2b.
  ///< One or more bearers shall be included for a
  ///< Handover/TAU/RAU with an SGW change.
} bearer_contexts_to_be_updated_t;
//-----------------
typedef struct bearer_context_modified_s {
  uint8_t eps_bearer_id;  ///< EPS Bearer ID
  gtpv2c_cause_t cause;
  fteid_t s1u_sgw_fteid;  ///< Sender F-TEID for user plane
} bearer_context_modified_t;

typedef struct bearer_contexts_modified_s {
#define MSG_MODIFY_BEARER_RESPONSE_MAX_BEARER_CONTEXTS 11
  uint8_t num_bearer_context;
  bearer_context_modified_t
      bearer_contexts[MSG_MODIFY_BEARER_RESPONSE_MAX_BEARER_CONTEXTS];
} bearer_contexts_modified_t;

//-----------------
typedef struct bearer_context_marked_for_removal_s {
  uint8_t eps_bearer_id;  ///< EPS bearer ID
  gtpv2c_cause_t cause;
} bearer_context_marked_for_removal_t;

typedef struct bearer_contexts_marked_for_removal_s {
  uint8_t num_bearer_context;
  bearer_context_marked_for_removal_t
      bearer_contexts[MSG_MODIFY_BEARER_RESPONSE_MAX_BEARER_CONTEXTS];
} bearer_contexts_marked_for_removal_t;

//-----------------
typedef struct bearer_context_to_be_modified_s {
  uint8_t eps_bearer_id;  ///< EPS Bearer ID
  fteid_t s1_eNB_fteid;   ///< S1 eNodeB F-TEID
} bearer_context_to_be_modified_t;

typedef struct bearer_contexts_to_be_modified_s {
#define MSG_MODIFY_BEARER_REQUEST_MAX_BEARER_CONTEXTS 11
  uint8_t num_bearer_context;
  bearer_context_to_be_modified_t
      bearer_contexts[MSG_MODIFY_BEARER_REQUEST_MAX_BEARER_CONTEXTS];
} bearer_contexts_to_be_modified_t;
//-----------------

typedef struct bearer_context_to_be_removed_s {
  uint8_t eps_bearer_id;   ///< EPS Bearer ID, Mandatory
  fteid_t s4u_sgsn_fteid;  ///< S4-U SGSN F-TEID, Conditional , redundant
  gtpv2c_cause_t cause;
} bearer_context_to_be_removed_t;  // Within Create Session Request, Modify
                                   // Bearer Request, Modify Access Bearers
                                   // Request

typedef struct bearer_contexts_to_be_removed_s {
  uint8_t num_bearer_context;
  bearer_context_to_be_removed_t
      bearer_contexts[MSG_CREATE_SESSION_REQUEST_MAX_BEARER_CONTEXTS];
} bearer_contexts_to_be_removed_t;

typedef struct ebi_list_s {
  uint32_t num_ebi;
#define RELEASE_ACCESS_BEARER_MAX_BEARERS 8
  ebi_t ebis[RELEASE_ACCESS_BEARER_MAX_BEARERS];
} ebi_list_t;

//-------------------------------------
// 7.2.16 Update Bearer Response
//-------------------------------------
// 7.2.16-2: Bearer Context within Update Bearer Response

typedef struct bearer_context_within_update_bearer_response_s {
  uint8_t eps_bearer_id;  ///< EBI
  gtpv2c_cause_t
      cause;  ///< This IE shall indicate if the bearer handling was successful,
              ///< and if not, it gives information on the reason.
  fteid_t s12_rnc_fteid;    ///< C This IE shall be sent on the S4 interface if
                            ///< the S12 interface is used. See NOTE 1.
  fteid_t s4_u_sgsn_fteid;  ///< C This IE shall be sent on the S4 interface if
                            ///< the S4-U interface is used. See NOTE1.
  protocol_configuration_options_t
      pco;  ///< If the UE includes the PCO IE in the corresponding
            ///< message, then the MME/SGSN shall copy the content of
            ///< this IE transparently from the PCO IE included by the UE.
            ///< If the SGW receives PCO from MME/SGSN, SGW shall
            ///< forward it to the PGW. This bearer level IE takes
            ///< precedence over the PCO IE in the message body if they
            ///< both exist.
} bearer_context_within_update_bearer_response_t;

typedef struct bearer_contexts_within_update_bearer_response_s {
#define MSG_UPDATE_BEARER_RESPONSE_MAX_BEARER_CONTEXTS 11
  uint8_t num_bearer_context;
  bearer_context_within_update_bearer_response_t
      bearer_context[MSG_UPDATE_BEARER_RESPONSE_MAX_BEARER_CONTEXTS];
} bearer_contexts_within_update_bearer_response_t;

//-------------------------------------
// 7.2.10-2: Bearer Context within Delete Bearer Response

typedef struct bearer_context_within_delete_bearer_response_s {
  uint8_t eps_bearer_id;  ///< EBI
  gtpv2c_cause_t
      cause;  ///< This IE shall indicate if the bearer handling was successful,
              ///< and if not, it gives information on the reason.
  protocol_configuration_options_t
      pco;  ///< If the UE includes the PCO IE in the corresponding
            ///< message, then the MME/SGSN shall copy the content of
            ///< this IE transparently from the PCO IE included by the UE.
            ///< If the SGW receives PCO from MME/SGSN, SGW shall
            ///< forward it to the PGW. This bearer level IE takes
            ///< precedence over the PCO IE in the message body if they
} bearer_context_within_delete_bearer_response_t;

#define MSG_DELETE_BEARER_REQUEST_MAX_FAILED_BEARER_CONTEXTS                   \
  11  // todo: find optimum number

typedef struct bearer_contexts_within_delete_bearer_response_s {
#define MSG_DELETE_BEARER_RESPONSE_MAX_BEARER_CONTEXTS 11
  uint8_t num_bearer_context;
  bearer_context_within_delete_bearer_response_t
      bearer_context[MSG_DELETE_BEARER_RESPONSE_MAX_BEARER_CONTEXTS];
} bearer_contexts_within_delete_bearer_response_t;

//-----------------

typedef struct bearer_contexts_within_create_bearer_request_s {
#define MSG_CREATE_BEARER_REQUEST_MAX_BEARER_CONTEXTS 11
  uint8_t num_bearer_context;
  bearer_context_within_create_bearer_request_t
      bearer_contexts[MSG_CREATE_BEARER_REQUEST_MAX_BEARER_CONTEXTS];
} bearer_contexts_within_create_bearer_request_t;

typedef struct bearer_contexts_within_create_bearer_response_s {
#define MSG_CREATE_BEARER_RESPONSE_MAX_BEARER_CONTEXTS 11
  uint8_t num_bearer_context;
  bearer_context_within_create_bearer_response_t
      bearer_contexts[MSG_CREATE_BEARER_RESPONSE_MAX_BEARER_CONTEXTS];
} bearer_contexts_within_create_bearer_response_t;

#endif /* FILE_SGW_IE_DEFS_SEEN */
