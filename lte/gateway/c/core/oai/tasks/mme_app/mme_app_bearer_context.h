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
