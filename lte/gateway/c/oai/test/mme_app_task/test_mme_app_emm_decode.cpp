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
#include "AttachRequest.h"
#include "dynamic_memory_check.h"
#include "log.h"
}

class EMMDecodeTest : public ::testing::Test {
  virtual void SetUp() {}
  virtual void TearDown() {}
};

TEST_F(EMMDecodeTest, TestDecodeAttachRequest1) {
  //   Combined attach, NAS message generated from s1ap tester
  uint8_t buffer[] = {0x72, 0x08, 0x09, 0x10, 0x10, 0x00, 0x00, 0x00,
                      0x00, 0x10, 0x02, 0xe0, 0xe0, 0x00, 0x04, 0x02,
                      0x01, 0xd0, 0x11, 0x40, 0x08, 0x04, 0x02, 0x60,
                      0x04, 0x00, 0x02, 0x1c, 0x00};
  uint32_t len     = 29;
  attach_request_msg attach_request;

  int rc = decode_attach_request(&attach_request, buffer, len);
  ASSERT_EQ(rc, len);
  ASSERT_EQ(attach_request.epsattachtype, 2);
  ASSERT_EQ(attach_request.naskeysetidentifier.naskeysetidentifier, 7);
  ASSERT_EQ(attach_request.naskeysetidentifier.tsc, 0);

  bdestroy_wrapper(&attach_request.esmmessagecontainer);
}

TEST_F(EMMDecodeTest, TestDecodeAttachRequest2) {
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
  ASSERT_EQ(attach_request.epsattachtype, 2);
  ASSERT_EQ(attach_request.naskeysetidentifier.naskeysetidentifier, 7);
  ASSERT_EQ(attach_request.naskeysetidentifier.tsc, 0);

  bdestroy_wrapper(&attach_request.esmmessagecontainer);
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  OAILOG_INIT("MME", OAILOG_LEVEL_DEBUG, MAX_LOG_PROTOS);
  return RUN_ALL_TESTS();
}
