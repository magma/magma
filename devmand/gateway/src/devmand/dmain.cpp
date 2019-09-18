// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <memory>

#include <folly/init/Init.h>

#include <devmand/Application.h>

int main(int argc, char* argv[]) {
  folly::init(&argc, &argv);

  devmand::Application app;
  app.run();
  return app.status();
}
