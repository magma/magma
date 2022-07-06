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

/*! \file s11_causes.h
 * \brief
 * \author Sebastien ROUX, Lionel Gauthier
 * \company Eurecom
 * \email: lionel.gauthier@eurecom.fr
 */
#ifndef FILE_S11_CAUSES_SEEN
#define FILE_S11_CAUSES_SEEN

#include <stdint.h>

typedef struct SGWCauseMapping_e {
  uint8_t value;
  /* Displayable cause name */
  char* name;
  /* Possible cause in message? */
  unsigned create_session_response : 1;
  unsigned create_bearer_response : 1;
  unsigned modify_bearer_response : 1;
  unsigned delete_session_response : 1;
} SGWCauseMapping_t;

char* sgw_cause_2_string(uint8_t cause_value);

#endif /* FILE_S11_CAUSES_SEEN */
