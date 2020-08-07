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
import type {MagmaGatewayFeature} from '@fbcnms/ui/insights/map/GeoJSON';

import React from 'react';
import {makeStyles} from '@material-ui/styles';

type Props = {
  feature: MagmaGatewayFeature,
  showLabel?: boolean,
  useLargeIcon?: boolean,
};

const useStyles = makeStyles(() => ({
  markerContainer: {
    cursor: 'pointer',
    width: '15px',
    height: '15px',
  },
  largeMarkerContainer: {
    cursor: 'pointer',
    width: '25px',
    height: '25px',
  },
  deviceLabel: {
    position: 'absolute',
    marginTop: '-10px',
    whiteSpace: 'nowrap',
  },
  iconStyle: {
    stroke: 'black',
    strokeOpacity: 0.8,
    fillOpacity: 0.45,
  },
}));

const CircleIcon = ({className = '', color = '#fff'}) => (
  <svg viewBox="-2 -2 4 4">
    <circle
      cx="0"
      cy="0"
      r="1.8"
      className={className}
      fill={color}
      strokeWidth={0.2}
    />
  </svg>
);

const DiamondIcon = ({className = '', color = '#fff'}) => (
  <svg viewBox="-2 -2 4 4">
    <path
      d="M -2,0 0,-2 2,0 0,2 z"
      className={className}
      fill={color}
      strokeWidth={0.5}
    />
  </svg>
);

export default function WifiMapMarkerIcon({
  feature,
  showLabel,
  useLargeIcon,
}: Props) {
  const classes = useStyles();
  let color = 'grey';
  const device = feature.properties.device;
  let Icon = CircleIcon;
  if (device && device.status) {
    if (device.upDanger) {
      color = 'orange';
    } else if (device.up) {
      color = 'lime';
    } else {
      color = 'red';
    }

    if (device.isGateway) {
      Icon = DiamondIcon;
    }
  }

  return (
    <div
      className={
        useLargeIcon ? classes.largeMarkerContainer : classes.markerContainer
      }>
      <Icon className={classes.iconStyle} color={color} />
      {showLabel && (
        <div className={classes.deviceLabel}>
          {feature.properties.device?.info}
        </div>
      )}
    </div>
  );
}
