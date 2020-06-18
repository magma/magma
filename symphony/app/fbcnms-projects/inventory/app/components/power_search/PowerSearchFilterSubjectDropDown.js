/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import type {EntityType} from '../comparison_view/ComparisonViewTypes';

import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import TextField from '@material-ui/core/TextField';
import {EntityTypeMap} from '../comparison_view/ComparisonViewTypes';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  selectMenu: {
    height: '14px',
    margin: '3px',
    marginRight: '5px',
  },
  beforeDropdown: {
    color: theme.palette.grey.A200,
    fontWeight: 'bold',
    paddingRight: '10px',
  },
  dropDownContainer: {
    display: 'flex',
    alignItems: 'center',
    paddingRight: '15px',
    paddingLeft: '15px',
  },
  input: {
    marginTop: '0px',
    marginBottom: '0px',
  },
}));

type Props = {
  subject: EntityType,
  onSubjectChange: EntityType => void,
};

const FILTER_SUBJECTS = {
  equipment: 'Equipment',
  link: 'Links',
  port: 'Ports',
  location: 'Locations',
};

const PowerSearchFilterSubjectDropDown = (props: Props) => {
  const classes = useStyles();
  const {subject, onSubjectChange} = props;

  return (
    <div className={classes.dropDownContainer}>
      <Text variant="body2" className={classes.beforeDropdown}>
        {'Filtering'}
      </Text>
      <TextField
        select
        variant="outlined"
        value={subject}
        onChange={event => {
          const v = event.target.value;
          if (Object.keys(EntityTypeMap).indexOf(v) !== -1) {
            // $FlowFixMe
            onSubjectChange(v);
          }
        }}
        className={classes.input}
        SelectProps={{
          classes: {selectMenu: classes.selectMenu},
        }}
        margin="dense">
        {Object.keys(FILTER_SUBJECTS).map(type => (
          <MenuItem key={type} value={type}>
            {FILTER_SUBJECTS[type]}
          </MenuItem>
        ))}
      </TextField>
    </div>
  );
};

export default PowerSearchFilterSubjectDropDown;
