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
import type {TextVariant} from '@fbcnms/ui/theme/symphony';
import type {WithStyles} from '@material-ui/core';

import Breadcrumbs from '@fbcnms/ui/components/Breadcrumbs';
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {createFragmentContainer, graphql} from 'react-relay';
import {withStyles} from '@material-ui/core/styles';

import nullthrows from '@fbcnms/util/nullthrows';

const styles = theme => ({
  breadcrumbs: {
    display: 'flex',
    alignItems: 'flex-start',
  },
  position: {
    display: 'flex',
    marginLeft: '4px',
  },
  positionName: {
    fontSize: theme.typography.pxToRem(13),
  },
  equipmentType: {
    fontSize: theme.typography.pxToRem(13),
    marginRight: '4px',
  },
  equipmentSubtext: {
    display: 'flex',
    alignItems: 'center',
  },
  seperator: {
    lineHeight: '16px',
  },
});

type Props = {
  equipment: Equipment,
  onParentLocationClicked?: (locationId: string) => void,
  onEquipmentClicked?: (equipmentId: string) => void,
  size?: 'default' | 'small' | 'large',
  showSelfEquipment: boolean,
  textClassName?: string,
  variant?: TextVariant,
} & WithStyles<typeof styles>;

const EquipmentBreadcrumbs = (props: Props) => {
  const {
    classes,
    equipment,
    onEquipmentClicked,
    onParentLocationClicked,
    size,
    showSelfEquipment,
    textClassName,
    variant,
  } = props;

  const positionSubText = pos => (
    <div className={classes.equipmentSubtext}>
      <Text className={classes.equipmentType}>
        {pos.parentEquipment.equipmentType.name}
      </Text>
      <Text className={classes.seperator} variant="body2">
        &#8226;
      </Text>
      <div className={classes.position}>
        <Text className={classes.positionName}>{pos.definition.name}</Text>
      </div>
    </div>
  );

  const onLocationClickedCallback = locationId => {
    ServerLogger.info(LogEvents.EQUIPMENT_CARD_LOCATION_BREADCRUMB_CLICKED, {
      locationId,
    });
    onParentLocationClicked && onParentLocationClicked(locationId);
  };
  const onEquipmentClickedCallback = id => {
    ServerLogger.info(LogEvents.EQUIPMENT_CARD_EQUIPMENT_BREADCRUMB_CLICKED, {
      equipmentId: id,
    });
    onEquipmentClicked && onEquipmentClicked(id);
  };
  const breadcrumbs = [
    ...equipment.locationHierarchy.map(l => ({
      id: l.id,
      name: l.name,
      subtext: size === 'small' ? null : l.locationType.name,
      ...(onParentLocationClicked && {
        onClick: () => onLocationClickedCallback(l.id),
      }),
    })),
    ...equipment.positionHierarchy.map(pos => ({
      id: pos.id,
      name: nullthrows(pos.parentEquipment).name,
      subtext: size === 'small' ? null : positionSubText(pos),
      ...(onEquipmentClicked && {
        onClick: () =>
          onEquipmentClickedCallback(nullthrows(pos.parentEquipment).id),
      }),
    })),
    ...(showSelfEquipment
      ? [
          {
            id: equipment.id,
            name: equipment.name,
            subtext: size === 'small' ? null : equipment.equipmentType.name,
            ...(onEquipmentClicked && {
              onClick: () => onEquipmentClickedCallback(equipment.id),
            }),
          },
        ]
      : []),
  ];
  return (
    <Breadcrumbs
      breadcrumbs={breadcrumbs}
      size={size}
      variant={variant}
      textClassName={textClassName}
    />
  );
};

EquipmentBreadcrumbs.defaultProps = {
  size: 'default',
  showTypes: true,
  showSelfEquipment: true,
};

export default withStyles(styles)(
  createFragmentContainer(EquipmentBreadcrumbs, {
    equipment: graphql`
      fragment EquipmentBreadcrumbs_equipment on Equipment {
        id
        name
        equipmentType {
          id
          name
        }
        locationHierarchy {
          id
          name
          locationType {
            name
          }
        }
        positionHierarchy {
          id
          definition {
            id
            name
            visibleLabel
          }
          parentEquipment {
            id
            name
            equipmentType {
              id
              name
            }
          }
        }
      }
    `,
  }),
);
