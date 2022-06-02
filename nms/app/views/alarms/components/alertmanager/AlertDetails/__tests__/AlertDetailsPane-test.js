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
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import AlertDetailsPane from '../AlertDetailsPane';
import {act, fireEvent, render} from '@testing-library/react';
import {alarmTestUtil} from '../../../../test/testHelpers';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {mockAlert, mockRuleInterface} from '../../../../test/testData';
import type {AlarmsWrapperProps} from '../../../../test/testHelpers';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {AlertViewerProps} from '../../../rules/RuleInterface';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {ApiUtil} from '../../../AlarmsApi';

describe('AlertDetailsPane', () => {
  let AlarmsWrapper: React.ComponentType<$Shape<AlarmsWrapperProps>>;
  let apiUtil: ApiUtil;

  const commonProps = {
    alert: mockAlert({labels: {alertname: '<<test alert>>'}}),
    onClose: jest.fn(),
  };

  beforeEach(() => {
    ({apiUtil, AlarmsWrapper} = alarmTestUtil());
  });

  describe('Basics', () => {
    it('renders with default props', () => {
      const {getByText, getByTestId} = render(
        <AlarmsWrapper>
          <AlertDetailsPane {...commonProps} />
        </AlarmsWrapper>,
      );
      expect(getByTestId('alert-details-pane')).toBeInTheDocument();
      expect(getByTestId('metric-alert-viewer')).toBeInTheDocument();
      expect(getByText('<<test alert>>')).toBeInTheDocument();
    });

    it('clicking the close button invokes onclose callback', () => {
      const {getByTestId} = render(
        <AlarmsWrapper>
          <AlertDetailsPane {...commonProps} />
        </AlarmsWrapper>,
      );

      const closeButton = getByTestId('alert-details-close');
      expect(closeButton).toBeInTheDocument();
      act(() => {
        fireEvent.click(closeButton);
      });
      expect(commonProps.onClose).toHaveBeenCalled();
    });

    it('shows extra labels', () => {
      const alert = mockAlert({labels: {testLabel: 'testValue'}});
      const {getByText} = render(
        <AlarmsWrapper>
          <AlertDetailsPane {...commonProps} alert={alert} />
        </AlarmsWrapper>,
      );
      expect(getByText(/testLabel/i)).toBeInTheDocument();
      expect(getByText(/testValue/i)).toBeInTheDocument();
    });

    it('shows extra annotations', () => {
      const alert = mockAlert({annotations: {testAnnotation: 'testValue'}});
      const {getByText} = render(
        <AlarmsWrapper>
          <AlertDetailsPane {...commonProps} alert={alert} />
        </AlarmsWrapper>,
      );
      expect(getByText(/testAnnotation/i)).toBeInTheDocument();
      expect(getByText(/testValue/i)).toBeInTheDocument();
    });
  });

  it('shows troubleshooting link', () => {
    const alert = mockAlert({labels: {testLabel: 'testValue'}});
    jest.spyOn(apiUtil, 'getTroubleshootingLink').mockReturnValue({
      link: 'www.example.com',
      title: 'View troubleshooting documentation',
    });

    const {getByText} = render(
      <AlarmsWrapper>
        <AlertDetailsPane {...commonProps} alert={alert} />
      </AlarmsWrapper>,
    );
    expect(
      getByText(/View troubleshooting documentation/i),
    ).toBeInTheDocument();
  });

  describe('Alert type selection', () => {
    let AlarmsWrapper: React.ComponentType<$Shape<AlarmsWrapperProps>>;

    beforeEach(() => {
      ({AlarmsWrapper} = alarmTestUtil());
    });

    it('by default, use the MetricAlertViewer', () => {
      const {getByTestId} = render(
        <AlarmsWrapper>
          <AlertDetailsPane {...commonProps} />
        </AlarmsWrapper>,
      );
      expect(getByTestId('metric-alert-viewer')).toBeInTheDocument();
    });
    it(
      'if getAlertType returns an unconfigured alert source, ' +
        'fallback to the default',
      () => {
        const getAlertTypeMock = jest.fn(() => 'unconfigured-source');
        const alert = mockAlert();
        const {getByTestId} = render(
          <AlarmsWrapper getAlertType={getAlertTypeMock}>
            <AlertDetailsPane {...commonProps} alert={alert} />
          </AlarmsWrapper>,
        );
        expect(getAlertTypeMock).toHaveBeenCalledWith(alert);
        expect(getByTestId('metric-alert-viewer')).toBeInTheDocument();
      },
    );
    it(
      'if getAlertType returns a alert source without an AlertViewer, ' +
        'fallback to default',
      () => {
        const getAlertTypeMock = jest.fn(() => 'prometheus');
        const alert = mockAlert();
        const {getByTestId} = render(
          <AlarmsWrapper getAlertType={getAlertTypeMock}>
            <AlertDetailsPane {...commonProps} alert={alert} />
          </AlarmsWrapper>,
        );
        expect(getAlertTypeMock).toHaveBeenCalledWith(alert);
        expect(getByTestId('metric-alert-viewer')).toBeInTheDocument();
      },
    );
    it(
      'if getAlertType returns an alert source with an AlertViewer, ' +
        'renders the AlertViewer',
      () => {
        const mockAlertType = 'test';
        const getAlertTypeMock = jest.fn(() => mockAlertType);
        function MockAlertViewer(_props: AlertViewerProps) {
          return <div data-testid="mock-alert-viewer" />;
        }
        const alert = mockAlert();
        const {getByTestId} = render(
          <AlarmsWrapper
            getAlertType={getAlertTypeMock}
            ruleMap={{
              [mockAlertType]: mockRuleInterface({
                AlertViewer: MockAlertViewer,
              }),
            }}>
            <AlertDetailsPane {...commonProps} alert={alert} />
          </AlarmsWrapper>,
        );
        expect(getAlertTypeMock).toHaveBeenCalledWith(alert);
        expect(getByTestId('mock-alert-viewer')).toBeInTheDocument();
      },
    );
  });
});
