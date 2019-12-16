/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {MutationCallbacks} from '../../mutations/MutationCallbacks.js';
import type {
  RemoveServiceMutationResponse,
  RemoveServiceMutationVariables,
} from '../../mutations/__generated__/RemoveServiceMutation.graphql';
import type {Service} from '../../common/Service';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithSnackbarProps} from 'notistack';
import type {WithStyles} from '@material-ui/core';

import DeleteOutlineIcon from '@material-ui/icons/DeleteOutline';
import React from 'react';
import RemoveServiceMutation from '../../mutations/RemoveServiceMutation';
import SymphonyTheme from '@fbcnms/ui/theme/symphony';
import classNames from 'classnames';
import nullthrows from '@fbcnms/util/nullthrows';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {withSnackbar} from 'notistack';
import {withStyles} from '@material-ui/core/styles';

const styles = _theme => ({
  deleteButton: {
    cursor: 'pointer',
    color: SymphonyTheme.palette.D400,
    width: '32px',
    height: '32px',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
  },
});

type Props = {
  className?: string,
  service: Service,
  onServiceRemoved: () => void,
} & WithStyles<typeof styles> &
  WithAlert &
  WithSnackbarProps;

class ServiceDeleteButton extends React.Component<Props> {
  render() {
    const {classes, className} = this.props;
    return (
      <div className={classNames(classes.deleteButton, className)}>
        <DeleteOutlineIcon onClick={this.removeService} />
      </div>
    );
  }

  removeService = () => {
    ServerLogger.info(LogEvents.DELETE_SERVICE_BUTTON_CLICKED, {
      source: 'service_details',
    });
    const {service} = this.props;
    const serviceId = service.id;
    this.props
      .confirm({
        message: 'Are you sure you want to delete this service?',
        confirmLabel: 'Delete',
      })
      .then(confirmed => {
        if (!confirmed) {
          return;
        }

        const variables: RemoveServiceMutationVariables = {
          id: nullthrows(serviceId),
        };

        const updater = store => {
          this.props.onServiceRemoved();
          store.delete(serviceId);
        };

        const callbacks: MutationCallbacks<RemoveServiceMutationResponse> = {
          onCompleted: (response, errors) => {
            if (errors && errors[0]) {
              this.props.alert('Failed removing service');
            }
          },
          onError: (_error: Error) => {
            this.props.alert('Failed removing service');
          },
        };

        RemoveServiceMutation(variables, callbacks, updater);
      });
  };
}

export default withStyles(styles)(withAlert(withSnackbar(ServiceDeleteButton)));
