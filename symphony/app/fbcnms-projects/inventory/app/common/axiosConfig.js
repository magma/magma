/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
'use strict';

import axios from 'axios';

axios.defaults.headers.common['X-CSRF-Token'] = window.CONFIG.appData.csrfToken;

export default axios;
