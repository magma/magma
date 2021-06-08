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
/*! \file gtpv1_u_messages_types.h
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#ifndef FILE_GTPV1_U_MESSAGES_TYPES_SEEN
#define FILE_GTPV1_U_MESSAGES_TYPES_SEEN

#include "common_types.h"
#include "3gpp_24.007.h"

typedef struct Gtpv1uCreateTunnelReq_s {
  teid_t context_teid;  ///< Tunnel Endpoint Identifier
  ebi_t eps_bearer_id;
} Gtpv1uCreateTunnelReq;

typedef struct Gtpv1uCreateTunnelResp_s {
  uint8_t status;       ///< Status of S1U endpoint creation (Failed = 0xFF or
                        ///< Success = 0x0)
  teid_t context_teid;  ///< local SGW S11 Tunnel Endpoint Identifier
  teid_t S1u_teid;      ///< Tunnel Endpoint Identifier
  ebi_t eps_bearer_id;
} Gtpv1uCreateTunnelResp;

typedef struct Gtpv1uUpdateTunnelReq_s {
  teid_t context_teid;  ///< S11 Tunnel Endpoint Identifier
  teid_t sgw_S1u_teid;  ///< SGW S1U local Tunnel Endpoint Identifier
  teid_t enb_S1u_teid;  ///< eNB S1U Tunnel Endpoint Identifier
  ip_address_t enb_ip_address_for_S1u;
  ebi_t eps_bearer_id;
} Gtpv1uUpdateTunnelReq;

typedef struct Gtpv1uUpdateTunnelResp_s {
  uint8_t status;       ///< Status (Failed = 0xFF or Success = 0x0)
  teid_t context_teid;  ///< S11 Tunnel Endpoint Identifier
  teid_t sgw_S1u_teid;  ///< SGW S1U local Tunnel Endpoint Identifier
  teid_t enb_S1u_teid;  ///< eNB S1U Tunnel Endpoint Identifier
  ebi_t eps_bearer_id;
} Gtpv1uUpdateTunnelResp;

typedef struct Gtpv1uDeleteTunnelReq_s {
  teid_t context_teid;  ///< local SGW S11 Tunnel Endpoint Identifier
  teid_t S1u_teid;      ///< local S1U Tunnel Endpoint Identifier to be deleted
} Gtpv1uDeleteTunnelReq;

typedef struct Gtpv1uDeleteTunnelResp_s {
  uint8_t status;       ///< Status of S1U endpoint deleteion (Failed = 0xFF or
                        ///< Success = 0x0)
  teid_t context_teid;  ///< local SGW S11 Tunnel Endpoint Identifier
  teid_t S1u_teid;      ///< local S1U Tunnel Endpoint Identifier to be deleted
} Gtpv1uDeleteTunnelResp;

typedef struct Gtpv1uTunnelDataInd_s {
  uint8_t* buffer;
  uint32_t length;
  uint32_t offset;        ///< start of message offset in buffer
  teid_t local_S1u_teid;  ///< Tunnel Endpoint Identifier
} Gtpv1uTunnelDataInd;

typedef struct Gtpv1uTunnelDataReq_s {
  uint8_t* buffer;
  uint32_t length;
  uint32_t offset;        ///< start of message offset in buffer
  teid_t local_S1u_teid;  ///< Tunnel Endpoint Identifier
  teid_t S1u_enb_teid;    ///< Tunnel Endpoint Identifier
} Gtpv1uTunnelDataReq;

#endif /* FILE_GTPV1_U_MESSAGES_TYPES_SEEN */
