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

#pragma once

namespace devmand {
namespace channels {

/*
 * This is the common interface all channels must implement. A channel is a way
 * to communicate with a device. The channel's engine maintains any state needed
 * outside of an individual connection.
 */
class Channel {
 public:
  Channel() = default;
  virtual ~Channel() = default;
  Channel(const Channel&) = delete;
  Channel& operator=(const Channel&) = delete;
  Channel(Channel&&) = delete;
  Channel& operator=(Channel&&) = delete;
};

} // namespace channels
} // namespace devmand
