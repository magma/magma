#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

__version__ = "0.0.1"

SCHEMA_FILE_NAME = "survey_schema.json"
SIMPLE_QUESTION_TYPE_TO_REQUIRED_PROPERTY_NAME = {
    "DATE": "dateData",
    "BOOL": "boolData",
    "EMAIL": "emailData",
    "TEXT": "textData",
    "FLOAT": "floatData",
    "INTEGER": "intData",
    "PHONE": "phoneData",
}
