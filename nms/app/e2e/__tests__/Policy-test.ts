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
  const page = await browser.newPage();
  await SimulateNMSLogin(page);
});

afterEach(async () => {
  await browser.close();
});

describe('NMS', () => {
  test('verifying policy dashboard', async () => {
    const page = await browser.newPage();
    try {
      await page.goto('https://magma-test.localhost/nms/test/traffic/policy');

      // check if the description is right
      await page.waitForXPath(`//span[text()='Policies']`);

      // verify if we have mocked policies installed
      await page.waitForXPath(`//button[text()='test1']`);
      await page.waitForXPath(`//button[text()='test2']`);
    } catch (err) {
      await page.screenshot({
        path: ARTIFACTS_DIR + 'policy_dashboard_failed.png',
      });
      await page.close();
      throw err;
    }

    await page.screenshot({
      path: ARTIFACTS_DIR + 'policy_dashboard.png',
    });
  }, 60000);
});

// TODO (andreilee): Get these tests working again without flakiness
// describe('NMS Policy Add', () => {
//   test('verifying policy dashboard', async () => {
//     const page = await browser.newPage();
//     try {
//       await page.goto('https://magma-test.localhost/nms/test/traffic/policy');

//       // check if the description is right
//       await page.waitForXPath(`//span[text()='Policies']`);

//       const buttonSelector = await page.$x(`//span[text()='Create New']`);
//       buttonSelector[0].click();

//       const policyButtonSelector = '[data-testid="newPolicyMenuItem"]';
//       await page.waitForSelector(policyButtonSelector);
//       page.click(policyButtonSelector);

//       await page.waitForXPath(`//span[text()='Add New Policy']`);

//       // verify if we have policy info tab active
//       await page.waitForXPath(`//span[text()='Basic policy rule fields']`);
//       const policySelector = '[data-testid="policyID"]';
//       const prioSelector = '[data-testid="policyPriority"]';

//       // add policy information attributes
//       await page.waitForSelector(policySelector);
//       await page.click(policySelector);
//       await page.type(policySelector, 'test_policy0');
//       await page.waitForSelector(prioSelector);
//       await page.click(prioSelector);
//       await page.type(prioSelector, '10');

//       let saveButtonSelector = await page.$x(
//         `//span[text()='Save And Continue']`,
//       );
//       saveButtonSelector[0].click();
//       await page.waitForXPath(`//span[text()='Policy saved successfully']`);

//       // wait for flow tab
//       await page.waitForXPath(
//         `//span[text()="A policy's flows determines how it routes traffic"]`,
//       );

//       const addFlowSelector = '[data-testid="addFlowButton"]';
//       await page.waitForSelector(addFlowSelector);
//       await page.click(addFlowSelector);
//       const ipSrc = '[data-testid="ipSrc"]';
//       const ipDst = '[data-testid="ipDest"]';

//       // add policy information attributes
//       await page.waitForSelector(ipSrc);
//       await page.click(ipSrc);
//       await page.type(ipSrc, '1.1.1.1');
//       await page.waitForSelector(ipDst);
//       await page.click(ipDst);
//       await page.type(ipDst, '2.2.2.2');

//       saveButtonSelector = await page.$x(`//span[text()='Save And Continue']`);
//       saveButtonSelector[0].click();
//       await page.waitForXPath(`//span[text()='Policy saved successfully']`);
//     } catch (err) {
//       await page.screenshot({path: ARTIFACTS_DIR + 'policy_add_failed.png'});
//       await page.close();
//       throw err;
//     }

//     await page.screenshot({
//       path: ARTIFACTS_DIR + 'policy_add.png',
//     });
//   }, 60000);
// });

// describe('NMS Policy Edit', () => {
//   test('verifying policy dashboard', async () => {
//     const page = await browser.newPage();
//     try {
//       await page.goto('https://magma-test.localhost/nms/test/traffic/policy');

//       // check if the description is right
//       await page.waitForXPath(`//span[text()='Policies']`);

//       const buttonSelector = await page.$x(`//button[text()='test1']`);
//       buttonSelector[0].click();

//       await page.waitForXPath(`//span[text()='Edit Policy']`);

//       // verify if we have policy info tab active
//       await page.waitForXPath(`//span[text()='Basic policy rule fields']`);
//       const prioSelector = '[data-testid="policyPriority"]';

//       // add policy information attributes
//       await page.waitForSelector(prioSelector);
//       await page.click(prioSelector);
//       await page.evaluate(prioSelector => {
//         document.querySelector(prioSelector).value = '';
//       }, prioSelector);
//       await page.type(prioSelector, '10');

//       const saveButtonSelector = await page.$x(`//span[text()='Save']`);
//       saveButtonSelector[0].click();

//       await page.waitForXPath(`//span[text()='Policy saved successfully']`);
//     } catch (err) {
//       await page.screenshot({path: ARTIFACTS_DIR + 'policy_edit_failed.png'});
//       await page.close();
//       throw err;
//     }

//     await page.screenshot({
//       path: ARTIFACTS_DIR + 'policy_edit.png',
//     });
//   }, 60000);
// });
