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

#ifndef EPS_BEARER_CONTEXT_STATUS_SEEN
#define EPS_BEARER_CONTEXT_STATUS_SEEN

#include <stdint.h>

#define EPS_BEARER_CONTEXT_STATUS_MINIMUM_LENGTH 4
#define EPS_BEARER_CONTEXT_STATUS_MAXIMUM_LENGTH 4

typedef uint16_t eps_bearer_context_status_t;

int encode_eps_bearer_context_status(
    eps_bearer_context_status_t* epsbearercontextstatus, uint8_t iei,
    uint8_t* buffer, uint32_t len);

int decode_eps_bearer_context_status(
    eps_bearer_context_status_t* epsbearercontextstatus, uint8_t iei,
    uint8_t* buffer, uint32_t len);

#endif /* EPS_BEARER_CONTEXT_STATUS_SEEN */
