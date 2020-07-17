/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import {Strategy} from 'passport-strategy';

import DynamicStrategy from '../DynamicStrategy';

import passport from 'passport';

class StubStrategy extends Strategy {
  constructor(success: boolean) {
    super();
    this.name = 'stub';
    this._success = success;
  }
  authenticate() {
    if (this._success) {
      return this.success({message: 'This is a success'});
    } else {
      return this.fail({message: 'This is a failure'});
    }
  }
}

function flushPromises() {
  return new Promise(resolve => setImmediate(resolve));
}

test('authenticate failure', async () => {
  const req = {
    body: {},
    logIn: jest.fn((user, options, next) => {
      next(null);
      return true;
    }),
  };

  const res = {end: jest.fn(), statusCode: null};
  const next = jest.fn();

  const p = passport.use(
    'dynamic',
    new DynamicStrategy(
      async _req => 'strategyID',
      async _req => new StubStrategy(false),
    ),
  );

  p.authenticate('dynamic')(req, res, next);

  await flushPromises();

  expect(req.logIn).not.toHaveBeenCalled();
  expect(res.statusCode).toBe(401);
  expect(res.end.mock.calls[0][0]).toBe('Unauthorized');
  expect(next.mock.calls.length).toBe(0);
});

test('authenticate success', async () => {
  const req = {
    body: {},
    logIn: jest.fn((user, options, next) => {
      next(null);
      return true;
    }),
  };

  const res = {end: jest.fn(), statusCode: null};
  const next = jest.fn();

  const p = passport.use(
    'dynamic',
    new DynamicStrategy(
      async _req => 'strategyID',
      async _req => new StubStrategy(true),
    ),
  );

  p.authenticate('dynamic', {})(req, res, next);

  await flushPromises();

  expect(req.logIn).toHaveBeenCalled();
  expect(res.statusCode).toBe(null);
  expect(res.end).not.toHaveBeenCalled();
  expect(next.mock.calls.length).toBe(1);
});
