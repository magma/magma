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
import LoginForm from '@fbcnms/ui/components/auth/LoginForm.js';
import React from 'react';
import ReactDOM from 'react-dom';
import nullthrows from '@fbcnms/util/nullthrows';
import {AppContextProvider} from '@fbcnms/ui/context/AppContext';
import {BrowserRouter} from 'react-router-dom';

import {} from './common/axiosConfig';
import {useRouter} from '@fbcnms/ui/hooks';

function LoginWrapper() {
  const {history} = useRouter();
  return (
    <LoginForm
      action={history.createHref({pathname: '/user/login'})}
      title="Magma"
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
