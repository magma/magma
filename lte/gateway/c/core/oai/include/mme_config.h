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
/*! \file mme_config.h
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#ifndef FILE_MME_CONFIG_SEEN
#define FILE_MME_CONFIG_SEEN

#include <pthread.h>
#include <stdint.h>
#include <arpa/inet.h>
#include <stdlib.h>
#include "lte/gateway/c/core/oai/common/mme_default_values.h"
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_23.003.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.008.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/include/service303.h"
#include "lte/gateway/c/core/oai/lib/hashtable/hashtable.h"
#include "lte/gateway/c/core/oai/lib/hashtable/obj_hashtable.h"
#include "orc8r/gateway/c/common/sentry/includes/SentryWrapper.h"

/* Currently supporting max 5 GUMMEI's in the mme configuration */
#define MIN_GUMMEI 1
#define MAX_GUMMEI 5
#define MIN_TAI_SUPPORTED 1
#define MAX_TAI_SUPPORTED 16
#define MAX_MCC_LENGTH 3
#define MAX_MNC_LENGTH 3
#define MIN_MNC_LENGTH 2
#define CIDR_SPLIT_LIST_COUNT 2
#define MAX_APN_CORRECTION_MAP_LIST 10
#define MAX_RESTRICTED_PLMN 10
#define MAX_LEN_TAC 8
#define MAX_LEN_SNR 6
#define MAX_LEN_IMEI 15
#define MAX_FED_MODE_MAP_CONFIG 10
#define MAX_IMSI_LENGTH 15
#define MAX_IMEI_HTBL_SZ 32
#define MAX_SAC_2_TACS_HTBL_SZ 32

#define MME_CONFIG_STRING_MME_CONFIG "MME"
#define MME_CONFIG_STRING_PID_DIRECTORY "PID_DIRECTORY"
#define MME_CONFIG_STRING_RUN_MODE "RUN_MODE"
#define MME_CONFIG_STRING_RUN_MODE_TEST "TEST"
#define MME_CONFIG_STRING_REALM "REALM"
#define MME_CONFIG_STRING_MAXENB "MAXENB"
#define MME_CONFIG_STRING_MAXUE "MAXUE"
#define MME_CONFIG_STRING_RELATIVE_CAPACITY "RELATIVE_CAPACITY"
#define MME_CONFIG_STRING_STATS_TIMER "STATS_TIMER_SEC"

#define MME_CONFIG_STRING_USE_STATELESS "USE_STATELESS"
#define MME_CONFIG_STRING_ENABLE5G_FEATURES "ENABLE5G_FEATURES"
#define MME_CONFIG_STRING_FULL_NETWORK_NAME "FULL_NETWORK_NAME"
#define MME_CONFIG_STRING_SHORT_NETWORK_NAME "SHORT_NETWORK_NAME"
#define MME_CONFIG_STRING_DAYLIGHT_SAVING_TIME "DAYLIGHT_SAVING_TIME"
#define MME_CONFIG_STRING_CSFB_CONFIG "CSFB"
#define MME_CONFIG_STRING_NON_EPS_SERVICE_CONTROL "NON_EPS_SERVICE_CONTROL"

#define MME_CONFIG_STRING_EMERGENCY_ATTACH_SUPPORTED                           \
  "EMERGENCY_ATTACH_SUPPORTED"
#define MME_CONFIG_STRING_UNAUTHENTICATED_IMSI_SUPPORTED                       \
  "UNAUTHENTICATED_IMSI_SUPPORTED"

#define EPS_NETWORK_FEATURE_SUPPORT_IMS_VOICE_OVER_PS_SESSION_IN_S1            \
  "EPS_NETWORK_FEATURE_SUPPORT_IMS_VOICE_OVER_PS_SESSION_IN_S1"
#define EPS_NETWORK_FEATURE_SUPPORT_EMERGENCY_BEARER_SERVICES_IN_S1_MODE       \
  "EPS_NETWORK_FEATURE_SUPPORT_EMERGENCY_BEARER_SERVICES_IN_S1_MODE"
#define EPS_NETWORK_FEATURE_SUPPORT_LOCATION_SERVICES_VIA_EPC                  \
  "EPS_NETWORK_FEATURE_SUPPORT_LOCATION_SERVICES_VIA_EPC"
#define EPS_NETWORK_FEATURE_SUPPORT_EXTENDED_SERVICE_REQUEST                   \
  "EPS_NETWORK_FEATURE_SUPPORT_EXTENDED_SERVICE_REQUEST"

#define MME_CONFIG_STRING_ACCEPT_COMBINED_ATTACH_TAU_WO_CSFB                   \
  "ACCEPT_COMBINED_ATTACH_TAU_WO_CSFB"

#define MME_CONFIG_STRING_INTERTASK_INTERFACE_CONFIG "INTERTASK_INTERFACE"
#define MME_CONFIG_STRING_INTERTASK_INTERFACE_QUEUE_SIZE "ITTI_QUEUE_SIZE"

#define MME_CONFIG_STRING_S6A_CONFIG "S6A"
#define MME_CONFIG_STRING_S6A_CONF_FILE_PATH "S6A_CONF"
#define MME_CONFIG_STRING_S6A_HSS_HOSTNAME "HSS_HOSTNAME"
#define MME_CONFIG_STRING_S6A_HSS_REALM "HSS_REALM"

#define MME_CONFIG_STRING_SCTP_CONFIG "SCTP"
#define MME_CONFIG_STRING_SCTP_UPSTREAM_SOCK "SCTP_UPSTREAM_SOCK"
#define MME_CONFIG_STRING_SCTP_UPSTREAM_SOCK_DEFAULT                           \
  "unix:///tmp/sctpd_upstream.sock"
#define MME_CONFIG_STRING_SCTP_DOWNSTREAM_SOCK "SCTP_DOWNSTREAM_SOCK"
#define MME_CONFIG_STRING_SCTP_DOWNSTREAM_SOCK_DEFAULT                         \
  "unix:///tmp/sctpd_downstream.sock"

#define MME_CONFIG_STRING_S1AP_CONFIG "S1AP"
#define MME_CONFIG_STRING_S1AP_OUTCOME_TIMER "S1AP_OUTCOME_TIMER"
#define MME_CONFIG_STRING_S1AP_PORT "S1AP_PORT"

#define MME_CONFIG_STRING_TAC_LIST "TAC_LIST"
#define MME_CONFIG_STRING_GUAMFI_LIST "GUAMFI_LIST"
#define MME_CONFIG_STRING_AMF_REGION_ID "AMF_REGION_ID"
#define MME_CONFIG_STRING_AMF_SET_ID "AMF_SET_ID"
#define MME_CONFIG_STRING_AMF_POINTER "AMF_POINTER"

#define MME_CONFIG_STRING_GUMMEI_LIST "GUMMEI_LIST"
#define MME_CONFIG_STRING_MME_CODE "MME_CODE"
#define MME_CONFIG_STRING_MME_GID "MME_GID"
#define MME_CONFIG_STRING_TAI_LIST "TAI_LIST"
#define MME_CONFIG_STRING_MCC "MCC"
#define MME_CONFIG_STRING_MNC "MNC"
#define MME_CONFIG_STRING_TAC "TAC"

#define MME_CONFIG_STRING_RESTRICTED_PLMN_LIST "RESTRICTED_PLMN_LIST"
#define MME_CONFIG_STRING_BLOCKED_IMEI_LIST "BLOCKED_IMEI_LIST"
#define MME_CONFIG_STRING_IMEI_TAC "IMEI_TAC"
#define MME_CONFIG_STRING_SNR "SNR"

#define MME_CONFIG_STRING_NETWORK_INTERFACES_CONFIG "NETWORK_INTERFACES"
#define MME_CONFIG_STRING_INTERFACE_NAME_FOR_S1_MME                            \
  "MME_INTERFACE_NAME_FOR_S1_MME"
#define MME_CONFIG_STRING_IPV4_ADDRESS_FOR_S1_MME "MME_IPV4_ADDRESS_FOR_S1_MME"
#define MME_CONFIG_STRING_IPV6_ADDRESS_FOR_S1_MME "MME_IPV6_ADDRESS_FOR_S1_MME"

#define MME_CONFIG_STRING_S1_IPV6_ENABLED "MME_S1_IPV6_ENABLED"
#define MME_CONFIG_STRING_INTERFACE_NAME_FOR_S11_MME                           \
  "MME_INTERFACE_NAME_FOR_S11_MME"
#define MME_CONFIG_STRING_IPV4_ADDRESS_FOR_S11_MME                             \
  "MME_IPV4_ADDRESS_FOR_S11_MME"
#define MME_CONFIG_STRING_MME_PORT_FOR_S11 "MME_PORT_FOR_S11_MME"
#define MME_CONFIG_STRING_SGW_INTERFACE_NAME_FOR_S11                           \
  "SGW_INTERFACE_NAME_FOR_S11"
#define MME_CONFIG_STRING_SGW_IPV4_ADDRESS_FOR_S11 "SGW_IPV4_ADDRESS_FOR_S11"

#define MME_CONFIG_STRING_NAS_CONFIG "NAS"
#define MME_CONFIG_STRING_NAS_SUPPORTED_INTEGRITY_ALGORITHM_LIST               \
  "ORDERED_SUPPORTED_INTEGRITY_ALGORITHM_LIST"
#define MME_CONFIG_STRING_NAS_SUPPORTED_CIPHERING_ALGORITHM_LIST               \
  "ORDERED_SUPPORTED_CIPHERING_ALGORITHM_LIST"

#define MME_CONFIG_STRING_NAS_T3402_TIMER "T3402"
#define MME_CONFIG_STRING_NAS_T3412_TIMER "T3412"
#define MME_CONFIG_STRING_NAS_T3422_TIMER "T3422"
#define MME_CONFIG_STRING_NAS_T3450_TIMER "T3450"
#define MME_CONFIG_STRING_NAS_T3460_TIMER "T3460"
#define MME_CONFIG_STRING_NAS_T3470_TIMER "T3470"
#define MME_CONFIG_STRING_NAS_T3485_TIMER "T3485"
#define MME_CONFIG_STRING_NAS_T3486_TIMER "T3486"
#define MME_CONFIG_STRING_NAS_T3489_TIMER "T3489"
#define MME_CONFIG_STRING_NAS_T3495_TIMER "T3495"
#define MME_CONFIG_STRING_NAS_FORCE_REJECT_TAU "FORCE_REJECT_TAU"
#define MME_CONFIG_STRING_NAS_FORCE_REJECT_SR "FORCE_REJECT_SR"
#define MME_CONFIG_STRING_NAS_DISABLE_ESM_INFORMATION_PROCEDURE                \
  "DISABLE_ESM_INFORMATION_PROCEDURE"
#define MME_CONFIG_STRING_NAS_FORCE_PUSH_DEDICATED_BEARER                      \
  "FORCE_PUSH_DEDICATED_BEARER"
#define MME_CONFIG_STRING_NAS_ENABLE_APN_CORRECTION "ENABLE_APN_CORRECTION"
#define MME_CONFIG_STRING_NAS_APN_CORRECTION_MAP_LIST "APN_CORRECTION_MAP_LIST"
#define MME_CONFIG_STRING_NAS_APN_CORRECTION_MAP_IMSI_PREFIX                   \
  "APN_CORRECTION_MAP_IMSI_PREFIX"
#define MME_CONFIG_STRING_NAS_APN_CORRECTION_MAP_APN_OVERRIDE                  \
  "APN_CORRECTION_MAP_APN_OVERRIDE"

#define MME_CONFIG_STRING_SGW_CONFIG "S-GW"

#define MME_CONFIG_STRING_SGS_CONFIG "SGS"
#define MME_CONFIG_STRING_SGS_TS6_1_TIMER "TS6_1"
#define MME_CONFIG_STRING_SGS_TS8_TIMER "TS8"
#define MME_CONFIG_STRING_SGS_TS9_TIMER "TS9"
#define MME_CONFIG_STRING_SGS_TS10_TIMER "TS10"
#define MME_CONFIG_STRING_SGS_TS13_TIMER "TS13"

#define MME_CONFIG_STRING_ASN1_VERBOSITY "ASN1_VERBOSITY"
#define MME_CONFIG_STRING_ASN1_VERBOSITY_NONE "none"
#define MME_CONFIG_STRING_ASN1_VERBOSITY_ANNOYING "annoying"
#define MME_CONFIG_STRING_ASN1_VERBOSITY_INFO "info"
#define MME_CONFIG_STRING_SGW_LIST_SELECTION "S-GW_LIST_SELECTION"
#define MME_CONFIG_STRING_ID "ID"

#define MAGMA_CONFIG_STRING "MAGMA"
#define MME_CONFIG_STRING_SERVICE303_CONFIG "SERVICE303"
#define MME_CONFIG_STRING_SERVICE303_CONF_SERVER_ADDRESS "SERVER_ADDRESS"
// CSFB
#define MME_CONFIG_STRING_CSFB_MCC "CSFB_MCC"
#define MME_CONFIG_STRING_CSFB_MNC "CSFB_MNC"
#define MME_CONFIG_STRING_LAC "LAC"

// HA
#define MME_CONFIG_STRING_USE_HA "USE_HA"
// Cloud Instances may utilize this to reach RAN behind NAT
#define MME_CONFIG_STRING_ENABLE_GTPU_PRIVATE_IP_CORRECTION                    \
  "ENABLE_GTPU_PRIVATE_IP_CORRECTION"

// Congestion Control
#define MME_CONFIG_STRING_CONGESTION_CONTROL_ENABLED                           \
  "CONGESTION_CONTROL_ENABLED"
#define MME_CONFIG_STRING_S1AP_ZMQ_TH "S1AP_ZMQ_TH"
#define MME_CONFIG_STRING_MME_APP_ZMQ_CONGEST_TH "MME_APP_ZMQ_CONGEST_TH"
#define MME_CONFIG_STRING_MME_APP_ZMQ_AUTH_TH "MME_APP_ZMQ_AUTH_TH"
#define MME_CONFIG_STRING_MME_APP_ZMQ_IDENT_TH "MME_APP_ZMQ_IDENT_TH"
#define MME_CONFIG_STRING_MME_APP_ZMQ_SMC_TH "MME_APP_ZMQ_SMC_TH"

// INBOUND ROAMING
#define MME_CONFIG_STRING_FED_MODE_MAP "FEDERATED_MODE_MAP"
#define MME_CONFIG_STRING_MODE "MODE"
#define MME_CONFIG_STRING_APN "APN"
#define MME_CONFIG_STRING_IMSI_RANGE "IMSI_RANGE"
#define MME_CONFIG_STRING_PLMN "PLMN"
#define MME_CONFIG_STRING_SERVICE_AREA_CODE "SAC"
#define MME_CONFIG_STRING_TAC_LIST_PER_SAC "TACS_PER_SAC"
#define MME_CONFIG_STRING_SRVC_AREA_CODE_2_TACS_MAP "SRVC_AREA_CODE_2_TACS_MAP"

// SENTRY CONFIGURATION
#define MME_CONFIG_STRING_SENTRY_CONFIG "SENTRY_CONFIG"
#define MME_CONFIG_STRING_SAMPLE_RATE "SAMPLE_RATE"
#define MME_CONFIG_STRING_UPLOAD_MME_LOG "UPLOAD_MME_LOG"
#define MME_CONFIG_STRING_URL_NATIVE "URL_NATIVE"

typedef enum { RUN_MODE_TEST = 0, RUN_MODE_OTHER } run_mode_t;

typedef struct eps_network_feature_config_s {
  uint8_t ims_voice_over_ps_session_in_s1;
  uint8_t emergency_bearer_services_in_s1_mode;
  uint8_t location_services_via_epc;
  uint8_t extended_service_request;
} eps_network_feature_config_t;

#define TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_NON_CONSECUTIVE_TACS 0x00
#define TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_CONSECUTIVE_TACS 0x01
#define TRACKING_AREA_IDENTITY_LIST_TYPE_MANY_PLMNS 0x02

typedef struct sctp_config_s {
  bstring upstream_sctp_sock;
  bstring downstream_sctp_sock;
} sctp_config_t;

typedef struct s1ap_config_s {
  uint16_t port_number;
  uint8_t outcome_drop_timer_sec;
} s1ap_config_t;

typedef struct ip_s {
  bstring if_name_s1_mme;
  struct in_addr s1_mme_v4;
  struct in6_addr s1_mme_v6;
  bool s1_ipv6_enabled;
  int netmask_s1_mme;

  bstring if_name_s11;
  struct in_addr s11_mme_v4;
  struct in6_addr s11_mme_v6;
  int netmask_s11;
  uint16_t port_s11;
} ip_t;

typedef struct s6a_config_s {
  bstring conf_file;
  bstring hss_host_name;
  bstring hss_realm;
} s6a_config_t;

typedef struct itti_config_s {
  uint32_t queue_size;
  bstring log_file;
} itti_config_t;

typedef struct apn_map_s {
  bstring imsi_prefix;
  bstring apn_override;
} apn_map_t;

typedef struct apn_map_config_s {
  int nb;
  apn_map_t apn_map[MAX_APN_CORRECTION_MAP_LIST];
} apn_map_config_t;

typedef struct nas_config_s {
  uint8_t prefered_integrity_algorithm[8];
  uint8_t prefered_ciphering_algorithm[8];
  uint32_t t3402_min;
  uint32_t t3412_min;
  uint32_t t3412_msec;  // keeping t3412_min as it is used to communicate to UE
                        // and to prevent back and forth conversions
  uint32_t t3422_msec;
  uint32_t t3450_msec;
  uint32_t t3460_msec;
  uint32_t t3470_msec;
  uint32_t t3485_msec;
  uint32_t t3486_msec;
  uint32_t t3489_msec;
  uint32_t t3495_msec;
  uint32_t ts6a_msec;
  uint32_t tics_msec;
  uint32_t tpaging_msec;
  // non standard features
  bool force_reject_tau;
  bool force_reject_sr;
  bool disable_esm_information;
  // apn correction
  bool enable_apn_correction;
  apn_map_config_t apn_map_config;
} nas_config_t;

typedef struct sgs_config_s {
  uint32_t ts6_1_msec;
  uint32_t ts8_msec;
  uint32_t ts9_msec;
  uint32_t ts10_msec;
  uint32_t ts13_msec;
} sgs_config_t;

#define MME_CONFIG_MAX_SGW 16
typedef struct e_dns_config_s {
  int nb_sgw_entries;
  bstring sgw_id[MME_CONFIG_MAX_SGW];
  struct in_addr sgw_ip_addr[MME_CONFIG_MAX_SGW];
} e_dns_config_t;

typedef struct gummei_config_s {
  int nb;
  gummei_t gummei[MAX_GUMMEI];
} gummei_config_t;

typedef struct restricted_plmn_s {
  int num;
  plmn_t plmn[MAX_RESTRICTED_PLMN];
} restricted_plmn_config_t;

typedef struct blocked_imei_list_s {
  int num;
  // data is NULL
  hash_table_uint64_ts_t* imei_htbl;
} blocked_imei_list_t;

typedef struct fed_mode_map_s {
  uint8_t mode;
  plmn_t plmn;
  // IMSI range
  uint8_t imsi_low[MAX_IMSI_LENGTH + 1];
  uint8_t imsi_high[MAX_IMSI_LENGTH + 1];
  bstring apn;
} fed_mode_map_t;

typedef struct fed_mode_map_config_s {
  int num;
  fed_mode_map_t mode_map[MAX_FED_MODE_MAP_CONFIG];
} fed_mode_map_config_t;

typedef struct sac_to_tacs_map_config_s {
  tac_list_per_sac_t* tac_list;
  obj_hash_table_t* sac_to_tacs_map_htbl;
} sac_to_tacs_map_config_t;

typedef struct mme_config_s {
  /* Reader/writer lock for this configuration */
  pthread_rwlock_t rw_lock;

  bstring config_file;
  bstring pid_dir;
  bstring realm;
  bstring full_network_name;
  bstring short_network_name;
  uint8_t daylight_saving_time;

  run_mode_t run_mode;

  uint32_t max_enbs;
  uint32_t max_ues;

  uint8_t relative_capacity;

  uint32_t stats_timer_sec;

  bstring ip_capability;
  bstring non_eps_service_control;

  uint8_t unauthenticated_imsi_supported;

  eps_network_feature_config_t eps_network_feature_support;

  gummei_config_t gummei;

  restricted_plmn_config_t restricted_plmn;

  blocked_imei_list_t blocked_imei;
  sac_to_tacs_map_config_t sac_to_tacs_map;

  served_tai_t served_tai;
  uint8_t num_par_lists;
  partial_list_t* partial_list;

  service303_data_t service303_config;
  sctp_config_t sctp_config;
  s1ap_config_t s1ap_config;
  s6a_config_t s6a_config;
  itti_config_t itti_config;
  nas_config_t nas_config;
  sgs_config_t sgs_config;
  log_config_t log_config;
  e_dns_config_t e_dns_emulation;
  sentry_config_t sentry_config;

  ip_t ip;

  lai_t lai;
  fed_mode_map_config_t mode_map_config;
  bool use_stateless;
  bool use_ha;
  bool enable_gtpu_private_ip_correction;
  bool enable5g_features;
  bool accept_combined_attach_tau_wo_csfb;

  bool enable_congestion_control;
  long s1ap_zmq_th;
  long mme_app_zmq_congest_th;
  long mme_app_zmq_auth_th;
  long mme_app_zmq_ident_th;
  long mme_app_zmq_smc_th;
} mme_config_t;

extern mme_config_t mme_config;

int mme_config_find_mnc_length(
    const char mcc_digit1P, const char mcc_digit2P, const char mcc_digit3P,
    const char mnc_digit1P, const char mnc_digit2P, const char mnc_digit3P);

void mme_config_init(mme_config_t*);
int mme_config_parse_opt_line(int argc, char* argv[], mme_config_t* mme_config);
int mme_config_parse_file(mme_config_t*);
int mme_config_parse_string(const char* config_string, mme_config_t* config_pP);
void mme_config_display(mme_config_t*);
void create_partial_lists(mme_config_t* config_pP);
void mme_config_exit(void);

void free_mme_config(mme_config_t* mme_config);
void clear_served_tai_config(served_tai_t* served_tai);

void free_partial_lists(partial_list_t* partialList, uint8_t num_par_lists);

#define mme_config_read_lock(mMEcONFIG)                                        \
  pthread_rwlock_rdlock(&(mMEcONFIG)->rw_lock)
#define mme_config_write_lock(mMEcONFIG)                                       \
  pthread_rwlock_wrlock(&(mMEcONFIG)->rw_lock)
#define mme_config_unlock(mMEcONFIG)                                           \
  pthread_rwlock_unlock(&(mMEcONFIG)->rw_lock)

#endif /* FILE_MME_CONFIG_SEEN */
