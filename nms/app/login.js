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
 * @flow strict-local
 * @format
 */
'use strict';

import './common/polyfill';

import {} from './common/axiosConfig';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import LoginForm from './views/login/LoginForm.tsx';
import React from 'react';
import ReactDOM from 'react-dom';
import ThemeProvider from '@material-ui/styles/ThemeProvider';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import defaultTheme from './theme/default';
// $FlowFixMe migrated to typescript
import nullthrows from '../shared/util/nullthrows';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {AppContextProvider} from './components/context/AppContext';
import {BrowserRouter} from 'react-router-dom';

function LoginWrapper() {
  return (
    <LoginForm
      // eslint-disable-next-line no-warning-comments
      // $FlowFixMe - createHref exists
      action="/user/login"
      title="Magma"
      // eslint-disable-next-line no-warning-comments
      // $FlowFixMe - createHref exists
      ssoAction="/user/login/saml"
      ssoEnabled={window.CONFIG.appData.ssoEnabled}
      csrfToken={window.CONFIG.appData.csrfToken}
    />
  );
}

ReactDOM.render(
  <AppContextProvider>
    <BrowserRouter>
      <ThemeProvider theme={defaultTheme}>
        <LoginWrapper />
      </ThemeProvider>
    </BrowserRouter>
  </AppContextProvider>,
  nullthrows(document.getElementById('root')),
);
