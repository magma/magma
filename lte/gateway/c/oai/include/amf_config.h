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
/*****************************************************************************
  Source      amf_message.h
  Date        2020/07/28
  Author      Ashish Prajapati
  Subsystem   Access and Mobility Management Function
  Description Defines Access and Mobility Management Messages
*****************************************************************************/
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

#define MIN_GUAMFI 1 /*minimum 1 Global Unique AMF Identifier is supported*/
#define MAX_GUAMFI 5 /*max 5 Global Unique AMF Identifiers are supported*/

typedef uint64_t imsi64_t;         /*holds the IMSI value*/
typedef uint32_t amf_ue_ngap_id_t; /*uniquely identifies the UE over the NG
                                      interface within the AMF*/

/*TAI list*/
typedef struct m5g_served_tai_s {
  uint8_t list_type;
  uint8_t nb_tai;
  uint16_t* plmn_mcc; /*Mobile Country Code*/
  uint16_t* plmn_mnc; /*Mobile Network Code*/
  uint16_t* plmn_mnc_len;
  uint16_t* tac; /*Tracking Area Code*/
} m5g_served_tai_t;

typedef struct ngap_config_s {
  uint16_t port_number; /*port #38412 for NGAP*/
} ngap_config_t;

/*Global Unique AMF Identifier*/
typedef struct guamfi_config_s {
  int nb;
  guamfi_t guamfi[MAX_GUAMFI];
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
  log_config_t log_config;
  bool use_stateless;
} amf_config_t;

extern amf_config_t amf_config; /*global*/

int amf_config_find_mnc_length(
    const char mcc_digit1P, const char mcc_digit2P, const char mcc_digit3P,
    const char mnc_digit1P, const char mnc_digit2P, const char mnc_digit3P);

void amf_config_init(amf_config_t*);
int amf_config_parse_opt_line(int argc, char* argv[], amf_config_t* amf_config);
int amf_config_parse_file(amf_config_t*);
void amf_config_display(amf_config_t*);

void amf_config_exit(void);
#define amf_config_read_lock(aMFcONFIG)                                        \
  pthread_rwlock_rdlock(&(aMFcONFIG)->rw_lock)
#define amf_config_write_lock(aMFcONFIG)                                       \
  pthread_rwlock_wrlock(&(aMFcONFIG)->rw_lock)
#define amf_config_unlock(aMFcONFIG)                                           \
  pthread_rwlock_unlock(&(aMFcONFIG)->rw_lock)
