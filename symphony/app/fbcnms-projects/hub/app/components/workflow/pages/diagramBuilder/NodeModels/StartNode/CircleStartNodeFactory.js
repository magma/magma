import * as SRD from "storm-react-diagrams";
import { CircleNodeStart } from "./CircleNodeStart";
import { CircleStartNodeModel } from "./CircleStartNodeModel";
import * as React from "react";

export class CircleStartNodeFactory extends SRD.AbstractNodeFactory {
  constructor() {
    super("start");
  }

  generateReactWidget(
    diagramEngine: SRD.DiagramEngine,
    node: SRD.NodeModel
  ): JSX.Element {
    return <CircleNodeStart node={node} />;
  }

  getNewInstance() {
    return new CircleStartNodeModel();
  }
}
