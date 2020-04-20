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

#include <boost/algorithm/string.hpp>
#include <devmand/channels/cli/datastore/BindingAwareDatastoreTransaction.h>

namespace devmand::channels::cli::datastore {

DiffResult BindingAwareDatastoreTransaction::diff(
    vector<DiffPath> registeredPaths) {
  return datastoreTransaction->diff(registeredPaths);
}

void BindingAwareDatastoreTransaction::delete_(Path path) {
  datastoreTransaction->delete_(path);
}

void BindingAwareDatastoreTransaction::overwrite(
    Path path,
    shared_ptr<Entity> entity) {
  datastoreTransaction->overwrite(path, codec->toDom(path, *entity));
}

void BindingAwareDatastoreTransaction::merge(
    Path path,
    shared_ptr<Entity> entity) {
  datastoreTransaction->merge(path, codec->toDom(path, *entity));
}

void BindingAwareDatastoreTransaction::commit() {
  datastoreTransaction->commit();
}

BindingAwareDatastoreTransaction::BindingAwareDatastoreTransaction(
    unique_ptr<DatastoreTransaction>&& _transaction,
    shared_ptr<BindingCodec> _codec)
    : datastoreTransaction(
          std::forward<unique_ptr<DatastoreTransaction>>(_transaction)),
      codec(_codec) {}

void BindingAwareDatastoreTransaction::isValid() {
  return datastoreTransaction->isValid();
}

void BindingAwareDatastoreTransaction::abort() {
  datastoreTransaction->abort();
}

void BindingAwareDatastoreTransaction::print() {
  datastoreTransaction->print();
}

} // namespace devmand::channels::cli::datastore
