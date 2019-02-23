/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#include <gtest/gtest.h>

#include "MagmaService.h"

using ::testing::Test;

namespace magma { namespace service303 {

const std::string MAGMA_SERVICE_NAME = "test_service";
const std::string MAGMA_SERVICE_VERSION = "0.0.0";
const std::string META_KEY = "key";
const std::string META_VALUE = "value";

TEST(test_magma_service, test_GetServiceInfo) {
  MagmaService magma_service(MAGMA_SERVICE_NAME, MAGMA_SERVICE_VERSION);
  ServiceInfo response;

  magma_service.GetServiceInfo(nullptr, nullptr, &response);
  EXPECT_EQ(response.name(), MAGMA_SERVICE_NAME);
  EXPECT_EQ(response.version(), MAGMA_SERVICE_VERSION);
  EXPECT_EQ(response.state(), ServiceInfo::ALIVE);
  EXPECT_EQ(response.health(), ServiceInfo::APP_UNKNOWN);
  auto start_time_1 = response.start_time_secs();
  EXPECT_TRUE(response.status().meta().empty());

  response.Clear();

  magma_service.setApplicationHealth(ServiceInfo::APP_HEALTHY);
  magma_service.GetServiceInfo(nullptr, nullptr, &response);
  EXPECT_EQ(response.name(), MAGMA_SERVICE_NAME);
  EXPECT_EQ(response.version(), MAGMA_SERVICE_VERSION);
  EXPECT_EQ(response.state(), ServiceInfo::ALIVE);
  EXPECT_EQ(response.health(), ServiceInfo::APP_HEALTHY);
  auto start_time_2 = response.start_time_secs();
  EXPECT_TRUE(response.status().meta().empty());

  EXPECT_EQ(start_time_1, start_time_2);
}

ServiceInfoMeta test_callback() {
  return ServiceInfoMeta {{META_KEY, META_VALUE}};
}

TEST(test_magma_service, test_GetServiceInfo_with_callback) {
  MagmaService magma_service(MAGMA_SERVICE_NAME, MAGMA_SERVICE_VERSION);
  ServiceInfo response;

  magma_service.GetServiceInfo(nullptr, nullptr, &response);
  EXPECT_TRUE(response.status().meta().empty());

  response.Clear();

  magma_service.SetServiceInfoCallback(test_callback);
  magma_service.GetServiceInfo(nullptr, nullptr, &response);
  auto meta = response.status().meta();
  EXPECT_FALSE(meta.empty());
  EXPECT_EQ(meta.size(), 1);
  EXPECT_EQ(meta[META_KEY], META_VALUE);

  response.Clear();

  magma_service.ClearServiceInfoCallback();
  magma_service.GetServiceInfo(nullptr, nullptr, &response);
  EXPECT_TRUE(response.status().meta().empty());

}

int main(int argc, char **argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

}}
