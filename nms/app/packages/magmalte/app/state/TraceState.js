/*
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

import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import type {EnqueueSnackbarOptions} from 'notistack';

import type {
  call_trace,
  call_trace_config,
  mutable_call_trace,
  network_id,
} from '@fbcnms/magma-api';

type InitTraceStateProps = {
  networkId: network_id,
  setTraceMap: ({[string]: call_trace}) => void,
  enqueueSnackbar?: (
    msg: string,
    cfg: EnqueueSnackbarOptions,
  ) => ?(string | number),
};

export async function InitTraceState(props: InitTraceStateProps) {
  const {networkId, setTraceMap, enqueueSnackbar} = props;
  try {
    const traces = await MagmaV1API.getNetworksByNetworkIdTracing({
      networkId,
    });
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

type CallTraceProps = {
  networkId: network_id,
  callTraces: {[string]: call_trace},
  setCallTraces: ({[string]: call_trace}) => void,
  key: string,
  value?: $Shape<call_trace_config & mutable_call_trace>,
};

/* SetCallTraceState
SetCallTraceState
if key and value are passed in,
if key is not present, a new trace is created (POST)
if key is present, existing trace is updated (PUT)
if value is not present, the trace is deleted (DELETE)
*/
export async function SetCallTraceState(props: CallTraceProps) {
  const {networkId, callTraces, setCallTraces, key, value} = props;
  if (value != null) {
    if (!(key in callTraces)) {
      await MagmaV1API.postNetworksByNetworkIdTracing({
        networkId: networkId,
        callTraceConfiguration: value, // call_trace_config
      });
    } else {
      await MagmaV1API.putNetworksByNetworkIdTracingByTraceId({
        networkId: networkId,
        traceId: key,
        callTraceConfiguration: value, // mutable_call_trace
      });
    }
    const callTrace = await MagmaV1API.getNetworksByNetworkIdTracingByTraceId({
      networkId: networkId,
      traceId: key,
    });

    if (callTrace) {
      const newTraces = {...callTraces, [key]: callTrace};
      setCallTraces(newTraces);
    }
  } else {
    await MagmaV1API.deleteNetworksByNetworkIdTracingByTraceId({
      networkId: networkId,
      traceId: key,
    });
    const newCallTraces = {...callTraces};
    delete newCallTraces[key];
    setCallTraces(newCallTraces);
  }
}
