/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import LinkEditDialog from './LinkEditDialog';
import OptionsPopoverButton from '../OptionsPopoverButton';
import PortEditDialog from './PortEditDialog';
import PortsConnectedStateDialog from './PortsConnectedStateDialog';
import React, {useState} from 'react';
import fbt from 'fbt';
import nullthrows from '@fbcnms/util/nullthrows';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {getInitialPropertyFromType} from '../../common/PropertyType';
import {getNonInstancePropertyTypes} from '../../common/Property';

import type {EquipmentPort, Link} from '../../common/Equipment';

type Props = {
  port: EquipmentPort,
  workOrderId: ?string,
};

const EquipmentPortsTableMenu = (props: Props) => {
  const {port, workOrderId} = props;
  const [selectedPort, setSelectedPort] = useState<?EquipmentPort>(null);
  const [editingLink, setEditingLink] = useState<?Link>(null);
  const [editingPort, setEditingPort] = useState<?EquipmentPort>(null);

  const relevantPropertyTypes = getNonInstancePropertyTypes(
    port.properties,
    port.definition.portType?.propertyTypes ?? [],
  );

  const portProperties = [
    ...(port.properties ?? []),
    ...relevantPropertyTypes.map(getInitialPropertyFromType),
  ];

  const menuOptions = [];
  if (portProperties.length > 0) {
    menuOptions.push({
      onClick: () => {
        ServerLogger.info(LogEvents.EDIT_EQUIPMENT_PORT_BUTTON_CLICKED);
        setEditingPort(port);
      },
      caption: fbt(
        'Edit Port Properties',
        'Caption for menu option for edit port properties dialog opening',
      ),
    });
  }
  if (!!port.link) {
    menuOptions.push({
      onClick: () => {
        ServerLogger.info(LogEvents.EDIT_LINK_CLICKED);
        setEditingLink(port.link);
      },
      caption: fbt(
        'Edit Link Properties',
        'Caption for menu option for editing port link properties',
      ),
    });
    menuOptions.push({
      onClick: () => {
        ServerLogger.info(LogEvents.DISCONNECT_PORTS_CLICKED);
        setSelectedPort(port);
      },
      caption: fbt(
        'Remove Link',
        'Caption for menu option for removing port link',
      ),
    });
  } else {
    menuOptions.push({
      onClick: () => {
        ServerLogger.info(LogEvents.ADD_LINK_CLICKED);
        setSelectedPort(port);
      },
      caption: fbt('Add Link', 'Caption for menu option for adding port link'),
    });
  }

  return (
    <>
      <OptionsPopoverButton options={menuOptions} />

      {selectedPort ? (
        <PortsConnectedStateDialog
          mode={!!selectedPort.link ? 'disconnect' : 'connect'}
          equipment={nullthrows(selectedPort.parentEquipment)}
          port={selectedPort}
          workOrderId={workOrderId}
          link={selectedPort.link}
          open={true}
          onClose={() => setSelectedPort(null)}
        />
      ) : null}
      {editingPort && (
        <PortEditDialog
          key={editingPort.id}
          port={editingPort}
          onClose={() => setEditingPort(null)}
        />
      )}
      {editingLink && (
        <LinkEditDialog
          key={editingLink.id}
          link={editingLink}
          onClose={() => setEditingLink(null)}
        />
      )}
    </>
  );
};

export default EquipmentPortsTableMenu;
