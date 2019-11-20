/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import DynamicPropertiesGrid from '../DynamicPropertiesGrid';
import React from 'react';
import type {Service} from '../../common/Service';

type Props = {
  className?: string,
  service: Service,
};

const ServiceDetails = (props: Props) => {
  const {className, service} = props;
  return (
    <div className={className}>
      <DynamicPropertiesGrid
        properties={service.properties}
        propertyTypes={service.serviceType.propertyTypes}
        hideTitle={true}
      />
    </div>
  );
};

export default ServiceDetails;
