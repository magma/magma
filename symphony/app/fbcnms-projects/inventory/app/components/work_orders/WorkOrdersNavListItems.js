/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import AssignmentIcon from '@material-ui/icons/Assignment';
import NavListItem from '@fbcnms/ui/components/NavListItem';
import ProjectsIcon from '@fbcnms/ui/icons/ProjectsIcon';
import React from 'react';
import WorkIcon from '@material-ui/icons/Work';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {useRouter} from '@fbcnms/ui/hooks';

export const WorkOrdersNavListItems = () => {
  const {relativeUrl} = useRouter();
  return [
    <NavListItem
      key={1}
      label="Work Orders"
      path={relativeUrl('/search')}
      icon={<WorkIcon />}
      onClick={() =>
        ServerLogger.info(LogEvents.WORK_ORDERS_SEARCH_NAV_CLICKED)
      }
    />,
    <NavListItem
      key={2}
      label="Projects"
      path={relativeUrl('/projects/search')}
      icon={<ProjectsIcon />}
      onClick={() => ServerLogger.info(LogEvents.PROJECTS_SEARCH_NAV_CLICKED)}
    />,
    <NavListItem
      key={3}
      label="Configure"
      path={relativeUrl('/configure')}
      icon={<AssignmentIcon />}
      onClick={() =>
        ServerLogger.info(LogEvents.WORK_ORDERS_CONFIGURE_NAV_CLICKED)
      }
    />,
  ];
};
