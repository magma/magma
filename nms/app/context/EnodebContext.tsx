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
import {
  Enodeb,
  EnodebState as EnodebStateResponse,
  NetworkRanConfigs,
} from '../../generated';
import {EnodebInfo} from '../components/lte/EnodebUtils';
import {EnqueueSnackbar, useEnqueueSnackbar} from '../hooks/useSnackbar';
import {NetworkId} from '../../shared/types/network';
import {useEffect, useState} from 'react';

type EnodebState = {
  enbInfo: Record<string, EnodebInfo>;
};
export type EnodebContextType = {
  state: EnodebState;
  lteRanConfigs?: NetworkRanConfigs;
  setState: (key: string, val?: EnodebInfo) => Promise<void>;
  refetch: (id?: string) => void;
};
type EnodebContextProps = {
  networkId: NetworkId;
  children: React.ReactNode;
};

const EnodebContext = React.createContext<EnodebContextType>(
  {} as EnodebContextType,
);

async function fetchEnodebs(params: {
  networkId: string;
  id?: string;
  enqueueSnackbar?: EnqueueSnackbar;
}): Promise<Record<string, EnodebInfo> | undefined> {
  const {networkId, id} = params;
  if (id !== undefined && id !== null) {
    let enb: Enodeb;
    try {
      enb = (
        await MagmaAPI.enodebs.lteNetworkIdEnodebsEnodebSerialGet({
          networkId: networkId,
          enodebSerial: id,
        })
      ).data;
      if (enb) {
        const newEnbSt = (
          await MagmaAPI.enodebs.lteNetworkIdEnodebsEnodebSerialStateGet({
            networkId: networkId,
            enodebSerial: id,
          })
        ).data;
        const newEnb = {
          [id]: {
            enb_state: newEnbSt,
            enb: enb,
          },
        };
        return newEnb;
      }
    } catch (e) {
      return {
        [id]: {
          enb_state: {},
          enb: enb!,
        } as EnodebInfo,
      };
    }
  } else {
    const resp = (
      await MagmaAPI.enodebs.lteNetworkIdEnodebsGet({
        networkId,
      })
    ).data;
    const enbs = resp.enodebs;
    if (!enbs) {
      return;
    }

    const requests = Object.keys(enbs).map(async k => {
      try {
        const {serial} = enbs[k];
        // eslint-disable-next-line max-len
        const enbSt = (
          await MagmaAPI.enodebs.lteNetworkIdEnodebsEnodebSerialStateGet({
            networkId: networkId,
            enodebSerial: serial,
          })
        ).data;
        return [enbs[k], enbSt ?? {}] as const;
      } catch (e) {
        return [enbs[k], {}] as const;
      }
    });

    const enbResp = await Promise.all(requests);
    const enbInfo: Record<string, EnodebInfo> = {};
    enbResp.filter(Boolean).forEach(r => {
      if (r.length > 0) {
        const [enb, enbSt] = r;
        if (enb != null && enbSt != null) {
          enbInfo[enb.serial] = {
            enb,
            enb_state: enbSt as EnodebStateResponse,
          };
        }
      }
    });
    return enbInfo;
  }
}

async function initEnodeState(params: {
  networkId: NetworkId;
  setEnbInfo: (enodebInfo: Record<string, EnodebInfo>) => void;
  enqueueSnackbar?: EnqueueSnackbar;
}) {
  const enodebInfo = await fetchEnodebs({
    networkId: params.networkId,
    enqueueSnackbar: params.enqueueSnackbar,
  });

  if (enodebInfo) {
    params.setEnbInfo(enodebInfo);
  }
}

async function setEnodebState(params: {
  networkId: NetworkId;
  enbInfo: Record<string, EnodebInfo>;
  setEnbInfo: (enodebInfo: Record<string, EnodebInfo>) => void;
  key: string;
  value?: EnodebInfo;
}) {
  const {networkId, enbInfo, setEnbInfo, key, value} = params;

  if (value != null) {
    // remove attached gateway id read only property
    if (value.enb.hasOwnProperty('attached_gateway_id')) {
      delete value.enb['attached_gateway_id'];
    }

    if (!(key in enbInfo)) {
      await MagmaAPI.enodebs.lteNetworkIdEnodebsPost({
        networkId: networkId,
        enodeb: value.enb,
      });
      setEnbInfo({...enbInfo, [key]: value});
    } else {
      await MagmaAPI.enodebs.lteNetworkIdEnodebsEnodebSerialPut({
        networkId: networkId,
        enodebSerial: key,
        enodeb: value.enb,
      });
      const prevEnbSt = enbInfo[key].enb_state;
      setEnbInfo({
        ...enbInfo,
        [key]: {
          enb_state: prevEnbSt,
          enb: value.enb,
        },
      });
    }
  } else {
    await MagmaAPI.enodebs.lteNetworkIdEnodebsEnodebSerialDelete({
      networkId: networkId,
      enodebSerial: key,
    });
    const newEnbInfo = {...enbInfo};
    delete newEnbInfo[key];
    setEnbInfo(newEnbInfo);
    return;
  }
}

export function EnodebContextProvider(props: EnodebContextProps) {
  const {networkId} = props;
  const [enbInfo, setEnbInfo] = useState<Record<string, EnodebInfo>>({});
  const [lteRanConfigs, setLteRanConfigs] = useState<NetworkRanConfigs>(
    {} as NetworkRanConfigs,
  );
  const [isLoading, setIsLoading] = useState(true);
  const enqueueSnackbar = useEnqueueSnackbar();
  useEffect(() => {
    const fetchState = async () => {
      try {
        if (networkId == null) {
          return;
        }
        const [lteRanConfigsResp] = await Promise.allSettled([
          MagmaAPI.lteNetworks.lteNetworkIdCellularRanGet({networkId}),
          initEnodeState({networkId, setEnbInfo, enqueueSnackbar}),
          Promise.reject('Bla'),
        ]);
        if (lteRanConfigsResp.status === 'fulfilled') {
          setLteRanConfigs(lteRanConfigsResp.value.data);
        }
      } catch (e) {
        enqueueSnackbar?.('failed fetching enodeb information', {
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
    <EnodebContext.Provider
      value={{
        state: {enbInfo},
        lteRanConfigs: lteRanConfigs,
        setState: (key: string, value?) => {
          return setEnodebState({
            enbInfo,
            setEnbInfo,
            networkId,
            key,
            value,
          });
        },
        refetch: id => {
          void fetchEnodebs({
            id: id,
            networkId,
            enqueueSnackbar,
          }).then(enodebs => {
            if (enodebs) {
              setEnbInfo(enodebState =>
                id ? {...enodebState, ...enodebs} : enodebs,
              );
            }
          });
        },
      }}>
      {props.children}
    </EnodebContext.Provider>
  );
}

export default EnodebContext;
