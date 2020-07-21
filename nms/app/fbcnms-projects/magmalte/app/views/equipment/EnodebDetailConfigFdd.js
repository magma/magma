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
import Divider from '@material-ui/core/Divider';
import ExpandLess from '@material-ui/icons/ExpandLess';
import ExpandMore from '@material-ui/icons/ExpandMore';
import Grid from '@material-ui/core/Grid';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import React from 'react';

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
  earfcnul: number,
};
export function EnodeConfigFdd(props: Props) {
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
    <List key="fddConfigs">
      <ListItem button onClick={() => setOpen(!open)}>
        <ListItemText primary="FDD" {...typographyProps} />
        {open ? <ExpandLess /> : <ExpandMore />}
      </ListItem>
      <Divider />
      <Collapse key="fdd" in={open} timeout="auto" unmountOnExit>
        <ListItem>
          <Grid container>
            <Grid item xs={6}>
              <ListItemText
                primary="EARFCNDL"
                secondary={props.earfcndl}
                {...typographyProps}
              />
            </Grid>
            <Grid item xs={6}>
              <ListItemText
                primary="EARFCNUL"
                secondary={props.earfcnul}
                {...typographyProps}
              />
              />
            </Grid>
          </Grid>
        </ListItem>
      </Collapse>
    </List>
  );
}

type EditProps = {
  earfcndl: string,
  earfcnul: string,
  setEarfcndl: string => void,
};
export default function EnodeConfigEditFdd(props: EditProps) {
  const classes = useStyles();

  return (
    <ListItem>
      <Grid container>
        <Grid item xs={6}>
          <Grid container>
            <Grid item xs={12}>
              EARFCNDL
            </Grid>
            <Grid item xs={12}>
              <OutlinedInput
                data-testid="earfcndl"
                className={classes.input}
                fullWidth={true}
                value={props.earfcndl}
                onChange={({target}) => props.setEarfcndl(target.value)}
              />
            </Grid>
          </Grid>
        </Grid>
        <Grid item xs={6}>
          <Grid container>
            <Grid item xs={12}>
              EARFCNUL
            </Grid>
            <Grid item xs={12}>
              <OutlinedInput
                className={classes.input}
                fullWidth={true}
                value={props.earfcnul}
                readOnly={true}
              />
            </Grid>
          </Grid>
        </Grid>
      </Grid>
    </ListItem>
  );
}
