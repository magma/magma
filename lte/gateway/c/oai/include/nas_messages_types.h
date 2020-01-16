/*
 * Copyright (c) 2015, EURECOM (www.eurecom.fr)
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice, this
 *    list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
 * ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 * The views and conclusions contained in the software and documentation are those
 * of the authors and should not be interpreted as representing official policies,
 * either expressed or implied, of the FreeBSD Project.
 */

/*! \file nas_messages_types.h
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#ifndef FILE_NAS_MESSAGES_TYPES_SEEN
#define FILE_NAS_MESSAGES_TYPES_SEEN

#include <stdint.h>

#include "3gpp_23.003.h"
#include "3gpp_29.274.h"
#include "nas/as_message.h"
#include "common_ies.h"
#include "nas/networkDef.h"

#define NAS_BEARER_PARAM(mSGpTR) (mSGpTR)->ittiMsg.nas_bearer_param
#define NAS_AUTHENTICATION_REQ(mSGpTR) (mSGpTR)->ittiMsg.nas_auth_req
#define NAS_AUTHENTICATION_PARAM_REQ(mSGpTR)                                   \
  (mSGpTR)->ittiMsg.nas_auth_param_req
#define NAS_DATA_LENGHT_MAX 256

typedef struct itti_nas_info_transfer_s {
  mme_ue_s1ap_id_t ue_id; /* UE lower layer identifier        */
  //nas_error_code_t err_code;     /* Transaction status               */
  bstring nas_msg; /* Uplink NAS message           */
} itti_nas_info_transfer_t;

typedef struct itti_nas_attach_req_s {
  /* TODO: Set the correct size */
  char apn[100];
  char imsi[16];
#define INITIAL_REQUEST (0x1)
  unsigned initial : 1;
  s1ap_initial_ue_message_t transparent;
} itti_nas_attach_req_t;

typedef struct itti_nas_auth_req_s {
  /* UE imsi */
  char imsi[16];

#define NAS_FAILURE_OK 0x0
#define NAS_FAILURE_IND 0x1
  unsigned failure : 1;
  int cause;
} itti_nas_auth_req_t;

typedef struct itti_nas_auth_rsp_s {
  char imsi[16];
} itti_nas_auth_rsp_t;

typedef struct itti_nas_auth_param_req_s {
  /* UE identifier */
  mme_ue_s1ap_id_t ue_id;

  /* Imsi of the UE (In case of initial request) */
  char imsi[16];
  uint8_t imsi_length;

  /* Indicates whether the procedure corresponds to a new connection or not */
  uint8_t initial_req : 1;

  uint8_t re_synchronization : 1;
  uint8_t auts[14];
  uint8_t num_vectors;
} itti_nas_auth_param_req_t;

#endif /* FILE_NAS_MESSAGES_TYPES_SEEN */
