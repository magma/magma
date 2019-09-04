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
import type {Subscriber} from './lte/AddEditSubscriberDialog';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithStyles} from '@material-ui/core';

import AddEditSubscriberDialog from './lte/AddEditSubscriberDialog';
import Button from '@material-ui/core/Button';
import DeleteIcon from '@material-ui/icons/Delete';
import EditIcon from '@material-ui/icons/Edit';
import IconButton from '@material-ui/core/IconButton';
import ImportSubscribersDialog from './ImportSubscribersDialog';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableFooter from '@material-ui/core/TableFooter';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Typography from '@material-ui/core/Typography';
import axios from 'axios';

import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {MagmaAPIUrls} from '../common/MagmaAPI';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

const styles = theme => ({
  header: {
    margin: '10px',
    display: 'flex',
    justifyContent: 'space-between',
  },
  buttons: {
    display: 'flex',
    justifyContent: 'flex-end',
    flexDirection: 'row',
  },
  paper: {
    margin: theme.spacing(3),
  },
});

type SubProfiles = {
  [string]: {max_dl_bit_rate?: number, max_ul_bit_rate?: number},
};

type Props = ContextRouter & WithAlert & WithStyles<typeof styles> & {};

type State = {
  subscribers: Array<Subscriber>,
  errorMessage: ?string,
  loading: boolean,
  showAddEditDialog: boolean,
  showImportDialog: boolean,
  editingSubscriber: ?Subscriber,
  subProfiles: SubProfiles,
};

class Subscribers extends React.Component<Props, State> {
  state = {
    subscribers: [],
    errorMessage: null,
    loading: true,
    showAddEditDialog: false,
    showImportDialog: false,
    editingSubscriber: null,
    subProfiles: {},
  };

  componentDidMount() {
    axios
      .all([
        axios.get(MagmaAPIUrls.subscribers(this.props.match), {
          params: {fields: 'all'},
        }),
        axios.get(
          MagmaAPIUrls.networkConfigsForType(this.props.match, 'cellular'),
        ),
      ])
      .then(
        axios.spread((response1, response2) => {
          let subProfiles = (response2.data.epc || {}).sub_profiles || {};
          subProfiles = {...subProfiles};
          if (!subProfiles.default) {
            subProfiles.default = {};
          }

          const subscribers = (Object.values(response1.data): Array<any>);
          this.setState({
            subscribers: (subscribers: Array<Subscriber>).map(s =>
              this._buildSubscriber(s, subProfiles),
            ),
            loading: false,
            subProfiles,
          });
        }),
      )
      .catch((error, _) =>
        this.setState({
          errorMessage: error.response.data.message.toString(),
          loading: false,
        }),
      );
  }

  render() {
    const rows = this.state.subscribers.map(row => (
      <SubscriberTableRow
        key={row.id}
        subscriber={row}
        onEdit={this.editSubscriber}
        onDelete={this.deleteSubscriber}
      />
    ));

    return (
      <div className={this.props.classes.paper}>
        <div className={this.props.classes.header}>
          <Typography variant="h5">Subscribers</Typography>
          <div className={this.props.classes.buttons}>
            <Button
              style={{marginRight: '10px'}}
              variant="contained"
              color="primary"
              onClick={this.showImportDialog}>
              Import
            </Button>
            <Button
              variant="contained"
              color="primary"
              onClick={this.showAddEditDialog}>
              Add Subscriber
            </Button>
          </div>
        </div>
        <Paper elevation={2}>
          {this.state.loading ? (
            <LoadingFiller />
          ) : (
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>IMSI</TableCell>
                  <TableCell>LTE Subscription State</TableCell>
                  <TableCell>Data Plan</TableCell>
                  <TableCell />
                </TableRow>
              </TableHead>
              <TableBody>{rows}</TableBody>
              <TableFooter
                style={
                  !this.state.loading &&
                  this.state.subscribers.length === 0 &&
                  this.state.errorMessage === null
                    ? {}
                    : {display: 'none'}
                }>
                <TableRow>
                  <TableCell colSpan="3">No subscribers found</TableCell>
                </TableRow>
              </TableFooter>
            </Table>
          )}
        </Paper>
        <div style={this.state.errorMessage !== null ? {} : {display: 'none'}}>
          <Typography color="error" variant="body2">
            {this.state.errorMessage}
          </Typography>
        </div>
        <AddEditSubscriberDialog
          key={(this.state.editingSubscriber || {}).id || 'new'}
          editingSubscriber={this.state.editingSubscriber}
          open={this.state.showAddEditDialog}
          onClose={this.hideDialogs}
          onSave={this.onSaveSubscriber}
          onSaveError={this.onSaveSubscriberError}
          subProfiles={Object.keys(this.state.subProfiles)}
        />
        <ImportSubscribersDialog
          open={this.state.showImportDialog}
          onClose={this.hideDialogs}
          onSave={this.onBulkUpload}
          onSaveError={this.onBulkUploadError}
        />
      </div>
    );
  }

  showAddEditDialog = () => this.setState({showAddEditDialog: true});
  showImportDialog = () => this.setState({showImportDialog: true});
  hideDialogs = () =>
    this.setState({
      showAddEditDialog: false,
      showImportDialog: false,
      editingSubscriber: null,
    });

  onSaveSubscriber = id => {
    axios
      .get(MagmaAPIUrls.subscriber(this.props.match, id))
      .then(response =>
        this.setState(state => {
          const subscribers = state.subscribers.slice(0);
          if (state.editingSubscriber) {
            const index = subscribers.indexOf(state.editingSubscriber);
            subscribers[index] = this._buildSubscriber(response.data);
          } else {
            subscribers.push(this._buildSubscriber(response.data));
          }
          return {subscribers, showAddEditDialog: false};
        }),
      )
      .catch(this.props.alert);
  };

  onSaveSubscriberError = (reason: any) => {
    this.props.alert(reason.response.data.message);
  };

  onBulkUpload = async (subscriberIDs: Array<string>) => {
    const responses = await Promise.all(
      subscriberIDs.map(id =>
        axios.get(MagmaAPIUrls.subscriber(this.props.match, id)),
      ),
    );
    this.setState(state => {
      const subscribers = [
        ...state.subscribers,
        ...responses.map(response => this._buildSubscriber(response.data)),
      ];
      return {subscribers, showImportDialog: false};
    });
  };

  onBulkUploadError = (failureIDs: Array<string>) => {
    this.props.alert(
      'Error adding the following subscribers: ' + failureIDs.join(', '),
    );
  };

  editSubscriber = subscriber =>
    this.setState({editingSubscriber: subscriber, showAddEditDialog: true});

  deleteSubscriber = sub =>
    this.props
      .confirm(`Are you sure you want to delete subscriber ${sub.id}?`)
      .then(confirmed => {
        if (confirmed) {
          axios
            .delete(MagmaAPIUrls.subscriber(this.props.match, 'IMSI' + sub.id))
            .then(_resp =>
              this.setState({
                subscribers: this.state.subscribers.filter(
                  s => s.id !== sub.id,
                ),
              }),
            )
            .catch(this.props.alert);
        }
      });

  _buildSubscriber(subscriber: Subscriber, subProfiles?: SubProfiles) {
    subProfiles = subProfiles || this.state.subProfiles;
    if (!(subscriber.sub_profile in subProfiles)) {
      subscriber.sub_profile = 'default';
    }

    subscriber.id = subscriber.id.replace(/^IMSI/, '');
    return subscriber;
  }
}

type Props2 = {
  subscriber: Subscriber,
  onEdit: Subscriber => void,
  onDelete: Subscriber => any,
};
class SubscriberTableRow extends React.Component<Props2> {
  render() {
    return (
      <TableRow>
        <TableCell>{this.props.subscriber.id}</TableCell>
        <TableCell>{this.props.subscriber.lte.state}</TableCell>
        <TableCell>{this.props.subscriber.sub_profile}</TableCell>
        <TableCell>
          <IconButton onClick={this.onEdit}>
            <EditIcon />
          </IconButton>
          <IconButton onClick={this.onDelete}>
            <DeleteIcon />
          </IconButton>
        </TableCell>
      </TableRow>
    );
  }

  onEdit = () => this.props.onEdit(this.props.subscriber);
  onDelete = () => this.props.onDelete(this.props.subscriber);
}

export default withStyles(styles)(withAlert(withRouter(Subscribers)));
