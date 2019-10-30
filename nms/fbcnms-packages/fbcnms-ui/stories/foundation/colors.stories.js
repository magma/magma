/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import ColorBlock from './ColorBlock';
import React from 'react';
import SymphonyTheme, {BLUE, DARK} from '../../theme/symphony';
import Text from '../../components/design-system/Text';
import {STORY_CATEGORIES} from '../storybookUtils';
import {storiesOf} from '@storybook/react';

import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(_theme => ({
  root: {
    padding: '74px',
  },
  section: {
    display: 'flex',
    flexDirection: 'row',
  },
  blocksRoot: {
    display: 'flex',
    flexDirection: 'column',
    marginBottom: '96px',
  },
  title: {
    textTransform: 'uppercase',
    marginBottom: '24px',
  },
  colors: {
    display: 'flex',
    flexDirection: 'row',
    flexWrap: 'wrap',
  },
  colorBlock: {
    marginRight: '20px',
    marginBottom: '20px',
  },
}));

const ColorBlocks = (props: {
  colors: Array<{color: string, name: string, code?: string}>,
  title: string,
}) => {
  const {colors, title} = props;
  const classes = useStyles();
  return (
    <div className={classes.blocksRoot}>
      <Text className={classes.title} weight="medium">
        {title}
      </Text>
      <div className={classes.colors}>
        {colors.map(color => (
          <ColorBlock
            key={color.color}
            className={classes.colorBlock}
            color={color.color}
            name={color.name}
            code={color.code}
          />
        ))}
      </div>
    </div>
  );
};

const ColorsRoot = () => {
  const classes = useStyles();
  return (
    <div className={classes.root}>
      <div className={classes.section}>
        <ColorBlocks
          title="primary"
          colors={[{color: SymphonyTheme.palette.B600, name: 'B600'}]}
        />
        <ColorBlocks
          title="secondary"
          colors={[{color: SymphonyTheme.palette.D900, name: 'D900'}]}
        />
        <ColorBlocks
          title="white"
          colors={[{color: SymphonyTheme.palette.white, name: 'White'}]}
        />
        <ColorBlocks
          title="background"
          colors={[{color: SymphonyTheme.palette.D10, name: 'D10'}]}
        />
        <ColorBlocks
          title="disabled"
          colors={[
            {
              color: SymphonyTheme.palette.disabled,
              name: 'DIS',
              code: '#303846 38%',
            },
          ]}
        />
      </div>
      <div className={classes.section}>
        <ColorBlocks
          title="blue"
          colors={Object.keys(BLUE)
            .reverse()
            .map(colorName => ({
              color: BLUE[colorName],
              name: colorName,
            }))}
        />
      </div>
      <div className={classes.section}>
        <ColorBlocks
          title="dark"
          colors={Object.keys(DARK)
            .reverse()
            .map(colorName => ({
              color: DARK[colorName],
              name: colorName,
            }))}
        />
      </div>
      <div className={classes.section}>
        <ColorBlocks
          title="other"
          colors={[
            {
              color: SymphonyTheme.palette.R600,
              name: 'R600',
            },
            {
              color: SymphonyTheme.palette.G600,
              name: 'G600',
            },
            {
              color: SymphonyTheme.palette.Y600,
              name: 'Y600',
            },
          ]}
        />
      </div>
    </div>
  );
};

storiesOf(`${STORY_CATEGORIES.FOUNDATION}`, module).add('1.1 Palette', () => (
  <ColorsRoot />
));
