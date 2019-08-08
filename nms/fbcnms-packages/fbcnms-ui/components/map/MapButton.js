/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import * as React from 'react';
import ToggleButton from '@material-ui/lab/ToggleButton';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  button: {
    background: theme.palette.background.default,
    color: 'black',
    borderRight: '1px solid #ddd',
    borderRadius: '4px',
    width: '30px',
    height: '30px',
    border: 0,
  },
}));

type Props = {
  onClick: () => void,
  children: ?React.Node,
};

const MapButton = (props: Props) => {
  const {onClick} = props;
  const classes = useStyles();

  return (
    <ToggleButton value={1} className={classes.button} onClick={onClick}>
      {props.children}
    </ToggleButton>
  );
};

export default MapButton;
