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
  EditEquipmentPortMutationResponse,
  EditEquipmentPortMutationVariables,
} from '../../mutations/__generated__/EditEquipmentPortMutation.graphql';
import type {EquipmentPort} from '../../common/Equipment';
import type {MutationCallbacks} from '../../mutations/MutationCallbacks.js';
import type {Theme} from '@material-ui/core';

import Button from '@fbcnms/ui/components/design-system/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import EditEquipmentPortMutation from '../../mutations/EditEquipmentPortMutation';
import PropertiesAddEditSection from '../form/PropertiesAddEditSection';
import React, {useCallback, useState} from 'react';
import SnackbarItem from '@fbcnms/ui/components/SnackbarItem';
import Text from '@fbcnms/ui/components/design-system/Text';
import update from 'immutability-helper';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {getInitialPropertyFromType} from '../../common/PropertyType';
import {
  getNonInstancePropertyTypes,
  sortPropertiesByIndex,
  toPropertyInput,
} from '../../common/Property';
import {makeStyles} from '@material-ui/styles';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';

const useStyles = makeStyles((theme: Theme) => ({
  button: {
    marginTop: theme.spacing(1),
    marginRight: theme.spacing(1),
  },
  content: {
    width: '80%',
    height: '60vh',
  },
  root: {
    display: 'flex',
    flexDirection: 'row',
  },
  titleText: {
    borderBottom: `1px solid ${theme.palette.gray1}`,
    padding: '8px 16px',
  },
  dialogTitle: {
    padding: 0,
    marginBottom: '8px',
    marginTop: '8px',
  },
}));

type Props = {
  port: EquipmentPort,
  onClose: void => void,
};

const getEditingPort = (port: EquipmentPort): EquipmentPort => {
  let initialProps = port.properties ?? [];
  const propertyTypes = port.definition.portType?.propertyTypes;
  if (propertyTypes) {
    initialProps = [
      ...initialProps,
      ...getNonInstancePropertyTypes(
        initialProps,
        propertyTypes,
      ).map(propType => getInitialPropertyFromType(propType)),
    ].sort(sortPropertiesByIndex);
  }

  return {
    ...port,
    properties: initialProps,
  };
};

const PortEditDialog = (props: Props) => {
  const {onClose} = props;
  const classes = useStyles();
  const [editingPort, setEditingPort] = useState(getEditingPort(props.port));
  const [isSubmitting, setIsSubmitting] = useState(false);
  const enqueueSnackbar = useEnqueueSnackbar();
  const propertyChangedHandler = useCallback(
    index => property =>
      setEditingPort(
        update(editingPort, {
          properties: {[index]: {$set: property}},
        }),
      ),
    [editingPort],
  );
  const onSave = useCallback(() => {
    setIsSubmitting(true);
    ServerLogger.info(LogEvents.SAVE_EQUIPMENT_PORT_BUTTON_CLICKED);
    const variables: EditEquipmentPortMutationVariables = {
      input: {
        side: {
          port: editingPort.definition.id,
          equipment: editingPort.parentEquipment.id,
        },
        properties: toPropertyInput(editingPort.properties),
      },
    };

    const callbacks: MutationCallbacks<EditEquipmentPortMutationResponse> = {
      onCompleted: (_, errors) => {
        if (errors && errors[0]) {
          enqueueSnackbar(errors[0].message, {
            children: key => (
              <SnackbarItem
                id={key}
                message={errors[0].message}
                variant="error"
              />
            ),
          });
        }
        onClose();
      },
      onError: () => {
        onClose();
      },
    };
    EditEquipmentPortMutation(variables, callbacks);
  }, [editingPort, enqueueSnackbar, onClose]);

  return (
    <Dialog fullWidth={true} maxWidth="md" open={true} onClose={props.onClose}>
      <DialogTitle disableTypography={true} className={classes.dialogTitle}>
        <Text variant="h6" className={classes.titleText}>
          {`Editing Port: ${editingPort.definition.name}`}
        </Text>
      </DialogTitle>
      <DialogContent>
        {editingPort.properties.length > 0 ? (
          <PropertiesAddEditSection
            properties={editingPort.properties}
            onChange={index => propertyChangedHandler(index)}
          />
        ) : null}
      </DialogContent>
      <DialogActions>
        <Button onClick={onSave} disabled={isSubmitting}>
          Save
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default PortEditDialog;
