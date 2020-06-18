/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {PermissionEnforcement} from './admin/userManagement/utils/usePermissions';
import type {ViewContainerProps} from '@fbcnms/ui/components/design-system/View/ViewContainer';

import * as React from 'react';
import ListAltIcon from '@material-ui/icons/ListAlt';
import MapIcon from '@material-ui/icons/Map';
import ViewContainer from '@fbcnms/ui/components/design-system/View/ViewContainer';
import {FormContextProvider} from '../common/FormContext';
import {VARIANTS} from '@fbcnms/ui/components/design-system/View/ViewBody';
import {useMemo, useState} from 'react';

export const DisplayOptions = {
  table: 'table',
  map: 'map',
};
export type DisplayOptionTypes = $Keys<typeof DisplayOptions>;

type ViewToggleProps = $ReadOnly<{|
  onViewToggleClicked?: (id: string) => void,
|}>;

type Props = $ReadOnly<{|
  ...ViewContainerProps,
  ...ViewToggleProps,
  permissions: PermissionEnforcement,
|}>;

const InventoryView = (props: Props) => {
  const {
    onViewToggleClicked,
    header: headerProp,
    bodyVariant: bodyVariantProp,
    permissions,
    ...restProps
  } = props;
  const [selectedDisplayOption, setSelectedDisplayOption] = useState(
    DisplayOptions.table,
  );
  const bodyVariant = useMemo(
    () =>
      selectedDisplayOption === DisplayOptions.map
        ? VARIANTS.plain
        : bodyVariantProp,
    [bodyVariantProp, selectedDisplayOption],
  );
  const header = useMemo(() => {
    if (headerProp == null || onViewToggleClicked == null) {
      return headerProp;
    }
    const onViewOptionClicked = displayOptionId => {
      setSelectedDisplayOption(displayOptionId);
      if (onViewToggleClicked) {
        onViewToggleClicked(displayOptionId);
      }
    };
    return {
      ...headerProp,
      viewOptions: {
        onItemClicked: onViewOptionClicked,
        selectedButtonId: selectedDisplayOption,
        buttons: [
          {
            item: <ListAltIcon />,
            id: DisplayOptions.table,
          },
          {
            item: <MapIcon />,
            id: DisplayOptions.map,
          },
        ],
      },
    };
  }, [headerProp, onViewToggleClicked, selectedDisplayOption]);

  const viewProps: ViewContainerProps = useMemo(
    () => ({
      ...restProps,
      header,
      bodyVariant,
    }),
    [bodyVariant, header, restProps],
  );
  return (
    <FormContextProvider permissions={permissions}>
      <ViewContainer {...viewProps} />
    </FormContextProvider>
  );
};

export default InventoryView;
