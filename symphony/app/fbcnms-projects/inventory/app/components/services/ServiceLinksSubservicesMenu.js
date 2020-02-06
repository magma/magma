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
import React, {useState} from 'react';
import ServiceMenu from './ServiceMenu';
import fbt from 'fbt';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';

type Props = {
  service: {id: string, name: string},
  onAddLink: (link: Link) => void,
};

const ServiceLinksSubservicesMenu = (props: Props) => {
  const {service, onAddLink} = props;
  const [addingEquipmentLink, setAddingEquipmentLink] = useState(false);

  return (
    <ServiceMenu
      isOpen={addingEquipmentLink}
      onClose={() => setAddingEquipmentLink(false)}
      items={[
        {
          caption: fbt(
            'Add Equipment Link',
            'Menu option to open a dialog to add link to a service',
          ),
          onClick: () => {
            ServerLogger.info(LogEvents.ADD_EQUIPMENT_LINK_BUTTON_CLICKED);
            setAddingEquipmentLink(true);
          },
        },
      ]}>
      <AddLinkToServiceDialog
        service={service}
        onClose={() => setAddingEquipmentLink(false)}
        onAddLink={link => {
          onAddLink(link);
          setAddingEquipmentLink(false);
        }}
      />
    </ServiceMenu>
  );
};

export default ServiceLinksSubservicesMenu;
