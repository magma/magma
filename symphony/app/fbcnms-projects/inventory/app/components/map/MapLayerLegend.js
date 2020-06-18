/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import Checkbox from '@material-ui/core/Checkbox';
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import update from 'immutability-helper';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  legend: {
    backgroundColor: theme.palette.grey[100],
    width: '100%',
    height: '100%',
    boxShadow: '1px 0px 0px 0px rgba(0, 0, 0, 0.1)',
    padding: '24px 16px',
    position: 'relative',
    zIndex: 2,
    overflowY: 'auto',
  },
  typeItem: {
    display: 'flex',
    alignItems: 'center',
    marginBottom: '12px',
  },
  typeItemName: {
    fontSize: '13px',
    flexGrow: 1,
  },
  checkboxRoot: {
    padding: '0px 9px 0px 0px',
  },
  title: {
    display: 'block',
    lineHeight: '19px',
    fontSize: '16px',
    fontWeight: 'bold',
    marginBottom: '20px',
  },
  marker: {
    height: '12px',
    width: '12px',
    borderRadius: '100%',
  },
}));

export type MapLayer = {
  id: string,
  name: string,
  color: string,
};

type Props = {
  layers: Array<MapLayer>,
  selection: Array<string>,
  onSelectionChanged: (selection: Array<string>) => void,
};

const MapLayerLegend = (props: Props) => {
  const classes = useStyles();
  const {layers, onSelectionChanged, selection} = props;
  return (
    <div className={classes.legend}>
      <Text variant="h6" className={classes.title}>
        Layers
      </Text>
      {layers.map(layer => (
        <div key={layer.id} className={classes.typeItem}>
          <Checkbox
            classes={{root: classes.checkboxRoot}}
            checked={selection.includes(layer.id)}
            disableRipple
            onChange={(_e, checked) => {
              const newSelection = update(
                selection,
                checked
                  ? {
                      $push: [layer.id],
                    }
                  : {
                      $splice: [[selection.indexOf(layer.id), 1]],
                    },
              );
              onSelectionChanged(newSelection);
            }}
            value={layer.id}
            color="primary"
          />
          <Text variant="body2" className={classes.typeItemName}>
            {layer.name}
          </Text>
          <div
            className={classes.marker}
            style={{
              backgroundColor: layer.color,
            }}
          />
        </div>
      ))}
    </div>
  );
};

export default MapLayerLegend;
