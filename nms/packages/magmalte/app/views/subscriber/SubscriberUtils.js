/*
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

import type {promql_return_object} from '../../../generated/MagmaAPIBindings';

const mBIT = 1000000;
const kBIT = 1000;
export function getLabelUnit(val: number) {
  if (val > mBIT) {
    return [(val / mBIT).toFixed(2), 'mb'];
  } else if (val > kBIT) {
    return [(val / kBIT).toFixed(2), 'kb'];
  }
  return [val.toFixed(2), 'bytes'];
}

/**
 * Converts bits to megabits
 * @param {number} val The value in bits to be converted
 * @returns {string} Megabits value of the number passed in
 */
export function convertBitToMbit(val: number) {
  return (val / mBIT).toFixed(2);
}

export function getPromValue(resp: promql_return_object) {
  const respArr = resp?.data?.result
    ?.map(item => {
      return parseFloat(item?.value?.[1]);
    })
    .filter(Boolean);
  return respArr && respArr.length ? respArr[0] : 0;
}

// default subscriber count in get subscriber query
export const DEFAULT_PAGE_SIZE = 25;

// susbcriber export colums title
export const SUBSCRIBER_EXPORT_COLUMNS = [
  {
    title: 'Name',
    field: 'name',
  },
  {title: 'IMSI', field: 'id'},

  {title: 'Auth Key', field: 'auth_key'},
  {title: 'Auth OPC', field: 'auth_opc'},
  {title: 'Service', field: 'state'},
  {title: 'Data Plan', field: 'sub_profile'},
  {title: 'Active APNs', field: 'active_apns'},
];
export const SUBSCRIBER_ADD_ERRORS = Object.freeze({
  INVALID_IMSI:
    'The IMSI should be a string IMSI followed by a number with 10-15 digits',
  INVALID_AUTH_KEY:
    'Auth key is not a valid hex (example: 000102030405060708090A0B0C0D0E0F)',
  INVALID_AUTH_OPC:
    'Auth opc is not a valid hex (example: 000102030405060708090A0B0C0D0E0F)',
  REQUIRED_SUB_PROFILE: 'Please select a data plan',
  DUPLICATE_IMSI: 'The IMSI is duplicated',
  REQUIRED_AUTH_KEY: 'Auth key is required',
});
