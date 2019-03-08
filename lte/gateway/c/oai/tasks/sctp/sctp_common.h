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

/*! \file sctp_common.h
 *  \brief MME SCTP related common procedures
 *  \author Sebastien ROUX, Lionel GAUTHIER
 *  \date 2013
 *  \version 1.0
 *  @ingroup _sctp
 */

#ifndef FILE_SCTP_COMMON_SEEN
#define FILE_SCTP_COMMON_SEEN

#include <stdint.h>

#include "common_types.h"

struct sockaddr;

int sctp_set_init_opt(
  const int sd,
  const sctp_stream_id_t instreams,
  const sctp_stream_id_t outstreams,
  const uint16_t max_attempts,
  const uint16_t init_timeout);

int sctp_get_sockinfo(
  int sock,
  sctp_stream_id_t *instream,
  sctp_stream_id_t *outstream,
  sctp_assoc_id_t *assoc_id);

int sctp_get_peeraddresses(
  int sock,
  struct sockaddr **remote_addr,
  int *nb_remote_addresses);

int sctp_get_localaddresses(
  int sock,
  struct sockaddr **local_addr,
  int *nb_local_addresses);

#endif /* FILE_SCTP_COMMON_SEEN */
