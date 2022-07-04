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

import './common/axiosConfig';
import './common/polyfill';

import ApplicationMain from './components/ApplicationMain';
import Main from './components/Main';
import MomentUtils from '@date-io/moment';
import React from 'react';
import ReactDOM from 'react-dom';
import nullthrows from '../shared/util/nullthrows';
import {BrowserRouter} from 'react-router-dom';
import {MuiPickersUtilsProvider} from '@material-ui/pickers';

ReactDOM.render(
  <BrowserRouter>
    <MuiPickersUtilsProvider utils={MomentUtils}>
      <ApplicationMain>
        <Main />
      </ApplicationMain>
    </MuiPickersUtilsProvider>
  </BrowserRouter>,
  nullthrows(document.getElementById('root')),
);
