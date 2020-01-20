/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

jest.mock('@fbcnms/sequelize-models');

const {
  appMiddleware,
  organizationMiddleware,
} = require('@fbcnms/express-middleware');
import express from 'express';
import request from 'supertest';

import {Organization} from '@fbcnms/sequelize-models';

const TEST_PAYLOAD = {
  receiver: 'webhook',
  status: 'firing',
  alerts: [
    {
      status: 'firing',
      labels: {
        alertname: 'My Alert',
        networkID: 'network1',
        gatewayID: 'gateway1',
      },
      annotations: {},
      startsAt: '2018-08-03T09:52:26.739266876+02:00',
      endsAt: '0001-01-01T00:00:00Z',
      generatorURL:
        'http://simon-laptop:9090/graph?g0.expr=go_memstats_alloc_bytes+%3E+0\u0026g0.tab=1',
    },
  ],
  groupLabels: {},
  commonLabels: {},
  commonAnnotations: {},
  externalURL: 'http://simon-laptop:9093',
  version: '4',
  groupKey: '{}:{networkID="network1", job="prometheus24"}',
};

beforeAll(async () => {
  await Organization.create({
    id: '1',
    name: 'master',
    customDomains: [],
    networkIDs: [],
    csvCharset: '',
    ssoCert: '',
    ssoEntrypoint: '',
    ssoIssuer: '',
  });
  await Organization.create({
    id: '2',
    name: 'notmaster',
    customDomains: [],
    networkIDs: ['network1'],
    csvCharset: '',
    ssoCert: '',
    ssoEntrypoint: '',
    ssoIssuer: '',
  });
});

describe('magma webhook endpoint', () => {
  const app = setupApp();

  test('basic alert works', async () => {
    const triggerMagmaAlert = jest
      .spyOn(require('../../graphgrpc/magmaalert'), 'triggerActionsAlert')
      .mockReturnValue(true);
    const res = await request(app)
      .post('/webhooks/magma')
      .send(TEST_PAYLOAD)
      .set('Host', 'master.myhost.com')
      .set('Content-Type', 'application/json')
      .set('Accept', 'application/json');
    expect(res.statusCode).toEqual(200);
    expect(res.body).toEqual({success: true, message: 'ok'});

    expect(triggerMagmaAlert).toHaveBeenCalledWith({
      tenantID: 'notmaster',
      alertname: 'My Alert',
      networkID: 'network1',
      labels: {
        alertname: 'My Alert',
        gatewayID: 'gateway1',
        networkID: 'network1',
      },
    });
  });

  test('only master allowed', async () => {
    const app = setupApp();
    const res = await request(app)
      .post('/webhooks/magma')
      .send(TEST_PAYLOAD)
      .set('Host', 'notmaster.myhost.com')
      .set('Content-Type', 'application/json')
      .set('Accept', 'application/json');

    expect(res.statusCode).toEqual(302);
    expect(res.header['location']).toEqual('/user/login');
  });
});

function setupApp() {
  const app = express();
  app.use(appMiddleware());
  app.use(organizationMiddleware());
  app.use('/webhooks', require('../routes').default);
  return app;
}
