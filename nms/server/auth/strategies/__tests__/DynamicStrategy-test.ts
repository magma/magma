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

import DynamicStrategy from '../DynamicStrategy';
import passport from 'passport';
import {NextFunction} from 'express';
import {Strategy} from 'passport-strategy';

class StubStrategy extends Strategy {
  name: string;
  _success: boolean;

  constructor(success: boolean) {
    super();
    this.name = 'stub';
    this._success = success;
  }

  authenticate() {
    if (this._success) {
      return this.success({message: 'This is a success'});
    } else {
      return this.fail({message: 'This is a failure'}, 401);
    }
  }
}

function flushPromises() {
  return new Promise(resolve => setImmediate(resolve));
}

test('authenticate failure', async () => {
  const req = {
    body: {},
    logIn: jest.fn((user, options, next: NextFunction) => {
      next(null);
      return true;
    }),
  };

  const res = {end: jest.fn<unknown, [string]>(), statusCode: null};
  const next = jest.fn();

  // eslint-disable-next-line @typescript-eslint/no-unsafe-call
  passport
    .use(
      'dynamic',
      new DynamicStrategy(
        () => Promise.resolve('strategyID'),
        () => Promise.resolve(new StubStrategy(false)),
      ),
    )
    .authenticate('dynamic')(req, res, next);

  await flushPromises();

  expect(req.logIn).not.toHaveBeenCalled();
  expect(res.statusCode).toBe(401);
  expect(res.end.mock.calls[0][0]).toBe('Unauthorized');
  expect(next.mock.calls.length).toBe(0);
});

test('authenticate success', async () => {
  const req = {
    body: {},
    logIn: jest.fn((user, options, next: NextFunction) => {
      next(null);
      return true;
    }),
  };

  const res = {end: jest.fn(), statusCode: null};
  const next = jest.fn();

  // eslint-disable-next-line @typescript-eslint/no-unsafe-call
  passport
    .use(
      'dynamic',
      new DynamicStrategy(
        () => Promise.resolve('strategyID'),
        () => Promise.resolve(new StubStrategy(true)),
      ),
    )
    .authenticate('dynamic', {})(req, res, next);

  await flushPromises();

  expect(req.logIn).toHaveBeenCalled();
  expect(res.statusCode).toBe(null);
  expect(res.end).not.toHaveBeenCalled();
  expect(next.mock.calls.length).toBe(1);
});
