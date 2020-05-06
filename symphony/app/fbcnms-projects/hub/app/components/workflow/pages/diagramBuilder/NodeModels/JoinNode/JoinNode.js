import * as React from "react";
import { PortWidget } from "@projectstorm/react-diagrams";
import { NodeContextMenu, NodeMenuProvider } from "../ContextMenu";

export class JoinNode extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      size: 80
    };
  }

  render() {
    return (
      <div
        className={"join-node"}
        style={{
          position: "relative",
          width: this.state.size,
          height: this.state.size
        }}
      >
        <NodeMenuProvider node={this.props.node}>
        <svg
          width={this.state.size}
          height={this.state.size}
          dangerouslySetInnerHTML={{
            __html: `
          <g id="Layer_1">
          </g>
          <g id="Layer_2">
           <polygon fill="${this.props.node.color}" points="50 15,15 15,15 65,50 65,65 40"/>
                <text x="26" y="45" fill="white" font-size="13px" >join</text>
          </g>
        `
          }}
        />
        </NodeMenuProvider>

        <div className="srd-node-glow"
             style={{
               position: "absolute",
               zIndex: -1,
               left: 35,
               top: 40
             }}
        />

        <div
          style={{
            position: "absolute",
            zIndex: 10,
            left: 5,
            top: this.state.size / 2 - 12
          }}
        >
          <PortWidget name="left" node={this.props.node} />
        </div>

        <div
          style={{
            position: "absolute",
            zIndex: 10,
            left: this.state.size - 28,
            top: this.state.size / 2 - 12
          }}
        >
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
