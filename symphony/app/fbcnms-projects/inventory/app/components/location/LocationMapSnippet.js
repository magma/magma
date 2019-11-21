/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import MapView from '../map/MapView';
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import classNames from 'classnames';
import {locationToGeoJson} from '../map/MapUtil';
import {makeStyles} from '@material-ui/styles';

type Props = {
  className?: string,
  location: {
    id: string,
    name: string,
    latitude: number,
    longitude: number,
    locationType?: {mapType?: ?string, mapZoomLevel?: ?string},
  },
};

const useStyles = makeStyles(_theme => ({
  container: {
    position: 'relative',
    height: '100%',
  },
  mapContainer: {
    height: '100%',
    marginBottom: '5px',
  },
  section: {
    marginBottom: '24px',
  },
  latLngLabel: {
    backgroundColor: 'rgba(255, 255, 255, 0.5)',
    display: 'inline',
    fontSize: '12px',
    lineHeight: '20px',
    fontFamily: `'Helvetica Neue', Arial, Helvetica, sans-serif`,
    position: 'absolute',
    bottom: '4px',
    left: '5px',
    padding: '2px 8px',
  },
}));

const LocationMapSnippet = (props: Props) => {
  const {location, className} = props;
  const classes = useStyles();

  const hasGeoLocation =
    location.latitude !== null &&
    location.longitude !== null &&
    (location.latitude >= -90 && location.latitude <= 90) &&
    (location.longitude >= -180 && location.longitude <= 180) &&
    (location.latitude !== 0 || location.longitude !== 0);

  if (!hasGeoLocation) {
    return null;
  }
  return (
    <div className={classNames(classes.container, className)}>
      <MapView
        id="mapView"
        center={{
          lat: location.latitude,
          lng: location.longitude,
        }}
        mode={
          location.locationType?.mapType === 'satellite'
            ? 'satellite'
            : 'streets'
        }
        zoomLevel={location.locationType?.mapZoomLevel ?? '8'}
        classes={{mapContainer: classes.mapContainer}}
        markers={locationToGeoJson(location)}
        showGeocoder={false}
        showMapSatelliteToggle={false}
      />
      <Text variant="body2" className={classes.latLngLabel}>
        {location.latitude.toFixed(3)} / {location.longitude.toFixed(3)}
      </Text>
    </div>
  );
};

export default LocationMapSnippet;
