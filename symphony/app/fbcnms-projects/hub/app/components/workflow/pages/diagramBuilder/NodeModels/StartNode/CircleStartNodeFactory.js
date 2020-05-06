import * as SRD from "@projectstorm/react-diagrams";
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
    return <CircleNodeStart node={node} diagramEngine={diagramEngine} />;
  }

  getNewInstance() {
    return new CircleStartNodeModel();
  }
}
