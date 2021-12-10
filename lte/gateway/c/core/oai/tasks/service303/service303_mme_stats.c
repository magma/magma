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
#define SERVICE303

#include <stddef.h>
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/include/service303.h"
#include "orc8r/gateway/c/common/service303/includes/MetricsHelpers.h"

void service303_mme_app_statistics_read(
    application_mme_app_stats_msg_t* stats_msg_p) {
  size_t label = 0;
  set_gauge("ue_registered", stats_msg_p->nb_ue_attached, label);
  set_gauge("ue_connected", stats_msg_p->nb_ue_connected, label);
  set_gauge("default_eps_bearers", stats_msg_p->nb_default_eps_bearers, label);
  set_gauge("s1u_bearers", stats_msg_p->nb_s1u_bearers, label);
  set_gauge(
      "mme_app_last_msg_latency", stats_msg_p->nb_mme_app_last_msg_latency,
      label);
}

void service303_s1ap_statistics_read(
    application_s1ap_stats_msg_t* stats_msg_p) {
  size_t label = 0;
  set_gauge("enb_connected", stats_msg_p->nb_enb_connected, label);
  set_gauge(
      "s1ap_last_msg_latency", stats_msg_p->nb_s1ap_last_msg_latency, label);
}

void service303_statistics_display(void) {
  size_t label = 0;
  OAILOG_DEBUG(
      LOG_SERVICE303,
      "======================================= STATISTICS "
      "============================================\n\n");
  OAILOG_DEBUG(LOG_SERVICE303, "               |   Current Status|\n");
  OAILOG_DEBUG(
      LOG_SERVICE303, "Attached UEs   | %10u      |\n",
      (uint32_t) get_gauge("ue_registered", label));
  OAILOG_DEBUG(
      LOG_SERVICE303, "Connected UEs  | %10u      |\n",
      (uint32_t) get_gauge("ue_connected", label));
  OAILOG_DEBUG(
      LOG_SERVICE303, "Connected eNBs | %10u      |\n",
      (uint32_t) get_gauge("enb_connected", label));
  OAILOG_DEBUG(
      LOG_SERVICE303, "Default Bearers| %10u      |\n",
      (uint32_t) get_gauge("default_eps_bearers", label));
  OAILOG_DEBUG(
      LOG_SERVICE303, "S1-U Bearers   | %10u      |\n\n",
      (uint32_t) get_gauge("s1u_bearers", label));

  OAILOG_DEBUG(
      LOG_SERVICE303,
      "======================================= STATISTICS "
      "============================================\n\n");
}
