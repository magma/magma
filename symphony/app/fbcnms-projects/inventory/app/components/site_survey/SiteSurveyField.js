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
import Text from '@fbcnms/ui/components/design-system/Text';
import classNames from 'classnames';
import {gray7} from '@fbcnms/ui/theme/colors';
import {makeStyles} from '@material-ui/styles';

type Props = {
  label: string,
  children: any,
  className: ?string,
};

const useStyles = makeStyles(theme => ({
  root: {
    flexGrow: 1,
    display: 'flex',
    flexDirection: 'row',
  },
  labelName: {
    '&&': {
      color: theme.palette.blueGrayDark,
      fontSize: '14px',
      lineHeight: '24px',
    },
  },
  question: {
    padding: '0px 25px 16px 25px',
    flexBasis: '25%',
  },
  reply: {
    flexBasis: '75%',
    backgroundColor: gray7,
    padding: '0px 24px',
    display: 'grid',
  },
}));

export default function SiteSurveyField(props: Props) {
  const {label, children, className} = props;
  const classes = useStyles();
  return (
    <div className={classes.root}>
      <div className={classNames(classes.question, className)}>
        <Text className={classes.labelName}>{label}</Text>
      </div>
      <div className={classNames(classes.reply, className)}>{children}</div>
    </div>
  );
}
