/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {Theme, WithStyles} from '@material-ui/core';

import * as React from 'react';
import Autosuggest from 'react-autosuggest';
import ClearIcon from '@material-ui/icons/Clear';
import Text from './design-system/Text';
import classNames from 'classnames';
import {blue05, gray10, gray11, gray12, gray9} from '@fbcnms/ui/theme/colors';
import {withStyles, withTheme} from '@material-ui/core/styles';

const styles = {
  root: {
    display: 'flex',
    borderRadius: '4px',
    alignItems: 'center',
    padding: '4px',
  },
  searchRoot: {
    backgroundColor: blue05,
  },
  enumRoot: {
    border: `1.43px solid ${gray9}`,
    width: '280px',
    '&:hover': {
      border: `1.43px solid ${gray10}`,
    },
    flexWrap: 'wrap',
    flexGrow: 1,
    flexDirection: 'row',
    minHeight: '32px',
  },
  chip: {
    display: 'flex',
    alignItems: 'center',
    borderRadius: '4px',
    backgroundColor: blue05,
    border: '1px solid white',
    padding: '0px 6px',
    height: '20px',
    marginRight: '8px',
  },
  enumChip: {
    marginRight: '4px',
    marginBottom: '4px',
  },
  chipDeleteIcon: {
    fontSize: '14px',
    color: gray11,
    margin: '0px',
    padding: '0px',
    cursor: 'pointer',
    '&:hover': {
      color: '#444950',
    },
  },
  enumChipDeleteIcon: {
    color: gray12,
  },
  chipLabel: {
    color: '#1d2129',
    fontSize: '12px',
    lineHeight: '16px',
    marginRight: '6px',
    padding: '0px',
  },
};

type SearchSource = 'Options' | 'UserInput';

export type Entry = {
  +id: string,
  +label: string,
};

const autoSuggestStyles = (theme: Theme, searchSource: SearchSource) => ({
  container: {
    position: 'relative',
    flexGrow: 1,
  },
  input: {
    width: '100%',
    fontFamily: theme.typography.subtitle1.fontFamily,
    fontWeight: theme.typography.subtitle1.fontWeight,
    fontSize: '14px',
    color: theme.typography.subtitle1.color,
    border: 'none',
    backgroundColor: searchSource == 'Options' ? blue05 : 'inherit',
  },
  inputFocused: {
    outlineWidth: 0,
  },
  suggestionsContainer: {
    display: 'none',
  },
  suggestionsContainerOpen: {
    backgroundColor: theme.palette.common.white,
    borderRadius: '2px',
    display: 'block',
    fontFamily: theme.typography.subtitle1.fontFamily,
    fontSize: theme.typography.subtitle1.fontSize,
    fontWeight: theme.typography.subtitle1.fontWeight,
    position: 'absolute',
    top: '42px',
    width: '100%',
    zIndex: '2',
    boxShadow: theme.shadows[2],
  },
  suggestionsList: {
    margin: 0,
    padding: 0,
    listStyleType: 'none',
  },
  suggestion: {
    color: theme.palette.grey[900],
    fontSize: '12px',
    lineHeight: '16px',
    cursor: 'pointer',
    padding: '6px 12px',
  },
  suggestionHighlighted: {
    backgroundColor: theme.palette.grey[100],
  },
});

type Props = {
  searchSource: SearchSource,
  tokens: Array<Entry>,
  searchEntries?: Array<Entry>,
  placeholder: string,
  onEntriesRequested: (searchTerm: string) => void,
  onChange: (entries: Array<Entry>) => void,
  onBlur?: () => void,
  theme: Theme,
} & WithStyles<typeof styles>;

type State = {
  searchTerm: string,
};

const BACKSPACE_KEY_CODE = 8;
const ENTER_KEY_CODE = 13;

class Tokenizer extends React.Component<Props, State> {
  static defaultProps = {
    placeholder: '',
  };

  state = {
    searchTerm: '',
  };

  render() {
    const {
      classes,
      theme,
      searchSource,
      tokens,
      searchEntries,
      onEntriesRequested,
      onChange,
      onBlur,
      placeholder,
    } = this.props;
    const {searchTerm} = this.state;
    const entries =
      searchSource === 'Options' && searchEntries
        ? searchEntries
        : [{id: searchTerm, label: searchTerm}];
    const unusedSearchEntries = entries.filter(entry =>
      tokens.every(token => token.id !== entry.id),
    );
    return (
      <div
        className={classNames({
          [classes.root]: true,
          [classes.enumRoot]: searchSource === 'UserInput',
          [classes.searchRoot]: searchSource === 'Options',
        })}>
        {tokens.map(token => (
          <div
            key={token.id}
            className={classNames({
              [classes.chip]: true,
              [classes.enumChip]: searchSource === 'UserInput',
            })}>
            <Text className={classes.chipLabel}>{token.label}</Text>
            <ClearIcon
              className={classNames({
                [classes.chipDeleteIcon]: true,
                [classes.enumChipDeleteIcon]: searchSource === 'UserInput',
              })}
              onMouseDown={e => {
                onChange(tokens.filter(t => t.id !== token.id));
                e.preventDefault();
              }}
            />
          </div>
        ))}
        <Autosuggest
          suggestions={unusedSearchEntries}
          getSuggestionValue={entry => entry.label}
          onSuggestionsFetchRequested={({value}) => onEntriesRequested(value)}
          renderSuggestion={entry => <div>{entry.label}</div>}
          onSuggestionSelected={(e, {suggestion}) => {
            this.setState({
              searchTerm: '',
            });
            onChange([...tokens, suggestion]);
          }}
          inputProps={{
            placeholder: placeholder,
            value: searchTerm,
            onKeyDown: e => {
              if (
                this.state.searchTerm === '' &&
                e.keyCode === ENTER_KEY_CODE
              ) {
                onBlur && onBlur();
                return;
              }

              if (
                this.state.searchTerm !== '' ||
                e.keyCode !== BACKSPACE_KEY_CODE
              ) {
                return;
              }
              onChange(tokens.slice(0, tokens.length - 1));
            },
            onChange: (_e, {newValue}) => this.setState({searchTerm: newValue}),
            onBlur: () => onBlur && onBlur(),
            autoFocus: true,
          }}
          theme={autoSuggestStyles(theme, searchSource)}
          highlightFirstSuggestion={true}
        />
      </div>
    );
  }
}

export default withTheme(withStyles(styles)(Tokenizer));
