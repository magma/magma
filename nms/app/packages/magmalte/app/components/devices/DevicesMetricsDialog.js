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

import AsyncMetric from '@fbcnms/ui/insights/AsyncMetric';
import Button from '@fbcnms/ui/components/design-system/Button';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import GridList from '@material-ui/core/GridList';
import GridListTile from '@material-ui/core/GridListTile';
import LoadingFillerBackdrop from '@fbcnms/ui/components/LoadingFillerBackdrop';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';

import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';
import {useRouter} from '@fbcnms/ui/hooks';

type Props = {
  onClose: () => void,
};

export default function DevicesMetricsDialog(props: Props) {
  const {match} = useRouter();
  const {deviceID, networkId} = match.params;
  const {isLoading, response} = useMagmaAPI(
    MagmaV1API.getNetworksByNetworkIdPrometheusSeries,
    {
      networkId,
    },
  );

  if (isLoading || !response) {
    return <LoadingFillerBackdrop />;
  }

  const metricIDs = new Set();
  response.forEach(item => {
    if (item.deviceID === deviceID) {
      metricIDs.add(item.__name__);
    }
  });

  return (
    <Dialog
      open={true}
      onClose={props.onClose}
      fullWidth={true}
      scroll="body"
      maxWidth="md">
      <DialogTitle>Device Metrics</DialogTitle>
      <DialogContent>
        <GridList cols={2} cellHeight={300}>
          {[...metricIDs].map((metric, i) => (
            <GridListTile key={i} cols={1}>
              <Card>
                <CardContent>
                  <Text component="h6" variant="h6">
                    {metric.replace(
                      '_openconfig_interfaces_interface_interface_',
                      '',
                    )}
                  </Text>
                  <div style={{height: 250}}>
                    <AsyncMetric
                      key={i}
                      label={metric}
                      unit=""
                      queries={[`${metric}{deviceID="${deviceID}"}`]}
                      timeRange="24_hours"
                    />
                  </div>
                </CardContent>
              </Card>
            </GridListTile>
          ))}
        </GridList>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose} skin="regular">
          Close
        </Button>
      </DialogActions>
    </Dialog>
  );
}
