import { NodeModel } from "storm-react-diagrams";
import { CircleEndPortModel } from "./CircleEndPortModel";
import * as _ from "lodash";
import { DiagramEngine } from "storm-react-diagrams";

export class CircleEndNodeModel extends NodeModel {
  name: string;
  color: string;
  ports: { [s: string]: CircleEndPortModel };

  constructor(name: string = "Untitled", color: string = "rgb(0,192,255)") {
    super("end", "end");
    this.name = name;
    this.color = color;

    this.addPort(new CircleEndPortModel(true, "left"));
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
}
