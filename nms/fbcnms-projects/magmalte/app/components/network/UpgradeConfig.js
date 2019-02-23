/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {
  CheckindGateway,
  NetworkUpgradeTier,
} from '../../common/MagmaAPIType';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithStyles} from '@material-ui/core';
import type {ContextRouter} from 'react-router-dom';

import axios from 'axios';
import React from 'react';
import Button from '@material-ui/core/Button';
import DeleteIcon from '@material-ui/icons/Delete';
import EditIcon from '@material-ui/icons/Edit';
import IconButton from '@material-ui/core/IconButton';
import LoadingFiller from '../LoadingFiller';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
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

import {Route, withRouter} from 'react-router-dom';

import {
  fetchAllNetworkUpgradeTiers,
  fetchAllGateways,
  MagmaAPIUrls,
} from '../../common/MagmaAPI';
import {get, map, merge, sortBy} from 'lodash-es';
import {withStyles} from '@material-ui/core/styles';

type State = {
  gateways: ?(CheckindGateway[]),
  errorMessage: ?string,
  loading: boolean,
  networkUpgradeTiers: ?(NetworkUpgradeTier[]),
  supportedVersions: ?(string[]),
};

type Props = WithAlert & WithStyles & ContextRouter & {};

const styles = _theme => ({
  header: {
    flexGrow: 1,
  },
});

const UpgradeTiersTable = (props: {
  onTierDeleteClick: (tierId: string) => void,
  tableData: Array<NetworkUpgradeTier>,
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
  tableData: Array<CheckindGateway>,
  networkUpgradeTiers: ?(NetworkUpgradeTier[]),
  onUpgradeTierChange: (gatewayID: string, tierID: string) => void,
}) => {
  const {networkUpgradeTiers, onUpgradeTierChange, tableData} = props;
  const sortedTableData = sortBy(tableData, row =>
    row.record.name.toLowerCase(),
  );

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
          <TableRow key={row.gateway_id}>
            <TableCell>{row.record.name}</TableCell>
            <TableCell>{get(row, 'record.hw_id.id')}</TableCell>
            <TableCell>
              <UpgradeStatusTierID
                onChange={onUpgradeTierChange}
                gatewayID={row.gateway_id}
                tierID={get(row, 'config.magmad_gateway.tier')}
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

class UpgradeConfig extends React.Component<Props, State> {
  state = {
    gateways: null,
    errorMessage: null,
    loading: true,
    networkUpgradeTiers: null,
    supportedVersions: null,
  };

  componentDidMount() {
    this.loadData();
  }

  loadData() {
    const {networkId} = this.props.match.params;
    Promise.all([
      axios.get(MagmaAPIUrls.upgradeChannel('stable')),
      fetchAllGateways(networkId || ''),
      fetchAllNetworkUpgradeTiers(networkId || ''),
    ])
      .then(([channelResp, gateways, networkUpgradeTiers]) => {
        this.setState({
          gateways,
          networkUpgradeTiers,
          supportedVersions: channelResp.data.supported_versions,
        });
      })
      .catch(this.props.alert);
  }

  render() {
    const {classes, match} = this.props;
    const {gateways, networkUpgradeTiers, supportedVersions} = this.state;

    if (!gateways) {
      return <LoadingFiller />;
    }

    return (
      <>
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
    this.handleGatewayUpgradeTierChangeAsync(gatewayID, newTierID).catch(
      error => this.props.alert(error.response?.data?.message || error),
    );
  };

  async handleGatewayUpgradeTierChangeAsync(gatewayID, newTierID) {
    const networkId = this.props.match.params.networkId || '';
    const url = MagmaAPIUrls.gatewayConfigs(this.props.match, gatewayID);
    const resp = await axios.get(url);
    const newData = merge({}, resp.data, {
      tier: newTierID,
    });
    await axios.put(url, newData);
    const gateways = await fetchAllGateways(networkId);
    this.setState({gateways});
  }

  handleUpgradeTierDelete = (tierId: string) => {
    this.props
      .confirm(`Are you sure you want to delete the tier ${tierId}?`)
      .then(confirmed => {
        if (!confirmed) {
          return;
        }
        axios
          .delete(MagmaAPIUrls.networkTier(this.props.match, tierId))
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
