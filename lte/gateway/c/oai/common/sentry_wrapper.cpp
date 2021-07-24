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

#include "sentry_wrapper.h"

#if SENTRY_ENABLED
#include <experimental/optional>
#include <yaml-cpp/yaml.h>  // IWYU pragma: keep

#include <cstdlib>
#include <fstream>
#include <sstream>
#include <string>

#include "sentry.h"
#include "includes/ServiceConfigLoader.h"

#define COMMIT_HASH_ENV "COMMIT_HASH"
#define CONTROL_PROXY_SERVICE_NAME "control_proxy"
#define SENTRY_NATIVE_URL "sentry_url_native"
#define SHOULD_UPLOAD_MME_LOG "sentry_upload_mme_log"
#define MME_LOG_PATH "/var/log/mme.log"
#define SNOWFLAKE_PATH "/etc/snowflake"
#define HWID "hwid"
#define SERVICE_NAME "service_name"

using std::experimental::optional;

bool should_upload_mme_log(YAML::Node control_proxy_config) {
  if (control_proxy_config[SHOULD_UPLOAD_MME_LOG].IsDefined()) {
    return control_proxy_config[SHOULD_UPLOAD_MME_LOG].as<bool>();
  }
  return false;
}

optional<std::string> get_sentry_url(YAML::Node control_proxy_config) {
  std::string sentry_url;
  if (control_proxy_config[SENTRY_NATIVE_URL].IsDefined()) {
    const std::string sentry_dns =
        control_proxy_config[SENTRY_NATIVE_URL].as<std::string>();
    if (sentry_dns.size()) {
      return sentry_dns;
    }
  }
  return {};
}

std::string get_snowflake() {
  std::ifstream ifs(SNOWFLAKE_PATH, std::ifstream::in);
  std::stringstream buffer;
  buffer << ifs.rdbuf();
  return buffer.str();
}

void initialize_sentry() {
  auto control_proxy_config = magma::ServiceConfigLoader{}.load_service_config(
      CONTROL_PROXY_SERVICE_NAME);
  auto op_sentry_url = get_sentry_url(control_proxy_config);
  if (op_sentry_url) {
    sentry_options_t* options = sentry_options_new();
    sentry_options_set_dsn(options, op_sentry_url->c_str());
    if (const char* commit_hash_p = std::getenv(COMMIT_HASH_ENV)) {
      sentry_options_set_release(options, commit_hash_p);
    }
    if (should_upload_mme_log(control_proxy_config)) {
      sentry_options_add_attachment(options, MME_LOG_PATH);
    }

    sentry_init(options);

    sentry_set_tag(SERVICE_NAME, "MME");
    sentry_set_tag(HWID, get_snowflake().c_str());
  }
}

void shutdown_sentry(void) {
  sentry_shutdown();
}
#else
void initialize_sentry(void) {}
void shutdown_sentry(void) {}
#endif
