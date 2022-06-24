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
import {addFegLteNetwork, addFegNetwork} from '../NetworkUtils';

const ADMIN_SELECTOR = `//span[text()='Administration']`;
const ADMIN_NW_SELECTOR = `//a[starts-with(@href, '/nms/test/admin/networks')]`;
const PROFILE_BUTTON_SELECTOR = `//*[@data-testid='profileButton']`;
const NETWORK_SELECTOR_SELECTOR = `//*[@data-testid='networkSelector']`;

let browser: Browser;
beforeEach(async () => {
  jest.setTimeout(60000);
  browser = await puppeteer.launch({
    args: ['--ignore-certificate-errors', '--window-size=1920,1080'],
    headless: true,
    // @ts-ignore
    defaultViewport: null,
  });
  const page = await browser.newPage();
  await SimulateNMSLogin(page);
});

afterEach(async () => {
  await browser.close();
});

describe('Admin component', () => {
  test('verifying addition of feg_lte networks', async () => {
    const page = await browser.newPage();
    await page.setViewport({width: 1280, height: 1024});
    try {
      await page.goto('https://magma-test.localhost/nms/test');
      await page.waitForXPath(`//span[text()='Dashboard']`, {
        timeout: 15000,
      });
      await page.waitForXPath(PROFILE_BUTTON_SELECTOR);
      const navSelector = await page.$x(PROFILE_BUTTON_SELECTOR);
      await navSelector[0].click();

      const adminSelector = await page.$x(ADMIN_SELECTOR);
      await page.waitForTimeout(500);
      await adminSelector[0].click();

      await page.waitForXPath(ADMIN_NW_SELECTOR);

      // It seems to take a while until the button becomes responsive
      await page.waitForTimeout(500);
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
  test('verifying LTE network addition from network selector', async () => {
    const page = await browser.newPage();
    try {
      // test_feg_lte_network is mocked out
      await page.goto('https://magma-test.localhost/nms/test');
      await page.waitForXPath(`//span[text()='Dashboard']`, {
        timeout: 15000,
      });

      await page.waitForXPath(NETWORK_SELECTOR_SELECTOR);
      const networkSelector = await page.$x(NETWORK_SELECTOR_SELECTOR);
      await networkSelector[0].click();

      await page.waitForXPath(`//span[text()='Create Network']`);
      const buttonSelector = await page.$x(`//span[text()='Create Network']`);
      await page.waitForTimeout(500);
      await buttonSelector[0].click();

      const networkIDSelector = '[data-testid="networkID"]';
      const networkNameSelector = '[data-testid="networkName"]';
      const networkDescriptionSelector = '[data-testid="networkDescription"]';

      // add network information attributes
      await page.waitForSelector(networkIDSelector);
      await page.click(networkIDSelector);
      await page.type(networkIDSelector, 'test_network');

      await page.waitForSelector(networkNameSelector);
      await page.click(networkNameSelector);
      await page.type(networkNameSelector, 'test_network');

      await page.waitForSelector(networkDescriptionSelector);
      await page.click(networkDescriptionSelector);
      await page.type(networkDescriptionSelector, 'test_network');

      // TODO need to figure out why we need to add this delay
      await page.waitForTimeout(500);

      const saveButtonSelector = '[data-testid="saveButton"]';
      await page.waitForSelector(saveButtonSelector);
      await page.click(saveButtonSelector);

      await page.waitForXPath(
        `//span[text()='Network test_network successfully created']`,
      );

      // use epc defaults
      const tacSelector = '[data-testid="tac"]';
      await page.waitForSelector(tacSelector);
      await page.click(tacSelector);
      await page.evaluate(tacSelector => {
        // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
        document.querySelector(tacSelector).value = '';
      }, tacSelector);
      await page.type(tacSelector, '2');

      const epcSaveButtonSelector = '[data-testid="epcSaveButton"]';
      await page.waitForSelector(epcSaveButtonSelector);
      await (await page.$(epcSaveButtonSelector))!.press('Enter');
      await page.waitForXPath(
        `//span[text()='EPC configs saved successfully']`,
      );

      const earfcndlSelector = '[data-testid="earfcndl"]';
      await page.waitForSelector(earfcndlSelector);
      await page.click(earfcndlSelector);
      await page.evaluate(earfcndlSelector => {
        // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
        document.querySelector(earfcndlSelector).value = '';
      }, earfcndlSelector);

      await page.type(earfcndlSelector, '44592');
      const ranSaveButtonSelector = '[data-testid="ranSaveButton"]';
      await page.waitForSelector(ranSaveButtonSelector);
      await (await page.$(ranSaveButtonSelector))!.press('Enter');

      await page.waitForXPath(
        `//span[text()='RAN configs saved successfully']`,
      );
    } catch (err) {
      await page.screenshot({
        path: ARTIFACTS_DIR + 'failed_networkselector_add.png',
      });
      await page.close();
      throw err;
    }
  }, 60000);

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

      const fedationTabSelector = '[data-testid="federationTab"]';
      await page.waitForSelector(fedationTabSelector);
      await page.click(fedationTabSelector);

      const fegPlaceholder = '[placeholder="Enter Federation Network ID"]';
      await page.waitForSelector(fegPlaceholder);
      await page.click(fegPlaceholder);

      await page.evaluate(selector => {
        // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
        document.querySelector(selector).value = '';
      }, fegPlaceholder);

      await page.type(fegPlaceholder, 'test_feg_network2');

      // @karthiksubraveti - TODO need to figure out why we need to add this delay
      await page.waitForTimeout(500);
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
