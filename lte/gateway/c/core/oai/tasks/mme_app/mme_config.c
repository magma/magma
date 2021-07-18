/*
 * Copyright (c) 2015, EURECOM (www.eurecom.fr)
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice,
 * this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 *
 * The views and conclusions contained in the software and documentation are
 * those of the authors and should not be interpreted as representing official
 * policies, either expressed or implied, of the FreeBSD Project.
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
#include "conversions.h"
#if EMBEDDED_SGW
#include "sgw_config.h"
#endif
static bool parse_bool(const char* str);

struct mme_config_s mme_config = {.rw_lock = PTHREAD_RWLOCK_INITIALIZER, 0};

//------------------------------------------------------------------------------
int mme_config_find_mnc_length(
    const char mcc_digit1P, const char mcc_digit2P, const char mcc_digit3P,
    const char mnc_digit1P, const char mnc_digit2P, const char mnc_digit3P) {
  uint16_t mcc   = 100 * mcc_digit1P + 10 * mcc_digit2P + mcc_digit3P;
  uint16_t mnc3  = 100 * mnc_digit1P + 10 * mnc_digit2P + mnc_digit3P;
  uint16_t mnc2  = 10 * mnc_digit1P + mnc_digit2P;
  int plmn_index = 0;

  if (mcc_digit1P < 0 || mcc_digit1P > 9 || mcc_digit2P < 0 ||
      mcc_digit2P > 9 || mcc_digit3P < 0 || mcc_digit3P > 9) {
    OAILOG_ERROR(
        LOG_MME_APP, "BAD MCC PARAMETER (%d%d%d)!\n", mcc_digit1P, mcc_digit2P,
        mcc_digit3P);
    return 0;
  }
  if (mnc_digit2P < 0 || mnc_digit2P > 9 || mnc_digit1P < 0 ||
      mnc_digit1P > 9) {
    OAILOG_ERROR(
        LOG_MME_APP, "BAD MNC PARAMETER (%d%d%d)!\n", mnc_digit1P, mnc_digit2P,
        mnc_digit3P);
    return 0;
  }

  while (plmn_index < mme_config.served_tai.nb_tai) {
    if (mme_config.served_tai.plmn_mcc[plmn_index] == mcc) {
      if ((mme_config.served_tai.plmn_mnc[plmn_index] == mnc2) &&
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

void log_config_init(log_config_t* log_conf) {
  memset(log_conf, 0, sizeof(*log_conf));

  log_conf->output                = NULL;
  log_conf->is_output_thread_safe = false;
  log_conf->color                 = false;

  log_conf->udp_log_level        = MAX_LOG_LEVEL;  // Means invalid TODO wtf
  log_conf->gtpv1u_log_level     = MAX_LOG_LEVEL;
  log_conf->gtpv2c_log_level     = MAX_LOG_LEVEL;
  log_conf->sctp_log_level       = MAX_LOG_LEVEL;
  log_conf->s1ap_log_level       = MAX_LOG_LEVEL;
  log_conf->nas_log_level        = MAX_LOG_LEVEL;
  log_conf->mme_app_log_level    = MAX_LOG_LEVEL;
  log_conf->s11_log_level        = MAX_LOG_LEVEL;
  log_conf->s6a_log_level        = MAX_LOG_LEVEL;
  log_conf->secu_log_level       = MAX_LOG_LEVEL;
  log_conf->util_log_level       = MAX_LOG_LEVEL;
  log_conf->itti_log_level       = MAX_LOG_LEVEL;
  log_conf->spgw_app_log_level   = MAX_LOG_LEVEL;
  log_conf->asn1_verbosity_level = 0;
}

void eps_network_feature_config_init(eps_network_feature_config_t* eps_conf) {
  eps_conf->emergency_bearer_services_in_s1_mode = 0;
  eps_conf->extended_service_request             = 0;
  eps_conf->ims_voice_over_ps_session_in_s1      = 0;
  eps_conf->location_services_via_epc            = 0;
}

void ipv4_config_init(ip_t* ip) {
  memset(ip, 0, sizeof(*ip));

  ip->if_name_s1_mme   = NULL;
  ip->s1_mme_v4.s_addr = INADDR_ANY;

  ip->if_name_s11       = NULL;
  ip->s11_mme_v4.s_addr = INADDR_ANY;

  ip->port_s11 = 2123;
}

void s1ap_config_init(s1ap_config_t* s1ap_conf) {
  s1ap_conf->port_number            = S1AP_PORT_NUMBER;
  s1ap_conf->outcome_drop_timer_sec = S1AP_OUTCOME_TIMER_DEFAULT;
}

void s6a_config_init(s6a_config_t* s6a_conf) {
  s6a_conf->hss_host_name = NULL;
  s6a_conf->hss_realm     = NULL;
  s6a_conf->conf_file     = bfromcstr(S6A_CONF_FILE);
}

void itti_config_init(itti_config_t* itti_conf) {
  itti_conf->queue_size = ITTI_QUEUE_MAX_ELEMENTS;
  itti_conf->log_file   = NULL;
}

void sctp_config_init(sctp_config_t* sctp_conf) {
  sctp_conf->in_streams  = SCTP_IN_STREAMS;
  sctp_conf->out_streams = SCTP_OUT_STREAMS;
}

void apn_map_config_init(apn_map_config_t* apn_map_config) {
  apn_map_config->nb                      = 0;
  apn_map_config->apn_map[0].imsi_prefix  = NULL;
  apn_map_config->apn_map[0].apn_override = NULL;
}

void nas_config_init(nas_config_t* nas_conf) {
  nas_conf->t3402_min               = T3402_DEFAULT_VALUE;
  nas_conf->t3412_min               = T3412_DEFAULT_VALUE;
  nas_conf->t3422_sec               = T3422_DEFAULT_VALUE;
  nas_conf->t3450_sec               = T3450_DEFAULT_VALUE;
  nas_conf->t3460_sec               = T3460_DEFAULT_VALUE;
  nas_conf->t3470_sec               = T3470_DEFAULT_VALUE;
  nas_conf->t3485_sec               = T3485_DEFAULT_VALUE;
  nas_conf->t3486_sec               = T3486_DEFAULT_VALUE;
  nas_conf->t3489_sec               = T3489_DEFAULT_VALUE;
  nas_conf->t3495_sec               = T3495_DEFAULT_VALUE;
  nas_conf->force_reject_tau        = true;
  nas_conf->force_reject_sr         = true;
  nas_conf->disable_esm_information = false;
  nas_conf->enable_apn_correction   = false;
  apn_map_config_init(&nas_conf->apn_map_config);
}

void gummei_config_init(gummei_config_t* gummei_conf) {
  gummei_conf->nb                        = 1;
  gummei_conf->gummei[0].mme_code        = MMEC;
  gummei_conf->gummei[0].mme_gid         = MMEGID;
  gummei_conf->gummei[0].plmn.mcc_digit1 = 0;
  gummei_conf->gummei[0].plmn.mcc_digit2 = 0;
  gummei_conf->gummei[0].plmn.mcc_digit3 = 1;
  gummei_conf->gummei[0].plmn.mcc_digit1 = 0;
  gummei_conf->gummei[0].plmn.mcc_digit2 = 1;
  gummei_conf->gummei[0].plmn.mcc_digit3 = 0x0F;
}

void served_tai_config_init(served_tai_t* served_tai) {
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

void service303_config_init(service303_data_t* service303_conf) {
  service303_conf->name    = bfromcstr(SERVICE303_MME_PACKAGE_NAME);
  service303_conf->version = bfromcstr(SERVICE303_MME_PACKAGE_VERSION);
}

void blocked_imei_config_init(blocked_imei_list_t* blocked_imeis) {
  blocked_imeis->num       = 0;
  blocked_imeis->imei_htbl = NULL;
}

void sac_to_tacs_map_config_init(sac_to_tacs_map_config_t* sac_to_tacs_map) {
  sac_to_tacs_map->sac_to_tacs_map_htbl = NULL;
}

//------------------------------------------------------------------------------
void mme_config_init(mme_config_t* config) {
  memset(config, 0, sizeof(*config));

  pthread_rwlock_init(&config->rw_lock, NULL);

  config->config_file                               = NULL;
  config->max_enbs                                  = 2;
  config->max_ues                                   = 2;
  config->unauthenticated_imsi_supported            = 0;
  config->relative_capacity                         = RELATIVE_CAPACITY;
  config->stats_timer_sec                           = 60;
  config->service303_config.stats_display_timer_sec = 60;
  config->enable_congestion_control                 = true;
  config->s1ap_zmq_th                               = LONG_MAX;
  config->mme_app_zmq_congest_th                    = LONG_MAX;
  config->mme_app_zmq_auth_th                       = LONG_MAX;
  config->mme_app_zmq_ident_th                      = LONG_MAX;
  config->mme_app_zmq_smc_th                        = LONG_MAX;

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
  blocked_imei_config_init(&config->blocked_imei);
  sac_to_tacs_map_config_init(&config->sac_to_tacs_map);
}

//------------------------------------------------------------------------------
void mme_config_exit(void) {
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

  free_wrapper((void**) &mme_config.served_tai.plmn_mcc);
  free_wrapper((void**) &mme_config.served_tai.plmn_mnc);
  free_wrapper((void**) &mme_config.served_tai.plmn_mnc_len);
  free_wrapper((void**) &mme_config.served_tai.tac);

  for (int i = 0; i < mme_config.e_dns_emulation.nb_sgw_entries; i++) {
    bdestroy_wrapper(&mme_config.e_dns_emulation.sgw_id[i]);
  }

  if (mme_config.blocked_imei.imei_htbl) {
    hashtable_uint64_ts_destroy(mme_config.blocked_imei.imei_htbl);
  }

  if (mme_config.sac_to_tacs_map.sac_to_tacs_map_htbl) {
    obj_hashtable_destroy(mme_config.sac_to_tacs_map.sac_to_tacs_map_htbl);
  }
}

//------------------------------------------------------------------------------
int mme_config_parse_file(mme_config_t* config_pP) {
  config_t cfg                  = {0};
  config_setting_t* setting_mme = NULL;
  config_setting_t* setting     = NULL;
  config_setting_t* subsetting  = NULL;
  config_setting_t* sub2setting = NULL;
  config_setting_t* sub3setting = NULL;
  int aint                      = 0;
  double adouble                = 0.0;
  int i = 0, n = 0, stop_index = 0, num = 0;
  const char* astring  = NULL;
  const char* tac      = NULL;
  const char* mcc      = NULL;
  const char* mnc      = NULL;
  char* if_name_s1_mme = NULL;
  char* s1_mme         = NULL;
  char* if_name_s11    = NULL;
  char* s11            = NULL;
  char* imsi_low_tmp   = NULL;
  char* imsi_high_tmp  = NULL;
#if !EMBEDDED_SGW
  char* sgw_ip_address_for_s11 = NULL;
#endif
  bool swap                  = false;
  bstring address            = NULL;
  bstring cidr               = NULL;
  bstring mask               = NULL;
  const char* imsi_prefix    = NULL;
  const char* apn_override   = NULL;
  struct in_addr in_addr_var = {0};
  const char* csfb_mcc       = NULL;
  const char* csfb_mnc       = NULL;
  const char* lac            = NULL;
  const char* tac_str        = NULL;
  const char* snr_str        = NULL;

  config_init(&cfg);

  if (config_pP->config_file != NULL) {
    /*
     * Read the file. If there is an error, report it and exit.
     */
    if (!config_read_file(&cfg, bdata(config_pP->config_file))) {
      OAILOG_CRITICAL(
          LOG_CONFIG, "Failed to parse MME configuration file: %s:%d - %s\n",
          bdata(config_pP->config_file), config_error_line(&cfg),
          config_error_text(&cfg));
      config_destroy(&cfg);
      Fatal(
          "Failed to parse MME configuration file %s!\n",
          bdata(config_pP->config_file));
    }
  } else {
    config_destroy(&cfg);
    Fatal("No MME configuration file provided!\n");
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
        config_pP->log_config.mme_app_log_level = OAILOG_LEVEL_STR2INT(astring);

      if (config_setting_lookup_string(
              setting, LOG_CONFIG_STRING_S6A_LOG_LEVEL,
              (const char**) &astring))
        config_pP->log_config.s6a_log_level = OAILOG_LEVEL_STR2INT(astring);
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
#if EMBEDDED_SGW
      if (config_setting_lookup_string(
              setting, LOG_CONFIG_STRING_GTPV1U_LOG_LEVEL,
              (const char**) &astring))
        config_pP->log_config.gtpv1u_log_level = OAILOG_LEVEL_STR2INT(astring);

      if (config_setting_lookup_string(
              setting, LOG_CONFIG_STRING_SPGW_APP_LOG_LEVEL,
              (const char**) &astring))
        config_pP->log_config.spgw_app_log_level =
            OAILOG_LEVEL_STR2INT(astring);

#else
      if (config_setting_lookup_string(
              setting, LOG_CONFIG_STRING_GTPV2C_LOG_LEVEL,
              (const char**) &astring))
        config_pP->log_config.gtpv2c_log_level = OAILOG_LEVEL_STR2INT(astring);

      if (config_setting_lookup_string(
              setting, LOG_CONFIG_STRING_S11_LOG_LEVEL,
              (const char**) &astring))
        config_pP->log_config.s11_log_level = OAILOG_LEVEL_STR2INT(astring);
#endif
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

    // GENERAL MME SETTINGS
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
            setting_mme, MME_CONFIG_STRING_STATS_TIMER, &aint))) {
      config_pP->stats_timer_sec                           = (uint32_t) aint;
      config_pP->service303_config.stats_display_timer_sec = (uint32_t) aint;
    }

    if ((config_setting_lookup_string(
            setting_mme, MME_CONFIG_STRING_USE_STATELESS,
            (const char**) &astring))) {
      config_pP->use_stateless = parse_bool(astring);
    }

    if ((config_setting_lookup_string(
            setting_mme, MME_CONFIG_STRING_ENABLE_CONVERGED_CORE,
            (const char**) &astring))) {
      config_pP->enable_converged_core = parse_bool(astring);
    }

    if ((config_setting_lookup_string(
            setting_mme, MME_CONFIG_STRING_USE_HA, (const char**) &astring))) {
      config_pP->use_ha = parse_bool(astring);
    }

    if ((config_setting_lookup_string(
            setting_mme, MME_CONFIG_STRING_ENABLE_GTPU_PRIVATE_IP_CORRECTION,
            (const char**) &astring))) {
      config_pP->enable_gtpu_private_ip_correction = parse_bool(astring);
    }

    if ((config_setting_lookup_string(
            setting_mme, MME_CONFIG_STRING_CONGESTION_CONTROL_ENABLED,
            (const char**) &astring))) {
      config_pP->enable_congestion_control = parse_bool(astring);
    }

    if ((config_setting_lookup_int(
            setting_mme, MME_CONFIG_STRING_S1AP_ZMQ_TH, &aint))) {
      config_pP->s1ap_zmq_th = (long) aint;
    }

    if ((config_setting_lookup_int(
            setting_mme, MME_CONFIG_STRING_MME_APP_ZMQ_CONGEST_TH, &aint))) {
      config_pP->mme_app_zmq_congest_th = (long) aint;
    }

    if ((config_setting_lookup_int(
            setting_mme, MME_CONFIG_STRING_MME_APP_ZMQ_AUTH_TH, &aint))) {
      config_pP->mme_app_zmq_auth_th = (long) aint;
    }

    if ((config_setting_lookup_int(
            setting_mme, MME_CONFIG_STRING_MME_APP_ZMQ_IDENT_TH, &aint))) {
      config_pP->mme_app_zmq_ident_th = (long) aint;
    }

    if ((config_setting_lookup_int(
            setting_mme, MME_CONFIG_STRING_MME_APP_ZMQ_SMC_TH, &aint))) {
      config_pP->mme_app_zmq_smc_th = (long) aint;
    }

    if ((config_setting_lookup_string(
            setting_mme,
            EPS_NETWORK_FEATURE_SUPPORT_EMERGENCY_BEARER_SERVICES_IN_S1_MODE,
            (const char**) &astring))) {
      config_pP->eps_network_feature_support
          .emergency_bearer_services_in_s1_mode = parse_bool(astring);
    }
    if ((config_setting_lookup_string(
            setting_mme, EPS_NETWORK_FEATURE_SUPPORT_EXTENDED_SERVICE_REQUEST,
            (const char**) &astring))) {
      config_pP->eps_network_feature_support.extended_service_request =
          parse_bool(astring);
    }
    if ((config_setting_lookup_string(
            setting_mme,
            EPS_NETWORK_FEATURE_SUPPORT_IMS_VOICE_OVER_PS_SESSION_IN_S1,
            (const char**) &astring))) {
      config_pP->eps_network_feature_support.ims_voice_over_ps_session_in_s1 =
          parse_bool(astring);
    }
    if ((config_setting_lookup_string(
            setting_mme, EPS_NETWORK_FEATURE_SUPPORT_LOCATION_SERVICES_VIA_EPC,
            (const char**) &astring))) {
      config_pP->eps_network_feature_support.location_services_via_epc =
          parse_bool(astring);
    }

    if ((config_setting_lookup_string(
            setting_mme, MME_CONFIG_STRING_UNAUTHENTICATED_IMSI_SUPPORTED,
            (const char**) &astring))) {
      config_pP->unauthenticated_imsi_supported = parse_bool(astring);
    }

    // ITTI SETTING
    setting = config_setting_get_member(
        setting_mme, MME_CONFIG_STRING_INTERTASK_INTERFACE_CONFIG);

    if (setting != NULL) {
      if ((config_setting_lookup_int(
              setting, MME_CONFIG_STRING_INTERTASK_INTERFACE_QUEUE_SIZE,
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
              setting, MME_CONFIG_STRING_S6A_CONF_FILE_PATH,
              (const char**) &astring))) {
        if (astring != NULL) {
          if (config_pP->s6a_config.conf_file) {
            bassigncstr(config_pP->s6a_config.conf_file, astring);
          } else {
            config_pP->s6a_config.conf_file = bfromcstr(astring);
          }
        }
      }

      if ((config_setting_lookup_string(
              setting, MME_CONFIG_STRING_S6A_HSS_HOSTNAME,
              (const char**) &astring))) {
        if (astring != NULL) {
          if (config_pP->s6a_config.hss_host_name) {
            bassigncstr(config_pP->s6a_config.hss_host_name, astring);
          } else {
            config_pP->s6a_config.hss_host_name = bfromcstr(astring);
          }
        } else
          Fatal(
              "You have to provide a valid HSS hostname %s=...\n",
              MME_CONFIG_STRING_S6A_HSS_HOSTNAME);
      }
      if ((config_setting_lookup_string(
              setting, MME_CONFIG_STRING_S6A_HSS_REALM,
              (const char**) &astring))) {
        if (astring != NULL) {
          if (config_pP->s6a_config.hss_realm) {
            bassigncstr(config_pP->s6a_config.hss_realm, astring);
          } else {
            config_pP->s6a_config.hss_realm = bfromcstr(astring);
          }
        } else
          Fatal(
              "You have to provide a valid HSS realm %s=...\n",
              MME_CONFIG_STRING_S6A_HSS_REALM);
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
            config_pP->gummei.gummei[i].plmn.mcc_digit1 = (uint8_t) atoi(c);
            c[0]                                        = mcc[1];
            config_pP->gummei.gummei[i].plmn.mcc_digit2 = (uint8_t) atoi(c);
            c[0]                                        = mcc[2];
            config_pP->gummei.gummei[i].plmn.mcc_digit3 = (uint8_t) atoi(c);
          }

          if ((config_setting_lookup_string(
                  sub2setting, MME_CONFIG_STRING_MNC, &mnc))) {
            AssertFatal(
                (strlen(mnc) == MIN_MNC_LENGTH) ||
                    (strlen(mnc) == MAX_MNC_LENGTH),
                "Bad MNC length (%ld), it must be %u or %u digit ex: 12 or 123",
                strlen(mnc), MIN_MNC_LENGTH, MAX_MNC_LENGTH);
            char c[2]                                   = {mnc[0], 0};
            config_pP->gummei.gummei[i].plmn.mnc_digit1 = (uint8_t) atoi(c);
            c[0]                                        = mnc[1];
            config_pP->gummei.gummei[i].plmn.mnc_digit2 = (uint8_t) atoi(c);
            if (3 == strlen(mnc)) {
              c[0]                                        = mnc[2];
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

    // RESTRICTED PLMN SETTING
    setting = config_setting_get_member(
        setting_mme, MME_CONFIG_STRING_RESTRICTED_PLMN_LIST);
    config_pP->restricted_plmn.num = 0;
    OAILOG_INFO(LOG_MME_APP, "MME_CONFIG_STRING_RESTRICTED_PLMN_LIST \n");
    if (setting != NULL) {
      num = config_setting_length(setting);
      OAILOG_INFO(
          LOG_MME_APP, "Number of restricted PLMNs configured =%d\n", num);
      AssertFatal(
          num <= MAX_RESTRICTED_PLMN,
          "Number of restricted PLMNs configured:%d exceeds number of "
          "restricted PLMNs supported :%d \n",
          num, MAX_RESTRICTED_PLMN);

      for (i = 0; i < num; i++) {
        sub2setting = config_setting_get_elem(setting, i);

        if (sub2setting != NULL) {
          if ((config_setting_lookup_string(
                  sub2setting, MME_CONFIG_STRING_MCC, &mcc))) {
            AssertFatal(
                strlen(mcc) == MAX_MCC_LENGTH,
                "Bad MCC length (%ld), it must be %u digit ex: 001\n",
                strlen(mcc), MAX_MCC_LENGTH);
            // NULL terminated string
            AssertFatal(
                mcc[0] >= '0' && mcc[0] <= '9',
                "MCC[0] is not a decimal digit\n");
            config_pP->restricted_plmn.plmn[i].mcc_digit1 = mcc[0] - '0';
            AssertFatal(
                mcc[1] >= '0' && mcc[1] <= '9',
                "MCC[1] is not a decimal digit\n");
            config_pP->restricted_plmn.plmn[i].mcc_digit2 = mcc[1] - '0';
            AssertFatal(
                mcc[2] >= '0' && mcc[2] <= '9',
                "MCC[2] is not a decimal digit\n");
            config_pP->restricted_plmn.plmn[i].mcc_digit3 = mcc[2] - '0';
          }

          if ((config_setting_lookup_string(
                  sub2setting, MME_CONFIG_STRING_MNC, &mnc))) {
            AssertFatal(
                (strlen(mnc) == MIN_MNC_LENGTH) ||
                    (strlen(mnc) == MAX_MNC_LENGTH),
                "Bad MNC length (%ld), it must be %u or %u digit ex: 12 or "
                "123\n",
                strlen(mnc), MIN_MNC_LENGTH, MAX_MNC_LENGTH);
            // NULL terminated string
            AssertFatal(
                mnc[0] >= '0' && mnc[0] <= '9',
                "MNC[0] is not a decimal digit\n");
            config_pP->restricted_plmn.plmn[i].mnc_digit1 = mnc[0] - '0';
            AssertFatal(
                mnc[1] >= '0' && mnc[1] <= '9',
                "MNC[1] is not a decimal digit\n");
            config_pP->restricted_plmn.plmn[i].mnc_digit2 = mnc[1] - '0';
            if (3 == strlen(mnc)) {
              AssertFatal(
                  mnc[2] >= '0' && mnc[2] <= '9',
                  "MNC[2] is not a decimal digit\n");
              config_pP->restricted_plmn.plmn[i].mnc_digit3 = mnc[2] - '0';
            } else {
              config_pP->restricted_plmn.plmn[i].mnc_digit3 = 0x0F;
            }
          }
          config_pP->restricted_plmn.num += 1;
        }
      }
    }

    // MODE MAP SETTING
    setting =
        config_setting_get_member(setting_mme, MME_CONFIG_STRING_FED_MODE_MAP);
    memset(&config_pP->mode_map_config, 0, sizeof(fed_mode_map_t));
    OAILOG_INFO(LOG_MME_APP, "MME_CONFIG_STRING_FED_MODE_MAP \n");
    if (setting != NULL) {
      num = config_setting_length(setting);
      OAILOG_INFO(LOG_MME_APP, "Number of mode maps configured =%d\n", num);
      AssertFatal(
          num <= MAX_FED_MODE_MAP_CONFIG,
          "Number of mode maps configured:%d exceeds number of "
          "mode maps supported :%d \n",
          num, MAX_FED_MODE_MAP_CONFIG);

      for (i = 0; i < num; i++) {
        sub2setting = config_setting_get_elem(setting, i);
        if (sub2setting != NULL) {
          OAILOG_INFO(LOG_MME_APP, "sub2setting\n");
          // MODE
          if ((config_setting_lookup_string(
                  sub2setting, MME_CONFIG_STRING_MODE, &astring))) {
            config_pP->mode_map_config.mode_map[i].mode = atoi(astring);
          }

          // PLMN
          if ((config_setting_lookup_string(
                  sub2setting, MME_CONFIG_STRING_PLMN, &astring))) {
            // NULL terminated string
            char fed_mode_mcc[MAX_MCC_LENGTH + 1],
                fed_mode_mnc[MAX_MNC_LENGTH + 1];
            // Convert to 3gpp PLMN (MCC and MNC) format
            // First 3 chars in astring is MCC next 3 or 2 chars is MNC
            memcpy(fed_mode_mcc, astring, MAX_MCC_LENGTH);
            fed_mode_mcc[MAX_MCC_LENGTH] = '\0';  // null terminated string
            n                            = strlen(astring) - MAX_MCC_LENGTH;
            memcpy(fed_mode_mnc, astring + MAX_MCC_LENGTH, n);
            fed_mode_mnc[n] = '\0';  // null terminated string
            AssertFatal(
                strlen(fed_mode_mcc) == MAX_MCC_LENGTH,
                "Bad MCC length (%ld), it must be %u digit ex: 001\n",
                strlen(fed_mode_mcc), MAX_MCC_LENGTH);
            AssertFatal(
                fed_mode_mcc[0] >= '0' && fed_mode_mcc[0] <= '9',
                "MCC[0] is not a decimal digit\n");
            config_pP->mode_map_config.mode_map[i].plmn.mcc_digit1 =
                fed_mode_mcc[0] - '0';
            AssertFatal(
                fed_mode_mcc[1] >= '0' && fed_mode_mcc[1] <= '9',
                "MCC[1] is not a decimal digit\n");
            config_pP->mode_map_config.mode_map[i].plmn.mcc_digit2 =
                fed_mode_mcc[1] - '0';
            AssertFatal(
                fed_mode_mcc[2] >= '0' && fed_mode_mcc[2] <= '9',
                "MCC[2] is not a decimal digit\n");
            config_pP->mode_map_config.mode_map[i].plmn.mcc_digit3 =
                fed_mode_mcc[2] - '0';

            // MNC
            AssertFatal(
                (strlen(fed_mode_mnc) == MIN_MNC_LENGTH) ||
                    (strlen(fed_mode_mnc) == MAX_MNC_LENGTH),
                "Bad MNC length (%ld), it must be %u or %u digit ex: 12 or "
                "123\n",
                strlen(fed_mode_mnc), MIN_MNC_LENGTH, MAX_MNC_LENGTH);

            // NULL terminated string
            AssertFatal(
                fed_mode_mnc[0] >= '0' && fed_mode_mnc[0] <= '9',
                "MNC[0] is not a decimal digit\n");
            config_pP->mode_map_config.mode_map[i].plmn.mnc_digit1 =
                fed_mode_mnc[0] - '0';
            AssertFatal(
                fed_mode_mnc[1] >= '0' && fed_mode_mnc[1] <= '9',
                "MNC[1] is not a decimal digit\n");
            config_pP->mode_map_config.mode_map[i].plmn.mnc_digit2 =
                fed_mode_mnc[1] - '0';
            if (3 == strlen(fed_mode_mnc)) {
              AssertFatal(
                  fed_mode_mnc[2] >= '0' && fed_mode_mnc[2] <= '9',
                  "MNC[2] is not a decimal digit\n");
              config_pP->mode_map_config.mode_map[i].plmn.mnc_digit3 =
                  fed_mode_mnc[2] - '0';
            } else {
              config_pP->mode_map_config.mode_map[i].plmn.mnc_digit3 = 0x0F;
            }
          }
          // IMSI range
          if ((config_setting_lookup_string(
                  sub2setting, MME_CONFIG_STRING_IMSI_RANGE, &astring))) {
            if (strlen(astring)) {
              imsi_high_tmp = strdup(astring);
              imsi_low_tmp  = strsep(&imsi_high_tmp, ":");
              memcpy(
                  (char*) config_pP->mode_map_config.mode_map[i].imsi_low,
                  imsi_low_tmp, strlen(imsi_low_tmp));
              AssertFatal(
                  strlen((char*) config_pP->mode_map_config.mode_map[i]
                             .imsi_low) <= MAX_IMSI_LENGTH,
                  "Invalid imsi_low length\n");
              memcpy(
                  (char*) config_pP->mode_map_config.mode_map[i].imsi_high,
                  imsi_high_tmp, strlen(imsi_high_tmp));
              AssertFatal(
                  strlen((char*) config_pP->mode_map_config.mode_map[i]
                             .imsi_high) <= MAX_IMSI_LENGTH,
                  "Invalid imsi_high length\n");
            }
          }
          // APN
          if ((config_setting_lookup_string(
                  sub2setting, MME_CONFIG_STRING_APN, &astring))) {
            config_pP->mode_map_config.mode_map[i].apn = bfromcstr(astring);
          }
          config_pP->mode_map_config.num += 1;
        }
      }
    }

    // BLOCKED IMEI LIST SETTING
    setting = config_setting_get_member(
        setting_mme, MME_CONFIG_STRING_BLOCKED_IMEI_LIST);
    char imei_str[MAX_LEN_IMEI + 1] = {0};
    imei64_t imei64                 = 0;
    config_pP->blocked_imei.num     = 0;
    OAILOG_INFO(LOG_MME_APP, "MME_CONFIG_STRING_BLOCKED_IMEI_LIST \n");
    if (setting != NULL) {
      num = config_setting_length(setting);
      OAILOG_INFO(LOG_MME_APP, "Number of blocked IMEIs configured =%d\n", num);
      if (num > 0) {
        // Create IMEI hashtable
        hashtable_rc_t h_rc = HASH_TABLE_OK;
        bstring b           = bfromcstr("mme_app_config_imei_htbl");
        config_pP->blocked_imei.imei_htbl =
            hashtable_uint64_ts_create(MAX_IMEI_HTBL_SZ, NULL, b);
        bdestroy_wrapper(&b);
        AssertFatal(
            config_pP->blocked_imei.imei_htbl != NULL,
            "Error creating IMEI hashtable\n");

        for (i = 0; i < num; i++) {
          memset(imei_str, 0, (MAX_LEN_IMEI + 1));
          sub2setting = config_setting_get_elem(setting, i);
          if (sub2setting != NULL) {
            if ((config_setting_lookup_string(
                    sub2setting, MME_CONFIG_STRING_IMEI_TAC, &tac_str))) {
              AssertFatal(
                  strlen(tac_str) == MAX_LEN_TAC,
                  "Bad TAC length (%ld), it must be %u digits\n",
                  strlen(tac_str), MAX_LEN_TAC);
              memcpy(imei_str, tac_str, strlen(tac_str));
            }
            if ((config_setting_lookup_string(
                    sub2setting, MME_CONFIG_STRING_SNR, &snr_str))) {
              if (strlen(snr_str)) {
                AssertFatal(
                    strlen(snr_str) == MAX_LEN_SNR,
                    "Bad SNR length (%ld), it must be %u digits\n",
                    strlen(snr_str), MAX_LEN_SNR);
                memcpy(&imei_str[strlen(tac_str)], snr_str, strlen(snr_str));
              }
            }
            // Store IMEI into hashlist
            imei64 = 0;
            IMEI_STRING_TO_IMEI64(imei_str, &imei64);
            h_rc = hashtable_uint64_ts_insert(
                config_pP->blocked_imei.imei_htbl, (const hash_key_t) imei64,
                0);
            AssertFatal(h_rc == HASH_TABLE_OK, "Hashtable insertion failed\n");

            config_pP->blocked_imei.num += 1;
          }
        }
      }
    }

    // SRVC_AREA_CODE_2_TACS_MAP
    setting = config_setting_get_member(
        setting_mme, MME_CONFIG_STRING_SRVC_AREA_CODE_2_TACS_MAP);
    OAILOG_INFO(LOG_MME_APP, "MME_CONFIG_STRING_SRVC_AREA_CODE_2_TACS_MAP \n");
    if (setting != NULL) {
      num = config_setting_length(setting);
      OAILOG_INFO(
          LOG_MME_APP, "Number of SRVC_AREA_CODE_2_TACS configured =%d\n", num);
      if (num > 0) {
        // Create SRVC_AREA_CODE_2_TACS hashtable
        hashtable_rc_t h_rc = HASH_TABLE_OK;
        bstring b           = bfromcstr("mme_app_config_sac_2_tacs_htbl");
        config_pP->sac_to_tacs_map.sac_to_tacs_map_htbl =
            obj_hashtable_create(MAX_SAC_2_TACS_HTBL_SZ, NULL, NULL, NULL, b);
        bdestroy_wrapper(&b);
        if (config_pP->sac_to_tacs_map.sac_to_tacs_map_htbl == NULL) {
          OAILOG_ERROR(
              LOG_MME_APP, "Error creating SAC_2_TACS_HTBL hashtable \n");
          return -1;
        }
        for (i = 0; i < num; i++) {
          sub2setting = config_setting_get_elem(setting, i);
          if (sub2setting != NULL) {
            if ((config_setting_lookup_int(
                    sub2setting, MME_CONFIG_STRING_SERVICE_AREA_CODE, &aint))) {
              // store in network byte order as SAC will come from
              // the network in ULA messsage.
              uint16_t sac_int = htons((uint16_t) aint);
              // TAC LIST
              sub3setting = config_setting_get_member(
                  sub2setting, MME_CONFIG_STRING_TAC_LIST_PER_SAC);
              if (sub3setting) {
                uint8_t num_tacs = config_setting_length(sub3setting);
                if (num_tacs > 0) {
                  config_pP->sac_to_tacs_map.tac_list =
                      calloc(1, sizeof(tac_list_per_sac_t));
                  AssertFatal(
                      config_pP->sac_to_tacs_map.tac_list != NULL,
                      "Memory allocation failed for tac_list\n");
                  config_pP->sac_to_tacs_map.tac_list->num_tac_entries =
                      num_tacs;
                  for (uint8_t itr = 0; itr < num_tacs; itr++) {
                    config_pP->sac_to_tacs_map.tac_list->tacs[itr] =
                        config_setting_get_int_elem(sub3setting, itr);
                  }
                  h_rc = obj_hashtable_insert(
                      config_pP->sac_to_tacs_map.sac_to_tacs_map_htbl,
                      (const void*) &sac_int, sizeof(uint16_t),
                      (void*) config_pP->sac_to_tacs_map.tac_list);
                  AssertFatal(
                      h_rc == HASH_TABLE_OK,
                      "SAC_2_TACS_HTBL hashtable insertion failed\n");
                }
              }
            }
          }
        }
      }
    }

    // NETWORK INTERFACE SETTING
    setting = config_setting_get_member(
        setting_mme, MME_CONFIG_STRING_NETWORK_INTERFACES_CONFIG);

    if (setting != NULL) {
      if ((config_setting_lookup_string(
               setting, MME_CONFIG_STRING_INTERFACE_NAME_FOR_S1_MME,
               (const char**) &if_name_s1_mme) &&
           config_setting_lookup_string(
               setting, MME_CONFIG_STRING_IPV4_ADDRESS_FOR_S1_MME,
               (const char**) &s1_mme) &&
           config_setting_lookup_string(
               setting, MME_CONFIG_STRING_INTERFACE_NAME_FOR_S11_MME,
               (const char**) &if_name_s11) &&
           config_setting_lookup_string(
               setting, MME_CONFIG_STRING_IPV4_ADDRESS_FOR_S11_MME,
               (const char**) &s11) &&
           config_setting_lookup_int(
               setting, MME_CONFIG_STRING_MME_PORT_FOR_S11, &aint))) {
        config_pP->ip.port_s11 = (uint16_t) aint;

        config_pP->ip.if_name_s1_mme = bfromcstr(if_name_s1_mme);
        cidr                         = bfromcstr(s1_mme);
        struct bstrList* list        = bsplit(cidr, '/');
        AssertFatal(
            list->qty == CIDR_SPLIT_LIST_COUNT, "Bad S1-MME CIDR address: %s",
            bdata(cidr));
        address = list->entry[0];
        mask    = list->entry[1];
        IPV4_STR_ADDR_TO_INADDR(
            bdata(address), config_pP->ip.s1_mme_v4,
            "BAD IP ADDRESS FORMAT FOR S1-MME !\n");
        config_pP->ip.netmask_s1_mme = atoi((const char*) mask->data);
        bstrListDestroy(list);
        in_addr_var.s_addr = config_pP->ip.s1_mme_v4.s_addr;
        OAILOG_INFO(
            LOG_MME_APP,
            "Parsing configuration file found S1-MME: %s/%d on %s\n",
            inet_ntoa(in_addr_var), config_pP->ip.netmask_s1_mme,
            bdata(config_pP->ip.if_name_s1_mme));
        bdestroy_wrapper(&cidr);

        bdestroy(cidr);
        config_pP->ip.if_name_s11 = bfromcstr(if_name_s11);
        cidr                      = bfromcstr(s11);
        list                      = bsplit(cidr, '/');
        AssertFatal(
            list->qty == CIDR_SPLIT_LIST_COUNT, "Bad MME S11 CIDR address: %s",
            bdata(cidr));
        address = list->entry[0];
        mask    = list->entry[1];
        IPV4_STR_ADDR_TO_INADDR(
            bdata(address), config_pP->ip.s11_mme_v4,
            "BAD IP ADDRESS FORMAT FOR S11 !\n");
        config_pP->ip.netmask_s11 = atoi((const char*) mask->data);
        bstrListDestroy(list);
        bdestroy_wrapper(&cidr);
        in_addr_var.s_addr = config_pP->ip.s11_mme_v4.s_addr;
        OAILOG_INFO(
            LOG_MME_APP, "Parsing configuration file found S11: %s/%d on %s\n",
            inet_ntoa(in_addr_var), config_pP->ip.netmask_s11,
            bdata(config_pP->ip.if_name_s11));
        bdestroy(cidr);
      }
    }

    // CSFB SETTING
    setting =
        config_setting_get_member(setting_mme, MME_CONFIG_STRING_CSFB_CONFIG);
    if (setting != NULL) {
      if ((config_setting_lookup_string(
              setting, MME_CONFIG_STRING_NON_EPS_SERVICE_CONTROL,
              (const char**) &astring))) {
        if (astring != NULL) {
          config_pP->non_eps_service_control = bfromcstr(astring);
        }
      }
      if (strcasecmp(
              (const char*) config_pP->non_eps_service_control->data, "OFF") !=
          0) {
        // Check CSFB MCC. MNC and LAC only if NON-EPS feature is enabled.
        if ((config_setting_lookup_string(
                setting, MME_CONFIG_STRING_CSFB_MCC, &csfb_mcc))) {
          AssertFatal(
              strlen(csfb_mcc) == MAX_MCC_LENGTH,
              "Bad MCC length(%ld), it must be %u digit ex: 001",
              strlen(csfb_mcc), MAX_MCC_LENGTH);
          char c[2]                = {csfb_mcc[0], 0};
          config_pP->lai.mccdigit1 = (uint8_t) atoi(c);
          c[0]                     = csfb_mcc[1];
          config_pP->lai.mccdigit2 = (uint8_t) atoi(c);
          c[0]                     = csfb_mcc[2];
          config_pP->lai.mccdigit3 = (uint8_t) atoi(c);
        }
        if ((config_setting_lookup_string(
                setting, MME_CONFIG_STRING_CSFB_MNC, &csfb_mnc))) {
          AssertFatal(
              (strlen(csfb_mnc) == MIN_MNC_LENGTH) ||
                  (strlen(csfb_mnc) == MAX_MNC_LENGTH),
              "Bad MNC length (%ld), it must be %u or %u digit ex: 12 or 123",
              strlen(csfb_mnc), MIN_MNC_LENGTH, MAX_MNC_LENGTH);
          char c[2]                = {csfb_mnc[0], 0};
          config_pP->lai.mncdigit1 = (uint8_t) atoi(c);
          c[0]                     = csfb_mnc[1];
          config_pP->lai.mncdigit2 = (uint8_t) atoi(c);
          if (3 == strlen(csfb_mnc)) {
            c[0]                     = csfb_mnc[2];
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
              setting, MME_CONFIG_STRING_NAS_FORCE_REJECT_TAU,
              (const char**) &astring))) {
        config_pP->nas_config.force_reject_tau = parse_bool(astring);
      }
      if ((config_setting_lookup_string(
              setting, MME_CONFIG_STRING_NAS_FORCE_REJECT_SR,
              (const char**) &astring))) {
        config_pP->nas_config.force_reject_sr = parse_bool(astring);
      }
      if ((config_setting_lookup_string(
              setting, MME_CONFIG_STRING_NAS_DISABLE_ESM_INFORMATION_PROCEDURE,
              (const char**) &astring))) {
        config_pP->nas_config.disable_esm_information = parse_bool(astring);
      }
      if ((config_setting_lookup_string(
              setting, MME_CONFIG_STRING_NAS_ENABLE_APN_CORRECTION,
              (const char**) &astring))) {
        config_pP->nas_config.enable_apn_correction = parse_bool(astring);
      }

      // Parsing APN CORRECTION MAP
      if (config_pP->nas_config.enable_apn_correction) {
        subsetting = config_setting_get_member(
            setting, MME_CONFIG_STRING_NAS_APN_CORRECTION_MAP_LIST);
        config_pP->nas_config.apn_map_config.nb = 0;
        if (subsetting != NULL) {
          num = config_setting_length(subsetting);
          OAILOG_INFO(
              LOG_MME_APP, "Number of apn correction map configured =%d\n",
              num);
          AssertFatal(
              num <= MAX_APN_CORRECTION_MAP_LIST,
              "Number of apn correction map configured:%d exceeds the maximum "
              "number supported"
              ":%d \n",
              num, MAX_APN_CORRECTION_MAP_LIST);

          for (i = 0; i < num; i++) {
            sub2setting = config_setting_get_elem(subsetting, i);
            if (sub2setting != NULL) {
              if ((config_setting_lookup_string(
                      sub2setting,
                      MME_CONFIG_STRING_NAS_APN_CORRECTION_MAP_IMSI_PREFIX,
                      (const char**) &imsi_prefix))) {
                if (config_pP->nas_config.apn_map_config.apn_map[i]
                        .imsi_prefix) {
                  bassigncstr(
                      config_pP->nas_config.apn_map_config.apn_map[i]
                          .imsi_prefix,
                      imsi_prefix);
                } else {
                  config_pP->nas_config.apn_map_config.apn_map[i].imsi_prefix =
                      bfromcstr(imsi_prefix);
                }
              }
              if ((config_setting_lookup_string(
                      sub2setting,
                      MME_CONFIG_STRING_NAS_APN_CORRECTION_MAP_APN_OVERRIDE,
                      (const char**) &apn_override))) {
                if (config_pP->nas_config.apn_map_config.apn_map[i]
                        .apn_override) {
                  bassigncstr(
                      config_pP->nas_config.apn_map_config.apn_map[i]
                          .apn_override,
                      apn_override);
                } else {
                  config_pP->nas_config.apn_map_config.apn_map[i].apn_override =
                      bfromcstr(apn_override);
                }
              }
              config_pP->nas_config.apn_map_config.nb += 1;
            }
          }
        }
      }
    }

    // Parsing Sentry Config
    setting =
        config_setting_get_member(setting_mme, MME_CONFIG_STRING_SENTRY_CONFIG);
    memset(&config_pP->sentry_config, 0, sizeof(sentry_config_t));
    config_pP->sentry_config.url_native = bfromcstr("");
    OAILOG_INFO(LOG_MME_APP, "MME_CONFIG_STRING_SENTRY_CONFIG \n");
    if (setting != NULL) {
      if ((config_setting_lookup_float(
              setting, MME_CONFIG_STRING_SAMPLE_RATE, &adouble))) {
        config_pP->sentry_config.sample_rate = (float) adouble;
      }
      if ((config_setting_lookup_string(
              setting, MME_CONFIG_STRING_UPLOAD_MME_LOG,
              (const char**) &astring))) {
        config_pP->sentry_config.upload_mme_log = parse_bool(astring);
      }
      if ((config_setting_lookup_string(
              setting, MME_CONFIG_STRING_URL_NATIVE,
              (const char**) &astring))) {
        bassigncstr(config_pP->sentry_config.url_native, astring);
      }
    }

    // SGS TIMERS
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
#if (!EMBEDDED_SGW)
    // S-GW Setting
    setting =
        config_setting_get_member(setting_mme, MME_CONFIG_STRING_SGW_CONFIG);

    if (setting != NULL) {
      if ((config_setting_lookup_string(
              setting, MME_CONFIG_STRING_SGW_IPV4_ADDRESS_FOR_S11,
              (const char**) &sgw_ip_address_for_s11))) {
        OAILOG_DEBUG(
            LOG_MME_APP, "sgw interface IP information %s\n",
            sgw_ip_address_for_s11);

        IPV4_STR_ADDR_TO_INADDR(
            sgw_ip_address_for_s11, config_pP->e_dns_emulation.sgw_ip_addr[0],
            "BAD IP ADDRESS FORMAT FOR SGW S11 !\n");

        OAILOG_INFO(
            LOG_SPGW_APP, "Parsing configuration file found S-GW S11: %s\n",
            inet_ntoa(config_pP->e_dns_emulation.sgw_ip_addr[0]));
      }
    }
#endif
  }

  config_destroy(&cfg);
  return 0;
}

//------------------------------------------------------------------------------
void mme_config_display(mme_config_t* config_pP) {
  int j;

  OAILOG_INFO(
      LOG_CONFIG, "==== EURECOM %s v%s ====\n", PACKAGE_NAME, PACKAGE_VERSION);
  OAILOG_DEBUG(
      LOG_CONFIG, "Built with EMBEDDED_SGW .................: %d\n",
      EMBEDDED_SGW);
  OAILOG_DEBUG(
      LOG_CONFIG, "Built with S6A_OVER_GRPC .....................: %d\n",
      S6A_OVER_GRPC);

#if DEBUG_IS_ON
  OAILOG_DEBUG(
      LOG_CONFIG, "Built with CMAKE_BUILD_TYPE ................: %s\n",
      CMAKE_BUILD_TYPE);
  OAILOG_DEBUG(
      LOG_CONFIG, "Built with PACKAGE_NAME ....................: %s\n",
      PACKAGE_NAME);
  OAILOG_DEBUG(
      LOG_CONFIG, "Built with S1AP_DEBUG_LIST .................: %d\n",
      S1AP_DEBUG_LIST);
  OAILOG_DEBUG(
      LOG_CONFIG, "Built with SCTP_DUMP_LIST ..................: %d\n",
      SCTP_DUMP_LIST);
  OAILOG_DEBUG(
      LOG_CONFIG, "Built with TRACE_HASHTABLE .................: %d\n",
      TRACE_HASHTABLE);
  OAILOG_DEBUG(
      LOG_CONFIG, "Built with TRACE_3GPP_SPEC .................: %d\n",
      TRACE_3GPP_SPEC);
#endif
  OAILOG_INFO(LOG_CONFIG, "Configuration:\n");
  OAILOG_INFO(
      LOG_CONFIG, "- File .................................: %s\n",
      bdata(config_pP->config_file));
  OAILOG_INFO(
      LOG_CONFIG, "- Realm ................................: %s\n",
      bdata(config_pP->realm));
  OAILOG_INFO(
      LOG_CONFIG, "  full network name ....................: %s\n",
      bdata(config_pP->full_network_name));
  OAILOG_INFO(
      LOG_CONFIG, "  short network name ...................: %s\n",
      bdata(config_pP->short_network_name));
  OAILOG_INFO(
      LOG_CONFIG, "  Daylight Saving Time..................: %d\n",
      config_pP->daylight_saving_time);
  OAILOG_INFO(
      LOG_CONFIG, "- Run mode .............................: %s\n",
      (RUN_MODE_TEST == config_pP->run_mode) ? "TEST" : "NORMAL");
  OAILOG_INFO(
      LOG_CONFIG, "- Max eNBs .............................: %u\n",
      config_pP->max_enbs);
  OAILOG_INFO(
      LOG_CONFIG, "- Max UEs ..............................: %u\n",
      config_pP->max_ues);
  OAILOG_INFO(
      LOG_CONFIG, "- IMS voice over PS session in S1 ......: %s\n",
      config_pP->eps_network_feature_support.ims_voice_over_ps_session_in_s1 ==
              0 ?
          "false" :
          "true");
  OAILOG_INFO(
      LOG_CONFIG, "- Emergency bearer services in S1 mode .: %s\n",
      config_pP->eps_network_feature_support
                  .emergency_bearer_services_in_s1_mode == 0 ?
          "false" :
          "true");
  OAILOG_INFO(
      LOG_CONFIG, "- Location services via epc ............: %s\n",
      config_pP->eps_network_feature_support.location_services_via_epc == 0 ?
          "false" :
          "true");
  OAILOG_INFO(
      LOG_CONFIG, "- Extended service request .............: %s\n",
      config_pP->eps_network_feature_support.extended_service_request == 0 ?
          "false" :
          "true");
  OAILOG_INFO(
      LOG_CONFIG, "- Unauth IMSI support ..................: %s\n",
      config_pP->unauthenticated_imsi_supported == 0 ? "false" : "true");
  OAILOG_INFO(
      LOG_CONFIG, "- Relative capa ........................: %u\n",
      config_pP->relative_capacity);
  OAILOG_INFO(
      LOG_CONFIG, "- Statistics timer .....................: %u (seconds)\n\n",
      config_pP->stats_timer_sec);
  OAILOG_INFO(
      LOG_CONFIG, "- Congestion control enabled ........................: %s\n",
      config_pP->enable_congestion_control ? "true" : "false");
  OAILOG_INFO(
      LOG_CONFIG,
      "- S1AP ZMQ Threshold ...........................: %10ld "
      "(microseconds)\n",
      config_pP->s1ap_zmq_th);
  OAILOG_INFO(
      LOG_CONFIG,
      "- MME APP ZMQ Congestion Threshold .............: %10ld "
      "(microseconds)\n",
      config_pP->mme_app_zmq_congest_th);
  OAILOG_INFO(
      LOG_CONFIG,
      "- MME APP ZMQ Auth Complete Threshold...........: %10ld "
      "(microseconds)\n",
      config_pP->mme_app_zmq_auth_th);
  OAILOG_INFO(
      LOG_CONFIG,
      "- MME APP ZMQ Identity Complete Threshold.......: %10ld "
      "(microseconds)\n",
      config_pP->mme_app_zmq_ident_th);
  OAILOG_INFO(
      LOG_CONFIG,
      "- MME APP ZMQ SMC Complete Threshold ...........: %10ld "
      "(microseconds)\n\n",
      config_pP->mme_app_zmq_smc_th);
  OAILOG_INFO(
      LOG_CONFIG, "- Use Stateless ........................: %s\n\n",
      config_pP->use_stateless ? "true" : "false");
  OAILOG_INFO(
      LOG_CONFIG, "- enable_converged_core .......: %s\n\n",
      config_pP->enable_converged_core ? "true" : "false");
  OAILOG_INFO(LOG_CONFIG, "- CSFB:\n");
  OAILOG_INFO(
      LOG_CONFIG,
      "    Non EPS Service Control ........................: %s\n\n",
      bdata(config_pP->non_eps_service_control));
  OAILOG_INFO(LOG_CONFIG, "- S1-MME:\n");
  OAILOG_INFO(
      LOG_CONFIG, "    port number ......: %d\n",
      config_pP->s1ap_config.port_number);
  OAILOG_INFO(LOG_CONFIG, "- IP:\n");
  OAILOG_INFO(
      LOG_CONFIG, "    s1-MME iface .....: %s\n",
      bdata(config_pP->ip.if_name_s1_mme));
  OAILOG_INFO(
      LOG_CONFIG, "    s1-MME ip ........: %s\n",
      inet_ntoa(*((struct in_addr*) &config_pP->ip.s1_mme_v4)));
  OAILOG_INFO(
      LOG_CONFIG, "    s11 MME iface ....: %s\n",
      bdata(config_pP->ip.if_name_s11));
  OAILOG_INFO(
      LOG_CONFIG, "    s11 MME port .....: %d\n", config_pP->ip.port_s11);
  OAILOG_INFO(
      LOG_CONFIG, "    s11 MME ip .......: %s\n",
      inet_ntoa(*((struct in_addr*) &config_pP->ip.s11_mme_v4)));

  if (config_pP->e_dns_emulation.sgw_ip_addr[0].s_addr == AF_INET) {
    OAILOG_INFO(
        LOG_CONFIG, " Address : %s\n",
        inet_ntoa(*((struct in_addr*) &config_pP->e_dns_emulation.sgw_ip_addr[0]
                        .s_addr)));

  } else if (config_pP->e_dns_emulation.sgw_ip_addr[0].s_addr == AF_INET6) {
    char strv6[16];
    OAILOG_INFO(
        LOG_CONFIG, " Address : %s\n",
        inet_ntop(
            AF_INET6, &config_pP->e_dns_emulation.sgw_ip_addr[0].s_addr, strv6,
            16));
  } else {
    OAILOG_INFO(
        LOG_CONFIG, "  Address : Unknown address family %d\n",
        config_pP->e_dns_emulation.sgw_ip_addr[0].s_addr);
  }

  OAILOG_INFO(LOG_CONFIG, "- Sentry:\n");
  OAILOG_INFO(
      LOG_CONFIG, "    sample rate ......: %f\n",
      config_pP->sentry_config.sample_rate);
  OAILOG_INFO(
      LOG_CONFIG, "    upload MME log ...: %d\n",
      config_pP->sentry_config.upload_mme_log);
  OAILOG_INFO(
      LOG_CONFIG, "    URL native .......: %s\n",
      bdata(config_pP->sentry_config.url_native));

  OAILOG_INFO(LOG_CONFIG, "- ITTI:\n");
  OAILOG_INFO(
      LOG_CONFIG, "    queue size .......: %u (bytes)\n",
      config_pP->itti_config.queue_size);
  OAILOG_INFO(
      LOG_CONFIG, "    log file .........: %s\n",
      bdata(config_pP->itti_config.log_file));
  OAILOG_INFO(LOG_CONFIG, "- SCTP:\n");
  OAILOG_INFO(
      LOG_CONFIG, "    in streams .......: %u\n",
      config_pP->sctp_config.in_streams);
  OAILOG_INFO(
      LOG_CONFIG, "    out streams ......: %u\n",
      config_pP->sctp_config.out_streams);
  OAILOG_INFO(LOG_CONFIG, "- GUMMEIs (PLMN|MMEGI|MMEC):\n");
  for (j = 0; j < config_pP->gummei.nb; j++) {
    OAILOG_INFO(
        LOG_CONFIG, "            " PLMN_FMT "|%u|%u \n",
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
          LOG_CONFIG, "            %3u.%3u:%u\n",
          config_pP->served_tai.plmn_mcc[j], config_pP->served_tai.plmn_mnc[j],
          config_pP->served_tai.tac[j]);
    } else {
      OAILOG_INFO(
          LOG_CONFIG, "            %3u.%03u:%u\n",
          config_pP->served_tai.plmn_mcc[j], config_pP->served_tai.plmn_mnc[j],
          config_pP->served_tai.tac[j]);
    }
  }
  for (j = 0; j < config_pP->mode_map_config.num; j++) {
    OAILOG_INFO(LOG_CONFIG, "- MODE MAP : \n");
    OAILOG_INFO(
        LOG_CONFIG, "  - MODE : %d \n",
        config_pP->mode_map_config.mode_map[j].mode);
    OAILOG_INFO(
        LOG_CONFIG, "  - MCC MNC : %u,%u,%u,%u,%u,%u\n",
        config_pP->mode_map_config.mode_map[j].plmn.mcc_digit1,
        config_pP->mode_map_config.mode_map[j].plmn.mcc_digit2,
        config_pP->mode_map_config.mode_map[j].plmn.mcc_digit3,
        config_pP->mode_map_config.mode_map[j].plmn.mnc_digit1,
        config_pP->mode_map_config.mode_map[j].plmn.mnc_digit2,
        config_pP->mode_map_config.mode_map[j].plmn.mnc_digit3);
    OAILOG_INFO(
        LOG_CONFIG, "  - IMSI_LOW : %s\n",
        config_pP->mode_map_config.mode_map[j].imsi_low);
    OAILOG_INFO(
        LOG_CONFIG, "  - IMSI_HIGH : %s\n",
        config_pP->mode_map_config.mode_map[j].imsi_high);
    OAILOG_INFO(
        LOG_CONFIG, "  - APN : %s\n",
        bdata(config_pP->mode_map_config.mode_map[j].apn));
  }
  OAILOG_INFO(LOG_CONFIG, "- NAS:\n");
  OAILOG_INFO(
      LOG_CONFIG,
      "    Preferred Integrity Algorithms .: EIA%d EIA%d EIA%d EIA%d "
      "(decreasing "
      "priority)\n",
      config_pP->nas_config.prefered_integrity_algorithm[0],
      config_pP->nas_config.prefered_integrity_algorithm[1],
      config_pP->nas_config.prefered_integrity_algorithm[2],
      config_pP->nas_config.prefered_integrity_algorithm[3]);
  OAILOG_INFO(
      LOG_CONFIG,
      "    Preferred Integrity Algorithms .: EEA%d EEA%d EEA%d EEA%d "
      "(decreasing "
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
      LOG_CONFIG, "      Force reject TAU ............: %s\n",
      (config_pP->nas_config.force_reject_tau) ? "true" : "false");
  OAILOG_INFO(
      LOG_CONFIG, "      Force reject SR .............: %s\n",
      (config_pP->nas_config.force_reject_sr) ? "true" : "false");
  OAILOG_INFO(
      LOG_CONFIG, "      Disable Esm information .....: %s\n",
      (config_pP->nas_config.disable_esm_information) ? "true" : "false");
  OAILOG_INFO(
      LOG_CONFIG, "      Enable APN Correction .......: %s\n",
      (config_pP->nas_config.enable_apn_correction) ? "true" : "false");

  OAILOG_INFO(
      LOG_CONFIG,
      "      APN CORRECTION MAP LIST (IMSI_PREFIX | "
      "APN_OVERRIDE):\n");
  for (j = 0; j < config_pP->nas_config.apn_map_config.nb; j++) {
    OAILOG_INFO(
        LOG_CONFIG, "                                %s | %s \n",
        bdata(config_pP->nas_config.apn_map_config.apn_map[j].imsi_prefix),
        bdata(config_pP->nas_config.apn_map_config.apn_map[j].apn_override));
  }
  OAILOG_INFO(LOG_CONFIG, "- S6A:\n");
#if S6A_OVER_GRPC
  OAILOG_INFO(LOG_CONFIG, "    protocol .........: gRPC\n");
#else
  OAILOG_INFO(LOG_CONFIG, "    protocol .........: diameter\n");
  OAILOG_INFO(
      LOG_CONFIG, "    conf file ........: %s\n",
      bdata(config_pP->s6a_config.conf_file));
#endif
  OAILOG_INFO(LOG_CONFIG, "- Service303:\n");
  OAILOG_INFO(
      LOG_CONFIG, "    service name ........: %s\n",
      bdata(config_pP->service303_config.name));
  OAILOG_INFO(
      LOG_CONFIG, "    version ........: %s\n",
      bdata(config_pP->service303_config.version));
  OAILOG_INFO(LOG_CONFIG, "- Logging:\n");
  OAILOG_INFO(
      LOG_CONFIG, "    Output ..............: %s\n",
      bdata(config_pP->log_config.output));
  OAILOG_INFO(
      LOG_CONFIG, "    Output thread safe ..: %s\n",
      (config_pP->log_config.is_output_thread_safe) ? "true" : "false");
  OAILOG_INFO(
      LOG_CONFIG, "    Output with color ...: %s\n",
      (config_pP->log_config.color) ? "true" : "false");
  OAILOG_INFO(
      LOG_CONFIG, "    UDP log level........: %s\n",
      OAILOG_LEVEL_INT2STR(config_pP->log_config.udp_log_level));
  OAILOG_INFO(
      LOG_CONFIG, "    GTPV1-U log level....: %s\n",
      OAILOG_LEVEL_INT2STR(config_pP->log_config.gtpv1u_log_level));
  OAILOG_INFO(
      LOG_CONFIG, "    GTPV2-C log level....: %s\n",
      OAILOG_LEVEL_INT2STR(config_pP->log_config.gtpv2c_log_level));
  OAILOG_INFO(
      LOG_CONFIG, "    SCTP log level.......: %s\n",
      OAILOG_LEVEL_INT2STR(config_pP->log_config.sctp_log_level));
  OAILOG_INFO(
      LOG_CONFIG, "    S1AP log level.......: %s\n",
      OAILOG_LEVEL_INT2STR(config_pP->log_config.s1ap_log_level));
  OAILOG_INFO(
      LOG_CONFIG, "    ASN1 Verbosity level : %d\n",
      config_pP->log_config.asn1_verbosity_level);
  OAILOG_INFO(
      LOG_CONFIG, "    NAS log level........: %s\n",
      OAILOG_LEVEL_INT2STR(config_pP->log_config.nas_log_level));
  OAILOG_INFO(
      LOG_CONFIG, "    MME_APP log level....: %s\n",
      OAILOG_LEVEL_INT2STR(config_pP->log_config.mme_app_log_level));
  OAILOG_INFO(
      LOG_CONFIG, "    SPGW_APP log level....: %s\n",
      OAILOG_LEVEL_INT2STR(config_pP->log_config.spgw_app_log_level));
  OAILOG_INFO(
      LOG_CONFIG, "    S11 log level........: %s\n",
      OAILOG_LEVEL_INT2STR(config_pP->log_config.s11_log_level));
  OAILOG_INFO(
      LOG_CONFIG, "    S6a log level........: %s\n",
      OAILOG_LEVEL_INT2STR(config_pP->log_config.s6a_log_level));
  OAILOG_INFO(
      LOG_CONFIG, "    UTIL log level.......: %s\n",
      OAILOG_LEVEL_INT2STR(config_pP->log_config.util_log_level));
  OAILOG_INFO(
      LOG_CONFIG, "    ITTI log level.......: %s (InTer-Task Interface)\n",
      OAILOG_LEVEL_INT2STR(config_pP->log_config.itti_log_level));
}

//------------------------------------------------------------------------------
static void usage(char* target) {
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
int mme_config_parse_opt_line(int argc, char* argv[], mme_config_t* config_pP) {
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
            "%s mme_config.config_file %s\n", __FUNCTION__,
            bdata(config_pP->config_file));
      } break;

      case 'v': {
        config_pP->log_config.asn1_verbosity_level = atoi(optarg);
      } break;

      case 'V': {
        OAI_FPRINTF_INFO(
            "==== EURECOM %s v%s ===="
            "Please report any bug to: %s\n",
            PACKAGE_NAME, PACKAGE_VERSION, PACKAGE_BUGREPORT);
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

static bool parse_bool(const char* str) {
  if (strcasecmp(str, "yes") == 0) return true;
  if (strcasecmp(str, "true") == 0) return true;
  if (strcasecmp(str, "no") == 0) return false;
  if (strcasecmp(str, "false") == 0) return false;
  if (strcasecmp(str, "") == 0) return false;

  Fatal("Error in config file: got \"%s\" but expected bool\n", str);
}
