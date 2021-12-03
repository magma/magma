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

#include <string>
#include <openssl/ssl.h>

namespace magma {
namespace lte {

class ProxyConnector {
 public:
  virtual int send_data(void* data, uint32_t size) = 0;
  virtual int setup_proxy_socket()                 = 0;
  virtual void cleanup()                           = 0;
  virtual ~ProxyConnector()                        = default;
};

class ProxyConnectorImpl : public ProxyConnector {
 public:
  ProxyConnectorImpl(
      const std::string& proxy_addr, const int port,
      const std::string& cert_file, const std::string& key_file);

  /**
   * setup_proxy_socket instantiate ssl library and opens a tls connection
   * @return return positif integer if it successeds.
   */
  int setup_proxy_socket();

  /**
   * export_record exports the x3 record over tls to a remote server.
   * @param data - x3 record packet
   * @param size - x3 record length
   * @return return positif integer if sending data successeds.
   */
  int send_data(void* data, uint32_t size);

  /**
   * cleanup cleans up all allocated ressources in proxy connector
   * @return void
   */
  void cleanup();

  /**
   * performs cleanup if cleanup not explicitly called
   */
  ~ProxyConnectorImpl() override;

 private:
  const std::string& proxy_addr_;
  const int proxy_port_;
  const std::string& cert_file_;
  const std::string& key_file_;
  SSL* ssl_;
  SSL_CTX* ctx_;
  int proxy_;

  int open_connection();
  int load_certificates(SSL_CTX* ctx);
  SSL_CTX* init_ctx();
};

}  // namespace lte
}  // namespace magma
