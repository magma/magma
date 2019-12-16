/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import React from 'react';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(theme => ({
  tableCell: {
    fontSize: '12px',
    lineHeight: '16px',
    color: theme.palette.blueGrayDark,
    fontWeight: 400,
    borderBottom: 'none',
  },
  link: {
    color: '#3984ff',
    fontSize: '12px',
    lineHeight: '16px',
    fontWeight: 500,
    textDecoration: 'none',
  },
}));

type Props = {
  field: {
    type: string,
    value: string,
  },
};

const FieldValue = (props: Props) => {
  const {field} = props;
  const {history} = useRouter();
  const classes = useStyles();

  switch (field.type) {
    case 'ID':
      return (
        <span onClick={() => history.push(`/id/${field.value}`)}>
          <a className={classes.link} href="#">
            {field.value}
          </a>
        </span>
      );
    default:
      return <span>{field.value}</span>;
  }
};

export default FieldValue;
