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

/*! \file pgw_config.h
 * \brief
 * \author Lionel Gauthier
 * \company Eurecom
 * \email: lionel.gauthier@eurecom.fr
 */

#ifndef FILE_PGW_CONFIG_SEEN
#define FILE_PGW_CONFIG_SEEN

#include <sys/socket.h>  // inet_aton
#include <netinet/in.h>  // inet_aton
#include <arpa/inet.h>   // inet_aton
#include <pthread.h>
#include <stdbool.h>
#include <stdint.h>

#include "queue.h"
#include "bstrlib.h"
#include "pgw_types.h"

#define PGW_CONFIG_STRING_PGW_CONFIG "P-GW"
#define PGW_CONFIG_STRING_NETWORK_INTERFACES_CONFIG "NETWORK_INTERFACES"
#define PGW_CONFIG_STRING_PGW_INTERFACE_NAME_FOR_S5_S8                         \
  "PGW_INTERFACE_NAME_FOR_S5_S8"
#define PGW_CONFIG_STRING_PGW_INTERFACE_NAME_FOR_SGI                           \
  "PGW_INTERFACE_NAME_FOR_SGI"
#define PGW_CONFIG_STRING_PGW_IPV4_ADDR_FOR_SGI "PGW_IPV4_ADDRESS_FOR_SGI"
#define PGW_CONFIG_STRING_PGW_MASQUERADE_SGI "PGW_MASQUERADE_SGI"
#define PGW_CONFIG_STRING_UE_TCP_MSS_CLAMPING "UE_TCP_MSS_CLAMPING"
#define PGW_CONFIG_STRING_NAS_FORCE_PUSH_PCO                                   \
  "FORCE_PUSH_PROTOCOL_CONFIGURATION_OPTIONS"

#define PGW_CONFIG_STRING_IP_ADDRESS_POOL "IP_ADDRESS_POOL"
#define PGW_CONFIG_STRING_IPV4_ADDRESS_LIST "IPV4_LIST"
#define PGW_CONFIG_STRING_IPV4_PREFIX_DELIMITER '/'
#define PGW_CONFIG_STRING_DEFAULT_DNS_IPV4_ADDRESS "DEFAULT_DNS_IPV4_ADDRESS"
#define PGW_CONFIG_STRING_DEFAULT_DNS_SEC_IPV4_ADDRESS                         \
  "DEFAULT_DNS_SEC_IPV4_ADDRESS"
#define PGW_CONFIG_STRING_UE_MTU "UE_MTU"
#define PGW_CONFIG_STRING_GTPV1U_REALIZATION "GTPV1U_REALIZATION"
#define PGW_CONFIG_STRING_NO_GTP_KERNEL_AVAILABLE "NO_GTP_KERNEL_AVAILABLE"
#define PGW_CONFIG_STRING_GTP_KERNEL_MODULE "GTP_KERNEL_MODULE"
#define PGW_CONFIG_STRING_GTP_KERNEL "GTP_KERNEL"

#define PGW_CONFIG_STRING_INTERFACE_DISABLED "none"

#define PGW_CONFIG_STRING_PCEF "PCEF"
#define PGW_CONFIG_STRING_PCEF_ENABLED "PCEF_ENABLED"
#define PGW_CONFIG_STRING_TRAFFIC_SHAPPING_ENABLED "TRAFFIC_SHAPPING_ENABLED"
#define PGW_CONFIG_STRING_TCP_ECN_ENABLED "TCP_ECN_ENABLED"
#define PGW_CONFIG_STRING_AUTOMATIC_PUSH_DEDICATED_BEARER_PCC_RULE             \
  "AUTOMATIC_PUSH_DEDICATED_BEARER_PCC_RULE"
#define PGW_CONFIG_STRING_DEFAULT_BEARER_STATIC_PCC_RULE                       \
  "DEFAULT_BEARER_STATIC_PCC_RULE"
#define PGW_CONFIG_STRING_PUSH_STATIC_PCC_RULES "PUSH_STATIC_PCC_RULES"
#define PGW_CONFIG_STRING_APN_AMBR_UL "APN_AMBR_UL"
#define PGW_CONFIG_STRING_APN_AMBR_DL "APN_AMBR_DL"
#define PGW_ABORT_ON_ERROR false
#define PGW_WARN_ON_ERROR true
#define PGW_CONFIG_P_CSCF_IPV4_ADDRESS "P_CSCF_IPV4_ADDRESS"
#define PGW_CONFIG_P_CSCF_IPV6_ADDRESS "P_CSCF_IPV6_ADDRESS"
#define PGW_CONFIG_DNS_SERVER_IPV6_ADDRESS "DNS_SERVER_IPV6_ADDRESS"

#define PGW_CONFIG_STRING_NAT_ENABLED "ENABLE_NAT"

// may be more
#define PGW_MAX_ALLOCATED_PDN_ADDRESSES 1024
#define PGW_NUM_UE_POOL_MAX 16

typedef struct pgw_config_s {
  /* Reader/writer lock for this configuration */
  pthread_rwlock_t rw_lock;
  bstring config_file;

  struct {
    bstring if_name_S5_S8;
    struct in_addr S5_S8;
    uint32_t mtu_S5_S8;         // read from system
    struct in_addr addr_S5_S8;  // read from system
    uint8_t mask_S5_S8;         // read from system

    bstring if_name_SGI;
    uint32_t mtu_SGI;         // read from system
    struct in_addr addr_sgi;  // read from system
    uint8_t mask_sgi;         // read from system

    struct in_addr default_dns;
    struct in_addr default_dns_sec;
  } ipv4;

  bool ue_tcp_mss_clamp;  // for UE TCP traffic
  bool masquerade_SGI;

  int num_ue_pool;
  uint8_t ue_pool_mask[PGW_NUM_UE_POOL_MAX];
  struct in_addr ue_pool_addr[PGW_NUM_UE_POOL_MAX];

  bool force_push_pco;
  uint16_t ue_mtu;
  bool enable_nat;

  struct {
    bool enabled;
    bool traffic_shaping_enabled;
    bool tcp_ecn_enabled;  // test for CoDel qdisc
    sdf_id_t default_bearer_sdf_identifier;
    sdf_id_t automatic_push_dedicated_bearer_sdf_identifier;
    sdf_id_t preload_static_sdf_identifiers[SDF_ID_MAX - 1];
    uint64_t apn_ambr_ul;
    uint64_t apn_ambr_dl;
  } pcef;

  struct {
    struct in_addr ipv4_addr;
    struct in6_addr ipv6_addr;
  } pcscf;

  struct {
    struct in6_addr dns_ipv6_addr;
  } ipv6;

  STAILQ_HEAD(ipv4_pool_head_s, conf_ipv4_list_elm_s) ipv4_pool_list;
} pgw_config_t;

int pgw_config_process(pgw_config_t* config_pP);
void pgw_config_init(pgw_config_t* config_pP);
int pgw_config_parse_file(pgw_config_t* config_pP);
void pgw_config_display(pgw_config_t* config_p);
void free_pgw_config(pgw_config_t* pgw_config_p);

#define pgw_config_read_lock(pGWcONFIG)                                        \
  pthread_rwlock_rdlock(&(pGWcONFIG)->rw_lock)
#define pgw_config_write_lock(pGWcONFIG)                                       \
  pthread_rwlock_wrlock(&(pGWcONFIG)->rw_lock)
#define pgw_config_unlock(pGWcONFIG)                                           \
  pthread_rwlock_unlock(&(pGWcONFIG)->rw_lock)

#endif /* FILE_PGW_CONFIG_SEEN */
