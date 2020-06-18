// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <devmand/channels/cli/datastore/DatastoreDiff.h>
#include <devmand/devices/cli/translation/ReaderRegistry.h>

namespace devmand {
namespace devices {
namespace cli {

using namespace folly;
using namespace devmand::channels::cli::datastore;
using namespace std;

class Writer {
 public:
  virtual Future<Unit>
  create(const Path& path, dynamic cfg, const DeviceAccess& device) const = 0;
  virtual Future<Unit> update(
      const Path& path,
      dynamic before,
      dynamic after,
      const DeviceAccess& device) const = 0;
  virtual Future<Unit> remove(
      const Path& path,
      dynamic before,
      const DeviceAccess& device) const = 0;
};

struct WriterWithDependencies {
  shared_ptr<Writer> writer;
  vector<Path> dependencies;
};

struct SortedPathComparator {
  vector<Path> topologicalOrder;

  SortedPathComparator(const vector<Path>& _topologicalOrder)
      : topologicalOrder(_topologicalOrder) {}

  long indexOf(const Path& path) const {
    auto itr = find(topologicalOrder.begin(), topologicalOrder.end(), path);

    if (itr != topologicalOrder.cend()) {
      return distance(topologicalOrder.begin(), itr);
    } else {
      return -1;
    }
  }

  bool operator()(const Path& a, const Path& b) const {
    return indexOf(a) > indexOf(b);
  }
};

class WriterRegistry {
 private:
  const map<Path, shared_ptr<Writer>, SortedPathComparator> orderedWriters;

 public:
  WriterRegistry(
      map<Path, shared_ptr<Writer>, SortedPathComparator> _orderedWriters)
      : orderedWriters(_orderedWriters){};
  ~WriterRegistry() = default;
  WriterRegistry(const WriterRegistry&) = delete;
  WriterRegistry& operator=(const WriterRegistry&) = delete;
  WriterRegistry(WriterRegistry&&) = delete;
  WriterRegistry& operator=(WriterRegistry&&) = delete;

  friend ostream& operator<<(ostream& os, const WriterRegistry& registry);

  typedef multimap<Path, DatastoreDiff> Diff;
  void write(const Diff& diff, const DeviceAccess& device) const;

  vector<Path> getWriterPaths() const;
};

class WriterRegistryException : public runtime_error {
 public:
  WriterRegistryException(string reason) : runtime_error(reason){};
};

class WriterRegistryBuilder {
 private:
  map<Path, WriterWithDependencies> allWriters;
  const SchemaContext& schemaContext;

 public:
  WriterRegistryBuilder(const SchemaContext& _schemaContext)
      : schemaContext(_schemaContext){};
  // No validation against schema will be performed with NO_MODELS context
  WriterRegistryBuilder() : WriterRegistryBuilder(SchemaContext::NO_MODELS){};
  ~WriterRegistryBuilder() = default;
  WriterRegistryBuilder(const WriterRegistryBuilder&) = delete;
  WriterRegistryBuilder& operator=(const WriterRegistryBuilder&) = delete;
  WriterRegistryBuilder(WriterRegistryBuilder&&) = delete;
  WriterRegistryBuilder& operator=(WriterRegistryBuilder&&) = delete;

  unique_ptr<WriterRegistry> build();

  // TODO support subtree / wildcarded Writers

  template <
      typename T,
      typename enable_if<is_base_of<Writer, T>{}, int>::type = false>
  void add(Path path, shared_ptr<T> writer, vector<Path> dependencies = {}) {
    addWriter(path, static_pointer_cast<Writer>(writer), dependencies);
  }

 private:
  void addWriter(
      Path path,
      shared_ptr<Writer> writer,
      vector<Path> dependencies = {});
};

} // namespace cli
} // namespace devices
} // namespace devmand
