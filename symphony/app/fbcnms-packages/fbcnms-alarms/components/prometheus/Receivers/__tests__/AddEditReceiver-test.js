/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import 'jest-dom/extend-expect';
import * as React from 'react';
import AddEditReceiver from '../AddEditReceiver';
import {act, cleanup, fireEvent, render} from '@testing-library/react';
import {alarmTestUtil} from '@fbcnms/alarms/test/testHelpers';

const {apiUtil, AlarmsWrapper} = alarmTestUtil();

const commonProps = {
  isNew: true,
  onExit: jest.fn(),
};

afterEach(() => {
  cleanup();
});

test('renders', () => {
  const {getByLabelText} = render(
    <AlarmsWrapper>
      <AddEditReceiver {...commonProps} receiver={{name: ''}} />
    </AlarmsWrapper>,
  );
  expect(getByLabelText(/receiver name/i)).toBeInTheDocument();
});

test('clicking the add button adds a new config entry', async () => {
  const {getByLabelText, getByTestId, queryByTestId} = render(
    <AlarmsWrapper>
      <AddEditReceiver {...commonProps} receiver={{name: ''}} />
    </AlarmsWrapper>,
  );
  expect(queryByTestId('slack-config-editor')).not.toBeInTheDocument();
  act(() => {
    fireEvent.click(getByLabelText('add new receiver configuration'));
  });
  expect(getByTestId('slack-config-editor')).toBeInTheDocument();
});

test('editing a config entry then submitting submits the form state', () => {
  const {getByLabelText, getByTestId} = render(
    <AlarmsWrapper>
      <AddEditReceiver {...commonProps} receiver={{name: ''}} />
    </AlarmsWrapper>,
  );
  act(() => {
    fireEvent.click(getByLabelText('add new receiver configuration'));
  });
  act(() => {
    fireEvent.change(getByLabelText(/receiver name/i), {
      target: {value: 'test receiver'},
    });
  });
  act(() => {
    fireEvent.change(getByLabelText(/api url/i), {
      target: {value: 'https://slack.com/hook'},
    });
  });
  act(() => {
    fireEvent.change(getByLabelText(/channel/i), {
      target: {value: '#testchannel'},
    });
  });
  act(() => {
    fireEvent.click(getByTestId('editor-submit-button'));
  });
  expect(apiUtil.createReceiver).toHaveBeenCalledWith({
    receiver: {
      name: 'test receiver',
      slack_configs: [
        {api_url: 'https://slack.com/hook', channel: '#testchannel'},
      ],
    },
    networkId: undefined,
  });
});
