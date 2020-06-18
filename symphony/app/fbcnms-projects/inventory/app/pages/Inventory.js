/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ContextRouter} from 'react-router-dom';
import type {EquipmentPosition} from '../common/Equipment';
import type {EquipmentType} from '../common/EquipmentType';
import type {Location} from '../common/Location';
import type {LocationMenu_location} from '../components/location/__generated__/LocationMenu_location.graphql';
import type {LocationType} from '../common/LocationType';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithSnackbarProps} from 'notistack';
import type {WithStyles} from '@material-ui/core';

import AddToLocationDialog from '../components/AddToLocationDialog';
import EquipmentCard from '../components/EquipmentCard';
import InventoryErrorBoundary from '../common/InventoryErrorBoundary';
import InventoryTopBar from '../components/InventoryTopBar';
import LocationCard from '../components/LocationCard';
import LocationsTree from '../components/LocationsTree';
import React from 'react';
import SnackbarItem from '@fbcnms/ui/components/SnackbarItem';
import fbt from 'fbt';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {InventoryAPIUrls} from '../common/InventoryAPI';
import {LogEvents, ServerLogger} from '../common/LoggingUtils';
import {extractEntityIdFromUrl} from '../common/RouterUtils';
import {withRouter} from 'react-router-dom';
import {withSnackbar} from 'notistack';
import {withStyles} from '@material-ui/core';

const styles = {
  header: {
    display: 'flex',
    justifyContent: 'space-between',
  },
  tree: {
    maxWidth: '',
  },
  addBtn: {
    position: 'absolute',
    bottom: '61px',
    right: '61px',
  },
  gridContainer: {
    display: 'flex',
    height: 'calc(100% - 60px)',
    flexWrap: 'nowrap',
  },
  propertiesCard: {
    padding: '24px',
    height: '100%',
    flexGrow: 1,
    minWidth: '75%',
  },
  tabsContainer: {
    padding: '20px',
  },
};

const ADD_LOCATION_CARD: Card = {mode: 'add', type: 'location'};
const ADD_EQUIPMENT_CARD: Card = {mode: 'add', type: 'equipment'};
const EDIT_LOCATION_CARD: Card = {mode: 'edit', type: 'location'};
const EDIT_EQUIPMENT_CARD: Card = {mode: 'edit', type: 'equipment'};
const SHOW_LOCATION_CARD: Card = {mode: 'show', type: 'location'};
const SHOW_EQUIPMENT_CARD: Card = {mode: 'show', type: 'equipment'};

type Card = {
  mode: 'add' | 'edit' | 'show',
  type: 'location' | 'equipment',
};

type Props = ContextRouter &
  WithStyles<typeof styles> &
  WithAlert &
  WithSnackbarProps & {};

type State = {
  card: Card,
  dialogMode: 'hidden' | 'location' | 'equipment',
  errorMessage: ?string,
  parentLocationId: ?string,
  selectedEquipmentId: ?string,
  selectedEquipmentPosition: ?EquipmentPosition,
  selectedEquipmentType: ?EquipmentType,
  selectedLocationId: ?string,
  selectedLocationType: ?LocationType,
  selectedWorkOrderId: ?string,
  openLocationHierarchy: Array<string>,
};

class Inventory extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);

    this.state = {
      card: SHOW_LOCATION_CARD,
      dialogMode: 'hidden',
      errorMessage: null,
      parentLocationId: null,
      selectedEquipmentId: null,
      selectedEquipmentPosition: null,
      selectedEquipmentType: null,
      selectedLocationId: null,
      selectedLocationType: null,
      selectedWorkOrderId: null,
      openLocationHierarchy: [],
    };
  }

  navigateToLocation(selectedLocationId: ?string, source: ?string) {
    ServerLogger.info(LogEvents.NAVIGATE_TO_LOCATION, {
      locationId: selectedLocationId,
      source,
    });

    if (selectedLocationId != null) {
      if (selectedLocationId != this.state.selectedLocationId) {
        this.props.history.push(InventoryAPIUrls.location(selectedLocationId));
      } else {
        this.setLocationCardState(selectedLocationId);
      }
    }
  }

  navigateToEquipment(selectedEquipmentId: ?string, source: ?string) {
    ServerLogger.info(LogEvents.NAVIGATE_TO_EQUIPMENT, {
      equipmentId: selectedEquipmentId,
      source,
    });
    if (selectedEquipmentId != null) {
      this.props.history.push(InventoryAPIUrls.equipment(selectedEquipmentId));
    }
  }

  navigateToWorkOrder(selectedWorkOrderCardId: ?string) {
    const {history} = this.props;
    if (selectedWorkOrderCardId) {
      history.push(`/workorders/search?workorder=${selectedWorkOrderCardId}`);
    }
  }

  setLocationCardState(locationId) {
    const {
      card,
      selectedLocationId,
      selectedEquipmentId,
      selectedEquipmentPosition,
    } = this.state;

    if (
      card === SHOW_LOCATION_CARD &&
      selectedLocationId === locationId &&
      selectedEquipmentId === null &&
      selectedEquipmentPosition === null
    ) {
      return;
    }

    this.setState({
      card: SHOW_LOCATION_CARD,
      selectedLocationType: null,
      selectedLocationId: locationId,
      selectedEquipmentId: null,
      selectedEquipmentPosition: null,
    });
  }

  render() {
    const {classes} = this.props;
    const {card} = this.state;

    const queryLocationId = extractEntityIdFromUrl(
      'location',
      this.props.location.search,
    );
    if (queryLocationId !== this.state.selectedLocationId) {
      this.setLocationCardState(queryLocationId);
    } else if (queryLocationId === null) {
      const queryEquipmentId = extractEntityIdFromUrl(
        'equipment',
        this.props.location.search,
      );
      if (queryEquipmentId !== this.state.selectedEquipmentId) {
        this.setState({
          card: SHOW_EQUIPMENT_CARD,
          selectedEquipmentId: queryEquipmentId,
          selectedLocationId: null,
          selectedEquipmentPosition: null,
        });
      } else if (
        (queryEquipmentId === null &&
          this.state.selectedLocationType === null) ||
        this.state.card === null
      ) {
        this.setLocationCardState(null);
      }
    }

    return (
      <>
        <InventoryTopBar
          onWorkOrderSelected={selectedWorkOrderId =>
            this.setState({selectedWorkOrderId})
          }
          onSearchEntitySelected={(entityId, entityType) => {
            switch (entityType) {
              case 'location':
                this.navigateToLocation(entityId, 'goto_search');
                break;
              case 'equipment':
                this.navigateToEquipment(entityId, 'goto_search');
                break;
            }
          }}
          onNavigateToWorkOrder={selectedWorkOrderCardId =>
            this.navigateToWorkOrder(selectedWorkOrderCardId)
          }
        />
        <div className={classes.gridContainer}>
          <LocationsTree
            selectedLocationId={this.state.selectedLocationId}
            onSelect={selectedLocationId =>
              this.navigateToLocation(selectedLocationId, 'tree')
            }
            onAddLocation={(parentLocation: ?Location) =>
              this.setState({parentLocationId: parentLocation?.id}, () =>
                this.showDialog('location'),
              )
            }
          />
          <div className={classes.propertiesCard}>
            <InventoryErrorBoundary>
              {card.type == 'location' && (
                <LocationCard
                  mode={card.mode}
                  onEdit={this.onLocationEdit}
                  onSave={this.onLocationSave}
                  onCancel={this.onLocationCancel}
                  parentLocationId={this.state.parentLocationId}
                  selectedLocationId={this.state.selectedLocationId}
                  selectedLocationType={this.state.selectedLocationType}
                  selectedWorkOrderId={this.state.selectedWorkOrderId}
                  onEquipmentSelected={selectedEquipment =>
                    this.navigateToEquipment(selectedEquipment.id)
                  }
                  onWorkOrderSelected={selectedWorkOrderCardId =>
                    this.navigateToWorkOrder(selectedWorkOrderCardId)
                  }
                  onAddEquipment={() => this.showDialog('equipment')}
                  onLocationMoved={this.onMoveLocation}
                  onLocationRemoved={this.onDeleteLocation}
                />
              )}
              {card.type == 'equipment' && (
                <EquipmentCard
                  mode={card.mode}
                  onSave={this.onEquipmentSave}
                  onEdit={this.onEquipmentEdit}
                  onCancel={this.onEquipmentCancel}
                  selectedEquipmentId={this.state.selectedEquipmentId}
                  selectedEquipmentPosition={
                    this.state.selectedEquipmentPosition
                  }
                  selectedLocationId={this.state.selectedLocationId}
                  selectedEquipmentType={this.state.selectedEquipmentType}
                  selectedWorkOrderId={this.state.selectedWorkOrderId}
                  onAttachingEquipmentToPosition={(
                    selectedEquipmentType,
                    selectedEquipmentPosition,
                  ) => {
                    ServerLogger.info(
                      LogEvents.ATTACH_EQUIPMENT_TO_POSITION_CLICKED,
                    );
                    this.setState({
                      selectedEquipmentType,
                      selectedEquipmentPosition,
                      card: ADD_EQUIPMENT_CARD,
                      dialogMode: 'hidden',
                    });
                  }}
                  onEquipmentClicked={equipmentId =>
                    this.navigateToEquipment(equipmentId)
                  }
                  onParentLocationClicked={selectedLocationId => {
                    this.navigateToLocation(selectedLocationId);
                  }}
                  onWorkOrderSelected={selectedWorkOrderCardId =>
                    this.navigateToWorkOrder(selectedWorkOrderCardId)
                  }
                />
              )}
            </InventoryErrorBoundary>
          </div>
        </div>
        <AddToLocationDialog
          key={`add_to_location_${this.state.dialogMode}`}
          show={
            this.state.dialogMode === 'hidden'
              ? 'location'
              : this.state.dialogMode
          }
          open={this.state.dialogMode !== 'hidden'}
          onClose={this.hideDialog}
          onEquipmentTypeSelected={selectedEquipmentType =>
            this.setState({
              selectedEquipmentType,
              card: ADD_EQUIPMENT_CARD,
              dialogMode: 'hidden',
            })
          }
          onLocationTypeSelected={selectedLocationType =>
            this.setState({
              selectedLocationType,
              card: ADD_LOCATION_CARD,
              dialogMode: 'hidden',
            })
          }
        />
      </>
    );
  }

  showDialog = (dialogMode: 'location' | 'equipment') => {
    ServerLogger.info(
      dialogMode === 'location'
        ? LogEvents.ADD_LOCATION_BUTTON_CLICKED
        : LogEvents.ADD_EQUIPMENT_BUTTON_CLICKED,
    );
    this.setState({dialogMode});
  };
  hideDialog = () => this.setState({dialogMode: 'hidden'});

  onEquipmentCancel = () => {
    this.setState(state => ({
      card: state.selectedEquipmentId
        ? SHOW_EQUIPMENT_CARD
        : SHOW_LOCATION_CARD,
    }));
  };

  onEquipmentSave = () => {
    ServerLogger.info(LogEvents.SAVE_EQUIPMENT_BUTTON_CLICKED);
    if (this.state.selectedEquipmentId) {
      this.setState({
        card: SHOW_EQUIPMENT_CARD,
      });
    } else if (this.state.selectedEquipmentPosition) {
      this.navigateToEquipment(
        this.state.selectedEquipmentPosition.parentEquipment.id,
      );
    } else {
      this.navigateToLocation(this.state.selectedLocationId);
    }
  };

  onEquipmentEdit = () => {
    ServerLogger.info(LogEvents.EDIT_EQUIPMENT_BUTTON_CLICKED);
    this.setState({
      card: EDIT_EQUIPMENT_CARD,
    });
  };

  onLocationCancel = () => {
    ServerLogger.info(LogEvents.LOCATION_CARD_CANCEL_BUTTON_CLICKED);
    this.setState(state => ({
      card: state.selectedEquipmentId
        ? SHOW_EQUIPMENT_CARD
        : SHOW_LOCATION_CARD,
    }));
  };

  onLocationEdit = () => {
    ServerLogger.info(LogEvents.EDIT_LOCATION_BUTTON_CLICKED);
    this.setState({
      card: EDIT_LOCATION_CARD,
    });
  };

  onLocationSave = (newLocationId: string) => {
    ServerLogger.info(LogEvents.SAVE_LOCATION_BUTTON_CLICKED);
    this.navigateToLocation(newLocationId);
  };

  onMoveLocation = (movedLocation: LocationMenu_location) => {
    ServerLogger.info(LogEvents.MOVE_LOCATION_BUTTON_CLICKED);
    this.props.enqueueSnackbar('Location moved successfuly', {
      children: key => (
        <SnackbarItem
          id={key}
          message={fbt(
            'Location moved successfuly',
            'Pop-over message when moving a location',
          )}
          variant="success"
        />
      ),
    });
    this.navigateToLocation(
      movedLocation?.parentLocation?.id || this.state.parentLocationId || '',
    );
  };

  onDeleteLocation = (deletedLocation: LocationMenu_location) => {
    ServerLogger.info(LogEvents.DELETE_LOCATION_BUTTON_CLICKED);
    this.props.enqueueSnackbar('Location removed successfuly', {
      children: key => (
        <SnackbarItem
          id={key}
          message={fbt(
            'Location removed successfuly',
            'Pop-over message when deleting a location',
          )}
          variant="success"
        />
      ),
    });
    this.navigateToLocation(
      deletedLocation?.parentLocation?.id || this.state.parentLocationId || '',
    );
  };
}

export default withStyles(styles)(
  withRouter(withAlert(withSnackbar(Inventory))),
);
