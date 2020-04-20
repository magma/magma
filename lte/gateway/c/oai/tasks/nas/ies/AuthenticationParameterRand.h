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

#ifndef AUTHENTICATION_PARAMETER_RAND_H_
#define AUTHENTICATION_PARAMETER_RAND_H_
#include <stdint.h>
#include "bstrlib.h"

#define AUTHENTICATION_PARAMETER_RAND_MINIMUM_LENGTH 16
#define AUTHENTICATION_PARAMETER_RAND_MAXIMUM_LENGTH 16

typedef bstring AuthenticationParameterRand;

int encode_authentication_parameter_rand(
    AuthenticationParameterRand authenticationparameterrand, uint8_t iei,
    uint8_t* buffer, uint32_t len);

void dump_authentication_parameter_rand_xml(
    AuthenticationParameterRand authenticationparameterrand, uint8_t iei);

int decode_authentication_parameter_rand(
    AuthenticationParameterRand* authenticationparameterrand, uint8_t iei,
    uint8_t* buffer, uint32_t len);

#endif /* AUTHENTICATION PARAMETER RAND_H_ */
