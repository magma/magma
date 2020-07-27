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

#ifndef EMM_CAUSE_SEEN
#define EMM_CAUSE_SEEN

#include <stdint.h>

#define EMM_CAUSE_MINIMUM_LENGTH 1
#define EMM_CAUSE_MAXIMUM_LENGTH 1

#define EMM_CAUSE_UE_SECURITY_CAP_MISMATCH 23

typedef uint8_t emm_cause_t;

int encode_emm_cause(
    emm_cause_t* emmcause, uint8_t iei, uint8_t* buffer, uint32_t len);

int decode_emm_cause(
    emm_cause_t* emmcause, uint8_t iei, uint8_t* buffer, uint32_t len);

#endif /* EMM CAUSE_SEEN */
