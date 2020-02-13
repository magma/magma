/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import Grid from '@material-ui/core/Grid';
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import symphony from '@fbcnms/ui/theme/symphony';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles({
  propValue: {
    color: symphony.palette.D300,
  },
});

type Props = {
  title: string,
  value: string,
};

const LocationDetailsCardProperty = (props: Props) => {
  const {title, value} = props;
  const classes = useStyles();
  return (
    <>
      <Grid item xs={4}>
        <Text variant="subtitle2" weight="regular">
          {title}:
        </Text>
      </Grid>
      <Grid item xs={8}>
        <Text
          variant="subtitle2"
          weight="regular"
          className={classes.propValue}>
          {value}
        </Text>
      </Grid>
    </>
  );
};

export default LocationDetailsCardProperty;
