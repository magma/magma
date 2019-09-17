// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

#include <memory>

#include <folly/init/Init.h>

#include <devmand/Application.h>

int main(int argc, char* argv[]) {
  folly::init(&argc, &argv);

  devmand::Application app;
  app.run();
  return app.status();
}
