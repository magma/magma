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
  CellularNetworkConfig,
  CellularNetworkProfile,
} from '../../common/MagmaAPIType';
import type {ContextRouter, Match} from 'react-router-dom';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithStyles} from '@material-ui/core';
import type {Theme} from '@material-ui/core';

import React from 'react';
import axios from 'axios';
import Button from '@material-ui/core/Button';
import DataPlanEditDialog from './DataPlanEditDialog';
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
import {Route, withRouter} from 'react-router-dom';

import {MagmaAPIUrls} from '../../common/MagmaAPI';
import {get, map} from 'lodash-es';
import {withStyles} from '@material-ui/core/styles';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';

import {
  DEFAULT_DATA_PLAN_ID,
  BITRATE_MULTIPLIER,
  DATA_PLAN_UNLIMITED_RATES,
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
  config: ?CellularNetworkConfig,
  loading: boolean,
};

type Props = WithStyles & ContextRouter & WithAlert & {};

class DataPlanConfig extends React.Component<Props, State> {
  state = {
    config: null,
    loading: true,
  };

  componentDidMount() {
    axios
      .get<null, CellularNetworkConfig>(
        MagmaAPIUrls.networkConfigsForType(this.props.match, 'cellular'),
      )
      .then(response =>
        this.setState({
          config: response.data,
          loading: false,
        }),
      )
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

    const rows = map(
      config.epc.sub_profiles,
      (profile: CellularNetworkProfile, id) => (
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
          {/* TODO: Subscriber Count */}
          <TableCell>0</TableCell>
          <TableCell>
            <div className={classes.rowIcon}>
              <NestedRouteLink to={`/edit/${encodeURIComponent(id)}/`}>
                <IconButton>
                  <EditIcon />
                </IconButton>
              </NestedRouteLink>
            </div>
            <div className={classes.rowIcon}>
              {id !== DEFAULT_DATA_PLAN_ID && (
                <IconButton onClick={() => this.handleDataPlanDelete(id)}>
                  <DeleteIcon />
                </IconButton>
              )}
            </div>
          </TableCell>
        </TableRow>
      ),
    );

    return (
      <>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Data Plan Name</TableCell>
              <TableCell>Download Speed</TableCell>
              <TableCell>Upload Speed</TableCell>
              <TableCell>Subscriber Count</TableCell>
              <TableCell>
                <NestedRouteLink to="/edit/">
                  <Button color="primary" variant="contained">
                    Add Data Plan
                  </Button>
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
        networkConfig={this.state.config}
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
        } = config.epc.sub_profiles;
        const newConfig = {
          ...config,
          epc: {
            ...config.epc,
            sub_profiles: newSubProfiles,
          },
        };
        return axios
          .put(
            MagmaAPIUrls.networkConfigsForType(this.props.match, 'cellular'),
            newConfig,
          )
          .then(_resp => this.setState({config: newConfig}));
      });
  };

  handleDataPlanEditCancel = () => {
    this.props.history.push(`${this.props.match.url}/`);
  };

  handleDataPlanEditSave = (
    dataPlanId: string,
    newNetworkConfig: CellularNetworkConfig,
  ) => {
    this.setState({
      config: newNetworkConfig,
    });
    this.props.history.push(`${this.props.match.url}/`);
  };
}

export default withStyles(styles)(withRouter(withAlert(DataPlanConfig)));
