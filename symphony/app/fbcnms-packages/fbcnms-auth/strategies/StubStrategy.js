/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import {Strategy} from 'passport-strategy';

//use this in place of the real openid strategies until discovery finishes
export default class StubStrategy extends Strategy {
  constructor() {
    super();
    this.name = 'stub';
  }
  authenticate() {
    return this.fail({message: 'No implementation found for strategy'});
  }
}
