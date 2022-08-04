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

Subsystem EPS Session Management

Author    Frederic Maurel, Sebastien Roux

Description Defines identifiers of the EPS Session Management messages

*****************************************************************************/
#ifndef __ESM_MSGDEF_H__
#define __ESM_MSGDEF_H__

#include <asm/byteorder.h>
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.007.h"
/****************************************************************************/
/*********************  G L O B A L    C O N S T A N T S  *******************/
/****************************************************************************/

/* Header length boundaries of EPS Session Management messages  */
#define ESM_HEADER_LENGTH sizeof(esm_msg_header_t)
#define ESM_HEADER_MINIMUM_LENGTH ESM_HEADER_LENGTH
#define ESM_HEADER_MAXIMUM_LENGTH ESM_HEADER_LENGTH

/****************************************************************************/
/************************  G L O B A L    T Y P E S  ************************/
/****************************************************************************/

/*
 * Header of EPS Session Management plain NAS message
 * --------------------------------------------------
 *   8     7      6      5     4      3      2      1
 *  +-----------------------+------------------------+
 *  | EPS bearer identity | Protocol discriminator |
 *  +-----------------------+------------------------+
 *  | Procedure transaction identity     |
 *  +-----------------------+------------------------+
 *  |     Message type       |
 *  +-----------------------+------------------------+
 */
typedef struct {
#ifdef __LITTLE_ENDIAN_BITFIELD
  uint8_t protocol_discriminator : 4;
  ebi_t eps_bearer_identity : 4;
#endif
#ifdef __BIG_ENDIAN_BITFIELD
  ebi_t eps_bearer_identity : 4;
  uint8_t protocol_discriminator : 4;
#endif
  pti_t procedure_transaction_identity;
  uint8_t message_type;
} __attribute__((__packed__)) esm_msg_header_t;

/****************************************************************************/
/********************  G L O B A L    V A R I A B L E S  ********************/
/****************************************************************************/

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

#endif /* __ESM_MSGDEF_H__ */
