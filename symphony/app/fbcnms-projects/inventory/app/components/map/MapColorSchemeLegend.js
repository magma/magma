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
import Text from '@fbcnms/ui/components/design-system/Text';

import {fade} from '@material-ui/core/styles/colorManipulator';
import {makeStyles} from '@material-ui/styles';

export type MapColorScheme = {
  lowLabel: string,
  highLabel: string,
  colors: Array<string>,
};

type Props = {
  title: string,
  colorScheme: MapColorScheme,
};

const MapColorSchemeLegend = (props: Props) => {
  const classes = useStyles();
  const {title, colorScheme} = props;
  return (
    <div className={classes.container}>
      <Text variant="subtitle2" className={classes.title}>
        {title}
      </Text>
      <div className={classes.colorScale}>
        {colorScheme.colors.map(color => (
          <div
            key={color}
            style={{backgroundColor: color}}
            className={classes.colorUnit}
          />
        ))}
      </div>
      <div className={classes.labelContainer}>
        <Text variant="body2">{colorScheme.lowLabel}</Text>
        <Text variant="body2">{colorScheme.highLabel}</Text>
      </div>
    </div>
  );
};

const useStyles = makeStyles(theme => ({
  container: {
    backgroundColor: fade(theme.palette.background.paper, 0.5),
    padding: '10px',
  },
  title: {
    marginBottom: '10px',
  },
  colorScale: {
    height: '20px',
    width: '200px',
    display: 'flex',
    flexDirection: 'row',
    marginBottom: '5px',
  },
  colorUnit: {
    flex: 1,
  },
  labelContainer: {
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'space-between',
  },
}));

export default MapColorSchemeLegend;
