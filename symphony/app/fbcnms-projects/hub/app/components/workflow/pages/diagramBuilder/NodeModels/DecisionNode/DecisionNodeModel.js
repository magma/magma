import { NodeModel } from "@projectstorm/react-diagrams";
import * as _ from "lodash";
import { DefaultPortModel } from "@projectstorm/react-diagrams";
import { DiagramEngine } from "@projectstorm/react-diagrams";
import { DecisionNodePortModel } from "./DecisionNodePortModel";

export class DecisionNodeModel extends NodeModel {
  name: string;
  color: string;
  inputs: {};

  constructor(
    name: string = "Untitled",
    color: string = "rgb(0,192,255)",
    inputs: {}
  ) {
    super("decision");
    this.name = name;
    this.color = color;
    super.extras = { inputs: inputs };

    this.addPort(new DecisionNodePortModel(true, "left", "inputPort"));
    this.addPort(new DecisionNodePortModel(false, "right", "neutralPort"));
    this.addPort(new DecisionNodePortModel(false, "bottom", "failPort"));
  }

  deSerialize(object, engine: DiagramEngine) {
    super.deSerialize(object, engine);
    this.name = object.name;
    this.color = object.color;
  }

  serialize() {
    return _.merge(super.serialize(), {
      name: this.name,
      color: this.color
    });
  }

  getInputs() {
    return this.inputs;
  }
}
