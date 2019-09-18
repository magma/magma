// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

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
