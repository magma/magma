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

import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {magmad_gateway, tier} from '@fbcnms/magma-api';

import Button from '@fbcnms/ui/components/design-system/Button';
import DeleteIcon from '@material-ui/icons/Delete';
import EditIcon from '@material-ui/icons/Edit';
import IconButton from '@material-ui/core/IconButton';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import LoadingFillerBackdrop from '@fbcnms/ui/components/LoadingFillerBackdrop';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Text from '@fbcnms/ui/components/design-system/Text';
import Toolbar from '@material-ui/core/Toolbar';
import UpgradeStatusTierID from './UpgradeStatusTierID';
import UpgradeTierEditDialog from './UpgradeTierEditDialog';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';

import nullthrows from '@fbcnms/util/nullthrows';
import {Route} from 'react-router-dom';
import {makeStyles} from '@material-ui/styles';
import {map, sortBy} from 'lodash';
import {useEffect, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(() => ({
  header: {
    flexGrow: 1,
  },
}));

const UpgradeTiersTable = (props: {
  onTierDeleteClick: (tierId: string) => void,
  tableData: Array<tier>,
}) => {
  return (
    <Table>
      <TableHead>
        <TableRow>
          <TableCell>Tier ID</TableCell>
          <TableCell>Tier Name</TableCell>
          <TableCell>Software Version</TableCell>
          <TableCell />
        </TableRow>
      </TableHead>
      <TableBody>
        {map(props.tableData, row => (
          <TableRow key={row.id}>
            <TableCell>{row.id}</TableCell>
            <TableCell>{row.name}</TableCell>
            <TableCell>{row.version}</TableCell>
            <TableCell>
              <NestedRouteLink to={`/tier/edit/${encodeURIComponent(row.id)}/`}>
                <IconButton>
                  <EditIcon />
                </IconButton>
              </NestedRouteLink>
              <IconButton onClick={() => props.onTierDeleteClick(row.id)}>
                <DeleteIcon />
              </IconButton>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
};

const SupportedVersionsTable = (props: {supportedVersions: string[]}) => {
  return (
    <Table>
      <TableBody>
        {map(props.supportedVersions, (version, i: number) => (
          <TableRow key={version}>
            <TableCell>
              {version}
              {i === props.supportedVersions.length - 1 && (
                <b> (Newest Version)</b>
              )}
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
};

const GatewayUpgradeStatusTable = (props: {
  tableData: {[string]: magmad_gateway},
  networkUpgradeTiers: ?(tier[]),
  onUpgradeTierChange: (gatewayID: string, tierID: string) => Promise<void>,
}) => {
  const {networkUpgradeTiers, onUpgradeTierChange, tableData} = props;
  const sortedTableData = sortBy(
    Object.keys(tableData).map(k => tableData[k]),
    row => row.name.toLowerCase(),
  );

  const getGatewayVersionString = (gateway): string => {
    const packages = gateway.status?.platform_info?.packages || [];
    return packages.find(p => p.name === 'magma')?.version || 'Not Reported';
  };

  return (
    <Table>
      <TableHead>
        <TableRow>
          <TableCell>Name</TableCell>
          <TableCell>Hardware UUID</TableCell>
          <TableCell>Tier ID</TableCell>
          <TableCell>Current Version</TableCell>
        </TableRow>
      </TableHead>
      <TableBody>
        {map(sortedTableData, row => (
          <TableRow key={row.id}>
            <TableCell>{row.name}</TableCell>
            <TableCell>{row.device.hardware_id}</TableCell>
            <TableCell>
              <UpgradeStatusTierID
                onChange={onUpgradeTierChange}
                gatewayID={row.id}
                tierID={row.tier}
                networkUpgradeTiers={networkUpgradeTiers}
              />
            </TableCell>
            <TableCell>{getGatewayVersionString(row)}</TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
};

async function fetchAllNetworkUpgradeTiers(
  networkId: string,
): Promise<Array<tier>> {
  const tiers = await MagmaV1API.getNetworksByNetworkIdTiers({networkId});
  const requests = tiers.map(tierId =>
    MagmaV1API.getNetworksByNetworkIdTiersByTierId({networkId, tierId}),
  );
  return await Promise.all(requests);
}

function UpgradeConfig(props: WithAlert & {}) {
  const classes = useStyles();
  const {match, history, relativePath, relativeUrl} = useRouter();
  const [gateways, setGateways] = useState();
  const [networkUpgradeTiers, setNetworkUpgradeTiers] = useState();
  const [supportedVersions, setSupportedVersions] = useState();
  const [saving, setSaving] = useState(false);
  const [lastFetchTime, setLastFetchTime] = useState(Date.now());
  const enqueueSnackbar = useEnqueueSnackbar();

  const networkId = nullthrows(match.params.networkId);
  useEffect(() => {
    async function fetchStableReleases() {
      let supportedVersions;
      try {
        supportedVersions = (
          await MagmaV1API.getChannelsByChannelId({
            channelId: 'stable',
          })
        ).supported_versions;
      } catch (e) {
        enqueueSnackbar('Unable to fetch stable releases', {variant: 'error'});
      }
      setSupportedVersions(supportedVersions);
    }

    async function fetchAllData() {
      const [networkUpgradeTiers, gateways] = await Promise.all([
        fetchAllNetworkUpgradeTiers(networkId),
        MagmaV1API.getNetworksByNetworkIdGateways({networkId}),
        fetchStableReleases(),
      ]);

      setGateways(gateways);
      setNetworkUpgradeTiers(networkUpgradeTiers);
    }

    fetchAllData().catch(e => enqueueSnackbar(e, {variant: 'error'}));
  }, [
    networkId,
    setGateways,
    setNetworkUpgradeTiers,
    setSupportedVersions,
    lastFetchTime,
    enqueueSnackbar,
  ]);

  if (!gateways) {
    return <LoadingFiller />;
  }

  const handleUpgradeTierDelete = (tierId: string) => {
    props
      .confirm(`Are you sure you want to delete the tier ${tierId}?`)
      .then(confirmed => {
        if (!confirmed) {
          return;
        }
        MagmaV1API.deleteNetworksByNetworkIdTiersByTierId({
          networkId,
          tierId,
        })
          .then(() => setLastFetchTime(Date.now()))
          .catch(e => enqueueSnackbar(e, {variant: 'error'}));
      });
  };

  const handleGatewayUpgradeTierChange = async (gatewayID, newTierID) => {
    setSaving(true);
    try {
      await MagmaV1API.putLteByNetworkIdGatewaysByGatewayIdTier({
        networkId,
        gatewayId: gatewayID,
        tierId: JSON.stringify(`"${newTierID}"`),
      });
      const gateways = await MagmaV1API.getNetworksByNetworkIdGateways({
        networkId,
      });
      setGateways(gateways);
    } catch (error) {
      enqueueSnackbar(error.response?.data?.message || error, {
        variant: 'error',
      });
    }
    setSaving(false);
  };

  return (
    <>
      {saving && <LoadingFillerBackdrop />}
      {gateways && (
        <>
          <Toolbar>
            <Text className={classes.header} variant="h5">
              Gateway Upgrade Status
            </Text>
          </Toolbar>
          <GatewayUpgradeStatusTable
            networkUpgradeTiers={networkUpgradeTiers}
            onUpgradeTierChange={handleGatewayUpgradeTierChange}
            tableData={gateways}
          />
        </>
      )}
      {supportedVersions && (
        <>
          <Toolbar>
            <Text className={classes.header} variant="h5">
              Current Supported Versions
            </Text>
          </Toolbar>
          <SupportedVersionsTable supportedVersions={supportedVersions} />
        </>
      )}
      {networkUpgradeTiers && (
        <>
          <Toolbar>
            <Text className={classes.header} variant="h5">
              Upgrade Tiers
            </Text>
            <div>
              <NestedRouteLink to={`/tier/new/`}>
                <Button>Add Tier</Button>
              </NestedRouteLink>
            </div>
          </Toolbar>
          <UpgradeTiersTable
            tableData={networkUpgradeTiers}
            onTierDeleteClick={handleUpgradeTierDelete}
          />
        </>
      )}
      <Route
        exact
        path={relativePath('/tier/new')}
        component={() => (
          <UpgradeTierEditDialog
            onCancel={() => history.push(relativeUrl(''))}
            onSave={() => {
              history.push(relativeUrl(''));
              setLastFetchTime(Date.now());
            }}
          />
        )}
      />
      <Route
        exact
        path={relativePath('/tier/edit/:tierId')}
        component={({match}) => (
          <UpgradeTierEditDialog
            tier={nullthrows(
              (networkUpgradeTiers || []).find(
                t => t.id === match.params.tierId,
              ),
            )}
            onCancel={() => history.push(relativeUrl(''))}
            onSave={() => {
              history.push(relativeUrl(''));
              setLastFetchTime(Date.now());
            }}
          />
        )}
      />
    </>
  );
}

export default withAlert(UpgradeConfig);
