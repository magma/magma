/**
 * Copyright 2022 The Magma Authors.
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

import {AlertRoutingTree} from '../AlarmAPIType';
import {filterRouteByRuleName, filterUpdatedFilterRoutes} from '../hooks';

const response = {
  receiver: 'test_tenant_base_route',
  match: {networkID: 'test'},
  routes: [
    {receiver: 'User1', match: {alertname: 'High Disk Usage Alert'}},
    {
      receiver: 'User1',
      match: {alertname: 'Certificate Expiring Soon'},
    },
    {
      receiver: 'User1',
      match: {alertname: 'Bootstrap Exception Alert'},
    },
    {
      receiver: 'User2',
      match: {alertname: 'Gateway Checkin Failure'},
    },
  ],
};

describe('Testing filterRouteByRuleName', () => {
  test('Remove only High Disk Usage Alert from User1', () => {
    const ruleName = 'High Disk Usage Alert';
    const initialReceiver = 'User1';
    const routes = filterRouteByRuleName(response, initialReceiver, ruleName);
    expect(routes.length).toBe(3);
  });

  test('Remove only Gateway Checkin Failure from User2', () => {
    const ruleName = 'Gateway Checkin Failure';
    const initialReceiver = 'User2';
    const routes = filterRouteByRuleName(response, initialReceiver, ruleName);
    expect(routes.length).toBe(3);
  });

  test('Return all routes if initialReceiver is empty', () => {
    const ruleName = 'Gateway Checkin Failure';
    const initialReceiver = '';
    const routes = filterRouteByRuleName(response, initialReceiver, ruleName);
    expect(routes.length).toBe(4);
  });

  test('Return all routes if ruleName is empty', () => {
    const ruleName = '';
    const initialReceiver = 'User2';
    const routes = filterRouteByRuleName(response, initialReceiver, ruleName);
    expect(routes.length).toBe(4);
  });
});

describe('Testing filterUpdatedFilterRoutes', () => {
  test('Update only High Disk Usage Alert from User1 to User2', () => {
    const ruleName = 'High Disk Usage Alert';
    const initialReceiver = 'User1';
    const receiver = 'User2';

    const routes: Array<AlertRoutingTree> = filterUpdatedFilterRoutes(
      response,
      initialReceiver,
      ruleName,
      receiver,
    );
    expect(routes[0].receiver).toBe(receiver);
    expect(routes[1].receiver).toBe(initialReceiver);
  });
});
