/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import type {Equipment} from '../../common/Equipment';

import Avatar from '@material-ui/core/Avatar';
import InventoryQueryRenderer from '../../components/InventoryQueryRenderer';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import React from 'react';
import RouterIcon from '@material-ui/icons/Router';

import {graphql} from 'react-relay';
import {sortLexicographically} from '@fbcnms/ui/utils/displayUtils';
import {useState} from 'react';

type Props = {
  locationId: string,
  onSelect: Equipment => void,
};

const locationEquipmentsQuery = graphql`
  query LocationEquipmentsQuery($locationId: ID!) {
    location: node(id: $locationId) {
      ... on Location {
        id
        equipments {
          id
          name
          equipmentType {
            id
            name
          }
        }
      }
    }
  }
`;

export default function(props: Props) {
  const [selectedEquipment, setSelectedEquipment] = useState('');
  return (
    <InventoryQueryRenderer
      query={locationEquipmentsQuery}
      variables={{locationId: props.locationId}}
      render={results => {
        const location = results.location;
        const listItems = location.equipments
          .slice()
          .sort((x, y) => sortLexicographically(x.name ?? '', y.name ?? ''))
          .map(equipment => (
            <ListItem
              dense
              button
              key={equipment.id}
              selected={selectedEquipment === equipment.id}
              onClick={() => {
                setSelectedEquipment(equipment.id);
                props.onSelect(equipment);
              }}>
              <Avatar>
                <RouterIcon />
              </Avatar>
              <ListItemText primary={equipment.name} />
            </ListItem>
          ));

        return <List>{listItems}</List>;
      }}
    />
  );
}
