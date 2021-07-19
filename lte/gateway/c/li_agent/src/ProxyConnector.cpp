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

#include <netinet/in.h>
#include <arpa/inet.h>
#include <openssl/ssl.h>
#include <openssl/err.h>

#include "ProxyConnector.h"
#include "magma_logging.h"

namespace magma {
namespace lte {

ProxyConnectorImpl::ProxyConnectorImpl(
    const std::string& proxy_addr, const int proxy_port,
    const std::string& cert_file, const std::string& key_file)
    : proxy_addr_(proxy_addr),
      proxy_port_(proxy_port),
      cert_file_(cert_file),
      key_file_(key_file) {}

int ProxyConnectorImpl::setup_proxy_socket() {
  SSL_library_init();

  ctx_ = init_ctx();
  if (ctx_ == NULL) {
    return -1;
  }
  if (load_certificates(ctx_) < 0) {
    return -1;
  }
  proxy_ = open_connection();
  if (proxy_ < 0) {
    return -1;
  }
  ssl_ = SSL_new(ctx_);
  if (ssl_ == NULL) {
    return -1;
  }
  SSL_set_options(ssl_, SSL_OP_DONT_INSERT_EMPTY_FRAGMENTS);
  SSL_set_fd(ssl_, proxy_);
  if (SSL_connect(ssl_) == -1) {
    ERR_print_errors_fp(stderr);
    return -1;
  }

  return 0;
}

int ProxyConnectorImpl::load_certificates(SSL_CTX* ctx) {
  if (SSL_CTX_use_certificate_file(ctx, cert_file_.c_str(), SSL_FILETYPE_PEM) <=
      0) {
    ERR_print_errors_fp(stderr);
    return -1;
  }
  if (SSL_CTX_use_PrivateKey_file(ctx, key_file_.c_str(), SSL_FILETYPE_PEM) <=
      0) {
    ERR_print_errors_fp(stderr);
    return -1;
  }
  if (!SSL_CTX_check_private_key(ctx)) {
    MLOG(MERROR) << "Private key does not match the public certificate";
    return -1;
  }
  return 0;
}

SSL_CTX* ProxyConnectorImpl::init_ctx(void) {
  SSL_CTX* ctx;

  OpenSSL_add_all_algorithms(); /* Load cryptos, et.al. */
  SSL_load_error_strings();     /* Bring in and register error messages */
  ctx = SSL_CTX_new(TLS_client_method()); /* Create new context */
  if (ctx == NULL) {
    ERR_print_errors_fp(stderr);
    return NULL;
  }
  return ctx;
}

int ProxyConnectorImpl::open_connection() {
  int sd;
  struct sockaddr_in serv_addr;

  sd                        = socket(AF_INET, SOCK_STREAM, 0);
  serv_addr.sin_family      = AF_INET;
  serv_addr.sin_addr.s_addr = INADDR_ANY;
  serv_addr.sin_port        = htons(proxy_port_);

  if (inet_pton(AF_INET, proxy_addr_.c_str(), &serv_addr.sin_addr) <= 0) {
    MLOG(MERROR) << "Invalid address/ Address not supported";
    return -1;
  }
  if (connect(sd, (struct sockaddr*) &serv_addr, sizeof(struct sockaddr_in)) !=
      0) {
    MLOG(MERROR) << "Can't connect to the proxy, exiting";
    close(sd);
    return -1;
  }
  return sd;
}

int ProxyConnectorImpl::send_data(void* data, uint32_t size) {
  return SSL_write(ssl_, data, size);
}

void ProxyConnectorImpl::cleanup() {
  SSL_free(ssl_);
  close(proxy_);
  SSL_CTX_free(ctx_);
}

}  // namespace lte
}  // namespace magma
