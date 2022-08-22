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
import {CallTrace, CallTraceConfig, MutableCallTrace} from '../../generated';
import {EnqueueSnackbar, useEnqueueSnackbar} from '../hooks/useSnackbar';
import {NetworkId} from '../../shared/types/network';
import {useEffect, useState} from 'react';

type TraceContextType = {
  state: Record<string, CallTrace>;
  setState?: (key: string, val?: MutableCallTrace) => Promise<void>;
};
type ContextProviderProps = {
  networkId: NetworkId;
  children: React.ReactNode;
};

const TraceContext = React.createContext<TraceContextType>(
  {} as TraceContextType,
);

async function initTraceState(params: {
  networkId: NetworkId;
  setTraceMap: (traceMap: Record<string, CallTrace>) => void;
  enqueueSnackbar?: EnqueueSnackbar;
}) {
  const {networkId, setTraceMap, enqueueSnackbar} = params;

  try {
    const traces = (
      await MagmaAPI.callTracing.networksNetworkIdTracingGet({
        networkId,
      })
    ).data;

    if (traces) {
      setTraceMap(traces);
    }
  } catch (e) {
    enqueueSnackbar?.('failed fetching call trace information', {
      variant: 'error',
    });
    return;
  }
}

/** setCallTraceState
 * if key and value are passed in,
 * if key is not present, a new trace is created (POST)
 * if key is present, existing trace is updated (PUT)
 * if value is not present, the trace is deleted (DELETE)
 */
async function setCallTraceState(params: {
  networkId: NetworkId;
  callTraces: Record<string, CallTrace>;
  setCallTraces: (callTraces: Record<string, CallTrace>) => void;
  key: string;
  value?: Partial<CallTraceConfig & MutableCallTrace>;
}) {
  const {networkId, callTraces, setCallTraces, key, value} = params;

  if (value != null) {
    if (!(key in callTraces)) {
      await MagmaAPI.callTracing.networksNetworkIdTracingPost({
        networkId: networkId,
        callTraceConfiguration: value as CallTraceConfig,
      });
    } else {
      await MagmaAPI.callTracing.networksNetworkIdTracingTraceIdPut({
        networkId: networkId,
        traceId: key,
        callTraceConfiguration: value as MutableCallTrace,
      });
    }

    const callTrace = (
      await MagmaAPI.callTracing.networksNetworkIdTracingTraceIdGet({
        networkId: networkId,
        traceId: key,
      })
    ).data;

    if (callTrace) {
      const newTraces = {...callTraces, [key]: callTrace};
      setCallTraces(newTraces);
    }
  } else {
    await MagmaAPI.callTracing.networksNetworkIdTracingTraceIdDelete({
      networkId: networkId,
      traceId: key,
    });
    const newCallTraces = {...callTraces};
    delete newCallTraces[key];
    setCallTraces(newCallTraces);
  }
}

export function TraceContextProvider(props: ContextProviderProps) {
  const {networkId} = props;
  const [traceMap, setTraceMap] = useState<Record<string, CallTrace>>({});
  const [isLoading, setIsLoading] = useState(true);
  const enqueueSnackbar = useEnqueueSnackbar();

  useEffect(() => {
    const fetchLteState = async () => {
      if (networkId == null) {
        return;
      }
      await initTraceState({
        networkId,
        setTraceMap,
        enqueueSnackbar,
      });
      setIsLoading(false);
    };
    void fetchLteState();
  }, [networkId, enqueueSnackbar]);

  if (isLoading) {
    return <LoadingFiller />;
  }

  return (
    <TraceContext.Provider
      value={{
        state: traceMap,
        setState: (key: string, value?: MutableCallTrace | CallTraceConfig) =>
          setCallTraceState({
            networkId,
            callTraces: traceMap,
            setCallTraces: setTraceMap,
            key,
            value,
          }),
      }}>
      {props.children}
    </TraceContext.Provider>
  );
}

export default TraceContext;
