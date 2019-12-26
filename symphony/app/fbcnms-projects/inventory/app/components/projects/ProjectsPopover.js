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
import ClearIcon from '@material-ui/icons/Clear';
import IconButton from '@material-ui/core/IconButton';
import InventoryQueryRenderer from '../InventoryQueryRenderer';
import MoreVertIcon from '@material-ui/icons/MoreVert';
import WorkOrderPopover from '../work_orders/WorkOrderPopover';
import emptyFunction from '@fbcnms/util/emptyFunction';
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
    maxWidth: '400px',
    width: '25%',
    maxHeight: `calc(100% - ${cardTop} - 40px)`,
    borderRadius: '8px',
    overflow: 'hidden',
  },
  cardHeader: {
    paddingBottom: '4px',
    fontSize: '18px',
    borderBottom: `2px solid ${symphony.palette.separator}`,
  },
  cardContent: {
    overflowY: 'auto',
    overflowX: 'hidden',
    paddingTop: '8px',
  },
  workOrderBlock: {
    overflowY: 'hidden',
    padding: '8px 0px',
    borderBottom: `2px solid ${symphony.palette.separator}`,
    '&:last-child': {
      borderBottom: 'none',
      paddingBottom: '0',
    },
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
  onClearButtonClicked: () => void,
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
  const {projectId, onClearButtonClicked} = props;
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
            const headerContent = `${project.name}${project.location &&
              ' (' +
                project.location.latitude +
                ' , ' +
                project.location.longitude +
                ')'}`;
            return (
              <Card className={classes.card}>
                <CardHeader
                  action={
                    ((
                      <IconButton>
                        <MoreVertIcon />
                      </IconButton>
                    ),
                    (
                      <IconButton
                        aria-label="clear"
                        onClick={onClearButtonClicked}>
                        <ClearIcon />
                      </IconButton>
                    ))
                  }
                  subheader={headerContent}
                  className={classes.cardHeader}
                />
                <CardContent className={classes.cardContent}>
                  {project.workOrders.map(workOrder => (
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
                  ))}
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
