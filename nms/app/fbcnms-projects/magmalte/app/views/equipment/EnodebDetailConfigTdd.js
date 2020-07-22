/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import Collapse from '@material-ui/core/Collapse';
import ExpandLess from '@material-ui/icons/ExpandLess';
import ExpandMore from '@material-ui/icons/ExpandMore';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import React from 'react';

import {AltFormField} from '../../components/FormField';
import {colors} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(_ => ({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '80%',
  },
  itemTitle: {
    color: colors.primary.comet,
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
  },
  itemValue: {
    color: colors.primary.brightGray,
  },
}));

type Props = {
  earfcndl: number,
  specialSubframePattern: number,
  subframeAssignment: number,
};
export function EnodeConfigTdd(props: Props) {
  const classes = useStyles();
  const [open, setOpen] = React.useState(true);
  const typographyProps = {
    primaryTypographyProps: {
      variant: 'caption',
      className: classes.itemTitle,
    },
    secondaryTypographyProps: {
      variant: 'h6',
      className: classes.itemValue,
    },
  };
  return (
    <List key="tddConfigs">
      <ListItem button onClick={() => setOpen(!open)}>
        <ListItemText primary="TDD" />
        {open ? <ExpandLess /> : <ExpandMore />}
      </ListItem>
      <Collapse key="tdd" in={open} timeout="auto" unmountOnExit>
        <ListItem>
          <ListItemText
            secondary={props.earfcndl}
            primary="EARFCNDL"
            {...typographyProps}
          />
        </ListItem>
        <ListItem>
          <ListItemText
            secondary={props.specialSubframePattern}
            primary="Special Subframe Pattern"
            {...typographyProps}
          />
        </ListItem>
        <ListItem>
          <ListItemText
            secondary={props.subframeAssignment}
            primary="Subframe Assignment"
            {...typographyProps}
          />
        </ListItem>
      </Collapse>
    </List>
  );
}

type EditProps = {
  earfcndl: string,
  specialSubframePattern: string,
  subframeAssignment: string,
  setEarfcndl: string => void,
  setSpecialSubframePattern: string => void,
  setSubframeAssignment: string => void,
};

export default function EnodeConfigEditTdd(props: EditProps) {
  const classes = useStyles();

  return (
    <>
      <AltFormField label={'EARFCNDL'}>
        <OutlinedInput
          data-testid="earfcndl"
          className={classes.input}
          fullWidth={true}
          value={props.earfcndl}
          onChange={({target}) => props.setEarfcndl(target.value)}
        />
      </AltFormField>
      <AltFormField label={'Special Subframe Pattern'}>
        <OutlinedInput
          className={classes.input}
          fullWidth={true}
          value={props.specialSubframePattern}
          onChange={({target}) => props.setSpecialSubframePattern(target.value)}
        />
      </AltFormField>
      <AltFormField label={'Subframe Assignment'}>
        <OutlinedInput
          className={classes.input}
          fullWidth={true}
          value={props.subframeAssignment}
          onChange={({target}) => props.setSubframeAssignment(target.value)}
        />
      </AltFormField>
    </>
  );
}
