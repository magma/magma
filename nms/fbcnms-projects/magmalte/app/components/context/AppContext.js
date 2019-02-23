/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {SessionUser} from '../../common/UserModel';

import React from 'react';

type Context = {
  csrfToken: ?string,
  version: ?string,
  networkIds: string[],
  user: SessionUser,
};

export default React.createContext<Context>({
  csrfToken: null,
  version: null,
  networkIds: [],
  user: {email: '', isSuperUser: false},
});
