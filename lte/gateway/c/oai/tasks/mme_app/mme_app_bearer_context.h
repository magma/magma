/*
 * Copyright (c) 2015, EURECOM (www.eurecom.fr)
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice,
 * this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 *
 * The views and conclusions contained in the software and documentation are
 * those of the authors and should not be interpreted as representing official
 * policies, either expressed or implied, of the FreeBSD Project.
 */

/*! \file mme_app_bearer_context.h
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#ifndef FILE_MME_APP_BEARER_CONTEXT_SEEN
#define FILE_MME_APP_BEARER_CONTEXT_SEEN

#include <stdbool.h>

#include "3gpp_24.007.h"
#include "bstrlib.h"
#include "common_types.h"
#include "mme_app_ue_context.h"

typedef uint8_t mme_app_bearer_state_t;

/*
 * @struct bearer_context_new_t
 * @brief Parameters that should be kept for an eps bearer. Used for stacked
 * memory
 *
 * Structure of an EPS bearer
 * --------------------------
 * An EPS bearer is a logical concept which applies to the connection
 * between two endpoints (UE and PDN Gateway) with specific QoS attri-
 * butes. An EPS bearer corresponds to one Quality of Service policy
 * applied within the EPC and E-UTRAN.
 */
typedef struct bearer_context_new_s {
  // EPS Bearer ID: An EPS bearer identity uniquely identifies an EP S bearer
  // for one UE accessing via E-UTRAN
  ebi_t ebi;
  ebi_t linked_ebi;

  // S-GW IP address for S1-u: IP address of the S-GW for the S1-u interfaces.
  // S-GW TEID for S1u: Tunnel Endpoint Identifier of the S-GW for the S1-u
  // interface.
  fteid_t s_gw_fteid_s1u;  // set by S11 CREATE_SESSION_RESPONSE

  // PDN GW TEID for S5/S8 (user plane): P-GW Tunnel Endpoint Identifier for the
  // S5/S8 interface for the user plane. (Used for S-GW change only). NOTE: The
  // PDN GW TEID is needed in MME context as S-GW relocation is triggered
  // without interaction with the source S-GW, e.g. when a TAU occurs. The
  // Target S-GW requires this Information Element, so it must be stored by the
  // MME. PDN GW IP address for S5/S8 (user plane): P GW IP address for user
  // plane for the S5/S8 interface for the user plane. (Used for S-GW change
  // only). NOTE: The PDN GW IP address for user plane is needed in MME context
  // as S-GW relocation is triggered without interaction with the source S-GW,
  // e.g. when a TAU occurs. The Target S GW requires this Information Element,
  // so it must be stored by the MME.
  fteid_t p_gw_fteid_s5_s8_up;

  // EPS bearer QoS: QCI and ARP, optionally: GBR and MBR for GBR bearer

  // extra 23.401 spec members
  pdn_cid_t pdn_cx_id;

  /*
   * Two bearer states, one mme_app_bearer_state (towards SAE-GW) and one
   * towards eNodeB (if activated in RAN). todo: setting one, based on the other
   * is possible?
   */
  mme_app_bearer_state_t
      bearer_state; /**< Need bearer state to establish them. */
  esm_ebr_context_t
      esm_ebr_context; /**< Contains the bearer level QoS parameters. */
  fteid_t enb_fteid_s1u;

  /* QoS for this bearer */
  bearer_qos_t bearer_level_qos;

  /** Add an entry field to make it part of a list (session or UE, no need to
   * save more lists). */
  // LIST_ENTRY(bearer_context_new_s) 	entries;
  STAILQ_ENTRY(bearer_context_new_s) entries;
} __attribute__((__packed__)) bearer_context_new_t;

bearer_context_t* mme_app_create_bearer_context(
    ue_mm_context_t* const ue_mm_context, const pdn_cid_t pdn_cid,
    const ebi_t ebi, const bool is_default);
void mme_app_free_bearer_context(bearer_context_t** const bearer_context);
bearer_context_t* mme_app_get_bearer_context(
    ue_mm_context_t* const ue_context, const ebi_t ebi);
void mme_app_add_bearer_context(
    ue_mm_context_t* const ue_context, bearer_context_t* const bc,
    const pdn_cid_t pdn_cid, const bool is_default);
ebi_t mme_app_get_free_bearer_id(ue_mm_context_t* const ue_context);
void mme_app_bearer_context_s1_release_enb_informations(
    bearer_context_t* const bc);

#endif
