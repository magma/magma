/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import * as colors from '@fbcnms/ui/theme/colors';
import CircularProgress from '@material-ui/core/CircularProgress';
import Grid from '@material-ui/core/Grid';
import InventoryErrorBoundary from '../../common/InventoryErrorBoundary';
import LocationPopout from './LocationPopout';
import MapLayerLegend from './MapLayerLegend';
import MapView from './MapView';
import React, {useEffect, useState} from 'react';
import RelayEnvironment from '../../common/RelayEnvironment';
import {fetchQuery, graphql} from 'relay-runtime';
import {locationsToGeoJSONSource} from './MapUtil';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  root: {
    height: '100%',
  },
  legend: {
    backgroundColor: theme.palette.grey[100],
    width: '100%',
    height: '100%',
    boxShadow: '1px 0px 0px 0px rgba(0, 0, 0, 0.1)',
  },
  loadingContainer: {
    height: '100%',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
  },
  legendContainer: {
    height: '100%',
  },
}));

const locationTypesQuery = graphql`
  query LocationsMapTypesQuery {
    locationTypes {
      edges {
        node {
          id
          name
          locations(enforceHasLatLong: true) {
            edges {
              node {
                id
                name
                latitude
                longitude
              }
            }
          }
        }
      }
    }
  }
`;

const colorPalette: Array<string> = [
  colors.cherry,
  colors.cyan,
  colors.orange,
  colors.pink,
  colors.red,
  colors.green,
  colors.darkPink,
  colors.darkPurple,
  colors.yellow,
  colors.darkGreen,
  colors.gray3,
  colors.lightBlue,
  colors.brown,
];

type Props = {};

const LocationsMap = (_props: Props) => {
  const classes = useStyles();
  const [locationTypes, setLocationsTypes] = useState([]);
  const [selectedTypeIds, setSelectedTypeIds] = useState([]);
  const [error, setError] = useState(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    setIsLoading(true);
    // $FlowFixMe (T62907961) Relay flow types
    fetchQuery(RelayEnvironment, locationTypesQuery)
      .then(
        data => {
          const locationTypes = data.locationTypes.edges
            .map(edge => edge.node)
            .filter(type => type.locations.edges.length > 0);
          setLocationsTypes(locationTypes);
          setSelectedTypeIds(locationTypes.map(l => l.id));
          setError(null);
        },
        _e => {
          setError('Failed fetching location types');
        },
      )
      .finally(() => {
        setIsLoading(false);
      });
  }, []);

  const layers = locationTypes.map(locationType => {
    const locations = locationType.locations.edges.map(l => l.node);
    const typeIndex = locationTypes.findIndex(
      type => type.id === locationType.id,
    );
    return {
      source: locationsToGeoJSONSource(locationType.id, locations, {
        primaryKey: locationType.id,
        color: colorPalette[typeIndex % colorPalette.length],
        name: locationType.id,
      }),
    };
  });
  if (error) {
    return <div>{error}</div>;
  }

  if (isLoading) {
    return (
      <div className={classes.loadingContainer}>
        <CircularProgress size={50} />
      </div>
    );
  }

  return (
    <Grid className={classes.root} container spacing={0}>
      <InventoryErrorBoundary>
        <Grid className={classes.legendContainer} item xs={2}>
          <MapLayerLegend
            layers={locationTypes.map((type, i) => ({
              id: type.id,
              name: type.name,
              color: colorPalette[i % colorPalette.length],
            }))}
            selection={selectedTypeIds}
            onSelectionChanged={setSelectedTypeIds}
          />
        </Grid>
        <Grid item xs={10}>
          <MapView
            mode="streets"
            layers={layers.filter(l => selectedTypeIds.includes(l.source.key))}
            getFeaturePopoutContent={feature => (
              <LocationPopout locationId={feature.properties.id} />
            )}
            showGeocoder={true}
            showMapSatelliteToggle={true}
          />
        </Grid>
      </InventoryErrorBoundary>
    </Grid>
  );
};

export default LocationsMap;
