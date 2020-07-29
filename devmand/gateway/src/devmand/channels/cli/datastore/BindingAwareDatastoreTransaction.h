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

#include <devmand/channels/cli/datastore/DatastoreState.h>
#include <devmand/channels/cli/datastore/DatastoreTransaction.h>
#include <devmand/devices/cli/schema/BindingContext.h>
#include <devmand/devices/cli/schema/Path.h>

namespace devmand::channels::cli::datastore {
using devmand::devices::cli::BindingCodec;
using devmand::devices::cli::Path;
using std::unique_ptr;

class BindingAwareDatastoreTransaction {
 private:
  unique_ptr<DatastoreTransaction> datastoreTransaction;
  shared_ptr<BindingCodec> codec;

 public:
  BindingAwareDatastoreTransaction(
      unique_ptr<DatastoreTransaction>&& _transaction,
      shared_ptr<BindingCodec> _codec);

 public:
  template <typename T>
  shared_ptr<T> read(Path path) {
    const shared_ptr<T>& ydkData = make_shared<T>();
    const dynamic& data = datastoreTransaction->read(path);
    return std::static_pointer_cast<T>(codec->decode(toJson(data), ydkData));
  }
  DiffResult diff(vector<DiffPath> registeredPaths);
  void delete_(Path path);
  void overwrite(Path path, shared_ptr<Entity> entity);
  void merge(Path path, shared_ptr<Entity> entity);
  void isValid();
  void print();
  void commit();
  void abort();
};

} // namespace devmand::channels::cli::datastore
