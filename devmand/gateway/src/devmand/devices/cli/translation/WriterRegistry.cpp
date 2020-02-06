// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <boost/exception/diagnostic_information.hpp>
#include <boost/graph/graphviz.hpp>
#include <boost/graph/lookup_edge.hpp>
#include <boost/graph/topological_sort.hpp>
#include <devmand/devices/cli/translation/WriterRegistry.h>

namespace devmand {
namespace devices {
namespace cli {

void WriterRegistryBuilder::addWriter(
    Path path,
    shared_ptr<Writer> writer,
    vector<Path> dependencies) {
  if (schemaContext != SchemaContext::NO_MODELS) {
    if (!schemaContext.isPathValid(path)) {
      throw WriterRegistryException(
          "Unable to register writer for path: " + path.str() +
          ". Path is not valid");
    }
    for (auto& depPath : dependencies) {
      if (!schemaContext.isPathValid(depPath)) {
        throw WriterRegistryException(
            "Unable to register writer for path: " + path.str() +
            ". Dependency path is not valid" + depPath.str());
      }
    }
  }

  allWriters.emplace(path, WriterWithDependencies{writer, dependencies});
}

static vector<Path> topologicalSort(YangHierarchy pathGraph) {
  vector<YangHierarchy::vertex_descriptor> topologicalOrderIndex;
  try {
    topological_sort(pathGraph, back_inserter(topologicalOrderIndex));
  } catch (boost::not_a_dag& e) {
    throw WriterRegistryException(
        "Unable to build writer registry, dependency cycle detected");
  } catch (exception& e) {
    throw WriterRegistryException(
        "Unable to build writer registry, unknown error: " + string(e.what()));
  }
  vector<Path> topologicalOrder;
  for (auto& index : topologicalOrderIndex) {
    topologicalOrder.push_back(pathGraph.graph()[index].path);
  }
  return topologicalOrder;
}

unique_ptr<WriterRegistry> WriterRegistryBuilder::build() {
  if (allWriters.empty()) {
    return make_unique<WriterRegistry>(
        map<Path, shared_ptr<Writer>, SortedPathComparator>(
            SortedPathComparator({})));
  }

  YangHierarchy pathGraph;

  for (auto& writer : allWriters) {
    Path childPath = writer.first;
    shared_ptr<Writer> childHandler = writer.second.writer;
    vector<Path> childDependencies = writer.second.dependencies;

    YangHierarchy::vertex_descriptor childVertex =
        boost::add_vertex(childPath, pathGraph);
    pathGraph[childPath].path = childPath;

    for (auto& depPath : childDependencies) {
      YangHierarchy::vertex_descriptor depVertex =
          boost::add_vertex(depPath, pathGraph);
      pathGraph[depPath].path = depPath;

      if (!boost::lookup_edge(depVertex, childVertex, pathGraph).second) {
        boost::add_edge_by_label(depPath, childPath, pathGraph);
      }
    }
  }

  // order writers by their dependencies
  vector<Path> topologicalOrder = topologicalSort(pathGraph);

  if (VLOG_IS_ON(MDEBUG)) {
    stringstream dotGraph;
    boost::write_graphviz(dotGraph, pathGraph, PathVertexWriter(pathGraph));
    MLOG(MDEBUG) << "Writer hierarchy calculated: " << dotGraph.str();
  }

  // put writers into and ordered map to preserve dependencies among them
  auto orderedWriters = map<Path, shared_ptr<Writer>, SortedPathComparator>(
      SortedPathComparator(topologicalOrder));
  for (auto& writerWithDeps : allWriters) {
    orderedWriters.insert(
        make_pair(writerWithDeps.first, writerWithDeps.second.writer));
  }

  return make_unique<WriterRegistry>(orderedWriters);
}

ostream& operator<<(ostream& os, const WriterRegistry& registry) {
  os << "WriterRegistry:";
  os << endl;
  for (auto& writer : registry.orderedWriters) {
    os << " " << writer.first << endl;
  }
  return os;
}

static bool writerHasUpdates(
    const WriterRegistry::Diff& diff,
    const Path& registeredPath) {
  return diff.find(registeredPath) != diff.end();
}

static void writeSingle(
    const DatastoreDiff& update,
    const Writer& writer,
    const DeviceAccess& device) {
  (void)update;
  (void)writer;
  (void)device;
//  switch (update.type) {
//    case DatastoreDiffType::create: {
//      writer.create(update.keyedPath, update.after, device).get();
//      break;
//    }
//    case DatastoreDiffType::update: {
//      writer.update(update.keyedPath, update.before, update.after, device).get();
//      break;
//    }
//    case DatastoreDiffType::deleted: {
//      writer.remove(update.keyedPath, update.before, device).get();
//      break;
//    }
//  }
}

void WriterRegistry::write(const Diff& diff, const DeviceAccess& device) const {
  for (auto& writer : orderedWriters) {
    auto& registeredPath = writer.first;
    if (writerHasUpdates(diff, registeredPath)) {
      auto updateRange = diff.equal_range(registeredPath);
      for (auto itr = updateRange.first; itr != updateRange.second; ++itr) {
        DatastoreDiff update = (*itr).second;
        // TODO handle exceptions and do rollback
        writeSingle(update, *writer.second, device);
      }
    }
  }
}

vector<Path> WriterRegistry::getWriterPaths() const {
  vector<Path> pathsCopy = vector<Path>();
  for (auto& writer : orderedWriters) {
    pathsCopy.push_back(writer.first);
  }
  return pathsCopy;
}

} // namespace cli
} // namespace devices
} // namespace devmand
