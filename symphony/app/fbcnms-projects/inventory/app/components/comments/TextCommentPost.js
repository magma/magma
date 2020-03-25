/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {TextCommentPost_comment} from './__generated__/TextCommentPost_comment.graphql.js';

import ChatBubbleOutlineIcon from '@material-ui/icons/ChatBubbleOutline';
import DateTimeFormat from '../../common/DateTimeFormat.js';
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';

import symphony from '@fbcnms/ui/theme/symphony';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {createFragmentContainer, graphql} from 'react-relay';
import {makeStyles} from '@material-ui/styles';
import {withSnackbar} from 'notistack';

type Props = {
  comment: TextCommentPost_comment,
};

const useStyles = makeStyles(() => ({
  textCommentPost: {
    minHeight: '24px',
    padding: '8px 4px 8px 0px',
    display: 'flex',
    flexDirection: 'row',
  },
  commentIndicator: {
    padding: '8px 12px 0px 0px',
  },
  commentTypeIcon: {
    fontSize: '16px',
    color: symphony.palette.D300,
  },
  commentBody: {
    flexGrow: 1,
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'start',
  },
  commentAuthor: {
    fontWeight: 'bold',
  },
  commentContent: {
    flexGrow: 1,
    backgroundColor: symphony.palette.D10,
    borderRadius: '4px',
    padding: '4px 8px',
  },
  commentTime: {
    paddingTop: '4px',
    color: symphony.palette.D300,
  },
}));

const TextCommentPost = (props: Props) => {
  const classes = useStyles();
  const {comment} = props;

  return (
    <div className={classes.textCommentPost}>
      <div className={classes.commentIndicator}>
        <ChatBubbleOutlineIcon className={classes.commentTypeIcon} />
      </div>
      <div className={classes.commentBody}>
        <Text variant="body2" className={classes.commentContent}>
          <span className={classes.commentAuthor}>
            {comment.authorName + ' '}
          </span>
          <span>{comment.text}</span>
        </Text>
        <Text color="light" variant="subtitle2" className={classes.commentTime}>
          {DateTimeFormat.commentTime(comment.createTime)}
        </Text>
      </div>
    </div>
  );
};

export default withAlert(
  withSnackbar(
    createFragmentContainer(TextCommentPost, {
      comment: graphql`
        fragment TextCommentPost_comment on Comment {
          id
          authorName
          text
          createTime
        }
      `,
    }),
  ),
);
