/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import {act, render} from '@testing-library/react';
import type {ApiUtil} from '../components/AlarmsApi';

/**
 * I don't understand how to properly type these mocks so using any for now.
 * The consuming code is all strongly typed, this shouldn't be much of an issue.
 */
// eslint-disable-next-line flowtype/no-weak-types
export const useMagmaAPIMock = jest.fn<any, any>(() => ({
  isLoading: false,
  response: [],
  error: null,
}));

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
