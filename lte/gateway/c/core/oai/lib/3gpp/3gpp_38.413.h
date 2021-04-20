/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#pragma once

#include "3gpp_23.003.h"
#include "bstrlib.h"
#include "n11_messages_types.h"
#include "TrackingAreaIdentity.h"

#define MAX_NO_OF_PDUSESSIONS 16
#define MAX_QosFlow 16

typedef enum Type_of_Message_s {
  initiating_message = 1,
  successful_outcome,
  unsuccessfull_outcome
} Type_of_Message_t;

// 9.3 Information Element Definitions
// 9.3.1 Radio Network Layer Related IEs
// 9.3.1.1 Message Type
typedef struct Ngap_Message_Type_s {
  uint32_t Procedure_Code;
  Type_of_Message_t Type_of_Message;
} Ngap_Message_Type_t;

// 9.3.3.15 RAN Paging Priority IE contains the service priority as defined in
// TS 23.501
typedef enum Ngap_Paging_Priority_s {
  Ngap_PagingPriority_priolevel1_t = 0,
  Ngap_PagingPriority_priolevel2_t = 1,
  Ngap_PagingPriority_priolevel3_t = 2,
  Ngap_PagingPriority_priolevel4_t = 3,
  Ngap_PagingPriority_priolevel5_t = 4,
  Ngap_PagingPriority_priolevel6_t = 5,
  Ngap_PagingPriority_priolevel7_t = 6,
  Ngap_PagingPriority_priolevel8_t = 7
  /*
   * Enumeration is extensible
   */
} Ngap_Paging_Priority_t;

typedef long Ngap_PDUSessionID_t;
/* Ngap_S-NSSAI */
typedef struct Ngap_S_NSSAI_s {
  bstring sST;
  bstring* sD; /* OPTIONAL */
} Ngap_SNSSAI_t;

// 9.3.1.58 UE Aggregate Maximum Bit Rate
// all Non-GBR QoS flows per UE which is defined for the downlink and the uplink
// direction and a subscription parameter provided by the AMF to the NG-RAN node
typedef uint64_t bit_rate_t;
typedef struct ngap_ue_aggregate_maximum_bit_rate_s {
  bit_rate_t dl;
  bit_rate_t ul;
} ngap_ue_aggregate_maximum_bit_rate_t;

typedef uint32_t amf_ue_ngap_id_ty;
typedef uint32_t ran_ue_ngap_id_t;

typedef struct pdusession_setup_item_s {
  bstring nas_pdu;  // 5GC – UE or UE – 5GC message that is transferred without
                    // interpretation in the NG-RAN node  /*optional*/
  Ngap_PDUSessionID_t
      Pdu_Session_ID;  // PDU Session for a UE. The definition and use of the
                       // PDU Session ID is specified in TS 23.501 [9].
  Ngap_SNSSAI_t Ngap_s_nssai;  // S-NSSAI as defined in TS 23.003 [23].
  pdu_session_resource_setup_request_transfer_t
      PDU_Session_Resource_Setup_Request_Transfer;  // Containing the PDU
                                                    // Session Resource
                                                    // Setup Request
} pdusession_setup_item_t;

typedef struct Ngap_PDUSession_Resource_Setup_Request_List_s {
  uint16_t no_of_items;
  pdusession_setup_item_t item[MAX_NO_OF_PDUSESSIONS];

} Ngap_PDUSession_Resource_Setup_Request_List_t;

typedef struct response_gtp_tunnel_s {
  char transportLayerAddress[4];
  char gTP_TEID[4];
} response_gtp_tunnel_t;

typedef struct AssociatedQosFlowList_s {
  int items;
  int QosFlowIdentifier[MAX_QosFlow];
} AssociatedQosFlowList_t;

typedef struct QosFlowPerTNLInformation_s {
  response_gtp_tunnel_t tunnel;
  AssociatedQosFlowList_t associatedQosFlowList;
} QosFlowPerTNLInformation_t;

typedef struct pdusession_setup_response_item_s {
  Ngap_PDUSessionID_t
      Pdu_Session_ID;  // PDU Session for a UE. The definition and use of the
                       // PDU Session ID is specified in TS 23.501 [9].
  QosFlowPerTNLInformation_t
      PDU_Session_Resource_Setup_Response_Transfer;  // Containing the PDU
                                                     // Session Resource
                                                     // Setup Request
} pdusession_setup_response_item_t;

typedef struct Ngap_PDUSession_Resource_Setup_Response_List_s {
  uint16_t no_of_items;
  pdusession_setup_response_item_t item[MAX_NO_OF_PDUSESSIONS];

} Ngap_PDUSession_Resource_Setup_Response_List_t;

typedef struct pdusession_resource_failed_To_setup_item_s {
  Ngap_PDUSessionID_t
      Pdu_Session_ID;  // PDU Session for a UE. The definition and use of the
                       // PDU Session ID is specified in TS 23.501 [9].
  bstring
      PDU_Session_Resource_Setup_Unsuccessful_Transfer;  // Containing the PDU
                                                         // Session Resource
                                                         // Setup Request
} pdusession_resource_failed_To_setup_item_t;

typedef struct Ngap_PDUSession_Resource_Failed_To_Setup_List_s {
  uint16_t no_of_items;
  pdusession_resource_failed_To_setup_item_t item[MAX_NO_OF_PDUSESSIONS];

} Ngap_PDUSession_Resource_Failed_To_Setup_List_t;

typedef struct pdusession_resource_released_item_t_s {
  Ngap_PDUSessionID_t
      Pdu_Session_ID;  // PDU Session for a UE. The definition and use of the
                       // PDU Session ID is specified in TS 23.501 [9].
  bstring PDU_Session_Resource_Release_Response_Transfer;  // Containing the PDU
                                                           // Session Resource
                                                           // Setup Request
} pdusession_resource_released_item_t;

typedef struct Ngap_PDUSession_Resource_Released_List_s {
  uint16_t no_of_items;
  pdusession_resource_released_item_t item[MAX_NO_OF_PDUSESSIONS];

} Ngap_PDUSession_Resource_Released_List_t;

typedef struct pdusession_resource_to_released_item_s {
  Ngap_PDUSessionID_t
      Pdu_Session_ID;  // PDU Session for a UE. The definition and use of the
                       // PDU Session ID is specified in TS 23.501 [9].
  pdu_session_resource_release_command_transfer
      PDU_Session_Resource_TO_Release_Command_Transfer;  // Containing the PDU
                                                         // Session Resource
                                                         // Setup Request
} pdusession_resource_to_released_item_t;

typedef struct Ngap_PDUSession_Resource_TO_Release_List_s {
  uint16_t no_of_items;
  pdusession_resource_to_released_item_t item[MAX_NO_OF_PDUSESSIONS];

} Ngap_PDUSession_Resource_TO_Release_List_t;

// 9.2 Message Functional Definition and Content
// 9.2.1 PDU Session Management Messages
// 9.2.1.1 PDU SESSION RESOURCE SETUP REQUEST
// Direction: AMF → NG-RAN node
typedef struct PDU_Session_resource_setup_request_s {
  Ngap_Message_Type_t Ngap_Message_Type;
  amf_ue_ngap_id_ty
      amf_ue_ngap_id;  // This IE uniquely identifies the UE association over
                       // the NG interface, as described in TS 38.401
  ran_ue_ngap_id_t
      ran_ue_ngap_id;  // This IE uniquely identifies the UE association over
                       // the NG interface within the NG-RAN node
  Ngap_Paging_Priority_t RAN_Paging_Priority; /*optional*/
  bstring nas_pdu;  // 5GC – UE or UE – 5GC message that is transferred without
                    // interpretation in the NG-RAN node  /*optional*/
  Ngap_PDUSession_Resource_Setup_Request_List_t pdusesssion_setup_list;
} PDU_Session_resource_setup_request_t;

// 9.2.1.2 PDU SESSION RESOURCE SETUP RESPONSE
// This message is sent by the NG-RAN node as a response to the request to
// assign resources on Uu and NG-U for one or several PDU session resources.
// Direction: NG-RAN node → AMF
typedef struct PDU_Session_resource_setup_response_s {
  Ngap_Message_Type_t Ngap_Message_Type;
  amf_ue_ngap_id_t
      amf_ue_ngap_id;  // This IE uniquely identifies the UE association over
                       // the NG interface, as described in TS 38.401
  ran_ue_ngap_id_t
      ran_ue_ngap_id;  // This IE uniquely identifies the UE association over
                       // the NG interface within the NG-RAN node
  Ngap_PDUSessionID_t
      Pdu_Session_ID;  // PDU Session for a UE. The definition and use of the
                       // PDU Session ID is specified in TS 23.501 [9].
} PDU_Session_resource_setup_response_t;

typedef struct ngap_plmn_s {
  uint8_t mcc_digit2 : 4;
  uint8_t mcc_digit1 : 4;
  uint8_t mnc_digit3 : 4;
  uint8_t mcc_digit3 : 4;
  uint8_t mnc_digit2 : 4;
  uint8_t mnc_digit1 : 4;
} Ngap_plmn_t;

typedef uint8_t Ngap_AMF_RegionID_t;
typedef uint16_t
    Ngap_AMF_SetID_t;  // 9.3.3.12 AMF Set ID is used to uniquely identify an
                       // AMF Set within the AMF Region.
typedef uint8_t Ngap_AMF_Pointer_t;  // 9.3.3.19 AMF Pointer is used to identify
                                     // one or more AMF(s) within the AMF Set.

typedef uint16_t NR_Encryption_Algo;
typedef uint16_t NR_Integrity_Protection_Algo;
typedef uint16_t E_UTRA_Encryption_Algo;
typedef uint16_t E_UTRA_Integrity_Protection_Algo;
// 9.3.1.86 UE Security Capabilities
// This IE defines the supported algorithms for encryption and integrity
// protection in the UE.

typedef struct Ngap_ue_security_capabilities_s {
  NR_Encryption_Algo m5g_encryption_algo;
  NR_Integrity_Protection_Algo m5g_integrity_protection_algo;
  E_UTRA_Encryption_Algo e_utra_encryption_algo;
  E_UTRA_Integrity_Protection_Algo e_utra_integrity_protection_algo;
} Ngap_ue_security_capabilities_t;

// 9.2.2 UE Context Management Messages
// 9.2.2.1 INITIAL CONTEXT SETUP REQUEST
// This message is sent by the AMF to request the setup of a UE context.
// Direction: AMF → NG-RAN node
typedef struct Ngap_initial_context_setup_request_s {
  Ngap_Message_Type_t Ngap_Message_Type;
  amf_ue_ngap_id_ty
      amf_ue_ngap_id;  // This IE uniquely identifies the UE association over
                       // the NG interface, as described in TS 38.401
  ran_ue_ngap_id_t
      ran_ue_ngap_id;  // This IE uniquely identifies the UE association over
                       // the NG interface within the NG-RAN node
  guamfi_t Ngap_guami;
  Ngap_PDUSessionID_t
      Pdu_Session_ID;  // PDU Session for a UE. The definition and use of the
                       // PDU Session ID is specified in TS 23.501 [9].
  Ngap_SNSSAI_t Ngap_s_nssai;  // S-NSSAI as defined in TS 23.003 [23].
  bstring
      PDU_Session_Resource_Setup_Transfer;  // Containing the PDU Session
                                            // Resource Setup Request Transfer
                                            // IE specified in subclause 9.3.4.1
  Ngap_SNSSAI_t
      allowed_nssai;  // 9.3.1.31 Allowed NSSAI contains the allowed NSSAI.
  Ngap_ue_security_capabilities_t ue_security_capabilities;
  unsigned char*
      Security_Key;  // 9.3.1.87 Security Key is used to apply security in the
                     // NG-RAN for different scenarios as defined in TS 33.501
  bstring nas_pdu;   // optional

} Ngap_initial_context_setup_request_t;

typedef unsigned char* Transport_Layer_Address_t;
typedef uint32_t Ngap_Gtp_Teid_t;

typedef union GTP_tunnel_s {
  Transport_Layer_Address_t
      endpoint_ip_address;        // Transport Layer Address 9.3.2.4
  Ngap_Gtp_Teid_t ngap_gtp_teid;  // 9.3.2.5 GTP-TEID IE is the GTP Tunnel
                                  // Endpoint Identifier to
                                  // be used for the user plane transport
                                  // between the NG-RAN node and the UPF.
} GTP_tunnel_t;

typedef struct UP_Transport_Layer_Info_s {
  GTP_tunnel_t GTP_tunnel;

} UP_Transport_Layer_Info_t;
// 9.3.1.99 Associated QoS Flow List
// This IE indicates the list of QoS flows associated with e.g. a DRB or UP TNL
// endpoint.

typedef enum QoS_Flow_Mapping_Indi {
  ul = 0,
  dl = 1,
} QoS_Flow_Mapping_Indi_t;

typedef struct Associated_QoS_Flow_List_s {
  uint32_t QoS_Flow_Identifier;  // 9.3.1.51
  QoS_Flow_Mapping_Indi_t QoS_Flow_Mapping_Indi;
} Associated_QoS_Flow_List_t;

// 9.3.2.8 QoS Flow per TNL Information
// This IE indicates the NG-U transport layer information and associated list of
// QoS flows.

typedef struct DL_QoS_Flow_per_TNL_Info_s {
  UP_Transport_Layer_Info_t up_transport_layer_info;    //
  Associated_QoS_Flow_List_t associated_qos_flow_list;  // 9.3.1.99
} DL_QoS_Flow_per_TNL_Info_t;

// 9.3.4.2 PDU Session Resource Setup Response Transfer
// This IE is transparent to the AMF.

typedef struct PDU_Session_Resource_Setup_Response_Transfer_s {
  DL_QoS_Flow_per_TNL_Info_t
      dl_qos_flow_per_tnl_info;  // QoS Flow per TNL Information 9. 3.2.8
  // Additional DL QoS Flow per TNL Information //optional TODO :will be part of
  // PDU Session establishment exchanges

} PDU_Session_Resource_Setup_Response_Transfer_t;

// 9.2.2.2 INITIAL CONTEXT SETUP RESPONSE
// This message is sent by the NG-RAN node to confirm the setup of a UE context.
// Direction: NG-RAN node → AMF
typedef struct Ngap_initial_context_setup_response_s {
  Ngap_Message_Type_t Ngap_Message_Type;
  amf_ue_ngap_id_ty
      amf_ue_ngap_id;  // This IE uniquely identifies the UE association over
                       // the NG interface, as described in TS 38.401
  ran_ue_ngap_id_t
      ran_ue_ngap_id;  // This IE uniquely identifies the UE association over
                       // the NG interface within the NG-RAN node
  Ngap_PDUSessionID_t
      Pdu_Session_ID;  // PDU Session for a UE. The definition and use of the
                       // PDU Session ID is specified in TS 23.501 [9].
  PDU_Session_Resource_Setup_Response_Transfer_t
      PDU_session_resource_setup_res_trans;

} Ngap_initial_context_setup_response_t;

// paging
typedef struct tai_5G_s {
  plmn_t plmn;
  uint32_t tac : 24;
} tai_5G_t;

typedef struct Ngap_TAI_List_For_Paging_s {
  uint16_t no_of_items;
  tai_5G_t tai_list[TRACKING_AREA_IDENTITY_MAX_NUM_OF_TAIS];
} Ngap_TAI_List_For_Paging_t;
