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

const DASHBOARD_SELECTOR = `//span[text()='Dashboard']`;
const LOGINFORM_SELECTOR = `//span[text()='Magma']`;
export const ARTIFACTS_DIR = `/tmp/nms_artifacts/`;

const user = {
  email: 'admin@magma.test',
  passwd: 'password1234',
};
export async function SimulateNMSLogin(page: puppeteer.Page) {
  await page.goto('https://magma-test.localhost/nms');
  await page.waitForXPath(LOGINFORM_SELECTOR);
  await page.click('input[name=email]');
  await page.type('input[name=email]', user.email);

  await page.click('input[name=password]');
  await page.type('input[name=password]', user.passwd);

  await page.click('button');
  await page.waitForXPath(DASHBOARD_SELECTOR, {
    timeout: 15000,
  });
}
