/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the Apache License, Version 2.0  (the "License"); you may not use this file
 * except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
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

/*! \file s1ap_mme.h
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#ifndef FILE_S1AP_MME_SEEN
#define FILE_S1AP_MME_SEEN

#if MME_CLIENT_TEST == 0
#include "intertask_interface.h"
#endif

#include <stdbool.h>
#include <stdint.h>

#include "hashtable.h"
#include "3gpp_36.401.h"
#include "common_types.h"

// Forward declarations
struct enb_description_s;

#define S1AP_TIMER_INACTIVE_ID (-1)
#define S1AP_UE_CONTEXT_REL_COMP_TIMER 1 // in seconds

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
  S1AP_INIT, /// The sctp association has been established but s1 hasn't been setup.
  S1AP_RESETING, /// The s1state is resetting due to an SCTP reset on the bound association.
  S1AP_READY,   ///< MME and eNB are S1 associated, UE contexts can be added
  S1AP_SHUTDOWN /// The S1 state is being torn down due to sctp shutdown.
};

enum s1_ue_state_s {
  S1AP_UE_INVALID_STATE,
  S1AP_UE_WAITING_CSR, ///< Waiting for Initial Context Setup Response
  S1AP_UE_HANDOVER,    ///< Handover procedure triggered
  S1AP_UE_CONNECTED,   ///< UE context ready
  S1AP_UE_WAITING_CRR, /// UE Context release Procedure initiated , waiting for UE context Release Complete
};

/** Main structure representing UE association over s1ap
 *  Generated every time a new InitialUEMessage is received
 **/
typedef struct ue_description_s {
  struct enb_description_s *enb; ///< Which eNB this UE is attached to

  enum s1_ue_state_s s1_ue_state; ///< S1AP UE state

  enb_ue_s1ap_id_t
    enb_ue_s1ap_id : 24;           ///< Unique UE id over eNB (24 bits wide)
  mme_ue_s1ap_id_t mme_ue_s1ap_id; ///< Unique UE id over MME (32 bits wide)

  /** SCTP stream on which S1 message will be sent/received.
   *  During an UE S1 connection, a pair of streams is
   *  allocated and is used during all the connection.
   *  Stream 0 is reserved for non UE signalling.
   *  @name sctp stream identifier
   **/
  /*@{*/
  sctp_stream_id_t sctp_stream_recv; ///< eNB -> MME stream
  sctp_stream_id_t sctp_stream_send; ///< MME -> eNB stream
  /*@}*/

  s11_teid_t s11_sgw_teid;

  // UE Context Release procedure guard timer
  struct s1ap_timer_t s1ap_ue_context_rel_timer;

} ue_description_t;

/* Main structure representing eNB association over s1ap
 * Generated (or updated) every time a new S1SetupRequest is received.
 */
typedef struct enb_description_s {
  enum mme_s1_enb_state_s
    s1_state; ///< State of the eNB specific S1AP association

  /** eNB related parameters **/
  /*@{*/
  char enb_name[150];         ///< Printable eNB Name
  uint32_t enb_id;            ///< Unique eNB ID
  uint8_t default_paging_drx; ///< Default paging DRX interval for eNB
  /*@}*/

  /** UE list for this eNB **/
  /*@{*/
  uint32_t nb_ue_associated; ///< Number of NAS associated UE on this eNB
  hash_table_ts_t
    ue_coll; // contains ue_description_s, key is ue_description_s.?;
  /*@}*/
  // Wait for associated UE clean-up timer during sctp shutdown
  struct s1ap_timer_t s1ap_enb_assoc_clean_up_timer;
  /** SCTP stuff **/
  /*@{*/
  sctp_assoc_id_t sctp_assoc_id;     ///< SCTP association id on this machine
  sctp_stream_id_t next_sctp_stream; ///< Next SCTP stream
  sctp_stream_id_t instreams;  ///< Number of streams avalaible on eNB -> MME
  sctp_stream_id_t outstreams; ///< Number of streams avalaible on MME -> eNB
  /*@}*/
} enb_description_t;

extern bool hss_associated;
extern uint32_t nb_enb_associated;
extern struct mme_config_s *global_mme_config_p;

/** \brief S1AP layer top init
 * @returns -1 in case of failure
 **/
int s1ap_mme_init(void);

/** \brief S1AP layer top exit
 **/
void s1ap_mme_exit(void);

/** \brief Look for given eNB id in the list
 * \param enb_id The unique eNB id to search in list
 * @returns NULL if no eNB matchs the eNB id, or reference to the eNB element in list if matches
 **/
enb_description_t *s1ap_is_enb_id_in_list(const uint32_t enb_id);

/** \brief Look for given eNB SCTP assoc id in the list
 * \param enb_id The unique sctp assoc id to search in list
 * @returns NULL if no eNB matchs the sctp assoc id, or reference to the eNB element in list if matches
 **/
enb_description_t *s1ap_is_enb_assoc_id_in_list(
  const sctp_assoc_id_t sctp_assoc_id);

/** \brief Look for given ue eNB id in the list
 * \param enb_id The unique ue_eNB id to search in list
 * @returns NULL if no UE matchs the ue_enb_id, or reference to the ue element in list if matches
 **/
ue_description_t *s1ap_is_ue_enb_id_in_list(
  enb_description_t *enb_ref,
  const enb_ue_s1ap_id_t enb_ue_s1ap_id);

/** \brief Look for given ue mme id in the list
 * \param enb_id The unique ue_mme_id to search in list
 * @returns NULL if no UE matchs the ue_mme_id, or reference to the ue element in list if matches
 **/
ue_description_t *s1ap_is_ue_mme_id_in_list(const mme_ue_s1ap_id_t ue_mme_id);
ue_description_t *s1ap_is_s11_sgw_teid_in_list(const s11_teid_t teid);

/** \brief associate mainly 2(3) identifiers in S1AP layer: {mme_ue_s1ap_id_t, sctp_assoc_id (,enb_ue_s1ap_id)}
 **/
void s1ap_notified_new_ue_mme_s1ap_id_association(
  const sctp_assoc_id_t sctp_assoc_id,
  const enb_ue_s1ap_id_t enb_ue_s1ap_id,
  const mme_ue_s1ap_id_t mme_ue_s1ap_id);

/** \brief Allocate and add to the list a new eNB descriptor
 * @returns Reference to the new eNB element in list
 **/
enb_description_t *s1ap_new_enb(void);

/** \brief Allocate and add to the right eNB list a new UE descriptor
 * \param sctp_assoc_id association ID over SCTP
 * \param enb_ue_s1ap_id ue ID over S1AP
 * @returns Reference to the new UE element in list
 **/
ue_description_t *s1ap_new_ue(
  const sctp_assoc_id_t sctp_assoc_id,
  enb_ue_s1ap_id_t enb_ue_s1ap_id);

/** \brief Dump the eNB related information.
 * hashtable callback. It is called by hashtable_ts_apply_funct_on_elements()
 * Calls s1ap_dump_enb
 **/
bool s1ap_dump_enb_hash_cb(
  const hash_key_t keyP,
  void *const enb_void,
  void *unused_parameterP,
  void **unused_resultP);

/** \brief Dump the eNB list
 * Calls dump_enb for each eNB in list
 **/
void s1ap_dump_enb_list(void);

/** \brief Dump eNB related information.
 * Calls dump_ue for each UE in list
 * \param enb_ref eNB structure reference to dump
 **/
void s1ap_dump_enb(const enb_description_t *const enb_ref);

/** \brief Dump the ue related information.
 * hashtable callback. It is called by hashtable_ts_apply_funct_on_elements()
 * Calls s1ap_dump_ue
 **/
bool s1ap_dump_ue_hash_cb(
  const hash_key_t keyP,
  void *const ue_void,
  void *unused_parameterP,
  void **unused_resultP);

/** \brief Dump UE related information.
 * \param ue_ref ue structure reference to dump
 **/
void s1ap_dump_ue(const ue_description_t *const ue_ref);

bool s1ap_enb_compare_by_enb_id_cb(
  const hash_key_t keyP,
  void *const elementP,
  void *parameterP,
  void __attribute__((unused)) * *unused_resultP);

/** \brief Remove target UE from the list
 * \param ue_ref UE structure reference to remove
 **/
void s1ap_remove_ue(ue_description_t *ue_ref);

/** \brief Remove target eNB from the list and remove any UE associated
 * \param enb_ref eNB structure reference to remove
 **/
void s1ap_remove_enb(enb_description_t *enb_ref);

#endif /* FILE_S1AP_MME_SEEN */
