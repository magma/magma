/**
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

#include <pthread.h>
#include <stdint.h>
#include <arpa/inet.h>
#include <stdlib.h>
#include "lte/gateway/c/core/oai/common/amf_default_values.h"
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_23.003.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.008.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/include/service303.hpp"
#include "lte/gateway/c/core/oai/lib/hashtable/hashtable.h"
#include "lte/gateway/c/core/oai/include/mme_config.hpp"

#define MIN_GUAMI 1
#define MAX_GUAMI 5
#define MAX_APN_CORRECTION_MAP_LIST 10
#define AMF_S_NSSAI_ST_DEFAULT_VALUE 1
#define AMF_S_NSSAI_SD_INVALID_VALUE 0xffffff
#define AUTHENTICATION_COUNTER_MAX_RETRY "AUTHENTICATION_MAX_RETRY"
#define AUTHENTICATION_RETRY_TIMER_EXPIRY_MSECS "AUTHENTICATION_TIMER_EXPIRY"

#define AMF_CONFIG_STRING_AMF_CONFIG "AMF"
#define AMF_CONFIG_STRING_DEFAULT_DNS_IPV4_ADDRESS "DEFAULT_DNS_IPV4_ADDRESS"
#define AMF_CONFIG_STRING_DEFAULT_PCSCF_IPV4_ADDRESS "P_CSCF_IPV4_ADDRESS"
#define AMF_CONFIG_STRING_DEFAULT_PCSCF_IPV6_ADDRESS "P_CSCF_IPV6_ADDRESS"
#define AMF_CONFIG_STRING_DEFAULT_DNS_SEC_IPV4_ADDRESS \
  "DEFAULT_DNS_SEC_IPV4_ADDRESS"
#define AMF_CONFIG_PLMN_SUPPORT_MCC "mcc"
#define AMF_CONFIG_PLMN_SUPPORT_MNC "mnc"
#define AMF_CONFIG_PLMN_SUPPORT_SST "AMF_DEFAULT_SLICE_SERVICE_TYPE"
#define AMF_CONFIG_PLMN_SUPPORT_SD "AMF_DEFAULT_SLICE_DIFFERENTIATOR"
#define CONFIG_DEFAULT_DNN "DEFAULT_DNN"
#define AMF_CONFIG_AMF_PLMN_SUPPORT_LIST "PLMN_SUPPORT_LIST"
#define AMF_CONFIG_AMF_NAME "AMF_NAME"

#define AMF_CONFIG_STRING_GUAMFI_LIST "GUAMFI_LIST"
#define AMF_CONFIG_STRING_AMF_REGION_ID "AMF_REGION_ID"
#define AMF_CONFIG_STRING_AMF_SET_ID "AMF_SET_ID"
#define AMF_CONFIG_STRING_AMF_POINTER "AMF_POINTER"
#define AMF_CONFIG_STRING_NAS_ENABLE_IMS_VoPS_3GPP "ENABLE_IMS_VoPS_3GPP"
#define AMF_CONFIG_STRING_NAS_T3512 "T3512"

typedef struct nas5g_config_s {
  uint8_t preferred_integrity_algorithm[8];
  uint8_t preferred_ciphering_algorithm[8];
  uint32_t t3502_min;
  uint32_t t3512_min;
  uint32_t t3522_sec;
  uint32_t t3550_sec;
  uint32_t t3560_sec;
  uint32_t t3570_sec;
  uint32_t t3585_sec;
  uint32_t t3586_sec;
  uint32_t t3589_sec;
  uint32_t t3595_sec;
  uint32_t implicit_dereg_sec;

  // non standard features
  bool force_reject_tau;
  bool force_reject_sr;
  bool disable_esm_information;
  bool enable_IMS_VoPS_3GPP;
} nas5g_config_t;

typedef struct m5g_apn_map_s {
  bstring imsi_prefix;
  bstring apn_override;
} m5g_apn_map_t;

typedef struct m5g_apn_map_config_s {
  int nb;
  m5g_apn_map_t apn_map[MAX_APN_CORRECTION_MAP_LIST];
} m5g_apn_map_config_t;

typedef struct m5g_nas_config_s {
  uint8_t preferred_integrity_algorithm[8];
  uint8_t preferred_ciphering_algorithm[8];
  uint32_t t3402_min;
  uint32_t t3412_min;
  uint32_t t3422_msec;
  uint32_t t3450_msec;
  uint32_t t3460_msec;
  uint32_t t3470_msec;
  uint32_t t3485_msec;
  uint32_t t3486_msec;
  uint32_t t3489_msec;
  uint32_t t3495_msec;
  // non standard features
  bool force_reject_tau;
  bool force_reject_sr;
  bool disable_esm_information;
  // apn correction
  bool enable_apn_correction;
  m5g_apn_map_config_t m5g_apn_map_config;
} m5g_nas_config_t;

typedef uint64_t imsi64_t;

typedef struct ngap_config_s {
  uint16_t port_number;
  uint8_t outcome_drop_timer_sec;
} ngap_config_t;

typedef struct guamfi_config_s {
  int nb;
  guamfi_t guamfi[MAX_GUAMI];
#define MIN_GUAMFI 1 /*minimum 1 Global Unique AMF Identifier is supported*/
#define MAX_GUAMFI 5 /*max 5 Global Unique AMF Identifiers are supported*/

#define amf_config_read_lock(aMFcONFIG) \
  pthread_rwlock_rdlock(&(aMFcONFIG)->rw_lock)
#define amf_config_write_lock(aMFcONFIG) \
  pthread_rwlock_wrlock(&(aMFcONFIG)->rw_lock)
#define amf_config_unlock(aMFcONFIG) \
  pthread_rwlock_unlock(&(aMFcONFIG)->rw_lock)

  uint64_t imsi64_t;         /*holds the IMSI value*/
  uint64_t amf_ue_ngap_id_t; /*uniquely identifies the UE over the NG
                                        interface within the AMF*/
} guamfi_config_t;

typedef struct amf_uint24_s {
  uint32_t v : 24;
} __attribute__((packed)) amf_uint24_t;

typedef struct amf_s_nssai_s {
  uint8_t sst;
  amf_uint24_t sd;
} __attribute__((packed)) amf_s_nssai_t;

typedef struct plmn_support_s {
  plmn_t plmn;
  amf_s_nssai_t s_nssai;
} plmn_support_t;

typedef struct plmn_support_list_s {
#define MIN_PLMN_SUPPORT 1
#define MAX_PLMN_SUPPORT 5
  uint8_t plmn_support_count;
  plmn_support_t plmn_support[MAX_PLMN_SUPPORT];
} plmn_support_list_t;

typedef struct amf_config_s {
  /* Reader/writer lock for this configuration */
  pthread_rwlock_t rw_lock;

  bstring config_file;
  bstring pid_dir;
  bstring realm;
  bstring full_network_name;
  bstring short_network_name;
  uint8_t daylight_saving_time;
  uint32_t max_gnbs;
  uint32_t max_ues;
  uint8_t relative_capacity;
  bstring ip_capability;
  uint8_t unauthenticated_imsi_supported;
  guamfi_config_t guamfi;
  plmn_support_list_t plmn_support_list;
  served_tai_t served_tai;
  uint8_t num_par_lists;
  partial_list_t* partial_list;
  service303_data_t service303_config;
  ngap_config_t ngap_config;
  m5g_nas_config_t m5g_nas_config;
  log_config_t log_config;
  uint32_t amf_statistic_timer;
  nas5g_config_t nas_config;
  bool use_stateless;
  struct {
    struct in_addr default_dns;
    struct in_addr default_dns_sec;
  } ipv4;
  struct {
    struct in_addr ipv4;
  } pcscf_addr;

  bstring amf_name;
  bstring default_dnn;
  uint32_t auth_retry_interval;
  uint32_t auth_retry_max_count;
} amf_config_t;

int amf_app_init(amf_config_t*);

extern amf_config_t amf_config; /*global*/

int amf_config_find_mnc_length(const char mcc_digit1P, const char mcc_digit2P,
                               const char mcc_digit3P, const char mnc_digit1P,
                               const char mnc_digit2P, const char mnc_digit3P);

void amf_config_init(amf_config_t*);
int amf_config_parse_opt_line(int argc, char* argv[], amf_config_t* amf_config);
int amf_config_parse_file(amf_config_t*, const mme_config_t*);
void amf_config_display(amf_config_t*);
#ifdef __cplusplus
extern "C" {
#endif
void copy_amf_config_from_mme_config(amf_config_t* dest,
                                     const mme_config_t* src);
void clear_amf_config(amf_config_t*);
#ifdef __cplusplus
}
#endif
void copy_served_tai_config_list(amf_config_t* dest, const mme_config_t* src);

void amf_config_exit(void);
void amf_config_free(amf_config_t* amf_config);
