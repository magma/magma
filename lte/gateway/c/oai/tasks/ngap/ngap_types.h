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

#pragma once

#include <stdint.h>

#include "3gpp_38.401.h"

#include "common_types.h"
#include "hashtable.h"

#define NGAP_MAX_BROADCAST_PLMNS 6
/* Maximum TAI Items configured, can be upto 256 */
#define NGAP_MAX_TAI_ITEMS 16

#define NGAP_TIMER_INACTIVE_ID (-1)
#define NGAP_UE_CONTEXT_REL_COMP_TIMER 1  // in seconds

// Forward declarations
struct gnb_description_s;

/* Supported TAI items includes TAC and Broadcast PLMNs */
typedef struct m5g_supported_tai_items_s {
  uint16_t tac;             ///< Supported TAC value
  uint8_t bplmnlist_count;  ///< Number of Broadcast PLMNs in the TAI
  plmn_t bplmns[NGAP_MAX_BROADCAST_PLMNS];  ///< List of Broadcast PLMNS
} m5g_supported_tai_items_t;

// The current n1 state of the AMF relating to the specific gNB.
enum amf_ng_gnb_state_s {
  NGAP_INIT,  /// The sctp association has been established but n1 hasn't been
              /// setup.
  NGAP_RESETING,  /// The n1state is resetting due to an SCTP reset on the bound
                  /// association.
  NGAP_READY,     ///< AMF and gNB are N1 associated, UE contexts can be added
  NGAP_SHUTDOWN   /// The N1 state is being torn down due to sctp shutdown.
};

enum ng_ue_state_s {
  NGAP_UE_INVALID_STATE,
  NGAP_UE_WAITING_CSR,  ///< Waiting for Initial Context Setup Response
  NGAP_UE_HANDOVER,     ///< Handover procedure triggered
  NGAP_UE_CONNECTED,    ///< UE context ready
  NGAP_UE_WAITING_CRR,  /// UE Context release Procedure initiated , waiting for
                        /// UE context Release Complete
};

/* Supported TAs by gNB received in NG Setup request message */
typedef struct m5g_supported_ta_list_s {
  uint8_t list_count;  ///< Number of TAIs in the list
  m5g_supported_tai_items_t
      supported_tai_items[NGAP_MAX_TAI_ITEMS];  ///< List of TAIs
} m5g_supported_ta_list_t;

typedef struct ngap_state_s {
  // contains gNB_description_s, key is gNB_description_s.gnb_id (uint32_t)
  hash_table_ts_t gnbs;
  // contains sctp association id, key is amf_ue_ngap_id
  hash_table_ts_t amfid2associd;
  uint32_t num_gnbs;
} ngap_state_t;

/* Timer structure */
struct ngap_timer_t {
  long id;  /* The timer identifier                 */
  long sec; /* The timer interval value in seconds  */
};

/** Main structure representing UE association over ngap
 *  Generated every time a new InitialUEMessage is received
 **/
typedef struct m5g_ue_description_s {
  enum ng_ue_state_s ng_ue_state;  ///< NGAP UE state

  gnb_ue_ngap_id_t
      gnb_ue_ngap_id : 24;          ///< Unique UE id over gNB (24 bits wide)
  amf_ue_ngap_id_t amf_ue_ngap_id;  ///< Unique UE id over AMF (32 bits wide)
  sctp_assoc_id_t sctp_assoc_id;  ///< Assoc id of gNB which this UE is attached
  uint64_t comp_ngap_id;          ///< Unique composite UE id (sctp_assoc_id &
                                  ///< gnb_ue_ngap_id)

  /** SCTP stream on which NG message will be sent/received.
   *  During an UE NG connection, a pair of streams is
   *  allocated and is used during all the connection.
   *  Stream 0 is reserved for non UE signalling.
   *  @name sctp stream identifier
   **/
  /*@{*/
  sctp_stream_id_t sctp_stream_recv;  ///< gNB -> AMF stream
  sctp_stream_id_t sctp_stream_send;  ///< AMF -> gNB stream
  /*@}*/

  // UE Context Release procedure guard timer
  struct ngap_timer_t ngap_ue_context_rel_timer;
} m5g_ue_description_t;

/// typedef struct ngap_imsi_map_s {
typedef struct ngap_imsi_map_s {
  hash_table_uint64_ts_t* amf_ue_id_imsi_htbl;
} ngap_imsi_map_t;

enum n1_timer_class_s {
  NGAP_INVALID_TIMER_CLASS,
  NGAP_GNB_TIMER,
  NGAP_UE_TIMER
};

/* Main structure representing gNB association over ngap
 * Generated (or updated) every time a new NGSetupRequest is received.
 */
typedef struct gnb_description_s {
  enum amf_ng_gnb_state_s
      ng_state;  ///< State of the gNB specific NGAP association

  /** gNB related parameters **/
  /*@{*/
  char gnb_name[150];          ///< Printable gNB Name
  uint32_t gnb_id;             ///< Unique gNB ID
  uint8_t default_paging_drx;  ///< Default paging DRX interval for gNB
  m5g_supported_ta_list_t supported_ta_list;  ///< Supported TAs by gNB
  /*@}*/

  /** UE list for this gNB **/
  /*@{*/
  uint32_t nb_ue_associated;  ///< Number of NAS associated UE on this gNB
  hash_table_uint64_ts_t ue_id_coll;  ///< Contains comp_ngap_id assoc to
                                      ///< enodeb, key is amf_ue_ngap_id;
  /*@}*/
  /** SCTP stuff **/
  /*@{*/
  sctp_assoc_id_t sctp_assoc_id;      ///< SCTP association id on this machine
  sctp_stream_id_t next_sctp_stream;  ///< Next SCTP stream
  sctp_stream_id_t instreams;   ///< Number of streams avalaible on gNB -> AMF
  sctp_stream_id_t outstreams;  ///< Number of streams avalaible on AMF -> gNB
  /*@}*/
} gnb_description_t;
