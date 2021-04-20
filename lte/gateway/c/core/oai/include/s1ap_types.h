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

#pragma once

#include <stdint.h>

#include "3gpp_36.401.h"
#include "3gpp_36.413.h"

#include "common_types.h"
#include "hashtable.h"

// Forward declarations
struct enb_description_s;

#define S1AP_TIMER_INACTIVE_ID (-1)
#define S1AP_UE_CONTEXT_REL_COMP_TIMER 1  // in seconds

typedef struct s1ap_state_s {
  // contains eNB_description_s, key is eNB_description_s.enb_id (uint32_t)
  hash_table_ts_t enbs;
  // contains sctp association id, key is mme_ue_s1ap_id
  hash_table_ts_t mmeid2associd;
  uint32_t num_enbs;
} s1ap_state_t;

typedef struct s1ap_imsi_map_s {
  hash_table_uint64_ts_t* mme_ue_id_imsi_htbl;
} s1ap_imsi_map_t;

enum s1_timer_class_s {
  S1AP_INVALID_TIMER_CLASS,
  S1AP_ENB_TIMER,
  S1AP_UE_TIMER
};

/* S1AP Timer argument */
typedef struct s1ap_timer_arg_s {
  enum s1_timer_class_s timer_class;
  uint32_t instance_id;
} s1ap_timer_arg_t;

/* Timer structure */
struct s1ap_timer_t {
  long id;  /* The timer identifier                 */
  long sec; /* The timer interval value in seconds  */
};

// The current s1 state of the MME relating to the specific eNB.
enum mme_s1_enb_state_s {
  S1AP_INIT,  /// The sctp association has been established but s1 hasn't been
              /// setup.
  S1AP_RESETING,  /// The s1state is resetting due to an SCTP reset on the bound
                  /// association.
  S1AP_READY,     ///< MME and eNB are S1 associated, UE contexts can be added
  S1AP_SHUTDOWN   /// The S1 state is being torn down due to sctp shutdown.
};

enum s1_ue_state_s {
  S1AP_UE_INVALID_STATE,
  S1AP_UE_WAITING_CSR,  ///< Waiting for Initial Context Setup Response
  S1AP_UE_HANDOVER,     ///< Handover procedure triggered
  S1AP_UE_CONNECTED,    ///< UE context ready
  S1AP_UE_WAITING_CRR,  /// UE Context release Procedure initiated , waiting for
                        /// UE context Release Complete
};

typedef struct s1ap_handover_state_s {
  mme_ue_s1ap_id_t mme_ue_s1ap_id;
  uint32_t source_enb_id;
  uint32_t target_enb_id;
  enb_ue_s1ap_id_t
      target_enb_ue_s1ap_id : 24;  ///< Unique UE id over eNB (24 bits wide)
  sctp_stream_id_t target_sctp_stream_recv;  ///< eNB -> MME stream
  sctp_stream_id_t target_sctp_stream_send;  ///< MME -> eNB stream
  e_rab_admitted_list_t e_rab_admitted_list;
} s1ap_handover_state_t;

/** Main structure representing UE association over s1ap
 *  Generated every time a new InitialUEMessage is received
 **/
typedef struct ue_description_s {
  enum s1_ue_state_s s1_ue_state;  ///< S1AP UE state

  enb_ue_s1ap_id_t
      enb_ue_s1ap_id : 24;          ///< Unique UE id over eNB (24 bits wide)
  mme_ue_s1ap_id_t mme_ue_s1ap_id;  ///< Unique UE id over MME (32 bits wide)
  sctp_assoc_id_t sctp_assoc_id;  ///< Assoc id of eNB which this UE is attached
  uint64_t comp_s1ap_id;          ///< Unique composite UE id (sctp_assoc_id &
                                  ///< enb_ue_s1ap_id)

  /** SCTP stream on which S1 message will be sent/received.
   *  During an UE S1 connection, a pair of streams is
   *  allocated and is used during all the connection.
   *  Stream 0 is reserved for non UE signalling.
   *  @name sctp stream identifier
   **/
  /*@{*/
  sctp_stream_id_t sctp_stream_recv;  ///< eNB -> MME stream
  sctp_stream_id_t sctp_stream_send;  ///< MME -> eNB stream
  /*@}*/

  // UE Context Release procedure guard timer
  struct s1ap_timer_t s1ap_ue_context_rel_timer;

  // Handover status. We intentionally do not persist all of this state since
  // it's time sensitive; if the MME restarts during a HO procedure the RAN
  // will abort the procedure due to timeouts, rendering this state useless.
  s1ap_handover_state_t s1ap_handover_state;
} ue_description_t;

/* Maximum no. of Broadcast PLMNs. Value is 6
 * 3gpp spec 36.413 section-9.1.8.4
 */
#define S1AP_MAX_BROADCAST_PLMNS 6
/* Maximum TAI Items configured, can be upto 256 */
#define S1AP_MAX_TAI_ITEMS 16

/* Supported TAI items includes TAC and Broadcast PLMNs */
typedef struct supported_tai_items_s {
  uint16_t tac;             ///< Supported TAC value
  uint8_t bplmnlist_count;  ///< Number of Broadcast PLMNs in the TAI
  plmn_t bplmns[S1AP_MAX_BROADCAST_PLMNS];  ///< List of Broadcast PLMNS
} supported_tai_items_t;

/* Supported TAs by eNB received in S1 Setup request message */
typedef struct supported_ta_list_s {
  uint8_t list_count;  ///< Number of TAIs in the list
  supported_tai_items_t
      supported_tai_items[S1AP_MAX_TAI_ITEMS];  ///< List of TAIs
} supported_ta_list_t;

/* Main structure representing eNB association over s1ap
 * Generated (or updated) every time a new S1SetupRequest is received.
 */
typedef struct enb_description_s {
  enum mme_s1_enb_state_s
      s1_state;  ///< State of the eNB specific S1AP association

  /** eNB related parameters **/
  /*@{*/
  char enb_name[150];          ///< Printable eNB Name
  uint32_t enb_id;             ///< Unique eNB ID
  uint8_t default_paging_drx;  ///< Default paging DRX interval for eNB
  supported_ta_list_t supported_ta_list;  ///< Supported TAs by eNB
  /*@}*/

  /** UE list for this eNB **/
  /*@{*/
  uint32_t nb_ue_associated;  ///< Number of NAS associated UE on this eNB
  hash_table_uint64_ts_t ue_id_coll;  ///< Contains comp_s1ap_id assoc to
                                      ///< enodeb, key is mme_ue_s1ap_id;
  /*@}*/
  /** SCTP stuff **/
  /*@{*/
  sctp_assoc_id_t sctp_assoc_id;      ///< SCTP association id on this machine
  sctp_stream_id_t next_sctp_stream;  ///< Next SCTP stream
  sctp_stream_id_t instreams;   ///< Number of streams avalaible on eNB -> MME
  sctp_stream_id_t outstreams;  ///< Number of streams avalaible on MME -> eNB
  char ran_cp_ipaddr[16];    ///< Network byte order IP address of eNB SCTP end
                             ///< point
  uint8_t ran_cp_ipaddr_sz;  ///< IP addr size for ran_cp_ipaddr
  /*@}*/
} enb_description_t;
