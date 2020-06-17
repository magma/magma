#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

import os
from enum import Enum
from typing import Optional


XML_OUTPUT_DIRECTORY: Optional[str] = os.getenv("XML_OUTPUT_DIRECTORY")
TESTS_PATTERN: Optional[str] = os.getenv("TESTS_PATTERN", "*")
TEST_USER_EMAIL = "fbuser@fb.com"
PLATFORM_SERVER_HEALTH_CHECK_URL: str = os.getenv(
    "PLATFORM_SERVER_HEALTH_CHECK_URL", "http://platform-server/healthz"
)


class TestMode(Enum):
    DEV = "dev"
    LOCAL = "local"
    REMOTE = "remote"
