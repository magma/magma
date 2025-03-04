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

#include "lte/gateway/c/core/oai/include/amf_config.hpp"

#include <libconfig.h>
#include "lte/gateway/c/core/oai/common/log.h"
#include <errno.h>
#include "lte/gateway/c/core/common/assertions.h"
#include "lte/gateway/c/core/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/common/amf_default_values.h"
#include "lte/gateway/c/core/oai/include/TrackingAreaIdentity.h"
#include "lte/gateway/c/core/oai/include/mme_config.hpp"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.501.h"

void served_tai_config_init(served_tai_t* served_tai);
void clear_served_tai_config(served_tai_t* served_tai);

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
  log_conf->ngap_log_level = MAX_LOG_LEVEL;
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
  nas_conf->t3502_min = T3502_DEFAULT_VALUE;
  nas_conf->t3512_min = T3512_DEFAULT_VALUE;
  nas_conf->t3522_sec = T3522_DEFAULT_VALUE;
  nas_conf->t3550_sec = T3550_DEFAULT_VALUE;
  nas_conf->t3560_sec = T3560_DEFAULT_VALUE;
  nas_conf->t3570_sec = T3570_DEFAULT_VALUE;
  nas_conf->t3585_sec = T3585_DEFAULT_VALUE;
  nas_conf->t3586_sec = T3586_DEFAULT_VALUE;
  nas_conf->t3589_sec = T3589_DEFAULT_VALUE;
  nas_conf->t3595_sec = T3595_DEFAULT_VALUE;
  nas_conf->implicit_dereg_sec = IMPLICIT_DEREG_TIMER_VALUE;
  nas_conf->force_reject_tau = true;
  nas_conf->force_reject_sr = true;
  nas_conf->disable_esm_information = false;
  nas_conf->enable_IMS_VoPS_3GPP = true;
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
  guamfi_conf->nb = 1;
  guamfi_conf->guamfi[0].amf_set_id = AMFC;
  guamfi_conf->guamfi[0].amf_regionid = AMFGID;
  guamfi_conf->guamfi[0].amf_pointer = AMFPOINTER;
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
  plmn_support_list->plmn_support_count = MIN_PLMN_SUPPORT;
  plmn_support_list->plmn_support[0].plmn.mcc_digit1 = 0;
  plmn_support_list->plmn_support[0].plmn.mcc_digit2 = 0;
  plmn_support_list->plmn_support[0].plmn.mcc_digit3 = 0;
  plmn_support_list->plmn_support[0].plmn.mcc_digit1 = 0;
  plmn_support_list->plmn_support[0].plmn.mcc_digit2 = 0;
  plmn_support_list->plmn_support[0].plmn.mcc_digit3 = 0x0F;
  plmn_support_list->plmn_support[0].s_nssai.sst = AMF_S_NSSAI_ST_DEFAULT_VALUE;
  plmn_support_list->plmn_support[0].s_nssai.sd.v =
      AMF_S_NSSAI_SD_INVALID_VALUE;
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
  ngap_conf->port_number = NGAP_PORT_NUMBER;
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
#ifdef __cplusplus
extern "C" {
#endif
void amf_config_init(amf_config_t* config) {
  memset(config, 0, sizeof(*config));

  pthread_rwlock_init(&config->rw_lock, NULL);

  config->max_gnbs = 2;
  config->max_ues = 2;
  config->unauthenticated_imsi_supported = 0;
  config->relative_capacity = RELATIVE_CAPACITY;
  config->amf_statistic_timer = AMF_STATISTIC_TIMER_S;
  config->use_stateless = false;
  ngap_config_init(&config->ngap_config);
  nas5g_config_init(&config->nas_config);
  guamfi_config_init(&config->guamfi);
  plmn_support_list_config_init(&config->plmn_support_list);
  served_tai_config_init(&config->served_tai);
}
#ifdef __cplusplus
}
#endif
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
** Name:    amf_config_parse_opt_line()                                   **
**                                                                        **
** Description: Invokes amf_config_init() to initialize                   **
**              default values of AMF                                     **
**                                                                        **
**                                                                        **
***************************************************************************/
#ifdef __cplusplus
extern "C" {
#endif
int amf_config_parse_file(amf_config_t* config_pP,
                          const mme_config_t* mme_config_p) {
  config_t cfg = {0};
  config_setting_t* setting = NULL;
  config_setting_t* sub2setting = NULL;
  config_setting_t* setting_amf = NULL;
  int i = 0, num = 0;
  const char* astring = NULL;
  const char* mcc = NULL;
  const char* mnc = NULL;
  const char* region_id = NULL;
  const char* set_id = NULL;
  const char* pointer = NULL;
  const char* default_dns = NULL;
  const char* default_pcscf = NULL;
  const char* default_dns_sec = NULL;
  const char* set_sst = NULL;
  const char* set_sd = NULL;
  int aint = 0;
  config_init(&cfg);

  if (config_pP->config_file != NULL) {
    /*
     * Read the file. If there is an error, report it and exit.
     */
    if (!config_read_file(&cfg, bdata(config_pP->config_file))) {
      OAILOG_CRITICAL(LOG_CONFIG,
                      "Failed to parse AMF configuration file: %s:%d - %s\n",
                      bdata(config_pP->config_file), config_error_line(&cfg),
                      config_error_text(&cfg));
      config_destroy(&cfg);
      AssertFatal(1 == 0, "Failed to parse AMF configuration file %s!\n",
                  bdata(config_pP->config_file));
    }
  } else {
    config_destroy(&cfg);
    AssertFatal(0, "No AMF configuration file provided!\n");
  }

  copy_amf_config_from_mme_config(config_pP, mme_config_p);

  setting_amf = config_lookup(&cfg, AMF_CONFIG_STRING_AMF_CONFIG);
  if (setting_amf != NULL) {
    if (config_setting_lookup_string(setting_amf,
                                     AMF_CONFIG_STRING_DEFAULT_DNS_IPV4_ADDRESS,
                                     (const char**)&default_dns) &&
        config_setting_lookup_string(
            setting_amf, AMF_CONFIG_STRING_DEFAULT_DNS_SEC_IPV4_ADDRESS,
            (const char**)&default_dns_sec)) {
      IPV4_STR_ADDR_TO_INADDR(default_dns, config_pP->ipv4.default_dns,
                              "BAD IPv4 ADDRESS FORMAT FOR DEFAULT DNS !\n");
      IPV4_STR_ADDR_TO_INADDR(default_dns_sec, config_pP->ipv4.default_dns_sec,
                              "BAD IPv4 ADDRESS FORMAT FOR DEFAULT DNS SEC!\n");
    }

    if (config_setting_lookup_string(
            setting_amf, AMF_CONFIG_STRING_DEFAULT_PCSCF_IPV4_ADDRESS,
            (const char**)&default_pcscf)) {
      IPV4_STR_ADDR_TO_INADDR(default_pcscf, config_pP->pcscf_addr.ipv4,
                              "BAD IPv4 ADDRESS FORMAT FOR DEFAULT PCSCF !\n");
    }

    // AMF NAME
    if ((config_setting_lookup_string(setting_amf, AMF_CONFIG_AMF_NAME,
                                      (const char**)&astring))) {
      config_pP->amf_name = bfromcstr(astring);
    }

    // DEFAULT_DNN
    if (config_setting_lookup_string(setting_amf, CONFIG_DEFAULT_DNN,
                                     (const char**)&astring)) {
      config_pP->default_dnn = bfromcstr(astring);
    }
    // DEFAULT AUTH MAX RETRY COUNT
    if (config_setting_lookup_string(
            setting_amf, AUTHENTICATION_COUNTER_MAX_RETRY, &astring)) {
      config_pP->auth_retry_max_count = (uint32_t)atoi(astring);
    }

    // DEFAULT AUTH RETRY TIMER EXPIRES MSECS
    if (config_setting_lookup_string(
            setting_amf, AUTHENTICATION_RETRY_TIMER_EXPIRY_MSECS, &astring)) {
      config_pP->auth_retry_interval = (uint32_t)atoi(astring);
    }

    // AMF_PLMN_SUPPORT SETTING
    setting = config_setting_get_member(setting_amf,
                                        AMF_CONFIG_AMF_PLMN_SUPPORT_LIST);
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
          if ((config_setting_lookup_string(sub2setting, MME_CONFIG_STRING_MCC,
                                            &mcc))) {
            AssertFatal(strlen(mcc) == MAX_MCC_LENGTH,
                        "Bad MCC length (%ld), it must be %u digit ex: 001",
                        strlen(mcc), MAX_MCC_LENGTH);
            char c[2] = {mcc[0], 0};
            config_pP->plmn_support_list.plmn_support[i].plmn.mcc_digit1 =
                (uint8_t)atoi(c);
            c[0] = mcc[1];
            config_pP->plmn_support_list.plmn_support[i].plmn.mcc_digit2 =
                (uint8_t)atoi(c);
            c[0] = mcc[2];
            config_pP->plmn_support_list.plmn_support[i].plmn.mcc_digit3 =
                (uint8_t)atoi(c);
          }

          if ((config_setting_lookup_string(sub2setting, MME_CONFIG_STRING_MNC,
                                            &mnc))) {
            AssertFatal(
                (strlen(mnc) == MIN_MNC_LENGTH) ||
                    (strlen(mnc) == MAX_MNC_LENGTH),
                "Bad MNC length (%ld), it must be %u or %u digit ex: 12 or 123",
                strlen(mnc), MIN_MNC_LENGTH, MAX_MNC_LENGTH);
            char c[2] = {mnc[0], 0};
            config_pP->plmn_support_list.plmn_support[i].plmn.mnc_digit1 =
                (uint8_t)atoi(c);
            c[0] = mnc[1];
            config_pP->plmn_support_list.plmn_support[i].plmn.mnc_digit2 =
                (uint8_t)atoi(c);
            if (3 == strlen(mnc)) {
              c[0] = mnc[2];
              config_pP->plmn_support_list.plmn_support[i].plmn.mnc_digit3 =
                  (uint8_t)atoi(c);
            } else {
              config_pP->plmn_support_list.plmn_support[i].plmn.mnc_digit3 =
                  0x0F;
            }
          }

          if (config_setting_lookup_string(
                  sub2setting, AMF_CONFIG_PLMN_SUPPORT_SST, &set_sst)) {
            config_pP->plmn_support_list.plmn_support[i].s_nssai.sst =
                (uint8_t)atoi(set_sst);
          }

          if (config_setting_lookup_string(
                  sub2setting, AMF_CONFIG_PLMN_SUPPORT_SD, &set_sd)) {
            uint64_t default_sd_val = 0;
            errno = 0;
            default_sd_val = strtoll(set_sd, NULL, 16);
            AssertFatal(!(errno == ERANGE && (default_sd_val == LONG_MAX ||
                                              default_sd_val == LONG_MIN)) ||
                            !(errno != 0 && default_sd_val == 0),
                        "Slice Descriptor out of Range/Invalid");
            config_pP->plmn_support_list.plmn_support[i].s_nssai.sd.v =
                default_sd_val;
          }
          config_pP->plmn_support_list.plmn_support_count += 1;
        }  // If MCC/MNC/Slice Information is found
      }  // For the number of entries in the list for PLMN SUPPORT
    }  // PLMN_SUPPORT LIST is present

    // enable VoNR support
    if ((config_setting_lookup_string(
            setting_amf, AMF_CONFIG_STRING_NAS_ENABLE_IMS_VoPS_3GPP,
            &astring))) {
      config_pP->nas_config.enable_IMS_VoPS_3GPP = parse_bool(astring);
    }

    // t3512
    if ((config_setting_lookup_int(setting_amf, AMF_CONFIG_STRING_NAS_T3512,
                                   &aint))) {
      config_pP->nas_config.t3512_min = (uint32_t)aint;
    }

    // guamfi SETTING
    setting =
        config_setting_get_member(setting_amf, AMF_CONFIG_STRING_GUAMFI_LIST);
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
          if ((config_setting_lookup_string(sub2setting, MME_CONFIG_STRING_MCC,
                                            &mcc))) {
            AssertFatal(strlen(mcc) == MAX_MCC_LENGTH,
                        "Bad MCC length (%ld), it must be %u digit ex: 001",
                        strlen(mcc), MAX_MCC_LENGTH);
            char c[2] = {mcc[0], 0};
            config_pP->guamfi.guamfi[i].plmn.mcc_digit1 = (uint8_t)atoi(c);
            c[0] = mcc[1];
            config_pP->guamfi.guamfi[i].plmn.mcc_digit2 = (uint8_t)atoi(c);
            c[0] = mcc[2];
            config_pP->guamfi.guamfi[i].plmn.mcc_digit3 = (uint8_t)atoi(c);
          }

          if ((config_setting_lookup_string(sub2setting, MME_CONFIG_STRING_MNC,
                                            &mnc))) {
            AssertFatal(
                (strlen(mnc) == MIN_MNC_LENGTH) ||
                    (strlen(mnc) == MAX_MNC_LENGTH),
                "Bad MNC length (%ld), it must be %u or %u digit ex: 12 or 123",
                strlen(mnc), MIN_MNC_LENGTH, MAX_MNC_LENGTH);
            char c[2] = {mnc[0], 0};
            config_pP->guamfi.guamfi[i].plmn.mnc_digit1 = (uint8_t)atoi(c);
            c[0] = mnc[1];
            config_pP->guamfi.guamfi[i].plmn.mnc_digit2 = (uint8_t)atoi(c);
            if (3 == strlen(mnc)) {
              c[0] = mnc[2];
              config_pP->guamfi.guamfi[i].plmn.mnc_digit3 = (uint8_t)atoi(c);
            } else {
              config_pP->guamfi.guamfi[i].plmn.mnc_digit3 = 0x0F;
            }
          }

          if ((config_setting_lookup_string(
                  sub2setting, AMF_CONFIG_STRING_AMF_REGION_ID, &region_id))) {
            config_pP->guamfi.guamfi[i].amf_regionid =
                (uint16_t)atoi(region_id);
          }
          if ((config_setting_lookup_string(
                  sub2setting, AMF_CONFIG_STRING_AMF_SET_ID, &set_id))) {
            config_pP->guamfi.guamfi[i].amf_set_id = (uint8_t)atoi(set_id);
          }
          if ((config_setting_lookup_string(
                  sub2setting, AMF_CONFIG_STRING_AMF_POINTER, &pointer))) {
            config_pP->guamfi.guamfi[i].amf_pointer = (uint8_t)atoi(pointer);
          }

          config_pP->guamfi.nb += 1;
        }
      }
    }
  }  // NGP Setting is not NULL
  config_destroy(&cfg);
  return 0;
}
#ifdef __cplusplus
}
#endif

/***************************************************************************
**                                                                        **
** Name:   amf_config_free()                                              **
**                                                                        **
** Description: de-initializes the amf config                             **
**                                                                        **
**                                                                        **
***************************************************************************/
void amf_config_free(amf_config_t* amf_config) {
  free_wrapper((void**)&amf_config->served_tai.plmn_mcc);
  free_wrapper((void**)&amf_config->served_tai.plmn_mnc);
  free_wrapper((void**)&amf_config->served_tai.plmn_mnc_len);
  free_wrapper((void**)&amf_config->served_tai.tac);
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

void clear_amf_config(amf_config_t* amf_config) {
  if (!amf_config) return;

  bdestroy_wrapper(&amf_config->log_config.output);
  bdestroy_wrapper(&amf_config->config_file);
  bdestroy_wrapper(&amf_config->pid_dir);
  bdestroy_wrapper(&amf_config->realm);
  bdestroy_wrapper(&amf_config->full_network_name);
  bdestroy_wrapper(&amf_config->short_network_name);
  bdestroy_wrapper(&amf_config->ip_capability);
  bdestroy_wrapper(&amf_config->amf_name);
  bdestroy_wrapper(&amf_config->default_dnn);
  clear_served_tai_config(&amf_config->served_tai);
  free_partial_lists(amf_config->partial_list, amf_config->num_par_lists);
  amf_config->num_par_lists = 0;
}

/***************************************************************************
**                                                                        **
** Name:   copy_served_tai_config_list()                                  **
**                                                                        **
** Description: copy tai list from mme_config to amf_config               **
**                                                                        **
**                                                                        **
***************************************************************************/
void copy_served_tai_config_list(amf_config_t* dest, const mme_config_t* src) {
  if (!dest || !src) return;

  // served_tai
  if (dest->served_tai.nb_tai != src->served_tai.nb_tai) {
    if (NULL != dest->served_tai.plmn_mcc)
      free_wrapper((void**)&dest->served_tai.plmn_mcc);

    if (NULL != dest->served_tai.plmn_mnc)
      free_wrapper((void**)&dest->served_tai.plmn_mnc);

    if (NULL != dest->served_tai.plmn_mnc_len)
      free_wrapper((void**)&dest->served_tai.plmn_mnc_len);

    if (NULL != dest->served_tai.tac)
      free_wrapper((void**)&dest->served_tai.tac);

    dest->served_tai.nb_tai = src->served_tai.nb_tai;
    dest->served_tai.plmn_mcc =
        calloc(dest->served_tai.nb_tai, sizeof(uint16_t));
    dest->served_tai.plmn_mnc =
        calloc(dest->served_tai.nb_tai, sizeof(uint16_t));
    dest->served_tai.plmn_mnc_len =
        calloc(dest->served_tai.nb_tai, sizeof(uint16_t));
    dest->served_tai.tac = calloc(dest->served_tai.nb_tai, sizeof(uint16_t));
  }
  memcpy(dest->served_tai.plmn_mcc, src->served_tai.plmn_mcc,
         (dest->served_tai.nb_tai) * sizeof(uint16_t));
  memcpy(dest->served_tai.plmn_mnc, src->served_tai.plmn_mnc,
         (dest->served_tai.nb_tai) * sizeof(uint16_t));
  memcpy(dest->served_tai.plmn_mnc_len, src->served_tai.plmn_mnc_len,
         (dest->served_tai.nb_tai) * sizeof(uint16_t));
  memcpy(dest->served_tai.tac, src->served_tai.tac,
         (dest->served_tai.nb_tai) * sizeof(uint16_t));

  // num_par_lists
  dest->num_par_lists = src->num_par_lists;

  // partial_list
  dest->partial_list = calloc(dest->num_par_lists, sizeof(partial_list_t));

  for (uint8_t itr = 0; itr < src->num_par_lists && dest->partial_list; ++itr) {
    dest->partial_list[itr].list_type = src->partial_list[itr].list_type;
    dest->partial_list[itr].nb_elem = src->partial_list[itr].nb_elem;

    dest->partial_list[itr].plmn =
        calloc(dest->partial_list[itr].nb_elem, sizeof(plmn_t));
    memcpy(dest->partial_list[itr].plmn, src->partial_list[itr].plmn,
           (dest->partial_list[itr].nb_elem) * sizeof(plmn_t));

    dest->partial_list[itr].tac =
        calloc(dest->partial_list[itr].nb_elem, sizeof(tac_t));
    memcpy(dest->partial_list[itr].tac, src->partial_list[itr].tac,
           (dest->partial_list[itr].nb_elem) * sizeof(tac_t));
  }
}

void copy_amf_config_from_mme_config(amf_config_t* dest,
                                     const mme_config_t* src) {
  OAILOG_DEBUG(LOG_AMF_APP, "copy_amf_config_from_mme_config");
  // LOGGING setting
  dest->log_config = src->log_config;
  if (src->log_config.output)
    dest->log_config.output = bstrcpy(src->log_config.output);
  dest->log_config.amf_app_log_level = src->log_config.mme_app_log_level;

  // GENERAL AMF SETTINGS
  dest->realm = bstrcpy(src->realm);
  if (src->full_network_name)
    dest->full_network_name = bstrcpy(src->full_network_name);
  if (src->short_network_name)
    dest->short_network_name = bstrcpy(src->short_network_name);
  dest->daylight_saving_time = src->daylight_saving_time;
  if (src->pid_dir) dest->pid_dir = bstrcpy(src->pid_dir);
  dest->max_gnbs = src->max_enbs;
  dest->max_ues = src->max_ues;
  dest->relative_capacity = src->relative_capacity;
  dest->use_stateless = src->use_stateless;
  dest->unauthenticated_imsi_supported = src->unauthenticated_imsi_supported;

  // NAS-5G setting
  for (int i = 0; i < 8; i++) {
    dest->nas_config.preferred_integrity_algorithm[i] =
        src->nas_config.prefered_integrity_algorithm[i];
    dest->nas_config.preferred_ciphering_algorithm[i] =
        src->nas_config.prefered_ciphering_algorithm[i];
  }

  // TAI list setting
  copy_served_tai_config_list(dest, src);
}

void amf_config_display(amf_config_t* config_pP) {
  if (!config_pP) return;

  OAILOG_INFO(LOG_CONFIG, "==========AMF Configuration Start==========\n");

  OAILOG_INFO(LOG_CONFIG, "- Realm ................................: %s\n",
              bdata(config_pP->realm));
  OAILOG_INFO(LOG_CONFIG, "  full network name ....................: %s\n",
              bdata(config_pP->full_network_name));
  OAILOG_INFO(LOG_CONFIG, "  short network name ...................: %s\n",
              bdata(config_pP->short_network_name));
  OAILOG_INFO(LOG_CONFIG, "  Daylight Saving Time..................: %d\n",
              config_pP->daylight_saving_time);

  OAILOG_INFO(LOG_CONFIG, "- Max gNBs .............................: %u\n",
              config_pP->max_gnbs);
  OAILOG_INFO(LOG_CONFIG, "- Max UEs ..............................: %u\n",
              config_pP->max_ues);

  OAILOG_INFO(LOG_CONFIG, "- Use Stateless ........................: %s\n\n",
              config_pP->use_stateless ? "true" : "false");

  OAILOG_DEBUG(LOG_CONFIG, "- PARTIAL TAIs\n");
  OAILOG_DEBUG(LOG_CONFIG, "- Num of partial lists=%d\n",
               config_pP->num_par_lists);
  for (uint8_t itr = 0; itr < config_pP->num_par_lists; itr++) {
    if (config_pP->partial_list) {
      switch (config_pP->partial_list[itr].list_type) {
        case TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_CONSECUTIVE_TACS:
          OAILOG_DEBUG(
              LOG_CONFIG,
              "- List [%d] - TAI list type one PLMN consecutive TACs\n", itr);
          break;
        case TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_NON_CONSECUTIVE_TACS:
          OAILOG_DEBUG(
              LOG_CONFIG,
              "- List [%d] - TAI list type one PLMN non consecutive TACs\n",
              itr);
          break;
        case TRACKING_AREA_IDENTITY_LIST_TYPE_MANY_PLMNS:
          OAILOG_DEBUG(LOG_CONFIG,
                       "- List [%d] - TAI list type multiple PLMNs\n", itr);
          break;
        default:
          OAILOG_ERROR(LOG_CONFIG,
                       "Invalid served TAI list type (%u) configured\n",
                       config_pP->partial_list[itr].list_type);
          break;
      }
    }
  }

  for (uint8_t itr = 0; itr < config_pP->num_par_lists; ++itr) {
    OAILOG_DEBUG(LOG_CONFIG, "- Num of elements in list[%d]=%d\n", itr,
                 config_pP->partial_list[itr].nb_elem);
    if (config_pP->partial_list) {
      for (uint8_t idx = 0; idx < config_pP->partial_list[itr].nb_elem; ++idx) {
        if (config_pP->partial_list[itr].plmn &&
            config_pP->partial_list[itr].tac) {
          OAILOG_DEBUG(LOG_CONFIG,
                       "            "
                       "MCC1=%d\tMCC2=%d\tMCC3=%d\tMNC1=%d\tMNC2=%d\tMNC3=%d\t"
                       "TAC=%d\n",
                       config_pP->partial_list[itr].plmn[idx].mcc_digit1,
                       config_pP->partial_list[itr].plmn[idx].mcc_digit2,
                       config_pP->partial_list[itr].plmn[idx].mcc_digit3,
                       config_pP->partial_list[itr].plmn[idx].mnc_digit1,
                       config_pP->partial_list[itr].plmn[idx].mnc_digit2,
                       config_pP->partial_list[itr].plmn[idx].mnc_digit3,
                       config_pP->partial_list[itr].tac[idx]);
        }
      }
    }
  }

  OAILOG_INFO(LOG_CONFIG, "- Logging:\n");
  OAILOG_INFO(LOG_CONFIG, "    Output ..............: %s\n",
              bdata(config_pP->log_config.output));
  OAILOG_INFO(LOG_CONFIG, "    Output thread safe ..: %s\n",
              (config_pP->log_config.is_output_thread_safe) ? "true" : "false");
  OAILOG_INFO(LOG_CONFIG, "    Output with color ...: %s\n",
              (config_pP->log_config.color) ? "true" : "false");
  OAILOG_INFO(LOG_CONFIG, "    UDP log level........: %s\n",
              OAILOG_LEVEL_INT2STR(config_pP->log_config.udp_log_level));
  OAILOG_INFO(LOG_CONFIG, "    GTPV1-U log level....: %s\n",
              OAILOG_LEVEL_INT2STR(config_pP->log_config.gtpv1u_log_level));
  OAILOG_INFO(LOG_CONFIG, "    GTPV2-C log level....: %s\n",
              OAILOG_LEVEL_INT2STR(config_pP->log_config.gtpv2c_log_level));
  OAILOG_INFO(LOG_CONFIG, "    SCTP log level.......: %s\n",
              OAILOG_LEVEL_INT2STR(config_pP->log_config.sctp_log_level));
  OAILOG_INFO(LOG_CONFIG, "    S1AP log level.......: %s\n",
              OAILOG_LEVEL_INT2STR(config_pP->log_config.s1ap_log_level));
  OAILOG_INFO(LOG_CONFIG, "    ASN1 Verbosity level : %d\n",
              config_pP->log_config.asn1_verbosity_level);
  OAILOG_INFO(LOG_CONFIG, "    NAS log level........: %s\n",
              OAILOG_LEVEL_INT2STR(config_pP->log_config.nas_log_level));
  OAILOG_INFO(LOG_CONFIG, "    AMF_APP log level....: %s\n",
              OAILOG_LEVEL_INT2STR(config_pP->log_config.amf_app_log_level));
  OAILOG_INFO(LOG_CONFIG, "    SPGW_APP log level....: %s\n",
              OAILOG_LEVEL_INT2STR(config_pP->log_config.spgw_app_log_level));
  OAILOG_INFO(LOG_CONFIG, "    S11 log level........: %s\n",
              OAILOG_LEVEL_INT2STR(config_pP->log_config.s11_log_level));
  OAILOG_INFO(LOG_CONFIG, "    S6a log level........: %s\n",
              OAILOG_LEVEL_INT2STR(config_pP->log_config.s6a_log_level));
  OAILOG_INFO(LOG_CONFIG, "    UTIL log level.......: %s\n",
              OAILOG_LEVEL_INT2STR(config_pP->log_config.util_log_level));
  OAILOG_INFO(LOG_CONFIG,
              "    ITTI log level.......: %s (InTer-Task Interface)\n",
              OAILOG_LEVEL_INT2STR(config_pP->log_config.itti_log_level));

  OAILOG_INFO(LOG_CONFIG, "==========AMF Configuration End==========\n");
}
