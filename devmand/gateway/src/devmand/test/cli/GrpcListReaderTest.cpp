// Copyright (c) 2020-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#define LOG_WITH_GLOG

#include <magma_logging.h>

#include <devmand/channels/cli/plugin/protocpp/ReaderPlugin.grpc.pb.h>
#include <devmand/devices/cli/translation/GrpcListReader.h>
#include <devmand/test/cli/utils/Log.h>
#include <devmand/test/cli/utils/MockCli.h>
#include <folly/executors/CPUThreadPoolExecutor.h>
#include <folly/json.h>
#include <grpc++/grpc++.h>
#include <gtest/gtest.h>
#include <thread>

namespace devmand {
namespace test {
namespace cli {

using namespace std;
using namespace folly;
using namespace grpc;
using namespace devmand::channels::cli;
using namespace devmand::devices::cli;
using namespace devmand::channels::cli::plugin;
using namespace devmand::test::utils::cli;

class GrpcListReaderTest : public ::testing::Test {
 protected:
  shared_ptr<CPUThreadPoolExecutor> testExec;
  shared_ptr<Cli> cli;

  void SetUp() override {
    devmand::test::utils::log::initLog();
    testExec = make_shared<CPUThreadPoolExecutor>(5);
    cli = make_shared<EchoCli>();
  }
};

static void sendActualResponse(
    string json,
    ServerReaderWriter<ReadResponse, ReadRequest>* stream) {
  ReadResponse readResponse;
  ActualReadResponse* actualReadResponse = new ActualReadResponse();
  actualReadResponse->set_json(json);
  readResponse.set_allocated_actualreadresponse(actualReadResponse);
  stream->Write(readResponse);
}

// request 3 commands, then send final response
class DummyListReader : public ReaderPlugin::Service {
  Status Read(
      ServerContext* context,
      ServerReaderWriter<ReadResponse, ReadRequest>* stream) {
    (void)context;
    ReadRequest readRequest;
    stream->Read(&readRequest);
    EXPECT_TRUE(readRequest.has_actualreadrequest());
    // send final response
    MLOG(MDEBUG) << "Sending actualReadResponse";
    sendActualResponse("[{\"x\":0},{\"x\":1},{\"x\":2}]", stream);
    return Status::OK;
  }
};

static std::unique_ptr<Server> startServer(
    const string& address,
    ReaderPlugin::Service& service) {
  std::string server_address(address);
  ServerBuilder builder;
  builder.AddListeningPort(server_address, grpc::InsecureServerCredentials());
  builder.RegisterService(&service);
  std::unique_ptr<Server> server(builder.BuildAndStart());
  return server;
}

TEST_F(GrpcListReaderTest, testDummyReader) {
  const unsigned int port = 50052;
  const string& address = "localhost:" + to_string(port);

  MLOG(MDEBUG) << "Starting server";
  DummyListReader service;
  std::unique_ptr<Server> server = startServer(address, service);

  auto grpcClientChannel =
      grpc::CreateChannel(address, grpc::InsecureChannelCredentials());
  GrpcListReader tested(grpcClientChannel, "tested", testExec);
  Path path = "/somepath";
  DeviceAccess deviceAccess = DeviceAccess(cli, "test", testExec);
  vector<dynamic> result = tested.readKeys(path, deviceAccess).get();
  EXPECT_EQ(result.size(), 3);
  for (int i = 0; i < 3; i++) {
    dynamic expected = dynamic::object;
    expected["x"] = i;
    dynamic item = *next(result.begin(), i);
    EXPECT_EQ(toJson(expected), toJson(item));
  }
}

} // namespace cli
} // namespace test
} // namespace devmand
