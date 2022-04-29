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

const KB = 1024;
const MB = 1024 * 1024;
const GB = 1024 * 1024 * 1024;

export const sortLexicographically = (a: string, b: string) =>
  a.localeCompare(b, 'en', {numeric: true});

export const sortMixed = <T: string | number>(a: ?T, b: ?T) => {
  if (a == null && b == null) {
    return 0;
  }
  if (a == null) {
    return -1;
  }
  if (b == null) {
    return 1;
  }
  if (typeof a == 'number' && typeof b == 'number') {
    return a - b;
  }
  return sortLexicographically(a.toString(), b.toString());
};

export const formatFileSize = (sizeInBytes: number) => {
  if (sizeInBytes === 0) {
    return '0MB';
  }

  if (sizeInBytes >= GB) {
    return `${(sizeInBytes / GB).toFixed(2)}GB`;
  } else if (sizeInBytes >= MB) {
    return `${(sizeInBytes / MB).toFixed(2)}MB`;
  } else if (sizeInBytes >= KB) {
    return `${Math.round(sizeInBytes / KB)}KB`;
  } else {
    return `${sizeInBytes}B`;
  }
};

export const isJSON = (text: ?string): boolean => {
  if (!text) {
    return false;
  }
  try {
    JSON.parse(text);
  } catch (e) {
    return false;
  }
  return true;
};

// formats server side timestamps (seonds from epoch)
// to text input required format dd-mm-yyyy
export const formatDateForTextInput = (dateValue: ?string) => {
  return !!dateValue ? dateValue.split('T')[0] : '';
};

export const formatMultiSelectValue = (
  options: Array<{value: string, label: string}>,
  value: string,
) => options.find(option => option.value === value)?.label;

export function hexToRgb(hexColor: string) {
  hexColor = hexColor.substr(1);

  const re = new RegExp(`.{1,${hexColor.length / 3}}`, 'g');
  let colors = hexColor.match(re);

  if (colors && colors[0].length === 1) {
    colors = colors.map(n => n + n);
  }

  return colors ? colors.map(n => parseInt(n, 16)).join(',') : '';
}
