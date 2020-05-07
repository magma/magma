/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import * as React from 'react';
import Checkbox from '@fbcnms/ui/components/design-system/Checkbox/Checkbox';
import useSideEffectCallback from './useSideEffectCallback';
import {
  HierarchyContextProvider,
  useHierarchyContext,
} from './HierarchyContext';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useEffect, useState} from 'react';

const useStyles = makeStyles(() => ({
  root: {
    display: 'flex',
    flexDirection: 'column',
  },
  children: {
    marginLeft: '24px',
  },
}));

type SubTreeProps = $ReadOnly<{|
  title: React.Node,
  onChange: (?boolean) => void,
  children?: React.Node,
|}>;

function CheckboxSubTree(props: SubTreeProps) {
  const {onChange, title, children} = props;
  const classes = useStyles();

  const hierarchyContext = useHierarchyContext();

  const propagateValue = useSideEffectCallback(() => {
    const allChildren = hierarchyContext.childrenValues;
    const hasFalseChild = allChildren.findKey(child => child === false) != null;
    const hasTrueChild = allChildren.findKey(child => child === true) != null;
    const hasNullChild = allChildren.findKey(child => child == null) != null;

    const childTypesCount =
      (hasFalseChild ? 1 : 0) + (hasTrueChild ? 1 : 0) + (hasNullChild ? 1 : 0);

    let aggregatedValue;
    if (childTypesCount === 0) {
      if (hierarchyContext.parentValue == null) {
        aggregatedValue = false;
      } else {
        return;
      }
    } else if (childTypesCount > 1 || hasNullChild) {
      aggregatedValue = null;
    } else if (hasFalseChild) {
      aggregatedValue = false;
    } else if (hasTrueChild) {
      aggregatedValue = true;
    } else {
      return;
    }

    if (aggregatedValue != hierarchyContext.parentValue) {
      onChange(aggregatedValue);
    }
  });

  useEffect(
    () => {
      propagateValue();
    }, // eslint-disable-next-line react-hooks/exhaustive-deps
    [hierarchyContext.childrenValues],
  );

  return (
    <div className={classes.root}>
      <Checkbox
        checked={hierarchyContext.parentValue === true}
        indeterminate={
          hierarchyContext.parentValue == null &&
          !hierarchyContext.childrenValues.isEmpty()
        }
        title={title}
        onChange={status => onChange(status === 'checked')}
      />
      <div className={classes.children}>{children}</div>
    </div>
  );
}

type Props = $ReadOnly<{|
  id: string,
  title: React.Node,
  value?: ?boolean,
  onChange?: ?(boolean) => void,
  children?: React.Node,
|}>;

export function HierarchicalCheckbox(props: Props) {
  const {id, value: propValue, title, children} = props;
  const [value, setValue] = useState<?boolean>(null);
  const hierarchyContext = useHierarchyContext();

  const updateMyValue = useCallback(
    newValue => {
      setValue(newValue);
      hierarchyContext.setChildValue(id, newValue);
    },
    [hierarchyContext, id],
  );

  useEffect(() => {
    if (hierarchyContext.childrenValues.has(id)) {
      if (hierarchyContext.parentValue != null) {
        updateMyValue(hierarchyContext.parentValue);
      }
    } else {
      updateMyValue(propValue ?? hierarchyContext.parentValue);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [hierarchyContext.parentValue, id, propValue]);

  return (
    <HierarchyContextProvider parentValue={value}>
      <CheckboxSubTree title={title} onChange={updateMyValue}>
        {children}
      </CheckboxSubTree>
    </HierarchyContextProvider>
  );
}

export default HierarchicalCheckbox;
