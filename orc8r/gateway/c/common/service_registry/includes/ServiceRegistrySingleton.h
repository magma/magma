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
#pragma once

#include <yaml-cpp/yaml.h>                 // IWYU pragma: keep
#include <memory>                          // for shared_ptr, unique_ptr
#include <string>                          // for string
#include "includes/ServiceConfigLoader.h"  // for ServiceConfigLoader
namespace grpc {
class Channel;
}
namespace grpc {
class ChannelCredentials;
}

using grpc::Channel;
using grpc::ChannelCredentials;

namespace magma {

typedef struct {
  std::string ip;
  std::string port;
} ip_port_pair_t;

typedef struct {
  std::string ip;
  std::string port;
  std::string authority;
  std::shared_ptr<ChannelCredentials> creds = nullptr;
} create_grpc_channel_args_t;
/*
 * ServiceRegistrySingleton is a singleton used to get a grpc channel
 * given a service the client wants to connect to.
 */
class ServiceRegistrySingleton {
 public:
  static constexpr char* CLOUD = (char*) "cloud";
  static constexpr char* LOCAL = (char*) "local";

 public:
  static ServiceRegistrySingleton* Instance();

  static void flush();  // destroy instance

  /*
   * Gets the grpc args to the specified service based on service name
   * and destination.
   * @param service: service name to where a connection should be open.
   * @param destination: str indicating if a connection to the cloud service
   * or local service. Can be either "local" or "cloud".
   * @return create_grpc_channel_args_t, which contains a str ip,
   * a str port, and a str authority to be passed to CreateGrpcChannel.
   */
  const create_grpc_channel_args_t GetCreateGrpcChannelArgs(
      const std::string& service, const std::string& destination);
  /*
   * Gets the ip:port of a service if the service is local.
   * @param service: service name to where a connection should be open.
   * @return ip_port_pair_t, which contains a str ip, and a str port.
   */
  ip_port_pair_t GetServiceAddr(const std::string& service);

  /*
   * Gets the ip:port of a service in string form like "127.0.0.1:8888"
   * @param service: service name to where a connection should be open.
   * @return string of the ip:port pairing
   */
  std::string GetServiceAddrString(const std::string& service);

  /*
   * Gets SSL credentials for requests directly to the cloud
   */
  std::shared_ptr<ChannelCredentials> GetSslCredentials();

  /*
   * Gets a grpc connection to the specified service based on service name
   * and destination.
   * @param service: service name to where a connection should be open.
   * @param destination: str indicating if a connection to the cloud service
   * or local service. Can be either "local" or "cloud".
   * @return grpc::Channel to the given service. If a connection to cloud
   * service is requested, a connection to the control_proxy will be returned.
   */
  const std::shared_ptr<Channel> GetGrpcChannel(
      const std::string& service, const std::string& destination);

  /*
   * Returns a grpc connection to the bootstrapper service
   */
  const std::shared_ptr<Channel> GetBootstrapperGrpcChannel();

 private:
  ServiceRegistrySingleton();  // Prevent construction
  // Prevent construction by copying
  ServiceRegistrySingleton(const ServiceRegistrySingleton&);
  // Prevent assignment
  ServiceRegistrySingleton& operator=(const ServiceRegistrySingleton&);
  YAML::Node GetProxyConfig();
  YAML::Node GetRegistry();
  // load a cert file as a string into the given buffer
  std::string LoadCertFile(const std::string& file);
  const std::shared_ptr<Channel> CreateGrpcChannel(
      const std::string& ip, const std::string& port,
      const std::string& authority, std::shared_ptr<ChannelCredentials> creds);

 private:
  ServiceConfigLoader service_config_loader_;
  std::unique_ptr<YAML::Node> proxy_config_;
  std::unique_ptr<YAML::Node> registry_;
  static ServiceRegistrySingleton* instance_;
};

}  // namespace magma
