/*
 * Copyright (c) 2015, EURECOM (www.eurecom.fr)
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice, this
 *    list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
 * ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 * The views and conclusions contained in the software and documentation are those
 * of the authors and should not be interpreted as representing official policies,
 * either expressed or implied, of the FreeBSD Project.
 */

/*! \file mme_config.c
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#if HAVE_CONFIG_H
#include "config.h"
#endif

#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <stdbool.h>
#include <unistd.h>
#include <string.h>
#include <arpa/inet.h> /* To provide inet_addr */
#include <pthread.h>
#include <libconfig.h>
#include <netinet/in.h>

#include "assertions.h"
#include "dynamic_memory_check.h"
#include "log.h"
#include "common_defs.h"
#include "mme_config.h"
#include "3gpp_33.401.h"
#include "intertask_interface_conf.h"
#include "3gpp_23.003.h"
#include "3gpp_24.008.h"
#include "3gpp_24.301.h"
#include "TrackingAreaIdentity.h"
#include "bstrlib.h"
#include "mme_default_values.h"
#include "service303.h"
#if EMBEDDED_SGW
#include "sgw_config.h"
#endif
static bool parse_bool(const char *str);

struct mme_config_s mme_config = {.rw_lock = PTHREAD_RWLOCK_INITIALIZER, 0};

//------------------------------------------------------------------------------
int mme_config_find_mnc_length(
  const char mcc_digit1P,
  const char mcc_digit2P,
  const char mcc_digit3P,
  const char mnc_digit1P,
  const char mnc_digit2P,
  const char mnc_digit3P)
{
  uint16_t mcc = 100 * mcc_digit1P + 10 * mcc_digit2P + mcc_digit3P;
  uint16_t mnc3 = 100 * mnc_digit1P + 10 * mnc_digit2P + mnc_digit3P;
  uint16_t mnc2 = 10 * mnc_digit1P + mnc_digit2P;
  int plmn_index = 0;

  if (
    mcc_digit1P < 0 || mcc_digit1P > 9 || mcc_digit2P < 0 || mcc_digit2P > 9 ||
    mcc_digit3P < 0 || mcc_digit3P > 9) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "BAD MCC PARAMETER (%d%d%d)!\n",
      mcc_digit1P,
      mcc_digit2P,
      mcc_digit3P);
    return 0;
  }
  if (
    mnc_digit2P < 0 || mnc_digit2P > 9 || mnc_digit1P < 0 || mnc_digit1P > 9) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "BAD MNC PARAMETER (%d%d%d)!\n",
      mnc_digit1P,
      mnc_digit2P,
      mnc_digit3P);
    return 0;
  }

  while (plmn_index < mme_config.served_tai.nb_tai) {
    if (mme_config.served_tai.plmn_mcc[plmn_index] == mcc) {
      if (
        (mme_config.served_tai.plmn_mnc[plmn_index] == mnc2) &&
        (mme_config.served_tai.plmn_mnc_len[plmn_index] == 2)) {
        return 2;
      } else if (
        (mme_config.served_tai.plmn_mnc[plmn_index] == mnc3) &&
        (mme_config.served_tai.plmn_mnc_len[plmn_index] == 3)) {
        return 3;
      }
    }

    plmn_index += 1;
  }

  return 0;
}

void log_config_init(log_config_t *log_conf)
{
  memset(log_conf, 0, sizeof(*log_conf));

  log_conf->output = NULL;
  log_conf->is_output_thread_safe = false;
  log_conf->color = false;

  log_conf->udp_log_level = MAX_LOG_LEVEL; // Means invalid TODO wtf
  log_conf->gtpv1u_log_level = MAX_LOG_LEVEL;
  log_conf->gtpv2c_log_level = MAX_LOG_LEVEL;
  log_conf->sctp_log_level = MAX_LOG_LEVEL;
  log_conf->s1ap_log_level = MAX_LOG_LEVEL;
  log_conf->nas_log_level = MAX_LOG_LEVEL;
  log_conf->mme_app_log_level = MAX_LOG_LEVEL;
  log_conf->s11_log_level = MAX_LOG_LEVEL;
  log_conf->s6a_log_level = MAX_LOG_LEVEL;
  log_conf->secu_log_level = MAX_LOG_LEVEL;
  log_conf->util_log_level = MAX_LOG_LEVEL;
  log_conf->itti_log_level = MAX_LOG_LEVEL;
  log_conf->spgw_app_log_level = MAX_LOG_LEVEL;
  log_conf->pgw_app_log_level = MAX_LOG_LEVEL;
  log_conf->asn1_verbosity_level = 0;
}

void eps_network_feature_config_init(eps_network_feature_config_t *eps_conf)
{
  eps_conf->emergency_bearer_services_in_s1_mode = 0;
  eps_conf->extended_service_request = 0;
  eps_conf->ims_voice_over_ps_session_in_s1 = 0;
  eps_conf->location_services_via_epc = 0;
}

void ipv4_config_init(ip_t *ip)
{
  memset(ip, 0, sizeof(*ip));

  ip->if_name_s1_mme = NULL;
  ip->s1_mme_v4.s_addr = INADDR_ANY;

  ip->if_name_s11 = NULL;
  ip->s11_mme_v4.s_addr = INADDR_ANY;

  ip->port_s11 = 2123;
}

void s1ap_config_init(s1ap_config_t *s1ap_conf)
{
  s1ap_conf->port_number = S1AP_PORT_NUMBER;
  s1ap_conf->outcome_drop_timer_sec = S1AP_OUTCOME_TIMER_DEFAULT;
}

void s6a_config_init(s6a_config_t *s6a_conf)
{
  s6a_conf->hss_host_name = NULL;
  s6a_conf->conf_file = bfromcstr(S6A_CONF_FILE);
}

void itti_config_init(itti_config_t *itti_conf)
{
  itti_conf->queue_size = ITTI_QUEUE_MAX_ELEMENTS;
  itti_conf->log_file = NULL;
}

void sctp_config_init(sctp_config_t *sctp_conf)
{
  sctp_conf->in_streams = SCTP_IN_STREAMS;
  sctp_conf->out_streams = SCTP_OUT_STREAMS;
}

void nas_config_init(nas_config_t *nas_conf)
{
  nas_conf->t3402_min = T3402_DEFAULT_VALUE;
  nas_conf->t3412_min = T3412_DEFAULT_VALUE;
  nas_conf->t3422_sec = T3422_DEFAULT_VALUE;
  nas_conf->t3450_sec = T3450_DEFAULT_VALUE;
  nas_conf->t3460_sec = T3460_DEFAULT_VALUE;
  nas_conf->t3470_sec = T3470_DEFAULT_VALUE;
  nas_conf->t3485_sec = T3485_DEFAULT_VALUE;
  nas_conf->t3486_sec = T3486_DEFAULT_VALUE;
  nas_conf->t3489_sec = T3489_DEFAULT_VALUE;
  nas_conf->t3495_sec = T3495_DEFAULT_VALUE;
  nas_conf->force_reject_tau = true;
  nas_conf->force_reject_sr = true;
  nas_conf->disable_esm_information = false;
}

void gummei_config_init(gummei_config_t *gummei_conf)
{
  gummei_conf->nb = 1;
  gummei_conf->gummei[0].mme_code = MMEC;
  gummei_conf->gummei[0].mme_gid = MMEGID;
  gummei_conf->gummei[0].plmn.mcc_digit1 = 0;
  gummei_conf->gummei[0].plmn.mcc_digit2 = 0;
  gummei_conf->gummei[0].plmn.mcc_digit3 = 1;
  gummei_conf->gummei[0].plmn.mcc_digit1 = 0;
  gummei_conf->gummei[0].plmn.mcc_digit2 = 1;
  gummei_conf->gummei[0].plmn.mcc_digit3 = 0x0F;
}

void served_tai_config_init(served_tai_t *served_tai)
{
  served_tai->nb_tai = 1;
  served_tai->plmn_mcc = calloc(1, sizeof(*served_tai->plmn_mcc));
  served_tai->plmn_mnc = calloc(1, sizeof(*served_tai->plmn_mnc));
  served_tai->plmn_mnc_len = calloc(1, sizeof(*served_tai->plmn_mnc_len));
  served_tai->tac = calloc(1, sizeof(*served_tai->tac));
  served_tai->plmn_mcc[0] = PLMN_MCC;
  served_tai->plmn_mnc[0] = PLMN_MNC;
  served_tai->plmn_mnc_len[0] = PLMN_MNC_LEN;
  served_tai->tac[0] = PLMN_TAC;
}

void service303_config_init(service303_data_t *service303_conf)
{
  service303_conf->name = bfromcstr(SERVICE303_MME_PACKAGE_NAME);
  service303_conf->version = bfromcstr(SERVICE303_MME_PACKAGE_VERSION);
}

//------------------------------------------------------------------------------
void mme_config_init(mme_config_t *config)
{
  memset(config, 0, sizeof(*config));

  pthread_rwlock_init(&config->rw_lock, NULL);

  config->config_file = NULL;
  config->max_enbs = 2;
  config->max_ues = 2;
  config->unauthenticated_imsi_supported = 0;
  config->relative_capacity = RELATIVE_CAPACITY;
  config->mme_statistic_timer = MME_STATISTIC_TIMER_S;

  log_config_init(&config->log_config);
  eps_network_feature_config_init(&config->eps_network_feature_support);
  ipv4_config_init(&config->ip);
  s1ap_config_init(&config->s1ap_config);
  s6a_config_init(&config->s6a_config);
  itti_config_init(&config->itti_config);
  sctp_config_init(&config->sctp_config);
  nas_config_init(&config->nas_config);
  gummei_config_init(&config->gummei);
  served_tai_config_init(&config->served_tai);
  service303_config_init(&config->service303_config);
}

//------------------------------------------------------------------------------
void mme_config_exit(void)
{
  pthread_rwlock_destroy(&mme_config.rw_lock);
  bdestroy_wrapper(&mme_config.log_config.output);
  bdestroy_wrapper(&mme_config.realm);
  bdestroy_wrapper(&mme_config.config_file);

  /*
   * IP configuration
   */
  bdestroy_wrapper(&mme_config.ip.if_name_s1_mme);
  bdestroy_wrapper(&mme_config.ip.if_name_s11);
  bdestroy_wrapper(&mme_config.s6a_config.conf_file);
  bdestroy_wrapper(&mme_config.itti_config.log_file);

  free_wrapper((void **) &mme_config.served_tai.plmn_mcc);
  free_wrapper((void **) &mme_config.served_tai.plmn_mnc);
  free_wrapper((void **) &mme_config.served_tai.plmn_mnc_len);
  free_wrapper((void **) &mme_config.served_tai.tac);

  for (int i = 0; i < mme_config.e_dns_emulation.nb_sgw_entries; i++) {
    bdestroy_wrapper(&mme_config.e_dns_emulation.sgw_id[i]);
  }
}

//------------------------------------------------------------------------------
int mme_config_parse_file(mme_config_t *config_pP)
{
  config_t cfg = {0};
  config_setting_t *setting_mme = NULL;
  config_setting_t *setting = NULL;
  config_setting_t *subsetting = NULL;
  config_setting_t *sub2setting = NULL;
  int aint = 0;
  int i = 0, n = 0, stop_index = 0, num = 0;
  const char *astring = NULL;
  const char *tac = NULL;
  const char *mcc = NULL;
  const char *mnc = NULL;
  char *if_name_s1_mme = NULL;
  char *s1_mme = NULL;
  char *if_name_s11 = NULL;
  char *s11 = NULL;
  #if EMBEDDED_SGW
  char *sgw_ip_address_for_s11 = NULL;
  #endif
  char *sgw_ip_address_for_s11 = NULL;
  bool swap = false;
  bstring address = NULL;
  bstring cidr = NULL;
  bstring mask = NULL;
  struct in_addr in_addr_var = {0};
  const char *csfb_mcc = NULL;
  const char *csfb_mnc = NULL;
  const char *lac = NULL;

  config_init(&cfg);

  if (config_pP->config_file != NULL) {
    /*
     * Read the file. If there is an error, report it and exit.
     */
    if (!config_read_file(&cfg, bdata(config_pP->config_file))) {
      OAILOG_CRITICAL(
        LOG_CONFIG,
        "Failed to parse MME configuration file: %s:%d - %s\n",
        bdata(config_pP->config_file),
        config_error_line(&cfg),
        config_error_text(&cfg));
      config_destroy(&cfg);
      AssertFatal(
        1 == 0,
        "Failed to parse MME configuration file %s!\n",
        bdata(config_pP->config_file));
    }
  } else {
    config_destroy(&cfg);
    AssertFatal(0, "No MME configuration file provided!\n");
  }

  setting_mme = config_lookup(&cfg, MME_CONFIG_STRING_MME_CONFIG);

  if (setting_mme != NULL) {
//OAILOG_DEBUG("reading mme config")    
         
    // LOGGING setting
    setting = config_setting_get_member(setting_mme, LOG_CONFIG_STRING_LOGGING);

    if (setting != NULL) {
      if (config_setting_lookup_string(
            setting, LOG_CONFIG_STRING_OUTPUT, (const char **) &astring)) {
        if (astring != NULL) {
          if (config_pP->log_config.output) {
            bassigncstr(config_pP->log_config.output, astring);
          } else {
            config_pP->log_config.output = bfromcstr(astring);
          }
        }
      }

      if (config_setting_lookup_string(
            setting,
            LOG_CONFIG_STRING_OUTPUT_THREAD_SAFE,
            (const char **) &astring)) {
        if (astring != NULL) {
          config_pP->log_config.is_output_thread_safe = parse_bool(astring);
        }
      }

      if (config_setting_lookup_string(
            setting, LOG_CONFIG_STRING_COLOR, (const char **) &astring)) {
        if (strcasecmp("yes", astring) == 0)
          config_pP->log_config.color = true;
        else
          config_pP->log_config.color = false;
      }

      if (config_setting_lookup_string(
            setting,
            LOG_CONFIG_STRING_SCTP_LOG_LEVEL,
            (const char **) &astring))
        config_pP->log_config.sctp_log_level = OAILOG_LEVEL_STR2INT(astring);

      if (config_setting_lookup_string(
            setting,
            LOG_CONFIG_STRING_S1AP_LOG_LEVEL,
            (const char **) &astring))
        config_pP->log_config.s1ap_log_level = OAILOG_LEVEL_STR2INT(astring);

      if (config_setting_lookup_string(
            setting, LOG_CONFIG_STRING_NAS_LOG_LEVEL, (const char **) &astring))
        config_pP->log_config.nas_log_level = OAILOG_LEVEL_STR2INT(astring);

      if (config_setting_lookup_string(
            setting,
            LOG_CONFIG_STRING_MME_APP_LOG_LEVEL,
            (const char **) &astring))
        config_pP->log_config.mme_app_log_level = OAILOG_LEVEL_STR2INT(astring);

      if (config_setting_lookup_string(
            setting, LOG_CONFIG_STRING_S6A_LOG_LEVEL, (const char **) &astring))
        config_pP->log_config.s6a_log_level = OAILOG_LEVEL_STR2INT(astring);
      if (config_setting_lookup_string(
            setting,
            LOG_CONFIG_STRING_SECU_LOG_LEVEL,
            (const char **) &astring))
        config_pP->log_config.secu_log_level = OAILOG_LEVEL_STR2INT(astring);

      if (config_setting_lookup_string(
            setting, LOG_CONFIG_STRING_UDP_LOG_LEVEL, (const char **) &astring))
        config_pP->log_config.udp_log_level = OAILOG_LEVEL_STR2INT(astring);

      if (config_setting_lookup_string(
            setting,
            LOG_CONFIG_STRING_UTIL_LOG_LEVEL,
            (const char **) &astring))
        config_pP->log_config.util_log_level = OAILOG_LEVEL_STR2INT(astring);

      if (config_setting_lookup_string(
            setting,
            LOG_CONFIG_STRING_ITTI_LOG_LEVEL,
            (const char **) &astring))
        config_pP->log_config.itti_log_level = OAILOG_LEVEL_STR2INT(astring);
#if EMBEDDED_SGW
      if (config_setting_lookup_string(
            setting,
            LOG_CONFIG_STRING_GTPV1U_LOG_LEVEL,
            (const char **) &astring))
        config_pP->log_config.gtpv1u_log_level = OAILOG_LEVEL_STR2INT(astring);

      if (config_setting_lookup_string(
            setting,
            LOG_CONFIG_STRING_SPGW_APP_LOG_LEVEL,
            (const char **) &astring))
        config_pP->log_config.spgw_app_log_level =
          OAILOG_LEVEL_STR2INT(astring);

      if (config_setting_lookup_string(
            setting,
            LOG_CONFIG_STRING_PGW_APP_LOG_LEVEL,
            (const char **) &astring))
        config_pP->log_config.pgw_app_log_level = OAILOG_LEVEL_STR2INT(astring);
#else
      if (config_setting_lookup_string(
            setting,
            LOG_CONFIG_STRING_GTPV2C_LOG_LEVEL,
            (const char **) &astring))
        config_pP->log_config.gtpv2c_log_level = OAILOG_LEVEL_STR2INT(astring);

      if (config_setting_lookup_string(
            setting, LOG_CONFIG_STRING_S11_LOG_LEVEL, (const char **) &astring))
        config_pP->log_config.s11_log_level = OAILOG_LEVEL_STR2INT(astring);
#endif
      if ((config_setting_lookup_string(
            setting_mme,
            MME_CONFIG_STRING_ASN1_VERBOSITY,
            (const char **) &astring))) {
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

    // GENERAL MME SETTINGS
    if ((config_setting_lookup_string(
          setting_mme, MME_CONFIG_STRING_REALM, (const char **) &astring))) {
      config_pP->realm = bfromcstr(astring);
    }

    if ((config_setting_lookup_string(
          setting_mme,
          MME_CONFIG_STRING_FULL_NETWORK_NAME,
          (const char **) &astring))) {
      config_pP->full_network_name = bfromcstr(astring);
    }

    if ((config_setting_lookup_string(
          setting_mme,
          MME_CONFIG_STRING_SHORT_NETWORK_NAME,
          (const char **) &astring))) {
      config_pP->short_network_name = bfromcstr(astring);
    }

    if ((config_setting_lookup_int(
          setting_mme, MME_CONFIG_STRING_DAYLIGHT_SAVING_TIME, &aint))) {
      config_pP->daylight_saving_time = (uint32_t) aint;
    }

    if ((config_setting_lookup_string(
          setting_mme,
          MME_CONFIG_STRING_PID_DIRECTORY,
          (const char **) &astring))) {
      config_pP->pid_dir = bfromcstr(astring);
    }

    if ((config_setting_lookup_int(
          setting_mme, MME_CONFIG_STRING_MAXENB, &aint))) {
      config_pP->max_enbs = (uint32_t) aint;
    }

    if ((config_setting_lookup_int(
          setting_mme, MME_CONFIG_STRING_MAXUE, &aint))) {
      config_pP->max_ues = (uint32_t) aint;
    }

    if ((config_setting_lookup_int(
          setting_mme, MME_CONFIG_STRING_RELATIVE_CAPACITY, &aint))) {
      config_pP->relative_capacity = (uint8_t) aint;
    }

    if ((config_setting_lookup_int(
          setting_mme, MME_CONFIG_STRING_STATISTIC_TIMER, &aint))) {
      config_pP->mme_statistic_timer = (uint32_t) aint;
    }

    if ((config_setting_lookup_string(
          setting_mme,
          MME_CONFIG_STRING_IP_CAPABILITY,
          (const char **) &astring))) {
      config_pP->ip_capability = bfromcstr(astring);
    }

    if ((config_setting_lookup_string(
          setting_mme,
          MME_CONFIG_STRING_USE_STATELESS,
          (const char **) &astring))) {
      config_pP->use_stateless = parse_bool(astring);
    }

    if ((config_setting_lookup_string(
          setting_mme,
          EPS_NETWORK_FEATURE_SUPPORT_EMERGENCY_BEARER_SERVICES_IN_S1_MODE,
          (const char **) &astring))) {
      config_pP->eps_network_feature_support
        .emergency_bearer_services_in_s1_mode = parse_bool(astring);
    }
    if ((config_setting_lookup_string(
          setting_mme,
          EPS_NETWORK_FEATURE_SUPPORT_EXTENDED_SERVICE_REQUEST,
          (const char **) &astring))) {
      config_pP->eps_network_feature_support.extended_service_request =
        parse_bool(astring);
    }
    if ((config_setting_lookup_string(
          setting_mme,
          EPS_NETWORK_FEATURE_SUPPORT_IMS_VOICE_OVER_PS_SESSION_IN_S1,
          (const char **) &astring))) {
      config_pP->eps_network_feature_support.ims_voice_over_ps_session_in_s1 =
        parse_bool(astring);
    }
    if ((config_setting_lookup_string(
          setting_mme,
          EPS_NETWORK_FEATURE_SUPPORT_LOCATION_SERVICES_VIA_EPC,
          (const char **) &astring))) {
      config_pP->eps_network_feature_support.location_services_via_epc =
        parse_bool(astring);
    }

    if ((config_setting_lookup_string(
          setting_mme,
          MME_CONFIG_STRING_UNAUTHENTICATED_IMSI_SUPPORTED,
          (const char **) &astring))) {
      config_pP->unauthenticated_imsi_supported = parse_bool(astring);
    }

    // ITTI SETTING
    setting = config_setting_get_member(
      setting_mme, MME_CONFIG_STRING_INTERTASK_INTERFACE_CONFIG);

    if (setting != NULL) {
      if ((config_setting_lookup_int(
            setting,
            MME_CONFIG_STRING_INTERTASK_INTERFACE_QUEUE_SIZE,
            &aint))) {
        config_pP->itti_config.queue_size = (uint32_t) aint;
      }
    }
#if !S6A_OVER_GRPC
    // S6A SETTING
    setting =
      config_setting_get_member(setting_mme, MME_CONFIG_STRING_S6A_CONFIG);

    if (setting != NULL) {
      if ((config_setting_lookup_string(
            setting,
            MME_CONFIG_STRING_S6A_CONF_FILE_PATH,
            (const char **) &astring))) {
        if (astring != NULL) {
          if (config_pP->s6a_config.conf_file) {
            bassigncstr(config_pP->s6a_config.conf_file, astring);
          } else {
            config_pP->s6a_config.conf_file = bfromcstr(astring);
          }
        }
      }

      if ((config_setting_lookup_string(
            setting,
            MME_CONFIG_STRING_S6A_HSS_HOSTNAME,
            (const char **) &astring))) {
        if (astring != NULL) {
          if (config_pP->s6a_config.hss_host_name) {
            bassigncstr(config_pP->s6a_config.hss_host_name, astring);
          } else {
            config_pP->s6a_config.hss_host_name = bfromcstr(astring);
          }
        } else
          AssertFatal(
            1 == 0,
            "You have to provide a valid HSS hostname %s=...\n",
            MME_CONFIG_STRING_S6A_HSS_HOSTNAME);
      }
    }
#endif /* !S6A_OVER_GRPC */
    // SCTP SETTING
    setting =
      config_setting_get_member(setting_mme, MME_CONFIG_STRING_SCTP_CONFIG);

    if (setting != NULL) {
      if ((config_setting_lookup_int(
            setting, MME_CONFIG_STRING_SCTP_INSTREAMS, &aint))) {
        config_pP->sctp_config.in_streams = (uint16_t) aint;
      }

      if ((config_setting_lookup_int(
            setting, MME_CONFIG_STRING_SCTP_OUTSTREAMS, &aint))) {
        config_pP->sctp_config.out_streams = (uint16_t) aint;
      }
    }
    // S1AP SETTING
    setting =
      config_setting_get_member(setting_mme, MME_CONFIG_STRING_S1AP_CONFIG);

    if (setting != NULL) {
      if ((config_setting_lookup_int(
            setting, MME_CONFIG_STRING_S1AP_OUTCOME_TIMER, &aint))) {
        config_pP->s1ap_config.outcome_drop_timer_sec = (uint8_t) aint;
      }

      if ((config_setting_lookup_int(
            setting, MME_CONFIG_STRING_S1AP_PORT, &aint))) {
        config_pP->s1ap_config.port_number = (uint16_t) aint;
      }
    }
    // TAI list setting
    setting =
      config_setting_get_member(setting_mme, MME_CONFIG_STRING_TAI_LIST);
    if (setting != NULL) {
      num = config_setting_length(setting);
      OAILOG_INFO(LOG_MME_APP, "Number of TAIs configured: %d\n", num);
      AssertFatal(
        num >= MIN_TAI_SUPPORTED,
        "Not even one TAI is configured, configure minimum one TAI\n");
      AssertFatal(
        num <= MAX_TAI_SUPPORTED,
        "Too many TAIs configured: %d (Maximum supported: %d)",
        num,
        MAX_TAI_SUPPORTED);

      if (config_pP->served_tai.nb_tai != num) {
        if (config_pP->served_tai.plmn_mcc != NULL)
          free_wrapper((void **) &config_pP->served_tai.plmn_mcc);

        if (config_pP->served_tai.plmn_mnc != NULL)
          free_wrapper((void **) &config_pP->served_tai.plmn_mnc);

        if (config_pP->served_tai.plmn_mnc_len != NULL)
          free_wrapper((void **) &config_pP->served_tai.plmn_mnc_len);

        if (config_pP->served_tai.tac != NULL)
          free_wrapper((void **) &config_pP->served_tai.tac);

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
            config_pP->served_tai.plmn_mnc[i] = (uint16_t) atoi(mnc);
            config_pP->served_tai.plmn_mnc_len[i] = strlen(mnc);

            AssertFatal(
              (config_pP->served_tai.plmn_mnc_len[i] == MIN_MNC_LENGTH) ||
                (config_pP->served_tai.plmn_mnc_len[i] == MAX_MNC_LENGTH),
              "Bad MNC length %u, must be %d or %d",
              config_pP->served_tai.plmn_mnc_len[i],
              MIN_MNC_LENGTH,
              MAX_MNC_LENGTH);
          }

          if ((config_setting_lookup_string(
                sub2setting, MME_CONFIG_STRING_TAC, &tac))) {
            config_pP->served_tai.tac[i] = (uint16_t) atoi(tac);

            AssertFatal(
              TAC_IS_VALID(config_pP->served_tai.tac[i]),
              "Invalid TAC value " TAC_FMT,
              config_pP->served_tai.tac[i]);
          }
        }
      }
      // sort TAI list
      n = config_pP->served_tai.nb_tai;
      do {
        stop_index = 0;
        for (i = 1; i < n; i++) {
          swap = false;
          if (
            config_pP->served_tai.plmn_mcc[i - 1] >
            config_pP->served_tai.plmn_mcc[i]) {
            swap = true;
          } else if (
            config_pP->served_tai.plmn_mcc[i - 1] ==
            config_pP->served_tai.plmn_mcc[i]) {
            if (
              config_pP->served_tai.plmn_mnc[i - 1] >
              config_pP->served_tai.plmn_mnc[i]) {
              swap = true;
            } else if (
              config_pP->served_tai.plmn_mnc[i - 1] ==
              config_pP->served_tai.plmn_mnc[i]) {
              if (
                config_pP->served_tai.tac[i - 1] >
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

            swap16 = config_pP->served_tai.tac[i - 1];
            config_pP->served_tai.tac[i - 1] = config_pP->served_tai.tac[i];
            config_pP->served_tai.tac[i] = swap16;

            stop_index = i;
          }
        }
        n = stop_index;
      } while (0 != n);
      // helper for determination of list type (global view), we could make sublists with different types, but keep things simple for now
      config_pP->served_tai.list_type =
        TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_CONSECUTIVE_TACS;
      for (i = 1; i < config_pP->served_tai.nb_tai; i++) {
        if (
          (config_pP->served_tai.plmn_mcc[i] !=
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
        if (
          config_pP->served_tai.tac[i] !=
          (config_pP->served_tai.tac[i - 1] + 1)) {
          config_pP->served_tai.list_type =
            TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_NON_CONSECUTIVE_TACS;
        }
      }
    }

    // GUMMEI SETTING
    setting =
      config_setting_get_member(setting_mme, MME_CONFIG_STRING_GUMMEI_LIST);
    config_pP->gummei.nb = 0;
    if (setting != NULL) {
      num = config_setting_length(setting);
      OAILOG_INFO(LOG_MME_APP, "Number of GUMMEIs configured =%d\n", num);
      AssertFatal(
        num >= MIN_GUMMEI,
        "Not even one GUMMEI is configured, configure minimum one GUMMEI \n");
      AssertFatal(
        num <= MAX_GUMMEI,
        "Number of GUMMEIs configured:%d exceeds number of GUMMEIs supported "
        ":%d \n",
        num,
        MAX_GUMMEI);

      for (i = 0; i < num; i++) {
        sub2setting = config_setting_get_elem(setting, i);

        if (sub2setting != NULL) {
          if ((config_setting_lookup_string(
                sub2setting, MME_CONFIG_STRING_MCC, &mcc))) {
            AssertFatal(
              strlen(mcc) == MAX_MCC_LENGTH,
              "Bad MCC length (%ld), it must be %u digit ex: 001",
              strlen(mcc),
              MAX_MCC_LENGTH);
            char c[2] = {mcc[0], 0};
            config_pP->gummei.gummei[i].plmn.mcc_digit1 = (uint8_t) atoi(c);
            c[0] = mcc[1];
            config_pP->gummei.gummei[i].plmn.mcc_digit2 = (uint8_t) atoi(c);
            c[0] = mcc[2];
            config_pP->gummei.gummei[i].plmn.mcc_digit3 = (uint8_t) atoi(c);
          }

          if ((config_setting_lookup_string(
                sub2setting, MME_CONFIG_STRING_MNC, &mnc))) {
            AssertFatal(
              (strlen(mnc) == MIN_MNC_LENGTH) ||
                (strlen(mnc) == MAX_MNC_LENGTH),
              "Bad MNC length (%ld), it must be %u or %u digit ex: 12 or 123",
              strlen(mnc),
              MIN_MNC_LENGTH,
              MAX_MNC_LENGTH);
            char c[2] = {mnc[0], 0};
            config_pP->gummei.gummei[i].plmn.mnc_digit1 = (uint8_t) atoi(c);
            c[0] = mnc[1];
            config_pP->gummei.gummei[i].plmn.mnc_digit2 = (uint8_t) atoi(c);
            if (strlen(mnc) == 3) {
              c[0] = mnc[2];
              config_pP->gummei.gummei[i].plmn.mnc_digit3 = (uint8_t) atoi(c);
            } else {
              config_pP->gummei.gummei[i].plmn.mnc_digit3 = 0x0F;
            }
          }

          if ((config_setting_lookup_string(
                sub2setting, MME_CONFIG_STRING_MME_GID, &mnc))) {
            config_pP->gummei.gummei[i].mme_gid = (uint16_t) atoi(mnc);
          }
          if ((config_setting_lookup_string(
                sub2setting, MME_CONFIG_STRING_MME_CODE, &mnc))) {
            config_pP->gummei.gummei[i].mme_code = (uint8_t) atoi(mnc);
          }
          config_pP->gummei.nb += 1;
        }
      }
    }
    // NETWORK INTERFACE SETTING
    setting = config_setting_get_member(
      setting_mme, MME_CONFIG_STRING_NETWORK_INTERFACES_CONFIG);

    if (setting != NULL) {
      if ((config_setting_lookup_string(
             setting,
             MME_CONFIG_STRING_INTERFACE_NAME_FOR_S1_MME,
             (const char **) &if_name_s1_mme) &&
           config_setting_lookup_string(
             setting,
             MME_CONFIG_STRING_IPV4_ADDRESS_FOR_S1_MME,
             (const char **) &s1_mme) &&
           config_setting_lookup_string(
             setting,
             MME_CONFIG_STRING_INTERFACE_NAME_FOR_S11_MME,
             (const char **) &if_name_s11) &&
           config_setting_lookup_string(
             setting,
             MME_CONFIG_STRING_IPV4_ADDRESS_FOR_S11_MME,
             (const char **) &s11) &&
           config_setting_lookup_int(
             setting, MME_CONFIG_STRING_MME_PORT_FOR_S11, &aint))) {
        config_pP->ip.port_s11 = (uint16_t) aint;

        config_pP->ip.if_name_s1_mme = bfromcstr(if_name_s1_mme);
        cidr = bfromcstr(s1_mme);
        struct bstrList *list = bsplit(cidr, '/');
        AssertFatal(
          list->qty == CIDR_SPLIT_LIST_COUNT,
          "Bad S1-MME CIDR address: %s",
          bdata(cidr));
        address = list->entry[0];
        mask = list->entry[1];
        IPV4_STR_ADDR_TO_INADDR(
          bdata(address),
          config_pP->ip.s1_mme_v4,
          "BAD IP ADDRESS FORMAT FOR S1-MME !\n");
        config_pP->ip.netmask_s1_mme = atoi((const char *) mask->data);
        bstrListDestroy(list);
        in_addr_var.s_addr = config_pP->ip.s1_mme_v4.s_addr;
        OAILOG_INFO(
          LOG_MME_APP,
          "Parsing configuration file found S1-MME: %s/%d on %s\n",
          inet_ntoa(in_addr_var),
          config_pP->ip.netmask_s1_mme,
          bdata(config_pP->ip.if_name_s1_mme));
        bdestroy_wrapper(&cidr);

        bdestroy(cidr);
        config_pP->ip.if_name_s11 = bfromcstr(if_name_s11);
        cidr = bfromcstr(s11);
        list = bsplit(cidr, '/');
        AssertFatal(
          list->qty == CIDR_SPLIT_LIST_COUNT,
          "Bad MME S11 CIDR address: %s",
          bdata(cidr));
        address = list->entry[0];
        mask = list->entry[1];
        IPV4_STR_ADDR_TO_INADDR(
          bdata(address),
          config_pP->ip.s11_mme_v4,
          "BAD IP ADDRESS FORMAT FOR S11 !\n");
        config_pP->ip.netmask_s11 = atoi((const char *) mask->data);
        bstrListDestroy(list);
        bdestroy_wrapper(&cidr);
        in_addr_var.s_addr = config_pP->ip.s11_mme_v4.s_addr;
        OAILOG_INFO(
          LOG_MME_APP,
          "Parsing configuration file found S11: %s/%d on %s\n",
          inet_ntoa(in_addr_var),
          config_pP->ip.netmask_s11,
          bdata(config_pP->ip.if_name_s11));
        bdestroy(cidr);
      }
    }

    // CSFB SETTING
    setting =
      config_setting_get_member(setting_mme, MME_CONFIG_STRING_CSFB_CONFIG);
    if (setting != NULL) {
      if ((config_setting_lookup_string(
            setting,
            MME_CONFIG_STRING_NON_EPS_SERVICE_CONTROL,
            (const char **) &astring))) {
        if (astring != NULL) {
          config_pP->non_eps_service_control = bfromcstr(astring);
        }
      }
      if (
        strcasecmp(
          (const char *) config_pP->non_eps_service_control->data, "OFF") !=
        0) {
        // Check CSFB MCC. MNC and LAC only if NON-EPS feature is enabled.
        if ((config_setting_lookup_string(
              setting, MME_CONFIG_STRING_CSFB_MCC, &csfb_mcc))) {
          AssertFatal(
            strlen(csfb_mcc) == MAX_MCC_LENGTH,
            "Bad MCC length(%ld), it must be %u digit ex: 001",
            strlen(csfb_mcc),
            MAX_MCC_LENGTH);
          char c[2] = {csfb_mcc[0], 0};
          config_pP->lai.mccdigit1 = (uint8_t) atoi(c);
          c[0] = csfb_mcc[1];
          config_pP->lai.mccdigit2 = (uint8_t) atoi(c);
          c[0] = csfb_mcc[2];
          config_pP->lai.mccdigit3 = (uint8_t) atoi(c);
        }
        if ((config_setting_lookup_string(
              setting, MME_CONFIG_STRING_CSFB_MNC, &csfb_mnc))) {
          AssertFatal(
            (strlen(csfb_mnc) == MIN_MNC_LENGTH) ||
              (strlen(csfb_mnc) == MAX_MNC_LENGTH),
            "Bad MNC length (%ld), it must be %u or %u digit ex: 12 or 123",
            strlen(csfb_mnc),
            MIN_MNC_LENGTH,
            MAX_MNC_LENGTH);
          char c[2] = {csfb_mnc[0], 0};
          config_pP->lai.mncdigit1 = (uint8_t) atoi(c);
          c[0] = csfb_mnc[1];
          config_pP->lai.mncdigit2 = (uint8_t) atoi(c);
          if (strlen(csfb_mnc) == 3) {
            c[0] = csfb_mnc[2];
            config_pP->lai.mncdigit3 = (uint8_t) atoi(c);
          } else {
            config_pP->lai.mncdigit3 = 0x0F;
          }
        }

        if ((config_setting_lookup_string(
              setting, MME_CONFIG_STRING_LAC, &lac))) {
          config_pP->lai.lac = (uint16_t) atoi(lac);
        }
      }
    }

    // NAS SETTING
    setting =
      config_setting_get_member(setting_mme, MME_CONFIG_STRING_NAS_CONFIG);

    if (setting != NULL) {
      subsetting = config_setting_get_member(
        setting, MME_CONFIG_STRING_NAS_SUPPORTED_INTEGRITY_ALGORITHM_LIST);

      if (subsetting != NULL) {
        num = config_setting_length(subsetting);

        if (num <= 8) {
          for (i = 0; i < num; i++) {
            astring = config_setting_get_string_elem(subsetting, i);

            if (strcmp("EIA0", astring) == 0)
              config_pP->nas_config.prefered_integrity_algorithm[i] =
                EIA0_ALG_ID;
            else if (strcmp("EIA1", astring) == 0)
              config_pP->nas_config.prefered_integrity_algorithm[i] =
                EIA1_128_ALG_ID;
            else if (strcmp("EIA2", astring) == 0)
              config_pP->nas_config.prefered_integrity_algorithm[i] =
                EIA2_128_ALG_ID;
            else
              config_pP->nas_config.prefered_integrity_algorithm[i] =
                EIA0_ALG_ID;
          }

          for (i = num; i < 8; i++) {
            config_pP->nas_config.prefered_integrity_algorithm[i] = EIA0_ALG_ID;
          }
        }
      }

      subsetting = config_setting_get_member(
        setting, MME_CONFIG_STRING_NAS_SUPPORTED_CIPHERING_ALGORITHM_LIST);

      if (subsetting != NULL) {
        num = config_setting_length(subsetting);

        if (num <= 8) {
          for (i = 0; i < num; i++) {
            astring = config_setting_get_string_elem(subsetting, i);

            if (strcmp("EEA0", astring) == 0)
              config_pP->nas_config.prefered_ciphering_algorithm[i] =
                EEA0_ALG_ID;
            else if (strcmp("EEA1", astring) == 0)
              config_pP->nas_config.prefered_ciphering_algorithm[i] =
                EEA1_128_ALG_ID;
            else if (strcmp("EEA2", astring) == 0)
              config_pP->nas_config.prefered_ciphering_algorithm[i] =
                EEA2_128_ALG_ID;
            else
              config_pP->nas_config.prefered_ciphering_algorithm[i] =
                EEA0_ALG_ID;
          }

          for (i = num; i < 8; i++) {
            config_pP->nas_config.prefered_ciphering_algorithm[i] = EEA0_ALG_ID;
          }
        }
      }
      if ((config_setting_lookup_int(
            setting, MME_CONFIG_STRING_NAS_T3402_TIMER, &aint))) {
        config_pP->nas_config.t3402_min = (uint32_t) aint;
      }
      if ((config_setting_lookup_int(
            setting, MME_CONFIG_STRING_NAS_T3412_TIMER, &aint))) {
        config_pP->nas_config.t3412_min = (uint32_t) aint;
      }
      if ((config_setting_lookup_int(
            setting, MME_CONFIG_STRING_NAS_T3422_TIMER, &aint))) {
        config_pP->nas_config.t3422_sec = (uint32_t) aint;
      }
      if ((config_setting_lookup_int(
            setting, MME_CONFIG_STRING_NAS_T3450_TIMER, &aint))) {
        config_pP->nas_config.t3450_sec = (uint32_t) aint;
      }
      if ((config_setting_lookup_int(
            setting, MME_CONFIG_STRING_NAS_T3460_TIMER, &aint))) {
        config_pP->nas_config.t3460_sec = (uint32_t) aint;
      }
      if ((config_setting_lookup_int(
            setting, MME_CONFIG_STRING_NAS_T3470_TIMER, &aint))) {
        config_pP->nas_config.t3470_sec = (uint32_t) aint;
      }
      if ((config_setting_lookup_int(
            setting, MME_CONFIG_STRING_NAS_T3485_TIMER, &aint))) {
        config_pP->nas_config.t3485_sec = (uint32_t) aint;
      }
      if ((config_setting_lookup_int(
            setting, MME_CONFIG_STRING_NAS_T3486_TIMER, &aint))) {
        config_pP->nas_config.t3486_sec = (uint32_t) aint;
      }
      if ((config_setting_lookup_int(
            setting, MME_CONFIG_STRING_NAS_T3489_TIMER, &aint))) {
        config_pP->nas_config.t3489_sec = (uint32_t) aint;
      }
      if ((config_setting_lookup_int(
            setting, MME_CONFIG_STRING_NAS_T3495_TIMER, &aint))) {
        config_pP->nas_config.t3495_sec = (uint32_t) aint;
      }
      if ((config_setting_lookup_string(
            setting,
            MME_CONFIG_STRING_NAS_FORCE_REJECT_TAU,
            (const char **) &astring))) {
        config_pP->nas_config.force_reject_tau = parse_bool(astring);
      }
      if ((config_setting_lookup_string(
            setting,
            MME_CONFIG_STRING_NAS_FORCE_REJECT_SR,
            (const char **) &astring))) {
        config_pP->nas_config.force_reject_sr = parse_bool(astring);
      }
      if ((config_setting_lookup_string(
            setting,
            MME_CONFIG_STRING_NAS_DISABLE_ESM_INFORMATION_PROCEDURE,
            (const char **) &astring))) {
        config_pP->nas_config.disable_esm_information = parse_bool(astring);
      }
    }

    //SGS TIMERS
    setting =
      config_setting_get_member(setting_mme, MME_CONFIG_STRING_SGS_CONFIG);

    if (setting != NULL) {
      if ((config_setting_lookup_int(
            setting, MME_CONFIG_STRING_SGS_TS6_1_TIMER, &aint))) {
        config_pP->sgs_config.ts6_1_sec = (uint8_t) aint;
      }
      if ((config_setting_lookup_int(
            setting, MME_CONFIG_STRING_SGS_TS8_TIMER, &aint))) {
        config_pP->sgs_config.ts8_sec = (uint8_t) aint;
      }
      if ((config_setting_lookup_int(
            setting, MME_CONFIG_STRING_SGS_TS9_TIMER, &aint))) {
        config_pP->sgs_config.ts9_sec = (uint8_t) aint;
      }
      if ((config_setting_lookup_int(
            setting, MME_CONFIG_STRING_SGS_TS10_TIMER, &aint))) {
        config_pP->sgs_config.ts10_sec = (uint8_t) aint;
      }
      if ((config_setting_lookup_int(
            setting, MME_CONFIG_STRING_SGS_TS13_TIMER, &aint))) {
        config_pP->sgs_config.ts13_sec = (uint8_t) aint;
      }
    }
  
  // S-GW Setting
  setting =  
     config_setting_get_member(setting_mme, MME_CONFIG_STRING_SGW_CONFIG);

     if (setting != NULL) {

     if ((config_setting_lookup_string(setting,SGW_CONFIG_STRING_SGW_IPV4_ADDRESS_FOR_S11,(const char **) &sgw_ip_address_for_s11))) 
                           {

          OAILOG_DEBUG (LOG_SPGW_APP, "sgw interface IP information %s\n", sgw_ip_address_for_s11);
          
          cidr = bfromcstr(sgw_ip_address_for_s11);
          struct bstrList *list = bsplit(cidr, '/');
          AssertFatal(2 == list->qty, "Bad CIDR address %s", bdata(cidr));
          address = list->entry[0];
          IPV4_STR_ADDR_TO_INADDR(
            bdata(address),
            config_pP->e_dns_emulation.sgw_ip_addr[0],
            "BAD IP ADDRESS FORMAT FOR SGW S11 !\n");

          bstrListDestroy(list);
          bdestroy_wrapper(&cidr);
          OAILOG_INFO(
            LOG_SPGW_APP,
            "Parsing configuration file found S-GW S11: %s\n",
            inet_ntoa(config_pP->e_dns_emulation.sgw_ip_addr[0]));
        }
     }
  
  }

/*
  num = config_setting_length(setting);

    AssertFatal(
      num <= MME_CONFIG_MAX_SGW,
      "Too many SGW entries (%d) defined (Maximum supported: %d)\n",
      num,
      MME_CONFIG_MAX_SGW);

    config_pP->e_dns_emulation.nb_sgw_entries = 0;

    for (i = 0; i < num; i++) {
      sub2setting = config_setting_get_elem(setting, i);
      if (sub2setting != NULL) {
        const char *id = NULL;
        if (!(config_setting_lookup_string(
              sub2setting, MME_CONFIG_STRING_ID, &id))) {
          OAILOG_ERROR(
            LOG_SPGW_APP,
            "Could not get SGW ID item %d in %s\n",
            i,
            SGW_CONFIG_STRING_SGW_IPV4_ADDRESS_FOR_S11);
          break;
        }
        config_pP->e_dns_emulation.sgw_id[i] = bfromcstr(id);


//#if EMBEDDED_SGW
        if ((config_setting_lookup_string(
              sub2setting,
              SGW_CONFIG_STRING_SGW_IPV4_ADDRESS_FOR_S11,
              (const char **) &sgw_ip_address_for_s11))) 
                           {
          OAILOG_DEBUG (LOG_SPGW_APP, "sgw interface IP information %s\n", sgw_ip_address_for_s11);
          
          cidr = bfromcstr(sgw_ip_address_for_s11);
          struct bstrList *list = bsplit(cidr, '/');
          AssertFatal(
            list->qty == CIDR_SPLIT_LIST_COUNT,
            "Bad SGW S11 CIDR address: %s",
            bdata(cidr));
          address = list->entry[0];
          IPV4_STR_ADDR_TO_INADDR(
            bdata(address),
            config_pP->e_dns_emulation.sgw_ip_addr[i],
            "BAD IP ADDRESS FORMAT FOR SGW S11 !\n");


          bstrListDestroy(list)
          bdestroy_wrapper(&cidr);
          OAILOG_INFO(
            LOG_SPGW_APP,
            "Parsing configuration file found S-GW S11: %s\n",
            inet_ntoa(config_pP->e_dns_emulation.sgw_ip_addr[i]));
        }
//#endif
      }
      config_pP->e_dns_emulation.nb_sgw_entries++;
    }
 */
 //}

  config_destroy(&cfg);
  return 0;
}

//------------------------------------------------------------------------------
void mme_config_display(mme_config_t *config_pP)
{
  int j;

  OAILOG_INFO(
    LOG_CONFIG, "==== EURECOM %s v%s ====\n", PACKAGE_NAME, PACKAGE_VERSION);
  OAILOG_DEBUG(
    LOG_CONFIG,
    "Built with EMBEDDED_SGW .................: %d\n", EMBEDDED_SGW);
  OAILOG_DEBUG(
    LOG_CONFIG,
    "Built with S6A_OVER_GRPC .....................: %d\n", S6A_OVER_GRPC);

#if DEBUG_IS_ON
  OAILOG_DEBUG(
    LOG_CONFIG,
    "Built with CMAKE_BUILD_TYPE ................: %s\n",
    CMAKE_BUILD_TYPE);
  OAILOG_DEBUG(
    LOG_CONFIG,
    "Built with PACKAGE_NAME ....................: %s\n",
    PACKAGE_NAME);
  OAILOG_DEBUG(
    LOG_CONFIG,
    "Built with S1AP_DEBUG_LIST .................: %d\n",
    S1AP_DEBUG_LIST);
  OAILOG_DEBUG(
    LOG_CONFIG,
    "Built with SCTP_DUMP_LIST ..................: %d\n",
    SCTP_DUMP_LIST);
  OAILOG_DEBUG(
    LOG_CONFIG,
    "Built with TRACE_HASHTABLE .................: %d\n",
    TRACE_HASHTABLE);
  OAILOG_DEBUG(
    LOG_CONFIG,
    "Built with TRACE_3GPP_SPEC .................: %d\n",
    TRACE_3GPP_SPEC);
#endif
  OAILOG_INFO(LOG_CONFIG, "Configuration:\n");
  OAILOG_INFO(
    LOG_CONFIG,
    "- File .................................: %s\n",
    bdata(config_pP->config_file));
  OAILOG_INFO(
    LOG_CONFIG,
    "- Realm ................................: %s\n",
    bdata(config_pP->realm));
  OAILOG_INFO(
    LOG_CONFIG,
    "  full network name ....................: %s\n",
    bdata(config_pP->full_network_name));
  OAILOG_INFO(
    LOG_CONFIG,
    "  short network name ...................: %s\n",
    bdata(config_pP->short_network_name));
  OAILOG_INFO(
    LOG_CONFIG,
    "  Daylight Saving Time..................: %d\n",
    config_pP->daylight_saving_time);
  OAILOG_INFO(
    LOG_CONFIG,
    "- Run mode .............................: %s\n",
    (RUN_MODE_TEST == config_pP->run_mode) ? "TEST" : "NORMAL");
  OAILOG_INFO(
    LOG_CONFIG,
    "- Max eNBs .............................: %u\n",
    config_pP->max_enbs);
  OAILOG_INFO(
    LOG_CONFIG,
    "- Max UEs ..............................: %u\n",
    config_pP->max_ues);
  OAILOG_INFO(
    LOG_CONFIG,
    "- IMS voice over PS session in S1 ......: %s\n",
    config_pP->eps_network_feature_support.ims_voice_over_ps_session_in_s1 ==
        0 ?
      "false" :
      "true");
  OAILOG_INFO(
    LOG_CONFIG,
    "- Emergency bearer services in S1 mode .: %s\n",
    config_pP->eps_network_feature_support
          .emergency_bearer_services_in_s1_mode == 0 ?
      "false" :
      "true");
  OAILOG_INFO(
    LOG_CONFIG,
    "- Location services via epc ............: %s\n",
    config_pP->eps_network_feature_support.location_services_via_epc == 0 ?
      "false" :
      "true");
  OAILOG_INFO(
    LOG_CONFIG,
    "- Extended service request .............: %s\n",
    config_pP->eps_network_feature_support.extended_service_request == 0 ?
      "false" :
      "true");
  OAILOG_INFO(
    LOG_CONFIG,
    "- Unauth IMSI support ..................: %s\n",
    config_pP->unauthenticated_imsi_supported == 0 ? "false" : "true");
  OAILOG_INFO(
    LOG_CONFIG,
    "- Relative capa ........................: %u\n",
    config_pP->relative_capacity);
  OAILOG_INFO(
    LOG_CONFIG,
    "- Statistics timer .....................: %u (seconds)\n\n",
    config_pP->mme_statistic_timer);
  OAILOG_INFO(
    LOG_CONFIG,
    "- IP Capability ........................: %s\n\n",
    bdata(config_pP->ip_capability));
  OAILOG_INFO(
    LOG_CONFIG,
    "- Use Stateless ........................: %s\n\n",
    config_pP->use_stateless ? "true" : "false");
  OAILOG_INFO(LOG_CONFIG, "- CSFB:\n");
  OAILOG_INFO(
    LOG_CONFIG,
    "    Non EPS Service Control ........................: %s\n\n",
    bdata(config_pP->non_eps_service_control));
  OAILOG_INFO(LOG_CONFIG, "- S1-MME:\n");
  OAILOG_INFO(
    LOG_CONFIG,
    "    port number ......: %d\n",
    config_pP->s1ap_config.port_number);
  OAILOG_INFO(LOG_CONFIG, "- IP:\n");
  OAILOG_INFO(
    LOG_CONFIG,
    "    s1-MME iface .....: %s\n",
    bdata(config_pP->ip.if_name_s1_mme));
  OAILOG_INFO(
    LOG_CONFIG,
    "    s1-MME ip ........: %s\n",
    inet_ntoa(*((struct in_addr *) &config_pP->ip.s1_mme_v4)));
  OAILOG_INFO(
    LOG_CONFIG,
    "    s11 MME iface ....: %s\n",
    bdata(config_pP->ip.if_name_s11));
  OAILOG_INFO(
    LOG_CONFIG, "    s11 MME port .....: %d\n", config_pP->ip.port_s11);
  OAILOG_INFO(
    LOG_CONFIG,
    "    s11 MME ip .......: %s\n",
    inet_ntoa(*((struct in_addr *) &config_pP->ip.s11_mme_v4)));
  OAILOG_INFO(LOG_CONFIG, "- ITTI:\n");
  OAILOG_INFO(
    LOG_CONFIG,
    "    queue size .......: %u (bytes)\n",
    config_pP->itti_config.queue_size);
  OAILOG_INFO(
    LOG_CONFIG,
    "    log file .........: %s\n",
    bdata(config_pP->itti_config.log_file));
  OAILOG_INFO(LOG_CONFIG, "- SCTP:\n");
  OAILOG_INFO(
    LOG_CONFIG,
    "    in streams .......: %u\n",
    config_pP->sctp_config.in_streams);
  OAILOG_INFO(
    LOG_CONFIG,
    "    out streams ......: %u\n",
    config_pP->sctp_config.out_streams);
  OAILOG_INFO(LOG_CONFIG, "- GUMMEIs (PLMN|MMEGI|MMEC):\n");
  for (j = 0; j < config_pP->gummei.nb; j++) {
    OAILOG_INFO(
      LOG_CONFIG,
      "            " PLMN_FMT "|%u|%u \n",
      PLMN_ARG(&config_pP->gummei.gummei[j].plmn),
      config_pP->gummei.gummei[j].mme_gid,
      config_pP->gummei.gummei[j].mme_code);
  }
  OAILOG_INFO(LOG_CONFIG, "- TAIs : (mcc.mnc:tac)\n");
  switch (config_pP->served_tai.list_type) {
    case TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_CONSECUTIVE_TACS:
      OAILOG_INFO(LOG_CONFIG, "- TAI list type one PLMN consecutive TACs\n");
      break;
    case TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_NON_CONSECUTIVE_TACS:
      OAILOG_INFO(
        LOG_CONFIG, "- TAI list type one PLMN non consecutive TACs\n");
      break;
    case TRACKING_AREA_IDENTITY_LIST_TYPE_MANY_PLMNS:
      OAILOG_INFO(LOG_CONFIG, "- TAI list type multiple PLMNs\n");
      break;
    default:
      Fatal(
        "Invalid served TAI list type (%u) configured\n",
        config_pP->served_tai.list_type);
      break;
  }
  for (j = 0; j < config_pP->served_tai.nb_tai; j++) {
    if (config_pP->served_tai.plmn_mnc_len[j] == 2) {
      OAILOG_INFO(
        LOG_CONFIG,
        "            %3u.%3u:%u\n",
        config_pP->served_tai.plmn_mcc[j],
        config_pP->served_tai.plmn_mnc[j],
        config_pP->served_tai.tac[j]);
    } else {
      OAILOG_INFO(
        LOG_CONFIG,
        "            %3u.%03u:%u\n",
        config_pP->served_tai.plmn_mcc[j],
        config_pP->served_tai.plmn_mnc[j],
        config_pP->served_tai.tac[j]);
    }
  }
  OAILOG_INFO(LOG_CONFIG, "- NAS:\n");
  OAILOG_INFO(
    LOG_CONFIG,
    "    Preferred Integrity Algorithms .: EIA%d EIA%d EIA%d EIA%d (decreasing "
    "priority)\n",
    config_pP->nas_config.prefered_integrity_algorithm[0],
    config_pP->nas_config.prefered_integrity_algorithm[1],
    config_pP->nas_config.prefered_integrity_algorithm[2],
    config_pP->nas_config.prefered_integrity_algorithm[3]);
  OAILOG_INFO(
    LOG_CONFIG,
    "    Preferred Integrity Algorithms .: EEA%d EEA%d EEA%d EEA%d (decreasing "
    "priority)\n",
    config_pP->nas_config.prefered_ciphering_algorithm[0],
    config_pP->nas_config.prefered_ciphering_algorithm[1],
    config_pP->nas_config.prefered_ciphering_algorithm[2],
    config_pP->nas_config.prefered_ciphering_algorithm[3]);
  OAILOG_INFO(
    LOG_CONFIG, "    T3402 ....: %d min\n", config_pP->nas_config.t3402_min);
  OAILOG_INFO(
    LOG_CONFIG, "    T3412 ....: %d min\n", config_pP->nas_config.t3412_min);
  OAILOG_INFO(
    LOG_CONFIG, "    T3422 ....: %d sec\n", config_pP->nas_config.t3422_sec);
  OAILOG_INFO(
    LOG_CONFIG, "    T3450 ....: %d sec\n", config_pP->nas_config.t3450_sec);
  OAILOG_INFO(
    LOG_CONFIG, "    T3460 ....: %d sec\n", config_pP->nas_config.t3460_sec);
  OAILOG_INFO(
    LOG_CONFIG, "    T3470 ....: %d sec\n", config_pP->nas_config.t3470_sec);
  OAILOG_INFO(
    LOG_CONFIG, "    T3485 ....: %d sec\n", config_pP->nas_config.t3485_sec);
  OAILOG_INFO(
    LOG_CONFIG, "    T3486 ....: %d sec\n", config_pP->nas_config.t3486_sec);
  OAILOG_INFO(
    LOG_CONFIG, "    T3489 ....: %d sec\n", config_pP->nas_config.t3489_sec);
  OAILOG_INFO(
    LOG_CONFIG, "    T3470 ....: %d sec\n", config_pP->nas_config.t3470_sec);
  OAILOG_INFO(
    LOG_CONFIG, "    T3495 ....: %d sec\n", config_pP->nas_config.t3495_sec);
  OAILOG_INFO(LOG_CONFIG, "    NAS non standard features .:\n");
  OAILOG_INFO(
    LOG_CONFIG,
    "      Force reject TAU ............: %s\n",
    (config_pP->nas_config.force_reject_tau) ? "true" : "false");
  OAILOG_INFO(
    LOG_CONFIG,
    "      Force reject SR .............: %s\n",
    (config_pP->nas_config.force_reject_sr) ? "true" : "false");
  OAILOG_INFO(
    LOG_CONFIG,
    "      Disable Esm information .....: %s\n",
    (config_pP->nas_config.disable_esm_information) ? "true" : "false");

  OAILOG_INFO(LOG_CONFIG, "- S6A:\n");
#if S6A_OVER_GRPC
  OAILOG_INFO(
    LOG_CONFIG, "    protocol .........: gRPC\n");
#else
  OAILOG_INFO(
    LOG_CONFIG, "    protocol .........: diameter\n");
  OAILOG_INFO(
    LOG_CONFIG,
    "    conf file ........: %s\n",
    bdata(config_pP->s6a_config.conf_file));
#endif
  OAILOG_INFO(LOG_CONFIG, "- Service303:\n");
  OAILOG_INFO(
    LOG_CONFIG,
    "    service name ........: %s\n",
    bdata(config_pP->service303_config.name));
  OAILOG_INFO(
    LOG_CONFIG,
    "    version ........: %s\n",
    bdata(config_pP->service303_config.version));
  OAILOG_INFO(LOG_CONFIG, "- Logging:\n");
  OAILOG_INFO(
    LOG_CONFIG,
    "    Output ..............: %s\n",
    bdata(config_pP->log_config.output));
  OAILOG_INFO(
    LOG_CONFIG,
    "    Output thread safe ..: %s\n",
    (config_pP->log_config.is_output_thread_safe) ? "true" : "false");
  OAILOG_INFO(
    LOG_CONFIG,
    "    Output with color ...: %s\n",
    (config_pP->log_config.color) ? "true" : "false");
  OAILOG_INFO(
    LOG_CONFIG,
    "    UDP log level........: %s\n",
    OAILOG_LEVEL_INT2STR(config_pP->log_config.udp_log_level));
  OAILOG_INFO(
    LOG_CONFIG,
    "    GTPV1-U log level....: %s\n",
    OAILOG_LEVEL_INT2STR(config_pP->log_config.gtpv1u_log_level));
  OAILOG_INFO(
    LOG_CONFIG,
    "    GTPV2-C log level....: %s\n",
    OAILOG_LEVEL_INT2STR(config_pP->log_config.gtpv2c_log_level));
  OAILOG_INFO(
    LOG_CONFIG,
    "    SCTP log level.......: %s\n",
    OAILOG_LEVEL_INT2STR(config_pP->log_config.sctp_log_level));
  OAILOG_INFO(
    LOG_CONFIG,
    "    S1AP log level.......: %s\n",
    OAILOG_LEVEL_INT2STR(config_pP->log_config.s1ap_log_level));
  OAILOG_INFO(
    LOG_CONFIG,
    "    ASN1 Verbosity level : %d\n",
    config_pP->log_config.asn1_verbosity_level);
  OAILOG_INFO(
    LOG_CONFIG,
    "    NAS log level........: %s\n",
    OAILOG_LEVEL_INT2STR(config_pP->log_config.nas_log_level));
  OAILOG_INFO(
    LOG_CONFIG,
    "    MME_APP log level....: %s\n",
    OAILOG_LEVEL_INT2STR(config_pP->log_config.mme_app_log_level));
  OAILOG_INFO(
    LOG_CONFIG,
    "    SPGW_APP log level....: %s\n",
    OAILOG_LEVEL_INT2STR(config_pP->log_config.spgw_app_log_level));
  OAILOG_INFO(
    LOG_CONFIG,
    "    PGW_APP log level....: %s\n",
    OAILOG_LEVEL_INT2STR(config_pP->log_config.pgw_app_log_level));
  OAILOG_INFO(
    LOG_CONFIG,
    "    S11 log level........: %s\n",
    OAILOG_LEVEL_INT2STR(config_pP->log_config.s11_log_level));
  OAILOG_INFO(
    LOG_CONFIG,
    "    S6a log level........: %s\n",
    OAILOG_LEVEL_INT2STR(config_pP->log_config.s6a_log_level));
  OAILOG_INFO(
    LOG_CONFIG,
    "    UTIL log level.......: %s\n",
    OAILOG_LEVEL_INT2STR(config_pP->log_config.util_log_level));
  OAILOG_INFO(
    LOG_CONFIG,
    "    ITTI log level.......: %s (InTer-Task Interface)\n",
    OAILOG_LEVEL_INT2STR(config_pP->log_config.itti_log_level));
}

//------------------------------------------------------------------------------
static void usage(char *target)
{
  OAI_FPRINTF_INFO(
    "==== EURECOM %s version: %s ====\n", PACKAGE_NAME, PACKAGE_VERSION);
  OAI_FPRINTF_INFO("Please report any bug to: %s\n", PACKAGE_BUGREPORT);
  OAI_FPRINTF_INFO("Usage: %s [options]\n", target);
  OAI_FPRINTF_INFO("Available options:\n");
  OAI_FPRINTF_INFO("-h      Print this help and return\n");
  OAI_FPRINTF_INFO("-c<path>\n");
  OAI_FPRINTF_INFO("        Set the configuration file for mme\n");
  OAI_FPRINTF_INFO("        See template in UTILS/CONF\n");
  OAI_FPRINTF_INFO("-V      Print %s version and return\n", PACKAGE_NAME);
  OAI_FPRINTF_INFO("-v[1-2] Debug level:\n");
  OAI_FPRINTF_INFO("            1 -> ASN1 XER printf on and ASN1 debug off\n");
  OAI_FPRINTF_INFO("            2 -> ASN1 XER printf on and ASN1 debug on\n");
}

//------------------------------------------------------------------------------
int mme_config_parse_opt_line(int argc, char *argv[], mme_config_t *config_pP)
{
  int c;

  mme_config_init(config_pP);

  /*
   * Parsing command line
   */
  while ((c = getopt(argc, argv, "c:s:h:v:V")) != -1) {
    switch (c) {
      case 'c': {
        /*
         * Store the given configuration file. If no file is given,
         * * * * then the default values will be used.
         */
        config_pP->config_file = blk2bstr(optarg, strlen(optarg));
        OAI_FPRINTF_INFO(
          "%s mme_config.config_file %s\n",
          __FUNCTION__,
          bdata(config_pP->config_file));
      } break;

      case 'v': {
        config_pP->log_config.asn1_verbosity_level = atoi(optarg);
      } break;

      case 'V': {
        OAI_FPRINTF_INFO(
          "==== EURECOM %s v%s ===="
          "Please report any bug to: %s\n",
          PACKAGE_NAME,
          PACKAGE_VERSION,
          PACKAGE_BUGREPORT);
      } break;
      
      case 's': {
       OAI_FPRINTF_INFO(
           "Ignoring command line option s as there is no embedded sgw \n");
      } break;

      case 'h': /* Fall through */
    
      default:
        OAI_FPRINTF_ERR("Unknown command line option %c\n", c);
        usage(argv[0]);
        exit(0);
            
    }
  }

  /*
   * Parse the configuration file using libconfig
   */
  if (!config_pP->config_file) {
    config_pP->config_file = bfromcstr("/usr/local/etc/oai/mme.conf");
  }
  if (mme_config_parse_file(config_pP) != 0) {
    return -1;
  }

  /*
   * Display the configuration
   */
  mme_config_display(config_pP);
  return 0;
}

static bool parse_bool(const char *str)
{
  if (strcasecmp(str, "yes") == 0) return true;
  if (strcasecmp(str, "true") == 0) return true;
  if (strcasecmp(str, "no") == 0) return false;
  if (strcasecmp(str, "false") == 0) return false;
  if (strcasecmp(str, "") == 0) return false;

  Fatal("Error in config file: got \"%s\" but expected bool\n", str);
}
