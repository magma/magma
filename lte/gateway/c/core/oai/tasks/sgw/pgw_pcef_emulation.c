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

/*! \file sgw_handlers.c
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
  */
#define SGW

#include "pgw_pcef_emulation.h"

#include <inttypes.h>
#include <netinet/in.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "assertions.h"
#include "async_system.h"
#include "bstrlib.h"
#include "common_defs.h"
#include "common_types.h"
#include "dynamic_memory_check.h"
#include "hashtable.h"
#include "intertask_interface_types.h"
#include "log.h"
#include "pgw_config.h"
#include "pgw_types.h"
#include "spgw_config.h"

/*
 * Function that adds predefined PCC rules to PGW struct,
 * it returns an error or success code after adding rules.
 */
int pgw_pcef_emulation_init(
    spgw_state_t* state_p, const pgw_config_t* const pgw_config_p) {
  int rc             = RETURNok;
  hashtable_rc_t hrc = HASH_TABLE_OK;

  //--------------------------
  // Predefined PCC rules
  //--------------------------
  pcc_rule_t* pcc_rule;
  // Initializing PCC rules only if PGW state doesn't already contain them
  hrc = hashtable_ts_is_key_exists(
      state_p->deactivated_predefined_pcc_rules, SDF_ID_GBR_VOLTE_40K);
  if (hrc == HASH_TABLE_KEY_NOT_EXISTS) {
    pcc_rule                 = (pcc_rule_t*) calloc(1, sizeof(pcc_rule_t));
    pcc_rule->name           = bfromcstr("VOLTE_40K_PCC_RULE");
    pcc_rule->is_activated   = false;
    pcc_rule->sdf_id         = SDF_ID_GBR_VOLTE_40K;
    pcc_rule->bearer_qos.pci = PRE_EMPTION_CAPABILITY_ENABLED;
    pcc_rule->bearer_qos.pl  = 2;
    pcc_rule->bearer_qos.pvi = PRE_EMPTION_VULNERABILITY_DISABLED;
    pcc_rule->bearer_qos.qci = 1;
    pcc_rule->bearer_qos.gbr.br_ul =
        40;  // kilobits per second (1 kbps = 1000 bps)
    pcc_rule->bearer_qos.gbr.br_dl =
        40;  // kilobits per second (1 kbps = 1000 bps)
    pcc_rule->bearer_qos.mbr.br_ul =
        40;  // kilobits per second (1 kbps = 1000 bps)
    pcc_rule->bearer_qos.mbr.br_dl =
        40;  // kilobits per second (1 kbps = 1000 bps)
    pcc_rule->sdf_template.sdf_filter[0].identifier = PF_ID_VOLTE;
    pcc_rule->sdf_template.sdf_filter[0].spare      = 0;
    pcc_rule->sdf_template.sdf_filter[0].direction =
        TRAFFIC_FLOW_TEMPLATE_BIDIRECTIONAL;
    pcc_rule->sdf_template.sdf_filter[0].eval_precedence = 2;
    pcc_rule->sdf_template.sdf_filter[0].length          = 9;
    pcc_rule->sdf_template.sdf_filter[0].packetfiltercontents.flags =
        TRAFFIC_FLOW_TEMPLATE_IPV4_REMOTE_ADDR_FLAG;
    pcc_rule->sdf_template.sdf_filter[0]
        .packetfiltercontents.ipv4remoteaddr[0]
        .addr = 216;
    pcc_rule->sdf_template.sdf_filter[0]
        .packetfiltercontents.ipv4remoteaddr[1]
        .addr = 58;
    pcc_rule->sdf_template.sdf_filter[0]
        .packetfiltercontents.ipv4remoteaddr[2]
        .addr = 210;
    pcc_rule->sdf_template.sdf_filter[0]
        .packetfiltercontents.ipv4remoteaddr[3]
        .addr = 212;
    pcc_rule->sdf_template.sdf_filter[0]
        .packetfiltercontents.ipv4remoteaddr[0]
        .mask = 255;
    pcc_rule->sdf_template.sdf_filter[0]
        .packetfiltercontents.ipv4remoteaddr[1]
        .mask = 255;
    pcc_rule->sdf_template.sdf_filter[0]
        .packetfiltercontents.ipv4remoteaddr[2]
        .mask = 255;
    pcc_rule->sdf_template.sdf_filter[0]
        .packetfiltercontents.ipv4remoteaddr[3]
        .mask                                       = 255;
    pcc_rule->sdf_template.number_of_packet_filters = 1;
    hrc                                             = hashtable_ts_insert(
        state_p->deactivated_predefined_pcc_rules, pcc_rule->sdf_id, pcc_rule);
    if (HASH_TABLE_OK != hrc) {
      return RETURNerror;
    }
  }

  hrc = hashtable_ts_is_key_exists(
      state_p->deactivated_predefined_pcc_rules, SDF_ID_GBR_VOLTE_64K);
  if (hrc == HASH_TABLE_KEY_NOT_EXISTS) {
    pcc_rule                 = (pcc_rule_t*) calloc(1, sizeof(pcc_rule_t));
    pcc_rule->name           = bfromcstr("VOLTE_64K_PCC_RULE");
    pcc_rule->is_activated   = false;
    pcc_rule->sdf_id         = SDF_ID_GBR_VOLTE_64K;
    pcc_rule->bearer_qos.pci = PRE_EMPTION_CAPABILITY_ENABLED;
    pcc_rule->bearer_qos.pl  = 2;
    pcc_rule->bearer_qos.pvi = PRE_EMPTION_VULNERABILITY_DISABLED;
    pcc_rule->bearer_qos.qci = 1;
    pcc_rule->bearer_qos.gbr.br_ul =
        64;  // kilobits per second (1 kbps = 1000 bps)
    pcc_rule->bearer_qos.gbr.br_dl =
        64;  // kilobits per second (1 kbps = 1000 bps)
    pcc_rule->bearer_qos.mbr.br_ul =
        64;  // kilobits per second (1 kbps = 1000 bps)
    pcc_rule->bearer_qos.mbr.br_dl =
        64;  // kilobits per second (1 kbps = 1000 bps)
    pcc_rule->sdf_template.sdf_filter[0].identifier = PF_ID_VOLTE;
    pcc_rule->sdf_template.sdf_filter[0].spare      = 0;
    pcc_rule->sdf_template.sdf_filter[0].direction =
        TRAFFIC_FLOW_TEMPLATE_BIDIRECTIONAL;
    pcc_rule->sdf_template.sdf_filter[0].eval_precedence = 2;
    pcc_rule->sdf_template.sdf_filter[0].length          = 9;
    pcc_rule->sdf_template.sdf_filter[0].packetfiltercontents.flags =
        TRAFFIC_FLOW_TEMPLATE_IPV4_REMOTE_ADDR_FLAG;
    pcc_rule->sdf_template.sdf_filter[0]
        .packetfiltercontents.ipv4remoteaddr[0]
        .addr = 216;
    pcc_rule->sdf_template.sdf_filter[0]
        .packetfiltercontents.ipv4remoteaddr[1]
        .addr = 58;
    pcc_rule->sdf_template.sdf_filter[0]
        .packetfiltercontents.ipv4remoteaddr[2]
        .addr = 210;
    pcc_rule->sdf_template.sdf_filter[0]
        .packetfiltercontents.ipv4remoteaddr[3]
        .addr = 212;
    pcc_rule->sdf_template.sdf_filter[0]
        .packetfiltercontents.ipv4remoteaddr[0]
        .mask = 255;
    pcc_rule->sdf_template.sdf_filter[0]
        .packetfiltercontents.ipv4remoteaddr[1]
        .mask = 255;
    pcc_rule->sdf_template.sdf_filter[0]
        .packetfiltercontents.ipv4remoteaddr[2]
        .mask = 255;
    pcc_rule->sdf_template.sdf_filter[0]
        .packetfiltercontents.ipv4remoteaddr[3]
        .mask                                       = 255;
    pcc_rule->sdf_template.number_of_packet_filters = 1;
    hrc                                             = hashtable_ts_insert(
        state_p->deactivated_predefined_pcc_rules, pcc_rule->sdf_id, pcc_rule);
    if (HASH_TABLE_OK != hrc) {
      return RETURNerror;
    }
  }

  hrc = hashtable_ts_is_key_exists(
      state_p->deactivated_predefined_pcc_rules, SDF_ID_GBR_VILTE_192K);
  if (hrc == HASH_TABLE_KEY_NOT_EXISTS) {
    pcc_rule                 = (pcc_rule_t*) calloc(1, sizeof(pcc_rule_t));
    pcc_rule->name           = bfromcstr("VILTE_192K_PCC_RULE");
    pcc_rule->is_activated   = false;
    pcc_rule->sdf_id         = SDF_ID_GBR_VILTE_192K;
    pcc_rule->bearer_qos.pci = PRE_EMPTION_CAPABILITY_ENABLED;
    pcc_rule->bearer_qos.pl  = 2;
    pcc_rule->bearer_qos.pvi = PRE_EMPTION_VULNERABILITY_DISABLED;
    pcc_rule->bearer_qos.qci = 2;
    pcc_rule->bearer_qos.gbr.br_ul =
        192;  // kilobits per second (1 kbps = 1000 bps)
    pcc_rule->bearer_qos.gbr.br_dl =
        192;  // kilobits per second (1 kbps = 1000 bps)
    pcc_rule->bearer_qos.mbr.br_ul =
        192;  // kilobits per second (1 kbps = 1000 bps)
    pcc_rule->bearer_qos.mbr.br_dl =
        192;  // kilobits per second (1 kbps = 1000 bps)
    pcc_rule->sdf_template.sdf_filter[0].identifier = PF_ID_VILTE;
    pcc_rule->sdf_template.sdf_filter[0].spare      = 0;
    pcc_rule->sdf_template.sdf_filter[0].direction =
        TRAFFIC_FLOW_TEMPLATE_BIDIRECTIONAL;
    pcc_rule->sdf_template.sdf_filter[0].eval_precedence = 2;
    pcc_rule->sdf_template.sdf_filter[0].length          = 9;
    pcc_rule->sdf_template.sdf_filter[0].packetfiltercontents.flags =
        TRAFFIC_FLOW_TEMPLATE_IPV4_REMOTE_ADDR_FLAG;
    pcc_rule->sdf_template.sdf_filter[0]
        .packetfiltercontents.ipv4remoteaddr[0]
        .addr = 216;
    pcc_rule->sdf_template.sdf_filter[0]
        .packetfiltercontents.ipv4remoteaddr[1]
        .addr = 58;
    pcc_rule->sdf_template.sdf_filter[0]
        .packetfiltercontents.ipv4remoteaddr[2]
        .addr = 210;
    pcc_rule->sdf_template.sdf_filter[0]
        .packetfiltercontents.ipv4remoteaddr[3]
        .addr = 213;
    pcc_rule->sdf_template.sdf_filter[0]
        .packetfiltercontents.ipv4remoteaddr[0]
        .mask = 255;
    pcc_rule->sdf_template.sdf_filter[0]
        .packetfiltercontents.ipv4remoteaddr[1]
        .mask = 255;
    pcc_rule->sdf_template.sdf_filter[0]
        .packetfiltercontents.ipv4remoteaddr[2]
        .mask = 255;
    pcc_rule->sdf_template.sdf_filter[0]
        .packetfiltercontents.ipv4remoteaddr[3]
        .mask                                       = 255;
    pcc_rule->sdf_template.number_of_packet_filters = 1;
    hrc                                             = hashtable_ts_insert(
        state_p->deactivated_predefined_pcc_rules, pcc_rule->sdf_id, pcc_rule);
    if (HASH_TABLE_OK != hrc) {
      return RETURNerror;
    }
  }
  hrc = hashtable_ts_is_key_exists(
      state_p->deactivated_predefined_pcc_rules, SDF_ID_TEST_PING);
  if (hrc == HASH_TABLE_KEY_NOT_EXISTS) {
    pcc_rule                 = (pcc_rule_t*) calloc(1, sizeof(pcc_rule_t));
    pcc_rule->name           = bfromcstr("TEST_PING_PCC_RULE");
    pcc_rule->is_activated   = false;
    pcc_rule->sdf_id         = SDF_ID_TEST_PING;
    pcc_rule->bearer_qos.pci = PRE_EMPTION_CAPABILITY_DISABLED;
    pcc_rule->bearer_qos.pl  = 15;
    pcc_rule->bearer_qos.pvi = PRE_EMPTION_VULNERABILITY_ENABLED;
    pcc_rule->bearer_qos.qci = 7;
    pcc_rule->bearer_qos.gbr.br_ul =
        0;  // kilobits per second (1 kbps = 1000 bps)
    pcc_rule->bearer_qos.gbr.br_dl =
        0;  // kilobits per second (1 kbps = 1000 bps)
    pcc_rule->bearer_qos.mbr.br_ul =
        8;  // kilobits per second (1 kbps = 1000 bps)
    pcc_rule->bearer_qos.mbr.br_dl =
        8;  // kilobits per second (1 kbps = 1000 bps)
    pcc_rule->sdf_template.sdf_filter[0].identifier = PF_ID_PING;
    pcc_rule->sdf_template.sdf_filter[0].spare      = 0;
    pcc_rule->sdf_template.sdf_filter[0].direction =
        TRAFFIC_FLOW_TEMPLATE_BIDIRECTIONAL;
    pcc_rule->sdf_template.sdf_filter[0].eval_precedence = 2;
    pcc_rule->sdf_template.sdf_filter[0].length          = 9;
    pcc_rule->sdf_template.sdf_filter[0].packetfiltercontents.flags =
        TRAFFIC_FLOW_TEMPLATE_PROTOCOL_NEXT_HEADER_FLAG;
    pcc_rule->sdf_template.sdf_filter[0]
        .packetfiltercontents.protocolidentifier_nextheader = IPPROTO_ICMP;
    pcc_rule->sdf_template.number_of_packet_filters         = 1;
    hrc = hashtable_ts_insert(
        state_p->deactivated_predefined_pcc_rules, pcc_rule->sdf_id, pcc_rule);
    if (HASH_TABLE_OK != hrc) {
      return RETURNerror;
    }
  }

  hrc = hashtable_ts_is_key_exists(
      state_p->deactivated_predefined_pcc_rules, SDF_ID_NGBR_DEFAULT);
  if (hrc == HASH_TABLE_KEY_NOT_EXISTS) {
    // really necessary ?
    pcc_rule                 = (pcc_rule_t*) calloc(1, sizeof(pcc_rule_t));
    pcc_rule->name           = bfromcstr("DEFAULT_PCC_RULE");
    pcc_rule->is_activated   = false;
    pcc_rule->sdf_id         = SDF_ID_NGBR_DEFAULT;
    pcc_rule->bearer_qos.pci = PRE_EMPTION_CAPABILITY_DISABLED;
    pcc_rule->bearer_qos.pl  = 15;
    pcc_rule->bearer_qos.pvi = PRE_EMPTION_VULNERABILITY_ENABLED;
    pcc_rule->bearer_qos.qci = 9;
    pcc_rule->bearer_qos.gbr.br_ul =
        0;  // kilobits per second (1 kbps = 1000 bps)
    pcc_rule->bearer_qos.gbr.br_dl =
        0;  // kilobits per second (1 kbps = 1000 bps)
    pcc_rule->bearer_qos.mbr.br_ul =
        1000;  // kilobits per second (1 kbps = 1000 bps)
    pcc_rule->bearer_qos.mbr.br_dl =
        1000;  // kilobits per second (1 kbps = 1000 bps)
    pcc_rule->sdf_template.sdf_filter[0].identifier = PF_ID_DEFAULT;
    pcc_rule->sdf_template.sdf_filter[0].spare      = 0;
    pcc_rule->sdf_template.sdf_filter[0].direction =
        TRAFFIC_FLOW_TEMPLATE_DOWNLINK_ONLY;
    pcc_rule->sdf_template.sdf_filter[0].eval_precedence = 2;
    pcc_rule->sdf_template.sdf_filter[0].length          = 9;
    pcc_rule->sdf_template.sdf_filter[0].packetfiltercontents.flags =
        TRAFFIC_FLOW_TEMPLATE_IPV4_REMOTE_ADDR_FLAG;
    pcc_rule->sdf_template.sdf_filter[0]
        .packetfiltercontents.ipv4remoteaddr[0]
        .addr = (uint8_t)((pgw_config_p->ue_pool_addr[0].s_addr) & 0x000000FF);
    pcc_rule->sdf_template.sdf_filter[0]
        .packetfiltercontents.ipv4remoteaddr[1]
        .addr =
        (uint8_t)((pgw_config_p->ue_pool_addr[0].s_addr >> 8) & 0x000000FF);
    pcc_rule->sdf_template.sdf_filter[0]
        .packetfiltercontents.ipv4remoteaddr[2]
        .addr =
        (uint8_t)((pgw_config_p->ue_pool_addr[0].s_addr >> 16) & 0x000000FF);
    pcc_rule->sdf_template.sdf_filter[0]
        .packetfiltercontents.ipv4remoteaddr[3]
        .addr =
        (uint8_t)((pgw_config_p->ue_pool_addr[0].s_addr >> 24) & 0x000000FF);
    struct in_addr addr_mask = {0};
    addr_mask.s_addr =
        htonl(0xFFFFFFFF << (32 - pgw_config_p->ue_pool_mask[0]));
    pcc_rule->sdf_template.sdf_filter[0]
        .packetfiltercontents.ipv4remoteaddr[0]
        .mask = (uint8_t)((addr_mask.s_addr) & 0x000000FF);
    pcc_rule->sdf_template.sdf_filter[0]
        .packetfiltercontents.ipv4remoteaddr[1]
        .mask = (uint8_t)((addr_mask.s_addr >> 8) & 0x000000FF);
    pcc_rule->sdf_template.sdf_filter[0]
        .packetfiltercontents.ipv4remoteaddr[2]
        .mask = (uint8_t)((addr_mask.s_addr >> 16) & 0x000000FF);
    pcc_rule->sdf_template.sdf_filter[0]
        .packetfiltercontents.ipv4remoteaddr[3]
        .mask = (uint8_t)((addr_mask.s_addr >> 24) & 0x000000FF);
    pcc_rule->sdf_template.number_of_packet_filters = 1;
    hrc                                             = hashtable_ts_insert(
        state_p->deactivated_predefined_pcc_rules, pcc_rule->sdf_id, pcc_rule);
    if (HASH_TABLE_OK != hrc) {
      return RETURNerror;
    }
  }

  for (int i = 0; i < (SDF_ID_MAX - 1); i++) {
    if (pgw_config_p->pcef.preload_static_sdf_identifiers[i]) {
      pgw_pcef_emulation_apply_rule(
          state_p, pgw_config_p->pcef.preload_static_sdf_identifiers[i],
          pgw_config_p);
    } else
      break;
  }

  if (pgw_config_p->pcef.automatic_push_dedicated_bearer_sdf_identifier) {
    pgw_pcef_emulation_apply_rule(
        state_p,
        pgw_config_p->pcef.automatic_push_dedicated_bearer_sdf_identifier,
        pgw_config_p);
  }
  return rc;
}

//------------------------------------------------------------------------------
// may change sdf_id to PCC_rule name ?
void pgw_pcef_emulation_apply_rule(
    spgw_state_t* state_p, const sdf_id_t sdf_id,
    const pgw_config_t* const pgw_config_p) {
  pcc_rule_t* pcc_rule = NULL;
  hashtable_rc_t hrc   = hashtable_ts_get(
      state_p->deactivated_predefined_pcc_rules, sdf_id, (void**) &pcc_rule);

  if (HASH_TABLE_OK == hrc) {
    if (!pcc_rule->is_activated) {
      OAILOG_INFO(LOG_SPGW_APP, "Loading PCC rule %s\n", bdata(pcc_rule->name));
      pcc_rule->is_activated = true;
      for (int sdff_i = 0;
           sdff_i < pcc_rule->sdf_template.number_of_packet_filters; sdff_i++) {
        pgw_pcef_emulation_apply_sdf_filter(
            &pcc_rule->sdf_template.sdf_filter[sdff_i], pcc_rule->sdf_id,
            pgw_config_p);
      }
    }
  }
}

//------------------------------------------------------------------------------
void pgw_pcef_emulation_apply_sdf_filter(
    sdf_filter_t* const sdf_f, const sdf_id_t sdf_id,
    const pgw_config_t* const pgw_config_p) {
  if ((TRAFFIC_FLOW_TEMPLATE_BIDIRECTIONAL == sdf_f->direction) ||
      (TRAFFIC_FLOW_TEMPLATE_DOWNLINK_ONLY == sdf_f->direction)) {
    bstring filter = pgw_pcef_emulation_packet_filter_2_iptable_string(
        &sdf_f->packetfiltercontents, sdf_f->direction);

    bstring marking_command = NULL;
    if ((TRAFFIC_FLOW_TEMPLATE_IPV4_REMOTE_ADDR_FLAG |
         TRAFFIC_FLOW_TEMPLATE_IPV6_REMOTE_ADDR_FLAG) &
        sdf_f->packetfiltercontents.flags) {
      marking_command = bformat(
          "iptables -I POSTROUTING -t mangle  %s -j MARK --set-mark %d",
          bdata(filter), sdf_id);
    } else {
      // marking_command = bformat("iptables -I PREROUTING -t mangle
      // --in-interface %s --dest %"PRIu8".%"PRIu8".%"PRIu8".%"PRIu8"/%"PRIu8"
      // %s -j MARK --set-mark %d",
      marking_command = bformat(
          "iptables -I POSTROUTING -t mangle  --dest %" PRIu8 ".%" PRIu8
          ".%" PRIu8 ".%" PRIu8 "/%" PRIu8 " %s -j MARK --set-mark %d",
          NIPADDR(pgw_config_p->ue_pool_addr[0].s_addr),
          pgw_config_p->ue_pool_mask[0], bdata(filter), sdf_id);
    }
    bdestroy_wrapper(&filter);
    async_system_command(TASK_ASYNC_SYSTEM, false, bdata(marking_command));
    bdestroy_wrapper(&marking_command);

    // for UE <-> PGW traffic
    filter = pgw_pcef_emulation_packet_filter_2_iptable_string(
        &sdf_f->packetfiltercontents, sdf_f->direction);

    marking_command = NULL;
    if ((TRAFFIC_FLOW_TEMPLATE_IPV4_REMOTE_ADDR_FLAG |
         TRAFFIC_FLOW_TEMPLATE_IPV6_REMOTE_ADDR_FLAG) &
        sdf_f->packetfiltercontents.flags) {
      marking_command = bformat(
          "iptables -I OUTPUT -t mangle  %s -j MARK --set-mark %d",
          bdata(filter), sdf_id);
    } else {
      marking_command = bformat(
          "iptables -I OUTPUT -t mangle  --dest %" PRIu8 ".%" PRIu8 ".%" PRIu8
          ".%" PRIu8 "/%" PRIu8 " %s -j MARK --set-mark %d",
          NIPADDR(pgw_config_p->ue_pool_addr[0].s_addr),
          pgw_config_p->ue_pool_mask[0], bdata(filter), sdf_id);
    }
    bdestroy_wrapper(&filter);
    async_system_command(TASK_ASYNC_SYSTEM, false, bdata(marking_command));
    bdestroy_wrapper(&marking_command);
  }
}

//------------------------------------------------------------------------------
bstring pgw_pcef_emulation_packet_filter_2_iptable_string(
    packet_filter_contents_t* const packetfiltercontents, uint8_t direction) {
  bstring bstr = bfromcstralloc(64, " ");

  if ((TRAFFIC_FLOW_TEMPLATE_DOWNLINK_ONLY == direction) ||
      (TRAFFIC_FLOW_TEMPLATE_BIDIRECTIONAL == direction)) {
    if (TRAFFIC_FLOW_TEMPLATE_IPV4_REMOTE_ADDR_FLAG &
        packetfiltercontents->flags) {
      bformata(
          bstr, "  --destination %d.%d.%d.%d/%d.%d.%d.%d",
          packetfiltercontents->ipv4remoteaddr[0].addr,
          packetfiltercontents->ipv4remoteaddr[1].addr,
          packetfiltercontents->ipv4remoteaddr[2].addr,
          packetfiltercontents->ipv4remoteaddr[3].addr,
          packetfiltercontents->ipv4remoteaddr[0].mask,
          packetfiltercontents->ipv4remoteaddr[1].mask,
          packetfiltercontents->ipv4remoteaddr[2].mask,
          packetfiltercontents->ipv4remoteaddr[3].mask);
    } else {
      bformata(
          bstr, " --source %d.%d.%d.%d/%d.%d.%d.%d",
          packetfiltercontents->ipv4remoteaddr[0].addr,
          packetfiltercontents->ipv4remoteaddr[1].addr,
          packetfiltercontents->ipv4remoteaddr[2].addr,
          packetfiltercontents->ipv4remoteaddr[3].addr,
          packetfiltercontents->ipv4remoteaddr[0].mask,
          packetfiltercontents->ipv4remoteaddr[1].mask,
          packetfiltercontents->ipv4remoteaddr[2].mask,
          packetfiltercontents->ipv4remoteaddr[3].mask);
    }
  }
  if (TRAFFIC_FLOW_TEMPLATE_IPV6_REMOTE_ADDR_FLAG &
      packetfiltercontents->flags) {
    Fatal("TODO Implement pgw_pcef_emulation_packet_filter_2_iptable_string");
  }
  if (TRAFFIC_FLOW_TEMPLATE_PROTOCOL_NEXT_HEADER_FLAG &
      packetfiltercontents->flags) {
    bformata(
        bstr, " --protocol %u",
        packetfiltercontents->protocolidentifier_nextheader);
  }
  if (TRAFFIC_FLOW_TEMPLATE_SINGLE_LOCAL_PORT_FLAG &
      packetfiltercontents->flags) {
    if ((TRAFFIC_FLOW_TEMPLATE_DOWNLINK_ONLY == direction) ||
        (TRAFFIC_FLOW_TEMPLATE_BIDIRECTIONAL == direction)) {
      bformata(
          bstr, " --destination-port %" PRIu16 " ",
          packetfiltercontents->singlelocalport);
    } else if (TRAFFIC_FLOW_TEMPLATE_UPLINK_ONLY == direction) {
      bformata(
          bstr, " --source-port %" PRIu16 " ",
          packetfiltercontents->singlelocalport);
    }
  }
  if (TRAFFIC_FLOW_TEMPLATE_LOCAL_PORT_RANGE_FLAG &
      packetfiltercontents->flags) {
    Fatal("TODO LOCAL_PORT_RANGE");
  }
  if (TRAFFIC_FLOW_TEMPLATE_SINGLE_REMOTE_PORT_FLAG &
      packetfiltercontents->flags) {
    if ((TRAFFIC_FLOW_TEMPLATE_DOWNLINK_ONLY == direction) ||
        (TRAFFIC_FLOW_TEMPLATE_BIDIRECTIONAL == direction)) {
      bformata(
          bstr, " --source-port %" PRIu16 " ",
          packetfiltercontents->singleremoteport);
    } else if (TRAFFIC_FLOW_TEMPLATE_UPLINK_ONLY == direction) {
      bformata(
          bstr, " --destination-port %" PRIu16 " ",
          packetfiltercontents->singleremoteport);
    }
  }
  if (TRAFFIC_FLOW_TEMPLATE_REMOTE_PORT_RANGE_FLAG &
      packetfiltercontents->flags) {
    Fatal("TODO REMOTE_PORT_RANGE");
  }
  if (TRAFFIC_FLOW_TEMPLATE_SECURITY_PARAMETER_INDEX_FLAG &
      packetfiltercontents->flags) {
    bformata(
        bstr, " -m esp --espspi %" PRIu32 " ",
        packetfiltercontents->securityparameterindex);
  }
  if (TRAFFIC_FLOW_TEMPLATE_TYPE_OF_SERVICE_TRAFFIC_CLASS_FLAG &
      packetfiltercontents->flags) {
    // TODO mask
    bformata(
        bstr, " -m tos --tos 0x%02X",
        packetfiltercontents->typdeofservice_trafficclass.value);
  }
  if (TRAFFIC_FLOW_TEMPLATE_FLOW_LABEL_FLAG & packetfiltercontents->flags) {
    Fatal("TODO Implement pgw_pcef_emulation_packet_filter_2_iptable_string");
  }
  return bstr;
}

//------------------------------------------------------------------------------
int pgw_pcef_get_sdf_parameters(
    spgw_state_t* state_p, const sdf_id_t sdf_id,
    bearer_qos_t* const bearer_qos, packet_filter_t* const packet_filter,
    uint8_t* const num_pf) {
  pcc_rule_t* pcc_rule = NULL;
  hashtable_rc_t hrc   = hashtable_ts_get(
      state_p->deactivated_predefined_pcc_rules, sdf_id, (void**) &pcc_rule);

  if (HASH_TABLE_OK == hrc) {
    if (pcc_rule->is_activated) {
      memcpy(bearer_qos, &pcc_rule->bearer_qos, sizeof(pcc_rule->bearer_qos));
      memcpy(
          packet_filter, &pcc_rule->sdf_template.sdf_filter,
          sizeof(pcc_rule->sdf_template.sdf_filter[0]) *
              pcc_rule->sdf_template.number_of_packet_filters);
      *num_pf = pcc_rule->sdf_template.number_of_packet_filters;
      return RETURNok;
    }
  }
  memset(bearer_qos, 0, sizeof(*bearer_qos));
  memset(packet_filter, 0, sizeof(*packet_filter));
  *num_pf = 0;
  return RETURNerror;
}
