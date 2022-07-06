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
import EnodebContext, {EnodebContextType} from './EnodebContext';
import FEGGatewayContext, {FEGGatewayContextType} from './FEGGatewayContext';
import FEGSubscriberContext, {
  FEGSubscriberContextType,
} from './FEGSubscriberContext';
import GatewayContext, {GatewayContextType} from './GatewayContext';
import SubscriberContext, {SubscriberContextType} from './SubscriberContext';
import {FederationGateway} from '../../../generated-ts';
import {FetchEnodebs, FetchGateways} from '../../state/lte/EquipmentState';
import {
  FetchFegGateways,
  getActiveFegGatewayId,
  getFegGatewaysHealthStatus,
} from '../../state/feg/EquipmentState';
import {FetchFegSubscriberState} from '../../state/feg/SubscriberState';
import {FetchSubscriberState} from '../../state/lte/SubscriberState';
import {GatewayId} from '../../../shared/types/network';
import {useContext, useEffect, useRef, useState} from 'react';
import type {OptionsObject} from 'notistack';

export const REFRESH_INTERVAL = 30000;

export const RefreshTypeEnum = {
  SUBSCRIBER: 'subscriber',
  GATEWAY: 'gateway',
  FEG_GATEWAY: 'feg_gateway',
  FEG_SUBSCRIBER: 'feg_subscriber',
  ENODEB: 'enodeb',
} as const;

type ContextMap = {
  [RefreshTypeEnum.SUBSCRIBER]: typeof SubscriberContext;
  [RefreshTypeEnum.GATEWAY]: typeof GatewayContext;
  [RefreshTypeEnum.FEG_GATEWAY]: typeof FEGGatewayContext;
  [RefreshTypeEnum.FEG_SUBSCRIBER]: typeof FEGSubscriberContext;
  [RefreshTypeEnum.ENODEB]: typeof EnodebContext;
};

type StateMap = {
  [RefreshTypeEnum.SUBSCRIBER]: {
    sessionState: SubscriberContextType['sessionState'];
  };
  [RefreshTypeEnum.GATEWAY]: GatewayContextType['state'];
  [RefreshTypeEnum.FEG_GATEWAY]: {
    fegGateways: FEGGatewayContextType['state'];
    health: FEGGatewayContextType['health'];
    activeFegGatewayId: FEGGatewayContextType['activeFegGatewayId'];
  };
  [RefreshTypeEnum.FEG_SUBSCRIBER]: {
    sessionState: FEGSubscriberContextType['sessionState'];
  };
  [RefreshTypeEnum.ENODEB]: EnodebContextType['state'];
};

type RefreshType = keyof ContextMap;

type Props<T extends RefreshType> = {
  type: T;
  context: ContextMap[T];
  networkId: string;
  interval?: number;
  id?: string;
  enqueueSnackbar?: (
    msg: string,
    cfg: OptionsObject,
  ) => string | number | undefined | null;
  refresh: boolean;
  lastRefreshTime?: string;
};

// TODO: This hook is not well designed and typing it correctly is nearly impossible,
//  it should be replaced with a simpler solution.
/* eslint-disable @typescript-eslint/no-unsafe-call, @typescript-eslint/no-unsafe-argument, @typescript-eslint/no-unsafe-assignment, @typescript-eslint/no-unsafe-member-access, @typescript-eslint/no-unsafe-return */
export function useRefreshingContext<T extends keyof ContextMap>(
  props: Props<T>,
): StateMap[T] {
  const ctx: any = useContext(props.context as any);
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

  function updateContext(id: string | undefined, state: any) {
    let newState = state;
    if (id !== null && id !== undefined) {
      if (props.type === 'subscriber') {
        newState = {
          sessionState: Object.keys(state.sessionState || {}).length
            ? {
                ...ctx.sessionState,
                [id]: state.sessionState?.[id],
              }
            : {},
        };
      } else if (props.type === 'enodeb') {
        newState = {...ctx.state, [id]: state.enbInfo?.[id]};
      } else if (props.type === RefreshTypeEnum.FEG_GATEWAY) {
        newState = {
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
      void fetchState({
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
/* eslint-enable @typescript-eslint/no-unsafe-call, @typescript-eslint/no-unsafe-argument, @typescript-eslint/no-unsafe-assignment, @typescript-eslint/no-unsafe-member-access, @typescript-eslint/no-unsafe-return */

type FetchProps = {
  type: RefreshType;
  networkId: string;
  id?: string;
  enqueueSnackbar?: (
    msg: string,
    cfg: OptionsObject,
  ) => string | number | undefined | null;
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
    const fegGateways = (await FetchFegGateways({
      id: id,
      networkId,
      enqueueSnackbar,
    })) as Record<GatewayId, FederationGateway>;

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
