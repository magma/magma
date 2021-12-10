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
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/AttachAccept.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/CsServiceNotification.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/EmmInformation.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/ExtendedServiceRequest.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/TrackingAreaUpdateRequest.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/TrackingAreaUpdateReject.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/TrackingAreaUpdateAccept.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/NASSecurityModeCommand.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/NASSecurityModeComplete.h"
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

#define FILL_EMM_GUTI(msg_guti)                                                \
  do {                                                                         \
    msg_guti.guti.typeofidentity = EPS_MOBILE_IDENTITY_GUTI;                   \
    msg_guti.guti.mcc_digit1     = 0;                                          \
    msg_guti.guti.mcc_digit2     = 0;                                          \
    msg_guti.guti.mcc_digit3     = 1;                                          \
    msg_guti.guti.mnc_digit1     = 0;                                          \
    msg_guti.guti.mnc_digit2     = 1;                                          \
    msg_guti.guti.mnc_digit3     = 0x0f;                                       \
    msg_guti.guti.mme_group_id   = 1;                                          \
    msg_guti.guti.mme_code       = 1;                                          \
    msg_guti.guti.m_tmsi         = 0x2bfb815f;                                 \
    msg_guti.guti.spare          = 0xf;                                        \
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

  FILL_EMM_GUTI(original_msg.oldguti);

  // Set optional IEs
  original_msg.noncurrentnativenaskeysetidentifier.naskeysetidentifier =
      NAS_KEY_SET_IDENTIFIER_NOT_AVAILABLE;
  original_msg.noncurrentnativenaskeysetidentifier.tsc = 0;
  original_msg.gprscipheringkeysequencenumber          = 1;
  original_msg.oldptmsisignature                       = 2;

  FILL_EMM_GUTI(original_msg.additionalguti);

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

TEST_F(EMMEncodeDecodeTest, TestEncodeDecodeTAUAccept) {
  tracking_area_update_accept_msg original_msg = {0};
  tracking_area_update_accept_msg decoded_msg  = {0};

  FILL_EMM_COMMON_MANDATORY_DEFAULTS(original_msg);

  // Set mandatory IEs
  original_msg.epsupdateresult = EPS_UPDATE_RESULT_COMBINED_TA_LA_UPDATED;

  // Set optional IEs
  original_msg.t3412value.unit       = GPRS_TIMER_UNIT_2S;
  original_msg.t3412value.timervalue = 1;

  FILL_EMM_GUTI(original_msg.guti);

  original_msg.epsbearercontextstatus = 1;

  original_msg.locationareaidentification.mccdigit1 = 0;
  original_msg.locationareaidentification.mccdigit2 = 0;
  original_msg.locationareaidentification.mccdigit3 = 1;
  original_msg.locationareaidentification.mncdigit1 = 0;
  original_msg.locationareaidentification.mncdigit2 = 1;
  original_msg.locationareaidentification.mncdigit3 = 1;
  original_msg.locationareaidentification.lac       = 0x01;

  original_msg.emmcause = 1;

  original_msg.epsnetworkfeaturesupport.b1 = 1;
  original_msg.epsnetworkfeaturesupport.b2 = 0;

  original_msg.additionalupdateresult = 0;

  original_msg.presencemask =
      TRACKING_AREA_UPDATE_ACCEPT_T3412_VALUE_PRESENT |
      TRACKING_AREA_UPDATE_ACCEPT_GUTI_PRESENT |
      TRACKING_AREA_UPDATE_ACCEPT_EPS_BEARER_CONTEXT_STATUS_PRESENT |
      TRACKING_AREA_UPDATE_ACCEPT_LOCATION_AREA_IDENTIFICATION_PRESENT |
      TRACKING_AREA_UPDATE_ACCEPT_EMM_CAUSE_PRESENT |
      TRACKING_AREA_UPDATE_ACCEPT_EPS_NETWORK_FEATURE_SUPPORT_PRESENT |
      TRACKING_AREA_UPDATE_ACCEPT_ADDITIONAL_UPDATE_RESULT_PRESENT;

  // Encode and decode back message
  int encoded = encode_tracking_area_update_accept(
      &original_msg, temp_buffer, BUFFER_LEN);
  int decoded =
      decode_tracking_area_update_accept(&decoded_msg, temp_buffer, encoded);

  ASSERT_EQ(decoded, encoded);
  ASSERT_EQ(original_msg.epsupdateresult, decoded_msg.epsupdateresult);
  ASSERT_EQ(original_msg.t3412value.unit, decoded_msg.t3412value.unit);
  ASSERT_EQ(
      original_msg.t3412value.timervalue, decoded_msg.t3412value.timervalue);

  ASSERT_TRUE(!memcmp(
      &original_msg.guti, &decoded_msg.guti, sizeof(original_msg.guti)));
  ASSERT_EQ(
      original_msg.epsbearercontextstatus, decoded_msg.epsbearercontextstatus);
  ASSERT_TRUE(!memcmp(
      &original_msg.locationareaidentification,
      &decoded_msg.locationareaidentification,
      sizeof(original_msg.locationareaidentification)));
  ASSERT_EQ(original_msg.emmcause, decoded_msg.emmcause);
  ASSERT_EQ(
      original_msg.epsnetworkfeaturesupport.b1,
      decoded_msg.epsnetworkfeaturesupport.b1);
  ASSERT_EQ(
      original_msg.epsnetworkfeaturesupport.b2,
      decoded_msg.epsnetworkfeaturesupport.b2);
  ASSERT_EQ(
      original_msg.additionalupdateresult, decoded_msg.additionalupdateresult);
}

TEST_F(EMMEncodeDecodeTest, TestEncodeDecodeTAUReject) {
  tracking_area_update_reject_msg original_msg = {0};
  tracking_area_update_reject_msg decoded_msg  = {0};

  FILL_EMM_COMMON_MANDATORY_DEFAULTS(original_msg);

  // Set mandatory IEs
  original_msg.emmcause = 1;

  // Encode and decode back message
  int encoded = encode_tracking_area_update_reject(
      &original_msg, temp_buffer, BUFFER_LEN);
  int decoded =
      decode_tracking_area_update_reject(&decoded_msg, temp_buffer, encoded);

  ASSERT_EQ(encoded, decoded);
  ASSERT_EQ(original_msg.emmcause, decoded_msg.emmcause);
}

TEST_F(EMMEncodeDecodeTest, TestEncodeDecodeCSServiceNotification) {
  cs_service_notification_msg original_msg = {0};
  cs_service_notification_msg decoded_msg  = {0};

  FILL_EMM_COMMON_MANDATORY_DEFAULTS(original_msg);

  // Set mandatory IEs
  original_msg.pagingidentity = 1;

  // Set optional IEs
  original_msg.cli               = bfromcstr("TEST_CLI");
  original_msg.sscode            = 2;
  original_msg.lcsindicator      = 1;
  original_msg.lcsclientidentity = bfromcstr("TEST_LCS");

  original_msg.presencemask =
      CS_SERVICE_NOTIFICATION_CLI_PRESENT |
      CS_SERVICE_NOTIFICATION_SS_CODE_PRESENT |
      CS_SERVICE_NOTIFICATION_LCS_INDICATOR_PRESENT |
      CS_SERVICE_NOTIFICATION_LCS_CLIENT_IDENTITY_PRESENT;

  // Encode and decode back message
  int encoded =
      encode_cs_service_notification(&original_msg, temp_buffer, BUFFER_LEN);
  int decoded =
      decode_cs_service_notification(&decoded_msg, temp_buffer, encoded);

  ASSERT_EQ(encoded, decoded);
  ASSERT_EQ(original_msg.pagingidentity, decoded_msg.pagingidentity);
  ASSERT_EQ(original_msg.lcsindicator, decoded_msg.lcsindicator);
  ASSERT_EQ(original_msg.sscode, decoded_msg.sscode);
  ASSERT_EQ(
      std::string(bdata(original_msg.cli)),
      std::string(bdata(decoded_msg.cli)));
  ASSERT_EQ(
      std::string(bdata(original_msg.lcsclientidentity)),
      std::string(bdata(decoded_msg.lcsclientidentity)));

  bdestroy_wrapper(&original_msg.cli);
  bdestroy_wrapper(&decoded_msg.cli);

  bdestroy_wrapper(&original_msg.lcsclientidentity);
  bdestroy_wrapper(&decoded_msg.lcsclientidentity);
}

TEST_F(EMMEncodeDecodeTest, TestEncodeDecodeAttachAcceptConsecutiveTACs) {
  attach_accept_msg original_msg = {0};
  attach_accept_msg decoded_msg  = {0};

  FILL_EMM_COMMON_MANDATORY_DEFAULTS(original_msg);

  // Set mandatory IEs
  original_msg.epsattachresult = EPS_ATTACH_RESULT_EPS;

  original_msg.t3412value.unit       = GPRS_TIMER_UNIT_2S;
  original_msg.t3412value.timervalue = 1;

  original_msg.esmmessagecontainer = bfromcstr("TEST_CONTAINER");

  original_msg.tailist.numberoflists = 1;
  original_msg.tailist.partial_tai_list[0].typeoflist =
      TRACKING_AREA_IDENTITY_LIST_ONE_PLMN_CONSECUTIVE_TACS;
  original_msg.tailist.partial_tai_list[0].numberofelements = 1;
  original_msg.tailist.partial_tai_list[0]
      .u.tai_one_plmn_consecutive_tacs.plmn.mcc_digit1 = 0;
  original_msg.tailist.partial_tai_list[0]
      .u.tai_one_plmn_consecutive_tacs.plmn.mcc_digit2 = 0;
  original_msg.tailist.partial_tai_list[0]
      .u.tai_one_plmn_consecutive_tacs.plmn.mcc_digit3 = 1;
  original_msg.tailist.partial_tai_list[0]
      .u.tai_one_plmn_consecutive_tacs.plmn.mnc_digit1 = 0;
  original_msg.tailist.partial_tai_list[0]
      .u.tai_one_plmn_consecutive_tacs.plmn.mnc_digit2 = 1;
  original_msg.tailist.partial_tai_list[0]
      .u.tai_one_plmn_consecutive_tacs.plmn.mnc_digit3 = 1;
  original_msg.tailist.partial_tai_list[0].u.tai_one_plmn_consecutive_tacs.tac =
      1;

  // Set optional IEs
  FILL_EMM_GUTI(original_msg.guti);

  original_msg.t3402value.unit       = GPRS_TIMER_UNIT_60S;
  original_msg.t3402value.timervalue = 1;

  original_msg.locationareaidentification.mccdigit1 = 0;
  original_msg.locationareaidentification.mccdigit2 = 0;
  original_msg.locationareaidentification.mccdigit3 = 1;
  original_msg.locationareaidentification.mncdigit1 = 0;
  original_msg.locationareaidentification.mncdigit2 = 1;
  original_msg.locationareaidentification.mncdigit3 = 1;
  original_msg.locationareaidentification.lac       = 0x01;

  original_msg.emmcause = 1;

  original_msg.epsnetworkfeaturesupport.b1 = 1;
  original_msg.epsnetworkfeaturesupport.b2 = 0;

  original_msg.additionalupdateresult = 0;

  original_msg.presencemask =
      ATTACH_ACCEPT_GUTI_PRESENT |
      ATTACH_ACCEPT_LOCATION_AREA_IDENTIFICATION_PRESENT |
      ATTACH_ACCEPT_EMM_CAUSE_PRESENT | ATTACH_ACCEPT_T3402_VALUE_PRESENT |
      ATTACH_ACCEPT_EPS_NETWORK_FEATURE_SUPPORT_PRESENT |
      ATTACH_ACCEPT_ADDITIONAL_UPDATE_RESULT_PRESENT;

  // Encode and decode back message
  int encoded = encode_attach_accept(&original_msg, temp_buffer, BUFFER_LEN);
  int decoded = decode_attach_accept(&decoded_msg, temp_buffer, encoded);

  ASSERT_EQ(encoded, decoded);

  ASSERT_EQ(original_msg.epsattachresult, decoded_msg.epsattachresult);
  ASSERT_EQ(original_msg.t3412value.unit, decoded_msg.t3412value.unit);
  ASSERT_EQ(
      original_msg.t3412value.timervalue, decoded_msg.t3412value.timervalue);
  ASSERT_EQ(
      original_msg.t3402value.timervalue, decoded_msg.t3402value.timervalue);
  ASSERT_EQ(original_msg.t3402value.unit, decoded_msg.t3402value.unit);
  ASSERT_EQ(
      std::string(bdata(original_msg.esmmessagecontainer)),
      std::string(bdata(decoded_msg.esmmessagecontainer)));

  ASSERT_TRUE(!memcmp(
      &original_msg.guti, &decoded_msg.guti, sizeof(original_msg.guti)));
  ASSERT_TRUE(!memcmp(
      &original_msg.tailist, &decoded_msg.tailist,
      sizeof(original_msg.tailist)));
  ASSERT_TRUE(!memcmp(
      &original_msg.locationareaidentification,
      &decoded_msg.locationareaidentification,
      sizeof(original_msg.locationareaidentification)));
  ASSERT_EQ(original_msg.emmcause, decoded_msg.emmcause);
  ASSERT_EQ(
      original_msg.epsnetworkfeaturesupport.b1,
      decoded_msg.epsnetworkfeaturesupport.b1);
  ASSERT_EQ(
      original_msg.epsnetworkfeaturesupport.b2,
      decoded_msg.epsnetworkfeaturesupport.b2);
  ASSERT_EQ(
      original_msg.additionalupdateresult, decoded_msg.additionalupdateresult);

  bdestroy_wrapper(&original_msg.esmmessagecontainer);
  bdestroy_wrapper(&decoded_msg.esmmessagecontainer);
}

TEST_F(EMMEncodeDecodeTest, TestEncodeDecodeSMCCommand) {
  security_mode_command_msg original_msg = {0};
  security_mode_command_msg decoded_msg  = {0};

  FILL_EMM_COMMON_MANDATORY_DEFAULTS(original_msg);

  // Set mandatory IEs
  original_msg.selectednassecurityalgorithms.typeofcipheringalgorithm =
      NAS_SECURITY_ALGORITHMS_EEA2;
  original_msg.selectednassecurityalgorithms.typeofintegrityalgorithm =
      NAS_SECURITY_ALGORITHMS_EIA2;
  original_msg.naskeysetidentifier.naskeysetidentifier =
      NAS_KEY_SET_IDENTIFIER_NATIVE;
  original_msg.naskeysetidentifier.tsc = 0;

  original_msg.replayeduesecuritycapabilities.eea = UE_SECURITY_CAPABILITY_EEA1;
  original_msg.replayeduesecuritycapabilities.eia = UE_SECURITY_CAPABILITY_EIA1;
  original_msg.replayeduesecuritycapabilities.umts_present = 1;
  original_msg.replayeduesecuritycapabilities.uea = UE_SECURITY_CAPABILITY_UEA1;
  original_msg.replayeduesecuritycapabilities.uia = UE_SECURITY_CAPABILITY_UIA1;
  original_msg.replayeduesecuritycapabilities.gprs_present = 1;
  original_msg.replayeduesecuritycapabilities.gea = UE_SECURITY_CAPABILITY_GEA1;

  // Set optional IEs
  original_msg.imeisvrequest = IMEISV_REQUESTED;

  original_msg.replayedueadditionalsecuritycapabilities._5g_ea = 1;
  original_msg.replayedueadditionalsecuritycapabilities._5g_ia = 1;

  original_msg.presencemask =
      SECURITY_MODE_COMMAND_IMEISV_REQUEST_PRESENT |
      SECURITY_MODE_COMMAND_REPLAYED_UE_ADDITIONAL_SECU_CAPABILITY_PRESENT;

  // Encode and decode back message
  int encoded =
      encode_security_mode_command(&original_msg, temp_buffer, BUFFER_LEN);
  int decoded =
      decode_security_mode_command(&decoded_msg, temp_buffer, encoded);

  ASSERT_EQ(encoded, decoded);

  ASSERT_TRUE(!memcmp(
      &original_msg.replayeduesecuritycapabilities,
      &decoded_msg.replayeduesecuritycapabilities,
      sizeof(original_msg.replayeduesecuritycapabilities)));
  ASSERT_EQ(
      original_msg.replayedueadditionalsecuritycapabilities._5g_ea,
      decoded_msg.replayedueadditionalsecuritycapabilities._5g_ea);
  ASSERT_EQ(
      original_msg.replayedueadditionalsecuritycapabilities._5g_ia,
      decoded_msg.replayedueadditionalsecuritycapabilities._5g_ia);
  ASSERT_EQ(
      original_msg.naskeysetidentifier.naskeysetidentifier,
      decoded_msg.naskeysetidentifier.naskeysetidentifier);
  ASSERT_EQ(
      original_msg.naskeysetidentifier.tsc,
      decoded_msg.naskeysetidentifier.tsc);
}

TEST_F(EMMEncodeDecodeTest, TestEncodeDecodeSMCComplete) {
  security_mode_complete_msg original_msg = {0};
  security_mode_complete_msg decoded_msg  = {0};

  FILL_EMM_COMMON_MANDATORY_DEFAULTS(original_msg);

  // Set optional IEs
  original_msg.imeisv.imeisv.typeofidentity = MOBILE_IDENTITY_IMEISV;
  original_msg.imeisv.imeisv.oddeven        = 1;
  original_msg.imeisv.imeisv.tac1           = 1;
  original_msg.imeisv.imeisv.tac2           = 1;
  original_msg.imeisv.imeisv.tac3           = 0;
  original_msg.imeisv.imeisv.tac4           = 0;
  original_msg.imeisv.imeisv.tac5           = 0;
  original_msg.imeisv.imeisv.tac6           = 0;
  original_msg.imeisv.imeisv.tac7           = 0;
  original_msg.imeisv.imeisv.tac8           = 0;
  original_msg.imeisv.imeisv.snr1           = 1;
  original_msg.imeisv.imeisv.snr2           = 1;
  original_msg.imeisv.imeisv.snr3           = 0;
  original_msg.imeisv.imeisv.snr4           = 0;
  original_msg.imeisv.imeisv.snr5           = 0;
  original_msg.imeisv.imeisv.snr6           = 0;
  original_msg.imeisv.imeisv.svn1           = 0;
  original_msg.imeisv.imeisv.svn2           = 0;
  original_msg.imeisv.imeisv.last           = 0;

  original_msg.presencemask = SECURITY_MODE_COMPLETE_IMEISV_PRESENT;

  // Encode and decode back message
  int encoded =
      encode_security_mode_complete(&original_msg, temp_buffer, BUFFER_LEN);
  int decoded =
      decode_security_mode_complete(&decoded_msg, temp_buffer, encoded);

  ASSERT_EQ(encoded, decoded);

  ASSERT_TRUE(!memcmp(
      &original_msg.imeisv, &decoded_msg.imeisv, sizeof(original_msg.imeisv)));
}

TEST_F(EMMEncodeDecodeTest, TestEncodeDecodeExtendedServiceRequest) {
  extended_service_request_msg original_msg = {0};
  extended_service_request_msg decoded_msg  = {0};

  FILL_EMM_COMMON_MANDATORY_DEFAULTS(original_msg);

  // Set mandatory IEs
  original_msg.servicetype                             = MT_CS_FB;
  original_msg.naskeysetidentifier.naskeysetidentifier = 1;
  original_msg.naskeysetidentifier.tsc                 = 1;

  original_msg.mtmsi.tmsi.typeofidentity = MOBILE_IDENTITY_TMSI;
  original_msg.mtmsi.tmsi.tmsi[0]        = 1;
  original_msg.mtmsi.tmsi.tmsi[1]        = 0;
  original_msg.mtmsi.tmsi.tmsi[2]        = 1;
  original_msg.mtmsi.tmsi.tmsi[3]        = 1;
  original_msg.mtmsi.tmsi.oddeven        = 0;
  original_msg.mtmsi.tmsi.f              = 0xf;

  original_msg.csfbresponse = 1;

  original_msg.presencemask = EMM_CSFB_RSP_PRESENT;

  // Encode and decode back message
  int encoded =
      encode_extended_service_request(&original_msg, temp_buffer, BUFFER_LEN);
  int decoded =
      decode_extended_service_request(&decoded_msg, temp_buffer, encoded);

  ASSERT_EQ(encoded, decoded);

  ASSERT_EQ(original_msg.servicetype, decoded_msg.servicetype);
  ASSERT_EQ(original_msg.csfbresponse, decoded_msg.csfbresponse);
  ASSERT_EQ(
      original_msg.naskeysetidentifier.naskeysetidentifier,
      decoded_msg.naskeysetidentifier.naskeysetidentifier);
  ASSERT_EQ(
      original_msg.naskeysetidentifier.tsc,
      decoded_msg.naskeysetidentifier.tsc);
  ASSERT_TRUE(!memcmp(
      &original_msg.mtmsi, &decoded_msg.mtmsi, sizeof(original_msg.mtmsi)));
}

}  // namespace lte
}  // namespace magma
