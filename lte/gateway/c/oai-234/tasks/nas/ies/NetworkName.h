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

#ifndef NETWORK_NAME_H_
#define NETWORK_NAME_H_
#include <stdint.h>
#include "bstrlib.h"

#define NETWORK_NAME_MINIMUM_LENGTH 3
#define NETWORK_NAME_MAXIMUM_LENGTH 255

typedef struct NetworkName_tag {
  uint8_t codingscheme : 3;
  uint8_t addci : 1;
  uint8_t numberofsparebitsinlastoctet : 3;
  bstring textstring;
} NetworkName;

int encode_network_name(
    NetworkName* networkname, uint8_t iei, uint8_t* buffer, uint32_t len);

int decode_network_name(
    NetworkName* networkname, uint8_t iei, uint8_t* buffer, uint32_t len);

void dump_network_name_xml(NetworkName* networkname, uint8_t iei);

#endif /* NETWORK NAME_H_ */
