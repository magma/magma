// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <devmand/devices/cli/translation/BindingReaderRegistry.h>
#include <devmand/devices/cli/translation/BindingWriterRegistry.h>
#include <devmand/devices/cli/translation/PluginRegistry.h>

namespace devmand {
namespace devices {
namespace cli {

using namespace devmand::channels::cli;

class UbntInterfacePlugin : public Plugin {
 public:
  UbntInterfacePlugin(BindingContext& openconfigContext);

  DeviceType getDeviceType() const override;
  void provideReaders(ReaderRegistryBuilder& registry) const override;
  void provideWriters(WriterRegistryBuilder& registry) const override;

 private:
  BindingContext& openconfigContext;
};

} // namespace cli
} // namespace devices
} // namespace devmand
