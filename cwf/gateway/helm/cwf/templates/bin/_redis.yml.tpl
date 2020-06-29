---
#
# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

port: {{ .Values.cwf.redis.port }}
bind: {{ .Values.cwf.redis.bind }}
redis_loglevel: notice
dir: /var/opt/magma
# How frequently to save/dump to disk.
# E.g. the first element indicates we save every 900 seconds
# if at least 1 key has changed.
save:
  - seconds: 900
    num_keys: 1
  - seconds: 300
    num_keys: 10
  - seconds: 60
    num_keys: 1000