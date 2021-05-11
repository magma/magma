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
/*! \file s11_messages_types.h
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#ifndef FILE_S11_MESSAGES_TYPES_SEEN
#define FILE_S11_MESSAGES_TYPES_SEEN

#include "sgw_ie_defs.h"

#define S11_CREATE_SESSION_REQUEST(mSGpTR)                                     \
  (mSGpTR)->ittiMsg.s11_create_session_request
#define S11_CREATE_SESSION_RESPONSE(mSGpTR)                                    \
  (mSGpTR)->ittiMsg.s11_create_session_response
#define S11_CREATE_BEARER_REQUEST(mSGpTR)                                      \
  (mSGpTR)->ittiMsg.s11_create_bearer_request
#define S11_CREATE_BEARER_RESPONSE(mSGpTR)                                     \
  (mSGpTR)->ittiMsg.s11_create_bearer_response
#define S11_MODIFY_BEARER_REQUEST(mSGpTR)                                      \
  (mSGpTR)->ittiMsg.s11_modify_bearer_request
#define S11_MODIFY_BEARER_RESPONSE(mSGpTR)                                     \
  (mSGpTR)->ittiMsg.s11_modify_bearer_response
#define S11_DELETE_SESSION_REQUEST(mSGpTR)                                     \
  (mSGpTR)->ittiMsg.s11_delete_session_request
#define S11_DELETE_BEARER_COMMAND(mSGpTR)                                      \
  (mSGpTR)->ittiMsg.s11_delete_bearer_command
#define S11_DELETE_SESSION_RESPONSE(mSGpTR)                                    \
  (mSGpTR)->ittiMsg.s11_delete_session_response
#define S11_RELEASE_ACCESS_BEARERS_REQUEST(mSGpTR)                             \
  (mSGpTR)->ittiMsg.s11_release_access_bearers_request
#define S11_RELEASE_ACCESS_BEARERS_RESPONSE(mSGpTR)                            \
  (mSGpTR)->ittiMsg.s11_release_access_bearers_response
#define S11_PAGING_REQUEST(mSGpTR) (mSGpTR)->ittiMsg.s11_paging_request
#define S11_PAGING_RESPONSE(mSGpTR) (mSGpTR)->ittiMsg.s11_paging_response
#define S11_SUSPEND_NOTIFICATION(mSGpTR)                                       \
  (mSGpTR)->ittiMsg.s11_suspend_notification
#define S11_SUSPEND_ACKNOWLEDGE(mSGpTR)                                        \
  (mSGpTR)->ittiMsg.s11_suspend_acknowledge
#define S11_MODIFY_UE_AMBR_REQUEST(mSGpTR)                                     \
  (mSGpTR)->ittiMsg.s11_modify_ue_ambr_request
#define S11_NW_INITIATED_ACTIVATE_BEARER_REQUEST(mSGpTR)                       \
  (mSGpTR)->ittiMsg.s11_nw_init_actv_bearer_request
#define S11_NW_INITIATED_ACTIVATE_BEARER_RESP(mSGpTR)                          \
  (mSGpTR)->ittiMsg.s11_nw_init_actv_bearer_rsp
#define S11_NW_INITIATED_DEACTIVATE_BEARER_REQUEST(mSGpTR)                     \
  (mSGpTR)->ittiMsg.s11_nw_init_deactv_bearer_request
#define S11_NW_INITIATED_DEACTIVATE_BEARER_RESP(mSGpTR)                        \
  (mSGpTR)->ittiMsg.s11_nw_init_deactv_bearer_rsp
#define S11_DOWNLINK_DATA_NOTIFICATION_ACKNOWLEDGE(mSGpTR)                     \
  (mSGpTR)->ittiMsg.s11_downlink_data_notification_acknowledge
//-----------------------------------------------------------------------------
/** @struct itti_s11_nw_initiated_ded_bearer_actv_request_t
 *  @brief PCRF initiated Dedicated Bearer Activation Request
 */
typedef struct itti_s11_nw_init_actv_bearer_request_s {
  teid_t context_teid;                   ///< not in specs for inner use;
  teid_t s11_mme_teid;                   ///< MME TEID
  bearer_qos_t eps_bearer_qos;           ///< Bearer QoS
  traffic_flow_template_t tft;           ///< Traffic Flow Template
  protocol_configuration_options_t pco;  ///< PCO protocol_configuration_options
  fteid_t s1_u_sgw_fteid;                /// S1U SGW FTEID
  ebi_t lbi;                             // Linked Bearer ID
} itti_s11_nw_init_actv_bearer_request_t;

//-----------------------------------------------------------------------------
/** @struct itti_s11_nw_initiated_ded_bearer_actv_rsp_t
 *  @brief PCRF initiated Dedicated Bearer Activation Rsp
 */
typedef struct itti_s11_nw_init_actv_bearer_rsp_s {
  gtpv2c_cause_t cause;         ///< M
  teid_t sgw_s11_teid;          // CP TEID
  bearer_qos_t eps_bearer_qos;  ///< Bearer QoS to be stored in SGW context
  traffic_flow_template_t tft;  ///< TFT to be stored in SGW context
  protocol_configuration_options_t pco;  ///< PCO protocol_configuration_options
  bearer_contexts_within_create_bearer_response_t
      bearer_contexts;  ///< Several IEs with this type and instance value shall
                        ///< be
} itti_s11_nw_init_actv_bearer_rsp_t;

//-----------------------------------------------------------------------------
/** @struct itti_s11_nw_init_deactv_bearer_request_t
 *  @brief PCRF initiated Dedicated Bearer Deactivation Request
 */
typedef struct itti_s11_nw_init_deactv_bearer_request_s {
  uint8_t no_of_bearers;
  ebi_t ebi[BEARERS_PER_UE];  // EPS Bearer ID
  teid_t s11_mme_teid;        ///< MME TEID
  bool delete_default_bearer;
} itti_s11_nw_init_deactv_bearer_request_t;

//-----------------------------------------------------------------------------
/** @struct itti_s11_nw_init_actv_bearer_rsp_t
 *  @brief PCRF initiated Dedicated Bearer dectivation Rsp
 */
typedef struct itti_s11_nw_init_deactv_bearer_rsp_s {
  imsi64_t imsi;
  gtpv2c_cause_t cause;  ///< M
  ebi_t* lbi;            // Default EPS Bearer ID
  bearer_contexts_within_create_bearer_response_t
      bearer_contexts;  ///< Several IEs with this type and instance value shall
                        ///< be
  bool delete_default_bearer;
  teid_t s_gw_teid_s11_s4;
} itti_s11_nw_init_deactv_bearer_rsp_t;

//-------------
/** typedef edns peer ip instead of defining it multiple times for each message.
 */
typedef union {
  struct sockaddr_in
      addr_v4;  ///< MME ipv4 address for S-GW or S-GW ipv4 address for MME
  struct sockaddr_in6
      addr_v6;  ///< MME ipv6 address for S-GW or S-GW ipv6 address for MME
} edns_peer_ip_t;

//-----------------------------------------------------------------------------
/** @struct itti_s11_create_session_request_t
 *  @brief Create Session Request
 *
 * Spec 3GPP TS 29.274, Universal Mobile Telecommunications System (UMTS);
 *                      LTE; 3GPP Evolved Packet System (EPS);
 *                      Evolved General Packet Radio Service (GPRS);
 *                      Tunnelling Protocol for Control plane (GTPv2-C); Stage 3
 * The Create Session Request will be sent on S11 interface as
 * part of these procedures:
 * - E-UTRAN Initial Attach
 * - UE requested PDN connectivity
 * - Tracking Area Update procedure with Serving GW change
 * - S1/X2-based handover with SGW change
 */
typedef struct itti_s11_create_session_request_s {
  teid_t teid;  ///< S11- S-GW Tunnel Endpoint Identifier

  Imsi_t imsi;  ///< The IMSI shall be included in the message on the S4/S11
  ///< interface, and on S5/S8 interface if provided by the
  ///< MME/SGSN, except for the case:
  ///<     - If the UE is emergency attached and the UE is UICCless.
  ///< The IMSI shall be included in the message on the S4/S11
  ///< interface, and on S5/S8 interface if provided by the
  ///< MME/SGSN, but not used as an identifier
  ///<     - if UE is emergency attached but IMSI is not authenticated.
  ///< The IMSI shall be included in the message on the S2b interface.

  Msisdn_t msisdn;  ///< For an E-UTRAN Initial Attach the IE shall be included
  ///< when used on the S11 interface, if provided in the
  ///< subscription data from the HSS.
  ///< For a PDP Context Activation procedure the IE shall be
  ///< included when used on the S4 interface, if provided in the
  ///< subscription data from the HSS.
  ///< The IE shall be included for the case of a UE Requested
  ///< PDN Connectivity, if the MME has it stored for that UE.
  ///< It shall be included when used on the S5/S8 interfaces if
  ///< provided by the MME/SGSN.
  ///< The ePDG shall include this IE on the S2b interface during
  ///< an Attach with GTP on S2b and a UE initiated Connectivity
  ///< to Additional PDN with GTP on S2b, if provided by the
  ///< HSS/AAA.

  Mei_t mei;  ///< The MME/SGSN shall include the ME Identity (MEI) IE on
  ///< the S11/S4 interface:
  ///<     - If the UE is emergency attached and the UE is UICCless
  ///<     - If the UE is emergency attached and the IMSI is not authenticated
  ///< For all other cases the MME/SGSN shall include the ME
  ///< Identity (MEI) IE on the S11/S4 interface if it is available.
  ///< If the SGW receives this IE, it shall forward it to the PGW
  ///< on the S5/S8 interface.

  Uli_t uli;  ///< This IE shall be included on the S11 interface for E-
  ///< UTRAN Initial Attach and UE-requested PDN Connectivity
  ///< procedures. It shall include ECGI&TAI. The MME/SGSN
  ///< shall also include it on the S11/S4 interface for
  ///< TAU/RAU/X2-Handover/Enhanced SRNS Relocation
  ///< procedure if the PGW has requested location information
  ///< change reporting and MME/SGSN support location
  ///< information change reporting. The SGW shall include this
  ///< IE on S5/S8 if it receives the ULI from MME/SGSN.

  ServingNetwork_t serving_network;  ///< This IE shall be included on the
                                     ///< S4/S11, S5/S8 and S2b
  ///< interfaces for an E-UTRAN initial attach, a PDP Context
  ///< Activation, a UE requested PDN connectivity, an Attach
  ///< with GTP on S2b, a UE initiated Connectivity to Additional
  ///< PDN with GTP on S2b and a Handover to Untrusted Non-
  ///< 3GPP IP Access with GTP on S2b.

  rat_type_t
      rat_type;  ///< This IE shall be set to the 3GPP access type or to the
  ///< value matching the characteristics of the non-3GPP access
  ///< the UE is using to attach to the EPS.
  ///< The ePDG may use the access technology type of the
  ///< untrusted non-3GPP access network if it is able to acquire
  ///< it; otherwise it shall indicate Virtual as the RAT Type.
  ///< See NOTE 3, NOTE 4.

  indication_flags_t indication_flags;  ///< This IE shall be included if any
                                        ///< one of the applicable flags
  ///< is set to 1.
  ///< Applicable flags are:
  ///<     - S5/S8 Protocol Type: This flag shall be used on
  ///<       the S11/S4 interfaces and set according to the
  ///<       protocol chosen to be used on the S5/S8
  ///<       interfaces.
  ///<
  ///<     - Dual Address Bearer Flag: This flag shall be used
  ///<       on the S2b, S11/S4 and S5/S8 interfaces and shall
  ///<       be set to 1 when the PDN Type, determined based
  ///<       on UE request and subscription record, is set to
  ///<       IPv4v6 and all SGSNs which the UE may be
  ///<       handed over to support dual addressing. This shall
  ///<       be determined based on node pre-configuration by
  ///<       the operator.
  ///<
  ///<     - Handover Indication: This flag shall be set to 1 on
  ///<       the S11/S4 and S5/S8 interface during an E-
  ///<       UTRAN Initial Attach or a UE Requested PDN
  ///<       Connectivity or aPDP Context Activation procedure
  ///<       if the PDN connection/PDP Context is handed-over
  ///<       from non-3GPP access.
  ///<       This flag shall be set to 1 on the S2b interface
  ///<       during a Handover to Untrusted Non-3GPP IP
  ///<       Access with GTP on S2b and IP address
  ///<       preservation is requested by the UE.
  ///<
  ///<       ....
  ///<     - Unauthenticated IMSI: This flag shall be set to 1
  ///<       on the S4/S11 and S5/S8 interfaces if the IMSI
  ///<       present in the message is not authenticated and is
  ///<       for an emergency attached UE.

  fteid_t sender_fteid_for_cp;  ///< Sender F-TEID for control plane (MME)

  fteid_t pgw_address_for_cp;  ///< PGW S5/S8 address for control plane or PMIP
  ///< This IE shall be sent on the S11 / S4 interfaces. The TEID
  ///< or GRE Key is set to "0" in the E-UTRAN initial attach, the
  ///< PDP Context Activation and the UE requested PDN
  ///< connectivity procedures.

  char apn[ACCESS_POINT_NAME_MAX_LENGTH + 1];  ///< Access Point Name

  SelectionMode_t selection_mode;  ///< Selection Mode
  ///< This IE shall be included on the S4/S11 and S5/S8
  ///< interfaces for an E-UTRAN initial attach, a PDP Context
  ///< Activation and a UE requested PDN connectivity.
  ///< This IE shall be included on the S2b interface for an Initial
  ///< Attach with GTP on S2b and a UE initiated Connectivity to
  ///< Additional PDN with GTP on S2b.
  ///< It shall indicate whether a subscribed APN or a non
  ///< subscribed APN chosen by the MME/SGSN/ePDG was
  ///< selected.
  ///< CO: When available, this IE shall be sent by the MME/SGSN on
  ///< the S11/S4 interface during TAU/RAU/HO with SGW
  ///< relocation.

  pdn_type_t pdn_type;  ///< PDN Type
  ///< This IE shall be included on the S4/S11 and S5/S8
  ///< interfaces for an E-UTRAN initial attach, a PDP Context
  ///< Activation and a UE requested PDN connectivity.
  ///< This IE shall be set to IPv4, IPv6 or IPv4v6. This is based
  ///< on the UE request and the subscription record retrieved
  ///< from the HSS (for MME see 3GPP TS 23.401 [3], clause
  ///< 5.3.1.1, and for SGSN see 3GPP TS 23.060 [35], clause
  ///< 9.2.1). See NOTE 1.

  paa_t paa;  ///< PDN Address Allocation
  ///< This IE shall be included the S4/S11, S5/S8 and S2b
  ///< interfaces for an E-UTRAN initial attach, a PDP Context
  ///< Activation, a UE requested PDN connectivity, an Attach
  ///< with GTP on S2b, a UE initiated Connectivity to Additional
  ///< PDN with GTP on S2b and a Handover to Untrusted Non-
  ///< 3GPP IP Access with GTP on S2b. For PMIP-based
  ///< S5/S8, this IE shall also be included on the S4/S11
  ///< interfaces for TAU/RAU/Handover cases involving SGW
  ///< relocation.
  ///< The PDN type field in the PAA shall be set to IPv4, or IPv6
  ///< or IPv4v6 by MME, based on the UE request and the
  ///< subscription record retrieved from the HSS.
  ///< For static IP address assignment (for MME see 3GPP TS
  ///< 23.401 [3], clause 5.3.1.1, for SGSN see 3GPP TS 23.060
  ///< [35], clause 9.2.1, and for ePDG see 3GPP TS 23.402 [45]
  ///< subclause 4.7.3), the MME/SGSN/ePDG shall set the IPv4
  ///< address and/or IPv6 prefix length and IPv6 prefix and
  ///< Interface Identifier based on the subscribed values
  ///< received from HSS, if available. The value of PDN Type
  ///< field shall be consistent with the value of the PDN Type IE,
  ///< if present in this message.
  ///< For a Handover to Untrusted Non-3GPP IP Access with
  ///< GTP on S2b, the ePDG shall set the IPv4 address and/or
  ///< IPv6 prefix length and IPv6 prefix and Interface Identifier
  ///< based on the IP address(es) received from the UE.
  ///< If static IP address assignment is not used, and for
  ///< scenarios other than a Handover to Untrusted Non-3GPP
  ///< IP Access with GTP on S2b, the IPv4 address shall be set
  ///< to 0.0.0.0, and/or the IPv6 Prefix Length and IPv6 prefix
  ///< and Interface Identifier shall all be set to zero.
  ///<
  ///< CO: This IE shall be sent by the MME/SGSN on S11/S4
  ///< interface during TAU/RAU/HO with SGW relocation.

  // APN Restriction Maximum_APN_Restriction ///< This IE shall be included on
  // the S4/S11 and S5/S8
  ///< interfaces in the E-UTRAN initial attach, PDP Context
  ///< Activation and UE Requested PDN connectivity
  ///< procedures.
  ///< This IE denotes the most stringent restriction as required
  ///< by any already active bearer context. If there are no
  ///< already active bearer contexts, this value is set to the least
  ///< restrictive type.

  ebi_t default_ebi;

  ambr_t ambr;  ///< Aggregate Maximum Bit Rate (APN-AMBR)
  ///< This IE represents the APN-AMBR. It shall be included on
  ///< the S4/S11, S5/S8 and S2b interfaces for an E-UTRAN
  ///< initial attach, UE requested PDN connectivity, the PDP
  ///< Context Activation procedure using S4, the PS mobility
  ///< from the Gn/Gp SGSN to the S4 SGSN/MME procedures,
  ///< Attach with GTP on S2b and a UE initiated Connectivity to
  ///< Additional PDN with GTP on S2b.

  ///< Charging Characteristics comes from UpdateLocationAnswer via S6a
  charging_characteristics_t charging_characteristics;

  // EBI Linked EPS Bearer ID             ///< This IE shall be included on
  // S4/S11 in RAU/TAU/HO
  ///< except in the Gn/Gp SGSN to MME/S4-SGSN
  ///< RAU/TAU/HO procedures with SGW change to identify the
  ///< default bearer of the PDN Connection

  protocol_configuration_options_t pco;  /// PCO protocol_configuration_options
  ///< This IE is not applicable to TAU/RAU/Handover. If
  ///< MME/SGSN receives PCO from UE (during the attach
  ///< procedures), the MME/SGSN shall forward the PCO IE to
  ///< SGW. The SGW shall also forward it to PGW.

  bearer_contexts_to_be_created_t
      bearer_contexts_to_be_created;  ///< Bearer Contexts to be created
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

  bearer_contexts_to_be_removed_t
      bearer_contexts_to_be_removed;  ///< This IE shall be included on the
                                      ///< S4/S11 interfaces for the
  ///< TAU/RAU/Handover cases where any of the bearers
  ///< existing before the TAU/RAU/Handover procedure will be
  ///< deactivated as consequence of the TAU/RAU/Handover
  ///< procedure.
  ///< For each of those bearers, an IE with the same type and
  ///< instance value shall be included.

  // Trace Information trace_information  ///< This IE shall be included on the
  // S4/S11 interface if an
  ///< SGW trace is activated, and/or on the S5/S8 and S2b
  ///< interfaces if a PGW trace is activated. See 3GPP TS
  ///< 32.422 [18].

  // Recovery Recovery                    ///< This IE shall be included on the
  // S4/S11, S5/S8 and S2b
  ///< interfaces if contacting the peer for the first time

  FQ_CSID_t mme_fq_csid;  ///< This IE shall be included by the MME on the S11
                          ///< interface
  ///< and shall be forwarded by an SGW on the S5/S8 interfaces
  ///< according to the requirements in 3GPP TS 23.007 [17].

  FQ_CSID_t sgw_fq_csid;  ///< This IE shall included by the SGW on the S5/S8
                          ///< interfaces
  ///< according to the requirements in 3GPP TS 23.007 [17].

  // FQ_CSID_t          epdg_fq_csid;      ///< This IE shall be included by the
  // ePDG on the S2b interface
  ///< according to the requirements in 3GPP TS 23.007 [17].

  UETimeZone_t
      ue_time_zone;  ///< This IE shall be included by the MME over S11 during
  ///< Initial Attach, UE Requested PDN Connectivity procedure.
  ///< This IE shall be included by the SGSN over S4 during PDP
  ///< Context Activation procedure.
  ///< This IE shall be included by the MME/SGSN over S11/S4
  ///< TAU/RAU/Handover with SGW relocation.
  ///< C: If SGW receives this IE, SGW shall forward it to PGW
  ///< across S5/S8 interface.

  UCI_t uci;  ///< User CSG Information
              ///< CO This IE shall be included on the S4/S11 interface for E-
              ///< UTRAN Initial Attach, UE-requested PDN Connectivity and
              ///< PDP Context Activation using S4 procedures if the UE is
              ///< accessed via CSG cell or hybrid cell. The MME/SGSN
              ///< shall also include it for TAU/RAU/Handover procedures if
              ///< the PGW has requested CSG info reporting and
              ///< MME/SGSN support CSG info reporting. The SGW shall
              ///< include this IE on S5/S8 if it receives the User CSG
              ///< information from MME/SGSN.

  // Charging Characteristics
  // MME/S4-SGSN LDN
  // SGW LDN
  // ePDG LDN
  // Signalling Priority Indication
  // MMBR Max MBR/APN-AMBR
  // Private Extension

  /* S11 stack specific parameter. Not used in standalone epc mode */
  void* trxn;  ///< Transaction identifier
  edns_peer_ip_t edns_peer_ip;
  uint16_t peer_port;  ///< MME port for S-GW or S-GW port for MME
} itti_s11_create_session_request_t;

//-----------------------------------------------------------------------------
/** @struct itti_s11_create_session_response_t
 *  @brief Create Session Response
 *
 * The Create Session Response will be sent on S11 interface as
 * part of these procedures:
 * - E-UTRAN Initial Attach
 * - UE requested PDN connectivity
 * - Tracking Area Update procedure with SGW change
 * - S1/X2-based handover with SGW change
 */
typedef struct itti_s11_create_session_response_s {
  teid_t teid;  ///< Tunnel Endpoint Identifier

  // here fields listed in 3GPP TS 29.274
  gtpv2c_cause_t cause;  ///< If the SGW cannot accept any of the "Bearer
                         ///< Context Created" IEs within Create Session Request
  ///< message, the SGW shall send the Create Session Response with appropriate
  ///< reject Cause value.

  // change_reporting_action                    ///< This IE shall be included
  // on the S5/S8 and S4/S11
  ///< interfaces with the appropriate Action field if the location
  ///< Change Reporting mechanism is to be started or stopped
  ///< for this subscriber in the SGSN/MME.

  // csg_Information_reporting_action           ///< This IE shall be included
  // on the S5/S8 and S4/S11
  ///< interfaces with the appropriate Action field if the CSG Info
  ///< reporting mechanism is to be started or stopped for this
  ///< subscriber in the SGSN/MME.

  fteid_t s11_sgw_fteid;  ///< Sender F-TEID for control plane
  ///< This IE shall be sent on the S11/S4 interfaces. For the
  ///< S5/S8/S2b interfaces it is not needed because its content
  ///< would be identical to the IE PGW S5/S8/S2b F-TEID for
  ///< PMIP based interface or for GTP based Control Plane
  ///< interface.

  fteid_t s5_s8_pgw_fteid;  ///< PGW S5/S8/S2b F-TEID for PMIP based interface
                            ///< or for GTP based Control Plane interface
  ///< PGW shall include this IE on the S5/S8 interfaces during
  ///< the Initial Attach, UE requested PDN connectivity and PDP
  ///< Context Activation procedures.
  ///< If SGW receives this IE it shall forward the IE to MME/S4-
  ///< SGSN on S11/S4 interface.
  ///< This IE shall include the TEID in the GTP based S5/S8
  ///< case and the GRE key in the PMIP based S5/S8 case.
  ///< In PMIP based S5/S8 case, same IP address is used for
  ///< both control plane and the user plane communication.
  ///<
  ///< PGW shall include this IE on the S2b interface during the
  ///< Attach with GTP on S2b, UE initiated Connectivity to
  ///< Additional PDN with GTP on S2b and Handover to
  ///< Untrusted Non-3GPP IP Access with GTP on S2b
  ///< procedures.

  paa_t paa;  ///< PDN Address Allocation
  ///< This IE shall be included on the S5/S8, S4/S11 and S2b
  ///< interfaces for the E-UTRAN initial attach, PDP Context
  ///< Activation, UE requested PDN connectivity, Attach with
  ///< GTP on S2b, UE initiated Connectivity to Additional PDN
  ///< with GTP on S2b and Handover to Untrusted Non-3GPP IP
  ///< Access with GTP on S2b procedures.
  ///< The PDN type field in the PAA shall be set to IPv4, or IPv6
  ///< or IPv4v6 by the PGW.
  ///< For the interfaces other than S2b, if the DHCPv4 is used
  ///< for IPv4 address allocation, the IPv4 address field shall be
  ///< set to 0.0.0.0.

  APNRestriction_t
      apn_restriction;  ///< This IE shall be included on the S5/S8 and S4/S11
  ///< interfaces in the E-UTRAN initial attach, PDP Context
  ///< Activation and UE Requested PDN connectivity
  ///< procedures.
  ///< This IE shall also be included on S4/S11 during the Gn/Gp
  ///< SGSN to S4 SGSN/MME RAU/TAU procedures.
  ///< This IE denotes the restriction on the combination of types
  ///< of APN for the APN associated with this EPS bearer
  ///< Context.

  ambr_t ambr;  ///< Aggregate Maximum Bit Rate (APN-AMBR)
  ///< This IE represents the APN-AMBR. It shall be included on
  ///< the S5/S8, S4/S11 and S2b interfaces if the received APN-
  ///< AMBR has been modified by the PCRF.

  // EBI Linked EPS Bearer ID                   ///< This IE shall be sent on
  // the S4/S11 interfaces during
  ///< Gn/Gp SGSN to S4-SGSN/MME RAU/TAU procedure to
  ///< identify the default bearer the PGW selects for the PDN
  ///< Connection.

  protocol_configuration_options_t pco;  // PCO protocol_configuration_options
  ///< This IE is not applicable for TAU/RAU/Handover. If PGW
  ///< decides to return PCO to the UE, PGW shall send PCO to
  ///< SGW. If SGW receives the PCO IE, SGW shall forward it
  ///< MME/SGSN.

  bearer_contexts_created_t
      bearer_contexts_created;  ///< EPS bearers corresponding to Bearer
                                ///< Contexts sent in
  ///< request message. Several IEs with the same type and
  ///< instance value may be included on the S5/S8 and S4/S11
  ///< as necessary to represent a list of Bearers. One single IE
  ///< shall be included on the S2b interface.
  ///< One bearer shall be included for E-UTRAN Initial Attach,
  ///< PDP Context Activation or UE Requested PDN
  ///< Connectivity , Attach with GTP on S2b, UE initiated
  ///< Connectivity to Additional PDN with GTP on S2b, and
  ///< Handover to Untrusted Non-3GPP IP Access with GTP on
  ///< S2b.
  ///< One or more created bearers shall be included for a
  ///< Handover/TAU/RAU with an SGW change. See NOTE 2.

  bearer_contexts_marked_for_removal_t
      bearer_contexts_marked_for_removal;  ///< EPS bearers corresponding to
                                           ///< Bearer Contexts to be
  ///< removed that were sent in the Create Session Request
  ///< message.
  ///< For each of those bearers an IE with the same type and
  ///< instance value shall be included on the S4/S11 interfaces.

  // Recovery Recovery                          ///< This IE shall be included
  // on the S4/S11, S5/S8 and S2b
  ///< interfaces if contacting the peer for the first time

  // FQDN charging_Gateway_name                 ///< When Charging Gateway
  // Function (CGF) Address is
  ///< configured, the PGW shall include this IE on the S5
  ///< interface.
  ///< NOTE 1: Both Charging Gateway Name and Charging Gateway Address shall not
  ///< be included at the same time. When both are available, the operator
  ///< configures a preferred value.

  // IP Address charging_Gateway_address        ///< When Charging Gateway
  // Function (CGF) Address is
  ///< configured, the PGW shall include this IE on the S5
  ///< interface. See NOTE 1.

  FQ_CSID_t
      pgw_fq_csid;  ///< This IE shall be included by the PGW on the S5/S8 and
  ///< S2b interfaces and, when received from S5/S8 be
  ///< forwarded by the SGW on the S11 interface according to
  ///< the requirements in 3GPP TS 23.007 [17].

  FQ_CSID_t sgw_fq_csid;  ///< This IE shall be included by the SGW on the S11
                          ///< interface
  ///< according to the requirements in 3GPP TS 23.007 [17].

  // Local Distinguished Name (LDN) SGW LDN     ///< This IE is optionally sent
  // by the SGW to the MME/SGSN
  ///< on the S11/S4 interfaces (see 3GPP TS 32.423 [44]),
  ///< when contacting the peer node for the first time.
  ///< Also:
  ///< This IE is optionally sent by the SGW to the MME/SGSN
  ///< on the S11/S4 interfaces (see 3GPP TS 32.423 [44]),
  ///< when communicating the LDN to the peer node for the first
  ///< time.

  // Local Distinguished Name (LDN) PGW LDN     ///< This IE is optionally
  // included by the PGW on the S5/S8
  ///< and S2b interfaces (see 3GPP TS 32.423 [44]), when
  ///< contacting the peer node for the first time.
  ///< Also:
  ///< This IE is optionally included by the PGW on the S5/S8
  ///< interfaces (see 3GPP TS 32.423 [44]), when
  ///< communicating the LDN to the peer node for the first time.

  // EPC_Timer pgw_back_off_time                ///< This IE may be included on
  // the S5/S8 and S4/S11
  ///< interfaces when the PDN GW rejects the Create Session
  ///< Request with the cause "APN congestion". It indicates the
  ///< time during which the MME or S4-SGSN should refrain
  ///< from sending subsequent PDN connection establishment
  ///< requests to the PGW for the congested APN for services
  ///< other than Service Users/emergency services.
  ///< See NOTE 3:
  ///< The last received value of the PGW Back-Off Time IE shall supersede any
  ///< previous values received from that PGW and for this APN in the MME/SGSN.

  // Private Extension                          ///< This IE may be sent on the
  // S5/S8, S4/S11 and S2b
  ///< interfaces.

  /* S11 stack specific parameter. Not used in standalone epc mode */
  void* trxn;              ///< Transaction identifier
  struct in_addr peer_ip;  ///< MME ipv4 address
} itti_s11_create_session_response_t;

//-----------------------------------------------------------------------------
/** @struct itti_s11_create_bearer_request_t
 *  @brief Create Bearer Request
 *
 * The direction of this message shall be from PGW to SGW and from SGW to
 * MME/S4-SGSN, and from PGW to ePDG The Create Bearer Request message shall be
 * sent on the S5/S8 interface by the PGW to the SGW and on the S11 interface by
 * the SGW to the MME as part of the Dedicated Bearer Activation procedure. The
 * message shall also be sent on the S5/S8 interface by the PGW to the SGW and
 * on the S4 interface by the SGW to the SGSN as part of the Secondary PDP
 * Context Activation procedure or the Network Requested Secondary PDP Context
 * Activation procedure. The message shall also be sent on the S2b interface by
 * the PGW to the ePDG as part of the Dedicated S2b bearer activation with GTP
 * on S2b.
 */
typedef struct itti_s11_create_bearer_request_s {
  teid_t local_teid;  ///< not in specs for inner use

  teid_t teid;  ///< S11 SGW Tunnel Endpoint Identifier

  pti_t pti;  ///< C: This IE shall be sent on the S5/S8 and S4/S11 interfaces
  ///< when the procedure was initiated by a UE Requested
  ///< Bearer Resource Modification Procedure or UE Requested
  ///< Bearer Resource Allocation Procedure (see NOTE 1) or
  ///< Secondary PDP Context Activation Procedure.
  ///< The PTI shall be the same as the one used in the
  ///< corresponding Bearer Resource Command.

  ebi_t linked_eps_bearer_id;  ///< M: This IE shall be included to indicate the
                               ///< default bearer
  ///< associated with the PDN connection.

  protocol_configuration_options_t
      pco;  ///< O: This IE may be sent on the S5/S8 and S4/S11 interfaces

  bearer_contexts_within_create_bearer_request_t
      bearer_contexts;  ///< M: Several IEs with this type and instance values
                        ///< shall be
  ///< included as necessary to represent a list of Bearers.

  FQ_CSID_t
      pgw_fq_csid;  ///< C: This IE shall be included by MME on S11 and shall be
  ///< forwarded by SGW on S5/S8 according to the
  ///< requirements in 3GPP TS 23.007 [17].

  FQ_CSID_t sgw_fq_csid;  ///< C:This IE shall be included by the SGW on the S11
                          ///< interface
  ///< according to the requirements in 3GPP TS 23.007 [17].

  // Change Reporting Action ///< This IE shall be included on the S5/S8 and
  // S4/S11
  ///< interfaces with the appropriate Action field If the location
  ///< Change Reporting mechanism is to be started or stopped
  ///< for this subscriber in the SGSN/MME.

  // CSG Information ///< This IE shall be included on the S5/S8 and S4/S11
  ///< interfaces with the appropriate Action field if the CSG Info Reporting
  ///< Action reporting mechanism is to be started or stopped for this
  ///< subscriber in the SGSN/MME.

  // Private Extension   Private Extension

  /* GTPv2-C specific parameters */
  void* trxn;  ///< Transaction identifier
  struct in_addr peer_ip;
} itti_s11_create_bearer_request_t;

//-----------------------------------------------------------------------------
/** @struct itti_s11_create_bearer_response_t
 *  @brief Create Bearer Response
 *
 * The Create Bearer Response message shall be sent on the S5/S8 interface by
 * the SGW to the PGW, and on the S11 interface by the MME to the SGW as part of
 * the Dedicated Bearer Activation procedure. The message shall also be sent on
 * the S5/S8 interface by the SGW to the PGW and on the S4 interface by the SGSN
 * to the SGW as part of Secondary PDP Context Activation procedure or the
 * Network Requested Secondary PDP Context Activation procedure. The message
 * shall also be sent on the S2b interface by the ePDG to the PGW as part of the
 * Dedicated S2b bearer activation with GTP on S2b. Possible Cause values are
 * specified in Table 8.4-1. Message specific cause values are:
 * - "Request accepted".
 * - "Request accepted partially".
 * - "Context not found".
 * - "Semantic error in the TFT operation".
 * - "Syntactic error in the TFT operation".
 * - "Semantic errors in packet filter(s)".
 * - "Syntactic errors in packet filter(s)".
 * - "Service not supported".
 * - "Unable to page UE".
 * - "UE not responding".
 * - "Unable to page UE due to Suspension".
 * - "UE refuses".
 * - "Denied in RAT".
 * - "UE context without TFT already activated".
 */
typedef struct itti_s11_create_bearer_response_s {
  teid_t local_teid;  ///< not in specs for inner MME use
  teid_t teid;        ///< S11 MME Tunnel Endpoint Identifier

  // here fields listed in 3GPP TS 29.274
  gtpv2c_cause_t cause;  ///< M

  bearer_contexts_within_create_bearer_response_t
      bearer_contexts;  ///< Several IEs with this type and instance value shall
                        ///< be
  ///< included on the S4/S11, S5/S8 and S2b interfaces as
  ///< necessary to represent a list of Bearers.

  // Recovery   C This IE shall be included on the S4/S11, S5/S8 and S2b
  // interfaces if contacting the peer for the first time

  FQ_CSID_t mme_fq_csid;  ///< C This IE shall be included by the MME on the S11
  ///< interface and shall be forwarded by the SGW on the S5/S8
  ///< interfaces according to the requirements in 3GPP TS
  ///< 23.007 [17].

  FQ_CSID_t sgw_fq_csid;  ///< C This IE shall be included by the MME on the S11
  ///< interface and shall be forwarded by the SGW on the S5/S8
  ///< interfaces according to the requirements in 3GPP TS
  ///< 23.007 [17].

  FQ_CSID_t epdg_fq_csid;  ///< C This IE shall be included by the ePDG on the
                           ///< S2b interface
  ///< according to the requirements in 3GPP TS 23.007 [17].

  protocol_configuration_options_t
      pco;  ///< C: If the UE includes the PCO IE, then the MME/SGSN shall
  ///< copy the content of this IE transparently from the PCO IE
  ///< included by the UE. If the SGW receives PCO from
  ///< MME/SGSN, SGW shall forward it to the PGW.

  UETimeZone_t ue_time_zone;  ///< O: This IE is optionally included by the MME
                              ///< on the S11
  ///< interface or by the SGSN on the S4 interface.
  ///< CO: The SGW shall forward this IE on the S5/S8 interface if the
  ///< SGW supports this IE and it receives it from the
  ///< MME/SGSN.

  Uli_t uli;  ///< O: This IE is optionally included by the MME on the S11
  ///< interface or by the SGSN on the S4 interface.
  ///< CO: The SGW shall forward this IE on the S5/S8 interface if the
  ///< SGW supports this IE and it receives it from the
  ///< MME/SGSN.

  // Private Extension Private Extension        ///< optional

  /* S11 stack specific parameter. Not used in standalone epc mode */
  void* trxn;  ///< Transaction identifier
} itti_s11_create_bearer_response_t;

//-----------------------------------------------------------------------------
/** @struct itti_s11_modify_bearer_request_t
 *  @brief Modify Bearer Request
 *
 * The Modify Bearer Request will be sent on S11 interface as
 * part of these procedures:
 * - E-UTRAN Tracking Area Update without SGW Change
 * - UE triggered Service Request
 * - S1-based Handover
 * - E-UTRAN Initial Attach
 * - UE requested PDN connectivity
 * - X2-based handover without SGWrelocation
 */
typedef struct itti_s11_modify_bearer_request_s {
  teid_t local_teid;  ///< not in specs for inner MME use

  teid_t teid;  ///< S11 SGW Tunnel Endpoint Identifier

  // MEI                    ME Identity (MEI)  ///< C:This IE shall be sent on
  // the S5/S8 interfaces for the Gn/Gp
  ///< SGSN to MME TAU.

  Uli_t uli;  ///< C: The MME/SGSN shall include this IE for
  ///< TAU/RAU/Handover procedures if the PGW has requested
  ///< location information change reporting and MME/SGSN
  ///< support location information change reporting.
  ///< An MME/SGSN which supports location information
  ///< change shall include this IE for UE-initiated Service
  ///< Request procedure if the PGW has requested location
  ///< information change reporting and the UE's location info
  ///< has changed.
  ///< The SGW shall include this IE on S5/S8 if it receives the
  ///< ULI from MME/SGSN.
  ///< CO:This IE shall also be included on the S4/S11 interface for a
  ///< TAU/RAU/Handover with MME/SGSN change without
  ///< SGW change procedure, if the level of support (User
  ///< Location Change Reporting and/or CSG Information
  ///< Change Reporting) changes the MME shall include the
  ///< ECGI/TAI in the ULI, the SGSN shall include the CGI/SAI
  ///< in the ULI.
  ///< The SGW shall include this IE on S5/S8 if it receives the
  ///< ULI from MME/SGSN.

  ServingNetwork_t serving_network;  ///< CO:This IE shall be included on S11/S4
                                     ///< interface during the
  ///< following procedures:
  ///< - TAU/RAU/handover if Serving Network is changed.
  ///< - TAU/RAU when the UE was ISR activated which is
  ///<   indicated by ISRAU flag.
  ///< - UE triggered Service Request when UE is ISR
  ///<   activated.
  ///< - UE initiated Service Request if ISR is not active, but
  ///<   the Serving Network has changed during previous
  ///<   mobility procedures, i.e. intra MME/S4-SGSN
  ///<   TAU/RAU and the change has not been reported to
  ///<   the PGW yet.
  ///< - TAU/RAU procedure as part of the optional network
  ///<   triggered service restoration procedure with ISR, as
  ///<   specified by 3GPP TS 23.007 [17].
  ///<
  ///< CO:This IE shall be included on S5/S8 if the SGW receives this
  ///< IE from MME/SGSN and if ISR is not active.
  ///< This IE shall be included on S5/S8 if the SGW receives this
  ///< IE from MME/SGSN and ISR is active and the Modify
  ///< Bearer Request message needs to be sent to the PGW as
  ///< specified in the 3GPP TS 23.401 [3].

  rat_type_t rat_type;  ///< C: This IE shall be sent on the S11 interface for a
                        ///< TAU with
  ///< an SGSN interaction, UE triggered Service Request or an I-
  ///< RAT Handover.
  ///< This IE shall be sent on the S4 interface for a RAU with
  ///< MME interaction, a RAU with an SGSN change, a UE
  ///< Initiated Service Request or an I-RAT Handover.
  ///< This IE shall be sent on the S5/S8 interface if the RAT type
  ///< changes.
  ///< CO: CO If SGW receives this IE from MME/SGSN during a
  ///< TAU/RAU/Handover with SGW change procedure, the
  ///< SGW shall forward it across S5/S8 interface to PGW.
  ///< CO: The IE shall be sent on the S11/S4 interface during the
  ///< following procedures:
  ///< - an inter MM TAU or inter SGSN RAU when UE was
  ///<   ISR activated which is indicated by ISRAU flag.
  ///< - TAU/RAU procedure as part of optional network
  ///<   triggered service restoration procedure with ISR, as
  ///<   specified by 3GPP TS 23.007 [17].
  ///< If ISR is active, this IE shall also be included on the S11
  ///< interface in the S1-U GTP-U tunnel setup procedure during
  ///< an intra-MME intra-SGW TAU procedure.

  indication_flags_t indication_flags;  ///< C:This IE shall be included if any
                                        ///< one of the applicable flags
  ///< is set to 1.
  ///< Applicable flags are:
  ///< -ISRAI: This flag shall be used on S4/S11 interface
  ///<   and set to 1 if the ISR is established between the
  ///<   MME and the S4 SGSN.
  ///< - Handover Indication: This flag shall be set to 1 on
  ///<   the S4/S11 and S5/S8 interfaces during an E-
  ///<   UTRAN Initial Attach or for a UE Requested PDN
  ///<   Connectivity or a PDP Context Activation
  ///<   procedure, if the PDN connection/PDP context is
  ///<   handed-over from non-3GPP access.
  ///< - Direct Tunnel Flag: This flag shall be used on the
  ///<   S4 interface and set to 1 if Direct Tunnel is used.
  ///< - Change Reporting support Indication: shall be
  ///<   used on S4/S11, S5/S8 and set if the SGSN/MME
  ///<   supports location Info Change Reporting. This flag
  ///<   should be ignored by SGW if no message is sent
  ///<   on S5/S8. See NOTE 4.
  ///< - CSG Change Reporting Support Indication: shall
  ///<   be used on S4/S11, S5/S8 and set if the
  ///<   SGSN/MME supports CSG Information Change
  ///<   Reporting. This flag shall be ignored by SGW if no
  ///<   message is sent on S5/S8. See NOTE 4.
  ///< - Change F-TEID support Indication: This flag shall
  ///<   be used on S4/S11 for an IDLE state UE initiated
  ///<   TAU/RAU procedure and set to 1 to allow the
  ///<   SGW changing the GTP-U F-TEID.

  fteid_t sender_fteid_for_cp;  ///< C: Sender F-TEID for control plane
  ///< This IE shall be sent on the S11 and S4 interfaces for a
  ///< TAU/RAU/ Handover with MME/SGSN change and without
  ///< any SGW change.
  ///< This IE shall be sent on the S5 and S8 interfaces for a
  ///< TAU/RAU/Handover with a SGW change.

  ambr_t apn_ambr;  ///< C: Aggregate Maximum Bit Rate (APN-AMBR)
  ///< The APN-AMBR shall be sent for the PS mobility from the
  ///< Gn/Gp SGSN to the S4 SGSN/MME procedures..

  /* Delay Value in integer multiples of 50 millisecs, or zero */
  DelayValue_t delay_dl_packet_notif_req;  ///< C:This IE shall be sent on the
                                           ///< S11 interface for a UE
  ///< triggered Service Request.
  ///< CO: This IE shall be sent on the S4 interface for a UE triggered
  ///< Service Request.

  bearer_contexts_to_be_modified_t
      bearer_contexts_to_be_modified;  ///< C: This IE shall be sent on the
                                       ///< S4/S11 interface and S5/S8
  ///< interface except on the S5/S8 interface for a UE triggered
  ///< Service Request.
  ///< When Handover Indication flag is set to 1 (i.e., for
  ///< EUTRAN Initial Attach or UE Requested PDN Connectivity
  ///< when the UE comes from non-3GPP access), the PGW
  ///< shall ignore this IE. See NOTE 1.
  ///< Several IEs with the same type and instance value may be
  ///< included as necessary to represent a list of Bearers to be
  ///< modified.
  ///< During a TAU/RAU/Handover procedure with an SGW
  ///< change, the SGW includes all bearers it received from the
  ///< MME/SGSN (Bearer Contexts to be created, or Bearer
  ///< Contexts to be modified and also Bearer Contexts to be
  ///< removed) into the list of 'Bearer Contexts to be modified'
  ///< IEs, which are then sent on the S5/S8 interface to the
  ///< PGW (see NOTE 2).

  bearer_contexts_to_be_removed_t
      bearer_contexts_to_be_removed;  ///< C: This IE shall be included on the
                                      ///< S4 and S11 interfaces for
  ///< the TAU/RAU/Handover and Service Request procedures
  ///< where any of the bearers existing before the
  ///< TAU/RAU/Handover procedure and Service Request
  ///< procedures will be deactivated as consequence of the
  ///< TAU/RAU/Handover procedure and Service Request
  ///< procedures. (NOTE 3)
  ///< For each of those bearers, an IE with the same type and
  ///< instance value, shall be included.

  // recovery_t(restart counter) recovery;      ///< C: This IE shall be
  // included if contacting the peer for the first
  ///< time.

  UETimeZone_t ue_time_zone;  ///< CO: This IE shall be included by the MME/SGSN
                              ///< on the S11/S4
  ///< interfaces if the UE Time Zone has changed in the case of
  ///< TAU/RAU/Handover.
  ///< C: If SGW receives this IE, SGW shall forward it to PGW
  ///< across S5/S8 interface.

  FQ_CSID_t
      mme_fq_csid;  ///< C: This IE shall be included by MME on S11 and shall be
  ///< forwarded by SGW on S5/S8 according to the
  ///< requirements in 3GPP TS 23.007 [17].

  FQ_CSID_t sgw_fq_csid;  ///< C: This IE shall be included by SGW on S5/S8
                          ///< according to
  ///< the requirements in 3GPP TS 23.007 [17].

  UCI_t uci;  ///< CO: The MME/SGSN shall include this IE for
  ///< TAU/RAU/Handover procedures and UE-initiated Service
  ///< Request procedure if the PGW has requested CSG Info
  ///< reporting and the MME/SGSN support the CSG
  ///< information reporting. The SGW shall include this IE on
  ///< S5/S8 if it receives the User CSG Information from
  ///< MME/SGSN.

  // Local Distinguished Name (LDN) MME/S4-SGSN LDN ///< O: This IE is
  // optionally sent by the MME to the SGW on the
  ///< S11 interface and by the SGSN to the SGW on the S4
  ///< interface (see 3GPP TS 32.423 [44]), when communicating
  ///< the LDN to the peer node for the first time.

  // Local Distinguished Name (LDN) SGW LDN     ///< O: This IE is optionally
  // sent by the SGW to the PGW on the
  ///< S5/S8 interfaces (see 3GPP TS 32.423 [44]), for inter-
  ///< SGW mobity, when communicating the LDN to the peer
  ///< node for the first time.

  // MMBR           Max MBR/APN-AMBR            ///< CO: If the S4-SGSN supports
  // Max MBR/APN-AMBR, this IE
  ///< shall be included by the S4-SGSN over S4 interface in the
  ///< following cases:
  ///< - during inter SGSN RAU/SRNS relocation without
  ///<   SGW relocation and inter SGSN SRNS relocation
  ///<   with SGW relocation if Higher bitrates than
  ///<   16 Mbps flag is not included in the MM Context IE
  ///<   in the Context Response message or in the MM
  ///<   Context IE in the Forward Relocation Request
  ///<   message from the old S4-SGSN, while it is
  ///<   received from target RNC or a local Max
  ///<   MBR/APN-AMBR is configured based on
  ///<   operator's policy.
  ///<   - during Service Request procedure if Higher
  ///<   bitrates than 16 Mbps flag is received but the S4-
  ///<   SGSN has not received it before from an old RNC
  ///<   or the S4-SGSN has not updated the Max
  ///<   MBR/APN-AMBR to the PGW yet.
  ///< If SGW receives this IE, SGW shall forward it to PGW
  ///< across S5/S8 interface.

  // Private Extension   Private Extension

  /* GTPv2-C specific parameters */
  edns_peer_ip_t edns_peer_ip;
  void* trxn;  ///< Transaction identifier
  uint8_t internal_flags;
} itti_s11_modify_bearer_request_t;

//-----------------------------------------------------------------------------
/** @struct itti_s11_modify_bearer_response_t
 *  @brief Modify Bearer Response
 *
 * The Modify Bearer Response will be sent on S11 interface as
 * part of these procedures:
 * - E-UTRAN Tracking Area Update without SGW Change
 * - UE triggered Service Request
 * - S1-based Handover
 * - E-UTRAN Initial Attach
 * - UE requested PDN connectivity
 * - X2-based handover without SGWrelocation
 */
typedef struct itti_s11_modify_bearer_response_s {
  teid_t teid;  ///< S11 MME Tunnel Endpoint Identifier

  // here fields listed in 3GPP TS 29.274
  gtpv2c_cause_t cause;  ///<

  ebi_t linked_eps_bearer_id;  ///< This IE shall be sent on S5/S8 when the UE
                               ///< moves from a
  ///< Gn/Gp SGSN to the S4 SGSN or MME to identify the
  ///< default bearer the PGW selects for the PDN Connection.
  ///< This IE shall also be sent by SGW on S11, S4 during
  ///< Gn/Gp SGSN to S4-SGSN/MME HO procedures to identify
  ///< the default bearer the PGW selects for the PDN
  ///< Connection.

  ambr_t apn_ambr;  ///< Aggregate Maximum Bit Rate (APN-AMBR)
  ///< This IE shall be included in the PS mobility from Gn/Gp
  ///< SGSN to the S4 SGSN/MME procedures if the received
  ///< APN-AMBR has been modified by the PCRF.

  APNRestriction_t apn_restriction;  ///< This IE denotes the restriction on the
                                     ///< combination of types
  ///< of APN for the APN associated with this EPS bearer
  ///< Context. This IE shall be included over S5/S8 interfaces,
  ///< and shall be forwarded over S11/S4 interfaces during
  ///< Gn/Gp SGSN to MME/S4-SGSN handover procedures.
  ///< This IE shall also be included on S5/S8 interfaces during
  ///< the Gn/Gp SGSN to S4 SGSN/MME RAU/TAU
  ///< procedures.
  ///< The target MME or SGSN determines the Maximum APN
  ///< Restriction using the APN Restriction.
  // PCO protocol_configuration_options         ///< If SGW receives this IE
  // from PGW on GTP or PMIP based
  ///< S5/S8, the SGW shall forward PCO to MME/S4-SGSN
  ///< during Inter RAT handover from the UTRAN or from the
  ///< GERAN to the E-UTRAN. See NOTE 2:
  ///< If MME receives the IE, but no NAS message is sent, MME discards the IE.

  bearer_contexts_modified_t
      bearer_contexts_modified;  ///< EPS bearers corresponding to Bearer
                                 ///< Contexts to be
  ///< modified that were sent in Modify Bearer Request
  ///< message. Several IEs with the same type and instance
  ///< value may be included as necessary to represent a list of
  ///< the Bearers which are modified.

  bearer_contexts_marked_for_removal_t
      bearer_contexts_marked_for_removal;  ///< EPS bearers corresponding to
                                           ///< Bearer Contexts to be
  ///< removed sent in the Modify Bearer Request message.
  ///< Shall be included if request message contained Bearer
  ///< Contexts to be removed.
  ///< For each of those bearers an IE with the same type and
  ///< instance value shall be included.

  // change_reporting_action                    ///< This IE shall be included
  // with the appropriate Action field If
  ///< the location Change Reporting mechanism is to be started
  ///< or stopped for this subscriber in the SGSN/MME.

  // csg_Information_reporting_action           ///< This IE shall be included
  // with the appropriate Action field if
  ///< the location CSG Info change reporting mechanism is to be
  ///< started or stopped for this subscriber in the SGSN/MME.

  // FQDN Charging Gateway Name                 ///< When Charging Gateway
  // Function (CGF) Address is
  ///< configured, the PGW shall include this IE on the S5
  ///< interface during SGW relocation and when the UE moves
  ///< from Gn/Gp SGSN to S4-SGSN/MME. See NOTE 1:
  ///< Both Charging Gateway Name and Charging Gateway Address shall not be
  ///< included at the same time. When both are available, the operator
  ///< configures a preferred value.

  // IP Address Charging Gateway Address        ///< When Charging Gateway
  // Function (CGF) Address is
  ///< configured, the PGW shall include this IE on the S5
  ///< interface during SGW relocation and when the UE moves
  ///< from Gn/Gp SGSN to S4-SGSN/MME. See NOTE 1:
  ///< Both Charging Gateway Name and Charging Gateway Address shall not be
  ///< included at the same time. When both are available, the operator
  ///< configures a preferred value.

  FQ_CSID_t
      pgw_fq_csid;  ///< This IE shall be included by PGW on S5/S8and shall be
  ///< forwarded by SGW on S11 according to the requirements
  ///< in 3GPP TS 23.007 [17].

  FQ_CSID_t sgw_fq_csid;  ///< This IE shall be included by SGW on S11 according
                          ///< to the
  ///< requirements in 3GPP TS 23.007 [17].

  // recovery_t(restart counter) recovery;      ///< This IE shall be included
  // if contacting the peer for the first
  ///< time.

  // Local Distinguished Name (LDN) SGW LDN     ///< This IE is optionally sent
  // by the SGW to the MME/SGSN
  ///< on the S11/S4 interfaces (see 3GPP TS 32.423 [44]),
  ///< when contacting the peer node for the first time.

  // Local Distinguished Name (LDN) PGW LDN     ///< This IE is optionally
  // included by the PGW on the S5/S8
  ///< and S2b interfaces (see 3GPP TS 32.423 [44]), when
  ///< contacting the peer node for the first time.

  // Private Extension Private Extension        ///< optional

  /* S11 stack specific parameter. Not used in standalone epc mode */
  void* trxn;  ///< Transaction identifier
  uint8_t internal_flags;
} itti_s11_modify_bearer_response_t;

//-----------------------------------------------------------------------------
typedef struct itti_s11_delete_session_request_s {
  teid_t local_teid;  ///< not in specs for inner MME use
  teid_t teid;        ///< Tunnel Endpoint Identifier
  ebi_t lbi;          ///< Linked EPS Bearer ID
  bool noDelete;
  fteid_t sender_fteid_for_cp;  ///< Sender F-TEID for control plane
  uint8_t internal_flags;

  /* Operation Indication: This flag shall be set over S4/S11 interface
   * if the SGW needs to forward the Delete Session Request message to
   * the PGW. This flag shall not be set if the ISR associated GTP
   * entity sends this message to the SGW in the Detach procedure.
   * This flag shall also not be set to 1 in the SRNS Relocation Cancel
   * Using S4 (6.9.2.2.4a in 3GPP TS 23.060 [4]), Inter RAT handover
   * Cancel procedure with SGW change TAU with Serving GW change,
   * Gn/Gb based RAU (see 5.5.2.5, 5.3.3.1, D.3.5 in 3GPP TS 23.401 [3],
   * respectively), S1 Based handover Cancel procedure with SGW change.
   */
  indication_flags_t indication_flags;

  /* GTPv2-C specific parameters */
  void* trxn;  ///< Transaction identifier
  edns_peer_ip_t edns_peer_ip;
  struct in_addr peer_ip;
  Uli_t uli;
  ServingNetwork_t serving_network;
} itti_s11_delete_session_request_t;

//-----------------------------------------------------------------------------
/** @struct itti_s11_delete_session_response_t
 *  @brief Delete Session Response
 *
 * The Delete Session Response will be sent on S11 interface as
 * part of these procedures:
 * - EUTRAN Initial Attach
 * - UE, HSS or MME Initiated Detach
 * - UE or MME Requested PDN Disconnection
 * - Tracking Area Update with SGW Change
 * - S1 Based Handover with SGW Change
 * - X2 Based Handover with SGW Relocation
 * - S1 Based handover cancel with SGW change
 */
typedef struct itti_s11_delete_session_response_s {
  teid_t teid;  ///< Remote Tunnel Endpoint Identifier
  gtpv2c_cause_t cause;
  uint8_t internal_flags;  ///< Flags used for response messages sent over
                           ///< gtpv2c stack
  // recovery_t recovery;              ///< This IE shall be included on the
  // S5/S8, S4/S11 and S2b
  ///< interfaces if contacting the peer for the first time
  protocol_configuration_options_t
      pco;  ///< PGW shall include Protocol Configuration Options (PCO)
            ///< IE on the S5/S8 interface, if available.
            ///< If SGW receives this IE, SGW shall forward it to
            ///< SGSN/MME on the S4/S11 interface.

  /* GTPv2-C specific parameters */
  void* trxn;
  struct in_addr peer_ip;
  ebi_t lbi;
} itti_s11_delete_session_response_t;

//-----------------------------------------------------------------------------
/** @struct itti_s11_release_access_bearers_request_t
 *  @brief Release AccessBearers Request
 *
 * The Release Access Bearers Request message shall sent on the S11 interface by
 * the MME to the SGW as part of the S1 release procedure.
 * The message shall also be sent on the S4 interface by the SGSN to the SGW as
 * part of the procedures:
 * -    RAB release using S4
 * -    Iu Release using S4
 * -    READY to STANDBY transition within the network
 */
typedef struct itti_s11_release_access_bearers_request_s {
  teid_t local_teid;        ///< not in specs for inner MME use
  teid_t teid;              ///< Tunnel Endpoint Identifier
  ebi_list_t list_of_rabs;  ///< Shall be present on S4 interface when this
                            ///< message is used to release a subset of all
                            ///< active RABs according to the RAB release
                            ///< procedure. Several IEs with this type and
                            ///< instance values shall be included as necessary
                            ///< to represent a list of RABs to be released.
  node_type_t originating_node;  ///< This IE shall be sent on S11 interface, if
                                 ///< ISR is active in the MME.
  ///< This IE shall be sent on S4 interface, if ISR is active in the SGSN
  // Private Extension Private Extension ///< optional
  /* GTPv2-C specific parameters */
  void* trxn;
  edns_peer_ip_t edns_peer_ip;
} itti_s11_release_access_bearers_request_t;

//-----------------------------------------------------------------------------
/** @struct itti_s11_release_access_bearers_response_t
 *  @brief Release AccessBearers Response
 *
 * The Release Access Bearers Response message is sent on the S11 interface by
 * the SGW to the MME as part of the S1 release procedure. The message shall
 * also be sent on the S4 interface by the SGW to the SGSN as part of the
 * procedures:
 * -  RAB release using S4
 * -  Iu Release using S4
 * -  READY to STANDBY transition within the network
 * Possible Cause values are specified in Table 8.4-1. Message specific cause
 * values are:
 * - "Request accepted".
 * - "Request accepted partially".
 * - "Context not found
 */
typedef struct itti_s11_release_access_bearers_response_s {
  teid_t teid;  ///< Tunnel Endpoint Identifier
  gtpv2c_cause_t cause;
  // Recovery           ///< optional This IE shall be included if contacting
  // the peer for the first time Private Extension  ///< optional
  /* GTPv2-C specific parameters */
  void* trxn;
  struct in_addr peer_ip;
} itti_s11_release_access_bearers_response_t;

//-----------------------------------------------------------------------------
/** @struct itti_s11_downlink_data_notification_t
 *  @brief Downlink Data Notification
 *
 * The Downlink Data Notification message is sent on the S11 interface by the
 * SGW to the MME as part of the S1 paging procedure.
 */
typedef struct itti_s11_downlink_data_notification_s {
  teid_t teid;  ///< Tunnel Endpoint Identifier
  /* GTPv2-C specific parameters */
  void* trxn;
  struct sockaddr* peer_ip;
} itti_s11_downlink_data_notification_t;

//-----------------------------------------------------------------------------
/** @struct itti_s11_downlink_data_notification_acknowledge_t
 *  @brief Downlink Data Notification Acknowledge
 *
 * The Downlink Data Notification Acknowledge message is sent on the S11
 * interface by the MME to the SGW as part of the S1 paging procedure.
 */
typedef struct itti_s11_downlink_data_notification_acknowledge_s {
  teid_t teid;        ///< Tunnel Endpoint Identifier
  teid_t local_teid;  ///< Tunnel Endpoint Identifier
  gtpv2c_cause_t cause;
  /* GTPv2-C specific parameters */
  void* trxn;
  struct sockaddr* peer_ip;
} itti_s11_downlink_data_notification_acknowledge_t;

//-----------------------------------------------------------------------------
/** @struct itti_s11_delete_bearer_command_t
 *  @brief Initiate Delete Bearer procedure
 *
 * A Delete Bearer Command message shall be sent on the S11 interface by the MME
 * to the SGW and on the S5/S8 interface by the SGW to the PGW as a part of the
 * eNodeB requested bearer release or MME-Initiated Dedicated Bearer
 * Deactivation procedure.
 * The message shall also be sent on the S4 interface by the SGSN to the SGW and
 * on the S5/S8 interface by the SGW to the PGW as part of the MS and SGSN
 * Initiated Bearer Deactivation procedure using S4.
 */
typedef struct itti_s11_delete_bearer_command_s {
  teid_t teid;        ///< Tunnel Endpoint Identifier
  teid_t local_teid;  ///< Tunnel Endpoint Identifier

  // TODO
  void* trxn;
  edns_peer_ip_t edns_peer_ip;
  ebi_list_t ebi_list;
} itti_s11_delete_bearer_command_t;

/**
 * Message used to notify MME that a paging message should be sent to the UE
 * at the given imsi
 * If imsi is available send imsi if not send ue_ip address. Mme app shall
 * fetch imsi from ue_ip address
 */
typedef struct itti_s11_paging_request_s {
  const char* imsi;
  struct in_addr ipv4_addr;
} itti_s11_paging_request_t;

/**
 * Message used to notify SPGW that a paging
 */
typedef struct itti_s11_paging_response_s {
  const char* imsi;
  bool successful;
} itti_s11_paging_response_t;

/**
 * Suspend Notification Message is sent on the S11 interface by MME to SPGW as
 * part of CSFB handling in case, UE or GERAN network are not available for PS
 * handover (DTM is not supported) MME shall indicate to SPGW that UE/bearer is
 * suspended and shall discard PS data.
 */
typedef struct itti_s11_suspend_notification_s {
  teid_t teid;  ///< S11- S-GW Tunnel Endpoint Identifier
  Imsi_t imsi;  ///< The IMSI shall be included in the message on the S11
  ebi_t lbi;    ///< Linked EPS Bearer ID
} itti_s11_suspend_notification_t;

typedef struct itti_s11_suspend_acknowledge_s {
  teid_t teid;           ///< S11 MME Tunnel Endpoint Identifier
  gtpv2c_cause_t cause;  ///< If IMSI is absent in Suspend Notification
                         ///< cause is set to "mandatory IE missing"
                         ///< else shall send as "Request accepted"
} itti_s11_suspend_acknowledge_t;

/**
 * Modify UE AMBR Request is used to send from SPGW to MME_APP,
 * need to modify if PCRF changes UE AMBR after UE is attached.
 */
typedef struct itti_s11_modify_ue_ambr_request_s {
  teid_t teid;     ///< Tunnel Endpoint Identifier
  ambr_t ue_ambr;  ///< Aggregate Maximum Bit Rate (UE-AMBR)
} itti_s11_modify_ue_ambr_request_t;

#endif /* FILE_S11_MESSAGES_TYPES_SEEN */
