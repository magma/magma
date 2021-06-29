#include <gtest/gtest.h>
#include "orc8r/protos/common.pb.h"

// Demonstrate some basic assertions. TODO delete later!
TEST(HelloTest, BasicAssertions) {
  auto network = magma::orc8r::NetworkID();
  network.set_id("networkID!");
  EXPECT_EQ("networkID!", network.id());
  // Expect two strings not to be equal.
  EXPECT_STRNE("hello", "world");
  // Expect equality.
  EXPECT_EQ(7 * 6, 42);
}
