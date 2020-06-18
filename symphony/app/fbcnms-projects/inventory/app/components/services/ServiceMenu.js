/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import type {MenuOption} from '../OptionsPopoverButton';

import * as React from 'react';
import AddCircleOutlineIcon from '@material-ui/icons/AddCircleOutline';
import Dialog from '@material-ui/core/Dialog';
import FormAction from '@fbcnms/ui/components/design-system/Form/FormAction';
import OptionsPopoverButton from '../OptionsPopoverButton';
import symphony from '@fbcnms/ui/theme/symphony';
import {makeStyles} from '@material-ui/styles';

type Props = {
  items: Array<MenuOption>,
  isOpen: boolean,
  onClose: () => void,
  children: React.Node,
};

const useStyles = makeStyles(() => ({
  addIcon: {
    fill: symphony.palette.primary,
    marginRight: '8px',
  },
  dialog: {
    width: '80%',
    maxWidth: '1280px',
    height: '90%',
    maxHeight: '800px',
  },
}));

const ServiceMenu = (props: Props) => {
  const classes = useStyles();
  const {items, isOpen, onClose, children} = props;

  return (
    <FormAction>
      <OptionsPopoverButton
        options={items}
        menuIcon={<AddCircleOutlineIcon className={classes.addIcon} />}
      />
      <Dialog
        open={isOpen}
        onClose={onClose}
        maxWidth={false}
        fullWidth={true}
        classes={{paperFullWidth: classes.dialog}}>
        {children}
      </Dialog>
    </FormAction>
  );
};

export default ServiceMenu;
