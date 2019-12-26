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
  EditWorkOrderMutationResponse,
  EditWorkOrderMutationVariables,
} from '../../mutations/__generated__/EditWorkOrderMutation.graphql';
import type {MutationCallbacks} from '../../mutations/MutationCallbacks.js';
import type {WorkOrderLocation, WorkOrderProperties} from '../map/MapUtil';

import * as React from 'react';
import Button from '@material-ui/core/Button';
import EditWorkOrderMutation from '../../mutations/EditWorkOrderMutation';
import Text from '@fbcnms/ui/components/design-system/Text';
import TextField from '@material-ui/core/TextField';
import Typography from '@material-ui/core/Typography';
import UserTypeahead from '../typeahead/UserTypeahead';
import {formatDateForTextInput} from '@fbcnms/ui/utils/displayUtils';
import {formatMultiSelectValue} from '@fbcnms/ui/utils/displayUtils';
import {makeStyles} from '@material-ui/styles';
import {priorityValues, statusValues} from '../../common/WorkOrder';

const useStyles = makeStyles(() => ({
  fullDetails: {
    width: '100%',
  },
  root: {
    marginTop: '8px',
    minWidth: '157px',
  },
  title: {
    fontSize: '18px',
  },
  assigneeDiv: {
    display: 'flex',
    justifyContent: 'flex-start',
    alignItems: 'center',
    flexGrow: 1,
  },
  assigneeTypography: {
    marginRight: '0.35em',
  },
  gridInput: {
    display: 'inline-flex',
    margin: '5px',
    width: '250px',
  },
  dueDiv: {
    display: 'flex',
    justifyContent: 'flex-end',
    flexWrap: 'wrap',
    alignItems: 'center',
    marginRight: '0.35em',
    marginTop: '20px',
  },
  couple: {
    marginTop: '20px',
  },
}));

type Props = {
  workOrder: WorkOrderProperties,
  onWorkOrderClick?: () => void,
  displayFullDetails?: boolean,
  containerClassName?: string,
  selectedView?: string,
  onWorkOrderChanged?: (
    key: 'assignee' | 'installDate',
    value: ?string,
    workOrderId: string,
  ) => void,
};

const WorkOrderPopover = (props: Props) => {
  const {
    workOrder,
    onWorkOrderClick,
    displayFullDetails,
    selectedView,
    onWorkOrderChanged,
    containerClassName,
  } = props;
  const classes = useStyles();
  const editAssignee = selectedView === 'status' || workOrder.status === 'DONE';

  const setWorkOrderDetails = (
    key: 'assignee' | 'installDate',
    value: ?string,
  ) => {
    const variables: EditWorkOrderMutationVariables = {
      input: {
        id: workOrder.id,
        name: workOrder.name,
        ownerName: workOrder.ownerName,
        status: workOrder.status,
        priority: workOrder.priority,
        assignee: workOrder.assignee,
      },
    };
    switch (key) {
      case 'assignee':
        variables.input.assignee = value;
        break;
      case 'installDate':
        variables.input.installDate = value;
    }
    const callbacks: MutationCallbacks<EditWorkOrderMutationResponse> = {
      onCompleted: () => {
        onWorkOrderChanged && onWorkOrderChanged(key, value, workOrder.id);
      },
    };
    EditWorkOrderMutation(variables, callbacks);
  };

  const showAssignee = (assignee: string) => {
    return assignee === '' ? 'Unassigned' : assignee;
  };

  const formatLocation = (location: 'string' | WorkOrderLocation) => {
    const WorkOrderlocation =
      typeof location === 'string' ? JSON.parse(location) : location;
    return (
      WorkOrderlocation.name +
      ' (' +
      WorkOrderlocation.latitude +
      ' , ' +
      WorkOrderlocation.longitude +
      ')'
    );
  };

  return (
    <div className={containerClassName}>
      {displayFullDetails ? (
        <div className={classes.fullDetails}>
          <Text variant="h6" className={classes.title} gutterBottom>
            {workOrder.name}
          </Text>
          <Typography gutterBottom>{workOrder.description}</Typography>
          <div className={classes.couple}>
            <Typography gutterBottom>
              <strong>Owner: </strong>
              {workOrder.ownerName}
            </Typography>
            <div>
              <Typography gutterBottom>
                <strong>Status: </strong>
                {formatMultiSelectValue(statusValues, workOrder.status)}
              </Typography>
              <Typography gutterBottom>
                <strong>Priority: </strong>
                {formatMultiSelectValue(priorityValues, workOrder.priority)}
              </Typography>
            </div>
            <div className={classes.couple}>
              {!!workOrder.location && (
                <Typography gutterBottom>
                  <strong>Location: </strong>
                  {formatLocation(workOrder.location)}
                </Typography>
              )}
              <div className={classes.assigneeDiv}>
                <Typography className={classes.assigneeTypography}>
                  <strong>Assignee: </strong>
                  {editAssignee && showAssignee(workOrder.assignee)}
                </Typography>
                {!editAssignee && (
                  <UserTypeahead
                    margin="dense"
                    selectedUser={workOrder.assignee}
                    onUserSelection={user =>
                      setWorkOrderDetails('assignee', user)
                    }
                  />
                )}
              </div>
            </div>
            <div className={classes.dueDiv}>
              <div className={classes.assigneeDiv}>
                <Typography className={classes.assigneeTypography}>
                  <strong>Due: </strong>
                </Typography>
                <TextField
                  type="date"
                  label=""
                  variant="outlined"
                  margin="dense"
                  disabled={
                    workOrder.status === 'DONE' || selectedView === 'status'
                  }
                  className={classes.gridInput}
                  InputLabelProps={{
                    shrink: true,
                  }}
                  defaultValue={formatDateForTextInput(workOrder.installDate)}
                  onChange={event => {
                    setWorkOrderDetails(
                      'installDate',
                      event.target.value !== ''
                        ? new Date(event.target.value).toISOString()
                        : null,
                    );
                  }}
                />
              </div>
              {onWorkOrderClick && (
                <Button
                  variant="contained"
                  color="primary"
                  onClick={onWorkOrderClick}>
                  View
                </Button>
              )}
            </div>
          </div>
        </div>
      ) : (
        <div className={classes.root}>
          <Text variant="h6" gutterBottom>
            {workOrder.name}
          </Text>
          {showAssignee(workOrder.assignee)}
        </div>
      )}
    </div>
  );
};

export default WorkOrderPopover;
