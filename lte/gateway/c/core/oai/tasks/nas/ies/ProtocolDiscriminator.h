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

#ifndef PROTOCOL_DISCRIMINATOR_H_
#define PROTOCOL_DISCRIMINATOR_H_
#include <stdint.h>

#define PROTOCOL_DISCRIMINATOR_MINIMUM_LENGTH 1
#define PROTOCOL_DISCRIMINATOR_MAXIMUM_LENGTH 1

typedef uint8_t ProtocolDiscriminator;

int encode_protocol_discriminator(
    ProtocolDiscriminator* protocoldiscriminator, uint8_t iei, uint8_t* buffer,
    uint32_t len);

void dump_protocol_discriminator_xml(
    ProtocolDiscriminator* protocoldiscriminator, uint8_t iei);

int decode_protocol_discriminator(
    ProtocolDiscriminator* protocoldiscriminator, uint8_t iei, uint8_t* buffer,
    uint32_t len);

#endif /* PROTOCOL DISCRIMINATOR_H_ */
