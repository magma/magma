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
import Typography from '@material-ui/core/Typography';
import {withStyles, withTheme} from '@material-ui/core/styles';

const styles = _theme => ({
  root: {
    display: 'flex',
    borderRadius: '4px',
    backgroundColor: '#ecf3ff',
    alignItems: 'center',
    padding: '4px',
  },
  chip: {
    display: 'flex',
    alignItems: 'center',
    borderRadius: '4px',
    backgroundColor: '#ecf3ff',
    border: '1px solid white',
    padding: '0px 6px',
    height: '20px',
    marginRight: '8px',
  },
  chipDeleteIcon: {
    fontSize: '14px',
    color: '#ccd0d5',
    margin: '0px',
    padding: '0px',
    cursor: 'pointer',
    '&:hover': {
      color: '#444950',
    },
  },
  chipLabel: {
    color: '#1d2129',
    fontSize: '12px',
    lineHeight: '16px',
    marginRight: '6px',
    padding: '0px',
  },
});

export type Entry = {
  id: string,
  label: string,
};

const autoSuggestStyles = theme => ({
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
    backgroundColor: '#ecf3ff',
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
  searchEntries: Array<Entry>,
  onEntriesRequested: (searchTerm: string) => void,
  onChange?: (entries: Array<Entry>) => void,
  onBlur?: () => void,
  theme: Theme,
} & WithStyles;

type State = {
  tokens: Array<Entry>,
  searchTerm: string,
};

const BACKSPACE_KEY_CODE = 8;
const ENTER_KEY_CODE = 13;

class Tokenizer extends React.Component<Props, State> {
  state = {
    tokens: [],
    searchTerm: '',
  };

  render() {
    const {
      classes,
      theme,
      searchEntries,
      onEntriesRequested,
      onChange,
      onBlur,
    } = this.props;
    const {tokens, searchTerm} = this.state;
    const unusedSearchEntries = searchEntries.filter(entry =>
      tokens.every(token => token.id !== entry.id),
    );
    return (
      <div className={classes.root}>
        {tokens.map(token => (
          <div key={token.id} className={classes.chip}>
            <Typography className={classes.chipLabel}>{token.label}</Typography>
            <ClearIcon
              className={classes.chipDeleteIcon}
              onMouseDown={e => {
                this.setState(
                  prevState => ({
                    tokens: prevState.tokens.filter(t => t.id !== token.id),
                  }),
                  () => onChange && onChange(this.state.tokens),
                );
                e.preventDefault();
              }}
            />
          </div>
        ))}
        <Autosuggest
          suggestions={unusedSearchEntries}
          getSuggestionValue={entry => entry.label}
          onSuggestionsFetchRequested={({value}) => onEntriesRequested(value)}
          renderSuggestion={entry => (
            <div className={classes.entryRoot}>
              <div>{entry.label}</div>
            </div>
          )}
          onSuggestionSelected={(e, {suggestion}) => {
            this.setState(
              prevState => ({
                tokens: [...prevState.tokens, suggestion],
                searchTerm: '',
              }),
              () => onChange && onChange(this.state.tokens),
            );
          }}
          inputProps={{
            placeholder: '',
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

              this.setState(
                prevState => ({
                  tokens: prevState.tokens.slice(
                    0,
                    prevState.tokens.length - 1,
                  ),
                }),
                () => onChange && onChange(this.state.tokens),
              );
            },
            onChange: (_e, {newValue}) => this.setState({searchTerm: newValue}),
            onBlur: () => onBlur && onBlur(),
            autoFocus: true,
          }}
          theme={autoSuggestStyles(theme)}
          highlightFirstSuggestion={true}
        />
      </div>
    );
  }
}

export default withTheme()(withStyles(styles)(Tokenizer));
