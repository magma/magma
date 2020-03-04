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
import GlobalConfig from '../GlobalConfig';
import React from 'react';
import {act, fireEvent, render} from '@testing-library/react';
import {alarmTestUtil} from '../../../../test/testHelpers';

import type {AlertManagerGlobalConfig} from '../../../AlarmAPIType';

const {AlarmsWrapper, apiUtil} = alarmTestUtil();

const commonProps = {
  onExit: jest.fn(),
};

const defaultResponse: AlertManagerGlobalConfig = {
  resolve_timeout: '10m',
  slack_api_url: 'slack.com',
  pagerduty_url: 'pagerduty.com',
  smtp_from: 'hello@terragraph.link',
  smtp_hello: 'terragraph.link',
  smtp_smarthost: 'smtp.terragraph.link:25',
  smtp_auth_username: 'tg',
  smtp_auth_password: 'password',
  smtp_auth_secret: 'smtp_auth_secret',
  smtp_auth_identity: 'smtp_auth_identity',
  smtp_require_tls: true,
  opsgenie_api_url: 'opsgenie.com',
  opsgenie_api_key: 'opsgenie_api_key',
  hipchat_api_url: 'hipchat.com',
  hipchat_auth_token: 'hipchat_auth_token',
  wechat_api_url: 'wechat.com',
  wechat_api_secret: 'abc123',
  wechat_api_corp_id: '12345',
  victorops_api_url: 'victorops.com',
  victorops_api_key: '',
  http_config: {
    bearer_token: 'bearer_token',
    proxy_url: 'http://proxy_url.terragraph.link',
    basic_auth: {
      username: 'basic_auth_username',
      password: 'basic_auth_password',
    },
  },
};

test('renders', () => {
  const {getByText} = render(
    <AlarmsWrapper>
      <GlobalConfig {...commonProps} />
    </AlarmsWrapper>,
  );
  expect(getByText(/global receiver settings/i)).toBeInTheDocument();
});

test('fills form inputs with values from backend', () => {
  jest.spyOn(apiUtil, 'getGlobalConfig').mockReturnValue(defaultResponse);
  const {getByTestId} = render(
    <AlarmsWrapper>
      <GlobalConfig {...commonProps} />
    </AlarmsWrapper>,
  );

  expect(getByTestId('resolve_timeout')).toHaveValue(
    defaultResponse.resolve_timeout,
  );
  expect(getByTestId('slack_api_url')).toHaveValue(
    defaultResponse.slack_api_url,
  );
  expect(getByTestId('pagerduty_url')).toHaveValue(
    defaultResponse.pagerduty_url,
  );
  expect(getByTestId('smtp_from')).toHaveValue(defaultResponse.smtp_from);
  expect(getByTestId('smtp_hello')).toHaveValue(defaultResponse.smtp_hello);
  expect(getByTestId('smtp_smarthost')).toHaveValue(
    defaultResponse.smtp_smarthost,
  );
  expect(getByTestId('smtp_auth_username')).toHaveValue(
    defaultResponse.smtp_auth_username,
  );
  expect(getByTestId('smtp_auth_password')).toHaveValue(
    defaultResponse.smtp_auth_password,
  );
  expect(getByTestId('smtp_auth_secret')).toHaveValue(
    defaultResponse.smtp_auth_secret,
  );
  expect(getByTestId('smtp_auth_identity')).toHaveValue(
    defaultResponse.smtp_auth_identity,
  );
  expect(getByTestId('opsgenie_api_url')).toHaveValue(
    defaultResponse.opsgenie_api_url,
  );
  expect(getByTestId('opsgenie_api_key')).toHaveValue(
    defaultResponse.opsgenie_api_key,
  );
  expect(getByTestId('hipchat_api_url')).toHaveValue(
    defaultResponse.hipchat_api_url,
  );
  expect(getByTestId('hipchat_auth_token')).toHaveValue(
    defaultResponse.hipchat_auth_token,
  );
  expect(getByTestId('wechat_api_url')).toHaveValue(
    defaultResponse.wechat_api_url,
  );
  expect(getByTestId('wechat_api_secret')).toHaveValue(
    defaultResponse.wechat_api_secret,
  );
  expect(getByTestId('wechat_api_corp_id')).toHaveValue(
    defaultResponse.wechat_api_corp_id,
  );
  expect(getByTestId('victorops_api_url')).toHaveValue(
    defaultResponse.victorops_api_url,
  );
  expect(getByTestId('victorops_api_key')).toHaveValue(
    defaultResponse.victorops_api_key,
  );
  expect(getByTestId('http_config_bearer_token')).toHaveValue(
    defaultResponse.http_config.bearer_token,
  );
  expect(getByTestId('http_config_proxy_url')).toHaveValue(
    defaultResponse.http_config?.proxy_url,
  );
  expect(getByTestId('basic_auth_username')).toHaveValue(
    defaultResponse.http_config?.basic_auth?.username,
  );
  expect(getByTestId('basic_auth_password')).toHaveValue(
    defaultResponse.http_config?.basic_auth?.password,
  );
});

test('submitting form submits updated values to backend', async () => {
  jest.spyOn(apiUtil, 'getGlobalConfig').mockReturnValue({});
  const editConfigMock = jest
    .spyOn(apiUtil, 'editGlobalConfig')
    .mockImplementationOnce(() => Promise.resolve());
  const {getByTestId} = render(
    <AlarmsWrapper>
      <GlobalConfig {...commonProps} />
    </AlarmsWrapper>,
  );
  act(() => {
    fireEvent.change(getByTestId('resolve_timeout'), {target: {value: '5m'}});
  });
  act(() => {
    fireEvent.change(getByTestId('slack_api_url'), {
      target: {value: 'https://hooks.slack.com'},
    });
  });
  expect(editConfigMock).not.toHaveBeenCalled();
  await act(async () => {
    fireEvent.click(getByTestId('editor-submit-button'));
  });

  expect(editConfigMock).toHaveBeenCalledWith({
    config: {resolve_timeout: '5m', slack_api_url: 'https://hooks.slack.com'},
  });
});

test('erasing values from form removes keys from request', async () => {
  jest.spyOn(apiUtil, 'getGlobalConfig').mockReturnValue({
    resolve_timeout: '5m',
    slack_api_url: 'https://hooks.slack.com',
  });
  const editConfigMock = jest
    .spyOn(apiUtil, 'editGlobalConfig')
    .mockImplementationOnce(() => Promise.resolve());
  const {getByTestId} = render(
    <AlarmsWrapper>
      <GlobalConfig {...commonProps} />
    </AlarmsWrapper>,
  );
  act(() => {
    fireEvent.change(getByTestId('resolve_timeout'), {target: {value: ''}});
  });
  await act(async () => {
    fireEvent.click(getByTestId('editor-submit-button'));
  });
  expect(editConfigMock).toHaveBeenCalledWith({
    config: {slack_api_url: 'https://hooks.slack.com'},
  });
});
