/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {ContextRouter} from 'react-router-dom';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithStyles} from '@material-ui/core';
import type {enodeb} from '@fbcnms/magma-api';

import AddEditEnodebDialog from './AddEditEnodebDialog';
import Button from '@fbcnms/ui/components/design-system/Button';
import DeleteIcon from '@material-ui/icons/Delete';
import EditIcon from '@material-ui/icons/Edit';
import IconButton from '@material-ui/core/IconButton';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Text from '@fbcnms/ui/components/design-system/Text';

import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import nullthrows from '@fbcnms/util/nullthrows';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

const styles = theme => ({
  header: {
    margin: '10px',
    display: 'flex',
    justifyContent: 'space-between',
  },
  paper: {
    margin: theme.spacing(3),
  },
});

type Props = ContextRouter & WithAlert & WithStyles<typeof styles> & {};

type State = {
  showDialog: boolean,
  enodebs: enodeb[],
  editingEnodeb: ?enodeb,
};

class Enodebs extends React.Component<Props, State> {
  state = {
    showDialog: false,
    enodebs: [],
    editingEnodeb: null,
  };

  componentDidMount() {
    const {match} = this.props;
    MagmaV1API.getLteByNetworkIdEnodebs({
      networkId: nullthrows(match.params.networkId),
    })
      .then(result => {
        const enodebs = Object.keys(result).map(key => result[key]);
        this.setState({enodebs});
      })
      .catch(error => {
        this.props.alert('Failed to get eNB for network: ' + error);
      });
  }

  render() {
    const {enodebs} = this.state;
    const rows = (enodebs || []).map(enodeb => (
      <TableRow key={enodeb.serial}>
        <TableCell>
          {status}
          {enodeb.serial}
        </TableCell>
        <TableCell>{enodeb.config.device_class}</TableCell>
        <TableCell>
          <IconButton onClick={() => this.editEnodeb(enodeb)}>
            <EditIcon />
          </IconButton>
          <IconButton onClick={() => this.deleteEnodeb(enodeb)}>
            <DeleteIcon />
          </IconButton>
        </TableCell>
      </TableRow>
    ));

    const addEditEnodebDialog = this.state.showDialog ? (
      <AddEditEnodebDialog
        editingEnodeb={this.state.editingEnodeb}
        onClose={this.hideDialog}
        onSave={this.onSave}
      />
    ) : null;

    return (
      <>
        <div className={this.props.classes.paper}>
          <div className={this.props.classes.header}>
            <Text variant="h5">Configure eNodeB Devices</Text>
            <Button onClick={this.showDialog}>Add eNodeB</Button>
          </div>
          <Paper elevation={2}>
            {enodebs ? (
              <Table>
                <TableHead>
                  <TableRow>
                    <TableCell>Serial ID</TableCell>
                    <TableCell>Device Class</TableCell>
                    <TableCell />
                  </TableRow>
                </TableHead>
                <TableBody>{rows}</TableBody>
              </Table>
            ) : (
              <LoadingFiller />
            )}
          </Paper>
          {addEditEnodebDialog}
        </div>
      </>
    );
  }

  showDialog = () => this.setState({showDialog: true});
  hideDialog = () => {
    this.setState({
      showDialog: false,
      editingEnodeb: null,
    });
  };
  editEnodeb = editingEnodeb => {
    this.showDialog();
    this.setState({editingEnodeb});
  };

  onSave = (enodeb: enodeb) => {
    const newEnodebs = nullthrows(this.state.enodebs).slice(0);
    if (this.state.editingEnodeb) {
      newEnodebs[newEnodebs.indexOf(this.state.editingEnodeb)] = enodeb;
    } else {
      newEnodebs.push(enodeb);
    }
    this.setState({
      enodebs: newEnodebs,
      showDialog: false,
      editingEnodeb: null,
    });
  };

  deleteEnodeb = enodeb => {
    const {match} = this.props;
    this.props
      .confirm(`Are you sure you want to delete ${enodeb.serial}?`)
      .then(confirmed => {
        if (!confirmed) {
          return;
        }
        MagmaV1API.deleteLteByNetworkIdEnodebsByEnodebSerial({
          networkId: nullthrows(match.params.networkId),
          enodebSerial: enodeb.serial,
        }).then(() =>
          this.setState({
            enodebs: this.state.enodebs.filter(
              enb => enb.serial != enodeb.serial,
            ),
          }),
        );
      });
  };
}

export default withStyles(styles)(withAlert(withRouter(Enodebs)));
