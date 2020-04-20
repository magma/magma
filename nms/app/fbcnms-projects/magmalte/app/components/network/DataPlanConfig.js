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
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';

import Button from '@fbcnms/ui/components/design-system/Button';
import DataPlanEditDialog from './DataPlanEditDialog';
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
import {Route} from 'react-router-dom';

import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';
import {useState} from 'react';

import {
  BITRATE_MULTIPLIER,
  DATA_PLAN_UNLIMITED_RATES,
  DEFAULT_DATA_PLAN_ID,
} from './DataPlanConst';

const useStyles = makeStyles((theme: Theme) => ({
  rowIcon: {
    display: 'inline-block',
    ...theme.mixins.toolbar,
  },
}));

type Props = WithAlert & {};

function DataPlanConfig(props: Props) {
  const classes = useStyles();
  const {match, history, relativePath, relativeUrl} = useRouter();
  const [config, setConfig] = useState();

  const {isLoading} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdCellularEpc,
    {networkId: nullthrows(match.params.networkId)},
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
          networkId: nullthrows(match.params.networkId),
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
            <NestedRouteLink to={`/edit/${encodeURIComponent(id)}`}>
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
    history.push(relativeUrl(''));
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
              <NestedRouteLink to="/add">
                <Button>Add Data Plan</Button>
              </NestedRouteLink>
            </TableCell>
          </TableRow>
        </TableHead>
        {rows && <TableBody>{rows}</TableBody>}
      </Table>
      <Route
        path={relativePath('/add')}
        component={() => (
          <DataPlanEditDialog
            dataPlanId={null}
            epcConfig={config}
            onCancel={() => history.push(relativeUrl(''))}
            onSave={onSave}
          />
        )}
      />
      <Route
        path={relativePath('/edit/:dataPlanId')}
        component={({match}) => (
          <DataPlanEditDialog
            dataPlanId={match.params.dataPlanId}
            epcConfig={config}
            onCancel={() => history.push(relativeUrl(''))}
            onSave={onSave}
          />
        )}
      />
    </>
  );
}

export default withAlert(DataPlanConfig);
