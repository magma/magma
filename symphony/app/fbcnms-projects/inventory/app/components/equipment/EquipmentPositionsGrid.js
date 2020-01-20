/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {Equipment, EquipmentPosition} from '../../common/Equipment';
import type {EquipmentType} from '../../common/EquipmentType';
import type {PositionDefinition} from '../../common/EquipmentType';
import type {WithStyles} from '@material-ui/core';

import CardSection from '../CardSection';
import EquipmentPositionItem from './EquipmentPositionItem';
import React from 'react';
import {createFragmentContainer, graphql} from 'react-relay';
import {getNonInstancePositionDefinitions} from '../../common/Equipment';
import {sortByIndex} from '../draggable/DraggableUtils';
import {withStyles} from '@material-ui/core/styles';

type Props = {
  equipment: Equipment,
  onAttachingEquipmentToPosition: (
    equipmentType: EquipmentType,
    position: EquipmentPosition,
  ) => void,
  onEquipmentPositionClicked: (equipmentId: string) => void,
  workOrderId: ?string,
  onWorkOrderSelected: (workOrderId: string) => void,
} & WithStyles<typeof styles>;

const styles = theme => ({
  root: {
    display: 'flex',
    flexWrap: 'wrap',
  },
  position: {
    marginBottom: theme.spacing(2),
    marginRight: theme.spacing(2),
  },
});

class EquipmentPositionsGrid extends React.PureComponent<Props> {
  render() {
    const {classes, equipment} = this.props;
    const {positions} = equipment;
    const {positionDefinitions} = equipment.equipmentType;
    const positionsToDisplay = [
      ...positions,
      ...getNonInstancePositionDefinitions(positions, positionDefinitions).map(
        this.getTemporaryPosition,
      ),
    ];
    if (positionsToDisplay.length === 0) {
      return null;
    }
    return (
      <CardSection title="Positions">
        <div className={classes.root}>
          {positionsToDisplay
            .slice()
            .sort((positionA, positionB) =>
              sortByIndex(positionA.definition, positionB.definition),
            )
            .map(position => (
              <div className={classes.position} key={position.id}>
                <EquipmentPositionItem
                  position={position}
                  onAttachingEquipmentToPosition={
                    this.props.onAttachingEquipmentToPosition
                  }
                  onEquipmentPositionClicked={
                    this.props.onEquipmentPositionClicked
                  }
                  workOrderId={this.props.workOrderId}
                  equipment={this.props.equipment}
                  onWorkOrderSelected={this.props.onWorkOrderSelected}
                />
              </div>
            ))}
        </div>
      </CardSection>
    );
  }

  getTemporaryPosition = (
    definition: PositionDefinition,
  ): EquipmentPosition => {
    const {equipment} = this.props;
    return {
      id: `PositionDefinition@tmp${definition.name}`,
      definition: definition,
      parentEquipment: equipment,
      attachedEquipment: null,
    };
  };
}

export default withStyles(styles)(
  createFragmentContainer(EquipmentPositionsGrid, {
    equipment: graphql`
      fragment EquipmentPositionsGrid_equipment on Equipment {
        id
        ...AddToEquipmentDialog_parentEquipment
        positions {
          id
          definition {
            id
            name
            index
            visibleLabel
          }
          attachedEquipment {
            id
            name
            futureState
            services {
              id
            }
          }
          parentEquipment {
            id
          }
        }
        equipmentType {
          positionDefinitions {
            id
            name
            index
            visibleLabel
          }
        }
      }
    `,
  }),
);
