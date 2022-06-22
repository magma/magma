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
 */

import type {MagmadGateway, Tier} from '../../../generated-ts';
import type {WithAlert} from '../Alert/withAlert';

import Button from '@material-ui/core/Button';
import DeleteIcon from '@material-ui/icons/Delete';
import EditIcon from '@material-ui/icons/Edit';
import IconButton from '@material-ui/core/IconButton';
import LoadingFiller from '../LoadingFiller';
import LoadingFillerBackdrop from '../LoadingFillerBackdrop';
import MagmaAPI from '../../../api/MagmaAPI';
import NestedRouteLink from '../NestedRouteLink';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Text from '../../theme/design-system/Text';
import Toolbar from '@material-ui/core/Toolbar';
import UpgradeStatusTierID from './UpgradeStatusTierID';
import UpgradeTierEditDialog from './UpgradeTierEditDialog';
import nullthrows from '../../../shared/util/nullthrows';
import withAlert from '../Alert/withAlert';
import {GatewayId} from '../../../shared/types/network';
import {Route, Routes, useNavigate, useParams} from 'react-router-dom';
import {getErrorMessage} from '../../util/ErrorUtils';
import {makeStyles} from '@material-ui/styles';
import {map, sortBy} from 'lodash';
import {useEffect, useState} from 'react';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';

const useStyles = makeStyles(() => ({
  header: {
    flexGrow: 1,
  },
}));

const UpgradeTiersTable = (props: {
  onTierDeleteClick: (tierId: string) => void;
  tableData: Array<Tier>;
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
              <NestedRouteLink to={`tier/edit/${encodeURIComponent(row.id)}/`}>
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

const SupportedVersionsTable = (props: {supportedVersions: Array<string>}) => {
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
  tableData: Record<string, MagmadGateway>;
  networkUpgradeTiers: Array<Tier> | null | undefined;
  onUpgradeTierChange: (gatewayID: string, tierID: string) => Promise<void>;
}) => {
  const {networkUpgradeTiers, onUpgradeTierChange, tableData} = props;
  const sortedTableData = sortBy(
    Object.keys(tableData).map(k => tableData[k]),
    row => row.name.toLowerCase(),
  );

  const getGatewayVersionString = (gateway: MagmadGateway): string => {
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
            <TableCell>{row.device?.hardware_id}</TableCell>
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

function EditDialog(props: {
  networkUpgradeTiers?: Array<Tier>;
  setLastFetchTime: (time: number) => void;
}) {
  const navigate = useNavigate();
  const params = useParams();

  return (
    <UpgradeTierEditDialog
      tier={nullthrows(
        (props.networkUpgradeTiers || []).find(t => t.id === params.tierId),
      )}
      onCancel={() => navigate('..')}
      onSave={() => {
        navigate('..');
        props.setLastFetchTime(Date.now());
      }}
    />
  );
}

async function fetchAllNetworkUpgradeTiers(
  networkId: string,
): Promise<Array<Tier>> {
  const tiers = (await MagmaAPI.upgrades.networksNetworkIdTiersGet({networkId}))
    .data;
  const requests = tiers.map(tierId =>
    MagmaAPI.upgrades
      .networksNetworkIdTiersTierIdGet({networkId, tierId})
      .then(({data}) => data),
  );
  return await Promise.all(requests);
}

function UpgradeConfig(props: WithAlert) {
  const classes = useStyles();
  const navigate = useNavigate();
  const params = useParams();
  const [gateways, setGateways] = useState<Record<string, MagmadGateway>>();
  const [networkUpgradeTiers, setNetworkUpgradeTiers] = useState<Array<Tier>>();
  const [supportedVersions, setSupportedVersions] = useState<Array<string>>();
  const [saving, setSaving] = useState(false);
  const [lastFetchTime, setLastFetchTime] = useState(Date.now());
  const enqueueSnackbar = useEnqueueSnackbar();

  const networkId = nullthrows(params.networkId);
  useEffect(() => {
    async function fetchStableReleases() {
      let supportedVersions;
      try {
        supportedVersions = (
          await MagmaAPI.upgrades.channelsChannelIdGet({
            channelId: 'stable',
          })
        ).data.supported_versions;
      } catch (e) {
        enqueueSnackbar('Unable to fetch stable releases', {variant: 'error'});
      }
      setSupportedVersions(supportedVersions);
    }

    async function fetchAllData() {
      const [networkUpgradeTiers, response] = await Promise.all([
        fetchAllNetworkUpgradeTiers(networkId),
        MagmaAPI.gateways.networksNetworkIdGatewaysGet({networkId}),
        fetchStableReleases(),
      ]);

      setGateways(response.data.gateways);
      setNetworkUpgradeTiers(networkUpgradeTiers);
    }

    fetchAllData().catch(e =>
      enqueueSnackbar(getErrorMessage(e), {variant: 'error'}),
    );
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
    void props
      .confirm(`Are you sure you want to delete the tier ${tierId}?`)
      .then(confirmed => {
        if (!confirmed) {
          return;
        }
        void MagmaAPI.upgrades
          .networksNetworkIdTiersTierIdDelete({
            networkId,
            tierId,
          })
          .then(() => setLastFetchTime(Date.now()))
          .catch(e => enqueueSnackbar(getErrorMessage(e), {variant: 'error'}));
      });
  };

  const handleGatewayUpgradeTierChange = async (
    gatewayID: GatewayId,
    newTierID: string,
  ) => {
    setSaving(true);
    try {
      await MagmaAPI.lteGateways.lteNetworkIdGatewaysGatewayIdTierPut({
        networkId,
        gatewayId: gatewayID,
        tierId: JSON.stringify(`"${newTierID}"`),
      });
      const paginated_gateways = (
        await MagmaAPI.gateways.networksNetworkIdGatewaysGet({
          networkId,
        })
      ).data;
      setGateways(paginated_gateways.gateways);
    } catch (error) {
      enqueueSnackbar(getErrorMessage(error), {
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
              <NestedRouteLink to={`tier/new/`}>
                <Button variant="contained" color="primary">
                  Add Tier
                </Button>
              </NestedRouteLink>
            </div>
          </Toolbar>
          <UpgradeTiersTable
            tableData={networkUpgradeTiers}
            onTierDeleteClick={handleUpgradeTierDelete}
          />
        </>
      )}
      <Routes>
        <Route
          path="tier/new"
          element={
            <UpgradeTierEditDialog
              onCancel={() => navigate('')}
              onSave={() => {
                navigate('');
                setLastFetchTime(Date.now());
              }}
            />
          }
        />
        <Route
          path="tier/edit/:tierId"
          element={
            <EditDialog
              setLastFetchTime={setLastFetchTime}
              networkUpgradeTiers={networkUpgradeTiers}
            />
          }
        />
      </Routes>
    </>
  );
}

export default withAlert(UpgradeConfig);
