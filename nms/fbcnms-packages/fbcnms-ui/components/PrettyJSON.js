/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import * as React from 'react';

type Json =
  | void
  | null
  | string
  | number
  | boolean
  | Array<Json>
  | $Shape<{[string]: Json}>;

type Props = {
  jsonObject: Json,
};

export default function PrettyJSON({jsonObject}: Props) {
  return (
    <pre style={{whiteSpace: 'pre-wrap'}}>
      {JSON.stringify(jsonObject, null, 2)}
    </pre>
  );
}
