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

import withAlert from '../../../fbc_js_core/ui/components/Alert/withAlert';
import type {Theme} from '@material-ui/core';
import type {WithAlert} from '../../../fbc_js_core/ui/components/Alert/withAlert';

import Button from '@material-ui/core/Button';
import DataPlanEditDialog from './DataPlanEditDialog';
import DeleteIcon from '@material-ui/icons/Delete';
import EditIcon from '@material-ui/icons/Edit';
import IconButton from '@material-ui/core/IconButton';
import LoadingFiller from '../../../fbc_js_core/ui/components/LoadingFiller';
import MagmaV1API from '../../../generated/WebClient';
import NestedRouteLink from '../../../fbc_js_core/ui/components/NestedRouteLink';
import React, {useState} from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import {Route, Routes, useNavigate, useParams} from 'react-router-dom';

import nullthrows from '../../../fbc_js_core/util/nullthrows';
import useMagmaAPI from '../../../api/useMagmaAPI';
import {makeStyles} from '@material-ui/styles';

import {
  BITRATE_MULTIPLIER,
  DATA_PLAN_UNLIMITED_RATES,
  DEFAULT_DATA_PLAN_ID,
} from './DataPlanConst';
import type {network_epc_configs} from '../../../generated/MagmaAPIBindings';

const useStyles = makeStyles((theme: Theme) => ({
  rowIcon: {
    display: 'inline-block',
    ...theme.mixins.toolbar,
  },
}));

type Props = WithAlert & {};

function EditDialog(props: {
  config: network_epc_configs,
  onSave: (editName: string, newConfig: network_epc_configs) => void,
}) {
  const navigate = useNavigate();
  const params = useParams();

  return (
    <DataPlanEditDialog
      dataPlanId={params.dataPlanId}
      epcConfig={props.config}
      onCancel={() => navigate('..')}
      onSave={props.onSave}
    />
  );
}

function DataPlanConfig(props: Props) {
  const classes = useStyles();
  const navigate = useNavigate();
  const params = useParams();
  const [config, setConfig] = useState();

  const {isLoading} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdCellularEpc,
    {networkId: nullthrows(params.networkId)},
    setConfig,
  );

  if (!config || isLoading) {
    return <LoadingFiller />;
  }

  const onDelete = (dataPlanId: string) => {
    props
      .confirm(`Are you sure you want to delete "${dataPlanId}"?`)
      .then(confirmed => {
        if (!confirmed) {
          return;
        }
        // Creates a new object without the deleted subprofiles
        const {[dataPlanId]: _deletedProfile, ...newSubProfiles} = nullthrows(
          config.sub_profiles,
        );
        const newConfig = {
          ...config,
          sub_profiles: newSubProfiles,
        };
        return MagmaV1API.putLteByNetworkIdCellularEpc({
          networkId: nullthrows(params.networkId),
          config: newConfig,
        }).then(() => setConfig(newConfig));
      });
  };

  const rows = Object.keys(config.sub_profiles || {}).map(id => {
    const profile = nullthrows(config.sub_profiles)[id];
    return (
      <TableRow key={id}>
        <TableCell>{id}</TableCell>
        <TableCell>
          {profile.max_dl_bit_rate === DATA_PLAN_UNLIMITED_RATES.max_dl_bit_rate
            ? 'Unlimited'
            : profile.max_dl_bit_rate / BITRATE_MULTIPLIER + ' Mbps'}
        </TableCell>
        <TableCell>
          {profile.max_ul_bit_rate === DATA_PLAN_UNLIMITED_RATES.max_ul_bit_rate
            ? 'Unlimited'
            : profile.max_ul_bit_rate / BITRATE_MULTIPLIER + ' Mbps'}
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
            {id !== DEFAULT_DATA_PLAN_ID && (
              <IconButton color="primary" onClick={() => onDelete(id)}>
                <DeleteIcon />
              </IconButton>
            )}
          </div>
        </TableCell>
      </TableRow>
    );
  });

  const onSave = (_, newConfig) => {
    navigate('');
    setConfig(newConfig);
  };
  return (
    <>
      <Table>
        <TableHead>
          <TableRow>
            <TableCell>Data Plan Name</TableCell>
            <TableCell>Download Speed</TableCell>
            <TableCell>Upload Speed</TableCell>
            <TableCell>
              <NestedRouteLink to="add">
                <Button variant="contained" color="primary">
                  Add Data Plan
                </Button>
              </NestedRouteLink>
            </TableCell>
          </TableRow>
        </TableHead>
        {rows && <TableBody>{rows}</TableBody>}
      </Table>
      <Routes>
        <Route
          path="add"
          element={
            <DataPlanEditDialog
              dataPlanId={null}
              epcConfig={config}
              onCancel={() => navigate('')}
              onSave={onSave}
            />
          }
        />
        <Route
          path="edit/:dataPlanId"
          element={<EditDialog config={config} onSave={onSave} />}
        />
      </Routes>
    </>
  );
}

export default withAlert(DataPlanConfig);
