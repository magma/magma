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
import React from 'react';
import Select from '@fbcnms/ui/components/design-system/Select/Select';
import {withStyles} from '@material-ui/core/styles';

type Props = {
  mapType: string,
  mapZoomLevel: string,
  onMapTypeChanged: (newType: string) => void,
  onMapZoomLevelChanged: (newZoomLevel: string) => void,
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
          <Select
            options={Object.keys(mapViewTypes).map(type => ({
              key: type,
              value: type,
              label: mapViewTypes[type],
            }))}
            selectedValue={this.props.mapType}
            onChange={this.props.onMapTypeChanged}
          />
        </FormField>
        <FormField label="Map Zoom Level" className={classes.input}>
          <Select
            options={Object.keys(mapZoomLevels).map(level => ({
              key: level,
              value: level,
              label:
                ++levelsCount +
                (mapZoomLevels[level] ? ' - ' + mapZoomLevels[level] : ''),
            }))}
            selectedValue={this.props.mapZoomLevel}
            onChange={this.props.onMapZoomLevelChanged}
          />
        </FormField>
      </div>
    );
  }
}

export default withStyles(styles)(LocationMapViewProperties);
