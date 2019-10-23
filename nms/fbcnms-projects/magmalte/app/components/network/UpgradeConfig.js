/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ContextRouter} from 'react-router-dom';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithStyles} from '@material-ui/core';
import type {magmad_gateway, tier} from '@fbcnms/magma-api';

import Button from '@material-ui/core/Button';
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
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import UpgradeStatusTierID from './UpgradeStatusTierID';
import UpgradeTierEditDialog from './UpgradeTierEditDialog';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';

import nullthrows from '@fbcnms/util/nullthrows';
import {Route, withRouter} from 'react-router-dom';
import {map, sortBy} from 'lodash';
import {withStyles} from '@material-ui/core/styles';

type State = {
  gateways: ?{[string]: magmad_gateway},
  errorMessage: ?string,
  saving: boolean,
  networkUpgradeTiers: ?(tier[]),
  supportedVersions: ?(string[]),
};

type Props = WithAlert & WithStyles<typeof styles> & ContextRouter & {};

const styles = _theme => ({
  header: {
    flexGrow: 1,
  },
});

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
  onUpgradeTierChange: (gatewayID: string, tierID: string) => void,
}) => {
  const {networkUpgradeTiers, onUpgradeTierChange, tableData} = props;
  const sortedTableData = sortBy(tableData, row => row.name.toLowerCase());

  const getGatewayVersionString = (state): string => {
    return (state.status && state.status.version) || 'Not Reported';
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

class UpgradeConfig extends React.Component<Props, State> {
  state = {
    gateways: null,
    errorMessage: null,
    saving: false,
    networkUpgradeTiers: null,
    supportedVersions: null,
  };

  componentDidMount() {
    this.loadData();
  }

  loadData() {
    const networkId = nullthrows(this.props.match.params.networkId);
    Promise.all([
      MagmaV1API.getChannelsByChannelId({channelId: 'stable'}),
      MagmaV1API.getNetworksByNetworkIdGateways({networkId}),
      fetchAllNetworkUpgradeTiers(networkId || ''),
    ])
      .then(([channelResp, gateways, networkUpgradeTiers]) => {
        this.setState({
          gateways,
          networkUpgradeTiers,
          supportedVersions: channelResp.supported_versions,
        });
      })
      .catch(this.props.alert);
  }

  render() {
    const {classes, match} = this.props;
    const {
      gateways,
      networkUpgradeTiers,
      supportedVersions,
      saving,
    } = this.state;

    if (!gateways) {
      return <LoadingFiller />;
    }

    return (
      <>
        {saving && <LoadingFillerBackdrop />}
        {gateways && (
          <>
            <Toolbar>
              <Typography className={classes.header} variant="h5">
                Gateway Upgrade Status
              </Typography>
            </Toolbar>
            <GatewayUpgradeStatusTable
              networkUpgradeTiers={networkUpgradeTiers}
              onUpgradeTierChange={this.handleGatewayUpgradeTierChange}
              tableData={gateways}
            />
          </>
        )}
        {supportedVersions && (
          <>
            <Toolbar>
              <Typography className={classes.header} variant="h5">
                Current Supported Versions
              </Typography>
            </Toolbar>
            <SupportedVersionsTable supportedVersions={supportedVersions} />
          </>
        )}
        {networkUpgradeTiers && (
          <>
            <Toolbar>
              <Typography className={classes.header} variant="h5">
                Upgrade Tiers
              </Typography>
              <div>
                <NestedRouteLink to={`/tier/edit/`}>
                  <Button color="primary" variant="contained">
                    Add Tier
                  </Button>
                </NestedRouteLink>
              </div>
            </Toolbar>
            <UpgradeTiersTable
              tableData={networkUpgradeTiers}
              onTierDeleteClick={this.handleUpgradeTierDelete}
            />
          </>
        )}
        <Route
          exact
          path={`${match.path}/tier/edit`}
          component={this.renderTierDialog}
        />
        <Route
          exact
          path={`${match.path}/tier/edit/:tierId`}
          component={this.renderTierDialog}
        />
      </>
    );
  }

  handleGatewayUpgradeTierChange = (gatewayID, newTierID) => {
    this.setState({saving: true});
    this.handleGatewayUpgradeTierChangeAsync(gatewayID, newTierID).catch(
      error => {
        this.props.alert(error.response?.data?.message || error);
        this.setState({saving: false});
      },
    );
  };

  async handleGatewayUpgradeTierChangeAsync(gatewayID, newTierID) {
    const networkId = nullthrows(this.props.match.params.networkId);
    await MagmaV1API.putLteByNetworkIdGatewaysByGatewayIdTier({
      networkId,
      gatewayId: gatewayID,
      tierId: JSON.stringify(`"${newTierID}"`),
    });
    const gateways = await MagmaV1API.getNetworksByNetworkIdGateways({
      networkId,
    });
    this.setState({gateways, saving: false});
  }

  handleUpgradeTierDelete = (tierId: string) => {
    this.props
      .confirm(`Are you sure you want to delete the tier ${tierId}?`)
      .then(confirmed => {
        if (!confirmed) {
          return;
        }
        MagmaV1API.deleteNetworksByNetworkIdTiersByTierId({
          networkId: nullthrows(this.props.match.params.networkId),
          tierId,
        })
          .then(_resp => this.loadData())
          .catch(this.props.alert);
      });
  };

  handleUpgradeTierEditCancel = () => {
    this.props.history.push(`${this.props.match.url}/`);
  };

  handleUpgradeTierEditSave = _config => {
    this.props.history.push(`${this.props.match.url}/`);
    this.loadData();
  };

  renderTierDialog = ({match}) => {
    const tierId = match.params.tierId;
    return (
      <UpgradeTierEditDialog
        tierId={tierId}
        onCancel={this.handleUpgradeTierEditCancel}
        onSave={this.handleUpgradeTierEditSave}
      />
    );
  };
}

export default withStyles(styles)(withRouter(withAlert(UpgradeConfig)));
