#! /bin/sh

set -e # exit on any error

TEMP_FILE=$(mktemp)
yarn --silent swagger2js gen swagger.yml -t flow -c MagmaAPIBindings -b > "$TEMP_FILE"

HEADER="/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @REPLACE_WITH_GENERATED_TOKEN
 */
"

OUTPUT=__generated__/MagmaAPIBindings.js

# adding this sed command to avoid having Phabricator think this file is
# generated since it looks for the "generated" keyword
(echo "$HEADER"; cat "$TEMP_FILE") | sed -e "s#REPLACE_WITH_GENERATED_TOKEN#generated#" >$OUTPUT
