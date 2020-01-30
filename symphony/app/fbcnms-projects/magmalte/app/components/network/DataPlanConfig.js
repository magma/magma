/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ContextRouter, Match} from 'react-router-dom';
import type {Theme} from '@material-ui/core';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithStyles} from '@material-ui/core';
import type {network_epc_configs} from '@fbcnms/magma-api';

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
import {Route, withRouter} from 'react-router-dom';

import nullthrows from '@fbcnms/util/nullthrows';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {get} from 'lodash';
import {withStyles} from '@material-ui/core/styles';

import {
  BITRATE_MULTIPLIER,
  DATA_PLAN_UNLIMITED_RATES,
  DEFAULT_DATA_PLAN_ID,
} from './DataPlanConst';

const styles = (theme: Theme) => ({
  rowIcon: {
    display: 'inline-block',
    ...theme.mixins.toolbar,
  },
});

type ErrorResponse = {
  response: {
    data: {
      message: string,
    },
  },
};

type State = {
  config: ?network_epc_configs,
  loading: boolean,
};

type Props = WithStyles<typeof styles> & ContextRouter & WithAlert & {};

class DataPlanConfig extends React.Component<Props, State> {
  state = {
    config: null,
    loading: true,
  };

  componentDidMount() {
    MagmaV1API.getLteByNetworkIdCellularEpc({
      networkId: nullthrows(this.props.match.params.networkId),
    })
      .then(response => this.setState({config: response, loading: false}))
      .catch((error: ErrorResponse) => {
        this.props.alert(get(error, 'response.data.message', error));
        this.setState({
          loading: false,
        });
      });
  }

  render() {
    const {classes, match} = this.props;
    const {config} = this.state;

    if (!config) {
      return <LoadingFiller />;
    }

    const rows = Object.keys(config.sub_profiles || {}).map(id => {
      const profile = nullthrows(config.sub_profiles)[id];
      return (
        <TableRow key={id}>
          <TableCell>{id}</TableCell>
          <TableCell>
            {profile.max_dl_bit_rate ===
            DATA_PLAN_UNLIMITED_RATES.max_dl_bit_rate
              ? 'Unlimited'
              : profile.max_dl_bit_rate / BITRATE_MULTIPLIER + ' Mbps'}
          </TableCell>
          <TableCell>
            {profile.max_ul_bit_rate ===
            DATA_PLAN_UNLIMITED_RATES.max_ul_bit_rate
              ? 'Unlimited'
              : profile.max_ul_bit_rate / BITRATE_MULTIPLIER + ' Mbps'}
          </TableCell>
          <TableCell>
            <div className={classes.rowIcon}>
              <NestedRouteLink to={`/edit/${encodeURIComponent(id)}/`}>
                <IconButton color="primary">
                  <EditIcon />
                </IconButton>
              </NestedRouteLink>
            </div>
            <div className={classes.rowIcon}>
              {id !== DEFAULT_DATA_PLAN_ID && (
                <IconButton
                  color="primary"
                  onClick={() => this.handleDataPlanDelete(id)}>
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
              <TableCell>Data Plan Name</TableCell>
              <TableCell>Download Speed</TableCell>
              <TableCell>Upload Speed</TableCell>
              <TableCell>
                <NestedRouteLink to="/edit/">
                  <Button>Add Data Plan</Button>
                </NestedRouteLink>
              </TableCell>
            </TableRow>
          </TableHead>
          {rows && <TableBody>{rows}</TableBody>}
        </Table>
        <Route path={`${match.path}/edit`} component={this.renderEditDialog} />
        <Route
          path={`${match.path}/edit/:dataPlanId`}
          component={this.renderEditDialog}
        />
      </>
    );
  }

  renderEditDialog = (props: {match: Match}) => {
    const dataPlanId = props.match.params.dataPlanId;
    return (
      <DataPlanEditDialog
        dataPlanId={dataPlanId}
        epcConfig={nullthrows(this.state.config)}
        onCancel={this.handleDataPlanEditCancel}
        onSave={this.handleDataPlanEditSave}
      />
    );
  };

  handleDataPlanDelete = (dataPlanId: string) => {
    const {config} = this.state;
    if (!config) {
      return;
    }

    this.props
      .confirm(`Are you sure you want to delete "${dataPlanId}"?`)
      .then(confirmed => {
        if (!confirmed) {
          return;
        }
        // Creates a new object without the deleted subprofiles
        const {
          [dataPlanId]: deletedProfile, // eslint-disable-line no-unused-vars
          ...newSubProfiles
        } = nullthrows(config.sub_profiles);
        const newConfig = {
          ...config,
          sub_profiles: newSubProfiles,
        };
        return MagmaV1API.putLteByNetworkIdCellularEpc({
          networkId: nullthrows(this.props.match.params.networkId),
          config: newConfig,
        }).then(() => this.setState({config: newConfig}));
      });
  };

  handleDataPlanEditCancel = () => {
    this.props.history.push(`${this.props.match.url}/`);
  };

  handleDataPlanEditSave = (
    dataPlanId: string,
    newNetworkConfig: network_epc_configs,
  ) => {
    this.setState({
      config: newNetworkConfig,
    });
    this.props.history.push(`${this.props.match.url}/`);
  };
}

export default withStyles(styles)(withRouter(withAlert(DataPlanConfig)));
