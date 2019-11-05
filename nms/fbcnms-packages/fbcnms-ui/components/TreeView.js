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

// $FlowFixMe - icon exists
import ArrowRightIcon from '@material-ui/icons/ArrowRight';

import 'react-perfect-scrollbar/dist/css/styles.css';
import * as React from 'react';
import ExpansionPanel from '@material-ui/core/ExpansionPanel';
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails';
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary';
import ListItem from '@material-ui/core/ListItem';
import PerfectScrollbar from 'react-perfect-scrollbar';
import Text from './design-system/Text';
import classNames from 'classnames';
import nullthrows from '@fbcnms/util/nullthrows';
import {withStyles, withTheme} from '@material-ui/core/styles';

type TreeNode = {
  name: string,
  subtitle?: string,
  children: TreeNode[],
};

type Props = WithStyles<typeof styles> & {
  /** The data to render as a tree view */
  tree: TreeNode[],
  /** ID of the selected node */
  selectedId: ?string,
  /** Title to be displayed **/
  title?: string,
  /** Content to the right on the title **/
  dummyRootTitle?: ?string,
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
  root: {
    display: 'flex',
    flexGrow: 1,
    flexDirection: 'column',
  },
  treeContainer: {
    backgroundColor: theme.palette.common.white,
    display: 'flex',
    flexDirection: 'column',
    flexGrow: 1,
  },
  titleContainer: {
    display: 'flex',
    alignItems: 'center',
    paddingRight: theme.spacing(3),
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
  },
  panelContent: {
    margin: '0px',
    '&$panelExpanded': {
      margin: '0px',
    },
  },
  panelSummary: {
    minHeight: '31px',
    paddingLeft: '20px',
    paddingRight: '0px',
    '&$panelExpanded': {
      minHeight: '31px',
    },
    '& $headerRoot': {
      paddingRight: '0px',
    },
  },
  panelDetails: {
    padding: 0,
    display: 'block',
  },
  panelExpanded: {
    '& > $headerRoot > $arrowRightIcon': {
      transform: 'rotate(90deg)',
    },
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
  headerRoot: {
    display: 'flex',
    alignItems: 'center',
    width: '100%',
    '&:before': {
      content: '""',
      paddingTop: '16px',
      paddingBottom: '16px',
      borderLeft: `1px solid rgba(0, 0, 0, 0.086)`,
    },
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
    display: 'none',
    color: theme.palette.grey[600],
  },
  selectedItem: {
    backgroundColor: 'rgba(0, 0, 0, 0.08)',
  },
  treeItem: {
    '&:hover': {
      backgroundColor: 'rgba(0, 0, 0, 0.04)',
    },
    '&:hover $headerRoot::before': {
      borderLeft: `2px solid ${theme.palette.primary.main}`,
    },
    '&:hover $hoverRightContent': {
      display: 'flex',
    },
  },
  headerContainer: {
    display: 'flex',
    flexGrow: 1,
  },
  headerLeftContent: {
    display: 'flex',
    flexGrow: 1,
    paddingBottom: '9px',
    paddingTop: '9px',
  },
  hoverRightContent: {
    display: 'none',
    alignItems: 'center',
    marginRight: '24px',
  },
  leafRoot: {
    paddingTop: '0px',
    paddingBottom: '0px',
    '& $headerContainer': {
      marginLeft: '10px',
    },
  },
  arrowRightIcon: {
    color: 'rgba(0, 0, 0, 0.54)',
    transition: 'transform 150ms cubic-bezier(0.4, 0, 0.2, 1) 0ms',
    marginLeft: '4px',
  },
  addLocationToRootTitle: {
    color: theme.palette.text.secondary,
    flexGrow: 1,
    fontSize: theme.typography.pxToRem(13),
  },
  dummyContainer: {
    alignItems: 'center',
    display: 'flex',
    marginRight: '24px',
    width: '100%',
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

  renderDummyTitleNode(dummyNodeTitle: string) {
    const {
      classes,
      theme: {spacing},
    } = this.props;
    const unit = spacing();
    const spacingPx = unit * 2;

    const hoverRightContent = this.getNodeHoverRightContent(null);
    return (
      <ListItem
        classes={{
          root: classes.leafRoot,
        }}
        disableGutters
        style={{paddingLeft: spacingPx + unit}}
        key={'dummy_node'}>
        <div className={classes.headerRoot}>
          <div className={classes.headerContainer}>
            <div className={classes.dummyContainer}>
              <Text className={classes.addLocationToRootTitle}>
                {dummyNodeTitle}
              </Text>
              {hoverRightContent}
            </div>
          </div>
        </div>
      </ListItem>
    );
  }

  renderNode = (node: Object, parent: ?Object, depth = 0) => {
    const {
      theme: {spacing},
      classes,
    } = this.props;
    const unit = spacing();
    const spacingPx = unit * 2;
    const {selectedId} = this.props;
    const id = this.getNodeId(node);
    const isLeaf = this.isLeaf(node);
    const title = this.getNodeTitle(node);
    const subtitle = this.getNodeSubtitle(node);
    const children = this.getNodeChildren(node);
    const key = typeof id !== 'undefined' ? id : title;
    const hoverRightContent = this.getNodeHoverRightContent(node);
    const paddingLeft = depth * spacingPx + spacingPx + unit;

    const treeItemClasses = classNames({
      [classes.treeItem]: true,
      [classes.selectedItem]: selectedId === id,
    });

    const treeItemHeader = (
      <div className={classes.headerContainer}>
        <div className={classes.headerLeftContent}>
          <Text className={classes.heading}>{title}</Text>
          <Text className={classes.secondaryHeading}>{subtitle}</Text>
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
          <div className={classes.headerRoot}>{treeItemHeader}</div>
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
          onClick={() => {
            this.props.onClick && this.props.onClick(node);
            this.expand(node.id);
          }}>
          <div className={classes.headerRoot}>
            <ArrowRightIcon className={classes.arrowRightIcon} />
            {treeItemHeader}
          </div>
        </ExpansionPanelSummary>
        {key && this.state.expanded[key] === true && (
          <ExpansionPanelDetails classes={{root: classes.panelDetails}}>
            {nullthrows(children).map(l => this.renderNode(l, node, depth + 1))}
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
    const {classes, tree, title, dummyRootTitle} = this.props;
    return (
      <div className={classes.root}>
        <div className={classes.titleContainer}>
          <Text variant="h6" className={classes.title}>
            {title}
          </Text>
        </div>
        <div className={classes.treeContainer}>
          <PerfectScrollbar>
            <div>
              {dummyRootTitle !== null && dummyRootTitle !== undefined
                ? this.renderDummyTitleNode(dummyRootTitle)
                : null}
              {tree.map(node => this.renderNode(node, null))}
            </div>
          </PerfectScrollbar>
        </div>
      </div>
    );
  }
}

export default withTheme(withStyles(styles)(TreeView));
