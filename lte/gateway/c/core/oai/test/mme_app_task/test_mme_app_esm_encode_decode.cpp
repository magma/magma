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

extern "C" {
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
#include "lte/gateway/c/core/oai/include/EpsQualityOfService.h"
#include "lte/gateway/c/core/oai/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/common/TLVEncoder.h"
#include "lte/gateway/c/core/oai/common/TLVDecoder.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/msg/ActivateDedicatedEpsBearerContextAccept.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/msg/ActivateDefaultEpsBearerContextAccept.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/msg/ActivateDedicatedEpsBearerContextReject.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/msg/ActivateDefaultEpsBearerContextReject.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/msg/ActivateDefaultEpsBearerContextRequest.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/msg/ActivateDedicatedEpsBearerContextRequest.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/msg/BearerResourceAllocationReject.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/msg/BearerResourceAllocationRequest.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/msg/BearerResourceModificationReject.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/msg/BearerResourceModificationRequest.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/msg/DeactivateEpsBearerContextAccept.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/msg/DeactivateEpsBearerContextRequest.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/msg/EsmInformationRequest.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/msg/EsmInformationResponse.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/msg/EsmStatus.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/msg/ModifyEpsBearerContextRequest.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/msg/ModifyEpsBearerContextAccept.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/msg/ModifyEpsBearerContextReject.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/msg/PdnConnectivityRequest.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/msg/PdnConnectivityReject.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/msg/PdnDisconnectRequest.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/msg/PdnDisconnectReject.h"
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

#define EXPECT_EQ_PCO()                                                        \
  do {                                                                         \
    EXPECT_EQ(                                                                 \
        std::string((const char*) original_msg.protocolconfigurationoptions    \
                        .protocol_or_container_ids[0]                          \
                        .contents->data),                                      \
        std::string((const char*) decoded_msg.protocolconfigurationoptions     \
                        .protocol_or_container_ids[0]                          \
                        .contents->data));                                     \
    EXPECT_EQ(                                                                 \
        std::string((const char*) original_msg.protocolconfigurationoptions    \
                        .protocol_or_container_ids[1]                          \
                        .contents->data),                                      \
        std::string((const char*) decoded_msg.protocolconfigurationoptions     \
                        .protocol_or_container_ids[1]                          \
                        .contents->data));                                     \
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

  void fill_epsqos(EpsQualityOfService* epsqos) {
    epsqos->bitRatesPresent               = 1;
    epsqos->bitRatesExtPresent            = 1;
    epsqos->bitRatesExt2Present           = 1;
    epsqos->bitRates.guarBitRateForDL     = 10;
    epsqos->bitRates.guarBitRateForUL     = 5;
    epsqos->bitRates.maxBitRateForDL      = 100;
    epsqos->bitRates.maxBitRateForUL      = 50;
    epsqos->bitRatesExt.guarBitRateForDL  = 10;
    epsqos->bitRatesExt.guarBitRateForUL  = 5;
    epsqos->bitRatesExt.maxBitRateForDL   = 100;
    epsqos->bitRatesExt.maxBitRateForUL   = 50;
    epsqos->bitRatesExt2.guarBitRateForDL = 10;
    epsqos->bitRatesExt2.guarBitRateForUL = 5;
    epsqos->bitRatesExt2.maxBitRateForDL  = 100;
    epsqos->bitRatesExt2.maxBitRateForUL  = 50;
    epsqos->qci                           = 5;
    return;
  }

  void fill_tft(traffic_flow_template_t* tft) {
    tft->tftoperationcode = TRAFFIC_FLOW_TEMPLATE_OPCODE_CREATE_NEW_TFT;
    tft->ebit = TRAFFIC_FLOW_TEMPLATE_PARAMETER_LIST_IS_NOT_INCLUDED;
    tft->numberofpacketfilters = 1;
    tft->packetfilterlist.createnewtft[0].direction =
        TRAFFIC_FLOW_TEMPLATE_UPLINK_ONLY;
    tft->packetfilterlist.createnewtft[0].eval_precedence = 250;
    tft->packetfilterlist.createnewtft[0].packetfiltercontents.flags |=
        (TRAFFIC_FLOW_TEMPLATE_IPV4_REMOTE_ADDR_FLAG |
         TRAFFIC_FLOW_TEMPLATE_SINGLE_REMOTE_PORT_FLAG);
    for (int i = 0; i < TRAFFIC_FLOW_TEMPLATE_IPV4_ADDR_SIZE; ++i) {
      tft->packetfilterlist.createnewtft[0]
          .packetfiltercontents.ipv4remoteaddr[i]
          .addr = 8;
      tft->packetfilterlist.createnewtft[0]
          .packetfiltercontents.ipv4remoteaddr[i]
          .mask = 255;
    }
    tft->packetfilterlist.createnewtft[0]
        .packetfiltercontents.singleremoteport = 80;
    return;
  }

  void fill_negotiated_qos(quality_of_service_t* negotiated_qos) {
    negotiated_qos->delayclass                 = 1;
    negotiated_qos->reliabilityclass           = 2;
    negotiated_qos->peakthroughput             = 15;
    negotiated_qos->precedenceclass            = 2;
    negotiated_qos->meanthroughput             = 16;
    negotiated_qos->trafficclass               = 7;
    negotiated_qos->deliveryorder              = 1;
    negotiated_qos->deliveryoferroneoussdu     = 2;
    negotiated_qos->maximumsdusize             = 255;
    negotiated_qos->maximumbitrateuplink       = 99;
    negotiated_qos->maximumbitratedownlink     = 169;
    negotiated_qos->residualber                = 4;
    negotiated_qos->sduratioerror              = 12;
    negotiated_qos->transferdelay              = 61;
    negotiated_qos->traffichandlingpriority    = 0;
    negotiated_qos->guaranteedbitrateuplink    = 10;
    negotiated_qos->guaranteedbitratedownlink  = 20;
    negotiated_qos->signalingindication        = 1;
    negotiated_qos->sourcestatisticsdescriptor = 13;
    return;
  }

  void fill_apnambr(ApnAggregateMaximumBitRate* apnambr) {
    apnambr->apnambrfordownlink           = 11;
    apnambr->apnambrforuplink             = 11;
    apnambr->apnambrfordownlink_extended  = 11;
    apnambr->apnambrforuplink_extended    = 11;
    apnambr->apnambrfordownlink_extended2 = 11;
    apnambr->apnambrforuplink_extended2   = 11;
    apnambr->extensions                   = 3;
  }

  uint8_t buffer[BUFFER_LEN];
};

TEST_F(ESMEncodeDecodeTest, TestActivateDefaultEpsBearerContextRequest) {
  activate_default_eps_bearer_context_request_msg original_msg = {0};
  activate_default_eps_bearer_context_request_msg decoded_msg  = {0};
  FILL_COMMON_MANDATORY_DEFAULTS(original_msg);

  fill_epsqos(&original_msg.epsqos);

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

  fill_negotiated_qos(&original_msg.negotiatedqos);

  original_msg.negotiatedllcsapi    = 10;
  original_msg.radiopriority        = 5;
  original_msg.packetflowidentifier = 118;

  fill_apnambr(&original_msg.apnambr);

  original_msg.esmcause = 102;

  fill_pco(&original_msg.protocolconfigurationoptions);

  int encoded = encode_activate_default_eps_bearer_context_request(
      &original_msg, buffer, BUFFER_LEN);
  int decoded = decode_activate_default_eps_bearer_context_request(
      &decoded_msg, buffer, encoded);

  EXPECT_EQ(encoded, decoded);
  // TODO(@ulaskozat): Header will be decoded separately; then the following
  // line can be uncommented.
  // COMPARE_COMMON_MANDATORY_DEFAULTS();
  EXPECT_EQ(
      original_msg.pdnaddress.pdntypevalue,
      decoded_msg.pdnaddress.pdntypevalue);
  EXPECT_TRUE(!memcmp(
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
  EXPECT_TRUE(!memcmp(
      &original_msg.negotiatedqos, &decoded_msg.negotiatedqos,
      sizeof(original_msg.negotiatedqos)));
  EXPECT_EQ(original_msg.negotiatedllcsapi, decoded_msg.negotiatedllcsapi);
  EXPECT_EQ(original_msg.radiopriority, decoded_msg.radiopriority);
  EXPECT_EQ(
      original_msg.packetflowidentifier, decoded_msg.packetflowidentifier);
  EXPECT_TRUE(!memcmp(
      &original_msg.apnambr, &decoded_msg.apnambr,
      sizeof(original_msg.apnambr)));

  EXPECT_EQ(original_msg.esmcause, decoded_msg.esmcause);
  EXPECT_EQ_PCO();

  DESTROY_PCO();
  bdestroy_wrapper(&original_msg.accesspointname);
  bdestroy_wrapper(&decoded_msg.accesspointname);
  bdestroy_wrapper(&original_msg.pdnaddress.pdnaddressinformation);
  bdestroy_wrapper(&decoded_msg.pdnaddress.pdnaddressinformation);
}

TEST_F(ESMEncodeDecodeTest, TestActivateDefaultEpsBearerContextAccept) {
  activate_default_eps_bearer_context_accept_msg original_msg = {0};
  activate_default_eps_bearer_context_accept_msg decoded_msg  = {0};
  FILL_COMMON_MANDATORY_DEFAULTS(original_msg);
  original_msg.presencemask =
      ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_ACCEPT_PROTOCOL_CONFIGURATION_OPTIONS_PRESENT;

  fill_pco(&original_msg.protocolconfigurationoptions);

  int encoded = encode_activate_default_eps_bearer_context_accept(
      &original_msg, buffer, BUFFER_LEN);
  int decoded = decode_activate_default_eps_bearer_context_accept(
      &decoded_msg, buffer, encoded);

  EXPECT_EQ(encoded, decoded);
  // TODO(@ulaskozat): Header will be decoded separately; then the following
  // line can be uncommented.
  // COMPARE_COMMON_MANDATORY_DEFAULTS();
  EXPECT_EQ(original_msg.presencemask, decoded_msg.presencemask);
  EXPECT_EQ_PCO();

  DESTROY_PCO();
}

TEST_F(ESMEncodeDecodeTest, TestActivateDefaultEpsBearerContextReject) {
  activate_default_eps_bearer_context_reject_msg original_msg = {0};
  activate_default_eps_bearer_context_reject_msg decoded_msg  = {0};
  FILL_COMMON_MANDATORY_DEFAULTS(original_msg);
  original_msg.presencemask =
      ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_REJECT_PROTOCOL_CONFIGURATION_OPTIONS_PRESENT;

  fill_pco(&original_msg.protocolconfigurationoptions);

  int encoded = encode_activate_default_eps_bearer_context_reject(
      &original_msg, buffer, BUFFER_LEN);
  int decoded = decode_activate_default_eps_bearer_context_reject(
      &decoded_msg, buffer, encoded);

  EXPECT_EQ(encoded, decoded);
  // TODO(@ulaskozat): Header will be decoded separately; then the following
  // line can be uncommented.
  // COMPARE_COMMON_MANDATORY_DEFAULTS();
  EXPECT_EQ(original_msg.presencemask, decoded_msg.presencemask);
  EXPECT_EQ_PCO();

  DESTROY_PCO();
}

TEST_F(ESMEncodeDecodeTest, TestActivateDedicatedEpsBearerContextRequest) {
  activate_dedicated_eps_bearer_context_request_msg original_msg = {0};
  activate_dedicated_eps_bearer_context_request_msg decoded_msg  = {0};
  FILL_COMMON_MANDATORY_DEFAULTS(original_msg);

  original_msg.linkedepsbeareridentity = 5;

  fill_epsqos(&original_msg.epsqos);

  fill_tft(&original_msg.tft);

  // TI_IE encoding/decoding is not implemented, hence its flag as well as
  // entries will be omitted.
  original_msg.presencemask =
      ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REQUEST_NEGOTIATED_QOS_PRESENT |
      ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REQUEST_NEGOTIATED_LLC_SAPI_PRESENT |
      ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REQUEST_RADIO_PRIORITY_PRESENT |
      ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REQUEST_PACKET_FLOW_IDENTIFIER_PRESENT |
      ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REQUEST_PROTOCOL_CONFIGURATION_OPTIONS_PRESENT;

  fill_negotiated_qos(&original_msg.negotiatedqos);

  original_msg.negotiatedllcsapi    = 10;
  original_msg.radiopriority        = 5;
  original_msg.packetflowidentifier = 118;

  fill_pco(&original_msg.protocolconfigurationoptions);

  int encoded = encode_activate_dedicated_eps_bearer_context_request(
      &original_msg, buffer, BUFFER_LEN);
  int decoded = decode_activate_dedicated_eps_bearer_context_request(
      &decoded_msg, buffer, encoded);

  EXPECT_EQ(encoded, decoded);
  // TODO(@ulaskozat): Header will be decoded separately; then the following
  // line can be uncommented.
  // COMPARE_COMMON_MANDATORY_DEFAULTS();
  EXPECT_EQ(
      original_msg.linkedepsbeareridentity,
      decoded_msg.linkedepsbeareridentity);
  EXPECT_TRUE(!memcmp(
      &original_msg.epsqos, &decoded_msg.epsqos, sizeof(original_msg.epsqos)));
  EXPECT_TRUE(
      !memcmp(&original_msg.tft, &decoded_msg.tft, sizeof(original_msg.tft)));
  EXPECT_EQ(original_msg.presencemask, decoded_msg.presencemask);
  EXPECT_TRUE(!memcmp(
      &original_msg.negotiatedqos, &decoded_msg.negotiatedqos,
      sizeof(original_msg.negotiatedqos)));
  EXPECT_EQ(original_msg.negotiatedllcsapi, decoded_msg.negotiatedllcsapi);
  EXPECT_EQ(original_msg.radiopriority, decoded_msg.radiopriority);
  EXPECT_EQ(
      original_msg.packetflowidentifier, decoded_msg.packetflowidentifier);
  EXPECT_EQ_PCO();

  DESTROY_PCO();
}

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
  // TODO(@ulaskozat): Header will be decoded separately; then the following
  // line can be uncommented.
  // COMPARE_COMMON_MANDATORY_DEFAULTS();
  EXPECT_EQ(original_msg.presencemask, decoded_msg.presencemask);
  EXPECT_EQ_PCO();

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
  // TODO(@ulaskozat): Header will be decoded separately; then the following
  // line can be uncommented.
  // COMPARE_COMMON_MANDATORY_DEFAULTS();
  EXPECT_EQ(original_msg.esmcause, decoded_msg.esmcause);
  EXPECT_EQ(original_msg.presencemask, decoded_msg.presencemask);
  EXPECT_EQ_PCO();

  DESTROY_PCO();
}

TEST_F(ESMEncodeDecodeTest, TestBearerResourceAllocationRequest) {
  bearer_resource_allocation_request_msg original_msg = {0};
  bearer_resource_allocation_request_msg decoded_msg  = {0};
  FILL_COMMON_MANDATORY_DEFAULTS(original_msg);

  original_msg.linkedepsbeareridentity = 5;

  fill_epsqos(&original_msg.requiredtrafficflowqos);

  fill_tft(&original_msg.trafficflowaggregate);

  original_msg.presencemask =
      BEARER_RESOURCE_ALLOCATION_REQUEST_PROTOCOL_CONFIGURATION_OPTIONS_PRESENT;

  fill_pco(&original_msg.protocolconfigurationoptions);

  int encoded = encode_bearer_resource_allocation_request(
      &original_msg, buffer, BUFFER_LEN);
  int decoded =
      decode_bearer_resource_allocation_request(&decoded_msg, buffer, encoded);

  EXPECT_EQ(encoded, decoded);
  // TODO(@ulaskozat): Header will be decoded separately; then the following
  // line can be uncommented.
  // COMPARE_COMMON_MANDATORY_DEFAULTS();
  EXPECT_EQ(original_msg.presencemask, decoded_msg.presencemask);
  EXPECT_EQ(
      original_msg.linkedepsbeareridentity,
      decoded_msg.linkedepsbeareridentity);
  EXPECT_TRUE(!memcmp(
      &original_msg.requiredtrafficflowqos, &decoded_msg.requiredtrafficflowqos,
      sizeof(original_msg.requiredtrafficflowqos)));
  EXPECT_TRUE(!memcmp(
      &original_msg.trafficflowaggregate, &decoded_msg.trafficflowaggregate,
      sizeof(original_msg.trafficflowaggregate)));
  EXPECT_EQ_PCO();

  DESTROY_PCO();
}

TEST_F(ESMEncodeDecodeTest, TestBearerResourceAllocationReject) {
  bearer_resource_allocation_reject_msg original_msg = {0};
  bearer_resource_allocation_reject_msg decoded_msg  = {0};
  FILL_COMMON_MANDATORY_DEFAULTS(original_msg);
  original_msg.presencemask =
      BEARER_RESOURCE_ALLOCATION_REJECT_PROTOCOL_CONFIGURATION_OPTIONS_PRESENT;

  fill_pco(&original_msg.protocolconfigurationoptions);

  int encoded = encode_bearer_resource_allocation_reject(
      &original_msg, buffer, BUFFER_LEN);
  int decoded =
      decode_bearer_resource_allocation_reject(&decoded_msg, buffer, encoded);

  EXPECT_EQ(encoded, decoded);
  // TODO(@ulaskozat): Header will be decoded separately; then the following
  // line can be uncommented.
  // COMPARE_COMMON_MANDATORY_DEFAULTS();
  EXPECT_EQ(original_msg.presencemask, decoded_msg.presencemask);
  EXPECT_EQ_PCO();

  DESTROY_PCO();
}

TEST_F(ESMEncodeDecodeTest, TestBearerResourceModificationRequest) {
  bearer_resource_modification_request_msg original_msg = {0};
  bearer_resource_modification_request_msg decoded_msg  = {0};
  FILL_COMMON_MANDATORY_DEFAULTS(original_msg);

  original_msg.epsbeareridentityforpacketfilter = 8;

  fill_epsqos(&original_msg.requiredtrafficflowqos);

  fill_tft(&original_msg.trafficflowaggregate);

  original_msg.presencemask =
      BEARER_RESOURCE_MODIFICATION_REQUEST_REQUIRED_TRAFFIC_FLOW_QOS_PRESENT |
      BEARER_RESOURCE_MODIFICATION_REQUEST_ESM_CAUSE_PRESENT |
      BEARER_RESOURCE_MODIFICATION_REQUEST_PROTOCOL_CONFIGURATION_OPTIONS_PRESENT;

  original_msg.esmcause = 102;

  fill_pco(&original_msg.protocolconfigurationoptions);

  int encoded = encode_bearer_resource_modification_request(
      &original_msg, buffer, BUFFER_LEN);
  int decoded = decode_bearer_resource_modification_request(
      &decoded_msg, buffer, encoded);

  EXPECT_EQ(encoded, decoded);
  // TODO(@ulaskozat): Header will be decoded separately; then the following
  // line can be uncommented.
  // COMPARE_COMMON_MANDATORY_DEFAULTS();
  EXPECT_EQ(original_msg.presencemask, decoded_msg.presencemask);
  EXPECT_EQ(
      original_msg.epsbeareridentityforpacketfilter,
      decoded_msg.epsbeareridentityforpacketfilter);
  EXPECT_TRUE(!memcmp(
      &original_msg.requiredtrafficflowqos, &decoded_msg.requiredtrafficflowqos,
      sizeof(original_msg.requiredtrafficflowqos)));
  EXPECT_TRUE(!memcmp(
      &original_msg.trafficflowaggregate, &decoded_msg.trafficflowaggregate,
      sizeof(original_msg.trafficflowaggregate)));
  EXPECT_EQ(original_msg.esmcause, decoded_msg.esmcause);
  EXPECT_EQ_PCO();

  DESTROY_PCO();
}

TEST_F(ESMEncodeDecodeTest, TestBearerResourceModificationReject) {
  bearer_resource_modification_reject_msg original_msg = {0};
  bearer_resource_modification_reject_msg decoded_msg  = {0};
  FILL_COMMON_MANDATORY_DEFAULTS(original_msg);

  original_msg.esmcause = 102;

  original_msg.presencemask =
      BEARER_RESOURCE_MODIFICATION_REJECT_PROTOCOL_CONFIGURATION_OPTIONS_PRESENT;

  fill_pco(&original_msg.protocolconfigurationoptions);

  int encoded = encode_bearer_resource_modification_reject(
      &original_msg, buffer, BUFFER_LEN);
  int decoded =
      decode_bearer_resource_modification_reject(&decoded_msg, buffer, encoded);

  EXPECT_EQ(encoded, decoded);
  // TODO(@ulaskozat): Header will be decoded separately; then the following
  // line can be uncommented.
  // COMPARE_COMMON_MANDATORY_DEFAULTS();
  EXPECT_EQ(original_msg.esmcause, decoded_msg.esmcause);
  EXPECT_EQ(original_msg.presencemask, decoded_msg.presencemask);
  EXPECT_EQ_PCO();

  DESTROY_PCO();
}

TEST_F(ESMEncodeDecodeTest, TestDeactivateEpsBearerContextRequest) {
  deactivate_eps_bearer_context_request_msg original_msg = {0};
  deactivate_eps_bearer_context_request_msg decoded_msg  = {0};
  FILL_COMMON_MANDATORY_DEFAULTS(original_msg);

  original_msg.esmcause = 102;

  original_msg.presencemask =
      DEACTIVATE_EPS_BEARER_CONTEXT_REQUEST_PROTOCOL_CONFIGURATION_OPTIONS_PRESENT;

  fill_pco(&original_msg.protocolconfigurationoptions);

  int encoded = encode_deactivate_eps_bearer_context_request(
      &original_msg, buffer, BUFFER_LEN);
  int decoded = decode_deactivate_eps_bearer_context_request(
      &decoded_msg, buffer, encoded);

  EXPECT_EQ(encoded, decoded);
  // TODO(@ulaskozat): Header will be decoded separately; then the following
  // line can be uncommented.
  // COMPARE_COMMON_MANDATORY_DEFAULTS();
  EXPECT_EQ(original_msg.presencemask, decoded_msg.presencemask);
  EXPECT_EQ(original_msg.esmcause, decoded_msg.esmcause);
  EXPECT_EQ_PCO();

  DESTROY_PCO();
}

TEST_F(ESMEncodeDecodeTest, TestDeactivateEpsBearerContextAccept) {
  deactivate_eps_bearer_context_accept_msg original_msg = {0};
  deactivate_eps_bearer_context_accept_msg decoded_msg  = {0};
  FILL_COMMON_MANDATORY_DEFAULTS(original_msg);

  original_msg.presencemask =
      DEACTIVATE_EPS_BEARER_CONTEXT_ACCEPT_PROTOCOL_CONFIGURATION_OPTIONS_PRESENT;

  fill_pco(&original_msg.protocolconfigurationoptions);

  int encoded = encode_deactivate_eps_bearer_context_accept(
      &original_msg, buffer, BUFFER_LEN);
  int decoded = decode_deactivate_eps_bearer_context_accept(
      &decoded_msg, buffer, encoded);

  EXPECT_EQ(encoded, decoded);
  // TODO(@ulaskozat): Header will be decoded separately; then the following
  // line can be uncommented.
  // COMPARE_COMMON_MANDATORY_DEFAULTS();
  EXPECT_EQ(original_msg.presencemask, decoded_msg.presencemask);
  EXPECT_EQ_PCO();

  DESTROY_PCO();
}

TEST_F(ESMEncodeDecodeTest, TestEsmInformationRequest) {
  esm_information_request_msg original_msg = {0};
  esm_information_request_msg decoded_msg  = {0};
  FILL_COMMON_MANDATORY_DEFAULTS(original_msg);

  int encoded =
      encode_esm_information_request(&original_msg, buffer, BUFFER_LEN);
  int decoded = decode_esm_information_request(&decoded_msg, buffer, encoded);

  EXPECT_EQ(encoded, decoded);
  // Currently encoding/decoding functions are dummy with no actual
  // encoding/decoding As they are the same as message header fields.
  // TODO(@ulaskozat): Header will be decoded separately; then the following
  // line can be uncommented.
  // COMPARE_COMMON_MANDATORY_DEFAULTS();
}

TEST_F(ESMEncodeDecodeTest, TestEsmInformationResponse) {
  esm_information_response_msg original_msg = {0};
  esm_information_response_msg decoded_msg  = {0};
  FILL_COMMON_MANDATORY_DEFAULTS(original_msg);

  original_msg.presencemask =
      ESM_INFORMATION_RESPONSE_PROTOCOL_CONFIGURATION_OPTIONS_PRESENT |
      ESM_INFORMATION_RESPONSE_ACCESS_POINT_NAME_PRESENT;

  original_msg.accesspointname = bfromcstr("magma.ipv4");

  fill_pco(&original_msg.protocolconfigurationoptions);

  int encoded =
      encode_esm_information_response(&original_msg, buffer, BUFFER_LEN);
  int decoded = decode_esm_information_response(&decoded_msg, buffer, encoded);

  EXPECT_EQ(encoded, decoded);
  // TODO(@ulaskozat): Header will be decoded separately; then the following
  // line can be uncommented.
  // COMPARE_COMMON_MANDATORY_DEFAULTS();
  EXPECT_EQ(original_msg.presencemask, decoded_msg.presencemask);
  EXPECT_EQ(
      std::string((const char*) original_msg.accesspointname->data),
      std::string((const char*) decoded_msg.accesspointname->data));
  EXPECT_EQ_PCO();

  DESTROY_PCO();
  bdestroy_wrapper(&original_msg.accesspointname);
  bdestroy_wrapper(&decoded_msg.accesspointname);
}

TEST_F(ESMEncodeDecodeTest, TestEsmStatus) {
  esm_status_msg original_msg = {0};
  esm_status_msg decoded_msg  = {0};
  FILL_COMMON_MANDATORY_DEFAULTS(original_msg);

  original_msg.esmcause = 102;

  int encoded = encode_esm_status(&original_msg, buffer, BUFFER_LEN);
  int decoded = decode_esm_status(&decoded_msg, buffer, encoded);

  EXPECT_EQ(encoded, decoded);
  // TODO(@ulaskozat): Header will be decoded separately; then the following
  // line can be uncommented.
  // COMPARE_COMMON_MANDATORY_DEFAULTS();
  EXPECT_EQ(original_msg.esmcause, decoded_msg.esmcause);
  EXPECT_EQ(encoded, decoded);
}

TEST_F(ESMEncodeDecodeTest, TestModifyEpsBearerContextRequest) {
  modify_eps_bearer_context_request_msg original_msg = {0};
  modify_eps_bearer_context_request_msg decoded_msg  = {0};
  FILL_COMMON_MANDATORY_DEFAULTS(original_msg);

  original_msg.presencemask =
      MODIFY_EPS_BEARER_CONTEXT_REQUEST_NEW_EPS_QOS_PRESENT |
      MODIFY_EPS_BEARER_CONTEXT_REQUEST_TFT_PRESENT |
      MODIFY_EPS_BEARER_CONTEXT_REQUEST_NEW_QOS_PRESENT |
      MODIFY_EPS_BEARER_CONTEXT_REQUEST_NEGOTIATED_LLC_SAPI_PRESENT |
      MODIFY_EPS_BEARER_CONTEXT_REQUEST_RADIO_PRIORITY_PRESENT |
      MODIFY_EPS_BEARER_CONTEXT_REQUEST_PACKET_FLOW_IDENTIFIER_PRESENT |
      MODIFY_EPS_BEARER_CONTEXT_REQUEST_APNAMBR_PRESENT |
      MODIFY_EPS_BEARER_CONTEXT_REQUEST_PROTOCOL_CONFIGURATION_OPTIONS_PRESENT;

  fill_epsqos(&original_msg.newepsqos);

  fill_tft(&original_msg.tft);

  fill_negotiated_qos(&original_msg.newqos);

  original_msg.negotiatedllcsapi    = 10;
  original_msg.radiopriority        = 5;
  original_msg.packetflowidentifier = 118;

  fill_apnambr(&original_msg.apnambr);

  fill_pco(&original_msg.protocolconfigurationoptions);

  int encoded = encode_modify_eps_bearer_context_request(
      &original_msg, buffer, BUFFER_LEN);
  int decoded =
      decode_modify_eps_bearer_context_request(&decoded_msg, buffer, encoded);

  EXPECT_EQ(encoded, decoded);
  // TODO(@ulaskozat): Header will be decoded separately; then the following
  // line can be uncommented.
  // COMPARE_COMMON_MANDATORY_DEFAULTS();
  EXPECT_EQ(original_msg.presencemask, decoded_msg.presencemask);
  EXPECT_TRUE(!memcmp(
      &original_msg.newepsqos, &decoded_msg.newepsqos,
      sizeof(original_msg.newepsqos)));
  EXPECT_TRUE(
      !memcmp(&original_msg.tft, &decoded_msg.tft, sizeof(original_msg.tft)));
  EXPECT_TRUE(!memcmp(
      &original_msg.newqos, &decoded_msg.newqos, sizeof(original_msg.newqos)));
  EXPECT_EQ(original_msg.negotiatedllcsapi, decoded_msg.negotiatedllcsapi);
  EXPECT_EQ(original_msg.radiopriority, decoded_msg.radiopriority);
  EXPECT_EQ(
      original_msg.packetflowidentifier, decoded_msg.packetflowidentifier);
  EXPECT_TRUE(!memcmp(
      &original_msg.apnambr, &decoded_msg.apnambr,
      sizeof(original_msg.apnambr)));
  EXPECT_EQ_PCO();

  DESTROY_PCO();
}

TEST_F(ESMEncodeDecodeTest, TestModifyEpsBearerContextAccept) {
  modify_eps_bearer_context_accept_msg original_msg = {0};
  modify_eps_bearer_context_accept_msg decoded_msg  = {0};
  FILL_COMMON_MANDATORY_DEFAULTS(original_msg);

  original_msg.presencemask =
      MODIFY_EPS_BEARER_CONTEXT_ACCEPT_PROTOCOL_CONFIGURATION_OPTIONS_PRESENT;

  fill_pco(&original_msg.protocolconfigurationoptions);

  int encoded = encode_modify_eps_bearer_context_accept(
      &original_msg, buffer, BUFFER_LEN);
  int decoded =
      decode_modify_eps_bearer_context_accept(&decoded_msg, buffer, encoded);

  EXPECT_EQ(encoded, decoded);
  // TODO(@ulaskozat): Header will be decoded separately; then the following
  // line can be uncommented.
  // COMPARE_COMMON_MANDATORY_DEFAULTS();
  EXPECT_EQ(original_msg.presencemask, decoded_msg.presencemask);
  EXPECT_EQ_PCO();

  DESTROY_PCO();
}

TEST_F(ESMEncodeDecodeTest, TestModifyEpsBearerContextReject) {
  modify_eps_bearer_context_reject_msg original_msg = {0};
  modify_eps_bearer_context_reject_msg decoded_msg  = {0};
  FILL_COMMON_MANDATORY_DEFAULTS(original_msg);

  original_msg.esmcause = 102;

  original_msg.presencemask =
      BEARER_RESOURCE_MODIFICATION_REJECT_PROTOCOL_CONFIGURATION_OPTIONS_PRESENT;

  fill_pco(&original_msg.protocolconfigurationoptions);

  int encoded = encode_modify_eps_bearer_context_reject(
      &original_msg, buffer, BUFFER_LEN);
  int decoded =
      decode_modify_eps_bearer_context_reject(&decoded_msg, buffer, encoded);

  EXPECT_EQ(encoded, decoded);
  // TODO(@ulaskozat): Header will be decoded separately; then the following
  // line can be uncommented.
  // COMPARE_COMMON_MANDATORY_DEFAULTS();
  EXPECT_EQ(original_msg.esmcause, decoded_msg.esmcause);
  EXPECT_EQ(original_msg.presencemask, decoded_msg.presencemask);
  EXPECT_EQ_PCO();

  DESTROY_PCO();
}

TEST_F(ESMEncodeDecodeTest, TestPdnConnectivityRequest) {
  pdn_connectivity_request_msg original_msg = {0};
  pdn_connectivity_request_msg decoded_msg  = {0};
  FILL_COMMON_MANDATORY_DEFAULTS(original_msg);

  original_msg.requesttype = 3;
  original_msg.pdntype     = 2;

  original_msg.presencemask =
      PDN_CONNECTIVITY_REQUEST_ESM_INFORMATION_TRANSFER_FLAG_PRESENT |
      PDN_CONNECTIVITY_REQUEST_PROTOCOL_CONFIGURATION_OPTIONS_PRESENT;

  original_msg.esminformationtransferflag = 1;

  fill_pco(&original_msg.protocolconfigurationoptions);

  int encoded =
      encode_pdn_connectivity_request(&original_msg, buffer, BUFFER_LEN);
  int decoded = decode_pdn_connectivity_request(&decoded_msg, buffer, encoded);

  EXPECT_EQ(encoded, decoded);
  // TODO(@ulaskozat): Header will be decoded separately; then the following
  // line can be uncommented.
  // COMPARE_COMMON_MANDATORY_DEFAULTS();
  EXPECT_EQ(original_msg.requesttype, decoded_msg.requesttype);
  EXPECT_EQ(original_msg.pdntype, decoded_msg.pdntype);
  EXPECT_EQ(original_msg.presencemask, decoded_msg.presencemask);
  EXPECT_EQ(
      original_msg.esminformationtransferflag,
      decoded_msg.esminformationtransferflag);
  EXPECT_EQ_PCO();

  DESTROY_PCO();
}

TEST_F(ESMEncodeDecodeTest, TestPdnConnectivityReject) {
  pdn_connectivity_reject_msg original_msg = {0};
  pdn_connectivity_reject_msg decoded_msg  = {0};
  FILL_COMMON_MANDATORY_DEFAULTS(original_msg);

  original_msg.esmcause = 102;

  original_msg.presencemask =
      PDN_CONNECTIVITY_REJECT_PROTOCOL_CONFIGURATION_OPTIONS_PRESENT;

  fill_pco(&original_msg.protocolconfigurationoptions);

  int encoded =
      encode_pdn_connectivity_reject(&original_msg, buffer, BUFFER_LEN);
  int decoded = decode_pdn_connectivity_reject(&decoded_msg, buffer, encoded);

  EXPECT_EQ(encoded, decoded);
  // TODO(@ulaskozat): Header will be decoded separately; then the following
  // line can be uncommented.
  // COMPARE_COMMON_MANDATORY_DEFAULTS();
  EXPECT_EQ(original_msg.esmcause, decoded_msg.esmcause);
  EXPECT_EQ(original_msg.presencemask, decoded_msg.presencemask);
  EXPECT_EQ_PCO();

  DESTROY_PCO();
}

TEST_F(ESMEncodeDecodeTest, TestPdnDisconnectRequest) {
  pdn_disconnect_request_msg original_msg = {0};
  pdn_disconnect_request_msg decoded_msg  = {0};
  FILL_COMMON_MANDATORY_DEFAULTS(original_msg);

  original_msg.linkedepsbeareridentity = 5;

  original_msg.presencemask =
      PDN_DISCONNECT_REQUEST_PROTOCOL_CONFIGURATION_OPTIONS_PRESENT;

  fill_pco(&original_msg.protocolconfigurationoptions);

  int encoded =
      encode_pdn_disconnect_request(&original_msg, buffer, BUFFER_LEN);
  int decoded = decode_pdn_disconnect_request(&decoded_msg, buffer, encoded);

  EXPECT_EQ(encoded, decoded);
  // TODO(@ulaskozat): Header will be decoded separately; then the following
  // line can be uncommented.
  // COMPARE_COMMON_MANDATORY_DEFAULTS();
  EXPECT_EQ(
      original_msg.linkedepsbeareridentity,
      decoded_msg.linkedepsbeareridentity);
  EXPECT_EQ(original_msg.presencemask, decoded_msg.presencemask);
  EXPECT_EQ_PCO();

  DESTROY_PCO();
}

TEST_F(ESMEncodeDecodeTest, TestPdnDisconnectReject) {
  pdn_disconnect_reject_msg original_msg = {0};
  pdn_disconnect_reject_msg decoded_msg  = {0};
  FILL_COMMON_MANDATORY_DEFAULTS(original_msg);

  original_msg.esmcause = 102;

  original_msg.presencemask =
      PDN_DISCONNECT_REJECT_PROTOCOL_CONFIGURATION_OPTIONS_PRESENT;

  fill_pco(&original_msg.protocolconfigurationoptions);

  int encoded = encode_pdn_disconnect_reject(&original_msg, buffer, BUFFER_LEN);
  int decoded = decode_pdn_disconnect_reject(&decoded_msg, buffer, encoded);

  EXPECT_EQ(encoded, decoded);
  // TODO(@ulaskozat): Header will be decoded separately; then the following
  // line can be uncommented.
  // COMPARE_COMMON_MANDATORY_DEFAULTS();
  EXPECT_EQ(original_msg.esmcause, decoded_msg.esmcause);
  EXPECT_EQ(original_msg.presencemask, decoded_msg.presencemask);
  EXPECT_EQ_PCO();

  DESTROY_PCO();
}

}  // namespace lte
}  // namespace magma
