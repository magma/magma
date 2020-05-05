/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ServiceEndpoint} from '../../common/Service';

import ActiveConsumerEndpointIcon from '@fbcnms/ui/icons/ActiveConsumerEndpointIcon';
import ActiveProviderEndpointIcon from '@fbcnms/ui/icons/ActiveProviderEndpointIcon';
import EndpointIcon from '@fbcnms/ui/icons/EndpointIcon';
import EquipmentBreadcrumbs from '../equipment/EquipmentBreadcrumbs';
import OptionsPopoverButton from '../OptionsPopoverButton';
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import camelCase from 'lodash/camelCase';
import classNames from 'classnames';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {makeStyles} from '@material-ui/styles';

type Props = {
  endpoint: ServiceEndpoint,
  onDeleteEndpoint: () => void,
};

const useStyles = makeStyles(_ => ({
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
  detail: {
    display: 'flex',
    alignItems: 'start',
  },
  icon: {
    padding: '0px',
    marginLeft: '8px',
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
  componentName: {
    display: 'block',
    textOverflow: 'ellipsis',
    width: 'calc(100% - 32px)',
    overflow: 'hidden',
  },
  portName: {
    color: symphony.palette.D500,
  },
  locationName: {
    color: symphony.palette.D500,
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

const ServiceEndpointDetails = (props: Props) => {
  const classes = useStyles();
  const {endpoint, onDeleteEndpoint} = props;
  return (
    <div className={classes.root}>
      <div className={classes.linkRow}>
        <div className={classes.detail}>
          <EndpointIcon className={classes.icon} />
          {endpoint.definition.role == 'CONSUMER' ? (
            <ActiveConsumerEndpointIcon
              variant="small"
              className={classes.activeIcon}
            />
          ) : (
            <ActiveProviderEndpointIcon
              variant="small"
              className={classes.activeIcon}
            />
          )}
          <div>
            <Text variant="subtitle2" className={classes.componentName}>
              {`${endpoint.equipment.name} (${camelCase(
                endpoint.definition.name,
              )})`}
            </Text>
            <Text
              variant="body2"
              className={classNames(classes.componentName, classes.portName)}>
              {endpoint.definition.role}
            </Text>
            <EquipmentBreadcrumbs
              equipment={endpoint.equipment}
              showSelfEquipment={false}
              variant="body2"
              className={classes.componentName}
              textClassName={classes.locationName}
            />
          </div>
        </div>
      </div>

      <OptionsPopoverButton
        options={[
          {
            caption: fbt(
              'Remove Endpoint',
              'Menu option to delete endpoint pressed',
            ),
            onClick: () => {
              ServerLogger.info(
                LogEvents.DELETE_SERVICE_ENDPOINT_BUTTON_CLICKED,
              );
              onDeleteEndpoint();
            },
          },
        ]}
      />
    </div>
  );
};

export default ServiceEndpointDetails;
