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
import type {subscriber} from '@fbcnms/magma-api';

import AddEditSubscriberDialog from './lte/AddEditSubscriberDialog';
import Button from '@fbcnms/ui/components/design-system/Button';
import DeleteIcon from '@material-ui/icons/Delete';
import EditIcon from '@material-ui/icons/Edit';
import IconButton from '@material-ui/core/IconButton';
import ImportSubscribersDialog from './ImportSubscribersDialog';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableFooter from '@material-ui/core/TableFooter';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Text from '@fbcnms/ui/components/design-system/Text';

import nullthrows from '@fbcnms/util/nullthrows';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {map} from 'lodash';
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
  importButton: {
    marginRight: '8px',
  },
});

type Props = ContextRouter & WithAlert & WithStyles<typeof styles> & {};

type State = {
  subscribers: Array<subscriber>,
  errorMessage: ?string,
  loading: boolean,
  showAddEditDialog: boolean,
  showImportDialog: boolean,
  editingSubscriber: ?subscriber,
  subProfiles: Set<string>,
};

class Subscribers extends React.Component<Props, State> {
  state = {
    subscribers: [],
    errorMessage: null,
    loading: true,
    showAddEditDialog: false,
    showImportDialog: false,
    editingSubscriber: null,
    subProfiles: new Set(),
  };

  componentDidMount() {
    const networkId = nullthrows(this.props.match.params.networkId);
    Promise.all([
      MagmaV1API.getLteByNetworkIdSubscribers({
        networkId,
      }),
      MagmaV1API.getLteByNetworkIdCellularEpc({networkId}),
    ])
      .then(([subscribers, epcConfigs]) => {
        const subProfiles = new Set(
          Object.keys(epcConfigs.sub_profiles || {}),
        ).add('default');

        this.setState({
          subscribers: map(subscribers, s =>
            this._buildSubscriber(s, subProfiles),
          ),
          loading: false,
          subProfiles,
        });
      })
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
          <Text variant="h5">Subscribers</Text>
          <div className={this.props.classes.buttons}>
            <Button
              className={this.props.classes.importButton}
              onClick={this.showImportDialog}>
              Import
            </Button>
            <Button onClick={this.showAddEditDialog}>Add Subscriber</Button>
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
          <Text color="error" variant="body2">
            {this.state.errorMessage ?? ''}
          </Text>
        </div>
        <AddEditSubscriberDialog
          key={(this.state.editingSubscriber || {}).id || 'new'}
          editingSubscriber={this.state.editingSubscriber}
          open={this.state.showAddEditDialog}
          onClose={this.hideDialogs}
          onSave={this.onSaveSubscriber}
          onSaveError={this.onSaveSubscriberError}
          subProfiles={Array.from(this.state.subProfiles)}
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
    MagmaV1API.getLteByNetworkIdSubscribersBySubscriberId({
      networkId: nullthrows(this.props.match.params.networkId),
      subscriberId: id,
    })
      .then(newSubscriber =>
        this.setState(state => {
          const subscribers = state.subscribers.slice(0);
          if (state.editingSubscriber) {
            const index = subscribers.indexOf(state.editingSubscriber);
            subscribers[index] = this._buildSubscriber(newSubscriber);
          } else {
            subscribers.push(this._buildSubscriber(newSubscriber));
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
    const results = await Promise.all(
      subscriberIDs.map(id =>
        MagmaV1API.getLteByNetworkIdSubscribersBySubscriberId({
          networkId: nullthrows(this.props.match.params.networkId),
          subscriberId: id,
        }),
      ),
    );
    this.setState(state => {
      const subscribers = [
        ...state.subscribers,
        ...results.map(subscriber => this._buildSubscriber(subscriber)),
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

  deleteSubscriber = sub => {
    this.props
      .confirm(`Are you sure you want to delete subscriber ${sub.id}?`)
      .then(confirmed => {
        if (confirmed) {
          MagmaV1API.deleteLteByNetworkIdSubscribersBySubscriberId({
            networkId: this.props.match.params.networkId || '',
            subscriberId: `IMSI${sub.id}`,
          })
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
  };

  _buildSubscriber(subscriber: subscriber, subProfiles?: Set<string>) {
    subProfiles = subProfiles || this.state.subProfiles;
    if (!subProfiles.has(subscriber.lte.sub_profile)) {
      subscriber.lte.sub_profile = 'default';
    }

    subscriber.id = subscriber.id.replace(/^IMSI/, '');
    return subscriber;
  }
}

type Props2 = {
  subscriber: subscriber,
  onEdit: subscriber => void,
  onDelete: subscriber => void,
};
class SubscriberTableRow extends React.Component<Props2> {
  render() {
    return (
      <TableRow>
        <TableCell>{this.props.subscriber.id}</TableCell>
        <TableCell>{this.props.subscriber.lte.state}</TableCell>
        <TableCell>{this.props.subscriber.lte.sub_profile}</TableCell>
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
