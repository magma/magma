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
Source      securityDef.h

Version     0.1

Date        2013/05/02

Product     NAS stack

Subsystem   include

Author      Frederic Maurel, Lionel GAUTHIER

Description Contains global security definitions

*****************************************************************************/
#ifndef FILE_SECURITYDEF_SEEN
#define FILE_SECURITYDEF_SEEN

/****************************************************************************/
/*********************  G L O B A L    C O N S T A N T S  *******************/
/****************************************************************************/

/*
 * Index of the first byte of each fields of the AUTN parameter
 */
#define AUTH_SQN_INDEX 0
#define AUTH_AMF_INDEX (AUTH_SQN_INDEX + AUTH_SQN_SIZE)
#define AUTH_MAC_INDEX (AUTH_AMF_INDEX + AUTH_AMF_SIZE)

/*
 * Size of the authentication challenge parameters in bytes
 */
#define AUTH_SQN_SIZE 6 /* Sequence number:          48 bits  */
#define AUTH_AK_SIZE 6  /* Anonymity key:            48 bits  */
#define AUTH_AMF_SIZE 2 /* Authentication Management Field:  16 bits  */
#define AUTH_MAC_SIZE 8 /* Message Authentication Code:  64 bits  */
#define AUTH_AUTN_SIZE                                                         \
  16                          /* Authentication token:     128 bits            \
                AUTN = (SQN ⊕ AK) || AMF || MAC        */
#define AUTH_MACS_SIZE 8      /* Re-synchronization MAC:       64 bits  */
#define AUTH_AUTS_SIZE 16     /* Re-synchronization AUT:       128 bits */
#define AUTH_RAND_SIZE 16     /* Random challenge:         128 bits     */
#define AUTH_CK_SIZE 16       /* Ciphering key:            128 bits     */
#define AUTH_IK_SIZE 16       /* Integrity key:            128 bits     */
#define AUTH_RES_SIZE 16      /* Authentication response:      128 bits */
#define AUTH_SNID_SIZE 3      /* Serving network's identity:   24 bits  */
#define AUTH_KASME_SIZE 32    /* KASME security key:        256 bits    */
#define AUTH_KAMF_SIZE 32     /* KAMF security key:        256 bits    */
#define AUTH_KNAS_INT_SIZE 16 /* NAS integrity key     */
#define AUTH_KNAS_ENC_SIZE 16 /* NAS cyphering key     */
#define AUTH_KENB_SIZE AUTH_KASME_SIZE     /* eNodeB security key   */
#define AUTH_NEXT_HOP_SIZE AUTH_KASME_SIZE /* Next Hop security parameter*/

#define AUTH_KGNB_SIZE AUTH_KASME_SIZE /* gNodeB security key   */

/* "Separation bit" of AMF field */
#define AUTH_AMF_SEPARATION_BIT(a) ((a) &0x80)
/*
 * --------------------------------------------------------------------------
 * EPS NAS security context handled by EPS Mobility Management sublayer in
 * the UE and in the MME
 * --------------------------------------------------------------------------
 */
/* Type of security context */
typedef enum {
  SECURITY_CTX_TYPE_NOT_AVAILABLE = 0,
  SECURITY_CTX_TYPE_PARTIAL_NATIVE,
  SECURITY_CTX_TYPE_FULL_NATIVE,
  SECURITY_CTX_TYPE_MAPPED  // UNUSED
} emm_sc_type_t;

/****************************************************************************/
/************************  G L O B A L    T Y P E S  ************************/
/****************************************************************************/

/*
 * EPS authentication vector
 */
typedef struct {
  /* ASME security key                */
  uint8_t kasme[AUTH_KASME_SIZE];
  /* Random challenge parameter           */
  uint8_t rand[AUTH_RAND_SIZE];
  /* Authentication token parameter       */
  uint8_t autn[AUTH_AUTN_SIZE];
  /* Expected Authentication response parameter   */
#define AUTH_XRES_SIZE AUTH_RES_SIZE
  uint8_t xres_size;
  uint8_t xres[AUTH_XRES_SIZE];
} auth_vector_t;

/****************************************************************************/
/********************  G L O B A L    V A R I A B L E S  ********************/
/****************************************************************************/

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

#endif /* FILE_SECURITYDEF_SEEN*/
