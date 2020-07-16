/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import * as React from 'react';
import Card from '@material-ui/core/Card';
import classNames from 'classnames';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(_theme => ({
  card: {
    margin: '16px 0px',
    padding: '24px',
    boxShadow: '0px 1px 4px 0px rgba(0,0,0,0.17)',
  },
}));

type Props = {
  className?: string,
  children: Array<React.Element<*>> | React.Element<*>,
};

const SectionedCard = (props: Props) => {
  const classes = useStyles();
  return (
    <Card className={classNames(props.className, classes.card)}>
      {props.children}
    </Card>
  );
};

export default SectionedCard;
