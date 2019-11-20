/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {Equipment, EquipmentPort, Link} from '../../common/Equipment';
import type {WithStyles} from '@material-ui/core';

import AvailableLinksTable_links from './__generated__/AvailableLinksTable_links.graphql';
import EquipmentBreadcrumbs from '../equipment/EquipmentBreadcrumbs';
import LinkComponent from '@fbcnms/ui/components/Link';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Text from '@fbcnms/ui/components/design-system/Text';
import {createFragmentContainer, graphql} from 'react-relay';
import {sortLexicographically} from '@fbcnms/ui/utils/displayUtils';
import {withStyles} from '@material-ui/core/styles';

const styles = theme => ({
  noResultsRoot: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    marginTop: '100px',
  },
  noResultsLabel: {
    color: theme.palette.grey[600],
  },
  futureState: {
    textTransform: 'capitalize',
    maxWidth: '50px',
  },
});

type Props = {
  equipment: Equipment,
  links: AvailableLinksTable_links,
  onLinkSelected: (link: Link) => void,
} & WithStyles<typeof styles>;

type LinkPorts = Link & {
  srcPort: EquipmentPort,
  dstPort: EquipmentPort,
};

const showLinksByOrder = (
  srcEquipment: Equipment,
  links: AvailableLinksTable_links,
): Array<LinkPorts> => {
  return links
    .map(link => ({
      ...link,
      srcPort: link.ports[0],
      dstPort: link.ports[1],
    }))
    .map(link => {
      if (
        link.srcPort.parentEquipment.id != srcEquipment.id &&
        !link.srcPort.parentEquipment.positionHierarchy
          .map(position => position.parentEquipment.id)
          .includes(srcEquipment.id)
      ) {
        return {
          ...link,
          srcPort: link.dstPort,
          dstPort: link.srcPort,
        };
      }
      return link;
    })
    .sort((linkA, linkB) =>
      sortLexicographically(
        linkA.srcPort.definition.name,
        linkB.srcPort.definition.name,
      ),
    );
};

const AvailableLinksTable = (props: Props) => {
  const {equipment, links, onLinkSelected, classes} = props;
  const linksByOrder = showLinksByOrder(equipment, links);
  if (linksByOrder.length === 0) {
    return (
      <div className={classes.noResultsRoot}>
        <Text variant="h6" className={classes.noResultsLabel}>
          No available links out of
          {` ${equipment.equipmentType.name} ${equipment.name}`}
        </Text>
      </div>
    );
  }

  return (
    <Table>
      <TableHead>
        <TableRow>
          <TableCell>Equipment A (Selected)</TableCell>
          <TableCell>Port A</TableCell>
          <TableCell>Equipment B</TableCell>
          <TableCell>Port B</TableCell>
        </TableRow>
      </TableHead>
      <TableBody>
        {linksByOrder.map(link => {
          return (
            <TableRow key={`link_${link.id}`}>
              <TableCell>
                <EquipmentBreadcrumbs
                  equipment={link.srcPort.parentEquipment}
                  size="small"
                />
              </TableCell>
              <TableCell>
                <LinkComponent onClick={() => onLinkSelected(link)}>
                  {link.srcPort.definition.name}
                </LinkComponent>
              </TableCell>
              <TableCell>
                <EquipmentBreadcrumbs
                  equipment={link.dstPort.parentEquipment}
                  size="small"
                />
              </TableCell>
              <TableCell>{link.dstPort.definition.name}</TableCell>
            </TableRow>
          );
        })}
      </TableBody>
    </Table>
  );
};

export default withStyles(styles)(
  createFragmentContainer(AvailableLinksTable, {
    links: graphql`
      fragment AvailableLinksTable_links on Link @relay(plural: true) {
        id
        ports {
          parentEquipment {
            id
            name
            positionHierarchy {
              parentEquipment {
                id
              }
            }
            ...EquipmentBreadcrumbs_equipment
          }
          definition {
            id
            name
            type
          }
        }
      }
    `,
  }),
);
