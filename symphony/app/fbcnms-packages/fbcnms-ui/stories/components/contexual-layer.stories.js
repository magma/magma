/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import BaseContexualLayer from '../../components/design-system/ContexualLayer/BaseContexualLayer';
import React, {useCallback, useRef, useState} from 'react';
import Text from '../../components/design-system/Text';
import {STORY_CATEGORIES} from '../storybookUtils';
import {makeStyles} from '@material-ui/styles';
import {storiesOf} from '@storybook/react';

const useStyles = makeStyles(_theme => ({
  root: {
    width: '100%',
  },
  card: {
    marginBottom: '16px',
  },
}));

type Props = {element: ?Element};

const BelowContexualContainer = ({element}: Props) => {
  console.log('element', element);
  if (element == null) {
    return null;
  }
  return (
    <BaseContexualLayer context={element} position="below">
      <div>
        <Text variant="body2">Below the input, with the same width</Text>
      </div>
    </BaseContexualLayer>
  );
};

const AboveContexualContainer = ({element}: Props) => {
  if (element == null) {
    return null;
  }
  return (
    <BaseContexualLayer context={element} position="above">
      <div>
        <Text variant="body2">
          Above the input, with the same width. Amazing.
        </Text>
      </div>
    </BaseContexualLayer>
  );
};

const ContextualLayersRoot = () => {
  const classes = useStyles();
  const [isVisible, setIsVisible] = useState(false);
  const contextRef = useRef<?HTMLElement>(null);
  const refCallback = useCallback((element: ?HTMLElement) => {
    contextRef.current = element;
    setIsVisible(true);
  }, []);

  return (
    <div className={classes.root}>
      <input
        ref={refCallback}
        style={{position: 'absolute', left: 400, top: 400}}
      />
      {isVisible ? (
        <BelowContexualContainer element={contextRef.current} />
      ) : null}
      {isVisible ? (
        <AboveContexualContainer element={contextRef.current} />
      ) : null}
    </div>
  );
};

storiesOf(`${STORY_CATEGORIES.COMPONENTS}`, module).add(
  'Contextual Layer',
  () => <ContextualLayersRoot />,
);
