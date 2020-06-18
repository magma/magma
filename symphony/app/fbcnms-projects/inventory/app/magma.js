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

import {init} from 'fbt';

import '@fbcnms/babel-register/polyfill';

import Main from '@fbcnms/magmalte/app/components/Main';
import React from 'react';
import ReactDOM from 'react-dom';
import nullthrows from '@fbcnms/util/nullthrows';
import translatedFbts from '../i18n/translatedFbts.json';
import {BrowserRouter} from 'react-router-dom';
import {hot} from 'react-hot-loader';

import {} from './common/axiosConfig';

init({translations: translatedFbts});

/* eslint-disable-next-line no-undef */
const HotMain = hot(module)(Main);

ReactDOM.render(
  <BrowserRouter>
    <HotMain />
  </BrowserRouter>,
  nullthrows(document.getElementById('root')),
);
