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

import {} from './common/axiosConfig';
import LoginForm from '../fbc_js_core/ui/components/auth/LoginForm.js';
import React from 'react';
import ReactDOM from 'react-dom';
import nullthrows from '../fbc_js_core/util/nullthrows';
import {AppContextProvider} from '../fbc_js_core/ui/context/AppContext';
import {BrowserRouter} from 'react-router-dom';
import {useHistory} from 'react-router-dom';

function LoginWrapper() {
  const history = useHistory();
  return (
    <LoginForm
      // eslint-disable-next-line no-warning-comments
      // $FlowFixMe - createHref exists
      action={history.createHref({pathname: '/user/login'})}
      title="Magma"
      // eslint-disable-next-line no-warning-comments
      // $FlowFixMe - createHref exists
      ssoAction={history.createHref({pathname: '/user/login/saml'})}
      ssoEnabled={window.CONFIG.appData.ssoEnabled}
      csrfToken={window.CONFIG.appData.csrfToken}
    />
  );
}

ReactDOM.render(
  <AppContextProvider>
    <BrowserRouter>
      <LoginWrapper />
    </BrowserRouter>
  </AppContextProvider>,
  nullthrows(document.getElementById('root')),
);
