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
import EventAlertChart from '../EventAlertChart';
import MagmaAPIBindings from '@fbcnms/magma-api';
import React from 'react';
import axiosMock from 'axios';
import moment from 'moment';
import {MemoryRouter, Route} from 'react-router-dom';
import {cleanup, render, wait} from '@testing-library/react';
import type {promql_return_object} from '@fbcnms/magma-api';

afterEach(cleanup);

const mockMetricSt: promql_return_object = {
  status: 'success',
  data: {
    resultType: 'matrix',
    result: [
      {
        metric: {},
        values: [['1588898968.042', '6']],
      },
    ],
  },
};

jest.mock('axios');
jest.mock('@fbcnms/magma-api');
jest.mock('@fbcnms/ui/hooks/useSnackbar');

// chart component was failing here so mocking this out
// this shouldn't affect the prop verification part in the react
// chart component
window.HTMLCanvasElement.prototype.getContext = () => {};

describe('<EventAlertChart/>', () => {
  beforeEach(() => {
    // eslint-disable-next-line max-len
    MagmaAPIBindings.getNetworksByNetworkIdPrometheusQueryRange.mockResolvedValue(
      mockMetricSt,
    );
  });

  afterEach(() => {
    axiosMock.get.mockClear();
    MagmaAPIBindings.getNetworksByNetworkIdPrometheusQueryRange.mockClear();
  });

  const testCases = [
    {
      startDate: moment().subtract(2, 'hours'),
      endDate: moment(),
      step: '30s',
      valid: true,
    },
    {
      startDate: moment().subtract(10, 'day'),
      endDate: moment(),
      step: '4h',
      valid: true,
    },
    {
      startDate: moment(),
      endDate: moment().subtract(10, 'day'),
      step: '4h',
      valid: false,
    },
  ];

  testCases.forEach((tc, _) => {
    it('renders', async () => {
      // const endDate = moment();
      // const startDate = moment().subtract(3, 'hours');
      const Wrapper = () => (
        <MemoryRouter initialEntries={['/nms/mynetwork']} initialIndex={0}>
          <Route
            path="/nms/:networkId"
            render={props => (
              <EventAlertChart
                {...props}
                startEnd={[tc.startDate, tc.endDate]}
              />
            )}
          />
        </MemoryRouter>
      );
      render(<Wrapper />);
      await wait();
      if (tc.valid) {
        expect(
          MagmaAPIBindings.getNetworksByNetworkIdPrometheusQueryRange,
        ).toHaveBeenCalledTimes(1);
        expect(
          MagmaAPIBindings.getNetworksByNetworkIdPrometheusQueryRange.mock
            .calls[0][0].start,
        ).toEqual(tc.startDate.toISOString());
        expect(
          MagmaAPIBindings.getNetworksByNetworkIdPrometheusQueryRange.mock
            .calls[0][0].end,
        ).toEqual(tc.endDate.toISOString());
        expect(
          MagmaAPIBindings.getNetworksByNetworkIdPrometheusQueryRange.mock
            .calls[0][0].step,
        ).toEqual(tc.step);
      } else {
        // negative test for invalid start end use default timerange
        const defaultStep = '30s';
        expect(
          MagmaAPIBindings.getNetworksByNetworkIdPrometheusQueryRange.mock
            .calls[0][0].step,
        ).toEqual(defaultStep);
      }
    });
  });
});
