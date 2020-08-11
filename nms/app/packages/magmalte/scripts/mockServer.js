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
 *
 * @flow
 * @format
 */

const jsonServer = require('json-server');
const https = require('https');
const fs = require('fs');

const keyFile = 'mock/.cache/mock_server.key';
const certFile = 'mock/.cache/mock_server.cert';

const server = jsonServer.create();
const router = jsonServer.router('mock/db.json');
const middlewares = jsonServer.defaults();
server.use(middlewares);

const buffer = fs.readFileSync('mock/db.json', 'utf-8');
const db = JSON.parse(buffer);

// add custom route handlers
server.get('/magma/v1/networks/test/gateways', (req, res) => {
  if (req.method === 'GET') {
    res.status(200).jsonp(db['networksFull']['test']['gateways']);
  }
});

server.get('/magma/v1/networks/test/type', (req, res) => {
  if (req.method === 'GET') {
    res.status(200).jsonp(db['networksFull']['test']['type']);
  }
});

server.get('/magma/v1/lte/test/gateways', (req, res) => {
  if (req.method === 'GET') {
    res.status(200).jsonp(db['lte']['gateways']);
  }
});

server.get('/magma/v1/lte/test/enodebs', (req, res) => {
  if (req.method === 'GET') {
    res.status(200).jsonp(db['lte']['enodebs']);
  }
});
server.use('/magma/v1', router);

https
  .createServer(
    {
      key: fs.readFileSync(keyFile),
      cert: fs.readFileSync(certFile),
    },
    server,
  )
  .listen(3001, '0.0.0.0', () => {
    console.log('Go to https://localhost:3001/');
  });
