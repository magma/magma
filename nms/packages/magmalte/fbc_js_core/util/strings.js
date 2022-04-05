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
 * @flow strict-local
 * @format
 */

export function hexToBase64(hexString: string): string {
  let parsedValue;

  parsedValue = hexString.toLowerCase();
  if (parsedValue.length % 2 === 1) {
    parsedValue = '0' + parsedValue;
  }
  // Raise an exception if any bad value is entered
  if (!isValidHex(hexString)) {
    throw new Error('is not valid hex');
  }
  return Buffer.from(parsedValue, 'hex').toString('base64');
}

export function base64ToHex(base64String: string): string {
  return Buffer.from(base64String, 'base64').toString('hex');
}

export function decodeBase64(base64String: string): string {
  return Buffer.from(base64String, 'base64').toString();
}

export function isValidHex(hexString: string): boolean {
  return hexString.match(/^[a-fA-F0-9]*$/) !== null;
}

export function capitalize(s: string) {
  return s.charAt(0).toUpperCase() + s.slice(1);
}

export function joinNullableStrings(
  strings: Array<?string>,
  separator?: ?string,
): ?string {
  const tooltipStrings = strings.filter(str => str != null);
  if (tooltipStrings.length === 0) {
    return undefined;
  }
  return tooltipStrings.filter(str => str != '').join(separator ?? undefined);
}
