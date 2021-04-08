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
import {FetchEnodebs, FetchGateways} from '../../state/lte/EquipmentState';
import {
  FetchSubscriberState,
  FetchSubscribers,
} from '../../state/lte/SubscriberState';
import {useContext, useEffect, useRef, useState} from 'react';

export const REFRESH_INTERVAL = 30000;

type Props = {
  context: typeof React.Context,
  networkId: string,
  type: refreshType,
  interval?: number,
  id?: string,
  enqueueSnackbar?: (msg: string, cfg: {}) => ?(string | number),
  refresh: boolean,
  lastRefreshTime?: string,
};

type refreshType = 'subscriber' | 'gateway' | 'enodeb';

export function useRefreshingContext(props: Props) {
  const ctx = useContext(props.context);
  const [state, setState] = useState(
    props.type === 'subscriber'
      ? {state: ctx.state, sessionState: ctx.sessionState}
      : ctx.state,
  );

  const [autoRefreshTime, setAutoRefreshTime] = useState(props.lastRefreshTime);
  async function fetchState(props: FetchProps) {
    const newState = await fetchRefreshState({
      type: props.type,
      networkId: props.networkId,
      id: props.id,
      enqueueSnackbar: props.enqueueSnackbar,
    });
    if (newState) {
      setState(() => {
        if (props.type === 'subscriber') {
          return {
            sessionState: newState?.sessionState || {},
            state: newState?.state || {},
          };
        } else {
          return newState;
        }
      });
    }
  }

  function updateContext(id?: string, state) {
    let newState = state;
    if (id !== null && id !== undefined) {
      if (props.type === 'subscriber') {
        newState = {
          state: {
            ...ctx.state,
            // $FlowIgnore
            [id]: state.state?.[id],
          },
          // $FlowIgnore
          sessionState: Object.keys(state.sessionState || {}).length
            ? {
                ...ctx.sessionState,
                // $FlowIgnore
                [id]: state.sessionState?.[id],
              }
            : {},
        };
      } else if (props.type === 'enodeb') {
        // $FlowIgnore
        newState = {...ctx.state, [id]: state.enbInfo?.[id]};
      } else {
        newState = {...ctx.state, [id]: state?.[id]};
      }
    }
    return ctx.setState(null, null, newState);
  }

  // Avoid using state as a dependency of useEffect
  const stateRef = useRef(null);
  useEffect(() => {
    stateRef.current = state;
  }, [state]);

  useEffect(() => {
    const intervalId = setInterval(
      () => setAutoRefreshTime(new Date().toLocaleString()),
      props.interval,
    );
    if (!props.refresh) {
      return clearInterval(intervalId);
    }
    return () => {
      updateContext(props.id, stateRef.current);
      clearInterval(intervalId);
    };
    // eslint-disable-next-line
  }, [props.interval, props.refresh]);

  useEffect(() => {
    if (props.lastRefreshTime != autoRefreshTime) {
      fetchState({
        type: props.type,
        networkId: props.networkId,
        id: props.id,
        enqueueSnackbar: props.enqueueSnackbar,
      });
    }
  }, [
    props.type,
    props.networkId,
    props.enqueueSnackbar,
    props.id,
    props.lastRefreshTime,
    autoRefreshTime,
  ]);
  return state;
}

type FetchProps = {
  type: refreshType,
  networkId: string,
  id?: string,
  enqueueSnackbar?: (msg: string, cfg: {}) => ?(string | number),
};
async function fetchRefreshState(props: FetchProps) {
  const {type, networkId, id, enqueueSnackbar} = props;
  if (type === 'subscriber') {
    const subscribers = await FetchSubscribers({
      id: id,
      networkId,
      enqueueSnackbar,
    });
    const sessions = await FetchSubscriberState({
      id: id,
      networkId,
      enqueueSnackbar,
    });
    if (id !== null && id !== undefined) {
      return {
        sessionState: {[id]: sessions || {}},
        state: {[id]: subscribers || {}},
      };
    }
    return {sessionState: sessions, state: subscribers};
  } else if (type === 'gateway') {
    const gateways = await FetchGateways({
      id: id,
      networkId,
      enqueueSnackbar,
    });

    return gateways;
  } else {
    const enodebs = await FetchEnodebs({
      id: id,
      networkId,
      enqueueSnackbar,
    });

    return {enbInfo: enodebs};
  }
}
