/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

#pragma once

#include <string>
#include <unordered_map>

#include <folly/dynamic.h>
#include <folly/futures/Future.h>
#include <folly/json.h>

#include <devmand/channels/Channel.h>
// #include <devmand/channels/cnmaestro/Response.h>

namespace devmand {
namespace channels {
namespace cnmaestro {

// using OutstandingRequests = std::map<unsigned int, folly::Promise<Response>>;

class Channel final : public channels::Channel {
 public:
  Channel(
      const std::string& ipAddr_,
      const std::string& clientId_,
      const std::string& clientSecret_);
  Channel() = delete;
  ~Channel() override = default;
  Channel(const Channel&) = delete;
  Channel& operator=(const Channel&) = delete;
  Channel(Channel&&) = delete;
  Channel& operator=(Channel&&) = delete;

 public:
  folly::dynamic setupChannel();
  folly::dynamic getDeviceInfo(const std::string& clientMac);
  void updateDevice(folly::dynamic& updateInfo, const std::string& clientMac);

 private:
  void curlPut(folly::dynamic& updateInfo, const std::string& clientMac);
  folly::dynamic connectChannel();
  folly::dynamic makeCall(std::vector<std::string>& cmd);

 private:
  folly::dynamic accessToken;
  std::string ipAddr;
  std::string clientId;
  std::string clientSecret;
  std::string accessTokenCommandPath;
  std::string allDevicesCommandPath;
  std::string clientIdString;
  std::string clientSecretString;
  folly::dynamic accessTokenPiece;
  bool connected;

  // unsigned int requestGuid{0};
  // OutstandingRequests outstandingRequests;
};

} // namespace cnmaestro
} // namespace channels
} // namespace devmand
