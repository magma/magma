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
  RemoveWorkOrderMutationResponse,
  RemoveWorkOrderMutationVariables,
} from '../../mutations/__generated__/RemoveWorkOrderMutation.graphql';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithSnackbarProps} from 'notistack';
import type {WithStyles} from '@material-ui/core';
import type {WorkOrderDetails_workOrder} from './__generated__/WorkOrderDetails_workOrder.graphql.js';

import DeleteOutlineIcon from '@material-ui/icons/DeleteOutline';
import React from 'react';
import RemoveWorkOrderMutation from '../../mutations/RemoveWorkOrderMutation';
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
  workOrder: WorkOrderDetails_workOrder,
  onWorkOrderRemoved: () => void,
} & WithStyles<typeof styles> &
  WithAlert &
  WithSnackbarProps;

class WorkOrderDeleteButton extends React.Component<Props> {
  render() {
    const {classes, className} = this.props;
    return (
      <div className={classNames(classes.deleteButton, className)}>
        <DeleteOutlineIcon onClick={this.removeWorkOrder} />
      </div>
    );
  }

  removeWorkOrder = () => {
    ServerLogger.info(LogEvents.DELETE_WORK_ORDER_BUTTON_CLICKED, {
      source: 'work_order_details',
    });
    const {workOrder} = this.props;
    const workOrderId = workOrder.id;
    this.props
      .confirm({
        message: 'Are you sure you want to delete this work order?',
        confirmLabel: 'Delete',
      })
      .then(confirmed => {
        if (!confirmed) {
          return;
        }

        const variables: RemoveWorkOrderMutationVariables = {
          id: nullthrows(workOrderId),
        };

        const updater = store => {
          this.props.onWorkOrderRemoved();
          store.delete(workOrderId);
        };

        const callbacks: MutationCallbacks<RemoveWorkOrderMutationResponse> = {
          onCompleted: (response, errors) => {
            if (errors && errors[0]) {
              this.props.alert('Failed removing work order');
            }
          },
          onError: (_error: Error) => {
            this.props.alert('Failed removing work order');
          },
        };

        RemoveWorkOrderMutation(variables, callbacks, updater);
      });
  };
}

export default withStyles(styles)(
  withAlert(withSnackbar(WorkOrderDeleteButton)),
);
