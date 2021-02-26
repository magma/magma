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

/*! \file spgw_config.c
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#define PGW
#define PGW_CONFIG_C

#include <pthread.h>
#include <string.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <stdlib.h>
#include <stdbool.h>
#include <libconfig.h>

#include "bstrlib.h"
#include "assertions.h"
#include "log.h"
#include "common_defs.h"
#include "sgw_config.h"

#ifdef LIBCONFIG_LONG
#define libconfig_int long
#else
#define libconfig_int int
#endif

#ifdef LIBCONFIG_LONG
#define libconfig_int long
#else
#define libconfig_int int
#endif

//------------------------------------------------------------------------------
void sgw_config_init(sgw_config_t* config_pP) {
  memset(config_pP, 0, sizeof(*config_pP));
  pthread_rwlock_init(&config_pP->rw_lock, NULL);
}
//------------------------------------------------------------------------------
int sgw_config_process(sgw_config_t* config_pP) {
  int ret = RETURNok;
  return ret;
}

//------------------------------------------------------------------------------
int sgw_config_parse_file(sgw_config_t* config_pP)

{
  config_t cfg                             = {0};
  config_setting_t* setting_sgw            = NULL;
  char* sgw_if_name_S1u_S12_S4_up          = NULL;
  char* S1u_S12_S4_up                      = NULL;
  char* sgw_if_name_S5_S8_up               = NULL;
  char* S5_S8_up                           = NULL;
  char* sgw_if_name_S11                    = NULL;
  char* S11                                = NULL;
  libconfig_int sgw_udp_port_S1u_S12_S4_up = 2152;
  config_setting_t* subsetting             = NULL;
#if (!EMBEDDED_SGW)
  const char* astring = NULL;
#endif
  bstring address            = NULL;
  bstring cidr               = NULL;
  bstring mask               = NULL;
  struct in_addr in_addr_var = {0};
  (void) in_addr_var;

  config_init(&cfg);

  if (config_pP->config_file) {
    /*
     * Read the file. If there is an error, report it and exit.
     */
    if (!config_read_file(&cfg, bdata(config_pP->config_file))) {
      OAILOG_ERROR(
          LOG_SPGW_APP, "%s:%d - %s\n", bdata(config_pP->config_file),
          config_error_line(&cfg), config_error_text(&cfg));
      config_destroy(&cfg);
      AssertFatal(
          1 == 0, "Failed to parse SP-GW configuration file %s!\n",
          bdata(config_pP->config_file));
    }
  } else {
    OAILOG_ERROR(LOG_SPGW_APP, "No SP-GW configuration file provided!\n");
    config_destroy(&cfg);
    AssertFatal(0, "No SP-GW configuration file provided!\n");
  }

  OAILOG_INFO(
      LOG_SPGW_APP, "Parsing configuration file provided %s\n",
      bdata(config_pP->config_file));
  setting_sgw = config_lookup(&cfg, SGW_CONFIG_STRING_SGW_CONFIG);
  if (setting_sgw) {
#if (!EMBEDDED_SGW)
    // LOGGING setting
    subsetting =
        config_setting_get_member(setting_sgw, LOG_CONFIG_STRING_LOGGING);

    config_pP->log_config.udp_log_level      = MAX_LOG_LEVEL;  // Means invalid
    config_pP->log_config.gtpv1u_log_level   = MAX_LOG_LEVEL;
    config_pP->log_config.gtpv2c_log_level   = MAX_LOG_LEVEL;
    config_pP->log_config.sctp_log_level     = MAX_LOG_LEVEL;
    config_pP->log_config.s1ap_log_level     = MAX_LOG_LEVEL;
    config_pP->log_config.nas_log_level      = MAX_LOG_LEVEL;
    config_pP->log_config.mme_app_log_level  = MAX_LOG_LEVEL;
    config_pP->log_config.spgw_app_log_level = MAX_LOG_LEVEL;
    config_pP->log_config.s11_log_level      = MAX_LOG_LEVEL;
    config_pP->log_config.s6a_log_level      = MAX_LOG_LEVEL;
    config_pP->log_config.util_log_level     = MAX_LOG_LEVEL;
    config_pP->log_config.itti_log_level     = MAX_LOG_LEVEL;
    config_pP->log_config.async_system_log_level = MAX_LOG_LEVEL;
    if (subsetting) {
      if (config_setting_lookup_string(
              subsetting, LOG_CONFIG_STRING_OUTPUT, (const char**) &astring)) {
        if (astring != NULL) {
          if (config_pP->log_config.output) {
            bassigncstr(config_pP->log_config.output, astring);
          } else {
            config_pP->log_config.output = bfromcstr(astring);
          }
        }
      }

      if (config_setting_lookup_string(
              subsetting, LOG_CONFIG_STRING_OUTPUT_THREAD_SAFE,
              (const char**) &astring)) {
        if (astring != NULL) {
          if (strcasecmp(astring, "yes") == 0) {
            config_pP->log_config.is_output_thread_safe = true;
          } else {
            config_pP->log_config.is_output_thread_safe = false;
          }
        }
      }

      if (config_setting_lookup_string(
              subsetting, LOG_CONFIG_STRING_COLOR, (const char**) &astring)) {
        if (!strcasecmp("yes", astring))
          config_pP->log_config.color = true;
        else
          config_pP->log_config.color = false;
      }
      if (config_setting_lookup_string(
              subsetting, LOG_CONFIG_STRING_UDP_LOG_LEVEL,
              (const char**) &astring)) {
        config_pP->log_config.udp_log_level = OAILOG_LEVEL_STR2INT(astring);
      }

      if (config_setting_lookup_string(
              subsetting, LOG_CONFIG_STRING_GTPV1U_LOG_LEVEL,
              (const char**) &astring)) {
        config_pP->log_config.gtpv1u_log_level = OAILOG_LEVEL_STR2INT(astring);
      }

      if (config_setting_lookup_string(
              subsetting, LOG_CONFIG_STRING_GTPV2C_LOG_LEVEL,
              (const char**) &astring)) {
        config_pP->log_config.gtpv2c_log_level = OAILOG_LEVEL_STR2INT(astring);
      }

      if (config_setting_lookup_string(
              subsetting, LOG_CONFIG_STRING_SPGW_APP_LOG_LEVEL,
              (const char**) &astring)) {
        config_pP->log_config.spgw_app_log_level =
            OAILOG_LEVEL_STR2INT(astring);
      }

      if (config_setting_lookup_string(
              subsetting, LOG_CONFIG_STRING_S11_LOG_LEVEL,
              (const char**) &astring)) {
        config_pP->log_config.s11_log_level = OAILOG_LEVEL_STR2INT(astring);
      }

      if (config_setting_lookup_string(
              subsetting, LOG_CONFIG_STRING_UTIL_LOG_LEVEL,
              (const char**) &astring)) {
        config_pP->log_config.util_log_level = OAILOG_LEVEL_STR2INT(astring);
      }

      if (config_setting_lookup_string(
              subsetting, LOG_CONFIG_STRING_ITTI_LOG_LEVEL,
              (const char**) &astring)) {
        config_pP->log_config.itti_log_level = OAILOG_LEVEL_STR2INT(astring);
      }

      if (config_setting_lookup_string(
              subsetting, LOG_CONFIG_STRING_ASYNC_SYSTEM_LOG_LEVEL,
              (const char**) &astring)) {
        config_pP->log_config.async_system_log_level =
            OAILOG_LEVEL_STR2INT(astring);
      }
    }
#endif
    subsetting = config_setting_get_member(
        setting_sgw, SGW_CONFIG_STRING_NETWORK_INTERFACES_CONFIG);

    if (subsetting) {
      if ((config_setting_lookup_string(
               subsetting,
               SGW_CONFIG_STRING_SGW_INTERFACE_NAME_FOR_S1U_S12_S4_UP,
               (const char**) &sgw_if_name_S1u_S12_S4_up) &&
           config_setting_lookup_string(
               subsetting, SGW_CONFIG_STRING_SGW_IPV4_ADDRESS_FOR_S1U_S12_S4_UP,
               (const char**) &S1u_S12_S4_up) &&
           config_setting_lookup_string(
               subsetting, SGW_CONFIG_STRING_SGW_INTERFACE_NAME_FOR_S5_S8_UP,
               (const char**) &sgw_if_name_S5_S8_up) &&
           config_setting_lookup_string(
               subsetting, SGW_CONFIG_STRING_SGW_IPV4_ADDRESS_FOR_S5_S8_UP,
               (const char**) &S5_S8_up) &&
           config_setting_lookup_string(
               subsetting, SGW_CONFIG_STRING_SGW_INTERFACE_NAME_FOR_S11,
               (const char**) &sgw_if_name_S11) &&
           config_setting_lookup_string(
               subsetting, SGW_CONFIG_STRING_SGW_IPV4_ADDRESS_FOR_S11,
               (const char**) &S11))) {
        config_pP->ipv4.if_name_S1u_S12_S4_up =
            bfromcstr(sgw_if_name_S1u_S12_S4_up);
        cidr                  = bfromcstr(S1u_S12_S4_up);
        struct bstrList* list = bsplit(cidr, '/');
        AssertFatal(2 == list->qty, "Bad CIDR address %s", bdata(cidr));
        address = list->entry[0];
        mask    = list->entry[1];
        IPV4_STR_ADDR_TO_INADDR(
            bdata(address), config_pP->ipv4.S1u_S12_S4_up,
            "BAD IP ADDRESS FORMAT FOR S1u_S12_S4 !\n");
        config_pP->ipv4.netmask_S1u_S12_S4_up = atoi((const char*) mask->data);
        bstrListDestroy(list);
        in_addr_var.s_addr = config_pP->ipv4.S1u_S12_S4_up.s_addr;
        OAILOG_INFO(
            LOG_SPGW_APP,
            "Parsing configuration file found S1u_S12_S4_up: %s/%d on %s\n",
            inet_ntoa(in_addr_var), config_pP->ipv4.netmask_S1u_S12_S4_up,
            bdata(config_pP->ipv4.if_name_S1u_S12_S4_up));
        bdestroy(cidr);
        config_pP->ipv4.if_name_S5_S8_up = bfromcstr(sgw_if_name_S5_S8_up);
        cidr                             = bfromcstr(S5_S8_up);
        list                             = bsplit(cidr, '/');
        AssertFatal(2 == list->qty, "Bad CIDR address %s", bdata(cidr));
        address = list->entry[0];
        mask    = list->entry[1];
        IPV4_STR_ADDR_TO_INADDR(
            bdata(address), config_pP->ipv4.S5_S8_up,
            "BAD IP ADDRESS FORMAT FOR S5_S8 !\n");
        config_pP->ipv4.netmask_S5_S8_up = atoi((const char*) mask->data);
        bstrListDestroy(list);
        in_addr_var.s_addr = config_pP->ipv4.S5_S8_up.s_addr;
        OAILOG_INFO(
            LOG_SPGW_APP,
            "Parsing configuration file found S5_S8_up: %s/%d on %s\n",
            inet_ntoa(in_addr_var), config_pP->ipv4.netmask_S5_S8_up,
            bdata(config_pP->ipv4.if_name_S5_S8_up));

        bdestroy(cidr);
        config_pP->ipv4.if_name_S11 = bfromcstr(sgw_if_name_S11);
        cidr                        = bfromcstr(S11);
        list                        = bsplit(cidr, '/');
        AssertFatal(2 == list->qty, "Bad CIDR address %s", bdata(cidr));
        address = list->entry[0];
        mask    = list->entry[1];
        IPV4_STR_ADDR_TO_INADDR(
            bdata(address), config_pP->ipv4.S11,
            "BAD IP ADDRESS FORMAT FOR S11 !\n");
        config_pP->ipv4.netmask_S11 = atoi((const char*) mask->data);
        bstrListDestroy(list);
        in_addr_var.s_addr = config_pP->ipv4.S11.s_addr;
        OAILOG_INFO(
            LOG_SPGW_APP, "Parsing configuration file found S11: %s/%d on %s\n",
            inet_ntoa(in_addr_var), config_pP->ipv4.netmask_S11,
            bdata(config_pP->ipv4.if_name_S11));
        bdestroy(cidr);
      }

      if (config_setting_lookup_int(
              subsetting, SGW_CONFIG_STRING_SGW_PORT_FOR_S1U_S12_S4_UP,
              &sgw_udp_port_S1u_S12_S4_up)) {
        config_pP->udp_port_S1u_S12_S4_up = sgw_udp_port_S1u_S12_S4_up;
      } else {
        config_pP->udp_port_S1u_S12_S4_up = sgw_udp_port_S1u_S12_S4_up;
      }
    }
#if ENABLE_OPENFLOW
    config_setting_t* ovs_settings =
        config_setting_get_member(setting_sgw, SGW_CONFIG_STRING_OVS_CONFIG);
    if (ovs_settings == NULL) {
      AssertFatal(false, "Couldn't find OVS subsetting in spgw config\n");
    }
    char* ovs_bridge_name                       = NULL;
    libconfig_int gtp_port_num                  = 0;
    libconfig_int mtr_port_num                  = 0;
    libconfig_int internal_sampling_port_num    = 0;
    libconfig_int internal_sampling_fwd_tbl_num = 0;
    libconfig_int uplink_port_num               = 0;
    char* multi_tunnel                          = NULL;
    char* uplink_mac                            = NULL;
    char* pipelined_config_enabled              = NULL;
    if (config_setting_lookup_string(
            ovs_settings, SGW_CONFIG_STRING_OVS_BRIDGE_NAME,
            (const char**) &ovs_bridge_name) &&
        config_setting_lookup_int(
            ovs_settings, SGW_CONFIG_STRING_OVS_GTP_PORT_NUM, &gtp_port_num) &&
        config_setting_lookup_int(
            ovs_settings, SGW_CONFIG_STRING_OVS_UPLINK_PORT_NUM,
            &uplink_port_num) &&
        config_setting_lookup_string(
            ovs_settings, SGW_CONFIG_STRING_OVS_UPLINK_MAC,
            (const char**) &uplink_mac) &&
        config_setting_lookup_int(
            ovs_settings, SGW_CONFIG_STRING_OVS_MTR_PORT_NUM, &mtr_port_num) &&
        config_setting_lookup_int(
            ovs_settings, SGW_CONFIG_STRING_OVS_INTERNAL_SAMPLING_PORT_NUM,
            &internal_sampling_port_num) &&
        config_setting_lookup_int(
            ovs_settings, SGW_CONFIG_STRING_OVS_INTERNAL_SAMPLING_FWD_TBL_NUM,
            &internal_sampling_fwd_tbl_num) &&
        config_setting_lookup_string(
            ovs_settings, SGW_CONFIG_STRING_OVS_MULTI_TUNNEL,
            (const char**) &multi_tunnel) &&
        config_setting_lookup_string(
            ovs_settings, SGW_CONFIG_STRING_OVS_PIPELINED_CONFIG_ENABLED,
            (const char**) &pipelined_config_enabled)) {
      config_pP->ovs_config.bridge_name  = bfromcstr(ovs_bridge_name);
      config_pP->ovs_config.gtp_port_num = gtp_port_num;
      config_pP->ovs_config.mtr_port_num = mtr_port_num;
      config_pP->ovs_config.internal_sampling_port_num =
          internal_sampling_port_num;
      config_pP->ovs_config.internal_sampling_fwd_tbl_num =
          internal_sampling_fwd_tbl_num;
      config_pP->ovs_config.uplink_port_num = uplink_port_num;
      config_pP->ovs_config.uplink_mac      = bfromcstr(uplink_mac);

      if (strcasecmp(pipelined_config_enabled, "false") == 0) {
        config_pP->ovs_config.pipelined_config_enabled = false;
      } else {
        config_pP->ovs_config.pipelined_config_enabled = true;
      }
      OAILOG_INFO(
          LOG_SPGW_APP, "Pipelined config enable: %s\n",
          pipelined_config_enabled);

      if (strcasecmp(multi_tunnel, "false") == 0) {
        config_pP->ovs_config.multi_tunnel = false;
      } else {
        config_pP->ovs_config.multi_tunnel = true;
      }
      OAILOG_INFO(LOG_SPGW_APP, "Multi tunnel enable: %s\n", multi_tunnel);
    } else {
      AssertFatal(false, "Couldn't find all ovs settings in spgw config\n");
    }
#endif
  }
  config_destroy(&cfg);
  return RETURNok;
}

//------------------------------------------------------------------------------
void sgw_config_display(sgw_config_t* config_p) {
  OAILOG_INFO(
      LOG_SPGW_APP, "==== EURECOM %s v%s ====\n", PACKAGE_NAME,
      PACKAGE_VERSION);
  OAILOG_INFO(LOG_SPGW_APP, "Configuration:\n");
  OAILOG_INFO(
      LOG_SPGW_APP, "- File .................................: %s\n",
      bdata(config_p->config_file));

  OAILOG_INFO(LOG_SPGW_APP, "- S1-U:\n");
  OAILOG_INFO(
      LOG_SPGW_APP, "    port number ......: %d\n",
      config_p->udp_port_S1u_S12_S4_up);
  OAILOG_INFO(
      LOG_SPGW_APP, "    S1u_S12_S4 iface .....: %s\n",
      bdata(config_p->ipv4.if_name_S1u_S12_S4_up));
  OAILOG_INFO(
      LOG_SPGW_APP, "    S1u_S12_S4 ip ........: %s/%u\n",
      inet_ntoa(config_p->ipv4.S1u_S12_S4_up),
      config_p->ipv4.netmask_S1u_S12_S4_up);
  OAILOG_INFO(LOG_SPGW_APP, "- S5-S8:\n");
  OAILOG_INFO(
      LOG_SPGW_APP, "    S5_S8 iface ..........: %s\n",
      bdata(config_p->ipv4.if_name_S5_S8_up));
  OAILOG_INFO(
      LOG_SPGW_APP, "    S5_S8 ip .............: %s/%u\n",
      inet_ntoa(config_p->ipv4.S5_S8_up), config_p->ipv4.netmask_S5_S8_up);
  OAILOG_INFO(LOG_SPGW_APP, "- S11:\n");
  OAILOG_INFO(
      LOG_SPGW_APP, "    S11 iface ............: %s\n",
      bdata(config_p->ipv4.if_name_S11));
  OAILOG_INFO(
      LOG_SPGW_APP, "    S11 ip ...............: %s/%u\n",
      inet_ntoa(config_p->ipv4.S11), config_p->ipv4.netmask_S11);
  OAILOG_INFO(LOG_SPGW_APP, "- ITTI:\n");
  OAILOG_INFO(
      LOG_SPGW_APP, "    queue size .......: %u (bytes)\n",
      config_p->itti_config.queue_size);
  OAILOG_INFO(
      LOG_SPGW_APP, "    log file .........: %s\n",
      bdata(config_p->itti_config.log_file));
#if (!EMBEDDED_SGW)
  OAILOG_INFO(LOG_SPGW_APP, "- Logging:\n");
  OAILOG_INFO(
      LOG_SPGW_APP, "    Output ..............: %s\n",
      bdata(config_p->log_config.output));
  OAILOG_INFO(
      LOG_SPGW_APP, "    Output thread-safe...: %s\n",
      (config_p->log_config.is_output_thread_safe) ? "true" : "false");
  OAILOG_INFO(
      LOG_SPGW_APP, "    UDP log level........: %s\n",
      OAILOG_LEVEL_INT2STR(config_p->log_config.udp_log_level));
  OAILOG_INFO(
      LOG_SPGW_APP, "    GTPV1-U log level....: %s\n",
      OAILOG_LEVEL_INT2STR(config_p->log_config.gtpv1u_log_level));
  OAILOG_INFO(
      LOG_SPGW_APP, "    GTPV2-C log level....: %s\n",
      OAILOG_LEVEL_INT2STR(config_p->log_config.gtpv2c_log_level));
  OAILOG_INFO(
      LOG_SPGW_APP, "    S/P-GW APP log level.: %s\n",
      OAILOG_LEVEL_INT2STR(config_p->log_config.spgw_app_log_level));
  OAILOG_INFO(
      LOG_SPGW_APP, "    S11 log level........: %s\n",
      OAILOG_LEVEL_INT2STR(config_p->log_config.s11_log_level));
  OAILOG_INFO(
      LOG_SPGW_APP, "    UTIL log level.......: %s\n",
      OAILOG_LEVEL_INT2STR(config_p->log_config.util_log_level));
  OAILOG_INFO(
      LOG_SPGW_APP, "    ITTI log level.......: %s (InTer-Task Interface)\n",
      OAILOG_LEVEL_INT2STR(config_p->log_config.itti_log_level));
#endif
}
