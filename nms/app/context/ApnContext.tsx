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
import LoadingFiller from '../components/LoadingFiller';
import MagmaAPI from '../api/MagmaAPI';
import {Apn} from '../../generated';
import {NetworkId} from '../../shared/types/network';
import {useEffect, useState} from 'react';
import {useEnqueueSnackbar} from '../hooks/useSnackbar';

export type ApnContextType = {
  state: Record<string, Apn>;
  setState: (key: string, val?: Apn) => Promise<void>;
};
type ApnProvideProps = {
  networkId: NetworkId;
  children: React.ReactNode;
};

const ApnContext = React.createContext<ApnContextType>({} as ApnContextType);

async function setApnState(params: {
  networkId: NetworkId;
  apns: Record<string, Apn>;
  setApns: (apns: Record<string, Apn>) => void;
  key: string;
  value?: Apn;
}) {
  const {networkId, apns, setApns, key, value} = params;

  if (value != null) {
    if (!(key in apns)) {
      await MagmaAPI.apns.lteNetworkIdApnsPost({
        networkId: networkId,
        apn: value,
      });
      setApns({...apns, [key]: value});
    } else {
      await MagmaAPI.apns.lteNetworkIdApnsApnNamePut({
        networkId: networkId,
        apnName: key,
        apn: value,
      });
      setApns({...apns, [key]: value});
    }

    const apn = (
      await MagmaAPI.apns.lteNetworkIdApnsApnNameGet({
        networkId: networkId,
        apnName: key,
      })
    ).data;

    if (apn) {
      const newApns = {...apns, [key]: apn};
      setApns(newApns);
    }
  } else {
    await MagmaAPI.apns.lteNetworkIdApnsApnNameDelete({
      networkId: networkId,
      apnName: key,
    });
    const newApns = {...apns};
    delete newApns[key];
    setApns(newApns);
  }
}

export function ApnContextProvider(props: ApnProvideProps) {
  const {networkId} = props;
  const [apns, setApns] = useState<Record<string, Apn>>({});
  const [isLoading, setIsLoading] = useState(true);
  const enqueueSnackbar = useEnqueueSnackbar();

  useEffect(() => {
    const fetchState = async () => {
      try {
        setApns(
          (
            await MagmaAPI.apns.lteNetworkIdApnsGet({
              networkId,
            })
          ).data,
        );
      } catch (e) {
        enqueueSnackbar?.('failed fetching APN information', {
          variant: 'error',
        });
      }
      setIsLoading(false);
    };
    void fetchState();
  }, [networkId, enqueueSnackbar]);

  if (isLoading) {
    return <LoadingFiller />;
  }

  return (
    <ApnContext.Provider
      value={{
        state: apns,
        setState: (key, value?) => {
          return setApnState({
            apns,
            setApns,
            networkId,
            key,
            value,
          });
        },
      }}>
      {props.children}
    </ApnContext.Provider>
  );
}

export default ApnContext;
