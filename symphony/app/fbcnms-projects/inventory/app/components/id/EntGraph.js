/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {WithStyles} from '@material-ui/core';

import * as d3 from 'd3';
import CircularProgress from '@material-ui/core/CircularProgress';
import GraphVertexDetails from './GraphVertexDetails';
import React from 'react';
import SideBar from '@fbcnms/ui/components/layout/SideBar';
import classNames from 'classnames';
import nullthrows from '@fbcnms/util/nullthrows';
import {withStyles} from '@material-ui/core/styles';

const styles = theme => ({
  root: {
    display: 'flex',
    backgroundColor: 'white',
    padding: '16px',
    boxShadow: '0px 1px 4px 0px rgba(0,0,0,0.17)',
    borderRadius: '4px',
    height: '700px',
    alignItems: 'center',
    justifyContent: 'center',
  },
  svg: {
    width: '100%',
    height: '100%',
    opacity: 1,
  },
  link: {
    fill: 'none',
    stroke: theme.palette.grey[300],
    strokeWidth: 1,
  },
  node: {
    fill: theme.palette.common.white,
    stroke: theme.palette.blueGrayDark,
    strokeWidth: 0,
    '&:hover': {
      cursor: 'pointer',
    },
  },
  nodeText: {
    ...theme.typography.caption,
    fill: theme.palette.blueGrayDark,
    fontSize: TEXT_FONT_SIZE,
    cursor: 'pointer',
    pointerEvents: 'none',
    stroke: 'none',
    textAnchor: 'middle',
  },
  edgePath: {
    fillOpacity: 0,
    strokeOpacity: 0,
    pointerEvents: 'none',
  },
  edgeLabel: {
    fontSize: 12,
    fill: theme.palette.blueGrayDark,
    background: 'white',
    pointerEvents: 'none',
    textAnchor: 'middle',
  },
  arrowHead: {
    stroke: theme.palette.grey[300],
    fill: theme.palette.grey[300],
  },
});

const LINK_DISTANCE = 250;
const NODE_RADIUS = 36;
const TEXT_FONT_SIZE = 10;
const EPSILON = 0.01;
const ALPHA_TARGET = 0.3;

type Props = {
  rootNode: Object,
} & WithStyles<typeof styles>;

type State = {
  loading: boolean,
  selectedNodeId: ?string,
};

class EntGraph extends React.Component<Props, State> {
  _topologyContainer = null;

  constructor(props) {
    super(props);
    this._topologyContainer = React.createRef();

    this.state = {
      loading: false,
      selectedNodeId: null,
    };
  }

  componentDidMount() {
    const {classes, rootNode} = this.props;
    const container = nullthrows(this._topologyContainer?.current);

    const height = container.clientHeight;
    const width = container.clientWidth;
    const outgoingNodes = rootNode.edges
      .map(edge =>
        edge.ids.map(id => ({
          id,
          type: edge.type,
        })),
      )
      .filter(nodes => nodes.length > 0)
      .flat();
    const d3RootNode = {
      id: rootNode.id,
      type: rootNode.type,
    };

    const nodes = [d3RootNode, ...outgoingNodes];
    const links = outgoingNodes.map(outNode => ({
      source: d3RootNode.id,
      target: outNode.id,
      type: outNode.type,
    }));

    const svg = d3
      .select(nullthrows(this._topologyContainer).current)
      .append('svg')
      .attr('width', width)
      .attr('height', height)
      .attr('viewBox', [-width / 2, -height / 2, width, height]);

    svg
      .append('defs')
      .append('marker')
      .attr('id', 'arrowhead')
      .attr('viewBox', '0 0 10 10')
      .attr('refX', NODE_RADIUS + NODE_RADIUS / 2 + 3)
      .attr('refY', 5)
      .attr('orient', 'auto')
      .attr('markerWidth', 8)
      .attr('markerHeight', 8)
      .attr('xoverflow', 'visible')
      .append('path')
      .attr('d', 'M 0 0 L 10 5 L 0 10 Z')
      .attr('class', classes.arrowHead);

    const g = svg.append('g');

    // Create force simulation which will place the nodes and links on screen
    // with the correct distances between them
    const simulation = d3
      .forceSimulation(nodes)
      .alphaTarget(ALPHA_TARGET)
      .force(
        'link',
        d3
          .forceLink()
          .id(d => d.id)
          .distance(LINK_DISTANCE)
          .links(links)
          .strength(link => (link.source.type === link.target.type ? 0.5 : 1)),
      )
      .force('collide', d3.forceCollide().radius(NODE_RADIUS));

    // For each link create a line element
    const link = g
      .selectAll('line')
      .data(links)
      .enter()
      .append('line')
      .attr('class', classes.link)
      .attr('marker-end', 'url(#arrowhead)');

    const edgePaths = svg
      .selectAll(`.${classes.edgePath}`)
      .data(links)
      .enter()
      .append('path')
      .attr('class', classes.edgePath)
      .attr('id', (d, i) => 'edgepath' + i);

    const edgeLabels = svg
      .selectAll(`.${classes.edgeLabel}`)
      .data(links)
      .enter()
      .append('text')
      .attr('dy', 4)
      .attr('class', classes.edgeLabel)
      .attr('id', (d, i) => 'edgelabel' + i);

    edgeLabels
      .append('textPath')
      .attr('startOffset', '50%')
      .attr('href', (d, i) => `#edgepath${i}`)
      .text(d => d.type);

    // For each node create a group and enable dragging it
    const color = d3
      .scaleOrdinal()
      .domain([0, rootNode.edges.length])
      .range([
        '#77e6e6',
        '#f7923b',
        '#fbd872',
        '#e0b8fc',
        '#f58796',
        '#caeef9',
        '#d1e6b9',
        '#faad9b',
      ]);

    const node = g
      .selectAll(`.${classes.node}`)
      .data(nodes)
      .enter()
      .append('g')
      .attr('class', classes.node)
      .attr('width', NODE_RADIUS * 2)
      .attr('height', NODE_RADIUS * 2)
      .call(this._drag(simulation))
      .on('click', this.onNodeClicked);

    // Add a circle and a label on each node
    node
      .append('circle')
      .attr('r', NODE_RADIUS)
      .attr('fill', d =>
        color(rootNode.edges.findIndex(n => n.type === d.type)),
      );

    node
      .append('text')
      .text(d => d.id)
      .attr('dy', TEXT_FONT_SIZE / 3)
      .attr('class', classes.nodeText);

    const positionNodes = () => {
      node.attr('transform', d => `translate(${d.x} ${d.y})`);

      link
        .attr('x1', d => d.source.x)
        .attr('y1', d => d.source.y)
        .attr('x2', d => d.target.x)
        .attr('y2', d => d.target.y);

      edgePaths.attr(
        'd',
        d =>
          'M ' +
          d.source.x +
          ' ' +
          d.source.y +
          ' L ' +
          d.target.x +
          ' ' +
          d.target.y,
      );

      edgeLabels.attr('transform', (d, i, nodes) => {
        if (d.target.x < d.source.x) {
          // eslint-disable-next-line no-warning-comments
          // $FlowFixMe - getBBox exists in this context
          const bbox = nodes[i].getBBox();
          const rx = bbox.x + bbox.width / 2;
          const ry = bbox.y + bbox.height / 2;
          return 'rotate(180 ' + rx + ' ' + ry + ')';
        } else {
          return 'rotate(0)';
        }
      });
    };

    simulation.on('tick', () => {
      if (simulation.alpha() - simulation.alphaTarget() < EPSILON) {
        return;
      }

      positionNodes();
    });
  }

  onNodeClicked = node => {
    const {classes} = this.props;
    this.setState({selectedNodeId: node.id});
    d3.selectAll(`.${classes.node}`)
      .style('stroke-width', d => (d.id === node.id ? 2 : 0))
      .select('text')
      .style('font-weight', d => (d.id === node.id ? 600 : 400));
  };

  _drag = simulation => {
    const dragstarted = d => {
      if (!d3.event.active) {
        simulation.alpha(0.5).alphaTarget(ALPHA_TARGET).restart();
      }
      d.fx = d.x;
      d.fy = d.y;
    };

    const dragged = d => {
      d.fx = d3.event.x;
      d.fy = d3.event.y;
    };

    const dragended = d => {
      d.fx = null;
      d.fy = null;
    };

    return d3
      .drag()
      .on('start', dragstarted)
      .on('drag', dragged)
      .on('end', dragended);
  };

  render() {
    const {classes} = this.props;
    const {loading, selectedNodeId} = this.state;
    return (
      <div className={classes.root}>
        {loading && <CircularProgress />}
        <div
          className={classNames({
            [classes.svg]: !loading,
          })}
          ref={this._topologyContainer}
        />
        <SideBar
          top={64}
          isShown={selectedNodeId !== null}
          onClose={() => this.setState({selectedNodeId: null})}>
          <GraphVertexDetails vertexId={selectedNodeId ?? ''} />
        </SideBar>
      </div>
    );
  }
}

export default withStyles(styles)(EntGraph);
