import * as React from 'react';
import {NodeContextMenu, NodeMenuProvider} from '../ContextMenu';
import {PortWidget} from '@projectstorm/react-diagrams';

export class ForkNode extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      size: 80,
    };
  }

  render() {
    return (
      <div>
        <NodeMenuProvider node={this.props.node}>
          <svg
            width={this.state.size}
            height={this.state.size}
            dangerouslySetInnerHTML={{
              __html: `
          <g id="Layer_2">
            <polygon fill="${this.props.node.color}" points="30 65,65 65,65 15,30 15,15 40"/>
                <text x="32" y="45" fill="white" font-size="13px" >fork</text>
          </g>
        `,
            }}
          />
        </NodeMenuProvider>

        <div
          className="srd-node-glow"
          style={{
            position: 'absolute',
            zIndex: -1,
            left: 45,
            top: 40,
          }}
        />

        <div
          style={{
            position: 'absolute',
            zIndex: 10,
            left: 7,
            top: this.state.size / 2 - 12,
          }}>
          <PortWidget name="left" node={this.props.node} />
        </div>

        <div
          style={{
            position: 'absolute',
            zIndex: 10,
            left: this.state.size - 25,
            top: this.state.size / 2 - 12,
          }}>
          <PortWidget name="right" node={this.props.node} />
        </div>

        <NodeContextMenu
          node={this.props.node}
          diagramEngine={this.props.diagramEngine}
        />
      </div>
    );
  }
}
