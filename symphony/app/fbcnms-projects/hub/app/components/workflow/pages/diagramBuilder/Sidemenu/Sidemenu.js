/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import React, {useCallback, useEffect, useState} from 'react';
import {
  Divider,
  Dropdown,
  Grid,
  Icon,
  Input,
  Menu,
  Popup,
  Sidebar,
} from 'semantic-ui-react';

import './Sidemenu.css';
import SideMenuItem from './SideMenuItem';
import {getTaskInputsRegex, getWfInputsRegex, hash} from '../builder-utils';

const sub_workflow = wf => ({
  name: wf.name,
  taskReferenceName: wf.name.toLowerCase().trim() + '_ref_' + hash(),
  inputParameters: getWfInputsRegex(wf),
  type: 'SUB_WORKFLOW',
  subWorkflowParam: {
    name: wf.name,
    version: wf.version,
  },
  optional: false,
  startDelay: 0,
});

const sub_task = t => ({
  name: t.name,
  taskReferenceName: t.name.toLowerCase().trim() + '_ref_' + hash(),
  inputParameters: getTaskInputsRegex(t),
  type: 'SIMPLE',
  optional: false,
  startDelay: 0,
});

const favorites = props => {
  return props.workflows
    .map((wf, i) => {
      const wfObject = sub_workflow(wf);
      if (wf.description && wf.description.includes('FAVOURITE')) {
        return (
          <SideMenuItem
            key={`wf${i}`}
            model={{
              type: 'default',
              wfObject,
              name: wf.name,
              description: wf.hasOwnProperty('description')
                ? wf.description
                : '',
            }}
            name={wf.name}
          />
        );
      }
    })
    .filter(item => item !== undefined);
};

const workflows = props => {
  return props.workflows.map((wf, i) => {
    const wfObject = sub_workflow(wf);
    return (
      <SideMenuItem
        key={`wf${i}`}
        model={{
          type: 'default',
          wfObject,
          name: wf.name,
          description: wf.hasOwnProperty('description') ? wf.description : '',
        }}
        name={wf.name}
      />
    );
  });
};

const tasks = props => {
  return props.tasks.map((task, i) => {
    const wfObject = sub_task(task);
    return (
      <SideMenuItem
        key={`wf${i}`}
        model={{
          type: 'default',
          wfObject,
          name: task.name,
          description: task.hasOwnProperty('description')
            ? task.description
            : '',
        }}
        name={task.name}
      />
    );
  });
};

const custom = (props, custom) => {
  return props.workflows
    .map((wf, i) => {
      const wfObject = sub_workflow(wf);
      if (wf.description && wf.description.includes(custom)) {
        return (
          <SideMenuItem
            key={`wf${i}`}
            model={{
              type: 'default',
              wfObject,
              name: wf.name,
              description: wf.hasOwnProperty('description')
                ? wf.description
                : '',
            }}
            name={wf.name}
          />
        );
      }
    })
    .filter(item => item !== undefined);
};

const getCustoms = props => {
  return [
    ...new Set(
      props.workflows
        .map(wf => {
          if (wf.hasOwnProperty('description')) {
            if (wf.description.match(/custom(?:\w+)?\b/gim)) {
              return wf.description.match(/custom(?:\w+)?\b/gim);
            }
          }
        })
        .flat()
        .filter(item => item !== undefined),
    ),
  ];
};

const Sidemenu = props => {
  const [visible, setVisible] = useState(false);
  const [expanded, setExpanded] = useState(false);
  const [content, setContent] = useState([]);
  const [customs, setCustoms] = useState([]);
  const [open, setOpen] = useState();

  const getContent = useCallback(
    which => {
      switch (which) {
        case 'Workflows':
          setContent(workflows(props));
          break;
        case 'Favorites':
          setContent(favorites(props));
          break;
        case 'Tasks':
          setContent(tasks(props));
          break;
        default:
          setContent(custom(props, which));
          break;
      }
    },
    [props],
  );
  
  useEffect(() => {
    setTimeout(() => setVisible(true), 1000);
    if (customs.length < 1) {
      setCustoms(getCustoms(props));
    }
    getContent(open);
  }, [props, customs.length, getContent, open]);

  const handleOpen = which => {
    if (which === open) {
      setOpen();
      return setExpanded(false);
    }
    getContent(which);
    props.openCard(which);
    setOpen(which);
    setExpanded(true);
  };

  const shortcutsInfo = () => {
    return (
      <Grid columns="equal" style={{width: '350px'}}>
        <Grid.Column style={{textAlign: 'right'}}>
          <p>
            Save <kbd>Ctrl</kbd>+<kbd>S</kbd>
          </p>
          <p>
            Zoom In <kbd>Ctrl</kbd>+<kbd>+</kbd>
          </p>
          <p>
            Zoom Out <kbd>Ctrl</kbd>+<kbd>-</kbd>
          </p>
          <p>
            Expand <kbd>Ctrl</kbd>+<kbd>X</kbd>
          </p>
        </Grid.Column>
        <Grid.Column style={{textAlign: 'right'}}>
          <p>
            Delete <kbd>LMB</kbd>+<kbd>Delete</kbd>
          </p>
          <p>
            Lock <kbd>Ctrl</kbd>+<kbd>L</kbd>
          </p>
          <p>
            Execute <kbd>Alt</kbd>+<kbd>Enter</kbd>
          </p>
        </Grid.Column>
      </Grid>
    );
  };

  const handleLabelChange = (e, {_, value}) => {
    props.updateQuery(null, value);
  };

  return (
    <div style={{zIndex: 11}}>
      <Sidebar
        id="sidebar-primary"
        as={Menu}
        animation="overlay"
        onHide={() => setVisible(true)}
        visible={visible}
        vertical
        icon>
        <Menu.Item
          as="a"
          title="Workflows"
          active={open === 'Workflows'}
          onClick={() => handleOpen('Workflows')}>
          <Icon name="folder open" />
        </Menu.Item>
        <Menu.Item
          as="a"
          title="Tasks"
          active={open === 'Tasks'}
          onClick={() => handleOpen('Tasks')}>
          <Icon name="tasks" />
        </Menu.Item>
        <Menu.Item
          as="a"
          title="Favorites"
          active={open === 'Favorites'}
          onClick={() => handleOpen('Favorites')}>
          <Icon name="favorite" />
        </Menu.Item>
        {customs.map((custom, i) => (
          <Menu.Item
            as="a"
            title={`${custom}`}
            active={open === custom}
            onClick={() => handleOpen(custom)}>
            {i + 1}
          </Menu.Item>
        ))}
        <div className="bottom">
          <Popup
            style={{transform: 'translate3d(60px, 82vh, 0px)'}}
            basic
            content={shortcutsInfo}
            on="click"
            trigger={
              <Menu.Item as="a" title="Shortcuts">
                <Icon name="keyboard" />
              </Menu.Item>
            }
          />
          <Menu.Item
            as="a"
            title="Help"
            onClick={() =>
              window.open(
                'https://docs.frinx.io/frinx-machine/workflow-builder/workflow-builder.html',
                '_blank',
              )
            }>
            <Icon name="help circle" />
          </Menu.Item>
          <Menu.Item>
            <small>{process.env.REACT_APP_VERSION}</small>
          </Menu.Item>
        </div>
      </Sidebar>

      <Sidebar
        id="sidebar-secondary"
        as={Menu}
        animation="overlay"
        direction="left"
        vertical
        visible={expanded}>
        <div className="sidebar-header">
          <h3>{open}</h3>
          <Input
            fluid
            onChange={e => props.updateQuery(e.target.value, null)}
            icon="search"
            placeholder="Search..."
          />
          <br />
          <Dropdown
            placeholder="Labels"
            fluid
            multiple
            search
            selection
            onChange={handleLabelChange}
            options={[
              ...new Set(
                [open === 'Tasks' ? props.tasks : props.workflows]
                  .flat()
                  .map(wf => {
                    return wf.description
                      ? wf.description
                          .split('-')
                          .pop()
                          .replace(/\s/g, '')
                          .split(',')
                      : null;
                  })
                  .flat()
                  .filter(item => item !== null),
              ),
            ].map((label, i) => {
              return {key: i, text: label, value: label};
            })}
          />
          <small>
            Combine search and labels to find{' '}
            {open ? open.toLowerCase() : 'workflows'}
          </small>
        </div>
        <Divider horizontal section>
          {content.length} results
        </Divider>
        <div className="sidebar-content">{content}</div>
      </Sidebar>
    </div>
  );
};

export default Sidemenu;
