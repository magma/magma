/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import * as React from 'react';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import CardHeader from '@material-ui/core/CardHeader';
import InventoryQueryRenderer from '../InventoryQueryRenderer';
import Text from '@fbcnms/ui/components/design-system/Text';
import WorkOrderPopover from '../work_orders/WorkOrderPopover';
import emptyFunction from '@fbcnms/util/emptyFunction';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import useRouter from '@fbcnms/ui/hooks/useRouter';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {graphql} from 'relay-runtime';
import {makeStyles} from '@material-ui/styles';

const cardTop = '116px';
const useStyles = makeStyles(theme => ({
  card: {
    display: 'flex',
    flexDirection: 'column',
    position: 'absolute',
    right: '10px',
    top: cardTop,
    maxWidth: '40%',
    width: '416px',
    maxHeight: `calc(100% - ${cardTop} - 40px)`,
    borderRadius: '8px',
    overflow: 'hidden',
  },
  cardHeader: {
    paddingBottom: '8px',
    paddingTop: '8px',
    alignItems: 'baseline',
    borderBottom: `2px solid ${symphony.palette.separator}`,
  },
  cardHeaderContent: {
    paddingTop: '4px',
  },
  cardContent: {
    overflowY: 'auto',
    overflowX: 'hidden',
    padding: '0',
  },
  workOrderBlock: {
    overflowY: 'hidden',
    padding: '0',
    borderBottom: `2px solid ${symphony.palette.separatorLight}`,
    '&:last-child': {
      borderBottom: 'none',
    },
  },
  noWorkordersPlaceholder: {
    marginTop: '16px',
    textAlign: 'center',
    fontStyle: 'italic',
  },
  media: {
    height: 0,
    paddingTop: '56.25%', // 16:9
  },
  expand: {
    transform: 'rotate(0deg)',
    marginLeft: 'auto',
    transition: theme.transitions.create('transform', {
      duration: theme.transitions.duration.shortest,
    }),
    padding: '0px',
  },
  expandOpen: {
    transform: 'rotate(180deg)',
  },
}));

type Props = {
  projectId: ?string,
};

const ProjectsPopoverQuery = graphql`
  query ProjectsPopoverQuery($projectId: ID!) {
    project(id: $projectId) {
      id
      name
      location {
        id
        name
        latitude
        longitude
      }
      workOrders {
        id
        name
        description
        ownerName
        status
        priority
        assignee
        installDate
        location {
          id
          name
          latitude
          longitude
        }
      }
    }
  }
`;

const ProjectsPopover = (props: Props) => {
  const {projectId} = props;
  const classes = useStyles();
  const router = useRouter();

  React.useEffect(() => {
    ServerLogger.info(LogEvents.PROJECTS_MAP_POPUP_OPENED, {projectId});
  }, [projectId]);

  return (
    <>
      {projectId && (
        <InventoryQueryRenderer
          query={ProjectsPopoverQuery}
          variables={{projectId}}
          render={props => {
            const {project} = props;
            const pLoc = project.location;
            const headerContent = (
              <div className={classes.cardHeaderContent}>
                <Text variant="subtitle1">{project.name}</Text>
                {pLoc && (
                  <Text variant="body2">
                    {` (${pLoc.latitude}, ${pLoc.longitude})`}
                  </Text>
                )}
              </div>
            );
            return (
              <Card className={classes.card}>
                <CardHeader
                  subheader={headerContent}
                  className={classes.cardHeader}
                />
                <CardContent className={classes.cardContent}>
                  {project.workOrders?.length ? (
                    project.workOrders.map(workOrder => (
                      <div className={classes.workOrderBlock}>
                        <WorkOrderPopover
                          onWorkOrderChanged={emptyFunction}
                          displayFullDetails={true}
                          selectedView={'status'}
                          workOrder={workOrder}
                          onWorkOrderClick={() => {
                            router.history.push(
                              `/workorders/search?workorder=${workOrder.id}`,
                            );
                          }}
                        />
                      </div>
                    ))
                  ) : (
                    <div className={classes.noWorkordersPlaceholder}>
                      <Text variant="subtitle1" color="gray">
                        {fbt(
                          'No work orders related to this project',
                          'Placeholder in ProjectsPopover card',
                        )}
                      </Text>
                    </div>
                  )}
                </CardContent>
              </Card>
            );
          }}
        />
      )}
    </>
  );
};

export default ProjectsPopover;
