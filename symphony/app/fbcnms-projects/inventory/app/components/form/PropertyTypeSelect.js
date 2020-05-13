/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {OptionProps} from '@fbcnms/ui/components/design-system/Select/SelectMenu';
import type {PropertyKind} from '../../mutations/__generated__/AddEquipmentPortTypeMutation.graphql';
import type {PropertyType} from '../../common/PropertyType';

import AppContext from '@fbcnms/ui/context/AppContext';
import PropertyTypesTableDispatcher from './context/property_types/PropertyTypesTableDispatcher';
import React, {useContext, useMemo} from 'react';
import Select from '@fbcnms/ui/components/design-system/Select/Select';
import inventoryTheme from '../../common/theme';
import {PropertyTypeLabels} from '../PropertyTypeLabels';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  input: {
    ...inventoryTheme.textField,
    marginTop: '0px',
    marginBottom: '0px',
    width: '100%',
  },
}));

type PropertyTypeOption = {|
  kind: PropertyKind,
  nodeType: ?string,
|};

type Props = $ReadOnly<{|
  propertyType: PropertyType,
|}>;

const PropertyTypeSelect = ({propertyType}: Props) => {
  const classes = useStyles();
  const context = useContext(AppContext);
  const dispatch = useContext(PropertyTypesTableDispatcher);

  const getOptionKey = (type: string) =>
    `${PropertyTypeLabels[type].kind}_${type}`;

  const options = useMemo(
    () =>
      Object.keys(PropertyTypeLabels)
        .filter(
          (type: string) =>
            !PropertyTypeLabels[type].featureFlag ||
            context.isFeatureEnabled(PropertyTypeLabels[type].featureFlag),
        )
        .map((type: string) => ({
          key: getOptionKey(type),
          value: {
            kind: PropertyTypeLabels[type].kind,
            nodeType:
              PropertyTypeLabels[type].kind === 'node' ? type : undefined,
          },
          label: PropertyTypeLabels[type].label,
        })),
    [context],
  );

  const selectedValueIndex = useMemo(
    () =>
      options.findIndex(
        (op: OptionProps<PropertyTypeOption>) =>
          op.key ===
          getOptionKey(
            propertyType.nodeType != null && propertyType.nodeType != ''
              ? propertyType.nodeType
              : propertyType.type,
          ),
      ),
    [options, propertyType],
  );

  return (
    <Select
      className={classes.input}
      options={options}
      selectedValue={
        selectedValueIndex > -1 ? options[selectedValueIndex].value : null
      }
      onChange={value => {
        dispatch({
          type: 'UPDATE_PROPERTY_TYPE_KIND',
          id: propertyType.id,
          kind: value.kind,
          nodeType: value.nodeType,
        });
      }}
    />
  );
};

export default PropertyTypeSelect;
