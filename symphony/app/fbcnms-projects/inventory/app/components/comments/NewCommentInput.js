/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {
  AddCommentMutationVariables,
  CommentEntity,
} from '../../mutations/__generated__/AddCommentMutation.graphql';

import AddCommentMutation from '../../mutations/AddCommentMutation';
import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import React, {useState} from 'react';
import TextInput from '@fbcnms/ui/components/design-system/Input/TextInput';
import classNames from 'classnames';
import {makeStyles} from '@material-ui/styles';

type Props = {
  relatedEntityId: string,
  relatedEntityType: CommentEntity,
  className?: string,
};

const useStyles = makeStyles(() => ({
  newCommentBox: {
    width: '100%',
    padding: '8px 0px',
  },
  newCommentInput: {
    width: '100%',
  },
}));

const onAddComment = (entityId, entityType: CommentEntity, commentText) => {
  const variables: AddCommentMutationVariables = {
    input: {
      entityType: entityType,
      id: entityId,
      text: commentText,
    },
  };

  const updater = store => {
    const newComment = store.getRootField('addComment');
    const entityProxy = store.get(entityId);

    const linkedComments = entityProxy.getLinkedRecords('comments') || [];
    entityProxy.setLinkedRecords([...linkedComments, newComment], 'comments');
  };

  AddCommentMutation(variables, null, updater);
};

const NewCommentInput = (props: Props) => {
  const classes = useStyles();
  const {relatedEntityType, relatedEntityId, className} = props;
  const [composedCommentText, setComposedComment] = useState('');

  const onSubmit = () => {
    onAddComment(
      relatedEntityId,
      relatedEntityType,
      composedCommentText.trim(),
    );
    setComposedComment('');
  };

  return (
    <div className={classNames(className, classes.newCommentBox)}>
      <FormField>
        <TextInput
          className={classes.newCommentInput}
          type="string"
          placeholder="Write a comment..."
          hint="Press Enter to send"
          onChange={({target}) => setComposedComment(target.value)}
          onEnterPressed={onSubmit}
          value={composedCommentText}
        />
      </FormField>
    </div>
  );
};

export default NewCommentInput;
