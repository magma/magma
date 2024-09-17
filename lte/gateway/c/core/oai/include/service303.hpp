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

#pragma once

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/oai/include/service303_messages_types.hpp"

#define SERVICE303_MME_PACKAGE_NAME "mme"
#define SERVICE303_MME_PACKAGE_VERSION "1.0"
#define SERVICE303_SPGW_PACKAGE_NAME "spgw"
#define SERVICE303_SPGW_PACKAGE_VERSION "1.0"

#define NO_BOUNDARIES 0
#define NO_LABELS 0

void service303_mme_app_statistics_read(
    application_mme_app_stats_msg_t* stats_msg_p);
void service303_s1ap_statistics_read(application_s1ap_stats_msg_t* stats_msg_p);
void service303_statistics_display(void);
void service303_amf_app_statistics_read(
    application_amf_app_stats_msg_t* stats_msg_p);
void service303_statistics_display_5G(void);
void service303_ngap_statistics_read(application_ngap_stats_msg_t* stats_msg_p);

// service303 conf type added to be able to use same task interface for MME and
// SPGW while passing configs from mme_config and spgw_config types
typedef struct {
  bstring name;
  bstring version;
  uint32_t stats_display_timer_sec;
} service303_data_t;

typedef enum application_health_e {
  APP_UNKNOWN = 0,
  APP_UNHEALTHY = 1,
  APP_HEALTHY = 2,
} application_health_t;

#ifdef __cplusplus
extern "C" {
#endif
status_code_e service303_init(service303_data_t* service303_data);
#ifdef __cplusplus
}
#endif

/**
 * Start the Service303 Server and blocks
 *
 * @param name: service name string
 * @param version: the version number of the service
 */
void start_service303_server(bstring name, bstring version);

/**
 * Stop the server and clean up
 *
 */
void stop_service303_server(void);

/**
 * Simple helper function to set application health in the service. Only needed
 * to be called from a .c file.
 * @param health: one of 0 (APP_UNKNOWN), 1 (APP_HEALTHY), 2 (APP_UNHEALTHY)
 */
void service303_set_application_health(application_health_t health);
