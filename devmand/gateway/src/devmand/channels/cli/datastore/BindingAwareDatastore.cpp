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

#include <devmand/channels/cli/datastore/BindingAwareDatastore.h>
#include <devmand/channels/cli/datastore/Datastore.h>

namespace devmand::channels::cli::datastore {

using devmand::channels::cli::datastore::Datastore;
BindingAwareDatastore::BindingAwareDatastore(
    const shared_ptr<Datastore> _datastore,
    const shared_ptr<BindingCodec> _codec)
    : datastore(_datastore), codec(_codec) {}

unique_ptr<BindingAwareDatastoreTransaction>
BindingAwareDatastore::newBindingTx() {
  std::unique_lock<std::mutex> lock(_mutex);
  return std::make_unique<BindingAwareDatastoreTransaction>(
      datastore->newTx(), codec);
}

} // namespace devmand::channels::cli::datastore
