// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <boost/graph/graphviz.hpp>
#include <boost/graph/lookup_edge.hpp>
#include <devmand/devices/cli/translation/ReaderRegistry.h>

namespace devmand {
namespace devices {
namespace cli {

using namespace std;
using namespace folly;

const shared_ptr<Reader> StructuralReader::INSTANCE =
    static_pointer_cast<Reader>(make_shared<StructuralReader>());

Future<dynamic> CompositeReader::read(
    const Path& path,
    bool config,
    const DeviceAccess& device) const {
  if (shouldDelegate(path)) {
    auto childName = path.unkeyed().getSegments()[registeredPath.getDepth()];
    if (children.find(childName) == children.end()) {
      return Future<dynamic>(ReadException(
          device.id(),
          path,
          "No reader registered at: " +
              registeredPath.getChild(childName).str()));
    }
    return children.at(childName)->read(path, config, device);
  }

  // If reading config and this is oper, skip
  if (!isConfig and config) {
    return Future<dynamic>(dynamic::object());
  }

  vector<Future<dynamic>> childrenFutures;
  childrenFutures.emplace_back(clientReader->read(path, device));
  for (auto& child : children) {
    childrenFutures.emplace_back(
        child.second->read(path.getChild(child.first), config, device));
  }

  return collect(childrenFutures.begin(), childrenFutures.end())
      .thenValue([](auto childrenData) {
        dynamic localData = childrenData[0];
        for (const auto& childData : childrenData) {
          if (childData == localData) {
            continue;
          }
          if (!childData.empty() &&
              !childData.items().begin()->second.empty()) {
            // add child data only if it's not empty
            localData.merge_patch(childData);
          }
        }
        return localData;
      })
      .thenValue([registeredPath = registeredPath](auto mergedData) -> dynamic {
        if (registeredPath == "/") {
          return mergedData;
        } else {
          return dynamic::object(registeredPath.getLastSegment(), mergedData);
        }
      });
}

Future<dynamic> CompositeListReader::read(
    const Path& path,
    bool config,
    const DeviceAccess& device) const {
  string registeredSegment = registeredPath.getLastSegment();
  if (!path.getKeysFromSegment(registeredSegment).empty()) {
    // Read single list item since we have the keys already
    return readSingle(path, config, device);
  }

  // Read entire list
  return static_pointer_cast<ListReader>(clientReader)
      ->readKeys(path, device)
      // FIXME passing "this" into lambda ... but a reader should exist
      // basically forever and only go away when program ends .. we don't
      // expect translation code to be registered and unregistered in runtime
      .thenValue([path, config, device, this](auto childrenKeys) {
        vector<Future<dynamic>> keyFutures;
        for (const auto& childKeys : childrenKeys) {
          keyFutures.emplace_back(
              this->readSingle(path.addKeys(childKeys), config, device));
        }

        return collect(keyFutures.begin(), keyFutures.end());
      })
      .thenValue([registeredSegment](auto childrenData) {
        auto listEntries = dynamic::array();
        for (const auto& childData : childrenData) {
          listEntries.push_back(childData.at(registeredSegment)[0]);
        }
        return dynamic::object(registeredSegment, listEntries);
      });
}

Future<dynamic> CompositeListReader::readSingle(
    const Path& path,
    bool config,
    const DeviceAccess& device) const {
  if (shouldDelegate(path)) {
    return CompositeReader::read(path, config, device);
  }

  // Do some encapsulation if reading a list entry
  return CompositeReader::read(path, config, device)
      .thenValue([registeredSegment =
                      registeredPath.getLastSegment()](auto singleEntryData) {
        return dynamic::object(
            registeredSegment,
            dynamic::array(singleEntryData.at(registeredSegment)));
      });
}

ostream& operator<<(ostream& os, const CompositeReader& reader) {
  os << string(reader.registeredPath.getDepth(), ' ') << reader.registeredPath;
  for (auto& child : reader.children) {
    os << endl << *child.second;
  }
  return os;
}

bool CompositeReader::shouldDelegate(const Path& path) const {
  return path.unkeyed() != registeredPath &&
      path.unkeyed().getDepth() > registeredPath.getDepth();
}

CompositeReader::CompositeReader(
    Path _registeredPath,
    shared_ptr<Reader> _clientReader,
    CompositeReader::Children&& _children,
    bool _isConfig)
    : registeredPath(_registeredPath),
      clientReader(_clientReader),
      children(forward<Children>(_children)),
      isConfig(_isConfig) {}

static bool isListReader(const shared_ptr<Reader>& reader) {
  return dynamic_pointer_cast<ListReader>(reader) != NULL;
}

CompositeListReader::CompositeListReader(
    Path _registeredPath,
    shared_ptr<Reader> _clientReader,
    CompositeReader::Children&& _children,
    bool _isConfig)
    : CompositeReader(
          _registeredPath,
          _clientReader,
          forward<Children>(_children),
          _isConfig) {
  if (!isListReader(_clientReader)) {
    throw ReaderRegistryException(
        "Unable to register non list reader for path: " +
        _registeredPath.str());
  }
}

Future<dynamic> ReaderRegistry::readConfiguration(
    const Path& path,
    const DeviceAccess& device) const {
  // TODO validate returned data
  return (rootReader->read(path, true, device)).thenValue([](auto data) {
    dynamic root = dynamic::object;
    root.merge_patch(data);
    return root;
  });
}

Future<dynamic> ReaderRegistry::readState(
    const Path& path,
    const DeviceAccess& device) const {
  // TODO validate returned data
  return (rootReader->read(path, false, device)).thenValue([](auto data) {
    dynamic root = dynamic::object;
    root.merge_patch(data);
    return root;
  });
}

ostream& operator<<(ostream& os, const ReaderRegistry& registry) {
  os << "ReaderRegistry: " << *registry.rootReader;
  return os;
}

unique_ptr<ReaderRegistry> ReaderRegistryBuilder::build() {
  if (allReaders.empty()) {
    ReaderRegistryException(
        "Unable to construct reader registry from 0 readers");
  }

  YangHierarchy pathGraph;

  for (auto& reader : allReaders) {
    Path childPath = reader.first;

    while (childPath.getDepth() > 0) {
      const Path& parentPath = childPath.getParent();

      YangHierarchy::vertex_descriptor parentVertex =
          boost::add_vertex(parentPath, pathGraph);
      pathGraph[parentPath].path = parentPath;

      YangHierarchy::vertex_descriptor childVertex =
          boost::add_vertex(childPath, pathGraph);
      pathGraph[childPath].path = childPath;

      if (!boost::lookup_edge(parentVertex, childVertex, pathGraph).second) {
        boost::add_edge_by_label(parentPath, childPath, pathGraph);
      }
      childPath = parentPath;
    }
  }

  if (VLOG_IS_ON(MDEBUG)) {
    stringstream dotGraph;
    boost::write_graphviz(dotGraph, pathGraph, PathVertexWriter(pathGraph));
    MLOG(MDEBUG) << "Reader hierarchy calculated: " << dotGraph.str();
  }

  PathVertex& rootVertex = pathGraph[Path::ROOT];
  return make_unique<ReaderRegistry>(
      createCompositeReader(pathGraph, rootVertex));
}

void ReaderRegistryBuilder::add(Path path, ReaderFunction lambda) {
  addReader(
      path, static_pointer_cast<Reader>(make_shared<LambdaReader>(lambda)));
}

void ReaderRegistryBuilder::addList(
    Path path,
    ListReaderFunction keysLambda,
    ReaderFunction lambda) {
  // schema validation
  if (schemaContext != SchemaContext::NO_MODELS) {
    if (!schemaContext.isPathValid(path) or !schemaContext.isList(path)) {
      throw ReaderRegistryException(
          "Unable to register list reader for path: " + path.str() +
          ". Path is not valid or does not point to a list node");
    }
  }

  addReader(
      path,
      static_pointer_cast<Reader>(
          make_shared<LambdaListReader>(keysLambda, lambda)));
}

void ReaderRegistryBuilder::addReader(Path path, shared_ptr<Reader> reader) {
  if (schemaContext != SchemaContext::NO_MODELS) {
    if (!schemaContext.isPathValid(path)) {
      throw ReaderRegistryException(
          "Unable to register reader for path: " + path.str() +
          ". Path is not valid");
    }
  }

  allReaders.insert_or_assign(path, reader);
}

unique_ptr<CompositeReader> ReaderRegistryBuilder::createCompositeReader(
    const YangHierarchy& pathGraph,
    const PathVertex& vertex) const {
  auto vertexIdx = pathGraph.vertex(vertex.path);

  // Create children composite nodes recursively
  CompositeReader::Children children;
  YangHierarchy::out_edge_iterator ei, ei_end;
  for (boost::tie(ei, ei_end) = boost::out_edges(vertexIdx, pathGraph);
       ei != ei_end;
       ++ei) {
    YangHierarchy::vertex_descriptor childVertexIdx =
        boost::target(*ei, pathGraph);
    const PathVertex& childVertex = pathGraph.graph()[childVertexIdx];
    children.emplace(
        childVertex.path.getLastSegment(),
        this->createCompositeReader(pathGraph, childVertex));
  }

  // Create composite reader for current path
  auto clientReader = allReaders.find(vertex.path);

  // Determine if this is config or operational node
  bool isConfig = true;
  if (schemaContext != SchemaContext::NO_MODELS) {
    isConfig = schemaContext.isConfig(vertex.path);
  }

  if (clientReader != allReaders.end()) {
    if (isListReader(clientReader->second)) {
      return make_unique<CompositeListReader>(
          vertex.path, clientReader->second, move(children), isConfig);
    } else {
      return make_unique<CompositeReader>(
          vertex.path, clientReader->second, move(children), isConfig);
    }
  } else {
    // schema validation
    if (schemaContext != SchemaContext::NO_MODELS) {
      if (schemaContext.isList(vertex.path)) {
        throw ReaderRegistryException(
            "Missing list reader for path: " + vertex.path.str() +
            ". List data cannot be generated. Reader hierarchy is invalid");
      }
    }
    return make_unique<CompositeReader>(
        vertex.path, StructuralReader::INSTANCE, move(children), isConfig);
  }
}

} // namespace cli
} // namespace devices
} // namespace devmand
