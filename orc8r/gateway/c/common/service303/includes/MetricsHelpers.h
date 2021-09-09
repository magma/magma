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

#include <stdio.h>  // for size_t

#ifdef __cplusplus
extern "C" {
#endif

/**
 * Remove the counter metric that matches name+labels given
 * @param name
 * @param n_labels number of labels
 * @param ... label args (name, value)
 */
void remove_counter(const char* name, size_t n_labels, ...);

/**
 * Increments value for Counter metric
 * @param name
 * @param increment value to increment
 * @param n_labels number of labels
 * @param ... label args (name, value)
 */
void increment_counter(const char* name, double increment, size_t n_labels,
                       ...);

/**
 * Remove the gauge metric that matches name+labels given
 * @param name
 * @param n_labels number of labels
 * @param ... label args (name, value)
 */
void remove_gauge(const char* name, size_t n_labels, ...);

/**
 * Increments value for Gauge metric
 * @param name
 * @param increment value to increment
 * @param n_labels number of labels
 * @param ... label args (name, value)
 */
void increment_gauge(const char* name, double increment, size_t n_labels, ...);

/**
 * Decrements value for Gauge metric
 * @param name
 * @param increment value to increment
 * @param n_labels number of labels
 * @param ... label args (name, value)
 */
void decrement_gauge(const char* name, double decrement, size_t n_labels, ...);

/**
 * Sets specific value for Gauge metric
 * @param name
 * @param value to set
 * @param n_labels number of labels
 * @param ... label args (name, value)
 */
void set_gauge(const char* name, double value, size_t n_labels, ...);

/**
 * Returns value for Gauge metric
 * @param name
 * @param n_labels number of labels
 * @param ... label args (name, value)
 */
double get_gauge(const char* name, size_t n_labels, ...);

/**
 * Updates value of Histogram metric
 * @param name
 * @param value to observe
 * @param n_labels number of labels
 * @param ... label args (name, value)
 */
void observe_histogram(const char* name, double observation, size_t n_labels,
                       ...);

#ifdef __cplusplus
}
#endif
