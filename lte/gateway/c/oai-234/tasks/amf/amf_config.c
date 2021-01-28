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

  Source      amf_config.c

  Version     0.1

  Date        2020/07/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#include "log.h"
#include "3gpp_24.501.h"
#include "amf_config.h"
#include "amf_default_values.h"
//----------------------------------------------------------------------------

void log_amf_config_init(log_config_t* log_conf) {
  // memset(log_conf, 0, sizeof(*log_conf));

  // log_conf->output                = NULL;
  // log_conf->is_output_thread_safe = false;
  // log_conf->color                 = false;

  log_conf->ngap_log_level    = MAX_LOG_LEVEL;
  log_conf->nas_amf_log_level = MAX_LOG_LEVEL;
  log_conf->amf_app_log_level = MAX_LOG_LEVEL;
}

void nas5g_config_init(nas5g_config_t* nas_conf) {
  nas_conf->t3502_min               = T3502_DEFAULT_VALUE;
  nas_conf->t3512_min               = T3512_DEFAULT_VALUE;
  nas_conf->t3522_sec               = T3522_DEFAULT_VALUE;
  nas_conf->t3550_sec               = T3550_DEFAULT_VALUE;
  nas_conf->t3560_sec               = T3560_DEFAULT_VALUE;
  nas_conf->t3570_sec               = T3570_DEFAULT_VALUE;
  nas_conf->t3585_sec               = T3585_DEFAULT_VALUE;
  nas_conf->t3586_sec               = T3586_DEFAULT_VALUE;
  nas_conf->t3589_sec               = T3589_DEFAULT_VALUE;
  nas_conf->t3595_sec               = T3595_DEFAULT_VALUE;
  nas_conf->force_reject_tau        = true;
  nas_conf->force_reject_sr         = true;
  nas_conf->disable_esm_information = false;
}
//-----------------------------------------------------------------------------
void guamfi_config_init(guamfi_config_t* guamfi_conf) {
  guamfi_conf->nb                        = 1;
  guamfi_conf->guamfi[0].amf_code        = AMFC;
  guamfi_conf->guamfi[0].amf_gid         = AMFGID;
  guamfi_conf->guamfi[0].amf_Pointer     = AMFPOINTER;
  guamfi_conf->guamfi[0].plmn.mcc_digit1 = 0;
  guamfi_conf->guamfi[0].plmn.mcc_digit2 = 0;
  guamfi_conf->guamfi[0].plmn.mcc_digit3 = 1;
  guamfi_conf->guamfi[0].plmn.mcc_digit1 = 0;
  guamfi_conf->guamfi[0].plmn.mcc_digit2 = 1;
  guamfi_conf->guamfi[0].plmn.mcc_digit3 = 0x0F;
}

void served_tai_config_init(m5g_served_tai_t* served_tai) {
  served_tai->nb_tai          = 1;
  served_tai->plmn_mcc        = calloc(1, sizeof(*served_tai->plmn_mcc));
  served_tai->plmn_mnc        = calloc(1, sizeof(*served_tai->plmn_mnc));
  served_tai->plmn_mnc_len    = calloc(1, sizeof(*served_tai->plmn_mnc_len));
  served_tai->tac             = calloc(1, sizeof(*served_tai->tac));
  served_tai->plmn_mcc[0]     = PLMN_MCC;
  served_tai->plmn_mnc[0]     = PLMN_MNC;
  served_tai->plmn_mnc_len[0] = PLMN_MNC_LEN;
  served_tai->tac[0]          = PLMN_TAC;
}

//-------------------------------------------------------------------------------
void ngap_config_init(ngap_config_t* ngap_conf) {
  ngap_conf->port_number            = NGAP_PORT_NUMBER;
  ngap_conf->outcome_drop_timer_sec = NGAP_OUTCOME_TIMER_DEFAULT;
}
//------------------------------------------------------------------------------
void amf_config_init(amf_config_t* config) {
  // memset(config, 0, sizeof(*config));

  pthread_rwlock_init(&config->rw_lock, NULL);

  // config->config_file                    = NULL;
  config->max_gnbs                       = 2;
  config->max_ues                        = 2;
  config->unauthenticated_imsi_supported = 0;
  config->relative_capacity              = RELATIVE_CAPACITY;
  config->amf_statistic_timer            = AMF_STATISTIC_TIMER_S;

  // log_amf_config_init(&config->log_config);
  ngap_config_init(&config->ngap_config);
  nas5g_config_init(&config->nas_config);
  guamfi_config_init(&config->guamfi);
  served_tai_config_init(&config->served_tai);
}

int amf_config_parse_opt_line(int argc, char* argv[], amf_config_t* config_pP) {
  // int c;

  amf_config_init(config_pP);

  return 0;
}
//---------------------------------------------------------------------------
