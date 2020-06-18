/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {Equipment, EquipmentPort} from '../../common/Equipment';
import type {WithStyles} from '@material-ui/core';

import * as React from 'react';
import Breadcrumbs from '@fbcnms/ui/components/Breadcrumbs';
import {fadedSea} from '@fbcnms/ui/theme/colors';
import {withStyles} from '@material-ui/core/styles';

const styles = theme => ({
  root: {
    display: 'flex',
  },
  breadcrumb: {
    display: 'flex',
    '&:first-child': {
      '& $equipmentName': {
        fontWeight: 'normal',
        color: theme.palette.common.black,
      },
    },
    '&:last-child': {
      '& $equipmentName': {
        color: theme.palette.primary.main,
      },
    },
  },
  equipmentName: {
    fontSize: '13px',
  },
  equipmentLink: {
    cursor: 'pointer',
    fontSize: '13px',
    color: fadedSea,
    fontWeight: 600,
  },
  chevronText: {
    marginLeft: '4px',
    marginRight: '4px',
    fontSize: '13px',
    color: theme.palette.common.black,
  },
});

export type EquipmentPortWithBreadcrumbs = EquipmentPort & {
  breadcrumbs: Array<Equipment>,
};

type Props = {
  port: EquipmentPortWithBreadcrumbs,
  onPortEquipmentClicked?: (equipmentId: string) => void,
} & WithStyles<typeof styles>;

const EquipmentBreadcrumbsTitle = (props: Props) => {
  const {classes, port, onPortEquipmentClicked} = props;

  return (
    <div className={classes.root}>
      <Breadcrumbs
        breadcrumbs={port.breadcrumbs.map(b => ({
          id: b.id,
          name: b.name,
        }))}
        onBreadcrumbClicked={id =>
          onPortEquipmentClicked && onPortEquipmentClicked(id)
        }
        size="small"
      />
    </div>
  );
};

export default withStyles(styles)(EquipmentBreadcrumbsTitle);
