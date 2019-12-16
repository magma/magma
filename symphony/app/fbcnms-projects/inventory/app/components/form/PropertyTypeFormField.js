/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import type {PropertyType} from '../../common/PropertyType';
import type {WithStyles} from '@material-ui/core';

import FormField from '@fbcnms/ui/components/FormField';
import React from 'react';
import {createFragmentContainer, graphql} from 'react-relay';
import {getPropertyDefaultValue} from '../../common/PropertyType';
import {withStyles} from '@material-ui/core/styles';

type Props = {
  propertyType: PropertyType,
} & WithStyles<typeof styles>;

const styles = _theme => ({});

class PropertyTypeFormField extends React.Component<Props> {
  render() {
    const {propertyType} = this.props;
    return (
      <FormField
        label={propertyType.name}
        value={getPropertyDefaultValue(propertyType) || '-'}
      />
    );
  }
}

export default withStyles(styles)(
  createFragmentContainer(PropertyTypeFormField, {
    propertyType: graphql`
      fragment PropertyTypeFormField_propertyType on PropertyType {
        id
        name
        type
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
      }
    `,
  }),
);
