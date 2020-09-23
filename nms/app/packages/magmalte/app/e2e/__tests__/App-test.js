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

const puppeteer = require('puppeteer');

const user = {
  email: 'admin@magma.test',
  passwd: 'password1234',
};

const DASHBOARD_SELECTOR = `//span[text()='Dashboard']`;
const LOGINFORM_SELECTOR = `//span[text()='Magma']`;
const ARTIFACTS_DIR = `/tmp/nms_artifacts/`;

describe('NMS dashboard', () => {
  test('NMS loads correctly', async () => {
    const browser = await puppeteer.launch({
      args: ['--ignore-certificate-errors'],
      headless: false,
      defaultViewport: null,
    });
    const page = await browser.newPage();
    try {
      await page.goto('https://magma-test.localhost/');
      await page.waitForXPath(LOGINFORM_SELECTOR);
      await page.click('input[name=email]');
      await page.type('input[name=email]', user.email);

      await page.click('input[name=password]');
      await page.type('input[name=password]', user.passwd);

      await page.click('button');
      await page.waitForXPath(DASHBOARD_SELECTOR, {
        timeout: 15000,
      });
    } catch (err) {
      await page.screenshot({path: ARTIFACTS_DIR + 'failed.png'});
      browser.close();
      throw err;
    }

    await page.screenshot({path: ARTIFACTS_DIR + 'dashboard.png'});
    browser.close();
  }, 30000);
});
