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

#ifndef MS_NETWORK_FEATURE_SUPPORT_H_
#define MS_NETWORK_FEATURE_SUPPORT_H_
#include <stdint.h>

#define MS_NETWORK_FEATURE_SUPPORT_MINIMUM_LENGTH 3
#define MS_NETWORK_FEATURE_SUPPORT_MAXIMUM_LENGTH 10

typedef struct MsNetworkFeatureSupport_tag {
  uint8_t spare_bits : 3;
  uint8_t extended_periodic_timers : 1;
} MsNetworkFeatureSupport;

int encode_ms_network_feature_support(
    MsNetworkFeatureSupport* msnetworkfeaturesupport, uint8_t iei,
    uint8_t* buffer, uint32_t len);

int decode_ms_network_feature_support(
    MsNetworkFeatureSupport* msnetworkfeaturesupport, uint8_t iei,
    uint8_t* buffer, uint32_t len);

void dump_ms_network_feature_support_xml(
    MsNetworkFeatureSupport* msnetworkfeaturesupport, uint8_t iei);

#endif /* MS NETWORK CAPABILITY_H_ */
