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
/*! \file udp_messages_types.h
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#ifndef FILE_UDP_MESSAGES_TYPES_SEEN
#define FILE_UDP_MESSAGES_TYPES_SEEN

#include <stdint.h>
#include <netinet/in.h>

#define UDP_INIT(msg_ptr) (msg_ptr)->ittiMsg.udp_init
#define UDP_DATA_REQ(msg_ptr) (msg_ptr)->ittiMsg.udp_data_req
#define UDP_DATA_IND(msg_ptr) (msg_ptr)->ittiMsg.udp_data_ind

typedef struct udp_init {
  struct in_addr address;
  uint16_t port;
} udp_init_t;

typedef struct udp_data_req {
  uint8_t *buffer;
  uint32_t buffer_length;
  uint32_t buffer_offset;
  struct in_addr peer_address;
  uint16_t peer_port;
} udp_data_req_t;

typedef struct udp_data_ind {
  uint8_t *buffer;
  uint32_t buffer_length;
  struct in_addr peer_address;
  uint16_t peer_port;
} udp_data_ind_t;

#endif /* FILE_UDP_MESSAGES_TYPES_SEEN */
