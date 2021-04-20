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

#ifndef NONCE_SEEN
#define NONCE_SEEN

#include <stdint.h>

#define NONCE_MINIMUM_LENGTH 5
#define NONCE_MAXIMUM_LENGTH 5

typedef uint32_t nonce_t;

int encode_nonce(nonce_t* nonce, uint8_t iei, uint8_t* buffer, uint32_t len);

int decode_nonce(nonce_t* nonce, uint8_t iei, uint8_t* buffer, uint32_t len);

#endif /* NONCE_SEEN */
