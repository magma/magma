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
// $FlowFixMe migrated to typescript
import {FetchEnodebs, FetchGateways} from '../../state/lte/EquipmentState';
import {
  FetchFegGateways,
  getActiveFegGatewayId,
  getFegGatewaysHealthStatus,
  // $FlowFixMe migrated to typescript
} from '../../state/feg/EquipmentState';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {FetchFegSubscriberState} from '../../state/feg/SubscriberState';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {FetchSubscriberState} from '../../state/lte/SubscriberState';
import {useContext, useEffect, useRef, useState} from 'react';
import type {EnqueueSnackbarOptions} from 'notistack';

export const REFRESH_INTERVAL = 30000;

export const RefreshTypeEnum = Object.freeze({
  SUBSCRIBER: 'subscriber',
  GATEWAY: 'gateway',
  FEG_GATEWAY: 'feg_gateway',
  FEG_SUBSCRIBER: 'feg_subscriber',
  ENODEB: 'enodeb',
});

type Props = {
  context: typeof React.Context,
  networkId: string,
  type: RefreshType,
  interval?: number,
  id?: string,
  enqueueSnackbar?: (
    msg: string,
    cfg: EnqueueSnackbarOptions,
  ) => ?(string | number),
  refresh: boolean,
  lastRefreshTime?: string,
};

type RefreshType = $Values<typeof RefreshTypeEnum>;

export function useRefreshingContext(props: Props) {
  const ctx = useContext(props.context);
  const initState = () => {
    if (props.type === RefreshTypeEnum.FEG_GATEWAY) {
      return {
        fegGateways: ctx?.state,
        health: ctx?.health,
        activeFegGatewayId: ctx?.activeFegGatewayId,
      };
    } else if (
      props.type === RefreshTypeEnum.SUBSCRIBER ||
      props.type === RefreshTypeEnum.FEG_SUBSCRIBER
    ) {
      return {sessionState: ctx.sessionState};
    }
    return ctx.state;
  };
  const [state, setState] = useState(initState());

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
          // $FlowIgnore because state may not contain sessionState for other refresh type like enodeb
          sessionState: Object.keys(state.sessionState || {}).length
            ? {
                ...ctx.sessionState,
                // $FlowIgnore because state may not contain sessionState for other refresh type like enodeb
                [id]: state.sessionState?.[id],
              }
            : {},
        };
      } else if (props.type === 'enodeb') {
        // $FlowIgnore because state may not contain enbInfo for other refresh type like feg_gateway
        newState = {...ctx.state, [id]: state.enbInfo?.[id]};
      } else if (props.type === RefreshTypeEnum.FEG_GATEWAY) {
        newState = {
          // $FlowIgnore because state may not contain fegGateways for other refresh type like subscriber
          fegGateways: {...ctx.fegGateways, [id]: state?.fegGateways?.[id]},
        };
      } else if (props.type === RefreshTypeEnum.FEG_SUBSCRIBER) {
        newState = {...ctx.sessionState};
      } else {
        newState = {...ctx.state, [id]: state?.[id]};
      }
    }
    if (props.type === 'subscriber') {
      // update subscriber session state
      return ctx.setState(null, null, null, newState);
    } else if (props.type === RefreshTypeEnum.FEG_GATEWAY) {
      return ctx.setState(null, null, newState?.fegGateways || {});
    } else if (props.type === RefreshTypeEnum.FEG_SUBSCRIBER) {
      return ctx.setSessionState(newState);
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
  type: RefreshType,
  networkId: string,
  id?: string,
  enqueueSnackbar?: (
    msg: string,
    cfg: EnqueueSnackbarOptions,
  ) => ?(string | number),
};
async function fetchRefreshState(props: FetchProps) {
  const {type, networkId, id, enqueueSnackbar} = props;
  if (type === 'subscriber') {
    const sessions = await FetchSubscriberState({
      id: id,
      networkId,
      enqueueSnackbar,
    });
    if (id !== null && id !== undefined) {
      return {
        sessionState: {[id]: sessions || {}},
      };
    }
    return {sessionState: sessions};
  } else if (type === 'gateway') {
    const gateways = await FetchGateways({
      id: id,
      networkId,
      enqueueSnackbar,
    });

    return gateways;
  } else if (type === RefreshTypeEnum.FEG_GATEWAY) {
    const fegGateways = await FetchFegGateways({
      id: id,
      networkId,
      enqueueSnackbar,
    });
    const [health, activeFegGatewayId] = await Promise.all([
      getFegGatewaysHealthStatus(networkId, fegGateways, enqueueSnackbar),
      getActiveFegGatewayId(networkId, fegGateways, enqueueSnackbar),
    ]);
    return {fegGateways, health, activeFegGatewayId};
  } else if (type === RefreshTypeEnum.FEG_SUBSCRIBER) {
    const sessions = await FetchFegSubscriberState({
      networkId,
      enqueueSnackbar,
    });
    return {sessionState: sessions};
  } else {
    const enodebs = await FetchEnodebs({
      id: id,
      networkId,
      enqueueSnackbar,
    });

    return {enbInfo: enodebs};
  }
}
