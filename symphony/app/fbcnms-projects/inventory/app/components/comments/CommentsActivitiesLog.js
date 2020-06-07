/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {CommentsActivitiesLog_activities} from './__generated__/CommentsActivitiesLog_activities.graphql.js';
import type {CommentsActivitiesLog_comments} from './__generated__/CommentsActivitiesLog_comments.graphql.js';

import ActivityPost from './ActivityPost';
import AppContext from '@fbcnms/ui/context/AppContext';
import CommentsLogEmptyState from './CommentsLogEmptyState';
import React, {useContext, useRef} from 'react';
import TextCommentPost from './TextCommentPost';
import classNames from 'classnames';
import useVerticalScrollingEffect from '@fbcnms/ui/components/design-system/hooks/useVerticalScrollingEffect';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {createFragmentContainer, graphql} from 'react-relay';
import {makeStyles} from '@material-ui/styles';
import {withSnackbar} from 'notistack';

type Props = $ReadOnly<{|
  comments: CommentsActivitiesLog_comments,
  activities: CommentsActivitiesLog_activities,
  className?: string,
  postClassName?: string,
|}>;

const useStyles = makeStyles(() => ({
  commentsLog: {
    flexGrow: 1,
    marginBottom: '8px',
    height: '100%',
    minHeight: '80px',
    overflowY: 'auto',
    display: 'flex',
    flexDirection: 'column',
  },
  singleComment: {
    flexBasis: 'auto',
  },
}));

const CommentsActivitiesLog = (props: Props) => {
  const classes = useStyles();
  const thisElement = useRef(null);
  const {isFeatureEnabled} = useContext(AppContext);

  const {comments, activities, className, postClassName} = props;
  let objectsList;
  const activityEnabled = isFeatureEnabled('work_order_activities');
  if (!activityEnabled) {
    const hasComments = Array.isArray(comments) && comments.length > 0;
    objectsList = hasComments ? (
      comments.map(comment => (
        <div
          key={comment.id}
          className={classNames(postClassName, classes.singleComment)}>
          <TextCommentPost comment={comment} />
        </div>
      ))
    ) : (
      <CommentsLogEmptyState />
    );
  } else {
    const commentObjects = comments.map(comment => {
      return {
        createTime: comment.createTime,
        component: (
          <div
            key={comment.id}
            className={classNames(postClassName, classes.singleComment)}>
            <TextCommentPost comment={comment} />
          </div>
        ),
      };
    });

    const activityObjects = activities
      ? activities.map(activity => {
          return {
            createTime: activity.createTime,
            component: (
              <div key={activity.id}>
                <ActivityPost activity={activity} />
              </div>
            ),
          };
        })
      : [];

    objectsList = objectsList = [...commentObjects, ...activityObjects]
      .sort((a, b) => {
        return a.createTime.localeCompare(b.createTime);
      })
      .map(x => x.component);
  }
  useVerticalScrollingEffect(thisElement);

  return (
    <div
      ref={thisElement}
      className={classNames(className, classes.commentsLog)}>
      {objectsList}
    </div>
  );
};

export default withAlert(
  withSnackbar(
    createFragmentContainer(CommentsActivitiesLog, {
      comments: graphql`
        fragment CommentsActivitiesLog_comments on Comment
          @relay(plural: true) {
          id
          createTime
          ...TextCommentPost_comment
        }
      `,
      activities: graphql`
        fragment CommentsActivitiesLog_activities on Activity
          @relay(plural: true) {
          id
          createTime
          ...ActivityPost_activity
        }
      `,
    }),
  ),
);
