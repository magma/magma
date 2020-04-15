/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {Property} from '../../common/Property';

import * as React from 'react';
import FormField from '@fbcnms/ui/components/FormField';
import NodePropertyValue from '../NodePropertyValue';
import {createFragmentContainer, graphql} from 'react-relay';
import {getPropertyValue} from '../../common/Property';

type Props = {
  property: Property,
};

class PropertyFormField extends React.Component<Props> {
  render() {
    const {property} = this.props;
    const propType = property.propertyType ? property.propertyType : property;
    return (
      <FormField
        label={property.propertyType.name}
        value={
          propType.type === 'node' && propType.nodeType != null ? (
            <NodePropertyValue
              type={propType.nodeType}
              value={property.nodeValue}
            />
          ) : (
            getPropertyValue(property) ?? ''
          )
        }
      />
    );
  }
}

export default createFragmentContainer(PropertyFormField, {
  property: graphql`
    fragment PropertyFormField_property on Property {
      id
      propertyType {
        id
        name
        type
        nodeType
        index
        stringValue
        intValue
        booleanValue
        floatValue
        latitudeValue
        longitudeValue
        rangeFromValue
        rangeToValue
        isEditable
        isInstanceProperty
        isMandatory
        category
        isDeleted
      }
      stringValue
      intValue
      floatValue
      booleanValue
      latitudeValue
      longitudeValue
      rangeFromValue
      rangeToValue
      nodeValue {
        id
        name
      }
    }
  `,
});
