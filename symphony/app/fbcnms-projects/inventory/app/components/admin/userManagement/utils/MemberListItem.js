/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ToggleButtonDisplay} from './ListItem';

import * as React from 'react';
import CheckIcon from '@fbcnms/ui/components/design-system/Icons/Indications/CheckIcon';
import ListItem from './ListItem';
import PlusIcon from '@fbcnms/ui/components/design-system/Icons/Actions/PlusIcon';
import Strings from '@fbcnms/strings/Strings';
import fbt from 'fbt';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  itemDetails: {
    flexBasis: '10px',
    flexGrow: 1,
    flexShrink: 1,
  },
}));

export type MemberItem<T> = $ReadOnly<{|
  item: T,
  isMember: boolean,
|}>;

export type AssigenmentButtonProp = $ReadOnly<{|
  assigmentButton?: ?ToggleButtonDisplay,
|}>;

type Props<T> = $ReadOnly<{|
  member: MemberItem<T>,
  className?: ?string,
  ...AssigenmentButtonProp,
  onAssignToggle: () => void,
  children: React.Node,
|}>;

export default function MemberListItem<T>(props: Props<T>) {
  const {member, assigmentButton, onAssignToggle, children, className} = props;
  const classes = useStyles();

  const toggleButton = {
    isOn: member.isMember,
    displayVariants: assigmentButton,
    onToggleClicked: onAssignToggle,
    onContent: {
      regularContent: (
        <>
          <CheckIcon color="inherit" size="small" />
          <fbt desc="">Added</fbt>
        </>
      ),
      hoverContent: Strings.common.removeButton,
      onProcessContent: <fbt desc="">Removing</fbt>,
      skin: 'gray',
    },
    offContent: {
      regularContent: (
        <>
          <PlusIcon color="inherit" size="small" />
          {Strings.common.addButton}
        </>
      ),
      onProcessContent: <fbt desc="">Adding</fbt>,
      skin: 'primary',
    },
  };
  return (
    <ListItem className={className} toggleButton={toggleButton}>
      <div className={classes.itemDetails}>{children}</div>
    </ListItem>
  );
}
