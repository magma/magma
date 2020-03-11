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

import AppContext from '@fbcnms/ui/context/AppContext';
import Card from '@fbcnms/ui/components/design-system/Card/Card';
import CardHeader from '@fbcnms/ui/components/design-system/Card/CardHeader';
import Grid from '@material-ui/core/Grid';
import LocationDetailsCardProperty from './LocationDetailsCardProperty';
import LocationMapSnippet from './LocationMapSnippet';
import React, {useContext} from 'react';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  detailsHeaderContainer: {
    display: 'flex',
    flexDirection: 'row',
  },
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

  return (
    <Card className={className}>
      <div className={classes.detailsHeaderContainer}>
        <Grid container spacing={2}>
          <Grid item xs={12} md={4}>
            <CardHeader className={classes.header}>Details</CardHeader>
            <Grid container>
              <LocationDetailsCardProperty
                title="Type"
                value={location.locationType.name}
              />
              {externalIDEnabled && location.externalId && (
                <LocationDetailsCardProperty
                  title="External ID"
                  value={location.externalId}
                />
              )}
              {location.latitude !== 0 && (
                <LocationDetailsCardProperty
                  title="Lat"
                  value={String(location.latitude)}
                />
              )}
              {location.longitude !== 0 && (
                <LocationDetailsCardProperty
                  title="Long"
                  value={String(location.longitude)}
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
      </div>
    </Card>
  );
};

export default LocationDetailsCard;
