#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

__version__ = "2.6.1"

EQUIPMENTS_TO_SEARCH = 10
LOCATIONS_TO_SEARCH = 5
USER_ROLE = 0
SUPERUSER_ROLE = 3
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
