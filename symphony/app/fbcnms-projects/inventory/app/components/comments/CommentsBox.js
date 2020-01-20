/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import CommentsLog from './CommentsLog';
import NewCommentInput from './NewCommentInput';
import type {CommentEntity} from '../../mutations/__generated__/AddCommentMutation.graphql';
import type {CommentsBox_comments} from './__generated__/CommentsBox_comments.graphql.js';

import React from 'react';
import classNames from 'classnames';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {createFragmentContainer, graphql} from 'react-relay';
import {makeStyles} from '@material-ui/styles';
import {withSnackbar} from 'notistack';

type Props = {
  relatedEntityId: string,
  relatedEntityType: CommentEntity,
  comments: CommentsBox_comments,
  boxElementsClass?: string,
  commentsLogClass?: string,
  newCommentInputClass?: string,
};

const useStyles = makeStyles({
  container: {
    display: 'flex',
    flexDirection: 'column',
  },
});

const CommentsBox = (props: Props) => {
  const classes = useStyles();
  const {
    comments,
    relatedEntityType,
    relatedEntityId,
    boxElementsClass,
    commentsLogClass,
    newCommentInputClass,
  } = props;

  return (
    <div className={classes.container}>
      <CommentsLog
        className={classNames(boxElementsClass, commentsLogClass)}
        comments={comments}
      />
      <NewCommentInput
        className={classNames(boxElementsClass, newCommentInputClass)}
        relatedEntityId={relatedEntityId}
        relatedEntityType={relatedEntityType}
      />
    </div>
  );
};

export default withAlert(
  withSnackbar(
    createFragmentContainer(CommentsBox, {
      comments: graphql`
        fragment CommentsBox_comments on Comment @relay(plural: true) {
          ...CommentsLog_comments
        }
      `,
    }),
  ),
);
