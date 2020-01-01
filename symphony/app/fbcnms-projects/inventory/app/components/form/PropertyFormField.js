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
import type {WithStyles} from '@material-ui/core';

import FormField from '@fbcnms/ui/components/FormField';
import PropertyValue from './PropertyValue';
import React from 'react';
import {createFragmentContainer, graphql} from 'react-relay';
import {withStyles} from '@material-ui/core/styles';

type Props = {
  property: Property,
} & WithStyles<typeof styles>;

const styles = _theme => ({});

class PropertyFormField extends React.Component<Props> {
  render() {
    const {property} = this.props;
    return (
      <FormField
        label={property.propertyType.name}
        value={<PropertyValue property={property} />}
      />
    );
  }
}

export default withStyles(styles)(
  createFragmentContainer(PropertyFormField, {
    property: graphql`
      fragment PropertyFormField_property on Property {
        id
        propertyType {
          id
          name
          type
          isEditable
          isMandatory
          isInstanceProperty
          stringValue
        }
        stringValue
        intValue
        floatValue
        booleanValue
        latitudeValue
        longitudeValue
        rangeFromValue
        rangeToValue
        equipmentValue {
          id
          name
        }
        locationValue {
          id
          name
        }
        serviceValue {
          id
          name
        }
      }
    `,
  }),
);
