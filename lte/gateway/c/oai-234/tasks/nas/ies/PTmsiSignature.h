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

#ifndef P_TMSI_SIGNATURE_H_
#define P_TMSI_SIGNATURE_H_
#include <stdint.h>
#include "bstrlib.h"

#define P_TMSI_SIGNATURE_MINIMUM_LENGTH 4
#define P_TMSI_SIGNATURE_MAXIMUM_LENGTH 4

typedef bstring PTmsiSignature;

int encode_p_tmsi_signature(
    PTmsiSignature ptmsisignature, uint8_t iei, uint8_t* buffer, uint32_t len);

void dump_p_tmsi_signature_xml(PTmsiSignature ptmsisignature, uint8_t iei);

int decode_p_tmsi_signature(
    PTmsiSignature* ptmsisignature, uint8_t iei, uint8_t* buffer, uint32_t len);

#endif /* P TMSI SIGNATURE_H_ */
