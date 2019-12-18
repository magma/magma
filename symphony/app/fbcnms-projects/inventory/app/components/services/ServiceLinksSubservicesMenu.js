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

import AddLinkToServiceDialog from './AddLinkToServiceDialog';
import Dialog from '@material-ui/core/Dialog';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import React, {useState} from 'react';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {makeStyles} from '@material-ui/styles';

type Props = {
  service: {id: string, name: string},
  anchorEl: ?HTMLElement,
  onClose: () => void,
  onAddLink: (link: Link) => void,
};

const useStyles = makeStyles({
  dialog: {
    width: '80%',
    maxWidth: '1280px',
    height: '90%',
    maxHeight: '800px',
  },
});

const ServiceLinksSubservicesMenu = (props: Props) => {
  const classes = useStyles();
  const {service, anchorEl, onClose, onAddLink} = props;
  const [addingEquipmentLink, setAddingEquipmentLink] = useState(false);

  return (
    <>
      <Menu anchorEl={anchorEl} keepMounted open={!!anchorEl} onClose={onClose}>
        <MenuItem
          onClick={() => {
            ServerLogger.info(LogEvents.ADD_EQUIPMENT_LINK_BUTTON_CLICKED);
            setAddingEquipmentLink(true);
            onClose();
          }}>
          Add Equipment Link
        </MenuItem>
      </Menu>
      <Dialog
        open={addingEquipmentLink}
        onClose={() => setAddingEquipmentLink(false)}
        maxWidth={false}
        fullWidth={true}
        classes={{paperFullWidth: classes.dialog}}>
        <AddLinkToServiceDialog
          service={service}
          onClose={() => setAddingEquipmentLink(false)}
          onAddLink={link => {
            onAddLink(link);
            setAddingEquipmentLink(false);
          }}
        />
      </Dialog>
    </>
  );
};

export default ServiceLinksSubservicesMenu;
