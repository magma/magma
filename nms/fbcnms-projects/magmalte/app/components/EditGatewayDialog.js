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
import type {WithStyles} from '@material-ui/core';
import type {Gateway} from './GatewayUtils';

import AppBar from '@material-ui/core/AppBar';
import Dialog from '@material-ui/core/Dialog';
import React from 'react';
import Tabs from '@material-ui/core/Tabs';
import Tab from '@material-ui/core/Tab';
import GatewayMagmadFields from './GatewayMagmadFields';
import GatewaySummaryFields from './GatewaySummaryFields';
import GatewayCellularFields from './GatewayCellularFields';
import GatewayCommandFields from './GatewayCommandFields';

import {fetchDevice} from '../common/MagmaAPI';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

const styles = {
  appBar: {
    backgroundColor: '#f5f5f5',
    marginBottom: '20px',
  },
};

type Props = ContextRouter &
  WithStyles & {
    onClose: () => void,
    onSave: (gateway: {[string]: any}) => void,
    gateway: ?Gateway,
  };

type State = {
  tab: number,
};

class EditGatewayDialog extends React.Component<Props, State> {
  state = {
    error: '',
    tab: 0,
  };

  render() {
    if (!this.props.gateway) {
      return null;
    }

    const {classes} = this.props;
    let content;
    const props = {
      onClose: this.props.onClose,
      gateway: this.props.gateway,
      onSave: this.onSave,
    };

    switch (this.state.tab) {
      case 0:
        content = <GatewaySummaryFields {...props} />;
        break;
      case 1:
        content = <GatewayCellularFields {...props} />;
        break;
      case 2:
        content = <GatewayMagmadFields {...props} />;
        break;
      case 3:
        content = (
          <GatewayCommandFields
            {...props}
            showRestartCommand={true}
            showPingCommand={true}
            showGenericCommand={true}
          />
        );
        break;
    }
    return (
      <Dialog
        open={true}
        onClose={this.props.onClose}
        maxWidth="md"
        scroll="body">
        <AppBar position="static" className={classes.appBar}>
          <Tabs
            indicatorColor="primary"
            textColor="primary"
            value={this.state.tab}
            onChange={this.onTabChange}>
            <Tab label="Summary" />
            <Tab label="LTE" />
            <Tab label="Magma" />
            <Tab label="Commands" />
          </Tabs>
        </AppBar>
        {content}
      </Dialog>
    );
  }

  onTabChange = (event, tab) => this.setState({tab});
  onSave = gatewayID => {
    fetchDevice(this.props.match, gatewayID).then(this.props.onSave);
  };
}

export default withStyles(styles)(withRouter(EditGatewayDialog));
