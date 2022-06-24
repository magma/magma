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
  const page = await browser.newPage();
  await SimulateNMSLogin(page);
});

afterEach(() => {
  browser.close();
});

describe('NMS', () => {
  test('verifying apn dashboard', async () => {
    const page = await browser.newPage();
    try {
      await page.goto('https://magma-test.localhost/nms/test/traffic/apn');

      // check if the description is right
      await page.waitForXPath(`//span[text()='APNs']`);

      // verify if we have mocked apns installed
      await page.waitForXPath(`//button[text()='internet']`);
    } catch (err) {
      await page.screenshot({
        path: ARTIFACTS_DIR + 'apn_dashboard_failed.png',
      });
      await page.close();
      throw err;
    }

    await page.screenshot({
      path: ARTIFACTS_DIR + 'apn_dashboard.png',
    });
  }, 60000);
});

describe('NMS Apn Add', () => {
  test('verifying apn dashboard', async () => {
    const page = await browser.newPage();
    try {
      await page.goto('https://magma-test.localhost/nms/test/traffic/apn');

      // check if the description is right
      await page.waitForXPath(`//span[text()='APNs']`);

      const buttonSelector = await page.$x(`//span[text()='Create New APN']`);
      buttonSelector[0].click();

      await page.waitForXPath(`//span[text()='Add New APN']`);

      const apnSelector = '[data-testid="apnID"]';

      // add apn information attributes
      await page.waitForSelector(apnSelector);
      await page.click(apnSelector);
      await page.type(apnSelector, 'test_apn0');

      // ksubraveti : TODO need to figure out why we need to add this delay
      await page.waitForTimeout(500);
      const saveButtonSelector = await page.$x(`//span[text()='Save']`);
      saveButtonSelector[0].click();
      await page.waitForXPath(`//span[text()='APN saved successfully']`);
    } catch (err) {
      await page.screenshot({path: ARTIFACTS_DIR + 'apn_add_failed.png'});
      await page.close();
      throw err;
    }

    await page.screenshot({
      path: ARTIFACTS_DIR + 'apn_add.png',
    });
  }, 60000);
});

describe('NMS APN Edit', () => {
  test('verifying apn dashboard', async () => {
    const page = await browser.newPage();
    try {
      await page.goto('https://magma-test.localhost/nms/test/traffic/apn');

      // check if the description is right
      await page.waitForXPath(`//span[text()='APNs']`);

      const buttonSelector = await page.$x(`//button[text()='internet']`);
      buttonSelector[0].click();

      await page.waitForXPath(`//span[text()='Edit APN']`);

      const prioSelector = '[data-testid="apnPriority"]';

      // edit apn information attributes
      await page.waitForSelector(prioSelector);
      await page.click(prioSelector, {clickCount: 3});
      await page.type(prioSelector, '10');

      // ksubraveti : TODO need to figure out why we need to add this delay
      await page.waitForTimeout(500);
      const saveButtonSelector = await page.$x(`//span[text()='Save']`);
      saveButtonSelector[0].click();

      await page.waitForXPath(`//span[text()='APN saved successfully']`);
    } catch (err) {
      await page.screenshot({path: ARTIFACTS_DIR + 'apn_edit_failed.png'});
      await page.close();
      throw err;
    }

    await page.screenshot({
      path: ARTIFACTS_DIR + 'apn_edit.png',
    });
  }, 60000);
});
