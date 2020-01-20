/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import BubbleChartIcon from '@material-ui/icons/BubbleChart';
import ListAltIcon from '@material-ui/icons/ListAlt';
import React from 'react';
import ToggleButton from '@material-ui/lab/ToggleButton';
import ToggleButtonGroup from '@material-ui/lab/ToggleButtonGroup';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(_theme => ({
  root: {
    boxShadow: 'none',
    borderRadius: '4px',
    padding: '6px 0px',
  },
  button: {
    boxShadow: 'none',
    padding: '0px 10px',
    minWidth: 'auto',
    color: '#64748C',
    '&:first-child': {
      borderRight: '0.5px solid #D2DAE7',
    },
    '&:last-child': {
      borderLeft: '0.5px solid #D2DAE7',
    },
    '&:hover': {
      background: 'white',
    },
  },
  selectedButton: {
    '&&': {
      background: 'white',
      color: '#3984FF',
      '&:hover': {
        background: 'white',
      },
    },
  },
}));

type Props = {
  selectedView: 'details' | 'graph',
  onViewSelected: (view: 'details' | 'graph') => void,
};

const EntViewToggleButton = (props: Props) => {
  const {selectedView, onViewSelected} = props;
  const classes = useStyles();
  return (
    <ToggleButtonGroup
      className={classes.root}
      size="small"
      value={selectedView}
      exclusive
      onChange={(_e, view) => onViewSelected(view)}>
      <ToggleButton
        disableRipple
        className={classes.button}
        classes={{selected: classes.selectedButton}}
        value="details">
        <ListAltIcon />
      </ToggleButton>
      <ToggleButton
        disableRipple
        className={classes.button}
        classes={{selected: classes.selectedButton}}
        value="graph">
        <BubbleChartIcon />
      </ToggleButton>
    </ToggleButtonGroup>
  );
};

export default EntViewToggleButton;
