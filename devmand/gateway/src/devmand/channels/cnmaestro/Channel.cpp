// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <stdio.h>
#include <iostream>

#include <folly/Conv.h>
#include <folly/Subprocess.h>
#include <folly/executors/Async.h>
#include <folly/json.h>

#include <devmand/channels/cnmaestro/Channel.h>
#include <devmand/error/ErrorHandler.h>

namespace devmand {
namespace channels {
namespace cnmaestro {

namespace {
// TODO: dynamically read in from local file.
std::string accessTokenLead = "Authorization: Bearer ";
const char* grantType = "grant_type=client_credentials";
const char* ipHead = "https://";
std::string curlPath = "/usr/bin/curl";
std::string accessTokenPath = "api/v1/access/token";
std::string devicesPath = "api/v1/devices";
std::chrono::seconds curlTimeout(5);
constexpr int operationTimeout = 28;
} // namespace

Channel::Channel(
    const std::string& ipAddr_,
    const std::string& clientId_,
    const std::string& clientSecret_)
    : ipAddr(ipAddr_), clientId(clientId_), clientSecret(clientSecret_) {
  std::string address = ipHead + ipAddr + "/";
  accessTokenCommandPath = address + accessTokenPath;
  allDevicesCommandPath = address + devicesPath;

  clientIdString = std::string("client_id=") + clientId;
  clientSecretString = std::string("client_secret=") + clientSecret;

  connected = false;
}

folly::dynamic Channel::setupChannel() {
  // TODO: Convert to async
  // auto promise = std::make_shared<folly::Promise<folly::dynamic>>();
  // std::thread([=]{
  //   promise->setWith(std::bind(connectChannel));
  // }).detach();
  //
  // return promise->getFuture();
  return connectChannel();
}

folly::dynamic Channel::connectChannel() {
  std::vector<std::string> accessTokenCommand = {
      curlPath,
      accessTokenCommandPath,
      "-X",
      "POST",
      "-k",
      "-d",
      grantType,
      "-d",
      clientIdString,
      "-d",
      clientSecretString,
      "-m",
      folly::to<std::string>(curlTimeout.count())};

  auto returnedData = makeCall(accessTokenCommand);

  folly::dynamic output;

  if (not returnedData.isNull()) {
    output = folly::parseJson(returnedData.asString());
    accessToken = output["access_token"];
    accessTokenPiece = accessTokenLead + accessToken;
    connected = true;

  } else {
    connected = false;
    output = returnedData;
  }
  return output;
}

folly::dynamic Channel::getDeviceInfo(const std::string& clientMac) {
  // TODO: Convert to async
  folly::dynamic output;
  if (connected) {
    std::string curlAccessTokenPath = allDevicesCommandPath + "/" + clientMac;
    std::vector<std::string> statusCommand = {
        curlPath,
        curlAccessTokenPath,
        "-X",
        "GET",
        "-k",
        "-m",
        folly::to<std::string>(curlTimeout.count()),
        "-H",
        accessTokenPiece.asString()};
    output = makeCall(statusCommand);
  } else {
    LOG(ERROR) << "Device not connected. Retrying connection.";
    output = connectChannel();
  }
  return output;
}

folly::dynamic Channel::makeCall(std::vector<std::string>& cmd) {
  folly::Subprocess proc(
      std::vector<std::string>{cmd}, folly::Subprocess::Options().pipeStdout());
  auto p = proc.communicate();
  folly::dynamic output;
  auto exitStatus = proc.wait().exitStatus();
  switch (exitStatus) {
    case 0: {
      output = p.first;
      break;
    }
    case operationTimeout: {
      LOG(ERROR) << "Timed out while trying to connect to cnmaestro channel: "
                 << ipAddr;
      connected = false;
      break;
    }
    default: {
      LOG(ERROR) << "Unexpectedely errored with exit status " << exitStatus
                 << "whie trying to connect to cnmaestro channel " << ipAddr;
      connected = false;
      break;
    }
  }
  return output;
}

void Channel::curlPut(
    folly::dynamic& updateInfo,
    const std::string& clientMac) {
  // updateInfo is a limited json that only contains fields to be updated [no
  // strict structure]
  std::string updateString = toJson(updateInfo);
  std::string putAccessStr = allDevicesCommandPath + "/" + clientMac;
  std::vector<std::string> putCommand = {
      curlPath,
      putAccessStr,
      "-X",
      "PUT",
      "-k",
      "-m",
      folly::to<std::string>(curlTimeout.count()),
      "--header",
      "Content-Type: application/json",
      "-d",
      updateString,
      "-H",
      accessTokenPiece.asString()};
  folly::dynamic output = makeCall(putCommand);

  return;
}

void Channel::updateDevice(
    folly::dynamic& updateInfo,
    const std::string& clientMac) {
  if (connected) {
    curlPut(updateInfo, clientMac);
  } else {
    LOG(ERROR) << "No connection to cnmaestro; attempting to reconnect.";
    auto output = connectChannel();
  }
}

// folly::Future<folly::dynamic>
// Channel::asyncMsg(std::function<std::shared_ptr<httplib::Response>()> send){
//   auto requestId = requestGuid++;
//   auto result = outstandingRequests.emplace(
//       std::piecewise_construct,
//       std::forward_as_tuple(requestId),
//       std::forward_as_tuple(folly::Promise<Response>{}));
//   if (result.second) {
//     std::thread t([this, send, requestId, ]()) {
//       ErrorHandler::executeWithCatch([this, &send, requestId, ]()) {
//         auto res = send();
//         auto pResult = outstandingRequests.find(requestId);
//         LOG(ERROR) << "Some error happened. More error handling needed.";
//       }
//     }
//   }
// }

} // namespace cnmaestro
} // namespace channels
} // namespace devmand
