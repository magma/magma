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
import express, {
  Request,
  Response,
  request as middlewareRequest,
} from 'express';
import {agent as request} from 'supertest';

import {rateLimitMiddleware} from '..';

const config = {
  RATE_LIMIT_CONFIG: {
    windowMs: 15 * 1000, // 15 seconds
    limit: 2,
  },
};

Object.defineProperty(global, 'performance', {
  writable: true,
});

jest.useFakeTimers();
jest.mock('../../../config/config', () => config);

describe('Rate limit test', () => {
  const app = express();
  app.use(rateLimitMiddleware);
  app.get('/', (_request: Request, res: Response) => res.send('Hi there!'));

  const client1 = '1.2.3.4';
  const client2 = '2.3.4.5';

  it('disallows request after 2 requests within 15 seconds for a given ip', async () => {
    jest.spyOn(middlewareRequest, 'ip', 'get').mockReturnValue(client1);
    await request(app).get('/').expect(200);

    jest.spyOn(middlewareRequest, 'ip', 'get').mockReturnValue(client2);
    await request(app).get('/').expect(200);

    jest.spyOn(middlewareRequest, 'ip', 'get').mockReturnValue(client1);
    await request(app).get('/').expect(200);

    jest.spyOn(middlewareRequest, 'ip', 'get').mockReturnValue(client2);
    await request(app).get('/').expect(200);

    jest.spyOn(middlewareRequest, 'ip', 'get').mockReturnValue(client1);
    await request(app).get('/').expect(429);

    jest.spyOn(middlewareRequest, 'ip', 'get').mockReturnValue(client2);
    await request(app).get('/').expect(429);

    jest.advanceTimersByTime(16_000);

    jest.spyOn(middlewareRequest, 'ip', 'get').mockReturnValue(client1);
    await request(app).get('/').expect(200);

    jest.spyOn(middlewareRequest, 'ip', 'get').mockReturnValue(client2);
    await request(app).get('/').expect(200);
  });
});
