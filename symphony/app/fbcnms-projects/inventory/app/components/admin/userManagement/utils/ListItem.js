/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ButtonSkin} from '@fbcnms/ui/components/design-system/Button';

import * as React from 'react';
import Button from '@fbcnms/ui/components/design-system/Button';
import classNames from 'classnames';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useState} from 'react';

const useStyles = makeStyles(() => ({
  root: {
    padding: '8px 0',
    display: 'flex',
    overflow: 'hidden',
    alignItems: 'center',
    '&:not(:hover):not($alwaysShowToggleButton) $toggleButton': {
      display: 'none',
    },
    flexShrink: 0,
  },
  alwaysShowToggleButton: {},
  userDetails: {
    flexBasis: '10px',
    flexGrow: 1,
    flexShrink: 1,
  },
  toggleButton: {
    '&:hover $shownContent': {
      '&$hoverContent': {
        maxHeight: 'unset',
        visibility: 'visible',
      },
    },
    '&:not(:hover) $shownContent': {
      '&$regularContent': {
        maxHeight: 'unset',
        visibility: 'visible',
      },
    },
  },
  toggleButtonContentContainer: {
    display: 'flex',
    flexDirection: 'column',
  },
  toggleButtonContent: {
    maxHeight: 0,
    visibility: 'hidden',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
  },
  shownContent: {
    '&$onProcessContent': {
      maxHeight: 'unset',
      visibility: 'visible',
    },
  },
  regularContent: {},
  hoverContent: {},
  onProcessContent: {},
  togglingProcess: {},
}));

export const TOGGLE_BUTTON_DISPLAY = {
  always: 'always',
  onHover: 'onHover',
};

export type ToggleButtonDisplay = $Keys<typeof TOGGLE_BUTTON_DISPLAY>;

type ToggleButtonContent = {
  regularContent: React.Node,
  hoverContent?: React.Node,
  onProcessContent: React.Node,
  skin?: ?ButtonSkin,
};

export type AssigenmentButtonProp = $ReadOnly<{|
  displayVariants?: ?ToggleButtonDisplay,
  onContent: ToggleButtonContent,
  offContent: ToggleButtonContent,
  isOn: boolean,
  onToggleClicked: () => Promise<void> | void,
|}>;

type Props = $ReadOnly<{|
  className?: ?string,
  toggleButton?: ?AssigenmentButtonProp,
  children: React.Node,
|}>;

export default function ListItem(props: Props) {
  const {children, toggleButton, className} = props;
  const classes = useStyles();

  const [isProcessed, setIsProcessed] = useState(false);

  const togglePressed = useCallback(() => {
    if (toggleButton == null) {
      return;
    }

    setIsProcessed(true);

    Promise.resolve(toggleButton.onToggleClicked()).finally(() =>
      setIsProcessed(false),
    );
  }, [toggleButton]);

  return (
    <div
      className={classNames(classes.root, className, {
        [classes.alwaysShowToggleButton]:
          toggleButton?.displayVariants == TOGGLE_BUTTON_DISPLAY.always ||
          isProcessed,
      })}>
      {children}
      {toggleButton == null ? null : (
        <Button
          className={classNames(classes.toggleButton, {
            [classes.togglingProcess]: isProcessed,
          })}
          disabled={isProcessed}
          onClick={() => togglePressed()}
          skin={
            toggleButton.isOn
              ? toggleButton.onContent.skin || undefined
              : toggleButton.offContent.skin || undefined
          }>
          <div className={classes.toggleButtonContentContainer}>
            <div
              className={classNames(
                classes.toggleButtonContent,
                classes.regularContent,
                {
                  [classes.hoverContent]:
                    toggleButton.offContent.hoverContent == null,
                  [classes.shownContent]: !toggleButton.isOn && !isProcessed,
                },
              )}>
              {toggleButton.offContent.regularContent}
            </div>
            {toggleButton.offContent.hoverContent != null && (
              <div
                className={classNames(
                  classes.toggleButtonContent,
                  classes.hoverContent,
                  {
                    [classes.shownContent]: !toggleButton.isOn && !isProcessed,
                  },
                )}>
                {toggleButton.offContent.hoverContent}
              </div>
            )}
            <div
              className={classNames(
                classes.toggleButtonContent,
                classes.onProcessContent,
                {
                  [classes.shownContent]: !toggleButton.isOn && isProcessed,
                },
              )}>
              {toggleButton.offContent.onProcessContent}
            </div>
            <div
              className={classNames(
                classes.toggleButtonContent,
                classes.regularContent,
                {
                  [classes.hoverContent]:
                    toggleButton.onContent.hoverContent == null,
                  [classes.shownContent]: toggleButton.isOn && !isProcessed,
                },
              )}>
              {toggleButton.onContent.regularContent}
            </div>
            {toggleButton.onContent.hoverContent && (
              <div
                className={classNames(
                  classes.toggleButtonContent,
                  classes.hoverContent,
                  {
                    [classes.shownContent]: toggleButton.isOn && !isProcessed,
                  },
                )}>
                {toggleButton.onContent.hoverContent}
              </div>
            )}
            <div
              className={classNames(
                classes.toggleButtonContent,
                classes.onProcessContent,
                {
                  [classes.shownContent]: toggleButton.isOn && isProcessed,
                },
              )}>
              {toggleButton.onContent.onProcessContent}
            </div>
          </div>
        </Button>
      )}
    </div>
  );
}
