/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {Location} from '../../common/Location.js';

import * as React from 'react';
import AppContext from '@fbcnms/ui/context/AppContext';
import Card from '@fbcnms/ui/components/design-system/Card/Card';
import CardHeader from '@fbcnms/ui/components/design-system/Card/CardHeader';
import Grid from '@material-ui/core/Grid';
import LocationDetailsCardProperty from './LocationDetailsCardProperty';
import LocationMapSnippet from './LocationMapSnippet';
import Text from '@fbcnms/ui/components/design-system/Text';
import Tooltip from '@material-ui/core/Tooltip';
import fbt from 'fbt';
import {makeStyles} from '@material-ui/styles';
import {useContext} from 'react';

const useStyles = makeStyles(() => ({
  header: {
    flexGrow: 1,
  },
  map: {
    minHeight: '232px',
  },
}));

type Props = {
  className?: string,
  location: Location,
};

const LocationDetailsCard = (props: Props) => {
  const {className, location} = props;
  const classes = useStyles();
  const externalIDEnabled = useContext(AppContext).isFeatureEnabled(
    'external_id',
  );

  const getCoordTitle = (title: string) => {
    return (
      <Tooltip
        arrow
        placement="left"
        title={fbt('Value taken from parent', '')}>
        <div>
          <Text>
            <fbt desc="">
              <fbt:param name="Latitude or longitude followed by a star">
                {title}
              </fbt:param>
            </fbt>
          </Text>
        </div>
      </Tooltip>
    );
  };

  const getCoordAndTitle = (coord: 'LAT' | 'LONG', location: Location) => {
    if (
      location.latitude === 0 &&
      location.longitude === 0 &&
      location.parentCoords !== null
    ) {
      const value =
        coord === 'LAT'
          ? location.parentCoords?.latitude ?? 0
          : location.parentCoords?.longitude ?? 0;
      const title =
        coord === 'LAT' ? getCoordTitle('Lat*') : getCoordTitle('Long*');
      return {title, value};
    } else {
      const value =
        coord === 'LAT' ? location.latitude ?? 0 : location.longitude ?? 0;
      const title = coord == 'LAT' ? fbt('Lat', '') : fbt('Long', '');
      return {title, value};
    }
  };

  const latDetails = getCoordAndTitle('LAT', location);
  const longDetails = getCoordAndTitle('LONG', location);

  return (
    <Card className={className}>
      <Grid container spacing={2}>
        <Grid item xs={12} md={4}>
          <CardHeader className={classes.header}>Details</CardHeader>
          <Grid container>
            <LocationDetailsCardProperty
              title={
                <Text>
                  <fbt desc="">Type</fbt>
                </Text>
              }
              value={location.locationType.name}
            />
            {externalIDEnabled && location.externalId && (
              <LocationDetailsCardProperty
                title={
                  <Text>
                    <fbt desc="">External ID</fbt>
                  </Text>
                }
                value={location.externalId}
              />
            )}
            {latDetails.value !== 0 && (
              <LocationDetailsCardProperty
                title={latDetails.title}
                value={String(latDetails.value)}
              />
            )}
            {longDetails.value !== 0 && (
              <LocationDetailsCardProperty
                title={longDetails.title}
                value={String(longDetails.value)}
              />
            )}
          </Grid>
        </Grid>
        <Grid item xs={12} md={8}>
          <LocationMapSnippet
            className={classes.map}
            location={{
              id: location.id,
              name: location.name,
              latitude: location.latitude,
              longitude: location.longitude,
              locationType: {
                mapType: location.locationType.mapType,
                mapZoomLevel: location.locationType.mapZoomLevel,
              },
            }}
          />
        </Grid>
      </Grid>
    </Card>
  );
};

export default LocationDetailsCard;
