// Copyright (c) 2020-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <devmand/channels/cli/datastore/DatastoreDiff.h>
#include <devmand/channels/cli/datastore/DatastoreState.h>
#include <devmand/devices/cli/schema/BindingContext.h>
#include <devmand/devices/cli/schema/ModelRegistry.h>
#include <devmand/devices/cli/schema/Path.h>
#include <folly/dynamic.h>
#include <folly/json.h>
#include <libyang/libyang.h>
#include <magma_logging.h>
#include <ydk/types.hpp>
#include <atomic>

using devmand::channels::cli::datastore::DatastoreDiff;
using devmand::channels::cli::datastore::DatastoreDiffType;
using devmand::channels::cli::datastore::DatastoreState;
using devmand::devices::cli::Model;
using devmand::devices::cli::ModelRegistry;
using devmand::devices::cli::Path;
using devmand::devices::cli::SchemaContext;
using folly::dynamic;
using folly::Optional;
using folly::parseJson;
using std::atomic_bool;
using std::make_shared;
using std::map;
using std::multimap;
using std::pair;
using std::runtime_error;
using std::shared_ptr;
using std::string;
using std::vector;
using ydk::Entity;

namespace devmand::channels::cli::datastore {

struct DiffPath {
  Path path;
  bool asterix;

  DiffPath() : path("/"), asterix(false) {}

  DiffPath(const Path _path, bool _asterix) : path(_path), asterix(_asterix) {}
};

struct DiffResult {
  multimap<Path, DatastoreDiff> diffs;
  vector<Path> unhandledDiffs;
  void appendUnhandledPath(Path pathToAppend) {
    unhandledDiffs.emplace_back(pathToAppend);
  }
};

class DatastoreTransaction {
 private:
  shared_ptr<DatastoreState> datastoreState;
  lllyd_node* root = nullptr;
  atomic_bool hasCommited = ATOMIC_VAR_INIT(false);
  SchemaContext& schemaContext;

  static lllyd_node* computeRoot(lllyd_node* n);
  int datastoreTypeToLydOption();
  lllyd_node* dynamic2lydNode(dynamic entity);
  static lllyd_node*
  getExistingNode(lllyd_node* a, lllyd_node* b, DatastoreDiffType type);
  static vector<LydPair> splitNodeToRoots(lllyd_node* node);
  static void printDiffType(LLLYD_DIFFTYPE type);
  static string buildFullPath(lllyd_node* node, string pathSoFar);
  static DatastoreDiffType getDiffType(LLLYD_DIFFTYPE type);
  void print(lllyd_node* nodeToPrint);
  void checkIfCommitted();
  string toJson(lllyd_node* initial);
  static void addKeysToPath(lllyd_node* node, std::stringstream& path);
  static string makePrefixedSegment(lllyd_node* node);
  static bool shouldHandleSubtree(
      const DiffPath& registeredPath,
      const Path& changedPath);
  static bool isExactPath(
      const DiffPath& registeredPath,
      const Path& changedPath);
  static bool isAboveChange(
      const DiffPath& registeredPath,
      const Path& changedPath);
  static bool segmentDiffOneOrLess(
      const Path& noTotifyPath,
      const Path& changedPath);
  static bool isPickableDelete(
      DatastoreDiffType type,
      const DiffPath& toNotifyPath,
      const Path& changedPath);
  static bool isPickableCreate(
      DatastoreDiffType type,
      const DiffPath& toNotifyPath,
      const Path& changedPath);
  static dynamic appendAllParents(Path path, const dynamic& aDynamic);
  static Path unifyLength(Path registeredPath, Path keyedPath);
  vector<DiffPath>
  pickClosestPath(Path path, vector<DiffPath> paths, DatastoreDiffType type);
  static vector<DiffPath> matchClosesUpdatePath(
      Path& modifiedPath,
      vector<DiffPath>& registeredPaths);
  map<Path, DatastoreDiff> splitDiff(DatastoreDiff diff);
  void
  splitToMany(Path p, dynamic input, vector<std::pair<string, dynamic>>& v);
  vector<Path> getRegisteredPath(
      vector<DiffPath> registeredPaths,
      Path path,
      DatastoreDiffType type);
  dynamic read(Path path, lllyd_node* node);
  dynamic readAlreadyCommitted(Path path);
  map<Path, DatastoreDiff> libyangDiff(lllyd_node* a, lllyd_node* b);
  string appendKey(dynamic data, string path);
  llly_set* findNode(lllyd_node* node, string path);

 public:
  map<Path, DatastoreDiff> diff();

  DatastoreTransaction(
      shared_ptr<DatastoreState> datastoreState,
      SchemaContext& _schemaContext);
  dynamic read(Path path);
  void print();
  DiffResult diff(vector<DiffPath> registeredPaths);

  void isValid(); // throws exception of invalid
  bool delete_(Path path);
  void merge(Path path, const dynamic& aDynamic);
  void overwrite(Path path, const dynamic& aDynamic);
  void commit();
  void abort();

  virtual ~DatastoreTransaction();
};
} // namespace devmand::channels::cli::datastore
