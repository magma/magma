/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <memory>

#include <glog/logging.h>
#include <gtest/gtest.h>

#include "RuleStore.h"
#include "SessionID.h"
#include "SessionState.h"
#include "MemoryStoreClient.h"
#include "magma_logging.h"

using ::testing::Test;

namespace magma {

class RedisStoreClientTest : public ::testing::Test {
 protected:
  SessionIDGenerator id_gen_;
};
}

/**
 * End to end test of the MemoryStoreClient.
 * 1) Create MemoryStoreClient
 * 2) Read in sessions for subscribers IMSI1 and IMSI2
 * 3) Create bare-bones session for IMSI1 and IMSI2
 * 4) Write and commit session state for IMSI1 and IMSI2
 * 5) Read for subscribers IMSI1 and IMSI2
 * 6) Verify that state was written for IMSI1/IMSI2 and has been retrieved.
 */
TEST_F(RedisStoreClientTest, test_read_and_write)
{
}

int main(int argc, char **argv)
{
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

} // namespace magma
