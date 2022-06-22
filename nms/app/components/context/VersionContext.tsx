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

import * as React from 'react';
import MagmaAPI from '../../../api/MagmaAPI';
import axios from 'axios';
import {useEffect, useState} from 'react';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';

export type VersionContextType = {
  nmsVersion: string;
  orc8rVersion: string;
};

const VersionContext = React.createContext<VersionContextType>({
  nmsVersion: 'vNMS',
  orc8rVersion: 'vORC8R',
});

type Props = {
  children: React.ReactNode;
};

export function VersionContextProvider(props: Props) {
  const enqueueSnackbar = useEnqueueSnackbar();
  const [nmsVersion, setNmsVersion] = useState<string>('vNMS');
  const [orc8rVersion, setOrc8rVersion] = useState<string>('vORC8R');

  useEffect(() => {
    const fetchNmsVersion = async () => {
      axios
        .get<string>('/version')
        .then(response => {
          setNmsVersion(`v${response.data}`);
        })
        .catch(() => {
          enqueueSnackbar?.('failed fetching NMS version information', {
            variant: 'error',
          });
        });

      const version = (await MagmaAPI.about.aboutVersionGet()).data;
      setOrc8rVersion(`v${version.container_image_version!}`);
    };

    void fetchNmsVersion();
  });

  return (
    <VersionContext.Provider
      value={{
        nmsVersion: nmsVersion,
        orc8rVersion: orc8rVersion,
      }}>
      {props.children}
    </VersionContext.Provider>
  );
}

export default VersionContext;
