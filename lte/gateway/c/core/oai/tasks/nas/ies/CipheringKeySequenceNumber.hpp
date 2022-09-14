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

#ifndef CIPHERING_KEY_SEQUENCE_NUMBER_H_
#define CIPHERING_KEY_SEQUENCE_NUMBER_H_
#include <stdint.h>

#define CIPHERING_KEY_SEQUENCE_NUMBER_MINIMUM_LENGTH 1
#define CIPHERING_KEY_SEQUENCE_NUMBER_MAXIMUM_LENGTH 1

typedef uint8_t CipheringKeySequenceNumber;

int encode_ciphering_key_sequence_number(
    CipheringKeySequenceNumber* cipheringkeysequencenumber, uint8_t iei,
    uint8_t* buffer, uint32_t len);

void dump_ciphering_key_sequence_number_xml(
    CipheringKeySequenceNumber* cipheringkeysequencenumber, uint8_t iei);

uint8_t encode_u8_ciphering_key_sequence_number(
    CipheringKeySequenceNumber* cipheringkeysequencenumber);

int decode_ciphering_key_sequence_number(
    CipheringKeySequenceNumber* cipheringkeysequencenumber, uint8_t iei,
    uint8_t* buffer, uint32_t len);

int decode_u8_ciphering_key_sequence_number(
    CipheringKeySequenceNumber* cipheringkeysequencenumber, uint8_t iei,
    uint8_t value, uint32_t len);

#endif /* CIPHERING KEY SEQUENCE NUMBER_H_ */
