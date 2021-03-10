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

/*! \file common_types.c
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#ifndef FILE_COMMON_TYPES_SEEN
#define FILE_COMMON_TYPES_SEEN

#include <inttypes.h>
#include <arpa/inet.h>
#include <netinet/in.h>
#include <stdint.h>

#include "bstrlib.h"
#include "3gpp_33.401.h"
#include "3gpp_36.401.h"
#include "security_types.h"
#include "common_dim.h"
#include "3gpp_24.008.h"

//------------------------------------------------------------------------------
typedef uint16_t sctp_stream_id_t;
typedef uint32_t sctp_assoc_id_t;
typedef uint32_t sctp_ppid_t;
typedef uint64_t enb_s1ap_id_key_t;
#define MME_APP_ENB_S1AP_ID_KEY(kEy, eNb_Id, eNb_Ue_S1Ap_Id)                   \
  do {                                                                         \
    kEy = (((enb_s1ap_id_key_t) eNb_Id) << 24) | eNb_Ue_S1Ap_Id;               \
  } while (0);
#define MME_APP_ENB_S1AP_ID_KEY2ENB_S1AP_ID(kEy)                               \
  (enb_ue_s1ap_id_t)(((enb_s1ap_id_key_t) kEy) & ENB_UE_S1AP_ID_MASK)
#define MME_APP_ENB_S1AP_ID_KEY_FORMAT "0x%16" PRIX64

#define M_TMSI_BIT_MASK UINT32_MAX

//------------------------------------------------------------------------------
// UE S1AP IDs

#define INVALID_ENB_UE_S1AP_ID_KEY 0xFFFFFFFFFFFFFFFF
#define ENB_UE_S1AP_ID_MASK 0x00FFFFFF
#define ENB_UE_S1AP_ID_FMT "0x%06" PRIX32

#define MME_UE_S1AP_ID_FMT "0x%08" PRIX32

#define COMP_S1AP_ID_FMT "0x%016" PRIX64

/* INVALID_MME_UE_S1AP_ID
 * Any value between 0..2^32-1, is allowed/valid as per 3GPP spec 36.413.
 * Here we are conisdering 0 as invalid. Don't allocate 0 and consider this as
 * invalid
 */
#define INVALID_MME_UE_S1AP_ID 0x0
#define INVALID_ENB_UE_S1AP_ID 0x0

//------------------------------------------------------------------------------
// TEIDs
typedef uint32_t teid_t;
#define TEID_FMT "0x%" PRIX32
#define TEID_SCAN_FMT SCNx32
typedef teid_t s11_teid_t;
typedef teid_t s1u_teid_t;
#define INVALID_TEID 0x00000000

//------------------------------------------------------------------------------
// IMSI

typedef uint64_t imsi64_t;
#define IMSI_64_FMT "%" SCNu64
#define IMSI_64_FMT_DYN_LEN "%.*lu"
#define INVALID_IMSI64 (imsi64_t) 0

//------------------------------------------------------------------------------
// PLMN

//------------------------------------------------------------------------------

#define LAC_FMT "0x%04X"
/* Checks LAC validity */
#define LAC_IS_VALID(lac)                                                      \
  (((lac) != INVALID_LAC_0000) && ((lac) != INVALID_LAC_FFFE))

//------------------------------------------------------------------------------
// GUTI
#define GUTI_FMT PLMN_FMT "|%04x|%02x|%08x"
#define GUTI_ARG(GuTi_PtR)                                                     \
  PLMN_ARG(&(GuTi_PtR)->gummei.plmn), (GuTi_PtR)->gummei.mme_gid,              \
      (GuTi_PtR)->gummei.mme_code, (GuTi_PtR)->m_tmsi
#define MSISDN_LENGTH (15)
#define IMEI_DIGITS_MAX (15)
#define IMEISV_DIGITS_MAX (16)
#define CHARGING_CHARACTERISTICS_LENGTH (4)  // 3GPP TS 29.061
#define APN_MAX_LENGTH (100)
#define PRIORITY_LEVEL_MAX (15)
#define PRIORITY_LEVEL_MIN (1)
#define BEARERS_PER_UE (11)
#define MAX_APN_PER_UE (10)

//------------------------------------------------------------------------------
// IPv6 Interface Identifier length in bytes
#define IPV6_INTERFACE_ID_LEN 8
// IPv6 Prefix length in bits
#define IPV6_PREFIX_LEN 64
//------------------------------------------------------------------------------
typedef uint8_t ksi_t;
#define KSI_NO_KEY_AVAILABLE 0x07

typedef uint8_t AcT_t; /* Access Technology    */

typedef enum {
  RAT_WLAN           = 0,
  RAT_VIRTUAL        = 1,
  RAT_UTRAN          = 1000,
  RAT_GERAN          = 1001,
  RAT_GAN            = 1002,
  RAT_HSPA_EVOLUTION = 1003,
  RAT_EUTRAN         = 1004,
  RAT_CDMA2000_1X    = 2000,
  RAT_HRPD           = 2001,
  RAT_UMB            = 2002,
  RAT_EHRPD          = 2003,
} rat_type_t;

#define NUMBER_OF_RAT_TYPE 11

typedef enum {
  SS_SERVICE_GRANTED             = 0,
  SS_OPERATOR_DETERMINED_BARRING = 1,
  SS_MAX,
} subscriber_status_t;

typedef enum {
  NAM_PACKET_AND_CIRCUIT = 0,
  NAM_RESERVED           = 1,
  NAM_ONLY_PACKET        = 2,
  NAM_MAX,
} network_access_mode_t;

typedef uint64_t bitrate_t;

typedef char* APN_t;
typedef uint8_t APNRestriction_t;
typedef uint8_t DelayValue_t;
typedef uint8_t priority_level_t;
#define PRIORITY_LEVEL_FMT "0x%" PRIu8
#define PRIORITY_LEVEL_SCAN_FMT SCNu8
typedef uint32_t SequenceNumber_t;
typedef uint32_t access_restriction_t;
typedef uint32_t context_identifier_t;
typedef uint32_t rau_tau_timer_t;

typedef uint32_t ard_t;
typedef int pdn_cid_t;  // pdn connexion identity, related to esm protocol,
                        // sometimes type is mixed with int return code!!...
typedef uint8_t
    proc_tid_t;  // procedure transaction identity, related to esm protocol
#define ARD_UTRAN_NOT_ALLOWED (1U)
#define ARD_GERAN_NOT_ALLOWED (1U << 1)
#define ARD_GAN_NOT_ALLOWED (1U << 2)
#define ARD_I_HSDPA_EVO_NOT_ALLOWED (1U << 3)
#define ARD_E_UTRAN_NOT_ALLOWED (1U << 4)
#define ARD_HO_TO_NON_3GPP_NOT_ALLOWED (1U << 5)
#define ARD_MAX (1U << 6)

typedef union {
  uint8_t imei[IMEI_DIGITS_MAX - 1];  // -1 =  remove CD/SD digit
  uint8_t imeisv[IMEISV_DIGITS_MAX];
} me_identity_t;

typedef struct {
  bitrate_t br_ul;
  bitrate_t br_dl;
} ambr_t;

typedef uint8_t pdn_type_t;

typedef enum {
  IPv4        = 0,
  IPv6        = 1,
  IPv4_AND_v6 = 2,
  IPv4_OR_v6  = 3,
  IP_MAX,
} pdn_type_value_t;

typedef struct paa_s {
  pdn_type_value_t pdn_type;
  struct in_addr ipv4_address;
  struct in6_addr ipv6_address;
  /* Note in rel.8 the ipv6 prefix length has a fixed value of /64 */
  uint8_t ipv6_prefix_length;
  int vlan;
} paa_t;

void copy_paa(paa_t* paa_dst, paa_t* paa_src);
bstring paa_to_bstring(const paa_t* paa);

//-----------------
typedef struct {
  pdn_type_value_t pdn_type;
  struct {
    struct in_addr ipv4_address;
    struct in6_addr ipv6_address;
  } address;
} ip_address_t;

struct fteid_s;
bstring fteid_ip_address_to_bstring(const struct fteid_s* const fteid);
void get_fteid_ip_address(
    const struct fteid_s* const fteid, ip_address_t* const ip_address);
bstring ip_address_to_bstring(const ip_address_t* ip_address);
void bstring_to_ip_address(bstring const bstr, ip_address_t* const ip_address);
void bstring_to_paa(bstring bstr, paa_t* paa);

//-----------------
typedef enum {
  QCI_1 = 1,
  QCI_2 = 2,
  QCI_3 = 3,
  QCI_4 = 4,
  QCI_5 = 5,
  QCI_6 = 6,
  QCI_7 = 7,
  QCI_8 = 8,
  QCI_9 = 9,
  /* Values from 128 to 254 are operator specific.
   * Other are reserved.
   */
  QCI_MAX,
} qci_e;

typedef uint8_t qci_t;
#define QCI_FMT "0x%" PRIu8
#define QCI_SCAN_FMT SCNu8

typedef enum {
  PRE_EMPTION_CAPABILITY_ENABLED  = 0,
  PRE_EMPTION_CAPABILITY_DISABLED = 1,
  PRE_EMPTION_CAPABILITY_MAX,
} pre_emption_capability_t;

#define PRE_EMPTION_CAPABILITY_FMT "0x%" PRIu8
#define PRE_EMPTION_CAPABILITY_SCAN_FMT SCNu8

typedef enum {
  PRE_EMPTION_VULNERABILITY_ENABLED  = 0,
  PRE_EMPTION_VULNERABILITY_DISABLED = 1,
  PRE_EMPTION_VULNERABILITY_MAX,
} pre_emption_vulnerability_t;

#define PRE_EMPTION_VULNERABILITY_FMT "0x%" PRIu8
#define PRE_EMPTION_VULNERABILITY_SCAN_FMT SCNu8

typedef struct {
  priority_level_t priority_level;
  pre_emption_vulnerability_t pre_emp_vulnerability;
  pre_emption_capability_t pre_emp_capability;
} allocation_retention_priority_t;

typedef struct eps_subscribed_qos_profile_s {
  qci_t qci;
  allocation_retention_priority_t allocation_retention_priority;
} eps_subscribed_qos_profile_t;

typedef struct {
  char value[CHARGING_CHARACTERISTICS_LENGTH + 1];
  size_t length;
} charging_characteristics_t;

typedef struct apn_configuration_s {
  context_identifier_t context_identifier;

  /* Each APN configuration can have 0, 1, or 2 ip address:
   * - 0 means subscribed is dynamically allocated by P-GW depending on the
   * pdn_type
   * - 1 Only one type of IP address is returned by HSS
   * - 2 IPv4 and IPv6 address are returned by HSS and are statically
   * allocated
   */
  uint8_t nb_ip_address;
  ip_address_t ip_address[2];

#ifdef ACCESS_POINT_NAME_MAX_LENGTH
#define SERVICE_SELECTION_MAX_LENGTH ACCESS_POINT_NAME_MAX_LENGTH
#else
#define SERVICE_SELECTION_MAX_LENGTH 100
#endif
  pdn_type_t pdn_type;
  char service_selection[SERVICE_SELECTION_MAX_LENGTH];
  int service_selection_length;
  eps_subscribed_qos_profile_t subscribed_qos;
  ambr_t ambr;
  charging_characteristics_t charging_characteristics;
} apn_configuration_t;

typedef enum {
  ALL_APN_CONFIGURATIONS_INCLUDED            = 0,
  MODIFIED_ADDED_APN_CONFIGURATIONS_INCLUDED = 1,
  ALL_APN_MAX,
} all_apn_conf_ind_t;

typedef struct {
  context_identifier_t context_identifier;
  all_apn_conf_ind_t all_apn_conf_ind;
  /* Number of APNs provided */
  uint8_t nb_apns;
  /* List of APNs configuration 1 to n elements */
  struct apn_configuration_s apn_configuration[MAX_APN_PER_UE];
} apn_config_profile_t;

typedef struct {
  subscriber_status_t subscriber_status;
  char msisdn[MSISDN_LENGTH + 1];
  uint8_t msisdn_length;
  network_access_mode_t access_mode;
  access_restriction_t access_restriction;
  ambr_t subscribed_ambr;
  apn_config_profile_t apn_config_profile;
  rau_tau_timer_t rau_tau_timer;
  charging_characteristics_t default_charging_characteristics;
} subscription_data_t;

typedef struct authentication_info_s {
  uint8_t nb_of_vectors;
  eutran_vector_t eutran_vector[MAX_EPS_AUTH_VECTORS];
} authentication_info_t;

typedef enum {
  DIAMETER_AUTHENTICATION_DATA_UNAVAILABLE = 4181,
  DIAMETER_ERROR_USER_UNKNOWN              = 5001,
  DIAMETER_ERROR_ROAMING_NOT_ALLOWED       = 5004,
  DIAMETER_ERROR_UNKNOWN_EPS_SUBSCRIPTION  = 5420,
  DIAMETER_ERROR_RAT_NOT_ALLOWED           = 5421,
  DIAMETER_ERROR_EQUIPMENT_UNKNOWN         = 5422,
  DIAMETER_ERROR_UNKOWN_SERVING_NODE       = 5423,
} s6a_experimental_result_t;

typedef enum {
  DIAMETER_SUCCESS          = 2001,
  DIAMETER_UNABLE_TO_COMPLY = 5012,
} s6a_base_result_t;

typedef struct {
#define S6A_RESULT_BASE 0x0
#define S6A_RESULT_EXPERIMENTAL 0x1
  unsigned present : 1;

  union {
    /* Experimental result as defined in 3GPP TS 29.272 */
    s6a_experimental_result_t experimental;
    /* Diameter basics results as defined in RFC 3588 */
    s6a_base_result_t base;
  } choice;
} s6a_result_t;

typedef enum {
  MME_UPDATE_PROCEDURE = 0,
  SGSN_UPDATE_PROCEDURE,
  SUBSCRIPTION_WITHDRAWL,
  UPDATE_PROCEDURE_IWF,
  INITIAL_ATTACH_PROCEDURE
} s6a_cancellation_type_t;

#include "nas/commonDef.h"

struct fteid_s;

#endif /* FILE_COMMON_TYPES_SEEN */
