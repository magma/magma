/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import * as React from 'react';
import Text from '../Text';
import classNames from 'classnames';
import symphony from '../../../theme/symphony';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles({
  root: {
    backgroundColor: symphony.palette.white,
    borderRadius: '4px',
    boxShadow: symphony.shadows.DP1,
    display: 'inline-flex',
    flexDirection: 'row',
    alignItems: 'center',
    height: '32px',
  },
  button: {
    padding: '6px',
    cursor: 'pointer',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    '&:first-child': {
      paddingLeft: '8px',
    },
    '&:last-child': {
      paddingRight: '8px',
    },
    '&:hover, &:hover $buttonText': {
      color: symphony.palette.primary,
    },
  },
  buttonText: {
    display: 'flex',
    alignItems: 'center',
  },
  iconButton: {
    paddingBottom: '4px',
    paddingTop: '4px',
  },
  separator: {
    width: '1px',
    height: '20px',
    backgroundColor: 'rgba(139, 152, 173, 0.3)',
  },
  selected: {
    color: symphony.palette.primary,
    '& $buttonText': {
      color: symphony.palette.primary,
    },
  },
});

export type ButtonItem = {
  id: string,
  item: React.Node,
};

type Props = {
  buttons: Array<ButtonItem>,
  onItemClicked: (id: string) => void,
  selectedButtonId?: ?string,
  className?: string,
};

const ToggleButton = (props: Props) => {
  const {buttons, selectedButtonId, onItemClicked, className} = props;
  const classes = useStyles();
  return (
    <div className={classNames(classes.root, className)}>
      {buttons.map((button, i) => (
        <>
          <div
            className={classNames(classes.button, {
              [classes.iconButton]: typeof button.item !== 'string',
              [classes.selected]: button.id === selectedButtonId,
            })}
            key={button.id}
            onClick={() => onItemClicked(button.id)}>
            <Text
              className={classes.buttonText}
              variant="caption"
              weight="bold">
              {button.item}
            </Text>
          </div>
          {i !== buttons.length - 1 && <div className={classes.separator} />}
        </>
      ))}
    </div>
  );
};

export default ToggleButton;
