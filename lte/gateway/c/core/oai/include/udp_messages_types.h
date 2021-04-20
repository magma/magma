/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
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

/*! \file udp_messages_types.h
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#ifndef FILE_UDP_MESSAGES_TYPES_SEEN
#define FILE_UDP_MESSAGES_TYPES_SEEN

#define UDP_INIT(mSGpTR) (mSGpTR)->ittiMsg.udp_init
#define UDP_DATA_MAX_MSG_LEN                                                   \
  (4096) /**< Maximum supported gtpv2c packet length including header */

typedef struct {
  struct in_addr* in_addr;
  struct in6_addr* in6_addr;
  uint16_t port;
} udp_init_t;

typedef struct {
  uint8_t* buffer;
  uint32_t buffer_length;
  uint32_t buffer_offset;
  uint16_t local_port;
  struct sockaddr* peer_address;
  uint16_t peer_port;
} udp_data_req_t;

typedef struct {
  uint8_t msgBuf[UDP_DATA_MAX_MSG_LEN];
  uint32_t buffer_length;
  uint16_t local_port;
  union {
    struct sockaddr_in addrv4;
    struct sockaddr_in6 addrv6;
  } sock_addr;
  uint16_t peer_port;
} udp_data_ind_t;

#endif /* FILE_UDP_MESSAGES_TYPES_SEEN */
