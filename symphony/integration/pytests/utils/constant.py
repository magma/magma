#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

import os


XML_OUTPUT_DIRECTORY = os.getenv("XML_OUTPUT_DIRECTORY")
TEST_USER_EMAIL = "fbuser@fb.com"
PLATFORM_SERVER_HEALTH_CHECK_URL = os.getenv(
    "PLATFORM_SERVER_HEALTH_CHECK_URL", "http://platform-server/healthz"
)
