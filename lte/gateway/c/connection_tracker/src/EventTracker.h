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

#include "PacketGenerator.h"

namespace magma {
namespace lte {

class EventTracker {
 public:
  EventTracker(std::shared_ptr<PacketGenerator> pkt_gen, int zone);

  int init_conntrack_event_loop();

  std::shared_ptr<PacketGenerator> pkt_gen_;
  int zone_;
};

}  // namespace lte
}  // namespace magma
