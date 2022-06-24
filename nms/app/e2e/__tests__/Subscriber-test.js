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
  test('verifying subscriber dashboard', async () => {
    const page = await browser.newPage();
    try {
      await page.goto(
        'https://magma-test.localhost/nms/test/subscribers/overview',
      );

      // check if the description is right
      await page.waitForXPath(`//span[text()='Subscribers']`);

      // verify if we have mocked apns installed
      await page.waitForXPath(`//button[text()='IMSI001010002220018']`);
    } catch (err) {
      await page.screenshot({
        path: ARTIFACTS_DIR + 'subscriber_dashboard_failed.png',
      });
      await page.close();
      throw err;
    }

    await page.screenshot({
      path: ARTIFACTS_DIR + 'subscriber_dashboard.png',
    });
  }, 60000);
});

// describe('NMS Subscriber Add', () => {
//   test('verifying subscriber dashboard', async () => {
//     const page = await browser.newPage();
//     try {
//       await page.goto(
//         'https://magma-test.localhost/nms/test/subscribers/overview',
//       );

//       // check if the description is right
//       await page.waitForXPath(`//span[text()='Subscribers']`);

//       //Add new subscriber
//       const addSubscriberSelector = await page.$x(
//         `//span[text()='Add Subscriber']`,
//       );
//       addSubscriberSelector[0].click();
//       await page.waitForXPath(`//span[text()='Add Subscribers']`);

//       // Add subscriber
//       const addButtonSelector = await page.$x(`//button[@title='Add']`);
//       addButtonSelector[0].click();

//       // Add subscriber information attributes
//       const name = '[data-testid="name"]';
//       await page.waitForSelector(name);
//       await page.click(name);
//       await page.type(name, 'IMSI001010002220022');

//       const IMSI = '[data-testid="IMSI"]';
//       await page.waitForSelector(IMSI);
//       await page.click(IMSI);
//       await page.type(IMSI, 'IMSI001010002220022');

//       const authKey = '[data-testid="authKey"]';
//       await page.waitForSelector(authKey);
//       await page.click(authKey);
//       await page.type(authKey, '8baf473f2f8fd09487cccbd7097c6862');

//       const authOpc = '[data-testid="authOpc"]';
//       await page.waitForSelector(authOpc);
//       await page.click(authOpc);
//       await page.type(authOpc, '8e27b6af0e692e750f32667a3b14605d');

//       // Add subscriber to save
//       const saveSubscriber = await page.$x(`//button[@title='Save']`);
//       saveSubscriber[0].click();

//       // Save subscriber
//       const saveNewSubscriber = '[data-testid="saveSubscriber"]';
//       await page.waitForSelector(saveNewSubscriber);
//       await page.click(saveNewSubscriber);
//       await page.waitForXPath(
//         `//span[text()='Subscriber(s) saved successfully']`,
//       );
//     } catch (err) {
//       await page.screenshot({
//         path: ARTIFACTS_DIR + 'subscriber_add_failed.png',
//       });
//       await page.close();
//       throw err;
//     }

//     await page.screenshot({
//       path: ARTIFACTS_DIR + 'subscriber_add.png',
//     });
//   }, 60000);
// });

// describe('NMS Subscriber Edit', () => {
//   test('verifying subscriber dashboard', async () => {
//     const page = await browser.newPage();
//     try {
//       await page.goto(
//         'https://magma-test.localhost/nms/test/subscribers/overview/config/IMSI001010002220018/config',
//       );
//       // edit subscriber info
//       const editButtonSelector = '[data-testid="subscriber"]';
//       await page.waitForSelector(editButtonSelector);
//       await page.click(editButtonSelector);

//       const subscriberName = '[data-testid="name"]';
//       await page.waitForSelector(subscriberName);
//       await page.click(subscriberName, {clickCount: 3});
//       await page.type(subscriberName, 'test_subscriber');

//       // Save subscriber
//       const saveSubscriberInfo = '[data-testid="subscriber-saveButton"]';
//       await page.waitForSelector(saveSubscriberInfo);
//       await page.click(saveSubscriberInfo);

//       await page.waitForXPath(`//span[text()='Subscriber saved successfully']`);

//       // edit subscriber traffic policy
//       // const trafficPolicySelector = '[data-testid="trafficPolicy"]';
//       // await page.waitForSelector(trafficPolicySelector);
//       // await page.click(trafficPolicySelector);

//       // const activeApns = '#activeApnTestId';
//       // await page.waitForSelector(activeApns);
//       // await page.click(activeApns);

//       // const apnCheckbox = 'input[type="checkbox"]';
//       // await page.waitForSelector(apnCheckbox);
//       // await page.click(apnCheckbox);

//       // Save subscriber
//       // const saveSubscriberTrafficPolicy =
//       //   '[data-testid="trafficPolicy-saveButton"]';
//       // await page.waitForSelector(saveSubscriberTrafficPolicy);
//       // await page.click(saveSubscriberTrafficPolicy, {clickCount: 2});

//       // await page.waitForXPath(`//span[text()='Subscriber saved successfully']`);
//     } catch (err) {
//       await page.screenshot({
//         path: ARTIFACTS_DIR + 'subscriber_edit_failed.png',
//       });
//       await page.close();
//       throw err;
//     }
//     await page.screenshot({
//       path: ARTIFACTS_DIR + 'subscriber_edit.png',
//     });
//   }, 60000);
// });
