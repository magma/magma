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

import type {magmad_gateway} from '@fbcnms/magma-api';

export type MagmaFeatureCollection = {
  type: 'FeatureCollection',
  features: MagmaGatewayFeature[],
};

export type MagmaGatewayFeature = {
  type: 'Feature',
  geometry: {
    type: 'Point',
    coordinates: [number, number],
  },
  properties: {
    id: string | number,
    name?: string,
    iconSize: IconSize,
    gateway?: magmad_gateway,
    // $FlowFixMe[unclear-type] TODO(andreilee): migrated from fbcnms-ui
    [key: string]: any,
  },
};

export type MagmaConnectionFeature = {
  type: 'Feature',
  geometry: {
    type: 'LineString',
    coordinates: Array<[number, number]>,
  },
  properties: {
    id: string | number,
    name: string,
    // $FlowFixMe[unclear-type] TODO(andreilee): migrated from fbcnms-ui
    [string]: any,
  },
};

export type IconSize = [number, number];
