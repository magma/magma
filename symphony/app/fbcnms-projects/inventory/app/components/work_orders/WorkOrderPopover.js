/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {BasicLocation} from '../../common/Location';
import type {
  EditWorkOrderMutationResponse,
  EditWorkOrderMutationVariables,
} from '../../mutations/__generated__/EditWorkOrderMutation.graphql';
import type {MutationCallbacks} from '../../mutations/MutationCallbacks.js';
import type {ShortUser} from '../../common/EntUtils';
import type {WorkOrderProperties} from '../map/MapUtil';

import * as React from 'react';
import DateTimeFormat from '../../common/DateTimeFormat';
import EditWorkOrderMutation from '../../mutations/EditWorkOrderMutation';
import Strings from '@fbcnms/strings/Strings';
import Text from '@fbcnms/ui/components/design-system/Text';
import UserTypeahead from '../typeahead/UserTypeahead';
import classNames from 'classnames';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import {InventoryAPIUrls} from '../../common/InventoryAPI';
import {Link} from 'react-router-dom';
import {formatMultiSelectValue} from '@fbcnms/ui/utils/displayUtils';
import {makeStyles} from '@material-ui/styles';
import {priorityValues, statusValues} from '../../common/WorkOrder';

const useStyles = makeStyles(() => ({
  fullDetails: {
    width: '100%',
    padding: '24px',
  },
  quickPeek: {
    marginTop: '8px',
    minWidth: '157px',
  },
  notUnderlinedLink: {
    textDecoration: 'none',
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
  section: {
    '&:not(:first-child)': {
      marginTop: '20px',
    },
    '&>*': {
      '&:not(:first-child)': {
        marginTop: '4px',
      },
    },
  },
  field: {
    display: 'flex',
    alignItems: 'baseline',
    '&>:not(:first-child)': {
      marginLeft: '2px',
    },
    '&>:last-child': {
      flexGrow: '1',
    },
  },
  trunckedContent: {
    '-webkit-line-clamp': '2',
    overflow: 'hidden',
    display: '-webkit-box',
    '-webkit-box-orient': 'vertical',
  },
  fieldBox: {
    display: 'block',
    '&:not(:first-child)': {
      marginTop: '8px',
    },
    '&>*': {
      display: 'inline-flex',
      background: symphony.palette.background,
      borderRadius: '4px',
      padding: '3px 8px',
    },
  },
}));

type Props = {
  workOrder: WorkOrderProperties,
  onWorkOrderClick?: () => void,
  displayFullDetails?: boolean,
  containerClassName?: string,
  selectedView?: string,
  onWorkOrderChanged?: (
    key: 'assigneeId' | 'installDate',
    value: ?string,
    workOrderId: string,
  ) => void,
};

const WorkOrderPopover = (props: Props) => {
  const {
    workOrder,
    displayFullDetails,
    selectedView,
    onWorkOrderChanged,
    containerClassName,
  } = props;
  const classes = useStyles();
  const viewMode = selectedView === 'status' || workOrder.status === 'DONE';

  const setWorkOrderDetails = (
    key: 'assigneeId' | 'installDate',
    value: ?string,
  ) => {
    const variables: EditWorkOrderMutationVariables = {
      input: {
        id: workOrder.id,
        name: workOrder.name,
        ownerId: workOrder.owner.id,
        status: workOrder.status,
        priority: workOrder.priority,
        assigneeId: workOrder.assignedTo?.id,
      },
    };
    switch (key) {
      case 'assigneeId':
        variables.input.assigneeId = value;
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

  const showAssignee = (assignee: ?ShortUser) => {
    return assignee?.email || Strings.common.unassignedItem;
  };

  const woHeader = (
    <Link
      className={classes.notUnderlinedLink}
      to={InventoryAPIUrls.workorder(workOrder.id)}>
      <Text variant="subtitle1" color="primary">
        {workOrder.name}
      </Text>
    </Link>
  );

  const nameAndCoordinates = (locationInput: BasicLocation) => {
    return `${locationInput.name} (${locationInput.latitude}, ${locationInput.longitude})`;
  };

  return (
    <div className={containerClassName}>
      {displayFullDetails ? (
        <div className={classes.fullDetails}>
          {woHeader}
          <div>
            <Text
              title={workOrder.description}
              variant="body2"
              className={classNames(classes.field, classes.trunckedContent)}>
              {workOrder.description}
            </Text>
          </div>
          <div className={classes.section}>
            <Text variant="body2" className={classes.field}>
              <strong>
                {fbt('Assignee:', 'Work Order card "Assignee" field title')}
              </strong>
              {!!viewMode ? (
                <span>{showAssignee(workOrder.assignedTo)}</span>
              ) : (
                <UserTypeahead
                  margin="dense"
                  selectedUser={workOrder.assignedTo}
                  onUserSelection={user =>
                    setWorkOrderDetails('assigneeId', user?.id)
                  }
                />
              )}
            </Text>
            {!!workOrder.location && (
              <Text variant="body2" className={classes.field}>
                <strong>
                  {fbt('Location:', 'Work Order card "Location" field title')}
                </strong>
                <span>{nameAndCoordinates(workOrder.location)}</span>
              </Text>
            )}
          </div>
          <div className={classes.section}>
            <div className={classes.fieldBox}>
              <Text variant="body2" className={classes.field}>
                <strong>
                  {fbt('Status:', 'Work Order card "Status" field title')}
                </strong>
                <span>
                  {formatMultiSelectValue(statusValues, workOrder.status)}
                </span>
              </Text>
            </div>
            <div className={classes.fieldBox}>
              <Text variant="body2" className={classes.field}>
                <strong>
                  {fbt('Priority:', 'Work Order card "Priority" field title')}
                </strong>
                <span>
                  {formatMultiSelectValue(priorityValues, workOrder.priority)}
                </span>
              </Text>
            </div>
            <div className={classes.fieldBox}>
              <Text variant="body2" className={classes.field}>
                <strong>
                  {fbt('Due:', 'Work Order card "Due" field title')}
                </strong>
                <span>
                  {DateTimeFormat.dateTime(
                    workOrder.installDate,
                    Strings.common.emptyField,
                  )}
                </span>
              </Text>
            </div>
          </div>
        </div>
      ) : (
        <div className={classes.quickPeek}>
          {woHeader}
          <div>{showAssignee(workOrder.assignedTo)}</div>
        </div>
      )}
    </div>
  );
};

export default WorkOrderPopover;
