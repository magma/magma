/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ForceNetworkTopology_topology} from './__generated__/ForceNetworkTopology_topology.graphql';
import type {WithStyles} from '@material-ui/core';

import * as React from 'react';
import * as d3 from 'd3';
import CircularProgress from '@material-ui/core/CircularProgress';
import classNames from 'classnames';
import nullthrows from '@fbcnms/util/nullthrows';
import symphony from '@fbcnms/ui/theme/symphony';
import {createFragmentContainer, graphql} from 'react-relay';
import {renderToStaticMarkup} from 'react-dom/server';
import {withStyles} from '@material-ui/core/styles';

const CHARGE_STRENGTH = -400;
const COLLIDE_RADIUS = 1;
const LINK_DISTANCE = 200;
const NODE_RADIUS = 20;
const PADDING = 24;

const styles = theme => ({
  root: {
    position: 'absolute',
    left: 0,
    right: 0,
    top: 0,
    bottom: 0,
    opacity: 0,
  },
  rootShown: {
    opacity: 1,
  },
  spinner: {
    position: 'absolute',
    left: 'calc(50% - 20px) ',
    top: 'calc(50% - 20px)',
  },
  link: {
    fill: 'none',
    stroke: symphony.palette.D300,
    strokeWidth: 1,
  },
  node: {
    fill: theme.palette.common.white,
    stroke: theme.palette.grey[300],
    strokeWidth: 1,
  },
});

type Props = {
  topology: ForceNetworkTopology_topology,
  className?: string,
  renderNode: (id: string) => React.Node,
} & WithStyles<typeof styles>;

type State = {
  loading: boolean,
};

class ForceNetworkTopology extends React.Component<Props, State> {
  _topologyContainer = null;

  constructor(props) {
    super(props);
    this._topologyContainer = React.createRef();

    this.state = {
      loading: true,
    };
  }

  componentDidMount() {
    this.calculateGraph();
  }

  componentDidUpdate(prevProps) {
    if (
      this.props.topology === prevProps.topology &&
      this.props.renderNode === prevProps.renderNode
    ) {
      return;
    }
    this.calculateGraph();
  }

  calculateGraph() {
    const {classes, topology, renderNode} = this.props;

    const container = nullthrows(this._topologyContainer?.current);
    const height = container.clientHeight - PADDING;
    const width = container.clientWidth - PADDING;
    const nodes = topology.nodes.map(node => ({
      id: node.id,
    }));
    const links = topology.links.map(l => ({
      source: l.source.id,
      target: l.target.id,
    }));

    d3.select(nullthrows(this._topologyContainer).current)
      .selectAll('svg')
      .remove();

    // Create SVG
    const svg = d3
      .select(nullthrows(this._topologyContainer).current)
      .append('svg')
      .attr('width', width)
      .attr('height', height)
      .attr('viewBox', [-width / 2, -height / 2, width, height]);

    const g = svg.append('g');

    // Create force simulation which will place the nodes and links on screen
    // with the correct distances between them
    const simulation = d3
      .forceSimulation(nodes)
      .force(
        'link',
        d3
          .forceLink()
          .id(d => d.id)
          .distance(LINK_DISTANCE)
          .links(links),
      )
      .force('x', d3.forceX())
      .force('y', d3.forceY())
      .force('charge', d3.forceManyBody().strength(CHARGE_STRENGTH))
      .force('collide', d3.forceCollide().radius(COLLIDE_RADIUS))
      .alphaDecay(0.5);

    // For each link create a line element
    const link = g
      .selectAll('line')
      .data(links)
      .enter()
      .append('line')
      .attr('class', classes.link);

    // For each node create a group and enable dragging it
    const node = g
      .selectAll()
      .data(nodes)
      .enter()
      .append('g')
      .attr('class', classes.node)
      .attr('width', NODE_RADIUS * 2)
      .attr('height', NODE_RADIUS * 2)
      .call(this._drag(simulation));

    const newNode = node.html(d => renderToStaticMarkup(renderNode(d.id)));
    const text = newNode.select('text');
    const textBoxes = text.nodes().map(node => node.getBBox());
    newNode
      .select('#textBox')
      .select('rect')
      .attr(
        'transform',
        (d, i) => `translate(-${textBoxes[i].width / 2 + 16} -16)`,
      )
      .data(textBoxes)
      .attr('width', (d, i) => textBoxes[i].width + 30)
      .attr('height', (d, i) => textBoxes[i].height + 12);

    const positionNodes = () => {
      node.attr('transform', d => `translate(${d.x} ${d.y})`);

      link
        .attr('x1', d => d.source.x)
        .attr('y1', d => d.source.y)
        .attr('x2', d => d.target.x)
        .attr('y2', d => d.target.y);
    };

    simulation.on('tick', () => {
      node
        .attr(
          'cx',
          d =>
            (d.x =
              Math.max(
                NODE_RADIUS,
                Math.min(width - NODE_RADIUS, d.x + width / 2),
              ) -
              width / 2),
        )
        .attr(
          'cy',
          d =>
            (d.y =
              Math.max(
                NODE_RADIUS,
                Math.min(height - NODE_RADIUS, d.y + height / 2),
              ) -
              height / 2),
        );

      link
        .attr('x1', d => d.source.x)
        .attr('y1', d => d.source.y)
        .attr('x2', d => d.target.x)
        .attr('y2', d => d.target.y);
      positionNodes();
    });

    simulation.on('end', () => {
      positionNodes();
      if (this.state.loading) {
        this.setState({loading: false});
      }
    });
  }

  _drag = simulation => {
    const dragstarted = d => {
      if (!d3.event.active) {
        simulation.alphaTarget(0.3).restart();
      }
      d.fx = d.x;
      d.fy = d.y;
    };

    const dragged = d => {
      d.fx = d3.event.x;
      d.fy = d3.event.y;
    };

    const dragended = d => {
      if (!d3.event.active) {
        simulation.alphaTarget(0);
      }
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
    return (
      <>
        {this.state.loading && <CircularProgress className={classes.spinner} />}
        <div
          className={classNames(
            {
              [classes.root]: true,
              [classes.rootShown]: !this.state.loading,
            },
            classes.className,
          )}
          ref={this._topologyContainer}
        />
      </>
    );
  }
}

export default withStyles(styles)(
  createFragmentContainer(ForceNetworkTopology, {
    topology: graphql`
      fragment ForceNetworkTopology_topology on NetworkTopology {
        nodes {
          id
        }
        links {
          source {
            id
          }
          target {
            id
          }
        }
      }
    `,
  }),
);
