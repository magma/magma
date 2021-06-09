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

#ifndef TRAFFIC_FLOW_AGGREGATE_DESCRIPTION_SEEN
#define TRAFFIC_FLOW_AGGREGATE_DESCRIPTION_SEEN

#define TRAFFIC_FLOW_AGGREGATE_DESCRIPTION_MINIMUM_LENGTH 1
#define TRAFFIC_FLOW_AGGREGATE_DESCRIPTION_MAXIMUM_LENGTH 1

#include "3gpp_24.008.h"
typedef traffic_flow_template_t traffic_flow_aggregate_description_t;

int encode_traffic_flow_aggregate_description(
    traffic_flow_aggregate_description_t* trafficflowaggregatedescription,
    uint8_t iei, uint8_t* buffer, uint32_t len);

int decode_traffic_flow_aggregate_description(
    traffic_flow_aggregate_description_t* trafficflowaggregatedescription,
    uint8_t iei, uint8_t* buffer, uint32_t len);

#endif /* TRAFFIC FLOW AGGREGATE DESCRIPTION_SEEN */
