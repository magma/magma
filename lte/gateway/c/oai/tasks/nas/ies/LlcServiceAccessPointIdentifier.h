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

#ifndef LLC_SERVICE_ACCESS_POINT_IDENTIFIER_H_
#define LLC_SERVICE_ACCESS_POINT_IDENTIFIER_H_
#include <stdint.h>

#define LLC_SERVICE_ACCESS_POINT_IDENTIFIER_MINIMUM_LENGTH 2
#define LLC_SERVICE_ACCESS_POINT_IDENTIFIER_MAXIMUM_LENGTH 2

typedef uint8_t LlcServiceAccessPointIdentifier;

int encode_llc_service_access_point_identifier(
    LlcServiceAccessPointIdentifier* llcserviceaccesspointidentifier,
    uint8_t iei, uint8_t* buffer, uint32_t len);

void dump_llc_service_access_point_identifier_xml(
    LlcServiceAccessPointIdentifier* llcserviceaccesspointidentifier,
    uint8_t iei);

int decode_llc_service_access_point_identifier(
    LlcServiceAccessPointIdentifier* llcserviceaccesspointidentifier,
    uint8_t iei, uint8_t* buffer, uint32_t len);

#endif /* LLC SERVICE ACCESS POINT IDENTIFIER_H_ */
