/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import type {WithStyles} from '@material-ui/core';

import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import TextField from '@material-ui/core/TextField';
import {withStyles} from '@material-ui/core/styles';

type Props = {
  mapType: string,
  mapZoomLevel: string,
  onMapTypeChanged: (event: SyntheticInputEvent<any>) => void,
  onMapZoomLevelChanged: (mapZoomLevel: SyntheticInputEvent<any>) => void,
} & WithStyles<typeof styles>;

const styles = theme => ({
  container: {
    display: 'flex',
  },
  input: {
    marginRight: theme.spacing(2),
    width: 305,
  },
  menu: {
    width: 305,
  },
  selectMenu: {
    height: '14px',
  },
});

const mapViewTypes = {
  map: 'Default',
  satellite: 'Satellite',
};

const mapZoomLevels = {
  '7': 'Country',
  '8': 'Wide area',
  '11': 'City',
  '13': 'Village or town',
  '16': 'Street',
  '18': 'Building',
};

class LocationMapViewProperties extends React.PureComponent<Props> {
  render() {
    const {classes} = this.props;
    let levelsCount = 0;
    return (
      <div className={classes.container}>
        <FormField label="Map Type" className={classes.input}>
          <TextField
            select
            variant="outlined"
            value={this.props.mapType}
            onChange={this.props.onMapTypeChanged}
            SelectProps={{
              classes: {selectMenu: classes.selectMenu},
              MenuProps: {
                className: classes.menu,
              },
            }}
            margin="dense">
            {Object.keys(mapViewTypes).map(type => (
              <MenuItem key={type} value={type}>
                {mapViewTypes[type]}
              </MenuItem>
            ))}
          </TextField>
        </FormField>
        <FormField label="Map Zoom Level" className={classes.input}>
          <TextField
            select
            variant="outlined"
            className={classes.input}
            value={this.props.mapZoomLevel}
            onChange={this.props.onMapZoomLevelChanged}
            SelectProps={{
              classes: {selectMenu: classes.selectMenu},
              MenuProps: {
                className: classes.menu,
              },
            }}
            margin="dense">
            {Object.keys(mapZoomLevels).map(mapZoomLevel => {
              levelsCount++;
              return (
                <MenuItem key={mapZoomLevel} value={parseInt(mapZoomLevel)}>
                  {levelsCount +
                    (mapZoomLevels[mapZoomLevel]
                      ? ' - ' + mapZoomLevels[mapZoomLevel]
                      : '')}
                </MenuItem>
              );
            })}
          </TextField>
        </FormField>
      </div>
    );
  }
}

export default withStyles(styles)(LocationMapViewProperties);
