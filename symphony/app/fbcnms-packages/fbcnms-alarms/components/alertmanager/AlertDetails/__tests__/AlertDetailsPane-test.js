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
import AlertDetailsPane from '../AlertDetailsPane';
import {act, cleanup, fireEvent, render} from '@testing-library/react';
import {alarmTestUtil} from '../../../../test/testHelpers';
import {mockAlert, mockRuleInterface} from '../../../../test/testData';
import type {AlertViewerProps} from '../../../rules/RuleInterface';

afterEach(() => {
  cleanup();
});

const {AlarmsWrapper} = alarmTestUtil();
const commonProps = {
  alert: mockAlert({labels: {alertname: '<<test alert>>'}}),
  onClose: jest.fn(),
};

describe('Basics', () => {
  test('renders with default props', () => {
    const {getByText, getByTestId} = render(
      <AlarmsWrapper>
        <AlertDetailsPane {...commonProps} />
      </AlarmsWrapper>,
    );
    expect(getByTestId('alert-details-pane')).toBeInTheDocument();
    expect(getByTestId('metric-alert-viewer')).toBeInTheDocument();
    expect(getByText('<<test alert>>')).toBeInTheDocument();
    expect(getByText(/labels/i)).toBeInTheDocument();
  });

  test('clicking the close button invokes onclose callback', () => {
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

  test('shows extra labels', () => {
    const alert = mockAlert({labels: {testLabel: 'testValue'}});
    const {getByText} = render(
      <AlarmsWrapper>
        <AlertDetailsPane {...commonProps} alert={alert} />
      </AlarmsWrapper>,
    );
    expect(getByText(/testLabel/i)).toBeInTheDocument();
    expect(getByText(/testValue/i)).toBeInTheDocument();
  });

  test('shows extra annotations', () => {
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

describe('Alert type selection', () => {
  test('by default, use the MetricAlertViewer', () => {
    const {getByTestId} = render(
      <AlarmsWrapper>
        <AlertDetailsPane {...commonProps} />
      </AlarmsWrapper>,
    );
    expect(getByTestId('metric-alert-viewer')).toBeInTheDocument();
  });
  test(
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
  test(
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
  test(
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
            [mockAlertType]: mockRuleInterface({AlertViewer: MockAlertViewer}),
          }}>
          <AlertDetailsPane {...commonProps} alert={alert} />
        </AlarmsWrapper>,
      );
      expect(getAlertTypeMock).toHaveBeenCalledWith(alert);
      expect(getByTestId('mock-alert-viewer')).toBeInTheDocument();
    },
  );
});
