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

import puppeteer from 'puppeteer';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {ARTIFACTS_DIR, SimulateNMSLogin} from '../LoginUtils';

let browser;
beforeEach(async () => {
  jest.setTimeout(60000);
  browser = await puppeteer.launch({
    args: ['--ignore-certificate-errors', '--window-size=1920,1080'],
    headless: true,
    defaultViewport: null,
  });
});

afterEach(() => {
  browser.close();
});

describe('NMS dashboard', () => {
  test('NMS loads correctly', async () => {
    const page = await browser.newPage();
    try {
      await SimulateNMSLogin(page);
    } catch (err) {
      await page.screenshot({path: ARTIFACTS_DIR + 'failed.png'});
      browser.close();
      throw err;
    }

    await page.screenshot({path: ARTIFACTS_DIR + 'dashboard.png'});
    browser.close();
  }, 60000);
});
