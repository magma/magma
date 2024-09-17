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
import axios from 'axios';
import {useAlarmContext} from './AlarmContext';
import {useEnqueueSnackbar} from '../../../hooks/useSnackbar';
import {useParams} from 'react-router-dom';
import type {AlertRoutingTree} from '../../../../generated';
import type {ApiUtil} from './AlarmsApi';
import type {GenericRule, RuleInterfaceMap} from './rules/RuleInterface';

/**
 * Loads alert rules for each rule type. Rules are loaded in parallel and if one
 * rule loader fails or exceeds the load timeout, it will be cancelled and a
 * snackbar will be enqueued.
 */
export function useLoadRules<TRuleUnion>({
  ruleMap,
  lastRefreshTime,
}: {
  ruleMap: RuleInterfaceMap<TRuleUnion>;
  lastRefreshTime: string;
}): {rules: Array<GenericRule<TRuleUnion>>; isLoading: boolean} {
  const networkId = useNetworkId();
  const enqueueSnackbar = useEnqueueSnackbar();
  const [isLoading, setIsLoading] = React.useState(true);
  const [rules, setRules] = React.useState<Array<GenericRule<TRuleUnion>>>([]);

  React.useEffect(() => {
    const promises = Object.keys(ruleMap || {}).map((ruleType: string) => {
      const cancelSource = axios.CancelToken.source();
      const request = {
        networkId,
        cancelToken: cancelSource.token,
      };

      const ruleInterface = ruleMap[ruleType];

      return new Promise<Array<GenericRule<TRuleUnion>>>(resolve => {
        ruleInterface
          .getRules(request)
          .then(response => resolve(response))
          .catch(error => {
            console.error(error);
            enqueueSnackbar(
              `An error occurred while loading ${ruleInterface.friendlyName} rules`,
              {
                variant: 'error',
              },
            );
            resolve([]);
          });
      });
    });

    void Promise.all(promises).then(results => {
      const allResults = results.flat();
      setRules(allResults);
      setIsLoading(false);
    });
  }, [enqueueSnackbar, lastRefreshTime, networkId, ruleMap]);

  return {
    rules,
    isLoading,
  };
}

/**
 * An alert rule can have a 1:N mapping from rule name to receivers. This
 * mapping is configured via the route tree. This hook loads the route tree
 * and returns all routes which contain a matcher for exactly this alertname.
 */
export function useAlertRuleReceiver({
  ruleName,
  apiUtil,
}: {
  ruleName: string;
  apiUtil: ApiUtil;
}) {
  const networkId = useNetworkId();
  const {response} = apiUtil.useAlarmsApi(apiUtil.getRouteTree, {
    networkId,
  });

  // find all the routes which contain an alertname matcher for this alert
  const routesForAlertRule = React.useMemo<Array<AlertRoutingTree>>(() => {
    if (!(response && response.routes)) {
      return [];
    }
    // only go one level deep for now
    const matchingRoutes = response.routes.filter(route => {
      const match = (route.match as unknown) as Record<string, string>;
      return match && match['alertname'] && match['alertname'] === ruleName;
    });
    return matchingRoutes;
  }, [response, ruleName]);

  const [initialReceiver, setInitialReceiver] = React.useState<string | null>();
  const [receiver, setReceiver] = React.useState<string | null>();

  /**
   * once the routes are loaded, set the initial receiver so we can determine
   * if we need to add/remove/update routes
   */
  React.useEffect(() => {
    if (routesForAlertRule && routesForAlertRule.length > 0) {
      const _initialReceiver = routesForAlertRule[0].receiver;
      setInitialReceiver(_initialReceiver);
      setReceiver(_initialReceiver);
    }
  }, [routesForAlertRule]);

  const saveReceiver = React.useCallback(async () => {
    let updatedRoutes: AlertRoutingTree = response || {
      routes: [],
      receiver: `${networkId || 'tg'}_tenant_base_route`,
    };
    if (
      (receiver == null || receiver.trim() === '') &&
      initialReceiver != null
    ) {
      // remove the route
      updatedRoutes = {
        ...updatedRoutes,
        routes: (response?.routes || []).filter(
          route => route.receiver !== initialReceiver,
        ),
      };
    } else if (receiver != null && initialReceiver == null) {
      // creating a new route
      const newRoute: AlertRoutingTree = {
        receiver: receiver,
        match: {
          alertname: ruleName,
        },
      };
      updatedRoutes = {
        ...updatedRoutes,
        routes: [...(response?.routes || []), newRoute],
      };
    } else {
      // update existing route
      updatedRoutes = {
        ...updatedRoutes,
        routes: (response?.routes || []).map(route => {
          if (route.receiver !== initialReceiver) {
            return route;
          }
          return {
            ...route,
            receiver: receiver || '',
          };
        }),
      };
    }
    await apiUtil.editRouteTree({
      networkId: networkId,
      route: updatedRoutes,
    });
  }, [receiver, initialReceiver, apiUtil, networkId, response, ruleName]);

  return {receiver, setReceiver, saveReceiver};
}

export function useNetworkId(): string {
  const params = useParams<{networkId: string}>();
  const {getNetworkId} = useAlarmContext();
  if (typeof getNetworkId === 'function') {
    return getNetworkId();
  }
  return params.networkId!;
}
