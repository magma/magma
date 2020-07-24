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

import * as React from 'react';
import * as imm from 'immutable';
import emptyFunction from '@fbcnms/util/emptyFunction';
import {useCallback, useState} from 'react';

type WizardContextType = {
  set: (id: string, data: ?mixed) => void,
  get: (id: string) => ?mixed,
};

const WizardContext = React.createContext<WizardContextType>({
  set: emptyFunction,
  get: emptyFunction,
});

type Props = {
  children: React.Node,
};

type DataStore = imm.Map<string, mixed>;

export function WizardContextProvider(props: Props) {
  const [dataStore, setDataStore] = useState<DataStore>(
    new imm.Map<string, mixed>(),
  );

  const setData = useCallback(
    (id, data: ?mixed) => {
      setDataStore(dataStore.set(id, data));
    },
    [dataStore],
  );
  const getData = useCallback(
    id => (dataStore.has(id) && dataStore.get(id)) || undefined,
    [dataStore],
  );

  const providerValue = {
    set: setData,
    get: getData,
  };

  return (
    <WizardContext.Provider value={providerValue}>
      {props.children}
    </WizardContext.Provider>
  );
}

export default WizardContext;
