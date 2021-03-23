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
#include <gtest/gtest-message.h>     // for Message
#include <gtest/gtest-param-test.h>  // for ParamIteratorInterface, Values
#include <gtest/gtest-test-part.h>   // for TestPartResult
#include <gtest/gtest.h>             // for InitGoogleTest, RUN_ALL_TESTS
#include <string>                    // for basic_string, string
#include "IMSIEncoder.h"             // for IMSIEncoder, openflow

using ::testing::Test;
using ::testing::Values;
using namespace openflow;

namespace {

class IMSIEncoderTest : public ::testing::TestWithParam<std::string> {};

/*
 * Test IMSI encoder within GTP application by encoding an IMSI to uint64_t and
 * back to see if the values match
 */
TEST_P(IMSIEncoderTest, TestCompactExpand) {
  std::string imsi_test =
      IMSIEncoder::expand_imsi(IMSIEncoder::compact_imsi(GetParam()));
  ASSERT_STREQ(GetParam().c_str(), imsi_test.c_str());
}

INSTANTIATE_TEST_CASE_P(
    TestLeadingZeros, IMSIEncoderTest,
    Values("001010000000013", "011010000000013", "111010000000013"));

INSTANTIATE_TEST_CASE_P(
    TestDifferentLengths, IMSIEncoderTest,
    Values("001010000000013", "01010000000013", "28950000000013"));

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

}  // namespace
