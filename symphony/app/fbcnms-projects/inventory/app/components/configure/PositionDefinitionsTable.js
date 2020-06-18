/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import {makeStyles} from '@material-ui/styles';
import type {PositionDefinitionsTable_positionDefinitions} from './__generated__/PositionDefinitionsTable_positionDefinitions.graphql';

import CardSection from '../CardSection';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Text from '@fbcnms/ui/components/design-system/Text';
import inventoryTheme from '../../common/theme';
import {createFragmentContainer, graphql} from 'react-relay';
import {sortByIndex} from '../draggable/DraggableUtils';

const useStyles = makeStyles(_theme => ({
  table: inventoryTheme.table,
  cell: {
    paddingLeft: '0px',
  },
}));

type Props = {
  positionDefinitions: PositionDefinitionsTable_positionDefinitions,
};

const PositionDefinitionsTable = (props: Props) => {
  const {positionDefinitions} = props;
  const classes = useStyles();

  if (positionDefinitions.length === 0) {
    return null;
  }

  return (
    <CardSection title="Positions">
      <Table component="div" className={classes.table}>
        <TableHead component="div">
          <TableRow component="div">
            <TableCell component="div" className={classes.cell}>
              Position Name
            </TableCell>
            <TableCell component="div" className={classes.cell}>
              Visible Label
            </TableCell>
          </TableRow>
        </TableHead>
        <TableBody component="div">
          {positionDefinitions
            .slice()
            .sort(sortByIndex)
            .map((definition, i) => (
              <TableRow component="div" key={`position_${i}`}>
                <TableCell className={classes.cell} component="div" scope="row">
                  <Text variant="body2">{definition.name}</Text>
                </TableCell>
                <TableCell className={classes.cell} component="div" scope="row">
                  <Text variant="body2">{definition.visibleLabel}</Text>
                </TableCell>
                <TableCell component="div" />
              </TableRow>
            ))}
        </TableBody>
      </Table>
    </CardSection>
  );
};

export default createFragmentContainer(PositionDefinitionsTable, {
  positionDefinitions: graphql`
    fragment PositionDefinitionsTable_positionDefinitions on EquipmentPositionDefinition
      @relay(plural: true) {
      id
      name
      index
      visibleLabel
    }
  `,
});
