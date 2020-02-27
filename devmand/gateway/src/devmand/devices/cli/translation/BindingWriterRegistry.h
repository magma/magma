// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <boost/graph/adjacency_list.hpp>
#include <boost/graph/labeled_graph.hpp>
#include <devmand/devices/cli/schema/BindingContext.h>
#include <devmand/devices/cli/schema/Path.h>
#include <devmand/devices/cli/schema/SchemaContext.h>
#include <devmand/devices/cli/translation/DeviceAccess.h>
#include <devmand/devices/cli/translation/WriterRegistry.h>
#include <folly/dynamic.h>
#include <folly/futures/Future.h>
#include <folly/json.h>
#include <ydk/types.hpp>
#include <ostream>

namespace devmand {
namespace devices {
namespace cli {

using namespace folly;
using namespace std;
using namespace ydk;

// TODO extract client-facing SPI into a separate *SPI header file

template <
    typename YDKTYPE,
    typename enable_if<is_base_of<Entity, YDKTYPE>{}, int>::type = false>
class BindingWriter {
 public:
  using EntityType = YDKTYPE;

  virtual Future<Unit> create(
      const Path& path,
      shared_ptr<YDKTYPE> cfg,
      const DeviceAccess& device) const = 0;
  virtual Future<Unit> update(
      const Path& path,
      shared_ptr<YDKTYPE> before,
      shared_ptr<YDKTYPE> after,
      const DeviceAccess& device) const = 0;
  virtual Future<Unit> remove(
      const Path& path,
      shared_ptr<YDKTYPE> before,
      const DeviceAccess& device) const = 0;
};

template <typename YDKTYPE>
class BindingWriterAdapter : public Writer {
 protected:
  shared_ptr<BindingWriter<YDKTYPE>> bindingWriter;
  BindingContext& context;

 public:
  BindingWriterAdapter(
      shared_ptr<BindingWriter<YDKTYPE>> _bindingWriter,
      BindingContext& _context)
      : bindingWriter(_bindingWriter), context(_context){};

  Future<Unit> create(const Path& path, dynamic cfg, const DeviceAccess& device) const {
    shared_ptr<YDKTYPE> ptr = makePtr();
    context.getCodec().fromDom(cfg, ptr);
    return bindingWriter->create(path, ptr, device);
  };

  Future<Unit> update(
      const Path& path,
      dynamic before,
      dynamic after,
      const DeviceAccess& device) const {
    shared_ptr<YDKTYPE> ptrBefore = static_pointer_cast<YDKTYPE>(
        context.getCodec().fromDom(before, makePtr()));
    shared_ptr<YDKTYPE> ptrAfter = static_pointer_cast<YDKTYPE>(
        context.getCodec().fromDom(after, makePtr()));
    return bindingWriter->update(path, ptrBefore, ptrAfter, device);
  };

  Future<Unit> remove(const Path& path, dynamic before, const DeviceAccess& device)
      const {
    shared_ptr<YDKTYPE> ptrBefore = makePtr();
    context.getCodec().fromDom(before, ptrBefore);
    return bindingWriter->remove(path, ptrBefore, device);
  };

 private:
  shared_ptr<YDKTYPE> makePtr() const {
    return make_shared<YDKTYPE>();
  }
};

class BindingWriterRegistryBuilder {
 private:
  WriterRegistryBuilder& domBuilder;
  BindingContext& context;

 public:
  BindingWriterRegistryBuilder(
      WriterRegistryBuilder& _domBuilder,
      BindingContext& _context);
  // No validation against schema will be performed with NO_MODELS context
  ~BindingWriterRegistryBuilder() = default;
  BindingWriterRegistryBuilder(const BindingWriterRegistryBuilder&) = delete;
  BindingWriterRegistryBuilder& operator=(const BindingWriterRegistryBuilder&) =
      delete;
  BindingWriterRegistryBuilder(BindingWriterRegistryBuilder&&) = delete;
  BindingWriterRegistryBuilder& operator=(BindingWriterRegistryBuilder&&) =
      delete;

  // TODO support subtree / wildcarded writers

  template <
      typename T,
      typename YDKENTITY = typename T::EntityType,
      typename enable_if<is_base_of<BindingWriter<YDKENTITY>, T>{}, int>::type =
          false>
  void add(Path path, shared_ptr<T> writer, vector<Path> dependencies = {}) {
    domBuilder.add(
        path,
        make_shared<BindingWriterAdapter<YDKENTITY>>(writer, context),
        dependencies);
  }
};

#define BINDING_W(reg, ctx) BindingWriterRegistryBuilder(reg, ctx)

} // namespace cli
} // namespace devices
} // namespace devmand
