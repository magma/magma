# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
#
# Sync folders.
module CloudConfigs
  require 'yaml'

  modules_file = ENV["MAGMA_MODULES_FILE"] || "../../modules.yml"
  if File.exist?(File.expand_path("../../fb/config/modules.yml"))
    modules_file = "../../fb/config/modules.yml"
  end
  module_config = YAML.load_file(modules_file)

  $repos = [
    {:host_path => "../../../magma", :mount_path => "/home/vagrant/magma"},
  ]
  module_config["external_modules"].each do |mod|
    $repos << {:host_path => mod["host_path"], :mount_path => mod["mount_path"]}
  end
end
