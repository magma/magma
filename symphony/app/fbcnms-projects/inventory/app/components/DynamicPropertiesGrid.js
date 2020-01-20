/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {Property} from '../common/Property';
import type {PropertyType} from '../common/PropertyType';
import type {WithStyles} from '@material-ui/core';

import PropertyFormField from './form/PropertyFormField';
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import {createFragmentContainer, graphql} from 'react-relay';
import {getInitialPropertyFromType} from '../common/PropertyType';
import {
  getNonInstancePropertyTypes,
  sortPropertiesByIndex,
} from '../common/Property';
import {withStyles} from '@material-ui/core/styles';

type Props = {
  properties: Array<Property>,
  propertyTypes: Array<PropertyType>,
  hideTitle: boolean,
} & WithStyles<typeof styles>;

const styles = theme => ({
  subheader: {
    marginTop: theme.spacing(2),
    marginBottom: theme.spacing(2),
  },
  titleText: {
    fontWeight: 500,
  },
  root: {
    display: 'flex',
    flexWrap: 'wrap',
  },
  field: {
    width: '50%',
    marginBottom: '12px',
    paddingRight: '16px',
  },
});

const DynamicPropertiesGrid = (props: Props) => {
  const {classes, properties, propertyTypes, hideTitle} = props;
  if (properties.length == 0 && propertyTypes.length == 0) {
    return null;
  }
  const relevantPropertyTypes = getNonInstancePropertyTypes(
    properties,
    propertyTypes,
  );

  const allProps = [
    ...properties,
    ...relevantPropertyTypes.map(x => getInitialPropertyFromType(x)),
  ];

  return (
    <div>
      {!hideTitle && (
        <div className={classes.subheader}>
          <Text variant="subtitle1" className={classes.titleText}>
            Properties
          </Text>
        </div>
      )}
      <div className={classes.root}>
        {allProps.sort(sortPropertiesByIndex).map((property, i) => (
          <div className={classes.field} key={`property_${i}`}>
            <PropertyFormField property={property} />
          </div>
        ))}
      </div>
    </div>
  );
};

export default withStyles(styles)(
  createFragmentContainer(DynamicPropertiesGrid, {
    propertyTypes: graphql`
      fragment DynamicPropertiesGrid_propertyTypes on PropertyType
        @relay(plural: true) {
        id
        name
        index
        isInstanceProperty
        type
        stringValue
        intValue
        booleanValue
        latitudeValue
        longitudeValue
        rangeFromValue
        rangeToValue
        floatValue
      }
    `,
    properties: graphql`
      fragment DynamicPropertiesGrid_properties on Property
        @relay(plural: true) {
        ...PropertyFormField_property
        propertyType {
          id
          index
        }
      }
    `,
  }),
);
