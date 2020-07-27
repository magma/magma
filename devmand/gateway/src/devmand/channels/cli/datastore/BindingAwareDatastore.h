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
#include <devmand/channels/cli/datastore/BindingAwareDatastoreTransaction.h>
#include <devmand/channels/cli/datastore/Datastore.h>
#include <devmand/devices/cli/schema/BindingContext.h>

namespace devmand::channels::cli::datastore {

using devmand::channels::cli::datastore::BindingAwareDatastoreTransaction;
using devmand::channels::cli::datastore::Datastore;
using devmand::channels::cli::datastore::DatastoreTransaction;
using devmand::devices::cli::BindingCodec;

class BindingAwareDatastore {
 private:
  shared_ptr<Datastore> datastore;
  shared_ptr<BindingCodec> codec;
  std::mutex _mutex;

 public:
  unique_ptr<BindingAwareDatastoreTransaction>
  newBindingTx(); // operations on transaction are NOT thread-safe
  BindingAwareDatastore(
      const shared_ptr<Datastore>,
      const shared_ptr<BindingCodec> codec);
};

} // namespace devmand::channels::cli::datastore
