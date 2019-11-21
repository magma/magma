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
  EditLinkMutationResponse,
  EditLinkMutationVariables,
} from '../../mutations/__generated__/EditLinkMutation.graphql';
import type {Link} from '../../common/Equipment';
import type {MutationCallbacks} from '../../mutations/MutationCallbacks.js';
import type {Theme} from '@material-ui/core';

import Button from '@fbcnms/ui/components/design-system/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import EditLinkMutation from '../../mutations/EditLinkMutation';
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
import {uniqBy} from 'lodash';
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
  link: Link,
  onClose: void => void,
};

const getEditingLink = (link: Link): Link => {
  let initialLinkProps = link.properties ?? [];
  const linkPropertyTypes = uniqBy(
    link.ports
      .slice()
      .map(port => port.definition.portType?.linkPropertyTypes ?? [])
      .flatMap(i => i),
    'id',
  );
  initialLinkProps = [
    ...initialLinkProps,
    ...getNonInstancePropertyTypes(initialLinkProps, linkPropertyTypes).map(
      propType => getInitialPropertyFromType(propType),
    ),
  ].sort(sortPropertiesByIndex);

  return {
    ...link,
    properties: initialLinkProps,
  };
};

const LinkEditDialog = (props: Props) => {
  const {onClose} = props;
  const classes = useStyles();
  const [editingLink, setEditingLink] = useState(getEditingLink(props.link));
  const [isSubmitting, setIsSubmitting] = useState(false);
  const enqueueSnackbar = useEnqueueSnackbar();
  const propertyChangedHandler = useCallback(
    index => property =>
      setEditingLink(
        update(editingLink, {
          properties: {[index]: {$set: property}},
        }),
      ),
    [editingLink],
  );
  const onSave = useCallback(() => {
    setIsSubmitting(true);
    ServerLogger.info(LogEvents.SAVE_LINK_BUTTON_CLICKED);
    const variables: EditLinkMutationVariables = {
      input: {
        id: editingLink.id,
        properties: toPropertyInput(editingLink.properties),
      },
    };

    const callbacks: MutationCallbacks<EditLinkMutationResponse> = {
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
    EditLinkMutation(variables, callbacks);
  }, [editingLink, enqueueSnackbar, onClose]);

  return (
    <Dialog fullWidth={true} maxWidth="md" open={true} onClose={onClose}>
      <DialogTitle disableTypography={true} className={classes.dialogTitle}>
        <Text variant="h6" className={classes.titleText}>
          Editing Link
        </Text>
      </DialogTitle>
      <DialogContent>
        {editingLink.properties.length > 0 ? (
          <PropertiesAddEditSection
            properties={editingLink.properties}
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

export default LinkEditDialog;
