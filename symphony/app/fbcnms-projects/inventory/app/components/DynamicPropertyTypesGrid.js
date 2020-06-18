/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {PropertyType} from '../common/PropertyType';

import PropertyTypeFormField from './form/PropertyTypeFormField';
import React from 'react';
import {createFragmentContainer, graphql} from 'react-relay';
import {makeStyles} from '@material-ui/styles';
import {sortByIndex} from '../components/draggable/DraggableUtils';

const useStyles = makeStyles(_theme => ({
  root: {
    display: 'flex',
    flexWrap: 'wrap',
  },
  field: {
    minWidth: '50%',
    paddingRight: '16px',
    marginBottom: '12px',
  },
}));

type Props = {
  propertyTypes: Array<PropertyType>,
};

const DynamicPropertyTypesGrid = (props: Props) => {
  const propertyTypes = props.propertyTypes.slice().sort(sortByIndex);
  const classes = useStyles();
  return (
    <div className={classes.root}>
      {propertyTypes.map((propertyType, i) => (
        <div key={`property_${propertyType.id}`} className={classes.field}>
          <PropertyTypeFormField
            key={`property_${i}`}
            propertyType={propertyType}
          />
        </div>
      ))}
    </div>
  );
};

export default createFragmentContainer(DynamicPropertyTypesGrid, {
  propertyTypes: graphql`
    fragment DynamicPropertyTypesGrid_propertyTypes on PropertyType
      @relay(plural: true) {
      ...PropertyTypeFormField_propertyType
      id
      index
    }
  `,
});
