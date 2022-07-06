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
import FiringAlerts from '../FiringAlerts';
import {MockApiUtil, alarmTestUtil} from '../../../test/testHelpers';
import {act, fireEvent, render} from '@testing-library/react';

import {PromFiringAlert} from '../../../../../../generated-ts';
import type {AlarmsWrapperProps} from '../../../test/testHelpers';

describe('FiringAlerts', () => {
  let AlarmsWrapper: React.ComponentType<Partial<AlarmsWrapperProps>>;
  let apiUtil: MockApiUtil;

  beforeEach(() => {
    ({apiUtil, AlarmsWrapper} = alarmTestUtil());
  });

  it('renders with default props', () => {
    const {getByText} = render(
      <AlarmsWrapper>
        <FiringAlerts />
      </AlarmsWrapper>,
    );
    expect(getByText(/Start creating alert rules/i)).toBeInTheDocument();
    expect(getByText(/Add Alert Rule/i)).toBeInTheDocument();
  });

  it('renders firing alerts', () => {
    const firingAlarms = [
      ({
        labels: {alertname: '<<testalert>>', severity: 'INFO'},
      } as unknown) as PromFiringAlert,
    ];
    jest
      .spyOn(apiUtil, 'viewFiringAlerts')
      .mockReturnValue({data: firingAlarms});
    const {getByText} = render(
      <AlarmsWrapper>
        <FiringAlerts />
      </AlarmsWrapper>,
    );
    expect(getByText('<<testalert>>')).toBeInTheDocument();
    expect(getByText(/info/i)).toBeInTheDocument();
  });

  it('clicking view alert shows alert details pane', () => {
    const firingAlarms = [
      ({
        labels: {alertname: '<<testalert>>', severity: 'INFO'},
      } as unknown) as PromFiringAlert,
    ];
    jest
      .spyOn(apiUtil, 'viewFiringAlerts')
      .mockReturnValue({data: firingAlarms});
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
});
