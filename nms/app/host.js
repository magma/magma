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

import '../fbc_js_core/babel_register/polyfill';

import Index from './components/host/Index';
import React from 'react';
import ReactDOM from 'react-dom';
import {BrowserRouter} from 'react-router-dom';

import {} from './common/axiosConfig';
import nullthrows from '../fbc_js_core/util/nullthrows';

ReactDOM.render(
  <BrowserRouter>
    <Index />
  </BrowserRouter>,
  nullthrows(document.getElementById('root')),
);
