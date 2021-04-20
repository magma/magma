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

/*****************************************************************************
Version   0.1

Date    2012/09/27

Product   NAS stack

Subsystem EPS Mobility Management

Author    Frederic Maurel, Sebastien Roux

Description Defines identifiers of the EPS Mobility Management messages.

*****************************************************************************/
#ifndef FILE_EMM_MSGDEF_SEEN
#define FILE_EMM_MSGDEF_SEEN

#include <asm/byteorder.h>

/****************************************************************************/
/*********************  G L O B A L    C O N S T A N T S  *******************/
/****************************************************************************/

/* Header length boundaries of EPS Mobility Management messages  */
#define EMM_HEADER_LENGTH sizeof(emm_msg_header_t)
#define EMM_HEADER_MINIMUM_LENGTH EMM_HEADER_LENGTH
#define EMM_HEADER_MAXIMUM_LENGTH EMM_HEADER_LENGTH

/*
 * Message identifiers for EMM messages that does not follow the structure
 * of a standard layer 3 message
 */
#define SERVICE_REQUEST 0b01001101 /* TODO: TBD - 77 = 0x4d */

/****************************************************************************/
/************************  G L O B A L    T Y P E S  ************************/
/****************************************************************************/

/*
 * Header of EPS Mobility Management plain NAS message
 * ---------------------------------------------------
 *   8     7      6      5     4      3      2      1
 *  +-----------------------+------------------------+
 *  | Security header type  | Protocol discriminator |
 *  +-----------------------+------------------------+
 *  |     Message type       |
 *  +-----------------------+------------------------+
 */
typedef struct emm_msg_header_s {
#ifdef __LITTLE_ENDIAN_BITFIELD
  uint8_t protocol_discriminator : 4;
  uint8_t security_header_type : 4;
#endif
#ifdef __BIG_ENDIAN_BITFIELD
  uint8_t security_header_type : 4;
  uint8_t protocol_discriminator : 4;
#endif
  uint8_t message_type;
} __attribute__((__packed__)) emm_msg_header_t;

/****************************************************************************/
/********************  G L O B A L    V A R I A B L E S  ********************/
/****************************************************************************/

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

#endif /* FILE_EMM_MSGDEF_SEEN */
