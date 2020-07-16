/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import Grid from '@material-ui/core/Grid';
import HelpIcon from '@material-ui/icons/Help';
import IconButton from '@material-ui/core/IconButton';
import InputAdornment from '@material-ui/core/InputAdornment';
import ListItem from '@material-ui/core/ListItem';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import Tooltip from '@material-ui/core/Tooltip';
import Typography from '@material-ui/core/Typography';
import Visibility from '@material-ui/icons/Visibility';
import VisibilityOff from '@material-ui/icons/VisibilityOff';

import {makeStyles} from '@material-ui/styles';
import {useState} from 'react';

const useStyles = makeStyles(_ => ({
  root: {
    display: 'flex',
    marginBottom: '5px',
    alignItems: 'center',
  },
  heading: {
    flexBasis: '33.33%',
    marginRight: '15px',
    textAlign: 'right',
  },
  secondaryHeading: {
    flexBasis: '66.66%',
  },
  icon: {
    marginLeft: '5px',
    paddingTop: '4px',
    verticalAlign: 'bottom',
    width: '15px',
  },
}));

type Props = {
  label: string,
  children?: any,
  tooltip?: string,
};

export default function FormField(props: Props) {
  const classes = useStyles();
  const {tooltip} = props;
  return (
    <div className={classes.root}>
      <Text className={classes.heading} variant="body2">
        {props.label}
        {tooltip && (
          <Tooltip title={tooltip} placement="bottom-start">
            <HelpIcon className={classes.icon} />
          </Tooltip>
        )}
      </Text>
      <Typography
        className={classes.secondaryHeading}
        component="div"
        variant="body2">
        {props.children}
      </Typography>
    </div>
  );
}

export function AltFormField(props: Props) {
  return (
    <ListItem>
      <Grid container>
        <Grid item xs={12}>
          {props.label}
        </Grid>
        <Grid item xs={12}>
          {props.children}
        </Grid>
      </Grid>
    </ListItem>
  );
}

type PasswordProps = {
  value: string,
  onChange: string => void,
};

export function PasswordInput(props: PasswordProps) {
  const [showPassword, setShowPassword] = useState(false);
  return (
    <OutlinedInput
      {...props}
      type={showPassword ? 'text' : 'password'}
      value={props.value}
      onChange={e => props.onChange(e.target.value)}
      endAdornment={
        <InputAdornment position="end">
          <IconButton
            aria-label="toggle password visibility"
            onClick={() => setShowPassword(true)}
            onMouseDown={() => setShowPassword(false)}
            edge="end">
            {showPassword ? <Visibility /> : <VisibilityOff />}
          </IconButton>
        </InputAdornment>
      }
    />
  );
}
