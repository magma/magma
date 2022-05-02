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

import type {Theme} from '@material-ui/core';
import type {WithAlert} from '../../../fbc_js_core/ui/components/Alert/withAlert';
import type {WithStyles} from '@material-ui/core';

import ApnEditDialog from './ApnEditDialog';
import Button from '@material-ui/core/Button';
import DeleteIcon from '@material-ui/icons/Delete';
import EditIcon from '@material-ui/icons/Edit';
import IconButton from '@material-ui/core/IconButton';
import LoadingFiller from '../../../fbc_js_core/ui/components/LoadingFiller';
import MagmaV1API from '../../../generated/WebClient';
import NestedRouteLink from '../../../fbc_js_core/ui/components/NestedRouteLink';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import nullthrows from '../../../fbc_js_core/util/nullthrows';
import useMagmaAPI from '../../../api/useMagmaAPI';
import withAlert from '../../../fbc_js_core/ui/components/Alert/withAlert';
import {Route, Routes, useNavigate, useParams} from 'react-router-dom';
import {useEnqueueSnackbar} from '../../../fbc_js_core/ui/hooks/useSnackbar';
import {useState} from 'react';
import {withStyles} from '@material-ui/core/styles';

import {BITRATE_MULTIPLIER, DATA_PLAN_UNLIMITED_RATES} from './DataPlanConst';

const styles = (theme: Theme) => ({
  rowIcon: {
    display: 'inline-block',
    ...theme.mixins.toolbar,
  },
});

type Props = WithStyles<typeof styles> & WithAlert & {};

function Apn(props: Props) {
  const params = useParams();
  const navigate = useNavigate();
  const enqueueSnackbar = useEnqueueSnackbar();
  const [lastRefreshTime, setLastRefreshTime] = useState(new Date().getTime());
  const {classes} = props;

  const {isLoading: apnsLoading, response: networkAPNs} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdApns,
    {networkId: nullthrows(params.networkId)},
    undefined,
    lastRefreshTime,
  );

  if (apnsLoading) {
    return <LoadingFiller />;
  }

  const handleApnEditCancel = () => {
    navigate('');
  };

  const handleApnEditSave = () => {
    navigate('');
    setLastRefreshTime(new Date().getTime());
  };

  const handleApnDelete = async (apnName: string) => {
    const confirmed = await props.confirm(
      `Are you sure you want to delete ${apnName}?`,
    );
    if (confirmed) {
      MagmaV1API.deleteLteByNetworkIdApnsByApnName({
        networkId: nullthrows(params.networkId),
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
            <NestedRouteLink to={`edit/${encodeURIComponent(id)}`}>
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
              <NestedRouteLink to="add">
                <Button variant="contained" color="primary">
                  Add APN
                </Button>
              </NestedRouteLink>
            </TableCell>
          </TableRow>
        </TableHead>
        {rows && <TableBody>{rows}</TableBody>}
      </Table>
      <Routes>
        <Route
          path="/add"
          element={
            <ApnEditDialog
              apnName={null}
              apnConfig={networkAPNs?.['']}
              onCancel={handleApnEditCancel}
              onSave={handleApnEditSave}
            />
          }
        />
        <Route
          path="/edit/:apnName"
          element={
            <ApnEditDialog
              apnName={params.apnName}
              apnConfig={networkAPNs?.[params.apnName || '']}
              onCancel={handleApnEditCancel}
              onSave={handleApnEditSave}
            />
          }
        />
      </Routes>
    </>
  );
}

export default withStyles(styles)(withAlert(Apn));
