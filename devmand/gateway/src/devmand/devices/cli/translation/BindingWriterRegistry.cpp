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

#include <devmand/devices/cli/translation/BindingWriterRegistry.h>

namespace devmand {
namespace devices {
namespace cli {

using namespace folly;
using namespace std;
using namespace ydk;

BindingWriterRegistryBuilder::BindingWriterRegistryBuilder(
    WriterRegistryBuilder& _domBuilder,
    BindingContext& _context)
    : domBuilder(_domBuilder), context(_context) {}

} // namespace cli
} // namespace devices
} // namespace devmand
