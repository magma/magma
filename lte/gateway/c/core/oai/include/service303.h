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

#include "service303_messages_types.h"

#include "bstrlib.h"
#define SERVICE303_MME_PACKAGE_NAME "mme"
#define SERVICE303_MME_PACKAGE_VERSION "1.0"
#define SERVICE303_SPGW_PACKAGE_NAME "spgw"
#define SERVICE303_SPGW_PACKAGE_VERSION "1.0"

#define NO_BOUNDARIES 0
#define NO_LABELS 0
#define EPC_STATS_TIMER_MSEC 60000  // In milliseconds

void service303_mme_app_statistics_read(
    application_mme_app_stats_msg_t* stats_msg_p);
void service303_s1ap_statistics_read(application_s1ap_stats_msg_t* stats_msg_p);

// service303 conf type added to be able to use same task interface for MME and
// SPGW while passing configs from mme_config and spgw_config types
typedef struct {
  bstring name;
  bstring version;
} service303_data_t;

typedef enum application_health_e {
  APP_UNKNOWN   = 0,
  APP_UNHEALTHY = 1,
  APP_HEALTHY   = 2,
} application_health_t;

int service303_init(service303_data_t* service303_data);

#ifdef __cplusplus
extern "C" {
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
 * Increment a counter defined by the name and label set. Metric is
 * initialized if it doesn't yet exist. Usage example:
 *    increment_counter("test", 1, NO_LABELS)
 *    increment_counter("test", 1, 2, "key1", "val1", "key2", "val2")
 *
 * @param name: the counter family name
 * @param increment: the amount to increment the counter by
 * @param n_labels: the number of label pairs used or NO_LABELS
 * @param key1: the key of the first label
 * @param value1: the value of the first label
 */
void increment_counter(
    const char* name, double increment, size_t n_labels, ...);

/**
 * Increment a gauge defined by the name and label set. Metric is
 * initialized if it doesn't yet exist. Usage example:
 *    increment_gauge("test", 1, NO_LABELS)
 *    increment_gauge("test", 1, 2, "key1", "val1", "key2", "val2")
 *
 * @param name: the counter family name
 * @param increment: the amount to increment the gauge by
 * @param n_labels: the number of label pairs used or NO_LABELS
 * @param key1: the key of the first label
 * @param value1: the value of the first label
 */
void increment_gauge(const char* name, double increment, size_t n_labels, ...);

/**
 * Decrement a gauge defined by the name and label set. Metric is
 * initialized if it doesn't yet exist. Usage example:
 *    decrement_gauge("test", 1, NO_LABELS)
 *    decrement_gauge("test", 1, 2, "key1", "val1", "key2", "val2")
 *
 * @param name: the counter family name
 * @param decrement: the amount to decrement the gauge by
 * @param n_labels: the number of label pairs used or NO_LABELS
 * @param key1: the key of the first label
 * @param value1: the value of the first label
 */
void decrement_gauge(const char* name, double decrement, size_t n_labels, ...);

/**
 * Set a gauge defined by the name and label set. Metric is
 * initialized if it doesn't yet exist. Usage example:
 *    set_gauge("test", 1, NO_LABELS)
 *    set_gauge("test", 1, 2, "key1", "val1", "key2", "val2")
 *
 * @param name: the counter family name
 * @param value: the amount to set the gauge to
 * @param n_labels: the number of label pairs used or NO_LABELS
 * @param key1: the key of the first label
 * @param value1: the value of the first label
 */
void set_gauge(const char* name, double value, size_t n_labels, ...);

/**
 * Record an observation in the histogram defined by the name and label set.
 * Metric is initialized if it doesn't yet exist. The bucket boundaries are
 * static and only set on the first observation. Usage example:
 *    observe_histogram("test", 1, NO_LABELS, NO_BOUNDARIES);
 *    observe_histogram("test", 1, 1, "key", "value", NO_BOUNDARIES);
 *    observe_histogram("test", 1, 2, "key1", "value1", "key2", "value2",
 NO_BOUNDARIES);
 *    observe_histogram("test", 50, 1, "key", "value", 2, 10., 100.);

 *
 * @param name: the histogram family name
 * @param observation: the histogram obervation to record*
 * @param n_labels: the number of label pairs used or NO_LABELS
 * @param key1: the key of the first label
 * @param value1: the value of the first label
 * @param n_boundaries: the number of boundary definitions or NO_BOUNDARIES.
 *    NOTE: This must have type size_t
 * @param boundary1: floating point value of the first boundary
 *    NOTE: This must be a float or double (ie. requires a decimal point)
 */
void observe_histogram(
    const char* name, double observation, size_t n_labels, ...);

/**
 * Simple helper function to set application health in the service. Only needed
 * to be called from a .c file.
 * @param health: one of 0 (APP_UNKNOWN), 1 (APP_HEALTHY), 2 (APP_UNHEALTHY)
 */
void service303_set_application_health(application_health_t health);

#ifdef __cplusplus
}
#endif
