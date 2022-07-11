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

import puppeteer, {Browser} from 'puppeteer';
import {ARTIFACTS_DIR, SimulateNMSLogin} from '../LoginUtils';

let browser: Browser;
beforeEach(async () => {
  jest.setTimeout(60000);
  browser = await puppeteer.launch({
    args: ['--ignore-certificate-errors', '--window-size=1920,1080'],
    headless: true,
    // @ts-ignore
    defaultViewport: null,
  });
});

afterEach(async () => {
  await browser.close();
});

describe('NMS dashboard', () => {
  test('NMS loads correctly', async () => {
    const page = await browser.newPage();
    try {
      await SimulateNMSLogin(page);
    } catch (err) {
      await page.screenshot({path: ARTIFACTS_DIR + 'failed.png'});
      await browser.close();
      throw err;
    }

    await page.screenshot({path: ARTIFACTS_DIR + 'dashboard.png'});
    await browser.close();
  }, 60000);
});
