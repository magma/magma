/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import * as React from 'react';
import AlarmContext from '../components/AlarmContext';
import getPrometheusRuleInterface from '../components/rules/PrometheusEditor/getRuleInterface';
import {SymphonyWrapper} from '@fbcnms/test/testHelpers';
import {act, render} from '@testing-library/react';
import type {AlarmContext as AlarmContextType} from '../components/AlarmContext';
import type {ApiUtil} from '../components/AlarmsApi';
import type {RuleInterfaceMap} from '../components/rules/RuleInterface';

/**
 * I don't understand how to properly type these mocks so using any for now.
 * The consuming code is all strongly typed, this shouldn't be much of an issue.
 */
// eslint-disable-next-line flowtype/no-weak-types
export const useMagmaAPIMock = jest.fn<any, any>(
  <TParams, TResponse>(
    func: TParams => Promise<TResponse>,
    params: TParams,
    _cacheCounter?: string | number,
  ) => ({
    isLoading: false,
    response: func(params),
    error: null,
  }),
);

/**
 * Make sure when adding new functions to ApiUtil to add their mocks here
 */
export function mockApiUtil(merge?: $Shape<ApiUtil>): ApiUtil {
  return Object.assign(
    {
      useAlarmsApi: useMagmaAPIMock,
      viewFiringAlerts: jest.fn(),
      viewMatchingAlerts: jest.fn(),
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
      getMetricSeries: jest.fn(),
    },
    merge || {},
  );
}

// eslint-disable-next-line flowtype/no-weak-types
export async function renderAsync(...renderArgs: Array<any>): Promise<any> {
  let result;
  await act(async () => {
    result = await render(...renderArgs);
  });
  return result;
}

type AlarmsWrapperProps = {|
  children: React.Node,
  ...$Shape<AlarmContextType>,
|};
export function AlarmsWrapper({children, ...contextProps}: AlarmsWrapperProps) {
  return (
    <SymphonyWrapper>
      <AlarmContext.Provider value={contextProps}>
        {children}
      </AlarmContext.Provider>
    </SymphonyWrapper>
  );
}

export type AlarmTestUtil = {|
  AlarmsWrapper: React.ComponentType<$Shape<AlarmsWrapperProps>>,
  apiUtil: ApiUtil,
  ruleMap: RuleInterfaceMap<*>,
|};

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
  overrides?: $Shape<AlarmTestUtil>,
): AlarmTestUtil {
  const apiUtil = overrides?.apiUtil ?? mockApiUtil();
  const ruleMap = overrides?.ruleMap ?? getPrometheusRuleInterface({apiUtil});
  return {
    apiUtil,
    ruleMap,
    AlarmsWrapper: (props: AlarmsWrapperProps) => (
      <AlarmsWrapper apiUtil={apiUtil} ruleMap={ruleMap} {...props} />
    ),
  };
}
