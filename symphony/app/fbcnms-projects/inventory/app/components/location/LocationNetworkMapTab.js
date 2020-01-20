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
import Card from '@fbcnms/ui/components/design-system/Card/Card';
import InventoryQueryRenderer from '../../components/InventoryQueryRenderer';
import LocationEquipmentTopology from './LocationEquipmentTopology';
import {graphql} from 'react-relay';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(_theme => ({
  root: {
    height: '100%',
  },
}));

type Props = {
  locationId: string,
};

const networkTopologyQuery = graphql`
  query LocationNetworkMapTabQuery($locationId: ID!) {
    location: node(id: $locationId) {
      ... on Location {
        equipments {
          ...LocationEquipmentTopology_equipment
        }
        topology {
          ...LocationEquipmentTopology_topology
        }
      }
    }
  }
`;

const LocationNetworkMapTab = (props: Props) => {
  const {locationId} = props;
  const classes = useStyles();
  return (
    <InventoryQueryRenderer
      query={networkTopologyQuery}
      variables={{
        locationId: locationId,
      }}
      render={props => {
        const location = props.location;
        return (
          <Card className={classes.root}>
            <LocationEquipmentTopology
              topology={location.topology}
              equipment={location.equipments}
            />
          </Card>
        );
      }}
    />
  );
};

export default LocationNetworkMapTab;
