// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

#include <devmand/cartography/Method.h>

namespace devmand {
namespace cartography {

void Method::setHandlers(
    const AddHandler& addHandler,
    const DeleteHandler& deleteHandler) {
  add = addHandler;
  del = deleteHandler;
}

} // namespace cartography
} // namespace devmand
