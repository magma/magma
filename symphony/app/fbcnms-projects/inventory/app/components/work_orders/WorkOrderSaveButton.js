/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import Button from '@fbcnms/ui/components/design-system/Button';
import EditWorkOrderMutation from '../../mutations/EditWorkOrderMutation';
import FormValidationContext from '@fbcnms/ui/components/design-system/Form/FormValidationContext';
import React, {useCallback, useContext} from 'react';
import SnackbarItem from '@fbcnms/ui/components/SnackbarItem';
import useRouter from '@fbcnms/ui/hooks/useRouter';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {removeTempIDs} from '../../common/EntUtils';
import {toPropertyInput} from '../../common/Property';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import type {
  ChecklistViewer_checkListItems,
  WorkOrderDetails_workOrder,
} from './__generated__/WorkOrderDetails_workOrder.graphql.js';
import type {
  EditWorkOrderMutationResponse,
  EditWorkOrderMutationVariables,
} from '../../mutations/__generated__/EditWorkOrderMutation.graphql';
import type {MutationCallbacks} from '../../mutations/MutationCallbacks.js';
import type {Property} from '../../common/Property';

type Props = {
  workOrder: WorkOrderDetails_workOrder,
  properties: Array<Property>,
  checklist: ChecklistViewer_checkListItems,
  locationId: ?string,
};

const WorkOrderSaveButton = (props: Props) => {
  const {workOrder, properties, checklist, locationId} = props;
  const enqueueSnackbar = useEnqueueSnackbar();
  const {history, match} = useRouter();

  const enqueueError = useCallback(
    (message: string) => {
      enqueueSnackbar(message, {
        children: key => (
          <SnackbarItem id={key} message={message} variant="error" />
        ),
      });
    },
    [enqueueSnackbar],
  );

  const saveWorkOrder = useCallback(() => {
    const {
      id,
      name,
      description,
      ownerName,
      installDate,
      status,
      priority,
      assignee,
      project,
    } = workOrder;
    const variables: EditWorkOrderMutationVariables = {
      input: {
        id,
        name,
        description,
        ownerName,
        installDate: installDate ? installDate.toString() : null,
        status,
        priority,
        assignee,
        projectId: project?.id,
        properties: toPropertyInput(properties),
        locationId,
        checkList: removeTempIDs(checklist).map(item => {
          return {
            id: item.id,
            title: item.title,
            type: item.type,
            index: item.index,
            helpText: item.helpText,
            enumValues: item.enumValues,
            stringValue: item.stringValue,
            checked: item.checked,
          };
        }),
      },
    };
    const callbacks: MutationCallbacks<EditWorkOrderMutationResponse> = {
      onCompleted: (response, errors) => {
        if (errors && errors[0]) {
          enqueueError(errors[0].message);
        } else {
          // navigate to main page
          history.push(match.url);
        }
      },
      onError: () => {
        enqueueError('Error saving work order');
      },
    };
    ServerLogger.info(LogEvents.SAVE_WORK_ORDER_BUTTON_CLICKED, {
      source: 'work_order_details',
    });
    EditWorkOrderMutation(variables, callbacks);
  }, [
    workOrder,
    properties,
    locationId,
    checklist,
    enqueueError,
    history,
    match.url,
  ]);

  const validationContext = useContext(FormValidationContext);

  return (
    <Button
      disabled={validationContext.error.detected}
      tooltip={validationContext.error.message}
      onClick={saveWorkOrder}>
      Save
    </Button>
  );
};

export default WorkOrderSaveButton;
