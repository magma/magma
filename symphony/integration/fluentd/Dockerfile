# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# Extend base fluentd image
FROM fluent/fluentd:v1.9-1

USER root

# Install prometheus plugin
RUN fluent-gem install fluent-plugin-prometheus
