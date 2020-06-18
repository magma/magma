/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {Service} from '../../common/Service';

import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';

type Props = {
  service: Service,
  anchorEl: ?HTMLElement,
  onClose: () => void,
  onViewService: (serviceId: string) => void,
};

const EquipmentServicesTableMenu = (props: Props) => {
  const {service, anchorEl, onClose, onViewService} = props;

  return (
    <>
      <Menu
        id="simple-menu"
        anchorEl={anchorEl}
        keepMounted
        open={!!anchorEl}
        onClose={onClose}>
        <MenuItem
          onClick={() => {
            ServerLogger.info(LogEvents.VIEW_EQUIPMENT_SERVICE_BUTTON_CLICKED);
            onViewService(service.id);
            onClose();
          }}>
          View Service
        </MenuItem>
      </Menu>
    </>
  );
};

export default EquipmentServicesTableMenu;
