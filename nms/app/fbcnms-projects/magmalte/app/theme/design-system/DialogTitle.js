/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import CloseIcon from '@material-ui/icons/Close';
import DialogTitle from '@material-ui/core/DialogTitle';
import IconButton from '@material-ui/core/IconButton';
import React from 'react';
import Text from './Text';

import {colors} from '../default';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  closeButton: {
    color: colors.primary.white,
    padding: 0,
  },
}));

type Props = {
  label: string,
  onClose: () => void,
};

export default function CustomDialogTitle(props: Props) {
  const classes = useStyles(props);
  return (
    <DialogTitle>
      <Text variant="subtitle1">{props.label}</Text>
      <IconButton
        aria-label="close"
        className={classes.closeButton}
        onClick={props.onClose}>
        <CloseIcon />
      </IconButton>
    </DialogTitle>
  );
}
