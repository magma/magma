/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ContextRouter} from 'react-router-dom';
import type {Location} from '../../common/Location';
import type {WithStyles} from '@material-ui/core';

import * as React from 'react';
import Breadcrumbs from '@fbcnms/ui/components/Breadcrumbs';
import {InventoryAPIUrls} from '../../common/InventoryAPI';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {createFragmentContainer, graphql} from 'react-relay';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

const styles = _theme => ({
  breadcrumbs: {
    display: 'flex',
    alignItems: 'flex-start',
  },
  locationNameContainer: {
    display: 'flex',
    alignItems: 'center',
    flexWrap: 'wrap',
  },
});

type Props = ContextRouter & {
  locationDetails: Location,
  hideTypes: boolean,
  navigateOnClick?: boolean,
  size?: 'default' | 'small' | 'large',
} & WithStyles<typeof styles>;

const LocationBreadcrumbsTitle = (props: Props) => {
  const {
    classes,
    locationDetails,
    hideTypes,
    size,
    navigateOnClick = true,
  } = props;

  const navigateToLocation = React.useCallback(
    (selectedLocationId: string) => {
      ServerLogger.info(LogEvents.NAVIGATE_TO_LOCATION, {
        locationId: selectedLocationId,
      });

      props.history.push(InventoryAPIUrls.location(selectedLocationId));
    },
    [props.history],
  );

  const onBreadcrumbClicked = React.useCallback(
    id => {
      ServerLogger.info(LogEvents.LOCATION_CARD_BREADCRUMB_CLICKED, {
        locationId: id,
      });
      if (id && navigateOnClick) {
        navigateToLocation(id);
      }
    },
    [navigateOnClick, navigateToLocation],
  );

  return (
    <div className={classes.breadcrumbs}>
      <div className={classes.locationNameContainer}>
        <Breadcrumbs
          breadcrumbs={[
            ...locationDetails.locationHierarchy,
            locationDetails,
          ].map(l => ({
            id: l.id,
            name: l.name,
            subtext: hideTypes ? null : l.locationType.name,
            onClick: () => onBreadcrumbClicked(l.id),
          }))}
          size={size}
        />
      </div>
    </div>
  );
};

LocationBreadcrumbsTitle.defaultProps = {
  size: 'default',
};

export default withStyles(styles)(
  withRouter(
    createFragmentContainer(LocationBreadcrumbsTitle, {
      locationDetails: graphql`
        fragment LocationBreadcrumbsTitle_locationDetails on Location {
          id
          name
          locationType {
            name
          }
          locationHierarchy {
            id
            name
            locationType {
              name
            }
          }
        }
      `,
    }),
  ),
);
