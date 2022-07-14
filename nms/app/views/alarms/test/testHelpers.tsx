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
import AlarmContext from '../components/AlarmContext';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import defaultTheme from '../../../theme/default';
import getPrometheusRuleInterface from '../components/rules/PrometheusEditor/getRuleInterface';
import {MemoryRouter} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {RenderResult, act, render} from '@testing-library/react';
import {SnackbarProvider} from 'notistack';
import type {AlarmContext as AlarmContextType} from '../components/AlarmContext';
import type {ApiUtil} from '../components/AlarmsApi';
import type {RuleInterfaceMap} from '../components/rules/RuleInterface';

type MakeSync<U extends object> = {
  [k in keyof U]: U[k] extends (...args: [...(infer ARGS)]) => Promise<infer V>
    ? (args: ARGS) => V
    : U[k];
};

export type MockApiUtil = MakeSync<ApiUtil>;

/**
 * Make sure when adding new functions to ApiUtil to add their mocks here
 */
export function mockApiUtil(merge?: Partial<ApiUtil>): MockApiUtil {
  /**
   * I don't understand how to properly type these mocks so using any for now.
   * The consuming code is all strongly typed, this shouldn't be much of an issue.
   */
  const useAlarmsApi = jest.fn<any, any>(
    <TParams, TResponse>(
      func: (params: TParams) => Promise<{data: TResponse}>,
      params: TParams,
    ) => ({
      isLoading: false,
      // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment,@typescript-eslint/no-unsafe-member-access
      response: (func(params) as any)?.data,
      error: null,
    }),
  );

  return Object.assign(
    {
      useAlarmsApi,
      viewFiringAlerts: jest.fn(),
      viewMatchingAlerts: jest.fn(),
      getTroubleshootingLink: jest.fn(),
      createAlertRule: jest.fn(),
      editAlertRule: jest.fn(),
      getAlertRules: jest.fn(),
      deleteAlertRule: jest.fn(),
      createReceiver: jest.fn(),
      editReceiver: jest.fn(),
      getReceivers: jest.fn(),
      deleteReceiver: jest.fn(),
      getRouteTree: jest.fn(),
      editRouteTree: jest.fn(),
      getSuppressions: jest.fn(),
      getMetricNames: jest.fn(),
      getMetricSeries: jest.fn(),
      getGlobalConfig: jest.fn(),
      editGlobalConfig: jest.fn(),
      getTenants: jest.fn(),
      getAlertmanagerTenancy: jest.fn(),
      getPrometheusTenancy: jest.fn(),
    },
    merge || {},
  ) as MockApiUtil;
}

export async function renderAsync(
  renderElement: JSX.Element,
): Promise<RenderResult> {
  let result!: RenderResult;
  await act(async () => {
    result = await Promise.resolve().then(() => render(renderElement));
  });
  return result;
}
export type AlarmsWrapperProps = AlarmContextType & {
  children?: React.ReactNode;
};
export function AlarmsWrapper({children, ...contextProps}: AlarmsWrapperProps) {
  return (
    <AlarmsTestWrapper>
      <AlarmContext.Provider value={contextProps}>
        {children}
      </AlarmContext.Provider>
    </AlarmsTestWrapper>
  );
}

export function AlarmsTestWrapper({
  route,
  children,
}: {
  route?: string;
  children: React.ReactNode;
}) {
  return (
    <MemoryRouter initialEntries={[route || '/']} initialIndex={0}>
      <MuiThemeProvider theme={defaultTheme}>
        <MuiStylesThemeProvider theme={defaultTheme}>
          <SnackbarProvider>{children}</SnackbarProvider>
        </MuiStylesThemeProvider>
      </MuiThemeProvider>
    </MemoryRouter>
  );
}

export type AlarmTestUtil = {
  AlarmsWrapper: React.ComponentType<Partial<AlarmsWrapperProps>>;
  apiUtil: MockApiUtil;
  ruleMap: RuleInterfaceMap<any>;
};

/**
 * All in one function to setup alarm tests.
 * * Constructs a mock apiUtil, mock rule map, and creates an AlarmsWrapper
 * function with both of these mocks passed in as props.
 *
 * Example:
 *
 * const {apiUtil, AlarmsWrapper} = alarmTestUtil()
 * test('my component', () => {
 *   render(
 *    <AlarmsWrapper>
 *      <MyComponent/>
 *    </AlarmsWrapper>
 *   )
 *   expect(apiUtil.someFunction).toHaveBeenCalled();
 * })
 */
export function alarmTestUtil(
  overrides?: Partial<AlarmTestUtil>,
): AlarmTestUtil {
  const apiUtil = overrides?.apiUtil ?? mockApiUtil();
  const ruleMap =
    overrides?.ruleMap ??
    getPrometheusRuleInterface({
      apiUtil: (apiUtil as unknown) as ApiUtil,
    });
  return {
    apiUtil,
    ruleMap,
    AlarmsWrapper: (props: Partial<AlarmsWrapperProps>) => (
      <AlarmsWrapper
        apiUtil={(apiUtil as unknown) as ApiUtil}
        ruleMap={ruleMap}
        {...props}
      />
    ),
  };
}
