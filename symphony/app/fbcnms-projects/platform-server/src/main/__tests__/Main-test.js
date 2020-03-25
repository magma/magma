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

import app from '../../app';
import request from 'supertest';
import {Organization, User} from '@fbcnms/sequelize-models';

it('Returns a health check', async () => {
  await request(app)
    .get('/healthz')
    .expect(200)
    .expect('OK');
});

const ORGS = [
  {
    id: 1,
    name: 'myorg',
    customDomains: [],
    networkIDs: [],
    csvCharset: '',
    ssoCert: '',
    ssoEntrypoint: '',
    ssoIssuer: '',
  },
];

describe('login csrf token tests', () => {
  beforeAll(async () => {
    ORGS.forEach(async organization => await Organization.create(organization));
    await User.create({
      email: 'user@test.com',
      organization: 'myorg',
      password: 'password',
      role: 0,
      readOnly: false,
    });
  });

  const agent = request.agent(app).set('host', 'myorg.phb.io');
  const getCsrfToken = async (): Promise<string> => {
    const resp = await agent.get('/user/login').expect(200);
    const csrfTokenRegex = /"csrfToken":"([a-zA-Z\d_-]+)"/;
    const match = csrfTokenRegex.exec(resp.text);
    return match?.[1] || '';
  };

  it('returns csrftoken on login page', async () => {
    const csrfToken = await getCsrfToken();
    expect(csrfToken).not.toBe('');
  });

  it('is forbidden when no csrf token', async () => {
    await agent
      .post('/user/login')
      .send('username=test')
      .send('password=test')
      .expect(403)
      .expect(/invalid csrf token/);
  });

  it('allows submission with csrftoken', async () => {
    const csrfToken = await getCsrfToken();
    await agent
      .post('/user/login')
      .send('_csrf=' + csrfToken)
      .send('username=invaliduser')
      .send('password=invalidpassword')
      .send('to=/nms')
      .expect(302)
      .expect('Location', '/user/login?invalid=true&to=%2Fnms');
  });
});
