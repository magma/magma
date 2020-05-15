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
import type {WorkOrderType} from '../../common/WorkOrder';

import Breadcrumbs from '@fbcnms/ui/components/Breadcrumbs';
import Button from '@fbcnms/ui/components/design-system/Button';
import DeleteOutlineIcon from '@material-ui/icons/DeleteOutline';
import ExpandingPanel from '@fbcnms/ui/components/ExpandingPanel';
import ExperimentalPropertyTypesTable from '../form/ExperimentalPropertyTypesTable';
import FormAction from '@fbcnms/ui/components/design-system/Form/FormAction';
import NameDescriptionSection from '@fbcnms/ui/components/NameDescriptionSection';
import PropertyTypesTableDispatcher from '../form/context/property_types/PropertyTypesTableDispatcher';
import React, {useCallback, useState} from 'react';
import SnackbarItem from '@fbcnms/ui/components/SnackbarItem';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {FormContextProvider} from '../../common/FormContext';
import {addWorkOrderType} from '../../mutations/AddWorkOrderTypeMutation';
import {deleteWorkOrderType} from '../../mutations/RemoveWorkOrderTypeMutation';
import {editWorkOrderType} from '../../mutations/EditWorkOrderTypeMutation';
import {generateTempId, isTempId} from '../../common/EntUtils';
import {makeStyles} from '@material-ui/styles';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {usePropertyTypesReducer} from '../form/context/property_types/PropertyTypesTableState';

const useStyles = makeStyles(() => ({
  root: {
    padding: '24px 16px',
    maxHeight: '100%',
    overflow: 'hidden',
    display: 'flex',
    flexDirection: 'column',
  },
  header: {
    display: 'flex',
    paddingBottom: '24px',
  },
  body: {
    overflowY: 'auto',
  },
  buttons: {
    display: 'flex',
  },
  cancelButton: {
    marginRight: '8px',
  },
  deleteButton: {
    cursor: 'pointer',
    color: symphony.palette.D400,
    width: '32px',
    height: '32px',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    marginRight: '8px',
  },
}));

type Props = $ReadOnly<{|
  open: boolean,
  onClose: () => void,
  onSave: (workOrderType: WorkOrderType) => void,
  workOrderType: ?WorkOrderType,
  ...WithAlert,
|}>;

const AddEditWorkOrderTypeCard = ({
  workOrderType,
  onClose,
  onSave,
  confirm,
}: Props) => {
  const classes = useStyles();
  const [
    editingWorkOrderType,
    setEditingWorkOrderType,
  ] = useState<WorkOrderType>({
    id: workOrderType?.id ?? generateTempId(),
    name: workOrderType?.name ?? '',
    description: workOrderType?.description,
    numberOfWorkOrders: workOrderType?.numberOfWorkOrders ?? 0,
    propertyTypes: [],
  });
  const [isSaving, setIsSaving] = useState(false);
  const [propertyTypes, propertyTypesDispatcher] = usePropertyTypesReducer(
    workOrderType?.propertyTypes ?? [],
  );
  const enqueueSnackbar = useEnqueueSnackbar();

  const onDelete = useCallback(() => {
    confirm(
      fbt(
        'Are you sure you want to delete ' +
          fbt.param('name', editingWorkOrderType.name),
        '',
      ).toString(),
    )
      .then(() => deleteWorkOrderType(editingWorkOrderType.id))
      .then(onClose)
      .catch((errorMessage: string) =>
        enqueueSnackbar(errorMessage, {
          children: key => (
            <SnackbarItem id={key} message={errorMessage} variant="error" />
          ),
        }),
      );
  }, [
    confirm,
    editingWorkOrderType.name,
    editingWorkOrderType.id,
    onClose,
    enqueueSnackbar,
  ]);

  const nameChanged = name =>
    setEditingWorkOrderType(workOrder => ({
      ...workOrder,
      name,
    }));

  const descriptionChanged = description =>
    setEditingWorkOrderType(workOrder => ({
      ...workOrder,
      description,
    }));

  const onSaveClicked = () => {
    setIsSaving(true);
    const workOrderToSave = {...editingWorkOrderType, propertyTypes};
    const saveAction = isTempId(editingWorkOrderType.id)
      ? addWorkOrderType
      : editWorkOrderType;
    saveAction(workOrderToSave)
      .then(() => onSave(workOrderToSave))
      .catch((errorMessage: string) =>
        enqueueSnackbar(errorMessage, {
          children: key => (
            <SnackbarItem id={key} message={errorMessage} variant="error" />
          ),
        }),
      )
      .finally(() => setIsSaving(false));
  };

  return (
    <FormContextProvider>
      <div className={classes.root}>
        <div className={classes.header}>
          <Breadcrumbs
            breadcrumbs={[
              {
                id: 'wo_templates',
                name: 'Work Order Templates',
                onClick: onClose,
              },
              workOrderType
                ? {
                    id: workOrderType.id,
                    name: workOrderType.name,
                  }
                : {
                    id: 'new_wo_type',
                    name: `${fbt('New work order template', '')}`,
                  },
            ]}
            size="large"
          />
          <div className={classes.buttons}>
            {workOrderType != null && (
              <FormAction>
                <Button
                  className={classes.deleteButton}
                  variant="text"
                  skin="gray"
                  onClick={onDelete}>
                  <DeleteOutlineIcon />
                </Button>
              </FormAction>
            )}
            <Button
              className={classes.cancelButton}
              skin="regular"
              onClick={onClose}>
              Cancel
            </Button>
            <FormAction disableOnFromError={true} disabled={isSaving}>
              <Button onClick={onSaveClicked}>Save</Button>
            </FormAction>
          </div>
        </div>
        <div className={classes.body}>
          <ExpandingPanel title="Details">
            <NameDescriptionSection
              title="Title"
              name={editingWorkOrderType.name ?? ''}
              namePlaceholder={`${fbt('New work order template', '')}`}
              description={editingWorkOrderType.description ?? ''}
              descriptionPlaceholder={`${fbt(
                'Write a description if you want it to appear whenever this template of work order is created',
                '',
              )}`}
              onNameChange={nameChanged}
              onDescriptionChange={descriptionChanged}
            />
          </ExpandingPanel>
          <ExpandingPanel title="Properties">
            <PropertyTypesTableDispatcher.Provider
              value={propertyTypesDispatcher}>
              <ExperimentalPropertyTypesTable
                supportDelete={true}
                propertyTypes={propertyTypes}
              />
            </PropertyTypesTableDispatcher.Provider>
          </ExpandingPanel>
        </div>
      </div>
    </FormContextProvider>
  );
};

export default withAlert(AddEditWorkOrderTypeCard);
