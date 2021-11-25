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
#include <stdint.h>
#include <stdbool.h>
#include <cstring>
#include <string>
#include "bstrlib.h"

extern "C" {
#include "lte/gateway/c/core/oai/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/common/TLVEncoder.h"
#include "lte/gateway/c/core/oai/common/TLVDecoder.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/msg/ActivateDedicatedEpsBearerContextAccept.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/msg/ActivateDedicatedEpsBearerContextReject.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/msg/ActivateDefaultEpsBearerContextRequest.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/msg/ActivateDedicatedEpsBearerContextRequest.h"
#include "lte/gateway/c/core/oai/common/common_defs.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.008.h"
}

namespace magma {
namespace lte {

#define BUFFER_LEN 200

#define FILL_COMMON_MANDATORY_DEFAULTS(msg)                                    \
  do {                                                                         \
    msg.protocoldiscriminator        = EPS_SESSION_MANAGEMENT_MESSAGE;         \
    msg.epsbeareridentity            = 8;                                      \
    msg.proceduretransactionidentity = 2;                                      \
    msg.messagetype                  = 44;                                     \
  } while (0)

#define COMPARE_COMMON_MANDATORY_DEFAULTS()                                    \
  do {                                                                         \
    EXPECT_EQ(                                                                 \
        original_msg.protocoldiscriminator,                                    \
        decoded_msg.protocoldiscriminator);                                    \
    EXPECT_EQ(original_msg.epsbeareridentity, decoded_msg.epsbeareridentity);  \
    EXPECT_EQ(                                                                 \
        original_msg.proceduretransactionidentity,                             \
        decoded_msg.proceduretransactionidentity);                             \
    EXPECT_EQ(original_msg.messagetype, decoded_msg.messagetype);              \
  } while (0)

#define DESTROY_PCO()                                                          \
  do {                                                                         \
    bdestroy_wrapper(&original_msg.protocolconfigurationoptions                \
                          .protocol_or_container_ids[0]                        \
                          .contents);                                          \
    bdestroy_wrapper(&original_msg.protocolconfigurationoptions                \
                          .protocol_or_container_ids[1]                        \
                          .contents);                                          \
    bdestroy_wrapper(                                                          \
        &decoded_msg.protocolconfigurationoptions.protocol_or_container_ids[0] \
             .contents);                                                       \
    bdestroy_wrapper(                                                          \
        &decoded_msg.protocolconfigurationoptions.protocol_or_container_ids[1] \
             .contents);                                                       \
  } while (0)

class ESMEncodeDecodeTest : public ::testing::Test {
  virtual void SetUp() {}
  virtual void TearDown() {}

 protected:
  void fill_pco(protocol_configuration_options_t* pco) {
    pco->num_protocol_or_container_id    = 2;
    pco->protocol_or_container_ids[0].id = PCO_CI_P_CSCF_IPV6_ADDRESS_REQUEST;
    bstring test_string1                 = bfromcstr("teststring");
    pco->protocol_or_container_ids[0].contents = test_string1;
    pco->protocol_or_container_ids[0].length   = blength(test_string1);
    pco->protocol_or_container_ids[1].id =
        PCO_CI_DSMIPV6_IPV4_HOME_AGENT_ADDRESS;
    bstring test_string2 = bfromcstr("longer.test.string");
    pco->protocol_or_container_ids[1].contents = test_string2;
    pco->protocol_or_container_ids[1].length   = blength(test_string2);
    return;
  }

  uint8_t buffer[BUFFER_LEN];
};

TEST_F(ESMEncodeDecodeTest, TestActivateDedicatedEpsBearerContextAccept) {
  activate_dedicated_eps_bearer_context_accept_msg original_msg = {0};
  activate_dedicated_eps_bearer_context_accept_msg decoded_msg  = {0};
  FILL_COMMON_MANDATORY_DEFAULTS(original_msg);
  original_msg.presencemask =
      ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_ACCEPT_PROTOCOL_CONFIGURATION_OPTIONS_PRESENT;
  fill_pco(&original_msg.protocolconfigurationoptions);

  int encoded = encode_activate_dedicated_eps_bearer_context_accept(
      &original_msg, buffer, BUFFER_LEN);
  int decoded = decode_activate_dedicated_eps_bearer_context_accept(
      &decoded_msg, buffer, encoded);

  EXPECT_EQ(encoded, decoded);
  // Header needs to be decoded separately and then can be compared
  // COMPARE_COMMON_MANDATORY_DEFAULTS();
  EXPECT_EQ(original_msg.presencemask, decoded_msg.presencemask);
  EXPECT_EQ(
      std::string((const char*) original_msg.protocolconfigurationoptions
                      .protocol_or_container_ids[0]
                      .contents->data),
      std::string((const char*) decoded_msg.protocolconfigurationoptions
                      .protocol_or_container_ids[0]
                      .contents->data));
  EXPECT_EQ(
      std::string((const char*) original_msg.protocolconfigurationoptions
                      .protocol_or_container_ids[1]
                      .contents->data),
      std::string((const char*) decoded_msg.protocolconfigurationoptions
                      .protocol_or_container_ids[1]
                      .contents->data));

  DESTROY_PCO();
}

TEST_F(ESMEncodeDecodeTest, TestActivateDedicatedEpsBearerContextReject) {
  activate_dedicated_eps_bearer_context_reject_msg original_msg = {0};
  activate_dedicated_eps_bearer_context_reject_msg decoded_msg  = {0};
  FILL_COMMON_MANDATORY_DEFAULTS(original_msg);
  original_msg.esmcause = 255;
  original_msg.presencemask =
      ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REJECT_PROTOCOL_CONFIGURATION_OPTIONS_PRESENT;
  fill_pco(&original_msg.protocolconfigurationoptions);

  int encoded = encode_activate_dedicated_eps_bearer_context_reject(
      &original_msg, buffer, BUFFER_LEN);
  int decoded = decode_activate_dedicated_eps_bearer_context_reject(
      &decoded_msg, buffer, encoded);

  EXPECT_EQ(encoded, decoded);
  // Header needs to be decoded separately and then can be compared
  // COMPARE_COMMON_MANDATORY_DEFAULTS();
  EXPECT_EQ(original_msg.esmcause, decoded_msg.esmcause);
  EXPECT_EQ(original_msg.presencemask, decoded_msg.presencemask);
  EXPECT_EQ(
      std::string((const char*) original_msg.protocolconfigurationoptions
                      .protocol_or_container_ids[0]
                      .contents->data),
      std::string((const char*) decoded_msg.protocolconfigurationoptions
                      .protocol_or_container_ids[0]
                      .contents->data));
  EXPECT_EQ(
      std::string((const char*) original_msg.protocolconfigurationoptions
                      .protocol_or_container_ids[1]
                      .contents->data),
      std::string((const char*) decoded_msg.protocolconfigurationoptions
                      .protocol_or_container_ids[1]
                      .contents->data));

  DESTROY_PCO();
}

TEST_F(ESMEncodeDecodeTest, TestActivateDefaultEpsBearerContextRequest) {
  activate_default_eps_bearer_context_request_msg original_msg = {0};
  activate_default_eps_bearer_context_request_msg decoded_msg  = {0};
  FILL_COMMON_MANDATORY_DEFAULTS(original_msg);

  original_msg.epsqos.bitRatesPresent               = 1;
  original_msg.epsqos.bitRatesExtPresent            = 1;
  original_msg.epsqos.bitRatesExt2Present           = 1;
  original_msg.epsqos.bitRates.guarBitRateForDL     = 10;
  original_msg.epsqos.bitRates.guarBitRateForUL     = 5;
  original_msg.epsqos.bitRates.maxBitRateForDL      = 100;
  original_msg.epsqos.bitRates.maxBitRateForUL      = 50;
  original_msg.epsqos.bitRatesExt.guarBitRateForDL  = 10;
  original_msg.epsqos.bitRatesExt.guarBitRateForUL  = 5;
  original_msg.epsqos.bitRatesExt.maxBitRateForDL   = 100;
  original_msg.epsqos.bitRatesExt.maxBitRateForUL   = 50;
  original_msg.epsqos.bitRatesExt2.guarBitRateForDL = 10;
  original_msg.epsqos.bitRatesExt2.guarBitRateForUL = 5;
  original_msg.epsqos.bitRatesExt2.maxBitRateForDL  = 100;
  original_msg.epsqos.bitRatesExt2.maxBitRateForUL  = 50;
  original_msg.epsqos.qci                           = 5;

  original_msg.accesspointname = bfromcstr("magma.ipv4");

  original_msg.pdnaddress.pdntypevalue = PDN_VALUE_TYPE_IPV4V6;
  original_msg.pdnaddress.pdnaddressinformation =
      bfromcstr("192.154.134.111,ef::fff");

  // TI_IE encoding/decoding is not implemented, hence its flag as well as
  // entries will be omitted.
  original_msg.presencemask =
      ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_REQUEST_NEGOTIATED_QOS_PRESENT |
      ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_REQUEST_NEGOTIATED_LLC_SAPI_PRESENT |
      ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_REQUEST_RADIO_PRIORITY_PRESENT |
      ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_REQUEST_PACKET_FLOW_IDENTIFIER_PRESENT |
      ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_REQUEST_APNAMBR_PRESENT |
      ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_REQUEST_ESM_CAUSE_PRESENT |
      ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_REQUEST_PROTOCOL_CONFIGURATION_OPTIONS_PRESENT;

  original_msg.negotiatedqos.delayclass                 = 1;
  original_msg.negotiatedqos.reliabilityclass           = 2;
  original_msg.negotiatedqos.peakthroughput             = 15;
  original_msg.negotiatedqos.precedenceclass            = 2;
  original_msg.negotiatedqos.meanthroughput             = 16;
  original_msg.negotiatedqos.trafficclass               = 7;
  original_msg.negotiatedqos.deliveryorder              = 1;
  original_msg.negotiatedqos.deliveryoferroneoussdu     = 2;
  original_msg.negotiatedqos.maximumsdusize             = 255;
  original_msg.negotiatedqos.maximumbitrateuplink       = 99;
  original_msg.negotiatedqos.maximumbitratedownlink     = 169;
  original_msg.negotiatedqos.residualber                = 4;
  original_msg.negotiatedqos.sduratioerror              = 12;
  original_msg.negotiatedqos.transferdelay              = 61;
  original_msg.negotiatedqos.traffichandlingpriority    = 0;
  original_msg.negotiatedqos.guaranteedbitrateuplink    = 10;
  original_msg.negotiatedqos.guaranteedbitratedownlink  = 20;
  original_msg.negotiatedqos.signalingindication        = 1;
  original_msg.negotiatedqos.sourcestatisticsdescriptor = 13;

  original_msg.negotiatedllcsapi    = 10;
  original_msg.radiopriority        = 5;
  original_msg.packetflowidentifier = 118;

  original_msg.apnambr.apnambrfordownlink           = 11;
  original_msg.apnambr.apnambrforuplink             = 11;
  original_msg.apnambr.apnambrfordownlink_extended  = 11;
  original_msg.apnambr.apnambrforuplink_extended    = 11;
  original_msg.apnambr.apnambrfordownlink_extended2 = 11;
  original_msg.apnambr.apnambrforuplink_extended2   = 11;
  original_msg.apnambr.extensions                   = 4;

  original_msg.esmcause = 102;

  fill_pco(&original_msg.protocolconfigurationoptions);

  int encoded = encode_activate_default_eps_bearer_context_request(
      &original_msg, buffer, BUFFER_LEN);
  int decoded = decode_activate_default_eps_bearer_context_request(
      &decoded_msg, buffer, encoded);

  EXPECT_EQ(encoded, decoded);
  // Header needs to be decoded separately and then can be compared
  // COMPARE_COMMON_MANDATORY_DEFAULTS();
  EXPECT_EQ(
      original_msg.pdnaddress.pdntypevalue,
      decoded_msg.pdnaddress.pdntypevalue);
  EXPECT_FALSE(memcmp(
      &original_msg.epsqos, &decoded_msg.epsqos, sizeof(original_msg.epsqos)));
  EXPECT_EQ(
      std::string((const char*) original_msg.accesspointname->data),
      std::string((const char*) decoded_msg.accesspointname->data));
  EXPECT_EQ(
      std::string(
          (const char*) original_msg.pdnaddress.pdnaddressinformation->data),
      std::string(
          (const char*) decoded_msg.pdnaddress.pdnaddressinformation->data));
  EXPECT_EQ(original_msg.presencemask, decoded_msg.presencemask);
  EXPECT_FALSE(memcmp(
      &original_msg.negotiatedqos, &decoded_msg.negotiatedqos,
      sizeof(original_msg.negotiatedqos)));
  // EXPECT_EQ(original_msg.negotiatedllcsapi, decoded_msg.negotiatedllcsapi);
  EXPECT_EQ(original_msg.radiopriority, decoded_msg.radiopriority);
  // EXPECT_EQ(
  //    original_msg.packetflowidentifier, decoded_msg.packetflowidentifier);
  //   EXPECT_FALSE(memcmp(
  //       &original_msg.apnambr, &decoded_msg.apnambr,
  //       sizeof(original_msg.apnambr)));
  EXPECT_EQ(original_msg.esmcause, decoded_msg.esmcause);
  EXPECT_EQ(
      std::string((const char*) original_msg.protocolconfigurationoptions
                      .protocol_or_container_ids[0]
                      .contents->data),
      std::string((const char*) decoded_msg.protocolconfigurationoptions
                      .protocol_or_container_ids[0]
                      .contents->data));
  EXPECT_EQ(
      std::string((const char*) original_msg.protocolconfigurationoptions
                      .protocol_or_container_ids[1]
                      .contents->data),
      std::string((const char*) decoded_msg.protocolconfigurationoptions
                      .protocol_or_container_ids[1]
                      .contents->data));

  DESTROY_PCO();
  bdestroy_wrapper(&original_msg.accesspointname);
  bdestroy_wrapper(&decoded_msg.accesspointname);
  bdestroy_wrapper(&original_msg.pdnaddress.pdnaddressinformation);
  bdestroy_wrapper(&decoded_msg.pdnaddress.pdnaddressinformation);
}

}  // namespace lte
}  // namespace magma
