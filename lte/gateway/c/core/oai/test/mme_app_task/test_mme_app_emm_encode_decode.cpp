/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
#include <gtest/gtest.h>

extern "C" {
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/AttachRequest.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/EmmInformation.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/TrackingAreaUpdateRequest.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/UplinkNasTransport.h"
#include "lte/gateway/c/core/oai/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/common/log.h"
}

#include "lte/gateway/c/core/oai/test/mme_app_task/mme_app_test_util.h"

namespace magma {
namespace lte {

#define BUFFER_LEN 200

#define FILL_EMM_COMMON_MANDATORY_DEFAULTS(msg)                                \
  do {                                                                         \
    msg.protocoldiscriminator = EPS_SESSION_MANAGEMENT_MESSAGE;                \
    msg.securityheadertype    = 0x0001;                                        \
    msg.messagetype           = 33;                                            \
  } while (0)

class EMMEncodeDecodeTest : public ::testing::Test {
  virtual void SetUp() {}
  virtual void TearDown() {}

 protected:
  uint8_t temp_buffer[BUFFER_LEN];
};

TEST_F(EMMEncodeDecodeTest, TestEncodeDecodeAttachRequest) {
  //   Combined attach, NAS message generated from s1ap tester
  uint8_t buffer[] = {0x72, 0x08, 0x09, 0x10, 0x10, 0x00, 0x00, 0x00,
                      0x00, 0x10, 0x02, 0xe0, 0xe0, 0x00, 0x04, 0x02,
                      0x01, 0xd0, 0x11, 0x40, 0x08, 0x04, 0x02, 0x60,
                      0x04, 0x00, 0x02, 0x1c, 0x00};
  uint32_t len     = 29;
  attach_request_msg attach_request = {0};

  // Decode and encode back message
  int decoded = decode_attach_request(&attach_request, buffer, len);
  int encoded = encode_attach_request(&attach_request, temp_buffer, decoded);

  ASSERT_EQ(encoded, decoded);
  // Check encoded buffer and original buffer are equal
  EXPECT_ARRAY_EQ(buffer, temp_buffer, encoded);

  bdestroy_wrapper(&attach_request.esmmessagecontainer);
  bdestroy_wrapper(&attach_request.supportedcodecs);
}

TEST_F(EMMEncodeDecodeTest, TestDecodeAttachRequestPixel) {
  //   Combined attach, NAS message generated from Pixel 4
  uint8_t buffer[] = {
      0x72, 0x08, 0x39, 0x51, 0x10, 0x00, 0x30, 0x09, 0x01, 0x07, 0x07, 0xf0,
      0x70, 0xc0, 0x40, 0x19, 0x00, 0x80, 0x00, 0x34, 0x02, 0x0c, 0xd0, 0x11,
      0xd1, 0x27, 0x2d, 0x80, 0x80, 0x21, 0x10, 0x01, 0x00, 0x00, 0x10, 0x81,
      0x06, 0x00, 0x00, 0x00, 0x00, 0x83, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00,
      0x0d, 0x00, 0x00, 0x0a, 0x00, 0x00, 0x05, 0x00, 0x00, 0x10, 0x00, 0x00,
      0x11, 0x00, 0x00, 0x1a, 0x01, 0x01, 0x00, 0x23, 0x00, 0x00, 0x24, 0x00,
      0x5c, 0x0a, 0x01, 0x31, 0x04, 0x65, 0xe0, 0x3e, 0x00, 0x90, 0x11, 0x03,
      0x57, 0x58, 0xa6, 0x20, 0x0d, 0x60, 0x14, 0x04, 0xef, 0x65, 0x23, 0x3b,
      0x88, 0x00, 0x92, 0xf2, 0x00, 0x00, 0x40, 0x08, 0x04, 0x02, 0x60, 0x04,
      0x00, 0x02, 0x1f, 0x00, 0x5d, 0x01, 0x03, 0xc1};

  uint32_t len = 116;
  attach_request_msg attach_request;

  int rc = decode_attach_request(&attach_request, buffer, len);
  ASSERT_EQ(rc, len);
  ASSERT_EQ(attach_request.epsattachtype, EPS_ATTACH_TYPE_COMBINED_EPS_IMSI);
  ASSERT_EQ(
      attach_request.naskeysetidentifier.naskeysetidentifier,
      NAS_KEY_SET_IDENTIFIER_NOT_AVAILABLE);
  ASSERT_EQ(attach_request.naskeysetidentifier.tsc, 0);

  bdestroy_wrapper(&attach_request.esmmessagecontainer);
  bdestroy_wrapper(&attach_request.supportedcodecs);
}

TEST_F(EMMEncodeDecodeTest, TestEncodeDecodeEmmInformation) {
  emm_information_msg original_msg = {0};
  emm_information_msg decoded_msg  = {0};

  FILL_EMM_COMMON_MANDATORY_DEFAULTS(original_msg);

  // Set EMMInformation optional IEs
  original_msg.localtimezone             = 2;
  original_msg.networkdaylightsavingtime = 1;

  original_msg.presencemask =
      EMM_INFORMATION_LOCAL_TIME_ZONE_PRESENT |
      EMM_INFORMATION_NETWORK_DAYLIGHT_SAVING_TIME_PRESENT;

  // Encode and decode back message
  int encoded = encode_emm_information(&original_msg, temp_buffer, BUFFER_LEN);
  int decoded = decode_emm_information(&decoded_msg, temp_buffer, encoded);

  ASSERT_EQ(encoded, decoded);
  ASSERT_EQ(original_msg.localtimezone, decoded_msg.localtimezone);
  ASSERT_EQ(
      original_msg.networkdaylightsavingtime,
      decoded_msg.networkdaylightsavingtime);
}

TEST_F(EMMEncodeDecodeTest, TestEncodeDecodeUplinkNasTransport) {
  uplink_nas_transport_msg original_msg = {0};
  uplink_nas_transport_msg decoded_msg  = {0};

  FILL_EMM_COMMON_MANDATORY_DEFAULTS(original_msg);

  // Set mandatory IEs
  original_msg.nasmessagecontainer = bfromcstr("TEST_NAS_MSG_CONTAINER");

  // Encode and decode back message
  int encoded =
      encode_uplink_nas_transport(&original_msg, temp_buffer, BUFFER_LEN);
  int decoded = decode_uplink_nas_transport(&decoded_msg, temp_buffer, encoded);

  ASSERT_EQ(encoded, decoded);
  ASSERT_EQ(
      std::string(bdata(original_msg.nasmessagecontainer)),
      std::string(bdata(decoded_msg.nasmessagecontainer)));

  bdestroy_wrapper(&original_msg.nasmessagecontainer);
  bdestroy_wrapper(&decoded_msg.nasmessagecontainer);
}

TEST_F(EMMEncodeDecodeTest, TestEncodeDecodeTAURequest) {
  tracking_area_update_request_msg original_msg = {0};
  tracking_area_update_request_msg decoded_msg  = {0};

  FILL_EMM_COMMON_MANDATORY_DEFAULTS(original_msg);

  // Set mandatory IEs
  original_msg.epsupdatetype.active_flag = 1;
  original_msg.epsupdatetype.eps_update_type_value =
      EPS_UPDATE_TYPE_TA_UPDATING;
  original_msg.naskeysetidentifier.naskeysetidentifier =
      NAS_KEY_SET_IDENTIFIER_NATIVE;
  original_msg.naskeysetidentifier.tsc = 0;

  original_msg.oldguti.guti.typeofidentity = EPS_MOBILE_IDENTITY_GUTI;
  original_msg.oldguti.guti.mcc_digit1     = 0;
  original_msg.oldguti.guti.mcc_digit2     = 0;
  original_msg.oldguti.guti.mcc_digit3     = 1;
  original_msg.oldguti.guti.mnc_digit1     = 0;
  original_msg.oldguti.guti.mnc_digit2     = 1;
  original_msg.oldguti.guti.mnc_digit3     = 0x0f;
  original_msg.oldguti.guti.mme_group_id   = 1;
  original_msg.oldguti.guti.mme_code       = 1;
  original_msg.oldguti.guti.m_tmsi         = 0x2bfb815f;
  original_msg.oldguti.guti.spare          = 0xf;

  // Set optional IEs
  original_msg.noncurrentnativenaskeysetidentifier.naskeysetidentifier =
      NAS_KEY_SET_IDENTIFIER_NOT_AVAILABLE;
  original_msg.noncurrentnativenaskeysetidentifier.tsc = 0;
  original_msg.gprscipheringkeysequencenumber          = 1;
  original_msg.oldptmsisignature                       = 2;

  original_msg.additionalguti.guti.typeofidentity = EPS_MOBILE_IDENTITY_GUTI;
  original_msg.additionalguti.guti.mcc_digit1     = 0;
  original_msg.additionalguti.guti.mcc_digit2     = 1;
  original_msg.additionalguti.guti.mcc_digit3     = 1;
  original_msg.additionalguti.guti.mnc_digit1     = 0;
  original_msg.additionalguti.guti.mnc_digit2     = 1;
  original_msg.additionalguti.guti.mnc_digit3     = 0x0f;
  original_msg.additionalguti.guti.mme_group_id   = 1;
  original_msg.additionalguti.guti.mme_code       = 1;
  original_msg.additionalguti.guti.m_tmsi         = 0x2bfb816f;
  original_msg.additionalguti.guti.spare          = 0xf;

  original_msg.nonceue = 3;

  original_msg.drxparameter.nondrxtimer      = 0x7;
  original_msg.drxparameter.splitonccch      = 1;
  original_msg.drxparameter.splitpgcyclecode = 1;

  original_msg.ueradiocapabilityinformationupdateneeded = 1;
  original_msg.epsbearercontextstatus                   = 1;
  original_msg.tmsistatus                               = 1;
  original_msg.supportedcodecs = bfromcstr("TEST_CODECS");

  original_msg.presencemask =
      TRACKING_AREA_UPDATE_REQUEST_NONCURRENT_NATIVE_NAS_KEY_SET_IDENTIFIER_PRESENT |
      TRACKING_AREA_UPDATE_REQUEST_GPRS_CIPHERING_KEY_SEQUENCE_NUMBER_PRESENT |
      TRACKING_AREA_UPDATE_REQUEST_OLD_PTMSI_SIGNATURE_PRESENT |
      TRACKING_AREA_UPDATE_REQUEST_ADDITIONAL_GUTI_PRESENT |
      TRACKING_AREA_UPDATE_REQUEST_NONCEUE_PRESENT |
      TRACKING_AREA_UPDATE_REQUEST_DRX_PARAMETER_PRESENT |
      TRACKING_AREA_UPDATE_REQUEST_UE_RADIO_CAPABILITY_INFORMATION_UPDATE_NEEDED_PRESENT |
      TRACKING_AREA_UPDATE_REQUEST_EPS_BEARER_CONTEXT_STATUS_PRESENT |
      TRACKING_AREA_UPDATE_REQUEST_TMSI_STATUS_PRESENT |
      TRACKING_AREA_UPDATE_REQUEST_SUPPORTED_CODECS_PRESENT;

  // Encode and decode back message
  int encoded = encode_tracking_area_update_request(
      &original_msg, temp_buffer, BUFFER_LEN);
  int decoded =
      decode_tracking_area_update_request(&decoded_msg, temp_buffer, encoded);

  ASSERT_EQ(encoded, decoded);

  ASSERT_TRUE(!memcmp(
      &original_msg.oldguti.guti, &decoded_msg.oldguti.guti,
      sizeof(original_msg.oldguti.guti)));
  ASSERT_TRUE(!memcmp(
      &original_msg.additionalguti.guti, &decoded_msg.additionalguti.guti,
      sizeof(original_msg.additionalguti.guti)));
  ASSERT_EQ(
      original_msg.gprscipheringkeysequencenumber,
      decoded_msg.gprscipheringkeysequencenumber);
  ASSERT_EQ(original_msg.oldptmsisignature, decoded_msg.oldptmsisignature);
  ASSERT_EQ(
      original_msg.naskeysetidentifier.naskeysetidentifier,
      decoded_msg.naskeysetidentifier.naskeysetidentifier);
  ASSERT_EQ(
      original_msg.noncurrentnativenaskeysetidentifier.naskeysetidentifier,
      decoded_msg.noncurrentnativenaskeysetidentifier.naskeysetidentifier);
  ASSERT_EQ(
      original_msg.noncurrentnativenaskeysetidentifier.tsc,
      decoded_msg.noncurrentnativenaskeysetidentifier.tsc);
  ASSERT_EQ(
      original_msg.epsupdatetype.eps_update_type_value,
      decoded_msg.epsupdatetype.eps_update_type_value);
  ASSERT_EQ(
      original_msg.epsupdatetype.active_flag,
      decoded_msg.epsupdatetype.active_flag);
  ASSERT_EQ(
      original_msg.drxparameter.nondrxtimer,
      decoded_msg.drxparameter.nondrxtimer);
  ASSERT_EQ(
      original_msg.drxparameter.splitonccch,
      decoded_msg.drxparameter.splitonccch);
  ASSERT_EQ(
      original_msg.drxparameter.splitpgcyclecode,
      decoded_msg.drxparameter.splitpgcyclecode);
  ASSERT_EQ(original_msg.nonceue, decoded_msg.nonceue);
  ASSERT_EQ(original_msg.tmsistatus, decoded_msg.tmsistatus);
  ASSERT_EQ(
      original_msg.epsbearercontextstatus, decoded_msg.epsbearercontextstatus);
  ASSERT_EQ(
      original_msg.ueradiocapabilityinformationupdateneeded,
      decoded_msg.ueradiocapabilityinformationupdateneeded);
  ASSERT_EQ(
      std::string(bdata(original_msg.supportedcodecs)),
      std::string(bdata(decoded_msg.supportedcodecs)));

  bdestroy_wrapper(&original_msg.supportedcodecs);
  bdestroy_wrapper(&decoded_msg.supportedcodecs);
}

}  // namespace lte
}  // namespace magma
