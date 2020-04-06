/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {LineMapLayer} from './WifiMapLayers';

import FormControl from '@material-ui/core/FormControl';
import FormHelperText from '@material-ui/core/FormHelperText';
import Input from '@material-ui/core/Input';
import InputLabel from '@material-ui/core/InputLabel';
import React from 'react';
import TypedSelect from '@fbcnms/ui/components/TypedSelect';

import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  formControl: {
    margin: theme.spacing(),
    minWidth: 120,
    width: 'calc(100% - 15px)',
  },
  selectEmpty: {
    marginTop: theme.spacing(2),
  },
}));

type Props = {
  onChange: (connType: LineMapLayer | '') => void,
  selectedConnType: LineMapLayer | '',
};

export default function WifiSelectConnType(props: Props) {
  const classes = useStyles();
  return (
    <FormControl className={classes.formControl}>
      <InputLabel htmlFor="conntype-helper">Connection Filter</InputLabel>
      <TypedSelect
        items={{
          '': 'All',
          defaultRoute: 'Default Routes',
          l3: 'L3 only',
          l2: 'L2 only',
          none: 'Visible (low signal)',
        }}
        value={props.selectedConnType}
        onChange={props.onChange}
        input={<Input name="connType" id="conntype-helper" />}
      />
      <FormHelperText>Filter by Connection Type</FormHelperText>
    </FormControl>
  );
}
