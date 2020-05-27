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
  InventoryEntName,
  InventoryPermissionEnforcement,
} from '../admin/userManagement/utils/usePermissions';

import * as React from 'react';
import FormActionWithPermissions from '../../common/FormActionWithPermissions';
import Grid from '@material-ui/core/Grid';
import IconButton from '@fbcnms/ui/components/design-system/IconButton';
import Text from '@fbcnms/ui/components/design-system/Text';
import Tooltip from '@material-ui/core/Tooltip';
import classNames from 'classnames';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import {DeleteIcon, EditIcon} from '@fbcnms/ui/components/design-system/Icons';
import {makeStyles} from '@material-ui/styles';

type Props = $ReadOnly<{|
  icon: React.Node,
  entityName: InventoryEntName,
  instanceCount: number,
  instanceNamePlural: string,
  instanceNameSingular: string,
  name: string,
  onDelete: () => void,
  onEdit?: ?() => void,
  allowDelete?: ?boolean,
|}>;

const useStyles = makeStyles(() => ({
  inline: {
    display: 'flex',
    alignItems: 'center',
    flexGrow: 1,
  },
  root: {
    flexGrow: 1,
  },
  iconContainer: {
    borderRadius: '50%',
    marginRight: '16px',
    backgroundColor: symphony.palette.D50,
    color: symphony.palette.D500,
    width: '48px',
    height: '48px',
    display: 'flex',
    flexShrink: 0,
    justifyContent: 'center',
    alignItems: 'center',
    ...symphony.typography.h5,
  },
  iconButton: {
    marginLeft: '16px',
  },
  boldText: {
    fontWeight: 'bold',
  },
  text: {
    color: '#4d4d4e',
  },
  counter: {
    marginRight: 'auto',
  },
  actionButtons: {
    display: 'flex',
    flexDirection: 'row-reverse',
    alignItems: 'center',
    flexGrow: 1,
  },
  checkbox: {
    margin: '12px',
  },
}));

function ConfigureExpansionPanel(props: Props) {
  const {
    icon,
    entityName,
    instanceCount,
    instanceNamePlural,
    instanceNameSingular,
    name,
    onEdit,
    allowDelete,
    onDelete,
  } = props;
  const classes = useStyles();

  const editButtonPermissions: InventoryPermissionEnforcement = {
    entity: entityName,
    action: 'update',
  };

  return (
    <Grid
      container
      className={classes.root}
      direction="row"
      justify="space-between"
      alignItems="center">
      <Grid item xs>
        <div className={classes.inline}>
          <div className={classes.iconContainer}>{icon}</div>
          <Text
            className={classNames(classes.text, classes.boldText)}
            variant="subtitle1">
            {name}
          </Text>
        </div>
      </Grid>
      <Grid item xs>
        <Text
          className={classNames(classes.text, classes.counter)}
          variant="body2">
          {`${instanceCount.toLocaleString()}
                ${
                  instanceCount == 1 ? instanceNameSingular : instanceNamePlural
                } of this type`}
        </Text>
      </Grid>
      <Grid item xs>
        <div className={classes.actionButtons}>
          <DeleteButton
            entityName={entityName}
            instanceCount={instanceCount}
            onDelete={onDelete}
            allowDelete={allowDelete}
          />
          {onEdit && (
            <FormActionWithPermissions permissions={editButtonPermissions}>
              <IconButton
                skin="primary"
                className={classes.iconButton}
                onClick={onEdit}
                icon={EditIcon}
              />
            </FormActionWithPermissions>
          )}
        </div>
      </Grid>
    </Grid>
  );
}

type DeleteButtonProps = $ReadOnly<{|
  entityName: InventoryEntName,
  instanceCount: number,
  allowDelete?: ?boolean,
  onDelete: () => void,
|}>;

function DeleteButton(props: DeleteButtonProps) {
  const {entityName, instanceCount, allowDelete, onDelete} = props;
  const classes = useStyles();

  const disabled =
    allowDelete !== undefined && allowDelete !== null
      ? !allowDelete
      : instanceCount > 0;

  const deleteButtonPermissions: InventoryPermissionEnforcement = {
    entity: entityName,
    action: 'delete',
  };

  const deleteButton = (
    <FormActionWithPermissions permissions={deleteButtonPermissions}>
      <IconButton
        className={classes.iconButton}
        skin="primary"
        disabled={disabled}
        onClick={onDelete}
        icon={DeleteIcon}
      />
    </FormActionWithPermissions>
  );
  const tooltip = fbt('Cannot delete a type that is in use', '');

  return disabled ? (
    <Tooltip title={tooltip}>{deleteButton}</Tooltip>
  ) : (
    deleteButton
  );
}

export default ConfigureExpansionPanel;
