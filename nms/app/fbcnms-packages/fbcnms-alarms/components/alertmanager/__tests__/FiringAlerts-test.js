/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import 'jest-dom/extend-expect';
import * as React from 'react';
import FiringAlerts from '../FiringAlerts';
import {act, cleanup, fireEvent, render} from '@testing-library/react';
import {alarmTestUtil} from '../../../test/testHelpers';

import type {FiringAlarm} from '../../AlarmAPIType';

const {apiUtil, AlarmsWrapper} = alarmTestUtil();

afterEach(() => {
  cleanup();
});

test('renders with default props', () => {
  const {getByText} = render(
    <AlarmsWrapper>
      <FiringAlerts />
    </AlarmsWrapper>,
  );
  expect(getByText(/name/i)).toBeInTheDocument();
  expect(getByText(/severity/i)).toBeInTheDocument();
  expect(getByText(/labels/i)).toBeInTheDocument();
});

test('renders firing alerts', () => {
  const firingAlarms: Array<$Shape<FiringAlarm>> = [
    {
      labels: {alertname: '<<testalert>>', severity: 'INFO'},
    },
  ];
  jest.spyOn(apiUtil, 'viewFiringAlerts').mockReturnValueOnce(firingAlarms);
  const {getByText} = render(
    <AlarmsWrapper>
      <FiringAlerts />
    </AlarmsWrapper>,
  );
  expect(getByText('<<testalert>>')).toBeInTheDocument();
  expect(getByText(/info/i)).toBeInTheDocument();
});

test('clicking view alert shows alert details pane', async () => {
  const firingAlarms: Array<$Shape<FiringAlarm>> = [
    {
      labels: {alertname: '<<testalert>>', severity: 'INFO'},
    },
  ];
  jest.spyOn(apiUtil, 'viewFiringAlerts').mockReturnValueOnce(firingAlarms);
  const {getByText, getByTestId} = render(
    <AlarmsWrapper>
      <FiringAlerts />
    </AlarmsWrapper>,
  );
  act(() => {
    fireEvent.click(getByText('<<testalert>>'));
  });

  expect(getByTestId('alert-details-pane')).toBeInTheDocument();
});
