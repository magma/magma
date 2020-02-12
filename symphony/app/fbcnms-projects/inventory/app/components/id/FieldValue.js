/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import React from 'react';
import {Link} from 'react-router-dom';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles({
  link: {
    color: '#3984ff',
    fontSize: '12px',
    lineHeight: '16px',
    fontWeight: 500,
    textDecoration: 'none',
  },
});

type Props = {
  field: {
    type: string,
    value: string,
  },
};

const FieldValue = (props: Props) => {
  const {field} = props;
  const classes = useStyles();

  switch (field.type) {
    case 'ID':
      return (
        <Link className={classes.link} to={`/id/${field.value}`}>
          {field.value}
        </Link>
      );
    default:
      return <span>{field.value}</span>;
  }
};

export default FieldValue;
