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

/*****************************************************************************
Source      mme_api.h

Version     0.1

Date        2013/02/28

Product     NAS stack

Subsystem   Application Programming Interface

Author      Frederic Maurel, Lionel GAUTHIER

Description Implements the API used by the NAS layer running in the MME
        to interact with a Mobility Management Entity

*****************************************************************************/
#ifndef FILE_MME_API_SEEN
#define FILE_MME_API_SEEN

#include <stdbool.h>
#include <stdint.h>

#include "TrackingAreaIdentity.h"
#include "TrackingAreaIdentityList.h"
#include "3gpp_23.003.h"
#include "common_types.h"
#include "3gpp_36.401.h"
#include "bstrlib.h"
#include "mme_config.h"

struct mme_config_s;

/****************************************************************************/
/*********************  G L O B A L    C O N S T A N T S  *******************/
/****************************************************************************/

/* Maximum number of UEs the MME may simultaneously support */
#define MME_API_NB_UE_MAX 256

/* Features supported by the MME */
typedef enum mme_api_feature_s {
  MME_API_NO_FEATURE_SUPPORTED = 0,
  MME_API_UNAUTHENTICATED_IMSI = (1 << 0),
  MME_API_IPV4                 = (1 << 1),
  MME_API_IPV6                 = (1 << 2),
  MME_API_SINGLE_ADDR_BEARERS  = (1 << 3),
  MME_API_SMS_SUPPORTED        = (1 << 4),
  MME_API_CSFB_SMS_SUPPORTED   = (1 << 5),
  MME_API_VOLTE_SUPPORTED      = (1 << 6),
  MME_API_SMS_ORC8R_SUPPORTED  = (1 << 7)
} mme_api_feature_t;

/* Network IP version capability */
typedef enum mme_api_ip_version_e {
  MME_API_IPV4_ADDR,
  MME_API_IPV6_ADDR,
  MME_API_IPV4V6_ADDR,
  MME_API_ADDR_MAX
} mme_api_ip_version_t;

typedef enum {
  UE_UNREGISTERED = 0,
  UE_REGISTERED,
} mm_state_t;

typedef struct gummei_list_s {
  uint8_t num_gummei;
  gummei_t gummei[MAX_GUMMEI];
} gummei_list_t;
/*
 * EPS Mobility Management configuration data
 * ------------------------------------------
 */
typedef struct mme_api_emm_config_s {
  mme_api_feature_t features; /* Supported features           */
  gummei_list_t gummei;       /* EPS Globally Unique MME Identity List*/
  uint8_t prefered_integrity_algorithm[8];  // choice in
                                            // NAS_SECURITY_ALGORITHMS_EIA0, etc
  uint8_t prefered_ciphering_algorithm[8];  // choice in
                                            // NAS_SECURITY_ALGORITHMS_EEA0, etc
  uint8_t eps_network_feature_support[2];
  bool force_push_pco;
  tai_list_t tai_list;
  bstring full_network_name;
  bstring short_network_name;
  uint8_t daylight_saving_time;
} mme_api_emm_config_t;

/*
 * EPS Session Management configuration data
 * -----------------------------------------
 */
typedef struct mme_api_esm_config_s {
  mme_api_feature_t features; /* Supported features           */
  uint8_t dns_prim_ipv4[4];   /* Network byte order */
  uint8_t dns_sec_ipv4[4];    /* Network byte order */
} mme_api_esm_config_t;

/****************************************************************************/
/************************  G L O B A L    T Y P E S  ************************/
/****************************************************************************/

/* EPS subscribed QoS profile  */
typedef struct mme_api_qos_s {
#define MME_API_UPLINK 0
#define MME_API_DOWNLINK 1
#define MME_API_DIRECTION 2
  int gbr[MME_API_DIRECTION]; /* Guaranteed Bit Rate          */
  int mbr[MME_API_DIRECTION]; /* Maximum Bit Rate         */
  int qci;                    /* QoS Class Identifier         */
} mme_api_qos_t;

/****************************************************************************/
/********************  G L O B A L    V A R I A B L E S  ********************/
/****************************************************************************/

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/
struct mme_config_s;

int mme_api_get_emm_config(
    mme_api_emm_config_t* config, const struct mme_config_s* mme_config_p);

#define REMOVE_OLD_CONTEXT true
#define REMOVE_NEW_CONTEXT false
void mme_api_duplicate_enb_ue_s1ap_id_detected(
    const enb_s1ap_id_key_t enb_ue_s1ap_id,
    const mme_ue_s1ap_id_t mme_ue_s1ap_id, const bool is_remove_old);

int mme_api_get_esm_config(mme_api_esm_config_t* config);

int mme_api_notify_imsi(const mme_ue_s1ap_id_t id, const imsi64_t imsi64);

int mme_api_notify_new_guti(const mme_ue_s1ap_id_t ueid, guti_t* const guti);

int mme_api_new_guti(
    const imsi_t* const imsi, const guti_t* const old_guti, guti_t* const guti,
    const tai_t* const originating_tai, tai_list_t* const tai_list);

int mme_api_subscribe(
    bstring* apn, mme_api_ip_version_t mme_pdn_index, bstring* pdn_addr,
    int is_emergency, mme_api_qos_t* qos);
int mme_api_unsubscribe(bstring apn);

void mme_ue_context_update_ue_emm_state(
    mme_ue_s1ap_id_t mme_ue_s1ap_id, mm_state_t new_emm_state);

bool mme_ue_context_get_ue_sgs_vlr_reliable(mme_ue_s1ap_id_t mme_ue_s1ap_id);

void mme_ue_context_update_ue_sgs_vlr_reliable(
    mme_ue_s1ap_id_t mme_ue_s1ap_id, bool vlr_reliable);

bool mme_ue_context_get_ue_sgs_neaf(mme_ue_s1ap_id_t mme_ue_s1ap_id);

void mme_ue_context_update_ue_sgs_neaf(
    mme_ue_s1ap_id_t mme_ue_s1ap_id, bool neaf);

/* Compares the given two PLMNs */
#define IS_PLMN_EQUAL(pLMN1, pLMN2)                                            \
  (((pLMN1.mcc_digit1 == pLMN2.mcc_digit1) &&                                  \
    (pLMN1.mcc_digit2 == pLMN2.mcc_digit2) &&                                  \
    (pLMN1.mcc_digit3 == pLMN2.mcc_digit3) &&                                  \
    (pLMN1.mnc_digit1 == pLMN2.mnc_digit1) &&                                  \
    (pLMN1.mnc_digit2 == pLMN2.mnc_digit2) &&                                  \
    (pLMN1.mnc_digit3 == pLMN2.mnc_digit3)) ?                                  \
       (true) :                                                                \
       (false))

#define COPY_GUMMEI(gUTI, gUMMEIvAl)                                           \
  do {                                                                         \
    gUTI->gummei.mme_gid         = gUMMEIvAl.mme_gid;                          \
    gUTI->gummei.mme_code        = gUMMEIvAl.mme_code;                         \
    gUTI->gummei.plmn.mcc_digit1 = gUMMEIvAl.plmn.mcc_digit1;                  \
    gUTI->gummei.plmn.mcc_digit2 = gUMMEIvAl.plmn.mcc_digit2;                  \
    gUTI->gummei.plmn.mcc_digit3 = gUMMEIvAl.plmn.mcc_digit3;                  \
    gUTI->gummei.plmn.mnc_digit1 = gUMMEIvAl.plmn.mnc_digit1;                  \
    gUTI->gummei.plmn.mnc_digit2 = gUMMEIvAl.plmn.mnc_digit2;                  \
    gUTI->gummei.plmn.mnc_digit3 = gUMMEIvAl.plmn.mnc_digit3;                  \
  } while (0)

#define COPY_PLMN(pLMN1, pLMN2)                                                \
  do {                                                                         \
    pLMN1.mcc_digit1 = pLMN2.mcc_digit1;                                       \
    pLMN1.mcc_digit2 = pLMN2.mcc_digit2;                                       \
    pLMN1.mcc_digit3 = pLMN2.mcc_digit3;                                       \
    pLMN1.mnc_digit1 = pLMN2.mnc_digit1;                                       \
    pLMN1.mnc_digit2 = pLMN2.mnc_digit2;                                       \
    pLMN1.mnc_digit3 = pLMN2.mnc_digit3;                                       \
  } while (0)

#define COPY_PLMN_IN_ARRAY_FMT(pLMN1, pLMN2)                                   \
  do {                                                                         \
    pLMN1.mcc[0] = pLMN2.mcc_digit1;                                           \
    pLMN1.mcc[1] = pLMN2.mcc_digit2;                                           \
    pLMN1.mcc[2] = pLMN2.mcc_digit3;                                           \
    pLMN1.mnc[0] = pLMN2.mnc_digit1;                                           \
    pLMN1.mnc[1] = pLMN2.mnc_digit2;                                           \
    pLMN1.mnc[2] = pLMN2.mnc_digit3;                                           \
  } while (0)

#define OAILOG_DEBUG_GUTI(gUTI_p)                                              \
  do {                                                                         \
    OAILOG_DEBUG(                                                              \
        LOG_MME_APP,                                                           \
        "mcc-1: %x mcc-2: %x mcc-3: %x mnc-1: %x mnc-2: %x mnc-3: %x"          \
        " mme_group_id: %x mme_code: %x m_tmsi: %x \n",                        \
        gUTI_p->gummei.plmn.mcc_digit1, gUTI_p->gummei.plmn.mcc_digit2,        \
        gUTI_p->gummei.plmn.mcc_digit3, gUTI_p->gummei.plmn.mnc_digit1,        \
        gUTI_p->gummei.plmn.mnc_digit2, gUTI_p->gummei.plmn.mnc_digit3,        \
        gUTI_p->gummei.mme_gid, gUTI_p->gummei.mme_code, gUTI_p->m_tmsi);      \
  } while (0)

#endif /* FILE_MME_API_SEEN*/
