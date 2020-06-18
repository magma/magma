/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import InventoryErrorBoundary from '../../common/InventoryErrorBoundary';
import React, {useMemo} from 'react';
import ServiceCardQueryRenderer from './ServiceCardQueryRenderer';
import ServiceComparisonView from './ServiceComparisonView';
import useRouter from '@fbcnms/ui/hooks/useRouter';
import {extractEntityIdFromUrl} from '../../common/RouterUtils';

const ServicesMain = () => {
  const {location} = useRouter();

  const selectedServiceCardId = useMemo(
    () => extractEntityIdFromUrl('service', location.search),
    [location],
  );

  if (selectedServiceCardId != null) {
    return (
      <InventoryErrorBoundary>
        <ServiceCardQueryRenderer serviceId={selectedServiceCardId} />
      </InventoryErrorBoundary>
    );
  }

  return (
    <InventoryErrorBoundary>
      <ServiceComparisonView />
    </InventoryErrorBoundary>
  );
};

export default ServicesMain;
