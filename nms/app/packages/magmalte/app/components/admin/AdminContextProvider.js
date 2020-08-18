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

import * as React from 'react';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import {AppContextProvider} from '@fbcnms/ui/context/AppContext';

import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';

export default function AdminContextProvider(props: {children: React.Node}) {
  const {error, isLoading, response} = useMagmaAPI(MagmaV1API.getNetworks, {});

  if (isLoading) {
    return <LoadingFiller />;
  }

  const networkIds = error || !response ? ['mpk_test'] : response.sort();

  return (
    <AppContextProvider networkIDs={networkIds}>
      {props.children}
    </AppContextProvider>
  );
}
