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

import ListItemText from '@material-ui/core/ListItemText';
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
        <TextField
          select
          label="Map Type"
          variant="outlined"
          className={classes.input}
          value={this.props.mapType}
          onChange={this.props.onMapTypeChanged}
          SelectProps={{
            MenuProps: {
              className: classes.menu,
            },
          }}
          margin="normal">
          {Object.keys(mapViewTypes).map(type => (
            <MenuItem key={type} value={type}>
              <ListItemText
                classes={{primary: classes.primary}}
                primary={mapViewTypes[type]}
              />
            </MenuItem>
          ))}
        </TextField>
        <TextField
          select
          label="Map Zoom Level"
          variant="outlined"
          className={classes.input}
          value={this.props.mapZoomLevel}
          onChange={this.props.onMapZoomLevelChanged}
          MenuProps={{
            className: classes.menu,
          }}
          margin="normal">
          {Object.keys(mapZoomLevels).map(mapZoomLevel => {
            levelsCount++;
            return (
              <MenuItem key={mapZoomLevel} value={parseInt(mapZoomLevel)}>
                <ListItemText
                  classes={{primary: classes.primary}}
                  primary={
                    levelsCount +
                    (mapZoomLevels[mapZoomLevel]
                      ? ' - ' + mapZoomLevels[mapZoomLevel]
                      : '')
                  }
                />
              </MenuItem>
            );
          })}
        </TextField>
      </div>
    );
  }
}

export default withStyles(styles)(LocationMapViewProperties);
