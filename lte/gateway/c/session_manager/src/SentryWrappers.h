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

// For Sentry documentation, see https://docs.sentry.io/platforms/native/

/**
 * @brief Initialize sentry if SENTRY_ENABLED flag is set and project slug is
 * configured in control_proxy.yml
 */
void initialize_sentry();

/**
 * @brief Shutdown sentry if SENTRY_ENABLED flag is set
 */
void shutdown_sentry();

/**
 * @brief Set the sentry transaction object
 *
 * @param name
 */
void set_sentry_transaction(const std::string& name);
