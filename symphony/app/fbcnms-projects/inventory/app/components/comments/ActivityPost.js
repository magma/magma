/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ActivityPost_activity} from './__generated__/ActivityPost_activity.graphql.js';

import ActivityCommentsIcon from './ActivityCommentsIcon';
import DateTimeFormat from '../../common/DateTimeFormat.js';
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import classNames from 'classnames';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import {createFragmentContainer, graphql} from 'react-relay';
import {makeStyles} from '@material-ui/styles';

type Props = $ReadOnly<{|
  activity: ActivityPost_activity,
|}>;

const useStyles = makeStyles(() => ({
  textActivityPost: {
    minHeight: '20px',
    padding: '4px 4px 12px 0px',
    display: 'flex',
    flexDirection: 'row',
  },
  activityBody: {
    flexGrow: 1,
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'start',
  },
  activityAuthor: {
    fontWeight: 'bold',
    paddingRight: '5px',
  },
  activityTime: {
    color: symphony.palette.D300,
  },
  nameText: {
    fontWeight: 'bold',
  },
  boldName: {
    fontWeight: 'bold',
    textTransform: 'capitalize',
  },
  fieldName: {
    textTransform: 'capitalize',
  },
}));

const ActivityPost = (props: Props) => {
  const classes = useStyles();
  const {activity} = props;

  const shouldCapitalizeValue = () => {
    return (
      activity.changedField === 'STATUS' || activity.changedField === 'PRIORITY'
    );
  };

  const genActivityValueComponent = (val: string) => {
    return (
      <span
        className={classNames({
          [classes.nameText]: !shouldCapitalizeValue(),
          [classes.boldName]: shouldCapitalizeValue(),
        })}>
        {val}
      </span>
    );
  };

  const genActivityMessage = () => {
    if (activity.changedField === 'CREATION_DATE') {
      return (
        <span>
          <fbt desc="">created this work order</fbt>
        </span>
      );
    }
    let oldVal = (activity.oldValue ?? '').toLowerCase();
    let newVal = (activity.newValue ?? '').toLowerCase();
    const oldValNode = activity.oldRelatedNode;
    if (oldValNode && oldValNode?.__typename === 'User') {
      oldVal = oldValNode.email;
    }
    const newValNode = activity.newRelatedNode;
    if (newValNode && newValNode?.__typename === 'User') {
      newVal = newValNode.email;
    }
    if (oldVal === '') {
      return (
        <span>
          <fbt desc="">
            set the{' '}
            <fbt:param name="changed field">
              <span className={classes.fieldName}>
                {activity.changedField.toLowerCase()}
              </span>
            </fbt:param>
            to be{' '}
            <fbt:param name="new value">
              {genActivityValueComponent(newVal)}
            </fbt:param>
          </fbt>
        </span>
      );
    }
    if (newVal === '') {
      return (
        <span>
          <fbt desc="">
            removed{' '}
            <fbt:param name="changed field">
              <span className={classes.fieldName}>
                {activity.changedField.toLowerCase()}
              </span>
            </fbt:param>
            value
          </fbt>
        </span>
      );
    }
    return (
      <span>
        <fbt desc="">
          changed the{' '}
          <fbt:param name="changed field">
            <span className={classes.fieldName}>
              {activity.changedField.toLowerCase()}
            </span>
          </fbt:param>
          from{' '}
          <fbt:param name="old value">
            {genActivityValueComponent(oldVal)}
          </fbt:param>
          to{' '}
          <fbt:param name="new value">
            {genActivityValueComponent(newVal)}
          </fbt:param>
        </fbt>
      </span>
    );
  };

  return (
    <div className={classes.textActivityPost}>
      <ActivityCommentsIcon field={activity.changedField} />
      <div className={classes.activityBody}>
        <Text variant="body2">
          <span className={classes.activityAuthor}>
            {activity.author?.email}
          </span>
          {genActivityMessage()}
        </Text>
        <Text
          color="light"
          variant="subtitle2"
          className={classes.activityTime}>
          {DateTimeFormat.commentTime(activity.createTime)}
        </Text>
      </div>
    </div>
  );
};

export default createFragmentContainer(ActivityPost, {
  activity: graphql`
    fragment ActivityPost_activity on Activity {
      id
      author {
        email
      }
      isCreate
      changedField
      newRelatedNode {
        __typename
        ... on User {
          id
          email
        }
      }
      oldRelatedNode {
        __typename
        ... on User {
          id
          email
        }
      }
      oldValue
      newValue
      createTime
    }
  `,
});
