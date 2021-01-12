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
import {useContext, useEffect, useState} from 'react';

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
  const [state, setState] = useState(ctx);

  async function fetchState(props: FetchProps) {
    const newState = await fetchRefreshState({
      context: props.context,
      type: props.type,
      networkId: props.networkId,
      id: props.id,
      enqueueSnackbar: props.enqueueSnackbar,
    });
    if (newState) {
      setState(prevState => {
        if (props.type === 'subscriber') {
          return {
            ...prevState,
            sessionState: newState?.sessionState || {},
            state: newState?.state || {},
          };
        } else {
          return {
            ...prevState,
            state: newState,
          };
        }
      });
    }
  }

  useEffect(() => {
    fetchState({
      context: ctx,
      type: props.type,
      networkId: props.networkId,
      id: props.id,
      enqueueSnackbar: props.enqueueSnackbar,
    });
    const intervalId = setInterval(
      () =>
        fetchState({
          context: ctx,
          type: props.type,
          networkId: props.networkId,
          id: props.id,
          enqueueSnackbar: props.enqueueSnackbar,
        }),
      props.interval,
    );
    if (!props.refresh) {
      return clearInterval(intervalId);
    }
    return () => {
      clearInterval(intervalId);
    };
  }, [
    ctx,
    props.type,
    props.networkId,
    props.enqueueSnackbar,
    props.id,
    props.interval,
    props.refresh,
  ]);
  return state;
}

type FetchProps = {
  context: typeof React.Context,
  type: refreshType,
  networkId: string,
  id?: string,
  enqueueSnackbar?: (msg: string, cfg: {}) => ?(string | number),
};
async function fetchRefreshState(props: FetchProps) {
  const {type, networkId, id, enqueueSnackbar, context} = props;
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
    if (id?.length) {
      return {
        sessionState: Object.keys(sessions).length
          ? {...context.sessionState, [id]: sessions}
          : context.sessionState,
        state: Object.keys(subscribers).length
          ? {...context.state, [id]: subscribers}
          : context.state,
      };
    }
    return {sessionState: sessions, state: subscribers};
  } else if (type === 'gateway') {
    const gateways = await FetchGateways({
      id: id,
      networkId,
      enqueueSnackbar: enqueueSnackbar,
    });

    if (id?.length) {
      return {...context.state, [id]: gateways};
    }
    return gateways;
  } else {
    const enodebs = await FetchEnodebs({
      id: id,
      networkId,
      enqueueSnackbar: enqueueSnackbar,
    });

    if (id?.length) {
      return {enbInfo: {...context.state.enbInfo, enodebs}};
    }
    return {enbInfo: enodebs};
  }
}
