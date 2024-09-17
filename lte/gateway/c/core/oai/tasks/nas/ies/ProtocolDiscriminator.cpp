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

#include "lte/gateway/c/core/oai/tasks/nas/ies/ProtocolDiscriminator.hpp"

#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/common/assertions.h"
#include "lte/gateway/c/core/oai/common/TLVDecoder.h"
#include "lte/gateway/c/core/oai/common/TLVEncoder.h"
#ifdef __cplusplus
}
#endif

int decode_protocol_discriminator(ProtocolDiscriminator* protocoldiscriminator,
                                  uint8_t iei, uint8_t* buffer, uint32_t len) {
  Fatal("TODO Implement decode_protocol_discriminator");
  return -1;
}

int encode_protocol_discriminator(ProtocolDiscriminator* protocoldiscriminator,
                                  uint8_t iei, uint8_t* buffer, uint32_t len) {
  Fatal("TODO Implement encode_protocol_discriminator");
  return -1;
}
