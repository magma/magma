/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import grey from '@material-ui/core/colors/grey';
import orange from '@material-ui/core/colors/orange';
import red from '@material-ui/core/colors/red';
import yellow from '@material-ui/core/colors/yellow';

type SeverityMap = {
  [string]: {
    name: string,
    order: number,
    color: string,
  },
};

export const SEVERITY: SeverityMap = {
  NOTICE: {
    name: 'NOTICE',
    order: 0,
    color: grey[500],
  },
  INFO: {
    name: 'INFO',
    order: 1,
    color: grey[500],
  },
  WARNING: {
    name: 'WARNING',
    order: 2,
    color: yellow.A400,
  },
  MINOR: {
    name: 'MINOR',
    order: 3,
    color: yellow.A400,
  },
  MAJOR: {
    name: 'MAJOR',
    order: 4,
    color: orange.A400,
  },
  CRITICAL: {
    name: 'CRITICAL',
    order: 5,
    color: red.A400,
  },
};
