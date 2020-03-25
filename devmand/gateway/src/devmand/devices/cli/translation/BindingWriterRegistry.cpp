// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/devices/cli/translation/BindingWriterRegistry.h>

namespace devmand {
namespace devices {
namespace cli {

using namespace folly;
using namespace std;
using namespace ydk;

BindingWriterRegistryBuilder::BindingWriterRegistryBuilder(
    WriterRegistryBuilder& _domBuilder,
    BindingContext& _context)
    : domBuilder(_domBuilder), context(_context) {}

} // namespace cli
} // namespace devices
} // namespace devmand
