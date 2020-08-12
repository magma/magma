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

import type {MagmaConnectionFeature} from '@fbcnms/ui/insights/map/GeoJSON';

import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Typography from '@material-ui/core/Typography';

import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  linkDetails: {
    margin: '8px',
    whiteSpace: 'nowrap',
  },
}));

type Props = {
  features: ?Array<MagmaConnectionFeature>,
};

export default function WifiMapFeatureDetail(props: Props) {
  const classes = useStyles();

  // features should all be lines
  const {features} = props;

  if (!features || features.length == 0) {
    return null;
  }

  // render only first element. TODO: smarter selection of feature
  const feature = features[0];
  if (!feature || !feature.properties) {
    return null;
  }
  const fProps = feature.properties;

  return (
    <>
      <Typography className={classes.linkDetails} variant="body2">
        <b>Connection Type: </b>
        {fProps.highestConnectionType}
        {fProps.unidirectional ? ' (unidirectional)' : ''}
        <br />
        <b>0: </b>
        {fProps.deviceInfo0} ({fProps.deviceId0})
        <br />
        <b>1: </b>
        {fProps.deviceInfo1} ({fProps.deviceId1})
        <br />
      </Typography>

      <Table>
        <TableHead>
          <TableRow>
            <TableCell />
            <TableCell>device0 view of device1</TableCell>
            <TableCell>device1 view of device0</TableCell>
          </TableRow>
        </TableHead>

        <TableBody>
          {(fProps.info0to1_isDefaultRoute ||
            fProps.info1to0_isDefaultRoute) && (
            <TableRow>
              <TableCell component="th">Default route</TableCell>
              <TableCell>
                {fProps.info0to1_isDefaultRoute ? 'yes' : ''}
              </TableCell>
              <TableCell>
                {fProps.info1to0_isDefaultRoute ? 'yes' : ''}
              </TableCell>
            </TableRow>
          )}

          <TableRow>
            <TableCell component="th">OpenR metric</TableCell>
            <TableCell>{fProps.info0to1_OpenrMetric}</TableCell>
            <TableCell>{fProps.info1to0_OpenrMetric}</TableCell>
          </TableRow>

          <TableRow>
            <TableCell component="th">L2 plink</TableCell>
            <TableCell>{fProps.info0to1_L2MeshPlink}</TableCell>
            <TableCell>{fProps.info1to0_L2MeshPlink}</TableCell>
          </TableRow>

          <TableRow>
            <TableCell component="th">L2 signal</TableCell>
            <TableCell>{fProps.info0to1_L2Signal}</TableCell>
            <TableCell>{fProps.info1to0_L2Signal}</TableCell>
          </TableRow>

          <TableRow>
            <TableCell component="th">L2 metric</TableCell>
            <TableCell>{fProps.info0to1_L2Metric}</TableCell>
            <TableCell>{fProps.info1to0_L2Metric}</TableCell>
          </TableRow>

          <TableRow>
            <TableCell component="th">L2 Tx Expected</TableCell>
            <TableCell>{fProps.info0to1_L2ExpectedThroughput}</TableCell>
            <TableCell>{fProps.info1to0_L2ExpectedThroughput}</TableCell>
          </TableRow>

          <TableRow>
            <TableCell component="th">L2 Rx Bitrate</TableCell>
            <TableCell>{fProps.info0to1_L2RxBitrate}</TableCell>
            <TableCell>{fProps.info1to0_L2RxBitrate}</TableCell>
          </TableRow>

          <TableRow>
            <TableCell component="th">L2 inactive t</TableCell>
            <TableCell>{fProps.info0to1_L2InactiveTime}</TableCell>
            <TableCell>{fProps.info1to0_L2InactiveTime}</TableCell>
          </TableRow>
        </TableBody>
      </Table>
      <br />
    </>
  );
}
