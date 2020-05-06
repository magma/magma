import * as React from "react";
import { PortWidget } from "@projectstorm/react-diagrams";

export class CircleNodeEnd extends React.Component {
  constructor(props) {
    super(props);
    this.state = {};
  }

  render() {
    return (
      <div className={"srd-circle-node"}>
        <svg width="60" height="60">
          <g>
            <circle cx="30" cy="30" r="30" fill="white" />
            <text x="17" y="35">
              End
            </text>
          </g>
        </svg>
        <div style={{ position: "absolute", zIndex: 10, left: -10, top: 21 }}>
          <PortWidget name="left" node={this.props.node} />
        </div>
      </div>
    );
  }
}
