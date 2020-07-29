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

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <devmand/devices/cli/schema/Model.h>
#include <devmand/devices/cli/schema/SchemaContext.h>
#include <folly/dynamic.h>
#include <ydk/codec_provider.hpp>
#include <ydk/codec_service.hpp>
#include <ydk/json_subtree_codec.hpp>
#include <ydk/path_api.hpp>
#include <mutex>
#include <sstream>

namespace devmand {
namespace devices {
namespace cli {

using namespace std;
using namespace folly;
using namespace ydk;
using namespace ydk::path;
using namespace devmand::devices::cli;

class BindingCodec {
 public:
  explicit BindingCodec(
      Repository& repo,
      const string& schemaDir,
      const SchemaContext& _schemaCtx);
  BindingCodec() = delete;
  ~BindingCodec() = default;
  BindingCodec(const BindingCodec&) = delete;
  BindingCodec& operator=(const BindingCodec&) = delete;
  BindingCodec(BindingCodec&&) = delete;
  BindingCodec& operator=(BindingCodec&&) = delete;

 private:
  mutex lock; // A codec is expected to be shared, protect it
  CodecServiceProvider codecServiceProvider;
  JsonSubtreeCodec jsonSubtreeCodec;
  const SchemaContext& schemaCtx;

 public:
  string encode(Entity& entity);
  shared_ptr<Entity> decode(const string& payload, shared_ptr<Entity> pointer);

  dynamic toDom(Path path, Entity& entity);
  shared_ptr<Entity> fromDom(
      const dynamic& payload,
      shared_ptr<Entity> pointer);
};

class BindingContext {
 public:
  explicit BindingContext(const Model& model, const SchemaContext& _schemaCtx);
  BindingContext() = delete;
  ~BindingContext() = default;
  BindingContext(const BindingContext&) = delete;
  BindingContext& operator=(const BindingContext&) = delete;
  BindingContext(BindingContext&&) = delete;
  BindingContext& operator=(BindingContext&&) = delete;

 private:
  Repository repo;
  BindingCodec bindingCodec;

 public:
  BindingCodec& getCodec();
};

class BindingSerializationException : public exception {
 private:
  string msg;

 public:
  BindingSerializationException(Entity& _entity, string _cause) {
    std::stringstream buffer;
    buffer << "Failed to encode: " << typeid(_entity).name() << " due to "
           << _cause;
    msg = buffer.str();
  };

  BindingSerializationException(shared_ptr<Entity>& _entity, string _cause) {
    std::stringstream buffer;
    buffer << "Failed to decode: " << typeid(*_entity).name() << " due to "
           << _cause;
    msg = buffer.str();
  };

 public:
  const char* what() const throw() override {
    return msg.c_str();
  }
};

} // namespace cli
} // namespace devices
} // namespace devmand
