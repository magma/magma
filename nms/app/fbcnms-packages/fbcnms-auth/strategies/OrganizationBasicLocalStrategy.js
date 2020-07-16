/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import {BasicStrategy} from 'passport-http';
import {validateUser} from './OrganizationLocalStrategy';

export default function () {
  return new BasicStrategy(
    {
      realm: 'Users',
      passReqToCallback: true,
    },
    validateUser,
  );
}
