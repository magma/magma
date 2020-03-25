/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {Link} from '../../common/Equipment';

import ActiveEquipmentIcon from '@fbcnms/ui/icons/ActiveEquipmentIcon';
import EquipmentIcon from '@fbcnms/ui/icons/EquipmentIcon';
import React from 'react';
import TableRowOptionsButton from '../TableRowOptionsButton';
import Text from '@fbcnms/ui/components/design-system/Text';
import classNames from 'classnames';
import symphony from '@fbcnms/ui/theme/symphony';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {makeStyles} from '@material-ui/styles';

type Props = {
  link: Link,
  onDeleteLink: () => void,
};

const useStyles = makeStyles(theme => ({
  root: {
    display: 'flex',
    '&:hover': {
      backgroundColor: symphony.palette.B50,
      '& $moreButton': {
        display: 'block',
      },
      '& $icon': {
        display: 'none',
      },
      '& $activeIcon': {
        display: 'block',
      },
    },
  },
  linkRow: {
    flexGrow: 1,
    padding: '6px 32px',
    position: 'relative',
  },
  line: {
    display: 'flex',
    alignItems: 'start',
  },
  icon: {
    padding: '0px',
    marginLeft: theme.spacing(),
  },
  separator: {
    borderBottom: `1px solid ${symphony.palette.B500}`,
    margin: '12px 24px 0px 24px',
    width: '24px',
  },
  moreButton: {
    position: 'absolute',
    right: '4px',
    top: '8px',
    padding: '4px',
    display: 'none',
    '&:hover': {
      color: symphony.palette.B600,
      backgroundColor: 'transparent',
    },
  },
  emptyIcon: {
    width: '24px',
    marginRight: '12px',
  },
  componentName: {
    display: 'block',
    textOverflow: 'ellipsis',
    width: 'calc(50% - 72px)',
    overflow: 'hidden',
  },
  portName: {
    color: symphony.palette.D500,
  },
  emptySeparator: {
    margin: '0px 24px',
    width: '24px',
  },
  icon: {
    display: 'block',
    marginRight: '12px',
  },
  activeIcon: {
    display: 'none',
    marginRight: '12px',
  },
}));

const ServiceLinkDetails = (props: Props) => {
  const classes = useStyles();
  const {link, onDeleteLink} = props;
  return (
    <div className={classes.root}>
      <div className={classes.linkRow}>
        <div className={classes.line}>
          <EquipmentIcon className={classes.icon} />
          <ActiveEquipmentIcon className={classes.activeIcon} />
          <Text variant="subtitle2" className={classes.componentName}>
            {link.ports[0].parentEquipment.name}
          </Text>
          <div className={classes.separator} />
          <EquipmentIcon className={classes.icon} />
          <ActiveEquipmentIcon className={classes.activeIcon} />
          <Text variant="subtitle2" className={classes.componentName}>
            {link.ports[1].parentEquipment.name}
          </Text>
        </div>
        <div className={classes.line}>
          <div className={classes.emptyIcon} />
          <Text
            variant="body2"
            className={classNames(classes.componentName, classes.portName)}>
            {link.ports[0].definition.name}
          </Text>
          <div className={classes.emptySeparator} />
          <div className={classes.emptyIcon} />
          <Text
            variant="body2"
            className={classNames(classes.componentName, classes.portName)}>
            {link.ports[1].definition.name}
          </Text>
        </div>
      </div>

      <TableRowOptionsButton
        options={[
          {
            caption: 'Remove Link',
            onClick: () => {
              ServerLogger.info(LogEvents.DELETE_SERVICE_LINK_BUTTON_CLICKED);
              onDeleteLink();
            },
          },
        ]}
      />
    </div>
  );
};

export default ServiceLinkDetails;
