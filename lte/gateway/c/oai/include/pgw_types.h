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

/*! \file pgw_pcef_emulation.h
 * \brief
 * \author Lionel Gauthier
 * \company Eurecom
 * \email: lionel.gauthier@eurecom.fr
 */

#pragma once

#include "3gpp_29.274.h"
#include "3gpp_24.008.h"
#include "3gpp_24.007.h"
#include "queue.h"

// Each service data flow template may contain any number of service data
// flow filters;
#define SERVICE_DATA_FLOW_TEMPLATE_NB_PACKET_FILTERS_MAX 8

typedef struct packet_filter_s sdf_filter_t;

typedef enum {
  PF_ID_MIN = 0,
  PF_ID_VOLTE,
  PF_ID_VILTE,
  PF_ID_VIDEO_CONVERSATIONAL,
  PF_ID_VIDEO_BUFFERED_STREAMING,
  PF_ID_SEC_WEB_BROWSING,
  PF_ID_WEB_BROWSING,
  PF_ID_DNS,
  PF_ID_PING,
  PF_ID_DEFAULT,
  PF_ID_MAX
} pf_id_t;

typedef enum {
  SDF_ID_MIN = (EPS_BEARER_IDENTITY_LAST + 1),
  SDF_ID_GBR_VOLTE_16K,
  SDF_ID_GBR_VOLTE_24K,
  SDF_ID_GBR_VOLTE_40K,
  SDF_ID_GBR_VOLTE_64K,
  SDF_ID_GBR_VILTE_192K,
  SDF_ID_GBR_VILTE_384K,
  SDF_ID_GBR_VILTE_768K,
  SDF_ID_GBR_VILTE_2M,
  SDF_ID_GBR_VILTE_4M,
  SDF_ID_GBR_NON_CONVERSATIONAL_VIDEO_256K,
  SDF_ID_GBR_NON_CONVERSATIONAL_VIDEO_512K,
  SDF_ID_GBR_NON_CONVERSATIONAL_VIDEO_1M,
  SDF_ID_NGBR_IMS_SIGNALLING,
  SDF_ID_NGBR_DEFAULT_PREMIUM,
  SDF_ID_NGBR_DEFAULT,
  SDF_ID_TEST_PING,
  SDF_ID_MAX
} sdf_id_t;

typedef struct sdf_template_s {
  uint8_t number_of_packet_filters;
  sdf_filter_t sdf_filter[SERVICE_DATA_FLOW_TEMPLATE_NB_PACKET_FILTERS_MAX];
  // Our understanding is the following: For non GBR different SDF
  // filters can map to same SDF if they have the same QCI and ARP.
  // SDFs (or aggregation of SDFs) with the same QCI and ARP can
  // be delivered through the same EPS bearer.
} sdf_template_t;

// Each PCC rule contains a service data flow template, which defines the data
// for the service data flow detection
typedef struct pcc_rule_s {
  bstring name;
  bool is_activated;
  sdf_id_t sdf_id;
  sdf_template_t sdf_template;
  bearer_qos_t bearer_qos;
  uint32_t precedence;
  STAILQ_ENTRY(pcc_rule_s) entries;
} pcc_rule_t;

typedef struct conf_ipv4_list_elm_s {
  STAILQ_ENTRY(conf_ipv4_list_elm_s) ipv4_entries;
  struct in_addr addr;
} conf_ipv4_list_elm_t;
