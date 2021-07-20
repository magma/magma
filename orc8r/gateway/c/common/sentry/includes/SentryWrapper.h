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

#ifdef __cplusplus
extern "C" {
#endif

#define MAX_URL_LENGTH 255

/**
 * @brief Struct to contain Sentry configuration relevant for C/C++ services
 */
typedef struct sentry_config {
  float sample_rate;
  bool upload_mme_log;
  char url_native[MAX_URL_LENGTH];
} sentry_config_t;

/**
 * @brief Initialize sentry if SENTRY_ENABLED flag is set and project slug is
 * configured in control_proxy.yml
 */
void initialize_sentry(const sentry_config_t* sentry_config);

/**
 * @brief Shutdown sentry if SENTRY_ENABLED flag is set
 */
void shutdown_sentry(void);

/**
 * @brief Set the sentry transaction object
 *
 * @param name
 */
void set_sentry_transaction(const char* name);

#ifdef __cplusplus
}
#endif
