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
 */

import MagmaAPI from '../../api/MagmaAPI';
import type {
  CallTrace,
  CallTraceConfig,
  MutableCallTrace,
} from '../../generated-ts';
import type {NetworkId} from '../../shared/types/network';
import type {OptionsObject} from 'notistack';

type InitTraceStateProps = {
  networkId: NetworkId;
  setTraceMap: (arg0: Record<string, CallTrace>) => void;
  enqueueSnackbar?: (
    msg: string,
    cfg: OptionsObject,
  ) => (string | number) | null | undefined;
};

export async function InitTraceState(props: InitTraceStateProps) {
  const {networkId, setTraceMap, enqueueSnackbar} = props;

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

type CallTraceProps = {
  networkId: NetworkId;
  callTraces: Record<string, CallTrace>;
  setCallTraces: (arg0: Record<string, CallTrace>) => void;
  key: string;
  value?: Partial<CallTraceConfig & MutableCallTrace>;
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
