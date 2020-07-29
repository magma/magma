/*
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

// Package serde contains the definition of a SERializer-DEserializer concept.
// This package also includes a global registry of serdes for applications to
// delegate implementation-agnostic serialization and deserialization to.
// Serdes are one of the primary plugin interfaces exposed by orc8r to extend
// services with domain-specific data models and logic.
package serde
