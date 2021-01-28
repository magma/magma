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

#ifndef KSI_AND_SEQUENCE_NUMBER_SEEN
#define KSI_AND_SEQUENCE_NUMBER_SEEN

#include <stdint.h>

#define KSI_AND_SEQUENCE_NUMBER_MINIMUM_LENGTH 2
#define KSI_AND_SEQUENCE_NUMBER_MAXIMUM_LENGTH 2

typedef struct KsiAndSequenceNumber_tag {
  uint8_t ksi : 3;
  uint8_t sequencenumber : 5;
} KsiAndSequenceNumber;

int encode_ksi_and_sequence_number(
    KsiAndSequenceNumber* ksiandsequencenumber, uint8_t iei, uint8_t* buffer,
    uint32_t len);

int decode_ksi_and_sequence_number(
    KsiAndSequenceNumber* ksiandsequencenumber, uint8_t iei, uint8_t* buffer,
    uint32_t len);

#endif /* KSI AND SEQUENCE NUMBER_H_ */
