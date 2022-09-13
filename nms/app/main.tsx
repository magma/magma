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
 */

import './util/axiosConfig';
import './util/chartjsSetup';
import './util/polyfill';

import ApplicationMain from './components/ApplicationMain';
import Main from './components/Main';
import React from 'react';
import ReactDOM from 'react-dom';
import nullthrows from '../shared/util/nullthrows';
import {AdapterDateFns} from '@mui/x-date-pickers/AdapterDateFns';
import {BrowserRouter} from 'react-router-dom';
import {LocalizationProvider} from '@mui/x-date-pickers';

ReactDOM.render(
  <BrowserRouter>
    <LocalizationProvider dateAdapter={AdapterDateFns}>
      <ApplicationMain>
        <Main />
      </ApplicationMain>
    </LocalizationProvider>
  </BrowserRouter>,
  nullthrows(document.getElementById('root')),
);
