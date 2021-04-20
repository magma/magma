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

#ifndef QUALITY_OF_SERVICE_H_
#define QUALITY_OF_SERVICE_H_
#include <stdint.h>

#define QUALITY_OF_SERVICE_MINIMUM_LENGTH 14
#define QUALITY_OF_SERVICE_MAXIMUM_LENGTH 14

typedef struct QualityOfService_tag {
  uint8_t delayclass : 3;
  uint8_t reliabilityclass : 3;
  uint8_t peakthroughput : 4;
  uint8_t precedenceclass : 3;
  uint8_t meanthroughput : 5;
  uint8_t trafficclass : 3;
  uint8_t deliveryorder : 2;
  uint8_t deliveryoferroneoussdu : 3;
  uint8_t maximumsdusize;
  uint8_t maximumbitrateuplink;
  uint8_t maximumbitratedownlink;
  uint8_t residualber : 4;
  uint8_t sduratioerror : 4;
  uint8_t transferdelay : 6;
  uint8_t traffichandlingpriority : 2;
  uint8_t guaranteedbitrateuplink;
  uint8_t guaranteedbitratedownlink;
  uint8_t signalingindication : 1;
  uint8_t sourcestatisticsdescriptor : 4;
} QualityOfService;

int encode_quality_of_service(
    QualityOfService* qualityofservice, uint8_t iei, uint8_t* buffer,
    uint32_t len);

int decode_quality_of_service(
    QualityOfService* qualityofservice, uint8_t iei, uint8_t* buffer,
    uint32_t len);

void dump_quality_of_service_xml(
    QualityOfService* qualityofservice, uint8_t iei);

#endif /* QUALITY OF SERVICE_H_ */
