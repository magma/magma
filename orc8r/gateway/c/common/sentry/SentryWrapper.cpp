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

#include "orc8r/gateway/c/common/sentry/includes/SentryWrapper.h"

#include <bits/local_lim.h>
#include <unistd.h>
#include <string>

#if SENTRY_ENABLED
#include <yaml-cpp/yaml.h>
#include <experimental/optional>
#include <cstdlib>
#include <fstream>
#include <cstring>

#include "sentry.h"
#include "orc8r/gateway/c/common/config/includes/ServiceConfigLoader.h"

#define COMMIT_HASH_ENV "COMMIT_HASH"
#define CONTROL_PROXY_SERVICE_NAME "control_proxy"
#define SENTRY_NATIVE_URL "sentry_url_native"
#define SENTRY_SAMPLE_RATE "sentry_sample_rate"
#define SHOULD_UPLOAD_MME_LOG "sentry_upload_mme_log"
#define MME_LOG_PATH "/var/log/mme.log"
#define SNOWFLAKE_PATH "/etc/snowflake"
#define HWID "hwid"
#define HOSTNAME "hostname"
#define SERVICE_NAME "service_name"
#define CLOUD_ADDRESS "cloud_address"
#define ORC8R_CLOUD_ADDRESS "orc8r_cloud_address"
#define DEFAULT_SAMPLE_RATE 0.5f

using std::experimental::optional;

bool should_upload_mme_log(bool sentry_upload_mme_log,
                           YAML::Node control_proxy_config) {
  if (control_proxy_config[SHOULD_UPLOAD_MME_LOG].IsDefined()) {
    return control_proxy_config[SHOULD_UPLOAD_MME_LOG].as<bool>();
  }
  return sentry_upload_mme_log;
}

optional<std::string> get_sentry_url(const char* sentry_url_native,
                                     YAML::Node control_proxy_config) {
  if (control_proxy_config[SENTRY_NATIVE_URL].IsDefined()) {
    const std::string dns_override =
        control_proxy_config[SENTRY_NATIVE_URL].as<std::string>();
    if (!dns_override.empty()) {
      return dns_override;
    }
  }
  const std::string sentry_url(sentry_url_native);
  if (!sentry_url.empty()) {
    return sentry_url;
  }
  return {};
}

optional<std::string> get_cloud_address(YAML::Node control_proxy_config) {
  if (control_proxy_config[SENTRY_NATIVE_URL].IsDefined()) {
    const std::string cloud_address =
        control_proxy_config[CLOUD_ADDRESS].as<std::string>();
    if (!cloud_address.empty()) {
      return cloud_address;
    }
  }
  return {};
}

float get_sentry_sample_rate(float sentry_sample_rate,
                             YAML::Node control_proxy_config) {
  if (control_proxy_config[SENTRY_SAMPLE_RATE].IsDefined()) {
    const auto sample_rate_override =
        control_proxy_config[SENTRY_SAMPLE_RATE].as<float>();
    if (sample_rate_override > 0) {
      return sample_rate_override;
    }
  }
  if (sentry_sample_rate > 0) {
    return sentry_sample_rate;
  }
  return DEFAULT_SAMPLE_RATE;
}

std::string get_snowflake() {
  std::ifstream ifs(SNOWFLAKE_PATH, std::ifstream::in);
  std::stringstream buffer;
  buffer << ifs.rdbuf();
  return buffer.str();
}

void initialize_sentry(const char* service_tag,
                       const sentry_config_t* sentry_config) {
  auto control_proxy_config = magma::ServiceConfigLoader{}.load_service_config(
      CONTROL_PROXY_SERVICE_NAME);
  auto op_sentry_url =
      get_sentry_url(sentry_config->url_native, control_proxy_config);
  if (!op_sentry_url) {
    return;
  }

  sentry_options_t* options = sentry_options_new();
  sentry_options_set_dsn(options, op_sentry_url->c_str());
  sentry_options_set_sample_rate(
      options,
      get_sentry_sample_rate(sentry_config->sample_rate, control_proxy_config));
  if (const char* commit_hash_p = std::getenv(COMMIT_HASH_ENV)) {
    sentry_options_set_release(options, commit_hash_p);
  }
  if (strncmp(service_tag, SENTRY_TAG_MME, SENTRY_TAG_LEN) == 0 &&
      should_upload_mme_log(sentry_config->upload_mme_log,
                            control_proxy_config)) {
    sentry_options_add_attachment(options, MME_LOG_PATH);
  }
  if (sentry_config->add_debug_logging) {
    sentry_options_set_debug(options, 1);
  }

  char db_path[SENTRY_DB_PREFIX_LEN + SENTRY_TAG_LEN];
  strcpy(db_path, SENTRY_DB_PREFIX);
  strcat(db_path, service_tag);
  sentry_options_set_database_path(options, db_path);

  sentry_init(options);
  char node_name[HOST_NAME_MAX];
  if (gethostname(node_name, HOST_NAME_MAX) == 0) {
    sentry_set_tag(HOSTNAME, node_name);
  }
  if (auto cloud_address = get_cloud_address(control_proxy_config)) {
    sentry_set_tag(ORC8R_CLOUD_ADDRESS, cloud_address->c_str());
  }
  sentry_set_tag(SERVICE_NAME, service_tag);
  sentry_set_tag(HWID, get_snowflake().c_str());
}

void shutdown_sentry(void) { sentry_shutdown(); }

void set_sentry_transaction(const char* name) { sentry_set_transaction(name); }

void sentry_log_error(const char* message) {
  sentry_value_t event =
      sentry_value_new_message_event(SENTRY_LEVEL_ERROR, NULL, message);
  sentry_event_value_add_stacktrace(event, NULL, 0);
  sentry_capture_event(event);
}

#else

void initialize_sentry(__attribute__((unused)) const char* service_tag,
                       __attribute__((unused))
                       const sentry_config_t* sentry_config) {}

void shutdown_sentry(void) {}

void set_sentry_transaction(__attribute__((unused)) const char* name) {}

void sentry_log_error(__attribute__((unused)) const char* message);

#endif
