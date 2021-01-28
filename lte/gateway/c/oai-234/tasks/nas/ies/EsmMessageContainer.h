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

#ifndef ESM_MESSAGE_CONTAINER_SEEN
#define ESM_MESSAGE_CONTAINER_SEEN
#include <stdint.h>

#include "bstrlib.h"

struct scenario_s;
struct scenario_player_msg_s;

#define ESM_MESSAGE_CONTAINER_MINIMUM_LENGTH 2  // [length]+[length]
#define ESM_MESSAGE_CONTAINER_MAXIMUM_LENGTH                                   \
  65538  // [IEI]+[length]+[length]+[ESM msg]

typedef bstring EsmMessageContainer;

int encode_esm_message_container(
    EsmMessageContainer esmmessagecontainer, uint8_t iei, uint8_t* buffer,
    uint32_t len);

int decode_esm_message_container(
    EsmMessageContainer* esmmessagecontainer, uint8_t iei, uint8_t* buffer,
    uint32_t len);

#endif /* ESM_MESSAGE_CONTAINER_SEEN */
