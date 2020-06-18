/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {EquipmentPort} from '../../common/Equipment';
import type {ServicePanel_service} from './__generated__/ServicePanel_service.graphql';

import AddEndpointToServiceDialog from './AddEndpointToServiceDialog';
import React, {useState} from 'react';
import ServiceMenu from './ServiceMenu';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';

type Props = {
  service: ServicePanel_service,
  onAddEndpoint: (port: EquipmentPort, role: string) => void,
};

export type EndpointDef = {
  name: string,
  equipmentTypeID: string,
  id: string,
};

const ServiceEndpointsMenu = (props: Props) => {
  const {service, onAddEndpoint} = props;
  const [addingEndpoint, setAddingEndpoint] = useState<?EndpointDef>(null);

  const serviceEPDefinitions = service.serviceType.endpointDefinitions;
  const items = serviceEPDefinitions
    .map(endpointDef => {
      return {
        caption: endpointDef.name,
        onClick: () => {
          ServerLogger.info(LogEvents.ADD_ENDPOINT_BUTTON_CLICKED);
          setAddingEndpoint({
            name: endpointDef.name,
            equipmentTypeID: endpointDef.equipmentType.id,
            id: endpointDef.id,
          });
        },
      };
    })
    .filter(x => !!x);

  const remainingItems = items.filter(
    item =>
      !service.endpoints.map(ep => ep?.definition?.name).includes(item.caption),
  );

  if (remainingItems.length == 0) {
    return null;
  }

  return (
    <ServiceMenu
      isOpen={!!addingEndpoint}
      onClose={() => setAddingEndpoint(null)}
      items={remainingItems}>
      <AddEndpointToServiceDialog
        service={service}
        onClose={() => setAddingEndpoint(null)}
        onAddEndpoint={port => {
          if (addingEndpoint) {
            onAddEndpoint(port, addingEndpoint.id);
            setAddingEndpoint(null);
          }
        }}
        endpointDef={addingEndpoint}
      />
    </ServiceMenu>
  );
};

export default ServiceEndpointsMenu;
