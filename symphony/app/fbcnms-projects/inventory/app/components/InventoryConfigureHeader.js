/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {ButtonProps} from '@fbcnms/ui/components/design-system/Button';
import type {PermissionEnforcement} from '../common/FormContext';
import type {PermissionHandlingProps} from '@fbcnms/ui/components/design-system/Form/FormAction';

import * as React from 'react';
import Button from '@fbcnms/ui/components/design-system/Button';
import ConfigureTitle from '@fbcnms/ui/components/ConfigureTitle';
import FormAction from '@fbcnms/ui/components/design-system/Form/FormAction';
import classNames from 'classnames';
import {FormContextProvider} from '../common/FormContext';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  headerRoot: {
    paddingBottom: '20px',
    display: 'flex',
    alignItems: 'flex-end',
  },
  title: {
    flexGrow: 1,
  },
  actionButton: {
    marginLeft: '8px',
  },
}));

export type DisplayOptionTypes = 'table' | 'map';
export const DisplayOptions = {
  table: 'table',
  map: 'map',
};

type ActionButtonProps = {
  title: string,
  action: () => void,
  ...PermissionHandlingProps,
  ...ButtonProps,
};

type Props = $ReadOnly<{|
  permissions: PermissionEnforcement,
  title: string,
  subtitle?: string,
  className?: ?string,
  onViewToggleClicked?: (id: string) => void,
  actionButtons?: Array<ActionButtonProps>,
|}>;

const InventoryConfigureHeader = (props: Props) => {
  const {permissions, className, actionButtons = []} = props;
  const classes = useStyles();

  return (
    <div className={classNames(classes.headerRoot, className)}>
      <FormContextProvider permissions={permissions}>
        <ConfigureTitle
          className={classes.title}
          title={props.title}
          subtitle={props.subtitle}
        />
        {actionButtons.map(actionButton => {
          const {
            title,
            action,
            ignorePermissions,
            hideOnMissingPermissions,
            ...otherButtonProps
          } = actionButton;
          return (
            <FormAction
              ignorePermissions={ignorePermissions}
              hideOnMissingPermissions={hideOnMissingPermissions}>
              <Button
                {...otherButtonProps}
                className={classes.actionButton}
                onClick={action}>
                {title}
              </Button>
            </FormAction>
          );
        })}
      </FormContextProvider>
    </div>
  );
};

export default InventoryConfigureHeader;
