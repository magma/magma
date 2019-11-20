/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {Link} from '../../common/Equipment';

import * as React from 'react';
import Button from '@fbcnms/ui/components/design-system/Button';
import DeleteIcon from '@material-ui/icons/Delete';
import IconButton from '@material-ui/core/IconButton';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import useRouter from '@fbcnms/ui/hooks/useRouter';
import {makeStyles} from '@material-ui/styles';

type Props = {
  links: Array<Link>,
  onDeleteLink?: (link: Link) => void,
};

const useStyles = makeStyles(theme => ({
  icon: {
    padding: '0px',
    marginLeft: theme.spacing(),
  },
}));

const ServiceLinksTable = (props: Props) => {
  const {links, onDeleteLink} = props;
  const classes = useStyles();
  const {history} = useRouter();

  const navigateToEquipment = (equipmentId: string) => {
    history.push(`/inventory/inventory?equipment=${equipmentId}`);
  };

  return (
    <Table>
      <TableHead>
        <TableRow>
          <TableCell>Equipment A</TableCell>
          <TableCell>Equipment A Port</TableCell>
          <TableCell>Equipment B</TableCell>
          <TableCell>Equipment B Port</TableCell>
          {onDeleteLink && <TableCell />}
        </TableRow>
      </TableHead>
      <TableBody>
        {links.map(link => {
          const portA = link.ports[0];
          const portB = link.ports[1];
          return (
            <TableRow key={`link_${link.id}`}>
              <TableCell>
                <Button
                  variant="text"
                  onClick={() => navigateToEquipment(portA.parentEquipment.id)}>
                  {portA.parentEquipment.name}
                </Button>
              </TableCell>
              <TableCell>{portA.definition.name}</TableCell>
              <TableCell>
                <Button
                  variant="text"
                  onClick={() => navigateToEquipment(portB.parentEquipment.id)}>
                  {portB.parentEquipment.name}
                </Button>
              </TableCell>
              <TableCell>{portB.definition.name}</TableCell>
              {onDeleteLink && (
                <TableCell>
                  <IconButton
                    onClick={() => onDeleteLink(link)}
                    color="primary"
                    className={classes.icon}>
                    <DeleteIcon />
                  </IconButton>
                </TableCell>
              )}
            </TableRow>
          );
        })}
      </TableBody>
    </Table>
  );
};

export default ServiceLinksTable;
