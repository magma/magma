/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import * as React from 'react';
import DoneIcon from '@material-ui/icons/Done';
import Text from './design-system/Text';
import symphony from '../theme/symphony';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(_theme => ({
  root: {
    display: 'flex',
    flexDirection: 'row',
    borderColor: symphony.palette.G600,
    borderWidth: '1px',
    borderRadius: '4px',
    borderStyle: 'solid',
    margin: 'auto 8px',
    alignItems: 'center',
  },
  icon: {
    margin: '2px 8px',
    color: symphony.palette.G600,
  },
}));

type Props = {
  message: string,
};

const DialogConfirm = ({message}: Props) => {
  const classes = useStyles();
  return (
    <div className={classes.root}>
      <DoneIcon className={classes.icon} />
      <Text variant="subtitle2" color="regular">
        {message}
      </Text>
    </div>
  );
};

export default DialogConfirm;
