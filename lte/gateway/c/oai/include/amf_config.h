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
#include "amf_default_values.h"
#include "common_types.h"
#include "3gpp_23.003.h"
#include "3gpp_24.008.h"
#include "log.h"
#include "service303.h"

#define MIN_GUMMEI 1
#define MAX_GUMMEI 5
#define MAX_APN_CORRECTION_MAP_LIST 10

extern char buf_plmn[3];

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
  // non standard features
  bool force_reject_tau;
  bool force_reject_sr;
  bool disable_esm_information;
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
  uint32_t t3422_sec;
  uint32_t t3450_sec;
  uint32_t t3460_sec;
  uint32_t t3470_sec;
  uint32_t t3485_sec;
  uint32_t t3486_sec;
  uint32_t t3489_sec;
  uint32_t t3495_sec;
  // non standard features
  bool force_reject_tau;
  bool force_reject_sr;
  bool disable_esm_information;
  // apn correction
  bool enable_apn_correction;
  m5g_apn_map_config_t m5g_apn_map_config;
} m5g_nas_config_t;
typedef uint64_t imsi64_t;
typedef uint32_t amf_ue_ngap_id_t;

typedef struct m5g_served_tai_s {
  uint8_t list_type;
  uint8_t nb_tai;
  uint16_t* plmn_mcc;
  uint16_t* plmn_mnc;
  uint16_t* plmn_mnc_len;
  uint16_t* tac;
} m5g_served_tai_t;

typedef struct ngap_config_s {
  uint16_t port_number;
  uint8_t outcome_drop_timer_sec;
} ngap_config_t;

typedef struct guamfi_config_s {
  int nb;
  guamfi_t guamfi[MAX_GUMMEI];
#define MIN_GUAMFI 1 /*minimum 1 Global Unique AMF Identifier is supported*/
#define MAX_GUAMFI 5 /*max 5 Global Unique AMF Identifiers are supported*/

#define amf_config_read_lock(aMFcONFIG)                                        \
  pthread_rwlock_rdlock(&(aMFcONFIG)->rw_lock)
#define amf_config_write_lock(aMFcONFIG)                                       \
  pthread_rwlock_wrlock(&(aMFcONFIG)->rw_lock)
#define amf_config_unlock(aMFcONFIG)                                           \
  pthread_rwlock_unlock(&(aMFcONFIG)->rw_lock)

  uint64_t imsi64_t;         /*holds the IMSI value*/
  uint32_t amf_ue_ngap_id_t; /*uniquely identifies the UE over the NG
                                        interface within the AMF*/
} guamfi_config_t;

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
  m5g_served_tai_t served_tai;
  service303_data_t service303_config;
  ngap_config_t ngap_config;
  m5g_nas_config_t m5g_nas_config;
  log_config_t log_config;
  uint32_t amf_statistic_timer;
  nas5g_config_t nas_config;
  bool use_stateless;
} amf_config_t;

int amf_app_init(amf_config_t*);

extern amf_config_t amf_config; /*global*/

int amf_config_find_mnc_length(
    const char mcc_digit1P, const char mcc_digit2P, const char mcc_digit3P,
    const char mnc_digit1P, const char mnc_digit2P, const char mnc_digit3P);

void amf_config_init(amf_config_t*);
int amf_config_parse_opt_line(int argc, char* argv[], amf_config_t* amf_config);
int amf_config_parse_file(amf_config_t*);
void amf_config_display(amf_config_t*);

void amf_config_exit(void);
