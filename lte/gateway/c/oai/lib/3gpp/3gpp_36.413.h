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

/*! \file 3gpp_36.413.h
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#ifndef FILE_3GPP_36_413_SEEN
#define FILE_3GPP_36_413_SEEN

#include "3gpp_24.007.h"
#include "3gpp_29.274.h"
#include "common_types.h"

// 9.2.1.60 Allocation and Retention Priority
// This IE specifies the relative importance compared to other E-RABs for
// allocation and retention of the E-UTRAN Radio Access Bearer.
typedef struct allocation_and_retention_priority_s {
  priority_level_t priority_level;  // INTEGER (0..15)
  pre_emption_capability_t pre_emption_capability;
  pre_emption_vulnerability_t pre_emption_vulnerability;
} allocation_and_retention_priority_t;

// 9.2.1.19 Bit Rate
// This IE indicates the number of bits delivered by E-UTRAN in UL or to E-UTRAN
// in DL within a period of time, divided by the duration of the period. It is
// used, for example, to indicate the maximum or guaranteed bit rate for a GBR
// bearer, or an aggregated maximum bit rate.

typedef uint64_t bit_rate_t;

// 9.2.1.20 UE Aggregate Maximum Bit Rate
// The UE Aggregate Maximum Bitrate is applicable for all Non-GBR bearers per UE
// which is defined for the Downlink
// and the Uplink direction and provided by the MME to the eNB.
// Applicable for non-GBR E-RABs
typedef struct ue_aggregate_maximum_bit_rate_s {
  bit_rate_t dl;
  bit_rate_t ul;
} ue_aggregate_maximum_bit_rate_t;

// 9.2.1.18 GBR QoS Information
// This IE indicates the maximum and guaranteed bit rates of a GBR bearer for
// downlink and uplink.
typedef struct gbr_qos_information_s {
  bit_rate_t e_rab_maximum_bit_rate_downlink;
  bit_rate_t e_rab_maximum_bit_rate_uplink;
  bit_rate_t e_rab_guaranteed_bit_rate_downlink;
  bit_rate_t e_rab_guaranteed_bit_rate_uplink;
} gbr_qos_information_t;

// 9.2.1.15 E-RAB Level QoS Parameters
// This IE defines the QoS to be applied to an E-RAB.
typedef struct e_rab_level_qos_parameters_s {
  qci_t qci;
  allocation_and_retention_priority_t allocation_and_retention_priority;
  gbr_qos_information_t gbr_qos_information;
} e_rab_level_qos_parameters_t;

// 9.2.1.2 E-RAB ID
typedef ebi_t e_rab_id_t;

// 9.1.3.1 E-RAB SETUP REQUEST
typedef struct e_rab_to_be_setup_item_s {
  e_rab_id_t e_rab_id;
  e_rab_level_qos_parameters_t e_rab_level_qos_parameters;
  bstring transport_layer_address;
  teid_t gtp_teid;
  bstring nas_pdu;
  // Correlation ID TODO if necessary
} e_rab_to_be_setup_item_t;

typedef struct e_rab_to_be_setup_list_s {
  uint16_t no_of_items;
#define MAX_NO_OF_E_RABS 16 /* Spec says 256 */
  e_rab_to_be_setup_item_t item[MAX_NO_OF_E_RABS];
} e_rab_to_be_setup_list_t;

// 9.1.3.2 E-RAB SETUP RESPONSE
typedef struct e_rab_setup_item_s {
  e_rab_id_t e_rab_id;
  bstring transport_layer_address;
  teid_t gtp_teid;
} e_rab_setup_item_t;

typedef struct e_rab_setup_list_s {
  uint16_t no_of_items;
  e_rab_setup_item_t item[MAX_NO_OF_E_RABS];
} e_rab_setup_list_t;

typedef struct e_rab_rel_item_s {
  e_rab_id_t e_rab_id;
} e_rab_rel_item_t;

typedef struct e_rab_rel_list_s {
  uint16_t no_of_items;
  e_rab_rel_item_t item[MAX_NO_OF_E_RABS];
} e_rab_rel_list_t;

typedef struct e_rab_switched_in_downlink_item_s {
  e_rab_id_t e_rab_id;
  bstring transport_layer_address;
  teid_t gtp_teid;
} e_rab_switched_in_downlink_item_t;

// 9.1.5.8 PATH SWITCH REQUEST
typedef struct e_rab_to_be_switched_in_downlink_list_s {
  uint16_t no_of_items;
#define MAX_NO_OF_E_RABS 16 /* Spec says 256 */
  e_rab_switched_in_downlink_item_t item[MAX_NO_OF_E_RABS];
} e_rab_to_be_switched_in_downlink_list_t;

// 9.1.5.4 HANDOVER REQUEST
typedef struct e_rab_to_be_setup_item_ho_req_s {
  e_rab_id_t e_rab_id;
  bstring transport_layer_address;
  teid_t gtp_teid;
  e_rab_level_qos_parameters_t e_rab_level_qos_parameters;
  // TODO: Include optional data-forwarding-not-possible IE
} e_rab_to_be_setup_item_ho_req_t;

typedef struct e_rab_to_be_setup_list_ho_req_s {
  uint16_t no_of_items;
  e_rab_to_be_setup_item_ho_req_t item[MAX_NO_OF_E_RABS];
} e_rab_to_be_setup_list_ho_req_t;

// 9.1.5.5 HANDOVER REQUEST ACK
typedef struct e_rab_admitted_item_s {
  e_rab_id_t e_rab_id;
  bstring transport_layer_address;
  teid_t gtp_teid;
  // TODO: Include optional UL and DL tunnels for indirect forwarding
} e_rab_admitted_item_t;

typedef struct e_rab_admitted_list_s {
  uint16_t no_of_items;
  e_rab_admitted_item_t item[MAX_NO_OF_E_RABS];
} e_rab_admitted_list_t;

// E-RAB TO BE MODIFIED ITEM BEARER MOD IND
typedef struct e_rab_to_be_modified_bearer_mod_ind_s {
  e_rab_id_t e_rab_id;
  fteid_t s1_xNB_fteid;  ///< S1 xNodeB F-TEID
} e_rab_to_be_modified_bearer_mod_ind_t;

typedef struct e_rab_not_to_be_modified_bearer_mod_ind_s {
  e_rab_id_t e_rab_id;
  fteid_t s1_xNB_fteid;  ///< S1 xNodeB F-TEID
} e_rab_not_to_be_modified_bearer_mod_ind_t;

typedef struct e_rab_to_be_modified_bearer_mod_ind_list_s {
  uint16_t no_of_items;
  e_rab_to_be_modified_bearer_mod_ind_t item[MAX_NO_OF_E_RABS];
} e_rab_to_be_modified_bearer_mod_ind_list_t;

typedef struct e_rab_not_to_be_modified_bearer_mod_ind_list_s {
  uint16_t no_of_items;
  e_rab_not_to_be_modified_bearer_mod_ind_t item[MAX_NO_OF_E_RABS];
} e_rab_not_to_be_modified_bearer_mod_ind_list_t;

typedef struct e_rab_modify_bearer_mod_conf_list_s {
  uint16_t no_of_items;
  e_rab_id_t e_rab_id[MAX_NO_OF_E_RABS];
} e_rab_modify_bearer_mod_conf_list_t;

#include "S1ap_Cause.h"

typedef struct e_rab_item_s {
  e_rab_id_t e_rab_id;
  S1ap_Cause_t cause;
} e_rab_item_t;

typedef struct e_rab_list_s {
  uint16_t no_of_items;
  e_rab_item_t item[MAX_NO_OF_E_RABS];
} e_rab_list_t;

#endif /* FILE_3GPP_36_413_SEEN */
