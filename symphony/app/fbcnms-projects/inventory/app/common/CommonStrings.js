/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import fbt from 'fbt';

const Strings = {
  common: {
    emptyField: `${fbt(
      'None',
      'Text to be displayed in case a user input field has no value',
    )}`,
    unassignedItem: `${fbt(
      'Unassigned',
      'Text to be displayed in case an assignable item was not assigned yet',
    )}`,
  },
};

export default Strings;
