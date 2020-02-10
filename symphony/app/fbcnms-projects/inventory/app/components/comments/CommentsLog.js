/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {CommentsLog_comments} from './__generated__/CommentsLog_comments.graphql.js';

import CommentsLogEmptyState from './CommentsLogEmptyState';
import React, {useRef} from 'react';
import TextCommentPost from './TextCommentPost';
import classNames from 'classnames';
import useVerticalScrollingEffect from '@fbcnms/ui/components/design-system/hooks/useVerticalScrollingEffect';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {createFragmentContainer, graphql} from 'react-relay';
import {makeStyles} from '@material-ui/styles';
import {withSnackbar} from 'notistack';

type Props = {
  comments: CommentsLog_comments,
  className?: string,
  postClassName?: string,
};

const useStyles = makeStyles({
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
});

const CommentsLog = (props: Props) => {
  const classes = useStyles();
  const thisElement = useRef(null);
  const {comments, className, postClassName} = props;

  const hasComments = Array.isArray(comments) && comments.length > 0;

  const commentObjects = hasComments ? (
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

  useVerticalScrollingEffect(thisElement);

  return (
    <div
      ref={thisElement}
      className={classNames(className, classes.commentsLog)}>
      {commentObjects}
    </div>
  );
};

export default withAlert(
  withSnackbar(
    createFragmentContainer(CommentsLog, {
      comments: graphql`
        fragment CommentsLog_comments on Comment @relay(plural: true) {
          id
          ...TextCommentPost_comment
        }
      `,
    }),
  ),
);
