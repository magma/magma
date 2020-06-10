/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {Equipment} from '../common/Equipment';
import type {LocationMenu_location} from '../components/location/__generated__/LocationMenu_location.graphql';
import type {LocationType} from '../common/LocationType';

import LocationAddEditCard from './location/LocationAddEditCard';
import LocationPropertiesCard from './location/LocationPropertiesCard';
import React from 'react';

type Props = {
  mode: 'show' | 'add' | 'edit',
  onEdit: () => void,
  onSave: (locationId: string) => void,
  onCancel: () => void,
  parentLocationId: ?string,
  selectedLocationId: ?string,
  selectedLocationType: ?LocationType,
  selectedWorkOrderId: ?string,
  onEquipmentSelected: Equipment => void,
  onWorkOrderSelected: (?string) => void,
  onAddEquipment: () => void,
  onLocationMoved: (movedLocation: LocationMenu_location) => void,
  onLocationRemoved: (removedLocation: LocationMenu_location) => void,
};

class LocationCard extends React.Component<Props> {
  render() {
    switch (this.props.mode) {
      case 'add':
        return (
          <LocationAddEditCard
            key={'new_location_type'}
            parentId={this.props.parentLocationId}
            type={this.props.selectedLocationType}
            onCancel={this.props.onCancel}
            onSave={this.props.onSave}
          />
        );
      case 'edit':
        return (
          <LocationAddEditCard
            key={'new_location_type'}
            editingLocationId={this.props.selectedLocationId}
            parentId={this.props.parentLocationId}
            type={this.props.selectedLocationType}
            onCancel={this.props.onCancel}
            onSave={this.props.onSave}
          />
        );
      case 'show':
        return (
          <LocationPropertiesCard
            key={this.props.selectedLocationId}
            locationId={this.props.selectedLocationId}
            selectedWorkOrderId={this.props.selectedWorkOrderId}
            onEquipmentSelected={this.props.onEquipmentSelected}
            onWorkOrderSelected={this.props.onWorkOrderSelected}
            onEdit={this.props.onEdit}
            onAddEquipment={this.props.onAddEquipment}
            onLocationMoved={this.props.onLocationMoved}
            onLocationRemoved={this.props.onLocationRemoved}
          />
        );
      default:
        return null;
    }
  }
}

export default LocationCard;
