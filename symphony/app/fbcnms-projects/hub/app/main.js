/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
'use strict';

import '@fbcnms/babel-register/polyfill';

import Main from './components/Main';
import React from 'react';
import ReactDOM from 'react-dom';
import nullthrows from '@fbcnms/util/nullthrows';
import {BrowserRouter} from 'react-router-dom';

ReactDOM.render(
  <BrowserRouter>
    <Main />
  </BrowserRouter>,
  nullthrows(document.getElementById('root')),
);
