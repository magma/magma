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
import puppeteer from 'puppeteer';

const ADD_NETWORK_SELECTOR = `//span[text()='Add Network']`;
const ADD_NETWORK_DIALOG = `//span[text()='Add Network']`;
const ADD_NETWORK_SAVE = `//span[text()='Save']`;

type AddParams = {
  name: string;
  desc: string;
  federation?: string;
};

export async function addFegNetwork(page: puppeteer.Page, params: AddParams) {
  // assuming that we are in administrative tools page
  await page.waitForXPath(ADD_NETWORK_SELECTOR);
  const buttonSelector = await page.$x(ADD_NETWORK_SELECTOR);
  await buttonSelector[0].click();

  await page.waitForXPath(ADD_NETWORK_DIALOG);

  await page.waitForSelector('input[name=networkId]');
  await page.click('input[name=networkId]');
  await page.type('input[name=networkId]', params.name);

  await page.waitForSelector('input[name=name]');
  await page.click('input[name=name]');
  await page.type('input[name=name]', params.name);

  await page.click('input[name=description]');
  await page.type('input[name=description]', params.desc);

  // select network type
  await page.click('#types');
  const [fegSelector] = await page.$x(`//span[text()='feg']`);
  await fegSelector.click();
  await page.waitForFunction(
    'document.getElementById("types").value === "feg"',
  );

  //  TODO  need to figure out why we need to add this delay
  await page.waitForTimeout(500);
  const [saveButton] = await page.$x(ADD_NETWORK_SAVE);
  await saveButton.click();
  await page.waitForXPath(`//td[text()='${params.name}']`);
}

export async function addFegLteNetwork(
  page: puppeteer.Page,
  params: AddParams,
) {
  // assuming that we are in administrative tools page
  await page.waitForXPath(ADD_NETWORK_SELECTOR);
  const buttonSelector = await page.$x(ADD_NETWORK_SELECTOR);
  await buttonSelector[0].click();

  await page.waitForXPath(ADD_NETWORK_DIALOG);

  await page.waitForSelector('input[name=networkId]');
  await page.click('input[name=networkId]');
  await page.type('input[name=networkId]', params.name);

  await page.waitForSelector('input[name=name]');
  await page.click('input[name=name]');
  await page.type('input[name=name]', params.name);

  await page.click('input[name=description]');
  await page.type('input[name=description]', params.desc);

  // select network type
  await page.click('#types');
  const [fegSelector] = await page.$x(`//span[text()='feg_lte']`);
  await fegSelector.click();
  await page.waitForFunction(
    'document.getElementById("types").value === "feg_lte"',
  );

  if (params.federation) {
    await page.waitForSelector('input[name=fegNetworkID]');
    await page.click('input[name=fegNetworkID]', {delay: 500});
    await page.type('input[name=fegNetworkID]', params.federation);
  }

  // need to figure out why we need to add this delay
  await page.waitForTimeout(500);
  const [saveButton] = await page.$x(ADD_NETWORK_SAVE);
  await saveButton.click();

  await page.waitForXPath(`//td[text()='${params.name}']`);
}
