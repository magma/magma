/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import '@fbcnms/babel-register/polyfill';

import Index from '@fbcnms/magmalte/app/components/master/Index';
import React from 'react';
import ReactDOM from 'react-dom';
import {BrowserRouter} from 'react-router-dom';
import {hot} from 'react-hot-loader';

import {} from './common/axiosConfig';
import nullthrows from '@fbcnms/util/nullthrows';

/* eslint-disable-next-line no-undef */
const HotIndex = hot(module)(Index);

ReactDOM.render(
  <BrowserRouter>
    <HotIndex />
  </BrowserRouter>,
  nullthrows(document.getElementById('root')),
);
