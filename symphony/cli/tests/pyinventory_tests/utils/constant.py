#!/usr/bin/env python3
# pyre-strict
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

import os
from typing import Optional


XML_OUTPUT_DIRECTORY: Optional[str] = os.getenv("XML_OUTPUT_DIRECTORY")
TEST_USER_EMAIL = "fbuser@fb.com"
PLATFORM_SERVER_HEALTH_CHECK_URL: str = os.getenv(
    "PLATFORM_SERVER_HEALTH_CHECK_URL", "http://platform-server/healthz"
)
