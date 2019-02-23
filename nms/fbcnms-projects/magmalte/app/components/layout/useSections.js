/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {Section} from '../layout/Section';

import {getLteSections} from '../lte/LteSections';

export default function useSections(): Section[] {
  return getLteSections();
}
