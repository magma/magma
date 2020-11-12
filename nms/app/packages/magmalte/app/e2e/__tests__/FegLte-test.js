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
import {ARTIFACTS_DIR, SimulateNMSLogin} from '../LoginUtils';
import {addFegLteNetwork, addFegNetwork} from '../NetworkUtils';

const ADMIN_SELECTOR = `//span[text()='Administrative Tools']`;
const ADMIN_NW_SELECTOR = `//a[starts-with(@href, '/admin/networks')]`;
const NAV_SELECTOR = `//body/div[1]/div/div/div[last()]`;

let browser;
beforeEach(async () => {
  jest.setTimeout(60000);
  browser = await puppeteer.launch({
    args: ['--ignore-certificate-errors'],
    headless: true,
    defaultViewport: null,
  });
  const page = await browser.newPage();
  await SimulateNMSLogin(page);
});

afterEach(() => {
  browser.close();
});

describe('Admin component', () => {
  test('verifying addition of feg_lte networks', async () => {
    const page = await browser.newPage();
    try {
      await page.goto('https://magma-test.localhost/');
      await page.waitForXPath(`//span[text()='Dashboard']`, {
        timeout: 15000,
      });
      await page.waitForXPath(NAV_SELECTOR);
      const navSelector = await page.$x(NAV_SELECTOR);
      await navSelector[0].click();

      const adminSelector = await page.$x(ADMIN_SELECTOR);
      await adminSelector[0].click();

      // wait for 'admin network page'
      await page.waitForNavigation();
      await page.waitForXPath(ADMIN_NW_SELECTOR);
      const adminNwSelector = await page.$x(ADMIN_NW_SELECTOR);
      await adminNwSelector[0].click();

      const fegNetwork = {
        name: 'test_feg_network2',
        desc: 'Test Feg Network Description',
      };
      await addFegNetwork(page, fegNetwork);

      const fegLteNetwork = {
        name: 'test_feg_lte_network2',
        desc: 'Test Feg LTE Network Description',
        federation: fegNetwork.name,
      };

      await addFegLteNetwork(page, fegLteNetwork);
    } catch (err) {
      await page.screenshot({path: ARTIFACTS_DIR + 'failed.png'});
      throw err;
    }
    await page.screenshot({
      path: ARTIFACTS_DIR + 'organization_network_list.png',
    });
    await page.close();
  }, 60000);
});

describe('NMS', () => {
  test('verifying feg_lte dashboard', async () => {
    const page = await browser.newPage();
    try {
      // test_feg_lte_network is mocked out
      await page.goto(
        'https://magma-test.localhost/nms/test_feg_lte_network/network/network',
      );

      // check if the description is right
      await page.waitForXPath(`//span[text()='Network']`);
      await page.waitForXPath(`//span[text()='test_feg_lte_network']`);
      await page.waitForXPath(
        `//span[text()='Test Feg LTE Network Description']`,
      );
      await page.waitForXPath(`//span[text()='test_feg_network']`);

      // edit description
      const editSelector = '[data-testid="infoEditButton"]';
      await page.waitForSelector(editSelector);
      await page.click(editSelector);

      const fegPlaceholder = '[placeholder="Enter Federation Network ID"]';
      await page.waitForSelector(fegPlaceholder);
      await page.click(fegPlaceholder);

      await page.evaluate(selector => {
        document.querySelector(selector).value = '';
      }, fegPlaceholder);

      await page.type(fegPlaceholder, 'test_feg_network2');

      // TODO need to figure out why we need to add this delay
      await page.waitFor(500);
      const [saveButton] = await page.$x(`//span[text()='Save']`);
      await saveButton.click();

      await page.waitForXPath(`//span[text()='Network']`);
    } catch (err) {
      await page.screenshot({path: ARTIFACTS_DIR + 'failed.png'});
      await page.close();
      throw err;
    }

    await page.screenshot({
      path: ARTIFACTS_DIR + 'feg_lte_network_dashboard.png',
    });
  }, 60000);
});
