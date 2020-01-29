/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import * as React from 'react';
import DynamicPropertiesGrid from '../DynamicPropertiesGrid';
import InventoryQueryRenderer from '../InventoryQueryRenderer';
import LocationBreadcrumbsTitle from '../location/LocationBreadcrumbsTitle';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {graphql} from 'relay-runtime';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(_theme => ({
  root: {
    marginTop: '8px',
    maxWidth: '600px',
    minWidth: '400px',
  },
  title: {
    marginBottom: '12px',
  },
}));

type Props = {
  locationId: string,
};

const locationPopoutQuery = graphql`
  query LocationPopoutQuery($locationId: ID!) {
    location: node(id: $locationId) {
      ... on Location {
        id
        name
        locationType {
          name
        }
        ...LocationBreadcrumbsTitle_locationDetails
        properties {
          ...DynamicPropertiesGrid_properties
        }
        locationType {
          propertyTypes {
            ...DynamicPropertiesGrid_propertyTypes
          }
        }
      }
    }
  }
`;

const LocationPopout = (props: Props) => {
  const {locationId} = props;
  const classes = useStyles();
  React.useEffect(() => {
    ServerLogger.info(LogEvents.LOCATIONS_MAP_POPUP_OPENED, {locationId});
  }, [locationId]);
  return (
    <InventoryQueryRenderer
      query={locationPopoutQuery}
      variables={{locationId}}
      render={props => {
        const {location} = props;
        return (
          <div className={classes.root}>
            <div className={classes.title}>
              <LocationBreadcrumbsTitle
                locationDetails={location}
                hideTypes={true}
                size="small"
              />
            </div>
            <DynamicPropertiesGrid
              properties={location.properties}
              propertyTypes={location.locationType.propertyTypes}
              hideTitle={true}
            />
          </div>
        );
      }}
    />
  );
};

export default LocationPopout;
