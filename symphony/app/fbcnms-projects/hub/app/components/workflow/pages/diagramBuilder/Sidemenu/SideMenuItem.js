/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import React from 'react';
import {Card, Menu} from 'semantic-ui-react';

const SideMenuItem = props => {
  let description = null;

  if (props.model.description) {
    description = props.model.description.split('-')[0];
  }

  return (
    <Menu.Item
      title={description}
      color="blue"
      fluid
      as={Card}
      draggable={true}
      onDragStart={e => {
        e.dataTransfer.setData(
          'storm-diagram-node',
          JSON.stringify(props.model),
        );
      }}
      style={{minHeight: '50px', cursor: 'grab', backgroundColor: 'white'}}>
      <div style={{float: 'left', maxWidth: '90%'}}>{props.name}</div>
      <div
        style={{
          float: 'right',
          marginTop: '8px',
          color: 'grey',
          opacity: '30%',
        }}>
        <i className="fas fa-grip-vertical" />
      </div>
    </Menu.Item>
  );
};

export default SideMenuItem;
