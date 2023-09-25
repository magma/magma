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

/*! \file spgw_config.cpp
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

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/common/assertions.h"
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/include/sgw_config.h"

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

static bool parse_bool(const char* str);

//------------------------------------------------------------------------------
void sgw_config_init(sgw_config_t* config_pP) {
  memset(config_pP, 0, sizeof(*config_pP));
  pthread_rwlock_init(&config_pP->rw_lock, NULL);
}
//------------------------------------------------------------------------------
status_code_e sgw_config_process(sgw_config_t* config_pP) {
  status_code_e ret = RETURNok;
  return ret;
}

//------------------------------------------------------------------------------
status_code_e sgw_config_parse_string(const char* config_string,
                                      sgw_config_t* config_pP)

{
  config_t cfg = {0};
  config_setting_t* setting_sgw = NULL;
  char* sgw_if_name_S1u_S12_S4_up = NULL;
  char* sgw_if_name_S1u_S12_S4_up_v6 = NULL;
  char* S1u_S12_S4_up = NULL;
  char* S1u_S12_S4_up_v6 = NULL;
  char* sgw_if_name_S5_S8_up = NULL;
  char* S5_S8_up = NULL;
  char* sgw_if_name_S11 = NULL;
  char* S11 = NULL;
  char* s1_ipv6_enabled = NULL;
  libconfig_int sgw_udp_port_S1u_S12_S4_up = 2152;
  libconfig_int sgw_udp_port_S1u_S12_S4_up_v6 = 2152;
  config_setting_t* subsetting = NULL;
  bstring address = NULL;
  bstring cidr = NULL;
  bstring mask = NULL;
  struct in_addr in_addr_var = {0};
  (void)in_addr_var;
  struct in6_addr in6_addr_var = {0};
  (void)in6_addr_var;

  config_init(&cfg);

  /*
   * Read the file. If there is an error, report it and exit.
   */
  if (!config_read_string(&cfg, config_string)) {
    Fatal("\n\n\n\\n\n\nFailed %s\n\n\n\n\n\n\n\n", config_error_text(&cfg));
    OAILOG_CRITICAL(LOG_CONFIG,
                    "Failed to parse SGW configuration file: %s:%d - %s\n",
                    bdata(config_pP->config_file), config_error_line(&cfg),
                    config_error_text(&cfg));
    config_destroy(&cfg);
    Fatal("Failed to parse SGW configuration file %s!\n",
          bdata(config_pP->config_file));
  }

  OAILOG_INFO(LOG_SPGW_APP, "Parsing configuration file provided %s\n",
              bdata(config_pP->config_file));
  setting_sgw = config_lookup(&cfg, SGW_CONFIG_STRING_SGW_CONFIG);
  if (setting_sgw) {
    subsetting = config_setting_get_member(
        setting_sgw, SGW_CONFIG_STRING_NETWORK_INTERFACES_CONFIG);

    if (subsetting) {
      if ((config_setting_lookup_string(
               subsetting,
               SGW_CONFIG_STRING_SGW_INTERFACE_NAME_FOR_S1U_S12_S4_UP,
               (const char**)&sgw_if_name_S1u_S12_S4_up) &&
           config_setting_lookup_string(
               subsetting, SGW_CONFIG_STRING_SGW_IPV4_ADDRESS_FOR_S1U_S12_S4_UP,
               (const char**)&S1u_S12_S4_up) &&
           config_setting_lookup_string(
               subsetting, SGW_CONFIG_STRING_SGW_INTERFACE_NAME_FOR_S5_S8_UP,
               (const char**)&sgw_if_name_S5_S8_up) &&
           config_setting_lookup_string(
               subsetting, SGW_CONFIG_STRING_SGW_IPV4_ADDRESS_FOR_S5_S8_UP,
               (const char**)&S5_S8_up) &&
           config_setting_lookup_string(
               subsetting, SGW_CONFIG_STRING_SGW_INTERFACE_NAME_FOR_S11,
               (const char**)&sgw_if_name_S11) &&
           config_setting_lookup_string(
               subsetting, SGW_CONFIG_STRING_SGW_IPV4_ADDRESS_FOR_S11,
               (const char**)&S11))) {
        config_pP->ipv4.if_name_S1u_S12_S4_up =
            bfromcstr(sgw_if_name_S1u_S12_S4_up);
        cidr = bfromcstr(S1u_S12_S4_up);
        struct bstrList* list = bsplit(cidr, '/');
        AssertFatal(2 == list->qty, "Bad CIDR address %s", bdata(cidr));
        address = list->entry[0];
        mask = list->entry[1];
        IPV4_STR_ADDR_TO_INADDR(bdata(address), config_pP->ipv4.S1u_S12_S4_up,
                                "BAD IP ADDRESS FORMAT FOR S1u_S12_S4 !\n");
        config_pP->ipv4.netmask_S1u_S12_S4_up = atoi((const char*)mask->data);
        bstrListDestroy(list);
        in_addr_var.s_addr = config_pP->ipv4.S1u_S12_S4_up.s_addr;
        OAILOG_INFO(
            LOG_SPGW_APP,
            "Parsing configuration file found S1u_S12_S4_up: %s/%d on %s\n",
            inet_ntoa(in_addr_var), config_pP->ipv4.netmask_S1u_S12_S4_up,
            bdata(config_pP->ipv4.if_name_S1u_S12_S4_up));
        bdestroy(cidr);
        config_pP->ipv4.if_name_S5_S8_up = bfromcstr(sgw_if_name_S5_S8_up);
        cidr = bfromcstr(S5_S8_up);
        list = bsplit(cidr, '/');
        AssertFatal(2 == list->qty, "Bad CIDR address %s", bdata(cidr));
        address = list->entry[0];
        mask = list->entry[1];
        IPV4_STR_ADDR_TO_INADDR(bdata(address), config_pP->ipv4.S5_S8_up,
                                "BAD IP ADDRESS FORMAT FOR S5_S8 !\n");
        config_pP->ipv4.netmask_S5_S8_up = atoi((const char*)mask->data);
        bstrListDestroy(list);
        in_addr_var.s_addr = config_pP->ipv4.S5_S8_up.s_addr;
        OAILOG_INFO(LOG_SPGW_APP,
                    "Parsing configuration file found S5_S8_up: %s/%d on %s\n",
                    inet_ntoa(in_addr_var), config_pP->ipv4.netmask_S5_S8_up,
                    bdata(config_pP->ipv4.if_name_S5_S8_up));

        bdestroy(cidr);
        config_pP->ipv4.if_name_S11 = bfromcstr(sgw_if_name_S11);
        cidr = bfromcstr(S11);
        list = bsplit(cidr, '/');
        AssertFatal(2 == list->qty, "Bad CIDR address %s", bdata(cidr));
        address = list->entry[0];
        mask = list->entry[1];
        IPV4_STR_ADDR_TO_INADDR(bdata(address), config_pP->ipv4.S11,
                                "BAD IP ADDRESS FORMAT FOR S11 !\n");
        config_pP->ipv4.netmask_S11 = atoi((const char*)mask->data);
        bstrListDestroy(list);
        in_addr_var.s_addr = config_pP->ipv4.S11.s_addr;
        OAILOG_INFO(LOG_SPGW_APP,
                    "Parsing configuration file found S11: %s/%d on %s\n",
                    inet_ntoa(in_addr_var), config_pP->ipv4.netmask_S11,
                    bdata(config_pP->ipv4.if_name_S11));
        bdestroy(cidr);
      }

      if (config_setting_lookup_string(
              subsetting,
              SGW_CONFIG_STRING_SGW_INTERFACE_NAME_FOR_S1U_S12_S4_UP,
              (const char**)&sgw_if_name_S1u_S12_S4_up_v6) &&
          config_setting_lookup_string(subsetting,
                                       SGW_CONFIG_STRING_S1_IPV6_ENABLED,
                                       (const char**)&s1_ipv6_enabled) &&
          config_setting_lookup_string(
              subsetting, SGW_CONFIG_STRING_SGW_IPV6_ADDRESS_FOR_S1U_S12_S4_UP,
              (const char**)&S1u_S12_S4_up_v6)) {
        // S1AP IPv6 address
        config_pP->ipv6.if_name_S1u_S12_S4_up =
            bfromcstr(sgw_if_name_S1u_S12_S4_up_v6);
        address = bfromcstr(S1u_S12_S4_up_v6);
        IPV6_STR_ADDR_TO_INADDR(bdata(address), config_pP->ipv6.S1u_S12_S4_up,
                                "BAD IPv6 ADDRESS FORMAT FOR S1u_S12_S4 !\n");
        memcpy(&in6_addr_var, &config_pP->ipv6.S1u_S12_S4_up,
               sizeof(in6_addr_var));
        bdestroy(address);
        char buf[INET6_ADDRSTRLEN];
        OAILOG_INFO(
            LOG_SPGW_APP,
            "Parsing configuration file found S1u_S12_S4_up: %s on %s\n",
            inet_ntop(AF_INET6, &in6_addr_var, buf, INET6_ADDRSTRLEN),
            bdata(config_pP->ipv6.if_name_S1u_S12_S4_up));

        config_pP->ipv6.s1_ipv6_enabled = parse_bool(s1_ipv6_enabled);
      }

      if (config_setting_lookup_int(
              subsetting, SGW_CONFIG_STRING_SGW_PORT_FOR_S1U_S12_S4_UP,
              &sgw_udp_port_S1u_S12_S4_up)) {
        config_pP->udp_port_S1u_S12_S4_up = sgw_udp_port_S1u_S12_S4_up;
      } else {
        config_pP->udp_port_S1u_S12_S4_up = sgw_udp_port_S1u_S12_S4_up;
      }

      if (config_setting_lookup_int(
              subsetting, SGW_CONFIG_STRING_SGW_V6_PORT_FOR_S1U_S12_S4_UP,
              &sgw_udp_port_S1u_S12_S4_up)) {
        config_pP->udp_port_S1u_S12_S4_up_v6 = sgw_udp_port_S1u_S12_S4_up_v6;
      } else {
        config_pP->udp_port_S1u_S12_S4_up_v6 = sgw_udp_port_S1u_S12_S4_up_v6;
      }
    }
    config_setting_t* ovs_settings =
        config_setting_get_member(setting_sgw, SGW_CONFIG_STRING_OVS_CONFIG);
    if (ovs_settings == NULL) {
      Fatal("Couldn't find OVS subsetting in spgw config\n");
    }
    char* ovs_bridge_name = NULL;
    libconfig_int gtp_port_num = 0;
    libconfig_int mtr_port_num = 0;
    libconfig_int internal_sampling_port_num = 0;
    libconfig_int internal_sampling_fwd_tbl_num = 0;
    libconfig_int uplink_port_num = 0;
    char* multi_tunnel = NULL;
    char* agw_l3_tunnel = NULL;
    char* gtp_echo = NULL;
    char* gtp_csum = NULL;

    char* uplink_mac = NULL;
    char* pipelined_managed_tbl0 = NULL;
    char* ebpf_enabled = NULL;
    if (config_setting_lookup_string(ovs_settings,
                                     SGW_CONFIG_STRING_OVS_BRIDGE_NAME,
                                     (const char**)&ovs_bridge_name) &&
        config_setting_lookup_int(
            ovs_settings, SGW_CONFIG_STRING_OVS_GTP_PORT_NUM, &gtp_port_num) &&
        config_setting_lookup_int(ovs_settings,
                                  SGW_CONFIG_STRING_OVS_UPLINK_PORT_NUM,
                                  &uplink_port_num) &&
        config_setting_lookup_string(ovs_settings,
                                     SGW_CONFIG_STRING_OVS_UPLINK_MAC,
                                     (const char**)&uplink_mac) &&
        config_setting_lookup_int(
            ovs_settings, SGW_CONFIG_STRING_OVS_MTR_PORT_NUM, &mtr_port_num) &&
        config_setting_lookup_int(
            ovs_settings, SGW_CONFIG_STRING_OVS_INTERNAL_SAMPLING_PORT_NUM,
            &internal_sampling_port_num) &&
        config_setting_lookup_int(
            ovs_settings, SGW_CONFIG_STRING_OVS_INTERNAL_SAMPLING_FWD_TBL_NUM,
            &internal_sampling_fwd_tbl_num) &&
        config_setting_lookup_string(ovs_settings,
                                     SGW_CONFIG_STRING_OVS_MULTI_TUNNEL,
                                     (const char**)&multi_tunnel) &&
        config_setting_lookup_string(ovs_settings,
                                     SGW_CONFIG_STRING_OVS_GTP_ECHO,
                                     (const char**)&gtp_echo) &&
        config_setting_lookup_string(ovs_settings,
                                     SGW_CONFIG_STRING_OVS_GTP_CHECKSUM,
                                     (const char**)&gtp_csum) &&
        config_setting_lookup_string(ovs_settings,
                                     SGW_CONFIG_STRING_AGW_L3_TUNNEL,
                                     (const char**)&agw_l3_tunnel) &&
        config_setting_lookup_string(ovs_settings,
                                     SGW_CONFIG_STRING_EBPF_ENABLED,
                                     (const char**)&ebpf_enabled) &&
        config_setting_lookup_string(
            ovs_settings, SGW_CONFIG_STRING_OVS_PIPELINED_CONFIG_ENABLED,
            (const char**)&pipelined_managed_tbl0)) {
      config_pP->ovs_config.bridge_name = bfromcstr(ovs_bridge_name);
      config_pP->ovs_config.gtp_port_num = gtp_port_num;
      config_pP->ovs_config.mtr_port_num = mtr_port_num;
      config_pP->ovs_config.internal_sampling_port_num =
          internal_sampling_port_num;
      config_pP->ovs_config.internal_sampling_fwd_tbl_num =
          internal_sampling_fwd_tbl_num;
      config_pP->ovs_config.uplink_port_num = uplink_port_num;
      config_pP->ovs_config.uplink_mac = bfromcstr(uplink_mac);

      if (strcasecmp(pipelined_managed_tbl0, "false") == 0) {
        config_pP->ovs_config.pipelined_managed_tbl0 = false;
      } else {
        config_pP->ovs_config.pipelined_managed_tbl0 = true;
      }
      OAILOG_INFO(LOG_SPGW_APP, "Pipelined config enable: %s\n",
                  pipelined_managed_tbl0);

      if (strcasecmp(multi_tunnel, "false") == 0) {
        config_pP->ovs_config.multi_tunnel = false;
      } else {
        config_pP->ovs_config.multi_tunnel = true;
      }
      OAILOG_INFO(LOG_SPGW_APP, "Multi tunnel enable: %s\n", multi_tunnel);
      if (strcasecmp(gtp_echo, "true") == 0) {
        config_pP->ovs_config.gtp_echo = true;
      } else {
        config_pP->ovs_config.gtp_echo = false;
      }
      OAILOG_INFO(LOG_SPGW_APP, "GTP-U echo response enable: %s\n", gtp_echo);
      if (strcasecmp(gtp_csum, "true") == 0) {
        config_pP->ovs_config.gtp_csum = true;
      } else {
        config_pP->ovs_config.gtp_csum = false;
      }
      OAILOG_INFO(LOG_SPGW_APP, "GTP-U checksum enable: %s\n", gtp_csum);

      if (strcasecmp(agw_l3_tunnel, "true") == 0) {
        config_pP->agw_l3_tunnel = true;
      } else {
        config_pP->agw_l3_tunnel = false;
      }
      OAILOG_INFO(LOG_SPGW_APP, "AGW L3 tunneling enable: %s\n", agw_l3_tunnel);

      if (strcasecmp(ebpf_enabled, "true") == 0) {
        config_pP->ebpf_enabled = true;
      } else {
        config_pP->ebpf_enabled = false;
      }
      OAILOG_INFO(LOG_SPGW_APP, "eBPF enabled: %s\n", ebpf_enabled);
    } else {
      Fatal("Couldn't find all ovs settings in spgw config\n");
    }
  }
  config_destroy(&cfg);
  return RETURNok;
}

status_code_e sgw_config_parse_file(sgw_config_t* config_pP) {
  FILE* fp = NULL;
  status_code_e ret_code = RETURNerror;
  fp = fopen(bdata(config_pP->config_file), "r");
  if (fp == NULL) {
    OAILOG_CRITICAL(LOG_CONFIG,
                    "Failed to open SGW configuration file at path: %s\n",
                    bdata(config_pP->config_file));
    Fatal("Failed to open SGW configuration file at path: %s\n",
          bdata(config_pP->config_file));
  }

  bstring buff = bread((bNread)fread, fp);
  if (buff == NULL) {
    fclose(fp);
    OAILOG_CRITICAL(LOG_CONFIG,
                    "Failed to read SGW configuration file at path: %s\n",
                    bdata(config_pP->config_file));
    Fatal("Failed to read SGW configuration file at path: %s:\n",
          bdata(config_pP->config_file));
  }
  ret_code = sgw_config_parse_string(bdata(buff), config_pP);
  bdestroy_wrapper(&buff);
  return ret_code;
}

//------------------------------------------------------------------------------
void sgw_config_display(sgw_config_t* config_p) {
  OAILOG_INFO(LOG_SPGW_APP, "==== EURECOM %s v%s ====\n", PACKAGE_NAME,
              PACKAGE_VERSION);
  OAILOG_INFO(LOG_SPGW_APP, "Configuration:\n");
  OAILOG_INFO(LOG_SPGW_APP, "- File .................................: %s\n",
              bdata(config_p->config_file));

  OAILOG_INFO(LOG_SPGW_APP, "- S1-U:\n");
  OAILOG_INFO(LOG_SPGW_APP, "    port number ......: %d\n",
              config_p->udp_port_S1u_S12_S4_up);
  OAILOG_INFO(LOG_SPGW_APP, "    S1u_S12_S4 iface .....: %s\n",
              bdata(config_p->ipv4.if_name_S1u_S12_S4_up));
  OAILOG_INFO(LOG_SPGW_APP, "    S1u_S12_S4 ip ........: %s/%u\n",
              inet_ntoa(config_p->ipv4.S1u_S12_S4_up),
              config_p->ipv4.netmask_S1u_S12_S4_up);
  if (config_p->ipv6.s1_ipv6_enabled) {
    char strv6[INET6_ADDRSTRLEN];
    OAILOG_INFO(LOG_CONFIG, "    S1u_S12_S4 ipv6 ......: %s\n",
                inet_ntop(AF_INET6, &config_p->ipv6.S1u_S12_S4_up, strv6,
                          INET6_ADDRSTRLEN));
  }
  OAILOG_INFO(LOG_SPGW_APP, "- S5-S8:\n");
  OAILOG_INFO(LOG_SPGW_APP, "    S5_S8 iface ..........: %s\n",
              bdata(config_p->ipv4.if_name_S5_S8_up));
  OAILOG_INFO(LOG_SPGW_APP, "    S5_S8 ip .............: %s/%u\n",
              inet_ntoa(config_p->ipv4.S5_S8_up),
              config_p->ipv4.netmask_S5_S8_up);
  OAILOG_INFO(LOG_SPGW_APP, "- S11:\n");
  OAILOG_INFO(LOG_SPGW_APP, "    S11 iface ............: %s\n",
              bdata(config_p->ipv4.if_name_S11));
  OAILOG_INFO(LOG_SPGW_APP, "    S11 ip ...............: %s/%u\n",
              inet_ntoa(config_p->ipv4.S11), config_p->ipv4.netmask_S11);
  OAILOG_INFO(LOG_SPGW_APP, "- ITTI:\n");
  OAILOG_INFO(LOG_SPGW_APP, "    queue size .......: %u (bytes)\n",
              config_p->itti_config.queue_size);
  OAILOG_INFO(LOG_SPGW_APP, "    log file .........: %s\n",
              bdata(config_p->itti_config.log_file));
}

void free_sgw_config(sgw_config_t* sgw_config) {
  bdestroy_wrapper(&sgw_config->config_file);
  bdestroy_wrapper(&sgw_config->ovs_config.bridge_name);
  bdestroy_wrapper(&sgw_config->ovs_config.uplink_mac);
  bdestroy_wrapper(&sgw_config->itti_config.log_file);
  bdestroy_wrapper(&sgw_config->ipv4.if_name_S1u_S12_S4_up);
  bdestroy_wrapper(&sgw_config->ipv4.if_name_S5_S8_up);
  bdestroy_wrapper(&sgw_config->ipv4.if_name_S11);
  bdestroy_wrapper(&sgw_config->ipv6.if_name_S1u_S12_S4_up);
}

static bool parse_bool(const char* str) {
  if (strcasecmp(str, "yes") == 0) return true;
  if (strcasecmp(str, "true") == 0) return true;
  if (strcasecmp(str, "no") == 0) return false;
  if (strcasecmp(str, "false") == 0) return false;
  if (strcasecmp(str, "") == 0) return false;

  Fatal("Error in config file: got \"%s\" but expected bool\n", str);
}
