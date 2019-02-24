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

import {withStyles, withTheme} from '@material-ui/core/styles';
import classNames from 'classnames';
import ExpansionPanel from '@material-ui/core/ExpansionPanel';
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails';
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary';
import KeyboardArrowDown from '@material-ui/icons/KeyboardArrowDown';
import ListItem from '@material-ui/core/ListItem';
import nullthrows from '@fbcnms/util/nullthrows';
import * as React from 'react';
import Typography from '@material-ui/core/Typography';

type TreeNode = {
  name: string,
  subtitle?: string,
  children: TreeNode[],
};

type Props = WithStyles & {
  /** The data to render as a tree view */
  tree: TreeNode[],
  /** ID of the selected node */
  selectedId: ?string,
  /** Title to be displayed **/
  title?: string,
  /** Content to the right on the title **/
  titleRightContent?: ?React.Node,
  /** Callback function fired when a tree leaf is clicked. */
  onClick: ?(any) => void,
  /** Property getter for each tree element's ID **/
  idPropertyGetter: Object => ?string,
  /** Property getter for each tree element's title **/
  titlePropertyGetter: ?(Object) => ?string,
  /** Property getter for each tree element's subtitle **/
  subtitlePropertyGetter: ?(Object) => ?string,
  /** Property getter for tree element's children **/
  childrenPropertyGetter: Object => ?Array<Object>,
  /** Property getter for tree element's children **/
  hoverRightContentGetter: ?(Object) => ?React.Node,
  /** Theme injected by withTheme **/
  theme: Theme,
};

type State = {
  expanded: any,
};

const styles = (theme: Theme) => ({
  treeContainer: {
    backgroundColor: theme.palette.common.white,
  },
  titleContainer: {
    display: 'flex',
    alignItems: 'center',
    paddingRight: '20px',
    marginBottom: '12px',
    marginLeft: '12px',
  },
  title: {
    lineHeight: '100%',
    flexGrow: 1,
    fontWeight: 'bold',
  },
  panel: {
    width: '100%',
    paddingRight: 0,
    paddingLeft: 0,
    '&:before': {
      opacity: 0,
    },
    '& $panelContent': {
      '& $headerContainer': {
        paddingRight: '16px',
      },
    },
  },
  panelExpanded: {},
  panelContent: {
    margin: '2px 0px',
    '&$panelExpanded': {
      margin: '2px 0px',
    },
  },
  panelSummary: {
    minHeight: '31px',
    paddingLeft: '20px',
    '&$panelExpanded': {
      minHeight: '31px',
    },
  },
  panelDetails: {
    padding: 0,
    display: 'block',
  },
  panelExpanded: {
    margin: 0,
    '&:before': {
      opacity: 0,
    },
  },
  childPanel: {
    '&:before': {
      opacity: 0,
    },
  },
  text: {
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'noWrap',
    maxWidth: '75vw',
  },
  heading: {
    fontSize: theme.typography.pxToRem(13),
    marginRight: '4px',
    lineHeight: '100%',
  },
  secondaryHeading: {
    marginTop: '0.2em',
    marginBottom: '0em',
    fontSize: theme.typography.pxToRem(11),
    color: theme.palette.text.secondary,
    lineHeight: '100%',
  },
  expandIcon: {
    color: theme.palette.grey[600],
  },
  selectedItem: {
    backgroundColor: theme.palette.action.selected,
  },
  treeItem: {
    '&:hover': {
      backgroundColor: theme.palette.action.hover,
    },
    '&:hover $headerContainer::before': {
      borderLeft: `2px solid ${theme.palette.primary.main}`,
    },
    '&:hover $hoverRightContent': {
      display: 'flex',
    },
  },
  headerContainer: {
    display: 'flex',
    flexGrow: 1,
    '&:before': {
      content: '""',
      marginRight: '6px',
      paddingTop: '8px',
      paddingBottom: '8px',
      borderLeft: `2px solid ${theme.palette.grey[200]}`,
    },
  },
  headerLeftContent: {
    display: 'flex',
    flexGrow: 1,
    paddingBottom: '9px',
    paddingTop: '9px',
  },
  hoverRightContent: {
    display: 'none',
    marginRight: '16px',
    alignItems: 'center',
  },
  leafRoot: {
    paddingTop: '2px',
    paddingBottom: '2px',
  },
});

/**
 * Render a tree view.
 */
class TreeView extends React.Component<Props, State> {
  state = {
    expanded: {},
  };

  static defaultProps = {
    idPropertyGetter: (node: Object) => node?.id ?? null,
    titlePropertyGetter: (node: Object) => node.name,
    subtitlePropertyGetter: null,
    childrenPropertyGetter: (node: Object) => node.children,
    hoverRightContentGetter: null,
  };

  renderNode = (node: Object, parent: ?Object, depth = 0) => {
    const {
      theme: {
        spacing: {unit},
      },
      classes,
    } = this.props;
    const spacing = unit * 1.5;
    const {selectedId} = this.props;
    const id = this.getNodeId(node);
    const isLeaf = this.isLeaf(node);
    const title = this.getNodeTitle(node);
    const subtitle = this.getNodeSubtitle(node);
    const children = this.getNodeChildren(node);
    const key = typeof id !== 'undefined' ? id : title;
    const hoverRightContent = this.getNodeHoverRightContent(node);
    const paddingLeft = depth * spacing + spacing + unit;

    const treeItemClasses = classNames({
      [classes.treeItem]: true,
      [classes.selectedItem]: selectedId === id,
    });

    const treeItemHeader = (
      <div className={classes.headerContainer}>
        <div className={classes.headerLeftContent}>
          <Typography className={classes.heading}>{title}</Typography>
          <Typography className={classes.secondaryHeading}>
            {subtitle}
          </Typography>
        </div>
        <div className={classes.hoverRightContent}>{hoverRightContent}</div>
      </div>
    );

    if (isLeaf) {
      return (
        <ListItem
          classes={{
            root: classes.leafRoot,
          }}
          className={treeItemClasses}
          disableGutters
          style={{paddingLeft}}
          key={key}
          value={title}
          onClick={() => this.props.onClick && this.props.onClick(node)}
          button>
          {treeItemHeader}
        </ListItem>
      );
    }

    const expansionPanelClasses = {
      expanded: classes.panelExpanded,
      ...(parent ? {root: classes.childPanel} : null),
    };

    return (
      <ExpansionPanel
        classes={expansionPanelClasses}
        key={key}
        elevation={0}
        className={classes.panel}
        onChange={() => this.props.onClick && this.props.onClick(node)}>
        <ExpansionPanelSummary
          className={treeItemClasses}
          classes={{
            expandIcon: classes.expandIcon,
            root: classes.panelSummary,
            expanded: classes.panelExpanded,
            content: classes.panelContent,
          }}
          style={{paddingLeft}}
          expandIcon={<KeyboardArrowDown />}
          onClick={() => {
            this.props.onClick && this.props.onClick(node);
            this.expand(node.id);
          }}>
          {treeItemHeader}
        </ExpansionPanelSummary>
        {key &&
          this.state.expanded[key] === true && (
            <ExpansionPanelDetails classes={{root: classes.panelDetails}}>
              {nullthrows(children).map(l =>
                this.renderNode(l, node, depth + 1),
              )}
            </ExpansionPanelDetails>
          )}
      </ExpansionPanel>
    );
  };

  expand(key: ?string) {
    key && this.setState({expanded: {...this.state.expanded, [key]: true}});
  }

  isLeaf(node) {
    const children = this.getNodeChildren(node);
    return !children || !children.length;
  }

  getNodeChildren(node) {
    return this._getProperty(node, this.props.childrenPropertyGetter);
  }

  getNodeTitle(node) {
    return this._getProperty(node, this.props.titlePropertyGetter);
  }

  getNodeSubtitle(node) {
    return this._getProperty(node, this.props.subtitlePropertyGetter);
  }

  getNodeId(node) {
    return this._getProperty(node, this.props.idPropertyGetter);
  }

  getNodeHoverRightContent(node) {
    return this._getProperty(node, this.props.hoverRightContentGetter);
  }

  _getProperty(node: Object, propertyGetter: ?(Object) => any) {
    return propertyGetter && propertyGetter(node);
  }

  render() {
    const {classes, tree, title, titleRightContent} = this.props;
    return (
      <div>
        <div className={classes.titleContainer}>
          <Typography variant="h6" className={classes.title}>
            {title}
          </Typography>
          {titleRightContent}
        </div>
        <div className={classes.treeContainer}>
          {tree.map(node => this.renderNode(node, null))}
        </div>
      </div>
    );
  }
}

export default withTheme()(withStyles(styles)(TreeView));
