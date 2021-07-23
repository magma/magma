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

#include "mme_config.h"

#ifdef __cplusplus
extern "C" {
#endif

/**
 * @brief Initialize sentry if SENTRY_ENABLED flag is set and project slug is
 * configured in control_proxy.yml
 */
void initialize_sentry(const sentry_config_t* sentry_config);

/**
 * @brief Shutdown sentry if SENTRY_ENABLED flag is set
 */
void shutdown_sentry(void);

#ifdef __cplusplus
}
#endif
