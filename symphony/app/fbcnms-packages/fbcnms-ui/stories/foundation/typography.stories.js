/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {TextVariant} from '../../theme/symphony';

import React from 'react';
import SymphonyTheme from '../../theme/symphony';
import Text from '../../components/design-system/Text';
import {STORY_CATEGORIES} from '../storybookUtils';
import {makeStyles} from '@material-ui/styles';
import {storiesOf} from '@storybook/react';

const useStyles = makeStyles(({symphony}) => {
  return {
    root: {
      padding: '52px',
    },
    textProps: {
      textAlign: 'right',
      marginBottom: '5px',
    },
    textContainer: {
      borderBottom: `1px solid ${symphony.palette.separator}`,
      paddingBottom: '26px',
      paddingTop: '25px',
    },
    capitalize: {
      textTransform: 'capitalize',
    },
  };
});

const TypographyBlock = (props: {variant: TextVariant}) => {
  const {variant} = props;
  const classes = useStyles();

  const variantObject = SymphonyTheme.typography[variant];
  const getVariantName = (variant: string) => {
    if (variantObject.textTransform === 'uppercase') {
      return variant;
    }

    return <span className={classes.capitalize}>{variant}</span>;
  };

  const getHumanReadableFontWeight = (variant: TextVariant) => {
    switch (SymphonyTheme.typography[variant].fontWeight) {
      case 300:
        return 'Light';
      case 400:
        return 'Regular';
      case 500:
        return 'Medium';
      case 600:
        return 'Bold';
      default:
        return 'Regular';
    }
  };

  return (
    <div className={classes.textContainer}>
      <div className={classes.textProps}>
        <Text variant="h6">
          {variantObject.fontSize.slice(0, -2)} / LH{' '}
          {variantObject.lineHeight === 'normal'
            ? 'auto'
            : Math.round(
                Number(variantObject.fontSize.slice(0, -2)) *
                  variantObject.lineHeight,
              )}{' '}
          / LS{' '}
          {variantObject.letterSpacing === 'normal'
            ? 'auto'
            : variantObject.letterSpacing.slice(0, -2)}
        </Text>
      </div>
      <div>
        <Text variant={variant}>
          {getVariantName(variant)} / Roboto{' '}
          {getHumanReadableFontWeight(variant)}
        </Text>
      </div>
    </div>
  );
};

const TypographyRoot = () => {
  const classes = useStyles();
  return (
    <div className={classes.root}>
      {Object.keys(SymphonyTheme.typography).map(variant => (
        <TypographyBlock key={variant} variant={variant} />
      ))}
    </div>
  );
};

storiesOf(`${STORY_CATEGORIES.FOUNDATION}`, module).add(
  '1.3 Typography',
  () => <TypographyRoot />,
);
