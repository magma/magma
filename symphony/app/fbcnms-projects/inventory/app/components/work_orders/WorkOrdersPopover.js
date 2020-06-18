/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {WithStyles} from '@material-ui/core';
import type {WorkOrderIdentifier} from '../../common/WorkOrder';

import ArrowDropDownIcon from '@material-ui/icons/ArrowDropDown';
import Button from '@fbcnms/ui/components/Button.js';
import CloseIcon from '@material-ui/icons/Close';
import DescriptionIcon from '@material-ui/icons/Description';
import Popover from '@material-ui/core/Popover';
import React from 'react';
import Typography from '@material-ui/core/Typography';
import WorkOrderDetailsPaneQueryRenderer from './WorkOrderDetailsPaneQueryRenderer';
import WorkOrdersPane from './WorkOrdersPane';
import classNames from 'classnames';
import {withStyles} from '@material-ui/core/styles';

const styles = theme => ({
  popover: {
    margin: '5px 10px 0px 0px',
  },
  popoverWrap: {
    display: 'flex',
    border: `1px solid ${theme.palette.grey[100]}`,
    borderRadius: '4px',
    cursor: 'pointer',
  },
  popoverBody: {
    display: 'flex',
    alignItems: 'center',
    backgroundColor: 'white',
    padding: '0px 6px',
    height: '32px',
  },
  popoverIconWrap: {
    borderLeft: `1px solid ${theme.palette.grey[100]}`,
  },
  popoverIcon: {
    fontSize: '18px',
    color: theme.palette.primary.dark,
    cursor: 'pointer',
    '&:hover': {
      color: theme.palette.primary.main,
    },
  },
  popoverLabel: {
    color: theme.palette.dark,
    fontSize: '14px',
    lineHeight: '18px',
    marginRight: '6px',
    padding: '2px',
  },
  actionIcon: {
    color: theme.palette.dark,
  },
});

type Props = WithStyles<typeof styles> & {
  onSelect: (workOrderId: ?string) => void,
  onNavigateToWorkOrder: (workOrderId: ?string) => void,
};

type State = {
  popoverOpen: boolean,
  selectedWorkOrder: ?WorkOrderIdentifier,
};

class WorkOrdersPopover extends React.Component<Props, State> {
  _anchorRef = React.createRef();

  state = {
    popoverOpen: false,
    selectedWorkOrder: null,
  };

  setSelected = selectedWorkOrder => {
    this.setState({selectedWorkOrder}, () => {
      this.props.onSelect && this.props.onSelect(selectedWorkOrder?.id);
    });
  };

  showPopover = () => this.setState({popoverOpen: true});
  hidePopover = () => this.setState({popoverOpen: false});

  render() {
    const {popoverOpen, selectedWorkOrder} = this.state;
    const {classes} = this.props;
    return (
      <>
        {!selectedWorkOrder ? (
          <div
            className={classes.popoverWrap}
            onClick={this.showPopover}
            ref={this._anchorRef}>
            <div className={classes.popoverBody}>
              <Typography
                className={classes.popoverLabel}
                onClick={this.showPopover}>
                Work Orders
              </Typography>
              <ArrowDropDownIcon className={classes.actionIcon} />
            </div>
          </div>
        ) : (
          <div
            className={classes.popoverWrap}
            onClick={this.showPopover}
            ref={this._anchorRef}>
            <div className={classes.popoverBody}>
              <Button
                size="small"
                className={classes.actionIcon}
                onClick={() =>
                  this.props.onNavigateToWorkOrder(
                    this.state.selectedWorkOrder?.id,
                  )
                }>
                <DescriptionIcon />
              </Button>
              <Typography
                className={classes.popoverLabel}
                onClick={this.showPopover}>
                {selectedWorkOrder.name}
              </Typography>
              <ArrowDropDownIcon className={classes.actionIcon} />
            </div>
            <div
              className={classNames(
                classes.popoverBody,
                classes.popoverIconWrap,
              )}>
              <CloseIcon
                className={classes.popoverIcon}
                onMouseDown={() => this.setSelected(null)}
              />
            </div>
          </div>
        )}
        <Popover
          id="simple-popper"
          className={classes.popover}
          open={popoverOpen}
          anchorEl={this._anchorRef.current}
          onClose={this.hidePopover}
          anchorOrigin={{
            vertical: 'bottom',
            horizontal: 'left',
          }}
          transformOrigin={{
            vertical: 'top',
            horizontal: 'left',
          }}>
          {this.state.selectedWorkOrder ? (
            <WorkOrderDetailsPaneQueryRenderer
              workOrderId={this.state.selectedWorkOrder.id}
            />
          ) : (
            <WorkOrdersPane
              onSelect={this.setSelected}
              onNavigateToWorkOrder={workOrderId =>
                this.props.onNavigateToWorkOrder(workOrderId)
              }
            />
          )}
        </Popover>
      </>
    );
  }
}

export default withStyles(styles)(WorkOrdersPopover);
