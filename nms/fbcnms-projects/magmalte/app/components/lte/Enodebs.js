/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {ContextRouter} from 'react-router-dom';
import type {WithStyles} from '@material-ui/core';
import type {Enodeb, EnodebPayload} from './EnodebUtils';

import axios from 'axios';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Paper from '@material-ui/core/Paper';
import Typography from '@material-ui/core/Typography';
import Button from '@material-ui/core/Button';
import AddEditEnodebDialog from './AddEditEnodebDialog';
import DeleteIcon from '@material-ui/icons/Delete';
import EditIcon from '@material-ui/icons/Edit';
import IconButton from '@material-ui/core/IconButton';
import MagmaTopBar from '../MagmaTopBar';

import {DEFAULT_ENODEB} from './EnodebUtils';
import {MagmaAPIUrls} from '../../common/MagmaAPI';
import nullthrows from '@fbcnms/util/nullthrows';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';

const styles = theme => ({
  header: {
    margin: '10px',
    display: 'flex',
    justifyContent: 'space-between',
  },
  paper: {
    margin: theme.spacing.unit * 3,
  },
});

type Props = ContextRouter & WithAlert & WithStyles & {};

type State = {
  showDialog: boolean,
  enodebs: Enodeb[],
  editingEnodeb: ?any,
};

class Enodebs extends React.Component<Props, State> {
  state = {
    showDialog: false,
    enodebs: [],
    editingEnodeb: null,
  };

  componentDidMount() {
    const {match} = this.props;
    const handleFunc = this._handleEnodebPayloads.bind(this);
    axios
      .get(MagmaAPIUrls.enodebs(match))
      .then(response => {
        const enbSerialArr = response.data;
        const enbReqArr = enbSerialArr.map(enbSerial =>
          axios.get(MagmaAPIUrls.enodeb(match, enbSerial)),
        );
        axios.all(enbReqArr).then(payloadArr => {
          handleFunc(enbSerialArr, payloadArr);
        });
      })
      .catch(error => {
        this.props.alert('Failed to get eNB for network: ' + error);
      });
  }

  _handleEnodebPayloads(enbSerialArr, payloadArr) {
    const handleFunc = this._buildEnodebFromPayload.bind(this);
    const enodebs = [];
    payloadArr.forEach(function(payload, i) {
      enodebs.push(handleFunc(enbSerialArr[i], payload.data));
    });
    this.setState({enodebs});
  }

  render() {
    const {enodebs} = this.state;
    const rows = (enodebs || []).map(enodeb => (
      <TableRow key={enodeb.serialId}>
        <TableCell>
          {status}
          {enodeb.serialId}
        </TableCell>
        <TableCell>{enodeb.deviceClass}</TableCell>
        <TableCell>
          <IconButton onClick={this.editEnodeb.bind(this, enodeb)}>
            <EditIcon />
          </IconButton>
          <IconButton onClick={this.deleteEnodeb.bind(this, enodeb)}>
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
        <MagmaTopBar title="eNodeB Devices" />
        <div className={this.props.classes.paper}>
          <div className={this.props.classes.header}>
            <Typography variant="h5">Configure eNodeB Devices</Typography>
            <Button
              variant="contained"
              color="primary"
              onClick={this.showDialog}>
              Add eNodeB
            </Button>
          </div>
          <Paper>
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

  onSave = (enbSerial: string, enodebPayload: EnodebPayload) => {
    const enodeb = this._buildEnodebFromPayload(enbSerial, enodebPayload);
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
      .confirm(`Are you sure you want to delete ${enodeb.serialId}?`)
      .then(confirmed => {
        if (!confirmed) {
          return;
        }
        axios.delete(MagmaAPIUrls.enodeb(match, enodeb.serialId)).then(_resp =>
          this.setState({
            enodebs: this.state.enodebs.filter(
              enb => enb.serialId != enodeb.serialId,
            ),
          }),
        );
      });
  };

  _buildEnodebFromPayload(enbSerial: string, enodeb: EnodebPayload): Enodeb {
    return {
      serialId: enbSerial,
      deviceClass: enodeb.device_class || DEFAULT_ENODEB.device_class,
      earfcndl: enodeb.earfcndl || DEFAULT_ENODEB.earfcndl,
      subframeAssignment:
        enodeb.subframe_assignment || DEFAULT_ENODEB.subframe_assignment,
      specialSubframePattern:
        enodeb.special_subframe_pattern ||
        DEFAULT_ENODEB.special_subframe_pattern,
      pci: enodeb.pci || DEFAULT_ENODEB.pci,
      bandwidthMhz: enodeb.bandwidth_mhz || DEFAULT_ENODEB.bandwidth_mhz,
      tac: enodeb.tac || DEFAULT_ENODEB.tac,
      cellId: enodeb.cell_id || DEFAULT_ENODEB.cell_id,
      transmitEnabled:
        enodeb.transmit_enabled || DEFAULT_ENODEB.transmit_enabled,
    };
  }
}

export default withStyles(styles)(withAlert(withRouter(Enodebs)));
