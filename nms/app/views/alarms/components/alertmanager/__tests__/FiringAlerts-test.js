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

import * as React from 'react';
import FiringAlerts from '../FiringAlerts';
import {act, fireEvent, render} from '@testing-library/react';
import {alarmTestUtil} from '../../../test/testHelpers';

import type {AlarmsWrapperProps} from '../../../test/testHelpers';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {ApiUtil} from '../../AlarmsApi';
// $FlowFixMe migrated to typescript
import type {FiringAlarm} from '../../AlarmAPIType';

describe('FiringAlerts', () => {
  let AlarmsWrapper: React.ComponentType<$Shape<AlarmsWrapperProps>>;
  let apiUtil: ApiUtil;

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
    const firingAlarms: Array<$Shape<FiringAlarm>> = [
      {
        labels: {alertname: '<<testalert>>', severity: 'INFO'},
      },
    ];
    jest.spyOn(apiUtil, 'viewFiringAlerts').mockReturnValue(firingAlarms);
    const {getByText} = render(
      <AlarmsWrapper>
        <FiringAlerts />
      </AlarmsWrapper>,
    );
    expect(getByText('<<testalert>>')).toBeInTheDocument();
    expect(getByText(/info/i)).toBeInTheDocument();
  });

  it('clicking view alert shows alert details pane', async () => {
    const firingAlarms: Array<$Shape<FiringAlarm>> = [
      {
        labels: {alertname: '<<testalert>>', severity: 'INFO'},
      },
    ];
    jest.spyOn(apiUtil, 'viewFiringAlerts').mockReturnValue(firingAlarms);
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
