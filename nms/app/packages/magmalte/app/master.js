/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
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
