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
/*! \file sctp_messages_types.h
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#ifndef FILE_SCTP_MESSAGES_TYPES_SEEN
#define FILE_SCTP_MESSAGES_TYPES_SEEN

#include <arpa/inet.h>
#include <stdint.h>

#include "bstrlib.h"

#include "common_types.h"

typedef uint32_t sctp_ppid_t;

#define SCTP_DATA_IND(msg) (msg)->ittiMsg.sctp_data_ind
#define SCTP_DATA_REQ(msg) (msg)->ittiMsg.sctp_data_req
#define SCTP_DATA_CNF(msg) (msg)->ittiMsg.sctp_data_cnf
#define SCTP_INIT_MSG(msg) (msg)->ittiMsg.sctpInit
#define SCTP_NEW_ASSOCIATION(msg) (msg)->ittiMsg.sctp_new_peer
#define SCTP_CLOSE_ASSOCIATION(msg) (msg)->ittiMsg.sctp_close_association
#define SCTP_MME_SERVER_INITIALIZED(msg)                                       \
  (msg)->ittiMsg.sctp_mme_server_initialized
#define SCTP_AMF_SERVER_INITIALIZED(msg)                                       \
  (msg)->ittiMsg.sctp_amf_server_initialized

typedef struct sctp_data_cnf_s {
  bstring payload;
  sctp_assoc_id_t assoc_id;
  sctp_stream_id_t stream;
  uint32_t mme_ue_s1ap_id;
  uint32_t amf_ue_ngap_id;
  sctp_ppid_t ppid;
  bool is_success;
} sctp_data_cnf_t;

typedef struct sctp_data_req_s {
  bstring payload;
  sctp_assoc_id_t assoc_id;
  sctp_stream_id_t stream;
  uint32_t mme_ue_s1ap_id;  // for helping data_rej
  uint32_t amf_ue_ngap_id;
  sctp_ppid_t ppid;
} sctp_data_req_t;

typedef struct sctp_data_ind_s {
  bstring payload;           ///< SCTP buffer
  sctp_assoc_id_t assoc_id;  ///< SCTP physical association ID
  sctp_stream_id_t stream;   ///< Stream number on which data had been received
  uint16_t instreams;   ///< Number of input streams for the SCTP connection
                        ///< between peers
  uint16_t outstreams;  ///< Number of output streams for the SCTP connection
                        ///< between peers
} sctp_data_ind_t;

typedef struct sctp_init_s {
  /* Request usage of ipv4 */
  unsigned ipv4 : 1;
  /* Request usage of ipv6 */
  unsigned ipv6 : 1;
  uint8_t nb_ipv4_addr;
  struct in_addr ipv4_address[10];
  uint8_t nb_ipv6_addr;
  struct in6_addr ipv6_address[10];
  uint16_t port;
  uint32_t ppid;
} sctp_init_t;

typedef struct sctp_close_association_s {
  sctp_assoc_id_t assoc_id;
  // True if the association is being closed down because of a reset.
  bool reset;
} sctp_close_association_t;

typedef struct sctp_new_peer_s {
  uint32_t instreams;
  uint32_t outstreams;
  sctp_assoc_id_t assoc_id;
  bstring ran_cp_ipaddr;
} sctp_new_peer_t;

typedef struct sctp_mme_server_initialized_s {
  bool successful;
} sctp_mme_server_initialized_t;

typedef struct sctp_amf_server_initialized_s {
  bool successful;
} sctp_amf_server_initialized_t;
#endif /* FILE_SCTP_MESSAGES_TYPES_SEEN */
