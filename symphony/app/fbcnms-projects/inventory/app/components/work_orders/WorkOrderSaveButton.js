/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import Button from '@fbcnms/ui/components/design-system/Button';
import EditWorkOrderMutation from '../../mutations/EditWorkOrderMutation';
import FormAction from '../../../../../fbcnms-packages/fbcnms-ui/components/design-system/Form/FormAction';
import FormValidationContext from '@fbcnms/ui/components/design-system/Form/FormValidationContext';
import React, {useCallback, useContext} from 'react';
import SnackbarItem from '@fbcnms/ui/components/SnackbarItem';
import useRouter from '@fbcnms/ui/hooks/useRouter';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {removeTempIDs} from '../../common/EntUtils';
import {toPropertyInput} from '../../common/Property';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import type {
  // $FlowFixMe (T62907961) Relay flow types
  CheckListCategoryTable_list,
  // $FlowFixMe (T62907961) Relay flow types
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
  checkListCategories: CheckListCategoryTable_list,
  locationId: ?string,
};

const WorkOrderSaveButton = (props: Props) => {
  const {
    workOrder,
    properties,
    checklist,
    checkListCategories,
    locationId,
  } = props;
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
        checkListCategories: removeTempIDs(checkListCategories).map(item => ({
          id: item.id,
          title: item.title,
          description: item.description,
          checkList: removeTempIDs(item.checkList).map(item => {
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
        })),
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
    checkListCategories,
    checklist,
    enqueueError,
    history,
    match.url,
  ]);

  const validationContext = useContext(FormValidationContext);

  return (
    <FormAction disabled={validationContext.error.detected}>
      <Button tooltip={validationContext.error.message} onClick={saveWorkOrder}>
        Save
      </Button>
    </FormAction>
  );
};

export default WorkOrderSaveButton;
