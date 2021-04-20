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

#ifndef DRX_PARAMETER_H_
#define DRX_PARAMETER_H_
#include <stdint.h>

#define DRX_PARAMETER_MINIMUM_LENGTH 3
#define DRX_PARAMETER_MAXIMUM_LENGTH 3

typedef struct DrxParameter_tag {
  uint8_t splitpgcyclecode;
  uint8_t cnspecificdrxcyclelengthcoefficientanddrxvaluefors1mode : 4;
  uint8_t splitonccch : 1;
  uint8_t nondrxtimer : 3;
} DrxParameter;

int encode_drx_parameter(
    DrxParameter* drxparameter, uint8_t iei, uint8_t* buffer, uint32_t len);

void dump_drx_parameter_xml(DrxParameter* drxparameter, uint8_t iei);

int decode_drx_parameter(
    DrxParameter* drxparameter, uint8_t iei, uint8_t* buffer, uint32_t len);

#endif /* DRX PARAMETER_H_ */
