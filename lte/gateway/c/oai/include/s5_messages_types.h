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
#ifndef FILE_S5_MESSAGES_TYPES_SEEN
#define FILE_S5_MESSAGES_TYPES_SEEN

#include "sgw_ie_defs.h"

#define S5_CREATE_BEARER_REQUEST(mSGpTR)                                       \
  (mSGpTR)->ittiMsg.s5_create_bearer_request
#define S5_CREATE_BEARER_RESPONSE(mSGpTR)                                      \
  (mSGpTR)->ittiMsg.s5_create_bearer_response

typedef struct itti_s5_create_bearer_request_s {
  teid_t context_teid; ///< local SGW S11 Tunnel Endpoint Identifier
  teid_t S1u_teid;     ///< Tunnel Endpoint Identifier
  ebi_t eps_bearer_id;
} itti_s5_create_bearer_request_t;

enum s5_failure_cause { S5_OK = 0, PCEF_FAILURE };

typedef struct itti_s5_create_bearer_response_s {
  teid_t context_teid; ///< local SGW S11 Tunnel Endpoint Identifier
  teid_t S1u_teid;     ///< Tunnel Endpoint Identifier
  ebi_t eps_bearer_id;
  itti_sgi_create_end_point_response_t sgi_create_endpoint_resp;
  enum s5_failure_cause failure_cause;
} itti_s5_create_bearer_response_t;

#endif /* FILE_S5_MESSAGES_TYPES_SEEN*/
