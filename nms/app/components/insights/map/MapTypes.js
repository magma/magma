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
 * @flow local
 * @format
 */

import type {MagmaGatewayFeature} from './GeoJSON';

import mapboxgl from 'mapbox-gl';

export type MapMarkerProps = {
  // $FlowFixMe[value-as-type] TODO(andreilee): migrated from fbcnms-ui
  map: mapboxgl.Map,
  feature: MagmaGatewayFeature,
  onClick?: (string | number) => void,
  showLabel?: boolean,
};
