// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/devices/cli/translation/BindingReaderRegistry.h>

namespace devmand {
namespace devices {
namespace cli {

using namespace folly;
using namespace std;
using namespace ydk;

static dynamic
entityToDom(BindingContext& context, Entity& entity, const Path& path) {
  dynamic asDynamic = context.getCodec().toDom(path, entity);
  // the object is prefixed with its name inside dynamic, extract
  return asDynamic.items().begin()->second;
}

Future<dynamic> BindingReaderAdapter::read(
    const Path& path,
    const DeviceAccess& device) const {
  return bindingReader->read(path, device)
      .thenValue([path, &context = this->context](auto entity) {
        return entityToDom(context, *entity, path);
      });
}

BindingReaderAdapter::BindingReaderAdapter(
    shared_ptr<BindingReader> _bindingReader,
    BindingContext& _context)
    : bindingReader(_bindingReader), context(_context) {}

BindingListReaderAdapter::BindingListReaderAdapter(
    shared_ptr<BindingListReader> _bindingReader,
    BindingContext& _context)
    : bindingReader(_bindingReader), context(_context) {}

Future<vector<dynamic>> BindingListReaderAdapter::readKeys(
    const Path& path,
    const DeviceAccess& device) const {
  return bindingReader->readKeys(path, device).thenValue([](auto entityKeys) {
    vector<dynamic> transformed;
    for (EntityKeys& entityKey : entityKeys) {
      dynamic transformedSingle = dynamic::object();
      for (YLeaf& entityKeyLeaf : entityKey) {
        transformedSingle[entityKeyLeaf.name] = entityKeyLeaf.value;
      }
      transformed.push_back(transformedSingle);
    }
    return transformed;
  });
}

Future<dynamic> BindingListReaderAdapter::read(
    const Path& path,
    const DeviceAccess& device) const {
  return bindingReader->read(path, device)
      .thenValue([path, &context = this->context](auto entity) {
        if (entity == nullptr) {
          // no additional data to add to list, just the keys
          return path.getKeys();
        } else {
          return entityToDom(context, *entity, path);
        }
      });
}
} // namespace cli
} // namespace devices
} // namespace devmand
