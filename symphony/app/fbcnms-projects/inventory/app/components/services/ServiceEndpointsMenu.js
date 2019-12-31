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

import AddEndpointToServiceDialog from './AddEndpointToServiceDialog';
import React, {useState} from 'react';
import ServiceMenu from './ServiceMenu';
import fbt from 'fbt';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';

type Props = {
  service: {id: string, name: string},
  onAddEndpoint: (port: EquipmentPort, role: 'consumer' | 'provider') => void,
};

const ServiceEndpointsMenu = (props: Props) => {
  const {service, onAddEndpoint} = props;
  const [addingEndpoint, setAddingEndpoint] = useState<?(
    | 'consumer'
    | 'provider'
  )>(null);

  return (
    <ServiceMenu
      isOpen={!!addingEndpoint}
      onClose={() => setAddingEndpoint(null)}
      items={[
        {
          label: fbt(
            'Add Consumer Endpoint',
            'Menu option to open a dialog to add customer endpoint to a service',
          ),
          onClick: () => {
            ServerLogger.info(LogEvents.ADD_CONSUMER_ENDPOINT_BUTTON_CLICKED);
            setAddingEndpoint('consumer');
          },
        },
        {
          label: fbt(
            'Add Provider Endpoint',
            'Menu option to open a dialog to add provider endpoint to a service',
          ),
          onClick: () => {
            ServerLogger.info(LogEvents.ADD_PROVIDER_ENDPOINT_BUTTON_CLICKED);
            setAddingEndpoint('provider');
          },
        },
      ]}>
      <AddEndpointToServiceDialog
        service={service}
        onClose={() => setAddingEndpoint(null)}
        onAddEndpoint={port => {
          if (addingEndpoint) {
            onAddEndpoint(port, addingEndpoint);
            setAddingEndpoint(null);
          }
        }}
        endpointRole={addingEndpoint ?? 'consumer'}
      />
    </ServiceMenu>
  );
};

export default ServiceEndpointsMenu;
