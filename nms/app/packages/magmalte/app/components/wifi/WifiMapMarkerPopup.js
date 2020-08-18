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

import type {MagmaGatewayFeature} from '@fbcnms/ui/insights/map/GeoJSON';

import React from 'react';

import {makeStyles} from '@material-ui/styles';

type Props = {
  feature: MagmaGatewayFeature,
};

const useStyles = makeStyles(() => ({
  container: {
    padding: 12,
  },
  detailList: {
    margin: '0px',
  },
  // we want to remove the padding of the MapBox popup container because it
  // intefers with the `mouseleave` event listener we install on the popup
  // element rendered inside the container
  '@global': {
    'div.mapboxgl-popup-content': {padding: 0},
  },
}));

export default function WifiMapMarkerPopup({feature}: Props) {
  const classes = useStyles();
  const device = feature.properties.device;
  if (!device) {
    return <></>;
  }
  return (
    <div className={classes.container}>
      <b>ID: </b>
      {device.id}
      <br />
      <b>Info: </b>
      {device.info}
      <br />
      {device.status && (
        <>
          <b>Mesh IP: </b>
          <ul className={classes.detailList}>
            {(device.status.meta?.mesh0_ip || '')
              .split(',')
              .filter(ip => ip !== '')
              .map((ip, i) => (
                <li key={i}>{ip}</li>
              ))}
          </ul>
        </>
      )}
    </div>
  );
}
