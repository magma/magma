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
/****************************************************************************
  Source      amf_app_ue_context.h
  Version     0.1
  Date        2020/07/28
  Product     AMF core
  Subsystem   Access and Mobility Management Function
  Author      Ashish Prajapati
  Description Defines Access and Mobility Management Messages

*****************************************************************************/

#ifndef FILE_AMF_APP_UE_CONTEXT_SEEN
#define FILE_AMF_APP_UE_CONTEXT_SEEN
#include <stdint.h>
#include <inttypes.h> /* For sscanf formats */
#include <time.h>     /* to provide time_t */

#include "tree.h"
#include "queue.h"
#include "hashtable.h"
#include "obj_hashtable.h"
#include "bstrlib.h"
#include "common_types.h"
#include "common_defs.h"
#include "ngap_messages_types.h"
#include "security_types.h"
#include "sgw_ie_defs.h"
#include "amf_app_sgs_fsm.h"
#include "intertask_interface_types.h"

typedef enum {
  M5GCM_IDLE = 0,
  M5GCM_CONNECTED,
} m5gcm_state_t;

#define IMSI_DIGITS_MAX 15

typedef struct {
  uint32_t length;
  char data[IMSI_DIGITS_MAX + 1];
} amf_app_imsi_t;

// TODO: (amar) only used in testing
#define IMSI_FORMAT "s"
#define IMSI_DATA(AMF_APP_IMSI) (AMF_APP_IMSI.data)

/* Convert the IMSI contained by a char string NULL terminated to uint64_t */
#if 0
//TODO:  
bool amf_app_is_imsi_empty(amf_app_imsi_t const* imsi);

bool amf_app_imsi_compare(
    amf_app_imsi_t const* imsi_a, amf_app_imsi_t const* imsi_b);

void amf_app_copy_imsi(
    amf_app_imsi_t* imsi_dst, amf_app_imsi_t const* imsi_src);

void amf_app_string_to_imsi(
    amf_app_imsi_t* const imsi_dst, char const* const imsi_string_src);

void amf_app_imsi_to_string(
    char* const imsi_dst, amf_app_imsi_t const* const imsi_src);

uint64_t amf_app_imsi_to_u64(amf_app_imsi_t imsi_src);

void amf_app_ue_context_uint_to_imsi(
    uint64_t imsi_src, amf_app_imsi_t* imsi_dst);

void amf_app_convert_imsi_to_imsi_amf(
    amf_app_imsi_t* imsi_dst, const imsi_t* imsi_src);

amf_ue_ngap_id_t amf_app_ctx_get_new_ue_id(
    amf_ue_ngap_id_t* amf_app_ue_ngap_id_generator_p);
#endif
/*
 * Timer identifier returned when in inactive state (timer is stopped or has
 * failed to be started)
 */
#define AMF_APP_TIMER_INACTIVE_ID (-1)

#define AMF_APP_DELTA_T3412_REACHABILITY_TIMER 4            // in minutes
#define AMF_APP_DELTA_REACHABILITY_IMPLICIT_DETACH_TIMER 0  // in minutes

#define AMF_APP_INITIAL_CONTEXT_SETUP_RSP_TIMER_VALUE 2  // In seconds
#define AMF_APP_UE_CONTEXT_MODIFICATION_TIMER_VALUE 2    // In seconds
#define AMF_APP_PAGING_RESPONSE_TIMER_VALUE 4            // In seconds
#define AMF_APP_ULR_RESPONSE_TIMER_VALUE 3               // In seconds
/* Timer structure */
struct amf_app_timer_t {
  long id;      /* The timer identifier                 */
  uint32_t sec; /* The timer interval value in seconds  */
};

/** @struct bearer_context_t
 *  @brief Parameters that should be kept for an eps bearer.
 */
typedef struct bearer_context_s {
  /* EPS Bearer ID: An EPS bearer identity uniquely identifies
   * an EP S bearer for one UE accessing via E-UTRAN
   */
  ebi_t ebi;

  /* Procedure Transaction Identifier */
  proc_tid_t transaction_identifier;

  /* S-GW IP address for S1-u interfaces.
   * S-GW TEID for S1-u interface.
   * set by N11,gc CREATE_SESSION_RESPONSE
   */
  fteid_t s_gw_fteid_s1u;

  /* PDN GW TEID for S5/S8 (user plane), Used for S-GW change only
   * PDN GW IP address for S5/S8 (user plane), Used for S-GW change only
   *
   * NOTE:
   *      The PDN GW TEID and PDN GW IP address for user plane are needed
   *      in AMF context as S-GW relocation is triggered without interaction
   *      with the source S-GW,
   *      e.g. when a TAU occurs. The Target SGW requires this Information
   *      Element, so it must be stored by the AMF.
   */
  fteid_t p_gw_fteid_s5_s8_up;

  pdn_cid_t pdn_cx_id;
  esm_ebr_context_t esm_ebr_context;
  fteid_t gnb_fteid_s1u;

  /* QoS for this bearer */
  qci_t qci;
  priority_level_t priority_level;
  pre_emption_vulnerability_t preemption_vulnerability;
  pre_emption_capability_t preemption_capability;
} bearer_context_t;

/** @struct pdn_context_s
 *  Parameters that should be kept for a subscribed apn by the UE.
 */
typedef struct pdn_context_s {
  context_identifier_t context_identifier;

  /* APN in Use:an ID at P-GW through which a user can access the Subscribed APN
   *            This APN shall be composed of the APN Network
   *            Identifier and the default APN Operator Identifier,
   *            as specified in TS 23.003 [9],
   *            clause 9.1.2 (EURECOM: "mnc<MNC>.mcc<MCC>.gprs").
   *            Any received value in the APN OI Replacement field is not
   *            applied here.
   */
  bstring apn_in_use;

  /* APN Subscribed: The subscribed APN received from the HSS */
  bstring apn_subscribed;

  /* PDN Type: IPv4, IPv6 or IPv4v6 */
  pdn_type_t pdn_type;

  /* paa: IPv4 address and/or IPv6 prefix of UE set by
   *          N11,gc CREATE_SESSION_RESPONSE
   *          NOTE:
   *          The AMF might not have information on the allocated IPv4 address.
   *          Alternatively, following mobility involving a pre-release 8 SGSN,
   *          This IPv4 address might not be the one allocated to the UE.
   */
  paa_t paa;

  /* APN-OI Replacement: APN level APN-OI Replacement which has same role as
   *            UE level APN-OI Replacement but with higher priority than
   *            UE level APN-OI Replacement. This is and optional parameter.
   *            When available, it shall be used to construct the PDN GW
   *            FQDN instead of UE level APN-OI Replacement.
   */
  bstring apn_oi_replacement;

  /* PDN GW Address in Use(control plane): The IP address of the PDN GW
   *           currently
   *           used for sending control plane signalling.
   */
  ip_address_t p_gw_address_s5_s8_cp;

  /* PDN GW TEID for S5/S8 (control plane), for GTP-based S5/S8 only */
  teid_t p_gw_teid_s5_s8_cp;

  /* EPS subscribed QoS profile:
   *            The bearer level QoS parameter values for that
   *            APN's default bearer's QCI and ARP (see clause 4.7.3).
   */
  eps_subscribed_qos_profile_t default_bearer_eps_subscribed_qos_profile;

  /* Subscribed APN-AMBR: The Maximum Aggregated uplink and downlink MBR values
   *            to be shared across all Non-GBR bearers,
   *             which are established for this APN, according to the
   *            subscription of the user.
   */
  ambr_t subscribed_apn_ambr;

  /* p_gw_apn_ambr: The Maximum Aggregated uplink and downlink MBR values to be
   *           shared across all Non-GBR bearers, which are established for this
   *           APN, as decided by the PDN GW.
   */
  ambr_t p_gw_apn_ambr;

  /* default_ebi: Identifies the EPS Bearer Id of the default bearer
   * within the given PDN connection.
   */
  ebi_t default_ebi;

  /* bearer_contexts[]: contains bearer indexes in
   *           ue_mm_context_t.bearer_contexts[], or -1
   */
  int bearer_contexts[BEARERS_PER_UE];

  /* S-GW teid and IP address for User-Plane
   * set by N11,gc CREATE_SESSION_RESPONSE
   */

  ip_address_t s_gw_address_s11_s4;
  teid_t s_gw_teid_s11_s4;

  esm_pdn_t esm_data;
  /* is_active == true indicates, PDN is active */
  bool is_active;

  protocol_configuration_options_t* pco;
} pdn_context_t;

typedef enum {
  GRANTED_SERVICE_EPS_ONLY,
  GRANTED_SERVICE_SMS_ONLY,
  GRANTED_SERVICE_CSFB_SMS
} granted_service_t;
#endif
/** @struct sgs_context_t
 *  @brief SGS Parameters that should be kept per UE.
 */
typedef struct sgs_context_s {
  sgs_fsm_state_t sgs_state;
  bool vlr_reliable;
#define SET_NEAF true;
#define RESET_NEAF false;
  /* Non EPS Alert Flag */
  bool neaf;
  /* SGS Location update timer */
  struct amf_app_timer_t ts6_1_timer;

#define EPS_DETACH_RETRANSMISSION_COUNTER_MAX 2
#define IMSI_DETACH_RETRANSMISSION_COUNTER_MAX 2
#define IMPLICIT_IMSI_DETACH_RETRANSMISSION_COUNTER_MAX 2
#define IMPLICIT_EPS_DETACH_RETRANSMISSION_COUNTER_MAX 2

  /* SGS EPS Detach indication timer */
  struct amf_app_timer_t ts8_timer;
  unsigned int ts8_retransmission_count;
  /* SGS IMSI Detach indication timer */
  struct amf_app_timer_t ts9_timer;
  unsigned int ts9_retransmission_count;
  /* SGS IMPLICIT IMSI DETACH INDICATION timer */
  struct amf_app_timer_t ts10_timer;
  unsigned int ts10_retransmission_count;
  /* SGS IMPLICIT EPS DETACH INDICATION timer */
  struct amf_app_timer_t ts13_timer;
  unsigned int ts13_retransmission_count;

  /* message_p: To store NGAP NAS DL DATA REQ in case of UE initiated IMSI or
   *             combined EPS/IMSI detach and
   *             send after recieving SGS IMSI Detach Ack
   */
  MessageDef* message_p;

  /* sgsap_msg: Received message over SGS interface */
  void* sgsap_msg;

  /* ongoing_procedure_t: SGS Location update procedure initiated due combined
   *             attach procedure or TAU procedure
   */
  ongoing_procedure_t ongoing_procedure;

  /* tau_active_flag: Value of active flagreceived in TAU Request */
  uint8_t tau_active_flag : 1;

  /* service_indicator: store the requested service (SMS or call),
   *             that shall be sent in SGS-Service Request
   *             while UE is in idle mode
   */
  uint8_t service_indicator;

  /* Indicates ongoing CSFB procedure */
  csfb_service_type_t csfb_service_type;

  /* Call Cancelled: is set on reception of SGS SERVICE ABORT message
   *             fom MSC to cancel the ongoing MT call
   */
  bool call_cancelled;

  /* mt_call_in_progress: If true, indicates MT call is in progress,
   *              used when SERVICE ABORT is received from MSC
   */
  bool mt_call_in_progress;

  /* is_emergency_call: True - if the call is of type Emergency call */
  bool is_emergency_call;
} sgs_context_t;

/** @struct ue_mm_context_t
 *  @brief Useful parameters to know in AMF application layer. They are set
 * according to 3GPP TS.23.401 #5.7.2
 */
typedef struct ue_mm_context_s {
  /* msisdn: The basic MSISDN of the UE. The presence is dictated by its storage
   *         in the HSS, set by S6A UPDATE LOCATION ANSWER
   */
  bstring msisdn;

  enum s1cause ue_context_rel_cause;
  mm_state_t mm_state;
  m5gcm_state_t m5gcm_state;

  /* Last known E-UTRAN cell, set by nas_attach_req_t */
  ecgi_t e_utran_cgi;

  /* cell_age: Time elapsed since the last E-UTRAN Cell Global Identity was
   *           acquired. Set by nas_auth_param_req_t
   */
  time_t cell_age;
  /* TODO: add csg_id */
  /* TODO: add csg_membership */
  /* TODO Access mode: Access mode of last known ECGI when the UE was active */

  /* apn_config_profile: set by S6A UPDATE LOCATION ANSWER */
  apn_config_profile_t apn_config_profile;

  /* access_restriction_data: The access restriction subscription information.
   *           set by S6A UPDATE LOCATION ANSWER
   */
  ard_t access_restriction_data;

  /* apn_oi_replacement: Indicates the domain name to replace the APN-OI
   *           when constructing the PDN GW FQDN upon which to perform a
   *           DNS resolution.
   *           This replacement applies for all the APNs provided in the
   *           subscriber's profile. See TS 23.003 [9] clause 9.1.2 for more
   */
  bstring apn_oi_replacement;
  teid_t amf_teid_s11;
  /* SCTP assoc id */
  sctp_assoc_id_t sctp_assoc_id_key;

  /* eNB UE NGAP ID,  Unique identity the UE within eNodeB */
  gnb_ue_ngap_id_t gnb_ue_ngap_id : 24;
  /* gnb_ngap_id_key = gnb-ue-ngap-id <24 bits> | gnb-id <8 bits> */
  gnb_ngap_id_key_t gnb_ngap_id_key;

  /* AMF UE NGAP ID, Unique identity the UE within AMF */
  amf_ue_ngap_id_t amf_ue_ngap_id;

  /* Subscribed UE-AMBR: The Maximum Aggregated uplink and downlink MBR values
   *           to be shared across all Non-GBR bearers according to the
   *           subscription of the user. Set by S6A UPDATE LOCATION ANSWER
   */
  ambr_t subscribed_ue_ambr;
  /* used_ue_ambr: The currently used Maximum Aggregated uplink and downlink
   *           MBR values to be shared across all Non-GBR bearers.
   *           Set by S6A UPDATE LOCATION ANSWER
   */
  ambr_t used_ue_ambr;
  /* rau_tau_timer: Indicates a subscribed Periodic RAU/TAU Timer value
   *           Set by S6A UPDATE LOCATION ANSWER
   */
  rau_tau_timer_t rau_tau_timer;

  int nb_active_pdn_contexts;
  pdn_context_t* pdn_contexts[MAX_APN_PER_UE];

  emm_context_t emm_context;
  bearer_context_t* bearer_contexts[BEARERS_PER_UE];

  /* ue_radio_capability: Store the radio capabilities as received in
   *           NGAP UE capability indication message
   */
  bstring ue_radio_capability;

  /* mobile_reachability_timer: Start when UE moves to idle state.
   *             Stop when UE moves to connected state
   */
  struct amf_app_timer_t mobile_reachability_timer;
  time_t time_mobile_reachability_timer_started;
  /* implicit_detach_timer: Start at the expiry of Mobile Reachability timer.
   * Stop when UE moves to connected state
   */
  struct amf_app_timer_t implicit_detach_timer;
  time_t time_implicit_detach_timer_started;
  /* Initial Context Setup Procedure Guard timer */
  struct amf_app_timer_t initial_context_setup_rsp_timer;
  /* UE Context Modification Procedure Guard timer */
  struct amf_app_timer_t ue_context_modification_timer;
  /* Timer for retrying paging messages */
  struct amf_app_timer_t paging_response_timer;
  /* send_ue_purge_request: If true AMF shall send S6a- Purge Req to
   * delete contexts at HSS
   */
  bool send_ue_purge_request;

  bool hss_initiated_detach;
  bool location_info_confirmed_in_hss;
  /* S6a- update location request guard timer */
  struct amf_app_timer_t ulr_response_timer;
  sgs_context_t* sgs_context;
  uint8_t attach_type;
  lai_t lai;
  int cs_fallback_indicator;
  uint8_t sgs_detach_type;
  /* granted_service_t: informs the granted service to UE */
  granted_service_t granted_service;
  /*  paging_proceeding_flag (PPF) shall set to true, when UE moves to
   * connected. Indicates that paging procedure can be prooceeded, Is set to
   * false, due to "Inactivity of UE including lack of periodic TAU"
   */
  bool ppf;

#define SUBSCRIPTION_UNKNOWN false
#define SUBSCRIPTION_KNOWN true
  bool subscription_known;
  ambr_t used_ambr;
  subscriber_status_t subscriber_status;
  network_access_mode_t network_access_mode;

  bool path_switch_req;
  LIST_HEAD(s11_procedures_s, amf_app_s11_proc_s) * s11_procedures;
} ue_mm_context_t;

typedef struct amf_ue_context_s {
  hash_table_uint64_ts_t* imsi_amf_ue_id_htbl;    // data is amf_ue_ngap_id_t
  hash_table_uint64_ts_t* tun11_ue_context_htbl;  // data is amf_ue_ngap_id_t
  hash_table_uint64_ts_t*
      gnb_ue_ngap_id_ue_context_htbl;             // data is amf_ue_ngap_id_t
  obj_hash_table_uint64_t* guti_ue_context_htbl;  // data is amf_ue_ngap_id_t
} amf_ue_context_t;

/** \brief Retrieve an UE context by selecting the provided IMSI
 * \param imsi Imsi to find in UE map
 * @returns an UE context matching the IMSI or NULL if the context doesn't
 *exists
 **/
ue_mm_context_t* amf_ue_context_exists_imsi(
    amf_ue_context_t* const amf_ue_context, const imsi64_t imsi);

/** \brief Retrieve an UE context by selecting the provided N11,gc teid
 * \param teid The tunnel endpoint identifier used between AMF and S-GW
 * @returns an UE context matching the teid or NULL if the context doesn't
 *exists
 **/
ue_mm_context_t* amf_ue_context_exists_s11_teid(
    amf_ue_context_t* const amf_ue_context, const s11_teid_t teid);

/** \brief Retrieve an UE context by selecting the provided amf_ue_ngap_id
 * \param amf_ue_ngap_id The UE id identifier used in NGAP AMF (and NAS)
 * @returns an UE context matching the amf_ue_ngap_id or NULL if the context
 *doesn't exists
 **/
ue_mm_context_t* amf_ue_context_exists_amf_ue_ngap_id(
    const amf_ue_ngap_id_t amf_ue_ngap_id);

/** \brief Retrieve an UE context by selecting the provided gnb_ue_ngap_id
 * \param gnb_ue_ngap_id The UE id identifier used in NGAP AMF
 * @returns an UE context matching the gnb_ue_ngap_id or NULL if the context
 *doesn't exists
 **/
ue_mm_context_t* amf_ue_context_exists_gnb_ue_ngap_id(
    amf_ue_context_t* const amf_ue_context_p, const gnb_ngap_id_key_t gnb_key);

/** \brief Retrieve an UE context by selecting the provided guti
 * \param guti The GUTI used by the UE
 * @returns an UE context matching the guti or NULL if the context doesn't
 *exists
 **/
ue_mm_context_t* amf_ue_context_exists_guti(
    amf_ue_context_t* const amf_ue_context, const guti_t* const guti);

/** \brief Move the content of a context to another context
 * \param dst            The destination context
 * \param src            The source context
 **/
void amf_app_move_context(ue_mm_context_t* dst, ue_mm_context_t* src);

/** \brief Notify the AMF_APP that a duplicated ue_context_t exist (both share
 * the same amf_ue_ngap_id)
 * \param gnb_key The UE id identifier used in NGAP and AMF_APP (agregated with
 * a gnb_id)
 * \param amf_ue_ngap_id The UE id identifier used in AMF_APP and NAS
 * \param is_remove_old  Remove old UE context or new UE context ?
 **/
void amf_ue_context_duplicate_gnb_ue_ngap_id_detected(
    const gnb_ngap_id_key_t gnb_key, const amf_ue_ngap_id_t amf_ue_ngap_id,
    const bool is_remove_old);

/** \brief Update an UE context by selecting the provided guti
 * \param amf_ue_context_p The AMF context
 * \param ue_context_p The UE context
 * \param gnb_ngap_id_key The eNB UE id identifier
 * \param amf_ue_ngap_id The UE id identifier used in NGAP AMF (and NAS)
 * \param imsi
 * \param len
 * \param amf_s11_teid The tunnel endpoint identifier used between AMF and S-GW
 * \param nas_ue_id The UE id identifier used in NGAP AMF and NAS
 * \param guti_p The GUTI used by the UE
 **/
void amf_ue_context_update_coll_keys(
    amf_ue_context_t* const amf_ue_context_p,
    ue_mm_context_t* const ue_context_p,
    const gnb_ngap_id_key_t gnb_ngap_id_key,
    const amf_ue_ngap_id_t amf_ue_ngap_id, const imsi64_t imsi,
    const s11_teid_t amf_s11_teid, const guti_t* const guti_p);

/** \brief dump AMF associative collections
 **/

void amf_ue_context_dump_coll_keys(const amf_ue_context_t* amf_ue_contexts_p);

/** \brief Insert a new UE context in the tree of known UEs.
 * At least the IMSI should be known to insert the context in the tree.
 * \param ue_context_p The UE context to insert
 * @returns 0 in case of success, -1 otherwise
 **/
int amf_insert_ue_context(
    amf_ue_context_t* const amf_ue_context,
    const struct ue_mm_context_s* const ue_context_p);

/** \brief TODO WORK HERE Remove UE context unnecessary information.
 * mark it as released. It is necessary to keep track of the association
 * (s_tmsi (guti), amf_ue_ngap_id)
 * \param ue_context_p The UE context to remove
 **/
void amf_notify_ue_context_released(
    amf_ue_context_t* const amf_ue_context_p,
    struct ue_mm_context_s* ue_context_p);

/** \brief Remove a UE context of the tree of known UEs.
 * \param ue_context_p The UE context to remove
 **/
void amf_remove_ue_context(
    amf_ue_context_t* const amf_ue_context,
    struct ue_mm_context_s* const ue_context_p);

/** \brief Allocate memory for a new UE context
 * @returns Pointer to the new structure, NULL if allocation failed
 **/
ue_mm_context_t* amf_create_new_ue_context(void);

void amf_app_free_pdn_connection(pdn_context_t** const pdn_connection);

void amf_app_ue_context_free_content(ue_mm_context_t* const amf_ue_context_p);

/**
 * Release memory allocated by MmeNasStateManager through MmeNasStateConverter
 * and NasStateConverter for each UE context, this is called by
 * hashtable_ts_destroy
 */
void amf_app_state_free_ue_context(void** ue_context_node);

/** \brief Dump the UE contexts present in the tree
 **/
void amf_app_dump_ue_contexts(void);

void amf_app_handle_ngap_ue_context_release_req(
    const itti_ngap_ue_context_release_req_t* ngap_ue_context_release_req);

bearer_context_t* amf_app_get_bearer_context(
    ue_mm_context_t* const ue_context, const ebi_t ebi);

void amf_app_handle_gnb_deregister_ind(
    const itti_ngap_gNB_deregistered_ind_t* gNB_deregistered_ind);

ebi_t amf_app_get_free_bearer_id(ue_mm_context_t* const ue_context);

void amf_app_free_bearer_context(bearer_context_t** bc);

void amf_app_send_delete_session_request(
    struct ue_mm_context_s* const ue_context_p, const ebi_t ebi,
    const pdn_cid_t cid);

void amf_app_handle_ngap_ue_context_modification_resp(
    const itti_ngap_ue_context_mod_resp_t* ngap_ue_context_mod_resp);
void amf_app_handle_ngap_ue_context_modification_fail(
    const itti_ngap_ue_context_mod_resp_fail_t* ngap_ue_context_mod_fail);

void amf_app_ue_sgs_context_free_content(
    sgs_context_t* const sgs_context_p, imsi64_t imsi);
bool is_amf_ue_context_network_access_mode_packet_only(
    ue_mm_context_t* ue_context_p);
int amf_app_send_s6a_update_location_req(
    struct ue_mm_context_s* const ue_context_pP);
void amf_app_recover_timers_for_all_ues(void);

#endif /* FILE_AMF_APP_UE_CONTEXT_SEEN */

/* @} */
