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

import type {ContextRouter} from 'react-router-dom';
import type {Theme} from '@material-ui/core';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithStyles} from '@material-ui/core';

import ApnEditDialog from './ApnEditDialog';
import Button from '@fbcnms/ui/components/design-system/Button';
import DeleteIcon from '@material-ui/icons/Delete';
import EditIcon from '@material-ui/icons/Edit';
import IconButton from '@material-ui/core/IconButton';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import {Route, withRouter} from 'react-router-dom';

import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';
import {useState} from 'react';
import {withStyles} from '@material-ui/core/styles';

import {BITRATE_MULTIPLIER, DATA_PLAN_UNLIMITED_RATES} from './DataPlanConst';

const styles = (theme: Theme) => ({
  rowIcon: {
    display: 'inline-block',
    ...theme.mixins.toolbar,
  },
});

type Props = WithStyles<typeof styles> & ContextRouter & WithAlert & {};

function Apn(props: Props) {
  const {history, relativeUrl} = useRouter();
  const enqueueSnackbar = useEnqueueSnackbar();
  const [lastRefreshTime, setLastRefreshTime] = useState(new Date().getTime());
  const {classes} = props;

  const {isLoading: apnsLoading, response: networkAPNs} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdApns,
    {networkId: nullthrows(props.match.params.networkId)},
    undefined,
    lastRefreshTime,
  );

  if (apnsLoading) {
    return <LoadingFiller />;
  }

  const handleApnEditCancel = () => {
    props.history.push(`${props.match.url}/`);
  };

  const handleApnEditSave = () => {
    history.push(relativeUrl(''));
    setLastRefreshTime(new Date().getTime());
  };

  const handleApnDelete = async (apnName: string) => {
    const confirmed = await props.confirm(
      `Are you sure you want to delete ${apnName}?`,
    );
    if (confirmed) {
      MagmaV1API.deleteLteByNetworkIdApnsByApnName({
        networkId: nullthrows(props.match.params.networkId),
        apnName: apnName,
      })
        .then(handleApnEditSave)
        .catch(error =>
          enqueueSnackbar(error.response.data.message, {variant: 'error'}),
        );
    }
  };

  const rows = Object.keys(networkAPNs || {}).map(id => {
    const profile = nullthrows(networkAPNs)[id];
    const apn_config = profile.apn_configuration;
    const qos = apn_config.qos_profile;
    return (
      <TableRow key={id}>
        <TableCell align="center">{id}</TableCell>
        <TableCell align="center">
          {apn_config.ambr.max_bandwidth_dl ===
          DATA_PLAN_UNLIMITED_RATES.max_dl_bit_rate
            ? 'Unlimited'
            : apn_config.ambr.max_bandwidth_dl / BITRATE_MULTIPLIER + ' Mbps'}
        </TableCell>
        <TableCell align="center">
          {apn_config.ambr.max_bandwidth_ul ===
          DATA_PLAN_UNLIMITED_RATES.max_ul_bit_rate
            ? 'Unlimited'
            : apn_config.ambr.max_bandwidth_ul / BITRATE_MULTIPLIER + ' Mbps'}
        </TableCell>
        <TableCell align="center">{qos.class_id}</TableCell>
        <TableCell align="center">{qos.priority_level}</TableCell>
        <TableCell align="center">
          {qos.preemption_capability === null ||
          qos.preemption_capability === false
            ? 0
            : 1}
        </TableCell>
        <TableCell align="center">
          {qos.preemption_vulnerability === null ||
          qos.preemption_vulnerability === false
            ? 0
            : 1}
        </TableCell>
        <TableCell>
          <div className={classes.rowIcon}>
            <NestedRouteLink to={`/edit/${encodeURIComponent(id)}`}>
              <IconButton color="primary">
                <EditIcon />
              </IconButton>
            </NestedRouteLink>
          </div>
          <div className={classes.rowIcon}>
            {id !== '' && (
              <IconButton color="primary" onClick={() => handleApnDelete(id)}>
                <DeleteIcon />
              </IconButton>
            )}
          </div>
        </TableCell>
      </TableRow>
    );
  });

  return (
    <>
      <Table>
        <TableHead>
          <TableRow>
            <TableCell align="center">Access Point Name (APN)</TableCell>
            <TableCell align="center">AMBR Downlink</TableCell>
            <TableCell align="center">AMBR Uplink</TableCell>
            <TableCell align="center">QoS Class ID</TableCell>
            <TableCell align="center">Priority Level</TableCell>
            <TableCell align="center">Preemption Capability</TableCell>
            <TableCell align="center">Preemption Vulnerability</TableCell>
            <TableCell>
              <NestedRouteLink to="/add">
                <Button>Add APN</Button>
              </NestedRouteLink>
            </TableCell>
          </TableRow>
        </TableHead>
        {rows && <TableBody>{rows}</TableBody>}
      </Table>
      <Route
        path={`${props.match.path}/add`}
        component={() => (
          <ApnEditDialog
            apnName={null}
            apnConfig={networkAPNs?.['']}
            onCancel={handleApnEditCancel}
            onSave={handleApnEditSave}
          />
        )}
      />
      <Route
        path={`${props.match.path}/edit/:apnName`}
        render={props => (
          <ApnEditDialog
            apnName={props.match.params.apnName}
            apnConfig={networkAPNs?.[props.match.params.apnName || '']}
            onCancel={handleApnEditCancel}
            onSave={handleApnEditSave}
          />
        )}
      />
    </>
  );
}

export default withStyles(styles)(withRouter(withAlert(Apn)));
