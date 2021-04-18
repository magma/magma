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

class ProxyConnector {
 public:
  virtual int send_data(void* data, uint32_t size) = 0;
  virtual int setup_proxy_socket()                 = 0;
};

class ProxyConnectorImpl : public ProxyConnector {
 public:
  ProxyConnectorImpl(
      const std::string& proxy_addr, const int port,
      const std::string& cert_file, const std::string& key_file);

  int setup_proxy_socket();
  int send_data(void* data, uint32_t size);
  void cleanup();

 private:
  int open_connection();
  int load_certificates(SSL_CTX* ctx);
  SSL_CTX* init_ctx();

  const std::string& proxy_addr_;
  const int proxy_port_;
  const std::string& cert_file_;
  const std::string& key_file_;
  SSL* ssl_;
  SSL_CTX* ctx_;
  int proxy_;
};

}  // namespace magma
