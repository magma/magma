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
import MomentUtils from '@date-io/moment';
import React from 'react';
import ReactDOM from 'react-dom';
import nullthrows from '@fbcnms/util/nullthrows';
import {BrowserRouter} from 'react-router-dom';
import {MuiPickersUtilsProvider} from '@material-ui/pickers';

import {} from './common/axiosConfig';

ReactDOM.render(
  <BrowserRouter>
    <MuiPickersUtilsProvider utils={MomentUtils}>
      <Main />
    </MuiPickersUtilsProvider>
  </BrowserRouter>,
  nullthrows(document.getElementById('root')),
);
