/*
  Copyright 2020 The Magma Authors.
  This source code is licensed under the BSD-style license found in the
  LICENSE file in the root directory of this source tree.
  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

#include <sstream>
#include <cstdint>
#include <cstring>
#include "QOSRules.h"
#include "CommonDefs.h"

using namespace std;
namespace magma5g {
NewQOSRulePktFilter::NewQOSRulePktFilter(){};
NewQOSRulePktFilter::~NewQOSRulePktFilter(){};
QOSRule::QOSRule(){};
QOSRulesMsg::QOSRulesMsg(){};
QOSRule::~QOSRule(){};
QOSRulesMsg::~QOSRulesMsg(){};

// Decode QOSRules IE
int QOSRulesMsg::DecodeQOSRulesMsg(
    QOSRulesMsg* qosrulesmsg, uint8_t iei, uint8_t* buffer, uint32_t len) {
  // Not yet Implemented, will be suppported POST MVC
  return (0);
};

// Encode QOSRules IE
int QOSRulesMsg::EncodeQOSRulesMsg(
    QOSRulesMsg* qosrulesmsg, uint8_t iei, uint8_t* buffer, uint32_t len) {
  uint16_t encoded = 0;
  uint8_t i        = 0;
  uint8_t j        = 0;

  // Checking IEI and pointer
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(buffer, QOSRULE_MIN_LEN, len);

  if (iei > 0) {
    CHECK_IEI_ENCODER((unsigned char) iei, qosrulesmsg->iei);
    *buffer = iei;
    MLOG(MDEBUG) << "In EncodeQOSRulesMsg: iei" << hex << int(*buffer) << endl;
    encoded++;
  }

  IES_ENCODE_U16(buffer, encoded, qosrulesmsg->length);
  MLOG(MDEBUG) << "Length : " << hex << int(qosrulesmsg->length) << endl;
  while (encoded < (qosrulesmsg->length) && i <= 255) {
    *(buffer + encoded) = qosrulesmsg->qosrule[i].qosruleid;
    MLOG(MDEBUG) << "qosruleid: " << hex << int(*(buffer + encoded)) << endl;
    encoded++;
    IES_ENCODE_U16(buffer, encoded, qosrulesmsg->qosrule[i].len);
    *(buffer + encoded) = 0x00 |
                          ((qosrulesmsg->qosrule[i].ruleopercode & 0x07) << 5) |
                          ((qosrulesmsg->qosrule[i].dqrbit & 0x01) << 4) |
                          (qosrulesmsg->qosrule[i].noofpktfilters & 0x0f);
    MLOG(MDEBUG) << "ruleopercode, dqrbit, noofpktfilters: " << hex
                 << int(*(buffer + encoded)) << endl;
    encoded++;
    for (j = 0; j < qosrulesmsg->qosrule[i].noofpktfilters; j++) {
      *(buffer + encoded) =
          0x00 |
          ((qosrulesmsg->qosrule[i].newqosrulepktfilter[j].spare & 0x03) << 6) |
          ((qosrulesmsg->qosrule[i].newqosrulepktfilter[j].pktfilterdir & 0x03)
           << 4) |
          (qosrulesmsg->qosrule[i].newqosrulepktfilter[j].pktfilterid & 0x0f);
      MLOG(MDEBUG) << "pktfilterdir, pktfilterid: " << hex
                   << int(*(buffer + encoded)) << endl;
      encoded++;
      *(buffer + encoded) = qosrulesmsg->qosrule[i].newqosrulepktfilter[j].len;
      MLOG(MDEBUG) << "len: " << hex << int(*(buffer + encoded)) << endl;
      encoded++;
      memcpy(
          buffer + encoded,
          qosrulesmsg->qosrule[i].newqosrulepktfilter[j].contents,
          qosrulesmsg->qosrule[i].newqosrulepktfilter[j].len);
      BUFFER_PRINT_LOG(
          buffer + encoded, qosrulesmsg->qosrule[i].newqosrulepktfilter[j].len);
      encoded = encoded + qosrulesmsg->qosrule[i].newqosrulepktfilter[j].len;
      encoded++;
    }

    *(buffer + encoded) = qosrulesmsg->qosrule[i].qosruleprecedence;
    MLOG(MDEBUG) << "qosruleprecedence: " << hex << int(*(buffer + encoded))
                 << endl;
    encoded++;
    *(buffer + encoded) = 0x00 | ((qosrulesmsg->qosrule[i].spare & 0x01) << 7) |
                          ((qosrulesmsg->qosrule[i].segregation & 0x01) << 6) |
                          (qosrulesmsg->qosrule[i].qfi & 0x3f);
    MLOG(MDEBUG) << "segregation, qfi: " << hex << int(*(buffer + encoded))
                 << endl;
    encoded++;
    i++;
  }

  return (encoded);
};
}  // namespace magma5g
