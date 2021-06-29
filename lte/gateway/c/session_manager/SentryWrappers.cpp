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

#include "SentryWrappers.h"

#if SENTRY_ENABLED
#include <experimental/optional>

#include <cstdlib>
#include <fstream>
#include <sstream>
#include <string>
#include <yaml-cpp/yaml.h>  // IWYU pragma: keep

#include "magma_logging_init.h"
#include "sentry.h"
#include "includes/ServiceConfigLoader.h"

#define COMMIT_HASH_ENV "COMMIT_HASH"
#define CONTROL_PROXY_SERVICE_NAME "control_proxy"
#define SENTRY_NATIVE_URL "sentry_url_native"
#define SENTRY_SAMPLE_RATE "sentry_sample_rate"
#define SNOWFLAKE_PATH "/etc/snowflake"
#define HWID "hwid"
#define SERVICE_NAME "service_name"

using std::experimental::optional;

// TODO(@themarwhal) pull common sentry functions into lib common

optional<std::string> get_sentry_url(
    const magma::mconfig::SentryConfig& sentry_config,
    YAML::Node control_proxy_config) {
  if (control_proxy_config[SENTRY_NATIVE_URL].IsDefined()) {
    const std::string dns_override =
        control_proxy_config[SENTRY_NATIVE_URL].as<std::string>();
    if (dns_override.size()) {
      return dns_override;
    }
  }
  const std::string sentry_url = sentry_config.url_native();
  if (sentry_url.size()) {
    return sentry_url;
  }
  return {};
}

float get_sentry_sample_rate(
    const magma::mconfig::SentryConfig& sentry_config,
    YAML::Node control_proxy_config) {
  if (control_proxy_config[SENTRY_SAMPLE_RATE].IsDefined()) {
    const float sample_rate_override =
        control_proxy_config[SENTRY_SAMPLE_RATE].as<float>();
    if (sample_rate_override) {
      return sample_rate_override;
    }
  }
  const float sample_rate = sentry_config.sample_rate();
  if (sample_rate) {
    return sample_rate;
  }
  return 1.0f;
}

std::string get_snowflake() {
  std::ifstream ifs(SNOWFLAKE_PATH, std::ifstream::in);
  std::stringstream buffer;
  buffer << ifs.rdbuf();
  return buffer.str();
}

void initialize_sentry(const magma::mconfig::SentryConfig& sentry_config) {
  auto control_proxy_config = magma::ServiceConfigLoader{}.load_service_config(
      CONTROL_PROXY_SERVICE_NAME);
  auto op_sentry_url = get_sentry_url(sentry_config, control_proxy_config);
  if (op_sentry_url) {
    MLOG(MINFO) << "Starting SessionD with Sentry!";
    sentry_options_t* options = sentry_options_new();
    sentry_options_set_dsn(options, op_sentry_url->c_str());
    sentry_options_set_sample_rate(
        options, get_sentry_sample_rate(sentry_config, control_proxy_config));
    if (const char* commit_hash_p = std::getenv(COMMIT_HASH_ENV)) {
      sentry_options_set_release(options, commit_hash_p);
    }

    sentry_init(options);
    sentry_set_tag(SERVICE_NAME, "SessionD");
    sentry_set_tag(HWID, get_snowflake().c_str());
  }
}

void shutdown_sentry() {
  sentry_shutdown();
}

void set_sentry_transaction(const std::string& name) {
  sentry_set_transaction(name.c_str());
}

#else
void initialize_sentry(__attribute__((unused))
                       const magma::mconfig::SentryConfig& sentry_config) {}
void shutdown_sentry() {}
void set_sentry_transaction(const std::string& name) {}
#endif
