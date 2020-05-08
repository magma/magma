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

describe('<EventAlertChart />', () => {
  beforeEach(() => {
    // eslint-disable-next-line max-len
    MagmaAPIBindings.getNetworksByNetworkIdPrometheusQueryRange.mockResolvedValue(
      mockMetricSt,
    );
  });

  afterEach(() => {
    axiosMock.get.mockClear();
  });

  const Wrapper = () => (
    <MemoryRouter initialEntries={['/nms/mynetwork']} initialIndex={0}>
      <Route path="/nms/:networkId" component={EventAlertChart} />
    </MemoryRouter>
  );
  it('renders', async () => {
    render(<Wrapper />);
    await wait();
    expect(
      MagmaAPIBindings.getNetworksByNetworkIdPrometheusQueryRange,
    ).toHaveBeenCalledTimes(1);
  });
});
