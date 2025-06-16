/*
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

import express, {Request, Response} from 'express';
import {agent as request} from 'supertest';
import {rateLimitMiddleware} from '..';

const config = {
  RATE_LIMIT_CONFIG: {
    windowMs: 15 * 1000, // 15 seconds
    limit: 2,
  },
};

jest.mock('../../../config/config', () => config);

describe('Rate limit test', () => {
  const client1 = '1.2.3.4';
  const client2 = '2.3.4.5';

  let app: express.Express;
  let now = Date.now();

  beforeEach(() => {
    jest.useFakeTimers({legacyFakeTimers: true});
    jest.spyOn(Date, 'now').mockImplementation(() => now);

    app = express();
    app.set('trust proxy', true); // Needed to respect X-Forwarded-For
    app.use(rateLimitMiddleware);
    app.get('/', (_req: Request, res: Response) => res.send('Hi there!'));
  });

  afterEach(() => {
    jest.clearAllTimers();
    jest.clearAllMocks();
  });

  it('disallows request after 2 requests within 15 seconds for a given IP', async () => {
    // Client 1
    await request(app).get('/').set('X-Forwarded-For', client1).expect(200);
    await request(app).get('/').set('X-Forwarded-For', client1).expect(200);
    await request(app).get('/').set('X-Forwarded-For', client1).expect(429);

    // Client 2
    await request(app).get('/').set('X-Forwarded-For', client2).expect(200);
    await request(app).get('/').set('X-Forwarded-For', client2).expect(200);
    await request(app).get('/').set('X-Forwarded-For', client2).expect(429);

    // Advance virtual time by 16 seconds
    now += 16000;
    jest.advanceTimersByTime(16000);

    // Should be allowed again after window resets
    await request(app).get('/').set('X-Forwarded-For', client1).expect(200);
    await request(app).get('/').set('X-Forwarded-For', client2).expect(200);
  });
});

