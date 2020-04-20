/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// File keys.go contains the config keynames in the state service's YAML config file.

package config

const (
	// EnableAutomaticReindexing is a parameter name in the state service config.
	// When value is true, state service handles automatically reindex state indexers.
	// When value is false, reindexing must be handled by the provided CLI.
	EnableAutomaticReindexing = "enable_automatic_reindexing"
)
