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
#include <libconfig.h>
#include "lte/gateway/c/core/oai/common/log.h"
#include <errno.h>
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.501.h"
#include "lte/gateway/c/core/oai/include/amf_config.h"
#include "lte/gateway/c/core/oai/common/amf_default_values.h"
#include "lte/gateway/c/core/oai/include/mme_config.h"
#include "lte/gateway/c/core/oai/include/TrackingAreaIdentity.h"
#include "lte/gateway/c/core/oai/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/common/assertions.h"

static bool parse_bool(const char* str);

struct amf_config_s amf_config = {.rw_lock = PTHREAD_RWLOCK_INITIALIZER, 0};

/***************************************************************************
**                                                                        **
** Name:    log_amf_config_init()                                         **
**                                                                        **
** Description: Initializes log level of AMF                              **
**                                                                        **
**                                                                        **
***************************************************************************/
void log_amf_config_init(log_config_t* log_conf) {
  log_conf->ngap_log_level    = MAX_LOG_LEVEL;
  log_conf->nas_amf_log_level = MAX_LOG_LEVEL;
  log_conf->amf_app_log_level = MAX_LOG_LEVEL;
}

/***************************************************************************
**                                                                        **
** Name:    nas5g_config_init()                                           **
**                                                                        **
** Description: Initializes default values for NAS5G                      **
**                                                                        **
**                                                                        **
***************************************************************************/
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

/***************************************************************************
**                                                                        **
** Name:    guamfi_config_init()                                          **
**                                                                        **
** Description: Initializes default values for guamfi                     **
**                                                                        **
**                                                                        **
***************************************************************************/
void guamfi_config_init(guamfi_config_t* guamfi_conf) {
  guamfi_conf->nb                        = 1;
  guamfi_conf->guamfi[0].amf_set_id      = AMFC;
  guamfi_conf->guamfi[0].amf_regionid    = AMFGID;
  guamfi_conf->guamfi[0].amf_pointer     = AMFPOINTER;
  guamfi_conf->guamfi[0].plmn.mcc_digit1 = 0;
  guamfi_conf->guamfi[0].plmn.mcc_digit2 = 0;
  guamfi_conf->guamfi[0].plmn.mcc_digit3 = 1;
  guamfi_conf->guamfi[0].plmn.mcc_digit1 = 0;
  guamfi_conf->guamfi[0].plmn.mcc_digit2 = 1;
  guamfi_conf->guamfi[0].plmn.mcc_digit3 = 0x0F;
}

/***************************************************************************
**                                                                        **
** Name:    plmn_support_list_config_init()                               **
**                                                                        **
** Description: Initializes default values for plmn_support_list          **
**                                                                        **
**                                                                        **
***************************************************************************/
void plmn_support_list_config_init(plmn_support_list_t* plmn_support_list) {
  plmn_support_list->plmn_support_count              = MIN_PLMN_SUPPORT;
  plmn_support_list->plmn_support[0].plmn.mcc_digit1 = 0;
  plmn_support_list->plmn_support[0].plmn.mcc_digit2 = 0;
  plmn_support_list->plmn_support[0].plmn.mcc_digit3 = 0;
  plmn_support_list->plmn_support[0].plmn.mcc_digit1 = 0;
  plmn_support_list->plmn_support[0].plmn.mcc_digit2 = 0;
  plmn_support_list->plmn_support[0].plmn.mcc_digit3 = 0x0F;
  plmn_support_list->plmn_support[0].s_nssai.sst =
      NGAP_S_NSSAI_ST_DEFAULT_VALUE;
  plmn_support_list->plmn_support[0].s_nssai.sd.v =
      NGAP_S_NSSAI_SD_INVALID_VALUE;
}

/***************************************************************************
**                                                                        **
** Name:   m5g_served_tai_config_init()                                   **
**                                                                        **
** Description: Initializes default values for served_tai                 **
**                                                                        **
**                                                                        **
***************************************************************************/
void m5g_served_tai_config_init(m5g_served_tai_t* served_tai) {
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

/***************************************************************************
**                                                                        **
** Name:    ngap_config_init()                                            **
**                                                                        **
** Description: Initializes default values for NGAP                       **
**                                                                        **
**                                                                        **
***************************************************************************/
void ngap_config_init(ngap_config_t* ngap_conf) {
  ngap_conf->port_number            = NGAP_PORT_NUMBER;
  ngap_conf->outcome_drop_timer_sec = NGAP_OUTCOME_TIMER_DEFAULT;
}

/***************************************************************************
**                                                                        **
** Name:    amf_config_init()                                             **
**                                                                        **
** Description: Initializes default values for AMF                        **
**                                                                        **
**                                                                        **
***************************************************************************/
void amf_config_init(amf_config_t* config) {
  memset(config, 0, sizeof(*config));

  pthread_rwlock_init(&config->rw_lock, NULL);

  config->max_gnbs                       = 2;
  config->max_ues                        = 2;
  config->unauthenticated_imsi_supported = 0;
  config->relative_capacity              = RELATIVE_CAPACITY;
  config->amf_statistic_timer            = AMF_STATISTIC_TIMER_S;
  config->use_stateless                  = false;
  ngap_config_init(&config->ngap_config);
  nas5g_config_init(&config->nas_config);
  guamfi_config_init(&config->guamfi);
  plmn_support_list_config_init(&config->plmn_support_list);
  m5g_served_tai_config_init(&config->served_tai);
}

/***************************************************************************
**                                                                        **
** Name:    amf_config_parse_opt_line()                                   **
**                                                                        **
** Description: Invokes amf_config_init() to initialize                   **
**              default values of AMF                                     **
**                                                                        **
**                                                                        **
***************************************************************************/
int amf_config_parse_opt_line(int argc, char* argv[], amf_config_t* config_pP) {
  amf_config_init(config_pP);
  return 0;
}

/***************************************************************************
**                                                                        **
** Name:    amf_config_parse_opt_line()                                   **
**                                                                        **
** Description: Invokes amf_config_init() to initialize                   **
**              default values of AMF                                     **
**                                                                        **
**                                                                        **
***************************************************************************/
int amf_config_parse_file(amf_config_t* config_pP) {
  config_t cfg                   = {0};
  config_setting_t* setting_mme  = NULL;
  config_setting_t* setting      = NULL;
  config_setting_t* sub2setting  = NULL;
  config_setting_t* setting_ngap = NULL;
  int aint                       = 0;
  int i = 0, n = 0, stop_index = 0, num = 0;
  const char* astring         = NULL;
  const char* tac             = NULL;
  const char* mcc             = NULL;
  const char* mnc             = NULL;
  bool swap                   = false;
  const char* region_id       = NULL;
  const char* set_id          = NULL;
  const char* pointer         = NULL;
  const char* default_dns     = NULL;
  const char* default_dns_sec = NULL;
  const char* set_sst         = NULL;
  const char* set_sd          = NULL;

  config_init(&cfg);

  if (config_pP->config_file != NULL) {
    /*
     * Read the file. If there is an error, report it and exit.
     */
    if (!config_read_file(&cfg, bdata(config_pP->config_file))) {
      OAILOG_CRITICAL(
          LOG_CONFIG, "Failed to parse AMF configuration file: %s:%d - %s\n",
          bdata(config_pP->config_file), config_error_line(&cfg),
          config_error_text(&cfg));
      config_destroy(&cfg);
      AssertFatal(
          1 == 0, "Failed to parse AMF configuration file %s!\n",
          bdata(config_pP->config_file));
    }
  } else {
    config_destroy(&cfg);
    AssertFatal(0, "No AMF configuration file provided!\n");
  }

  setting_mme = config_lookup(&cfg, MME_CONFIG_STRING_MME_CONFIG);

  if (setting_mme != NULL) {
    // LOGGING setting
    setting = config_setting_get_member(setting_mme, LOG_CONFIG_STRING_LOGGING);

    if (setting != NULL) {
      if (config_setting_lookup_string(
              setting, LOG_CONFIG_STRING_OUTPUT, (const char**) &astring)) {
        if (astring != NULL) {
          if (config_pP->log_config.output) {
            bassigncstr(config_pP->log_config.output, astring);
          } else {
            config_pP->log_config.output = bfromcstr(astring);
          }
        }
      }

      if (config_setting_lookup_string(
              setting, LOG_CONFIG_STRING_OUTPUT_THREAD_SAFE,
              (const char**) &astring)) {
        if (astring != NULL) {
          config_pP->log_config.is_output_thread_safe = parse_bool(astring);
        }
      }

      if (config_setting_lookup_string(
              setting, LOG_CONFIG_STRING_COLOR, (const char**) &astring)) {
        if (strcasecmp("yes", astring) == 0)
          config_pP->log_config.color = true;
        else
          config_pP->log_config.color = false;
      }

      if (config_setting_lookup_string(
              setting, LOG_CONFIG_STRING_SCTP_LOG_LEVEL,
              (const char**) &astring))
        config_pP->log_config.sctp_log_level = OAILOG_LEVEL_STR2INT(astring);

      if (config_setting_lookup_string(
              setting, LOG_CONFIG_STRING_S1AP_LOG_LEVEL,
              (const char**) &astring))
        config_pP->log_config.s1ap_log_level = OAILOG_LEVEL_STR2INT(astring);

      if (config_setting_lookup_string(
              setting, LOG_CONFIG_STRING_NAS_LOG_LEVEL,
              (const char**) &astring))
        config_pP->log_config.nas_log_level = OAILOG_LEVEL_STR2INT(astring);

      if (config_setting_lookup_string(
              setting, LOG_CONFIG_STRING_MME_APP_LOG_LEVEL,
              (const char**) &astring))
        config_pP->log_config.amf_app_log_level = OAILOG_LEVEL_STR2INT(astring);

      if (config_setting_lookup_string(
              setting, LOG_CONFIG_STRING_SECU_LOG_LEVEL,
              (const char**) &astring))
        config_pP->log_config.secu_log_level = OAILOG_LEVEL_STR2INT(astring);

      if (config_setting_lookup_string(
              setting, LOG_CONFIG_STRING_UDP_LOG_LEVEL,
              (const char**) &astring))
        config_pP->log_config.udp_log_level = OAILOG_LEVEL_STR2INT(astring);

      if (config_setting_lookup_string(
              setting, LOG_CONFIG_STRING_UTIL_LOG_LEVEL,
              (const char**) &astring))
        config_pP->log_config.util_log_level = OAILOG_LEVEL_STR2INT(astring);

      if (config_setting_lookup_string(
              setting, LOG_CONFIG_STRING_ITTI_LOG_LEVEL,
              (const char**) &astring))
        config_pP->log_config.itti_log_level = OAILOG_LEVEL_STR2INT(astring);
      if (config_setting_lookup_string(
              setting, LOG_CONFIG_STRING_GTPV1U_LOG_LEVEL,
              (const char**) &astring))
        config_pP->log_config.gtpv1u_log_level = OAILOG_LEVEL_STR2INT(astring);

      if ((config_setting_lookup_string(
              setting_mme, MME_CONFIG_STRING_ASN1_VERBOSITY,
              (const char**) &astring))) {
        if (strcasecmp(astring, MME_CONFIG_STRING_ASN1_VERBOSITY_NONE) == 0)
          config_pP->log_config.asn1_verbosity_level = 0;
        else if (
            strcasecmp(astring, MME_CONFIG_STRING_ASN1_VERBOSITY_ANNOYING) == 0)
          config_pP->log_config.asn1_verbosity_level = 2;
        else if (
            strcasecmp(astring, MME_CONFIG_STRING_ASN1_VERBOSITY_INFO) == 0)
          config_pP->log_config.asn1_verbosity_level = 1;
        else
          config_pP->log_config.asn1_verbosity_level = 0;
      }
    }

    // GENERAL AMF SETTINGS
    if ((config_setting_lookup_string(
            setting_mme, MME_CONFIG_STRING_REALM, (const char**) &astring))) {
      config_pP->realm = bfromcstr(astring);
    }

    if ((config_setting_lookup_string(
            setting_mme, MME_CONFIG_STRING_FULL_NETWORK_NAME,
            (const char**) &astring))) {
      config_pP->full_network_name = bfromcstr(astring);
    }

    if ((config_setting_lookup_string(
            setting_mme, MME_CONFIG_STRING_SHORT_NETWORK_NAME,
            (const char**) &astring))) {
      config_pP->short_network_name = bfromcstr(astring);
    }

    if ((config_setting_lookup_int(
            setting_mme, MME_CONFIG_STRING_DAYLIGHT_SAVING_TIME, &aint))) {
      config_pP->daylight_saving_time = (uint32_t) aint;
    }

    if ((config_setting_lookup_string(
            setting_mme, MME_CONFIG_STRING_PID_DIRECTORY,
            (const char**) &astring))) {
      config_pP->pid_dir = bfromcstr(astring);
    }

    if ((config_setting_lookup_int(
            setting_mme, MME_CONFIG_STRING_MAXENB, &aint))) {
      config_pP->max_gnbs = (uint32_t) aint;
    }

    if ((config_setting_lookup_int(
            setting_mme, MME_CONFIG_STRING_MAXUE, &aint))) {
      config_pP->max_ues = (uint32_t) aint;
    }

    if ((config_setting_lookup_int(
            setting_mme, MME_CONFIG_STRING_RELATIVE_CAPACITY, &aint))) {
      config_pP->relative_capacity = (uint8_t) aint;
    }

    if ((config_setting_lookup_string(
            setting_mme, MME_CONFIG_STRING_USE_STATELESS,
            (const char**) &astring))) {
      config_pP->use_stateless = parse_bool(astring);
    }

    if ((config_setting_lookup_string(
            setting_mme, MME_CONFIG_STRING_UNAUTHENTICATED_IMSI_SUPPORTED,
            (const char**) &astring))) {
      config_pP->unauthenticated_imsi_supported = parse_bool(astring);
    }

    // guamfi SETTING
    setting =
        config_setting_get_member(setting_mme, MME_CONFIG_STRING_GUAMFI_LIST);
    config_pP->guamfi.nb = 0;
    if (setting != NULL) {
      num = config_setting_length(setting);
      AssertFatal(
          num >= MIN_GUMMEI,
          "Not even one guamfi is configured, configure minimum one guamfi \n");
      AssertFatal(
          num <= MAX_GUMMEI,
          "Number of guamfis configured:%d exceeds number of guamfis supported "
          ":%d \n",
          num, MAX_GUMMEI);

      for (i = 0; i < num; i++) {
        sub2setting = config_setting_get_elem(setting, i);
        if (sub2setting != NULL) {
          if ((config_setting_lookup_string(
                  sub2setting, MME_CONFIG_STRING_MCC, &mcc))) {
            AssertFatal(
                strlen(mcc) == MAX_MCC_LENGTH,
                "Bad MCC length (%ld), it must be %u digit ex: 001",
                strlen(mcc), MAX_MCC_LENGTH);
            char c[2]                                   = {mcc[0], 0};
            config_pP->guamfi.guamfi[i].plmn.mcc_digit1 = (uint8_t) atoi(c);
            c[0]                                        = mcc[1];
            config_pP->guamfi.guamfi[i].plmn.mcc_digit2 = (uint8_t) atoi(c);
            c[0]                                        = mcc[2];
            config_pP->guamfi.guamfi[i].plmn.mcc_digit3 = (uint8_t) atoi(c);
          }

          if ((config_setting_lookup_string(
                  sub2setting, MME_CONFIG_STRING_MNC, &mnc))) {
            AssertFatal(
                (strlen(mnc) == MIN_MNC_LENGTH) ||
                    (strlen(mnc) == MAX_MNC_LENGTH),
                "Bad MNC length (%ld), it must be %u or %u digit ex: 12 or 123",
                strlen(mnc), MIN_MNC_LENGTH, MAX_MNC_LENGTH);
            char c[2]                                   = {mnc[0], 0};
            config_pP->guamfi.guamfi[i].plmn.mnc_digit1 = (uint8_t) atoi(c);
            c[0]                                        = mnc[1];
            config_pP->guamfi.guamfi[i].plmn.mnc_digit2 = (uint8_t) atoi(c);
            if (3 == strlen(mnc)) {
              c[0]                                        = mnc[2];
              config_pP->guamfi.guamfi[i].plmn.mnc_digit3 = (uint8_t) atoi(c);
            } else {
              config_pP->guamfi.guamfi[i].plmn.mnc_digit3 = 0x0F;
            }
          }

          if ((config_setting_lookup_string(
                  sub2setting, MME_CONFIG_STRING_AMF_REGION_ID, &region_id))) {
            config_pP->guamfi.guamfi[i].amf_regionid =
                (uint16_t) atoi(region_id);
          }
          if ((config_setting_lookup_string(
                  sub2setting, MME_CONFIG_STRING_AMF_SET_ID, &set_id))) {
            config_pP->guamfi.guamfi[i].amf_set_id = (uint8_t) atoi(set_id);
          }
          if ((config_setting_lookup_string(
                  sub2setting, MME_CONFIG_STRING_AMF_POINTER, &pointer))) {
            config_pP->guamfi.guamfi[i].amf_pointer = (uint8_t) atoi(pointer);
          }

          config_pP->guamfi.nb += 1;
        }
      }
    }

    // TAI list setting
    setting =
        config_setting_get_member(setting_mme, MME_CONFIG_STRING_TAI_LIST);
    if (setting != NULL) {
      num = config_setting_length(setting);
      if (num < MIN_TAI_SUPPORTED) {
        fprintf(
            stderr,
            "ERROR: No TAI is configured.  At least one TAI must be "
            "configured.\n");
      }

      if (config_pP->served_tai.nb_tai != num) {
        if (config_pP->served_tai.plmn_mcc != NULL)
          free_wrapper((void**) &config_pP->served_tai.plmn_mcc);

        if (config_pP->served_tai.plmn_mnc != NULL)
          free_wrapper((void**) &config_pP->served_tai.plmn_mnc);

        if (config_pP->served_tai.plmn_mnc_len != NULL)
          free_wrapper((void**) &config_pP->served_tai.plmn_mnc_len);

        if (config_pP->served_tai.tac != NULL)
          free_wrapper((void**) &config_pP->served_tai.tac);

        config_pP->served_tai.plmn_mcc =
            calloc(num, sizeof(*config_pP->served_tai.plmn_mcc));
        config_pP->served_tai.plmn_mnc =
            calloc(num, sizeof(*config_pP->served_tai.plmn_mnc));
        config_pP->served_tai.plmn_mnc_len =
            calloc(num, sizeof(*config_pP->served_tai.plmn_mnc_len));
        config_pP->served_tai.tac =
            calloc(num, sizeof(*config_pP->served_tai.tac));
      }

      config_pP->served_tai.nb_tai = num;

      for (i = 0; i < num; i++) {
        sub2setting = config_setting_get_elem(setting, i);

        if (sub2setting != NULL) {
          if ((config_setting_lookup_string(
                  sub2setting, MME_CONFIG_STRING_MCC, &mcc))) {
            config_pP->served_tai.plmn_mcc[i] = (uint16_t) atoi(mcc);
          }

          if ((config_setting_lookup_string(
                  sub2setting, MME_CONFIG_STRING_MNC, &mnc))) {
            config_pP->served_tai.plmn_mnc[i]     = (uint16_t) atoi(mnc);
            config_pP->served_tai.plmn_mnc_len[i] = strlen(mnc);

            AssertFatal(
                (config_pP->served_tai.plmn_mnc_len[i] == MIN_MNC_LENGTH) ||
                    (config_pP->served_tai.plmn_mnc_len[i] == MAX_MNC_LENGTH),
                "Bad MNC length %u, must be %d or %d",
                config_pP->served_tai.plmn_mnc_len[i], MIN_MNC_LENGTH,
                MAX_MNC_LENGTH);
          }

          if ((config_setting_lookup_string(
                  sub2setting, MME_CONFIG_STRING_TAC, &tac))) {
            config_pP->served_tai.tac[i] = (uint16_t) atoi(tac);

            if (!TAC_IS_VALID(config_pP->served_tai.tac[i])) {
              fprintf(
                  stderr, "ERROR: Invalid TAC value " TAC_FMT,
                  config_pP->served_tai.tac[i]);
            }
          }
        }
      }

      // sort TAI list
      n = config_pP->served_tai.nb_tai;
      do {
        stop_index = 0;
        for (i = 1; i < n; i++) {
          swap = false;
          if (config_pP->served_tai.plmn_mcc[i - 1] >
              config_pP->served_tai.plmn_mcc[i]) {
            swap = true;
          } else if (
              config_pP->served_tai.plmn_mcc[i - 1] ==
              config_pP->served_tai.plmn_mcc[i]) {
            if (config_pP->served_tai.plmn_mnc[i - 1] >
                config_pP->served_tai.plmn_mnc[i]) {
              swap = true;
            } else if (
                config_pP->served_tai.plmn_mnc[i - 1] ==
                config_pP->served_tai.plmn_mnc[i]) {
              if (config_pP->served_tai.tac[i - 1] >
                  config_pP->served_tai.tac[i]) {
                swap = true;
              }
            }
          }
          if (true == swap) {
            uint16_t swap16;
            swap16 = config_pP->served_tai.plmn_mcc[i - 1];
            config_pP->served_tai.plmn_mcc[i - 1] =
                config_pP->served_tai.plmn_mcc[i];
            config_pP->served_tai.plmn_mcc[i] = swap16;

            swap16 = config_pP->served_tai.plmn_mnc[i - 1];
            config_pP->served_tai.plmn_mnc[i - 1] =
                config_pP->served_tai.plmn_mnc[i];
            config_pP->served_tai.plmn_mnc[i] = swap16;

            swap16 = config_pP->served_tai.plmn_mnc_len[i - 1];
            config_pP->served_tai.plmn_mnc_len[i - 1] =
                config_pP->served_tai.plmn_mnc_len[i];
            config_pP->served_tai.plmn_mnc_len[i] = swap16;

            swap16                           = config_pP->served_tai.tac[i - 1];
            config_pP->served_tai.tac[i - 1] = config_pP->served_tai.tac[i];
            config_pP->served_tai.tac[i]     = swap16;

            stop_index = i;
          }
        }
        n = stop_index;
      } while (0 != n);

      // helper for determination of list type (global view), we could make
      // sublists with different types, but keep things simple for now
      config_pP->served_tai.list_type =
          TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_CONSECUTIVE_TACS;
      for (i = 1; i < config_pP->served_tai.nb_tai; i++) {
        if ((config_pP->served_tai.plmn_mcc[i] !=
             config_pP->served_tai.plmn_mcc[0]) ||
            (config_pP->served_tai.plmn_mnc[i] !=
             config_pP->served_tai.plmn_mnc[0])) {
          config_pP->served_tai.list_type =
              TRACKING_AREA_IDENTITY_LIST_TYPE_MANY_PLMNS;
          break;
        } else if (
            (config_pP->served_tai.plmn_mcc[i] !=
             config_pP->served_tai.plmn_mcc[i - 1]) ||
            (config_pP->served_tai.plmn_mnc[i] !=
             config_pP->served_tai.plmn_mnc[i - 1])) {
          config_pP->served_tai.list_type =
              TRACKING_AREA_IDENTITY_LIST_TYPE_MANY_PLMNS;
          break;
        }
        if (config_pP->served_tai.tac[i] !=
            (config_pP->served_tai.tac[i - 1] + 1)) {
          config_pP->served_tai.list_type =
              TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_NON_CONSECUTIVE_TACS;
        }
      }
    }
  }

  // GUAMFI SETTING
  setting =
      config_setting_get_member(setting_mme, MME_CONFIG_STRING_GUAMFI_LIST);
  config_pP->guamfi.nb = 0;
  if (setting != NULL) {
    num = config_setting_length(setting);
    OAILOG_DEBUG(LOG_AMF_APP, "Number of GUAMFIs configured =%d\n", num);
    AssertFatal(
        num >= MIN_GUAMI,
        "Not even one GUAMI is configured, configure minimum one GUMMEI \n");
    AssertFatal(
        num <= MAX_GUAMI,
        "Number of GUAMIs configured:%d exceeds number of GUMMEIs supported "
        ":%d \n",
        num, MAX_GUMMEI);

    for (i = 0; i < num; i++) {
      sub2setting = config_setting_get_elem(setting, i);

      if (sub2setting != NULL) {
        if ((config_setting_lookup_string(
                sub2setting, MME_CONFIG_STRING_MCC, &mcc))) {
          AssertFatal(
              strlen(mcc) == MAX_MCC_LENGTH,
              "Bad MCC length (%ld), it must be %u digit ex: 001", strlen(mcc),
              MAX_MCC_LENGTH);
          char c[2]                                   = {mcc[0], 0};
          config_pP->guamfi.guamfi[i].plmn.mcc_digit1 = (uint8_t) atoi(c);
          c[0]                                        = mcc[1];
          config_pP->guamfi.guamfi[i].plmn.mcc_digit2 = (uint8_t) atoi(c);
          c[0]                                        = mcc[2];
          config_pP->guamfi.guamfi[i].plmn.mcc_digit3 = (uint8_t) atoi(c);
        }

        if ((config_setting_lookup_string(
                sub2setting, MME_CONFIG_STRING_MNC, &mnc))) {
          AssertFatal(
              (strlen(mnc) == MIN_MNC_LENGTH) ||
                  (strlen(mnc) == MAX_MNC_LENGTH),
              "Bad MNC length (%ld), it must be %u or %u digit ex: 12 or 123",
              strlen(mnc), MIN_MNC_LENGTH, MAX_MNC_LENGTH);
          char c[2]                                   = {mnc[0], 0};
          config_pP->guamfi.guamfi[i].plmn.mnc_digit1 = (uint8_t) atoi(c);
          c[0]                                        = mnc[1];
          config_pP->guamfi.guamfi[i].plmn.mnc_digit2 = (uint8_t) atoi(c);
          if (3 == strlen(mnc)) {
            c[0]                                        = mnc[2];
            config_pP->guamfi.guamfi[i].plmn.mnc_digit3 = (uint8_t) atoi(c);
          } else {
            config_pP->guamfi.guamfi[i].plmn.mnc_digit3 = 0x0F;
          }
        }

        if ((config_setting_lookup_string(
                sub2setting, MME_CONFIG_STRING_AMF_REGION_ID, &mnc))) {
          config_pP->guamfi.guamfi[i].amf_regionid = (uint16_t) atoi(mnc);
        }
        if ((config_setting_lookup_string(
                sub2setting, MME_CONFIG_STRING_AMF_SET_ID, &mnc))) {
          config_pP->guamfi.guamfi[i].amf_set_id = (uint8_t) atoi(mnc);
        }
        if ((config_setting_lookup_string(
                sub2setting, MME_CONFIG_STRING_AMF_POINTER, &mnc))) {
          config_pP->guamfi.guamfi[i].amf_pointer = (uint8_t) atoi(mnc);
        }
        config_pP->guamfi.nb += 1;
      }
    }
  }

  setting_ngap = config_lookup(&cfg, NGAP_CONFIG_STRING_NGAP_CONFIG);
  if (setting_ngap != NULL) {
    if (config_setting_lookup_string(
            setting_ngap, NGAP_CONFIG_STRING_DEFAULT_DNS_IPV4_ADDRESS,
            (const char**) &default_dns) &&
        config_setting_lookup_string(
            setting_ngap, NGAP_CONFIG_STRING_DEFAULT_DNS_IPV4_ADDRESS,
            (const char**) &default_dns_sec)) {
      IPV4_STR_ADDR_TO_INADDR(
          default_dns, config_pP->ipv4.default_dns,
          "BAD IPv4 ADDRESS FORMAT FOR DEFAULT DNS !\n");
      IPV4_STR_ADDR_TO_INADDR(
          default_dns_sec, config_pP->ipv4.default_dns_sec,
          "BAD IPv4 ADDRESS FORMAT FOR DEFAULT DNS SEC!\n");
    }

    // AMF NAME
    if ((config_setting_lookup_string(
            setting_ngap, NGAP_CONFIG_AMF_NAME, (const char**) &astring))) {
      config_pP->amf_name = bfromcstr(astring);
    }

    // AMF_PLMN_SUPPORT SETTING
    setting = config_setting_get_member(
        setting_ngap, NGAP_CONFIG_AMF_PLMN_SUPPORT_LIST);
    config_pP->plmn_support_list.plmn_support_count = 0;
    if (setting != NULL) {
      num = config_setting_length(setting);
      OAILOG_DEBUG(LOG_AMF_APP, "Number of PLMN SUPPORT configured =%d\n", num);
      AssertFatal(
          num >= MIN_PLMN_SUPPORT,
          "Not even one PLMN SUPPORT is configured, configure minimum one PLMN "
          "LIST \n");
      AssertFatal(
          num <= MAX_PLMN_SUPPORT,
          "Number of PLMN SUPPPORT configured:%d exceeds number of PLMN_SUPPORT"
          ":%d \n",
          num, MAX_PLMN_SUPPORT);

      for (i = 0; i < num; i++) {
        sub2setting = config_setting_get_elem(setting, i);

        if (sub2setting != NULL) {
          if ((config_setting_lookup_string(
                  sub2setting, MME_CONFIG_STRING_MCC, &mcc))) {
            AssertFatal(
                strlen(mcc) == MAX_MCC_LENGTH,
                "Bad MCC length (%ld), it must be %u digit ex: 001",
                strlen(mcc), MAX_MCC_LENGTH);
            char c[2] = {mcc[0], 0};
            config_pP->plmn_support_list.plmn_support[i].plmn.mcc_digit1 =
                (uint8_t) atoi(c);
            c[0] = mcc[1];
            config_pP->plmn_support_list.plmn_support[i].plmn.mcc_digit2 =
                (uint8_t) atoi(c);
            c[0] = mcc[2];
            config_pP->plmn_support_list.plmn_support[i].plmn.mcc_digit3 =
                (uint8_t) atoi(c);
          }

          if ((config_setting_lookup_string(
                  sub2setting, MME_CONFIG_STRING_MNC, &mnc))) {
            AssertFatal(
                (strlen(mnc) == MIN_MNC_LENGTH) ||
                    (strlen(mnc) == MAX_MNC_LENGTH),
                "Bad MNC length (%ld), it must be %u or %u digit ex: 12 or 123",
                strlen(mnc), MIN_MNC_LENGTH, MAX_MNC_LENGTH);
            char c[2] = {mnc[0], 0};
            config_pP->plmn_support_list.plmn_support[i].plmn.mnc_digit1 =
                (uint8_t) atoi(c);
            c[0] = mnc[1];
            config_pP->plmn_support_list.plmn_support[i].plmn.mnc_digit2 =
                (uint8_t) atoi(c);
            if (3 == strlen(mnc)) {
              c[0] = mnc[2];
              config_pP->plmn_support_list.plmn_support[i].plmn.mnc_digit3 =
                  (uint8_t) atoi(c);
            } else {
              config_pP->plmn_support_list.plmn_support[i].plmn.mnc_digit3 =
                  0x0F;
            }
          }

          if (config_setting_lookup_string(
                  sub2setting, NGAP_CONFIG_PLMN_SUPPORT_SST, &set_sst)) {
            config_pP->plmn_support_list.plmn_support[i].s_nssai.sst =
                (uint8_t) atoi(set_sst);
          }

          if (config_setting_lookup_string(
                  sub2setting, NGAP_CONFIG_PLMN_SUPPORT_SD, &set_sd)) {
            uint64_t default_sd_val = 0;
            errno                   = 0;
            default_sd_val          = strtoll(set_sd, NULL, 16);
            AssertFatal(
                !(errno == ERANGE &&
                  (default_sd_val == LONG_MAX || default_sd_val == LONG_MIN)) ||
                    !(errno != 0 && default_sd_val == 0),
                "Slice Descriptor out of Range/Invalid");
            config_pP->plmn_support_list.plmn_support[i].s_nssai.sd.v =
                default_sd_val;
          }
          config_pP->plmn_support_list.plmn_support_count += 1;
        }  // If MCC/MNC/Slice Information is found
      }    // For the number of entries in the list for PLMN SUPPORT
    }      // PLMN_SUPPORT LIST is present
  }        // NGP Setting is not NULL

  config_destroy(&cfg);
  return 0;
}

static bool parse_bool(const char* str) {
  if (strcasecmp(str, "yes") == 0) return true;
  if (strcasecmp(str, "true") == 0) return true;
  if (strcasecmp(str, "no") == 0) return false;
  if (strcasecmp(str, "false") == 0) return false;
  if (strcasecmp(str, "") == 0) return false;

  Fatal("Error in config file: got \"%s\" but expected bool\n", str);
}

/***************************************************************************
**                                                                        **
** Name:   amf_config_free()                                              **
**                                                                        **
** Description: de-initializes the amf config                             **
**                                                                        **
**                                                                        **
***************************************************************************/
void amf_config_free(amf_config_t* amf_config) {
  free_wrapper((void**) &amf_config->served_tai.plmn_mcc);
  free_wrapper((void**) &amf_config->served_tai.plmn_mnc);
  free_wrapper((void**) &amf_config->served_tai.plmn_mnc_len);
  free_wrapper((void**) &amf_config->served_tai.tac);
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_config_exit()                                             **
 **                                                                        **
 ** Description: Cleanup configuration                                     **
 **                                                                        **
 ** Inputs:  void: no arguments                                            **
 **                                                                        **
 ***************************************************************************/
void amf_config_exit(void) {
  pthread_rwlock_destroy(&amf_config.rw_lock);
  amf_config_free(&amf_config);
}
