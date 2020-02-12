/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import MagmaAPIBindings from '@fbcnms/magma-api';
import axiosMock from 'axios';
import {parseFileAndSave} from '../ImportSubscribersDialog';

const match = {params: {}, isExact: true, path: '', url: ''};
jest.mock('axios');
jest.mock('@fbcnms/magma-api');

describe('ImportSubscribers parseFileAndSave', () => {
  const setErrorMsg = jest.fn();
  const props = {onSave: jest.fn(), onSaveError: jest.fn()};
  const id = 311;
  const CSV_HEADER = 'imsi,lte_state,lte_auth_key,lte_auth_opc,sub_profile';
  beforeEach(() => {
    MagmaAPIBindings.postLteByNetworkIdSubscribers.mockResolvedValueOnce(id);
  });

  afterEach(() => {
    axiosMock.get.mockClear();
    setErrorMsg.mockClear();
    props.onSave.mockClear();
    props.onSaveError.mockClear();
  });

  it('fails on binary inputs', async () => {
    const csvText = new ArrayBuffer(8);
    await parseFileAndSave(csvText, setErrorMsg, match, props);
    expect(setErrorMsg.mock.calls.length).toBe(1);
    expect(props.onSave.mock.calls.length).toBe(0);
    expect(props.onSaveError.mock.calls.length).toBe(0);
  });

  it('parses for mac files', async () => {
    const csvText = CSV_HEADER + '\n' + id + ',ACTIVE,1,1,default';
    await parseFileAndSave(csvText, setErrorMsg, match, props);
    expect(setErrorMsg.mock.calls.length).toBe(0);
    expect(props.onSave.mock.calls.length).toBe(1);
    // The first argument to the function was [id]
    const onSaveArg = props.onSave.mock.calls[0][0];
    expect(onSaveArg.length).toBe(1);
    expect(onSaveArg[0]).toBe(id);
    expect(props.onSaveError.mock.calls.length).toBe(0);
  });

  it('parses for windows files', async () => {
    const csvText = CSV_HEADER + '\r\n' + id + ',ACTIVE,1,1,default';
    await parseFileAndSave(csvText, setErrorMsg, match, props);
    expect(setErrorMsg.mock.calls.length).toBe(0);
    expect(props.onSave.mock.calls.length).toBe(1);
    // The first argument to the function was [id]
    const onSaveArg = props.onSave.mock.calls[0][0];
    expect(onSaveArg.length).toBe(1);
    expect(onSaveArg[0]).toBe(id);
    expect(props.onSaveError.mock.calls.length).toBe(0);
  });
});
